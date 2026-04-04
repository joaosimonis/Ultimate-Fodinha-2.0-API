package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

type Ultimate struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type Player struct {
	Name      string
	AvatarKey string
	Title     string
	Lore      string
	Ultimate  Ultimate
}

var players = []Player{
	{
		Name: "Ângela", AvatarKey: "angela", Title: "Guerreira Lendária",
		Lore: "Ângela combina força e estratégia com habilidade extraordinária. Sua presença inspira aliados e intimida inimigos.\n\nSeu nome ecoa em histórias de batalhas lendárias.",
		Ultimate: Ultimate{Name: "Tempestade de Lâminas", Description: "Ângela executa uma sequência devastadora de ataques rápidos."},
	},
	{
		Name: "Bruno", AvatarKey: "bruno", Title: "Mago de Gelo",
		Lore: "Bruno domina os segredos das geadas eternas, invocando tempestades congelantes e lanças de cristal puro. Onde ele caminha, o ar se torna frio e o solo se cobre de gelo.\n\nSua mente é tão calma quanto o inverno profundo, e sua magia reflete esse controle absoluto sobre o elemento da geada.",
		Ultimate: Ultimate{Name: "Prisão Glacial", Description: "Uma tempestade congelante envolve o campo de batalha, prendendo inimigos em blocos de gelo."},
	},
	{
		Name: "Duda", AvatarKey: "duda", Title: "Tecelã dos Sonhos",
		Lore: "Duda caminha entre sonhos e memórias.\n\nEla transforma pesadelos em esperança.",
		Ultimate: Ultimate{Name: "Reino dos Sonhos", Description: "Inimigos ficam desorientados enquanto aliados recebem regeneração."},
	},
	{
		Name: "Fabri", AvatarKey: "fabri", Title: "Maestro da Tempestade Sonora",
		Lore: "Fabri controla ondas sonoras como um verdadeiro maestro do caos. Cada movimento cria vibrações poderosas.\n\nBatalhas se tornam concertos de energia sonora sob seu comando.",
		Ultimate: Ultimate{Name: "Sinfonia do Trovão", Description: "Uma onda sonora massiva explode causando dano em área."},
	},
	{
		Name: "Fran", AvatarKey: "fran", Title: "Mestra dos Conhecimentos Antigos",
		Lore: "Fran é guardiã de saberes esquecidos e guia novas gerações.\n\nSeu conhecimento ilumina caminhos.",
		Ultimate: Ultimate{Name: "Biblioteca Arcana", Description: "Runas antigas surgem fortalecendo aliados e enfraquecendo inimigos."},
	},
	{
		Name: "Geladeira", AvatarKey: "geladeira", Title: "Engenheiro do Gelo",
		Lore: "Geladeira combina ciência e magia para criar estruturas congeladas complexas. Ele molda o gelo como um arquiteto molda pedra.\n\nSuas invenções transformam o campo de batalha.",
		Ultimate: Ultimate{Name: "Fortaleza de Cristal", Description: "Uma muralha gigante de gelo surge protegendo aliados."},
	},
	{
		Name: "Grade", AvatarKey: "grade", Title: "Domador de Dragão",
		Lore: "Grade é conhecido por possuir uma ligação rara e poderosa com as criaturas mais temidas dos céus: os dragões.\n\nEnquanto muitos guerreiros tentaram enfrentá-los, ele aprendeu a compreendê-los e conquistar sua confiança, tornando-se um verdadeiro domador dessas feras lendárias.",
		Ultimate: Ultimate{Name: "Fúria do Dragão Ancestral", Description: "Grade invoca seu dragão para sobrevoar o campo de batalha, liberando uma poderosa rajada de fogo que causa dano massivo em área."},
	},
	{
		Name: "Guto", AvatarKey: "guto", Title: "Druida",
		Lore: "Guto é um guardião da natureza que vive em harmonia com as forças da terra. Árvores, animais e ventos respondem ao seu chamado.\n\nSeu vínculo com a natureza é tão profundo que ele pode invocar raízes gigantes e criaturas da floresta para lutar ao seu lado.",
		Ultimate: Ultimate{Name: "Ira da Floresta", Description: "Raízes gigantes emergem do solo prendendo inimigos e causando dano contínuo."},
	},
	{
		Name: "Habekost", AvatarKey: "habekost", Title: "Ferreiro dos Deuses",
		Lore: "Habekost forja armas capazes de rivalizar com relíquias divinas.\n\nCada criação carrega poder ancestral.",
		Ultimate: Ultimate{Name: "Arsenal Divino", Description: "Armas lendárias surgem fortalecendo todos os aliados."},
	},
	{
		Name: "João", AvatarKey: "joao", Title: "Deus da Luz Celestial",
		Lore: "João canaliza o poder da luz divina. Sua presença ilumina até os lugares mais sombrios.\n\nEle representa esperança e justiça.",
		Ultimate: Ultimate{Name: "Julgamento Celestial", Description: "Feixes de luz descem do céu atingindo todos os inimigos."},
	},
	{
		Name: "Kevin", AvatarKey: "kevin", Title: "Capitão da Rocinha",
		Lore: "Kevin é um líder nato que inspira coragem e união entre seus aliados. Empunhando seu escudo lendário, ele avança sempre à frente da batalha, protegendo os seus e enfrentando qualquer ameaça com determinação inabalável.\n\nSeu espírito heroico transforma qualquer confronto em um símbolo de resistência e esperança.",
		Ultimate: Ultimate{Name: "Investida do Capitão", Description: "Kevin arremessa seu escudo lendário que ricocheteia entre os inimigos causando dano, retornando em seguida para conceder proteção temporária aos aliados."},
	},
	{
		Name: "Knak", AvatarKey: "knak", Title: "Mensageira Celestial",
		Lore: "Knak leva a palavra dos deuses aos mortais.\n\nSua presença inspira esperança.",
		Ultimate: Ultimate{Name: "Chamado Divino", Description: "Aliados recebem velocidade e cura instantânea."},
	},
	{
		Name: "Kreutzer", AvatarKey: "kreutzer", Title: "Mestre das Marés",
		Lore: "Kreutzer controla correntes marítimas e tempestades oceânicas.\n\nO mar responde ao seu chamado.",
		Ultimate: Ultimate{Name: "Tsunami Ancestral", Description: "Uma onda gigantesca atravessa o campo causando dano massivo."},
	},
	{
		Name: "Laura", AvatarKey: "laura", Title: "Feiticeira de Fogo",
		Lore: "Laura é uma manipuladora das chamas elementais, capaz de invocar labaredas que dançam ao seu redor como criaturas vivas. Seu poder nasce da intensidade de sua própria energia interior.\n\nQuando entra em combate, o campo de batalha se ilumina com o brilho das chamas que obedecem ao seu comando.",
		Ultimate: Ultimate{Name: "Inferno Carmesim", Description: "Uma explosão de fogo cobre o campo de batalha causando dano massivo em área."},
	},
	{
		Name: "Laura S.", AvatarKey: "lauraleoa", Title: "Feiticeira Arcana Sombria",
		Lore: "Laura S. domina energias arcanas ocultas que poucos ousam estudar.\n\nSua magia mistura beleza e perigo.",
		Ultimate: Ultimate{Name: "Abismo Arcano", Description: "Uma explosão de energia sombria consome os inimigos próximos."},
	},
	{
		Name: "Leo", AvatarKey: "leo", Title: "Alquimista Culinário",
		Lore: "Léo é um alquimista incomum que domina a arte de transformar ingredientes raros em criações mágicas. Em sua cozinha-fornalha, especiarias arcanas, ervas místicas e essências elementais se combinam para produzir pratos capazes de fortalecer aliados e restaurar energias perdidas.\n\nSuas receitas são tão precisas quanto poderosas, cada prato carregando propriedades únicas.",
		Ultimate: Ultimate{Name: "Banquete dos Heróis", Description: "Léo prepara um prato lendário que cura todos os aliados e concede força e resistência temporária para a equipe."},
	},
	{
		Name: "Lily", AvatarKey: "lily", Title: "Arquiteta de Atena",
		Lore: "Discípula da deusa da sabedoria, Lily cria estruturas mágicas e fortalezas estratégicas.\n\nSeu conhecimento transforma o campo de batalha.",
		Ultimate: Ultimate{Name: "Fortaleza de Atena", Description: "Uma cidadela mágica surge protegendo aliados."},
	},
	{
		Name: "Luís S.", AvatarKey: "luis", Title: "Caçador de Monstros Lendários",
		Lore: "Luís S. dedicou sua vida a caçar criaturas míticas. Seu conhecimento sobre monstros é incomparável.\n\nEle enfrenta perigos que poucos sobreviveriam.",
		Ultimate: Ultimate{Name: "Golpe do Caçador Supremo", Description: "Um ataque devastador que causa dano extra em criaturas poderosas."},
	},
	{
		Name: "Messi", AvatarKey: "messi", Title: "Elfo Arcano",
		Lore: "Messi é um antigo elfo que domina os fluxos da magia arcana com precisão absoluta. Desde jovem demonstrou um talento raro para manipular energias invisíveis, tornando-se um dos magos mais refinados de sua era.\n\nSua magia é elegante e rápida, quase impossível de prever.",
		Ultimate: Ultimate{Name: "Tempestade Arcana", Description: "Messi libera uma explosão massiva de energia arcana que atinge todos os inimigos próximos."},
	},
	{
		Name: "Murilo", AvatarKey: "murilo", Title: "Guerreiro Lendário",
		Lore: "Murilo é um herói reconhecido por sua coragem e habilidade em batalha. Sua armadura carrega marcas de inúmeras guerras vencidas.\n\nEle lidera seus aliados com determinação inabalável.",
		Ultimate: Ultimate{Name: "Investida do Herói", Description: "Murilo avança devastando inimigos em seu caminho."},
	},
	{
		Name: "Pelotas", AvatarKey: "pelotas", Title: "Guardião das Sombras",
		Lore: "Pelotas é uma figura enigmática que caminha entre luz e escuridão, dominando as sombras como poucos no mundo.\n\nUsando a própria escuridão como arma e proteção, ele manipula sombras para confundir inimigos e desaparecer no campo de batalha.",
		Ultimate: Ultimate{Name: "Manto da Noite", Description: "Pelotas envolve o campo em sombras profundas, reduzindo a visão dos inimigos enquanto se move livremente entre as trevas atacando múltiplos alvos."},
	},
	{
		Name: "Rogus", AvatarKey: "rogus", Title: "Berserker",
		Lore: "Rogus entra em batalha com uma fúria quase sobrenatural. Quanto mais luta, mais forte se torna.\n\nSeu estilo brutal de combate é temido em todos os campos de batalha.",
		Ultimate: Ultimate{Name: "Fúria Incontrolável", Description: "Rogus entra em estado de fúria causando dano massivo e ganhando resistência."},
	},
	{
		Name: "Thaís", AvatarKey: "thais", Title: "Sacerdotista",
		Lore: "Thaís canaliza energias divinas para proteger e curar aliados. Sua fé é sua maior força.\n\nSua presença traz esperança mesmo nas batalhas mais difíceis.",
		Ultimate: Ultimate{Name: "Bênção Sagrada", Description: "Todos os aliados são curados e recebem proteção divina."},
	},
	{
		Name: "Thiago", AvatarKey: "thiago", Title: "Monge",
		Lore: "Thiago alcançou a maestria do corpo e da mente através de anos de disciplina. Seus movimentos são precisos e harmoniosos.\n\nCada golpe seu é executado com foco absoluto, transformando combate em arte marcial perfeita.",
		Ultimate: Ultimate{Name: "Golpe do Dragão Interior", Description: "Um poderoso ataque de energia espiritual atinge todos os inimigos à frente."},
	},
	{
		Name: "Vogt", AvatarKey: "vogt", Title: "Guardião do Tempo",
		Lore: "Vogt vigia o fluxo do tempo e protege o destino.\n\nEle vê ecos do passado e vislumbres do futuro.",
		Ultimate: Ultimate{Name: "Ruptura Temporal", Description: "O tempo desacelera para inimigos enquanto aliados agem mais rápido."},
	},
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("sem .env, usando variáveis de ambiente")
	}

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_SSLMODE"),
	)

	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		log.Fatalf("erro ao conectar: %v", err)
	}
	defer pool.Close()

	for _, p := range players {
		ultimateJSON, err := json.Marshal(p.Ultimate)
		if err != nil {
			log.Fatalf("erro ao serializar ultimate de %s: %v", p.Name, err)
		}

		_, err = pool.Exec(context.Background(), `
			INSERT INTO players (name, avatar_key, title, lore, ultimate)
			VALUES ($1, $2, $3, $4, $5)
			ON CONFLICT (name) DO NOTHING
		`, p.Name, p.AvatarKey, p.Title, p.Lore, ultimateJSON)
		if err != nil {
			log.Printf("❌ erro ao inserir %s: %v", p.Name, err)
			continue
		}
		log.Printf("✅ %s inserido", p.Name)
	}

	log.Println("🎉 seed concluído!")
}
