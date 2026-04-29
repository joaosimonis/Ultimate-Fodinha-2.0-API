package players

import (
	"context" // Para saber o contexto da requisição até as querys
	"encoding/json" // Necessário para converter de JSONB para struct e vice-versa

	"github.com/jackc/pgx/v5" // Biblioteca PostgreSQL para GO, usada aqui para lidar com erros específicos como pgx.ErrNoRows
	"github.com/jackc/pgx/v5/pgxpool" // Pool de conexões com o PostgreSQL
)

/*
	É a camada que conversa diretamente com o PostegreSQL, usando o pgxpool para fazer as queries.
	É a camada mais baixa, ou seja, é a que tem mais acoplamento com a implementação do banco.
	Por isso, ela é a única que sabe como fazer as queries, como usar o pgxpool, etc.
	As outras camadas (service e handler) não sabem nada disso, elas só sabem que o repository tem os métodos que elas precisam (FindAll, FindByID, Create, Delete, ExistsByName).
*/

// 1. Struct do repository, que guarda o acesso ao banco (pool de conexões)
type Repository struct {
	pool *pgxpool.Pool
}

// 2. Construtor do repository, que recebe o pool de conexões e devolve um ponteiro para Repository
func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

/*
	3. Funções do repository, que fazem as queries no banco usando o pool de conexões
	   Petencem ao tipo Repository (receiver) e implementam os métodos da interface RepositoryReader, que é a interface que o service conhece.
	   Elas recebem um contexto (context.Context) para respeitar o ciclo de vida da requisição
	   Elas usam o pgxpool para fazer as queries, e retornam os resultados ou erros para o service
*/

func (r *Repository) FindAll(ctx context.Context) ([]Player, error) {
	// 1. Executa SELECT no banco para buscar todos os jogadores, ordenando por nome (guardado na variável rows)
	rows, err := r.pool.Query(ctx, `
		SELECT id, name, avatar_key, title, lore, ultimate, created_at
		FROM players
		ORDER BY name ASC
	`)
	// 2. Se der erro na query, devolve o erro para o service tratar (pode ser erro de conexão, erro de sintaxe, etc)
	if err != nil {
		return nil, err
	}
	// 3. Garante que as linhas sejam fechadas depois de usar, para evitar vazamento de conexões
	defer rows.Close()

	// 4. Cria um slice nil de Player, para guardar os jogadores que vierem do bancom
	var players []Player
	// 5. Percorre as linhas com rows.Next()
	for rows.Next() {
		var p Player // Variável para guardar os dados de cada jogador
		var ultimateJSON []byte // Variável para guardar o campo ultimate, que é JSONB no banco e precisa ser convertido para a struct Ultimate
		/* 6. Faz o scan dos campos da linha para a struct Player e ultimateJSON
			  "&"" é importante porque o Scan precisa modificar essas variáveis. Ele recebe o endereço delas
			  A ordem precisa ser a mesma da query SELECT
		*/
		if err := rows.Scan(&p.ID, &p.Name, &p.AvatarKey, &p.Title, &p.Lore, &ultimateJSON, &p.CreatedAt); err != nil {
			return nil, err
		}
		// 7. Converte o ultimateJSON (que é um array de bytes) para a struct Ultimate, usando json.Unmarshal
		if ultimateJSON != nil {
			// Se der erro, devolve o erro para o service tratar
			if err := json.Unmarshal(ultimateJSON, &p.Ultimate); err != nil {
				return nil, err
			}
		}
		// 8. Adiciona o jogador ao slice de players
		players = append(players, p)
	}
	// 9. Devolve o slice de players para o service
	return players, nil
}

func (r *Repository) FindByID(ctx context.Context, id string) (Player, error) {
	var p Player
	var ultimateJSON []byte

	err := r.pool.QueryRow(ctx, `
		SELECT id, name, avatar_key, title, lore, ultimate, created_at
		FROM players
		WHERE id = $1
	`, id).Scan(&p.ID, &p.Name, &p.AvatarKey, &p.Title, &p.Lore, &ultimateJSON, &p.CreatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return Player{}, ErrNotFound
		}
		return Player{}, err
	}

	if ultimateJSON != nil {
		if err := json.Unmarshal(ultimateJSON, &p.Ultimate); err != nil {
			return Player{}, err
		}
	}

	return p, nil
}

func (r *Repository) Create(ctx context.Context, req CreatePlayerRequest) (Player, error) {
	ultimateJSON, err := json.Marshal(req.Ultimate)
	if err != nil {
		return Player{}, err
	}

	var p Player
	var ultimateBytes []byte
	err = r.pool.QueryRow(ctx, `
		INSERT INTO players (name, avatar_key, title, lore, ultimate)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, name, avatar_key, title, lore, ultimate, created_at
	`, req.Name, req.AvatarKey, req.Title, req.Lore, ultimateJSON).
		Scan(&p.ID, &p.Name, &p.AvatarKey, &p.Title, &p.Lore, &ultimateBytes, &p.CreatedAt)
	if err != nil {
		return Player{}, err
	}

	if ultimateBytes != nil {
		if err := json.Unmarshal(ultimateBytes, &p.Ultimate); err != nil {
			return Player{}, err
		}
	}

	return p, nil
}

func (r *Repository) Delete(ctx context.Context, id string) error {
	result, err := r.pool.Exec(ctx, `
		DELETE FROM players WHERE id = $1
	`, id)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *Repository) ExistsByName(ctx context.Context, name string) (bool, error) {
	var exists bool
	err := r.pool.QueryRow(ctx, `
		SELECT EXISTS(SELECT 1 FROM players WHERE LOWER(name) = LOWER($1))
	`, name).Scan(&exists)
	return exists, err
}
