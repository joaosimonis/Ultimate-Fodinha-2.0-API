package players

import "context" // Pacote padrão do GO para lidar com contextos, como o contexto de uma requisição HTTP (cancelamento, timeout, etc)

/*
	Service é a camada de lógica de negócio, onde ficam as regras da minha aplicação
	Fica entre o handler e o repository
*/

// 1. Interface que define os métodos que o service precisa do repository.
type RepositoryReader interface {
	FindAll(ctx context.Context) ([]Player, error) // FindAll é o método que busca todos os jogadores no banco
	FindByID(ctx context.Context, id string) (Player, error) // FindByID é o método que busca um jogador pelo ID no banco
	Create(ctx context.Context, req CreatePlayerRequest) (Player, error) // Create é o método que cria um jogador no banco
	Delete(ctx context.Context, id string) error // Delete é o método que deleta um jogador pelo ID no banco
	ExistsByName(ctx context.Context, name string) (bool, error) // ExistsByName é o método que verifica se já existe um jogador com o mesmo nome no banco
}

/* 
	2. Struct do service, que tem como dependência a interface RepositoryReader
	   Com isso, o service sabe buscar, criar, deletar e verificar nome
	   Mas não sabe como isso é feito, ou seja, não tem acoplamento com a implementação do repository
*/
type Service struct {
	repo RepositoryReader
}

/* 
	3. Construtor, que vai receber um repository e devolver um ponteiro para Service
	   Ele é chamado em routes.go, onde o repository é criado e passado para o service
	   Isso se chama INJEÇÃO DE DEPENDÊNCIA
	   Ou seja, em vez do service criar o repository, ele recebe o repository pronto
*/
func NewService(repo RepositoryReader) *Service {
	return &Service{repo: repo}
}

/* 
	4. Essas são minhas funções, que são chamadas pelo handler. 
	   Elas usam os métodos do repository para fazer o que precisam
	   Podem ter regras de negócio, como verificar se o nome já existe antes de criar um jogador.
*/

/* 
	Podemos ver um padrão de como funciona as funções do service:
	- Pertence ao tipo Service (receiver)
	- Recebe um contexto (context.Context) como primeiro parâmetro, para lidar com cancelamento, timeout, etc
	- O contexto é importante para a query respeitar o ciclo de vida da requisição!
	- Recebe os parâmetros necessários para a operação (como id, CreatePlayerRequest, etc)
	- Retorna o resultado da operação (como Player, []Player, ou apenas error)
*/
func (s *Service) ListPlayers(ctx context.Context) ([]Player, error) {
	playersList, err := s.repo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	// Garante que envie como um array vazio, e não null, caso não haja jogadores
	if playersList == nil {
		return []Player{}, nil
	}

	return playersList, nil
}

func (s *Service) GetPlayerByID(ctx context.Context, id string) (Player, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *Service) CreatePlayer(ctx context.Context, req CreatePlayerRequest) (Player, error) {
	exists, err := s.repo.ExistsByName(ctx, req.Name)
	// Esse tipo de erro pode ser que o banco caiu, query falhou, conexão falhou, etc...
	if err != nil {
		return Player{}, err
	}

	// Esse erro é específico de que o nome já existe
	if exists {
		return Player{}, ErrNameTaken
	}

	return s.repo.Create(ctx, req)
}

func (s *Service) DeletePlayer(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

