package main

import (
	"fmt"
	"image/color"
	"log"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type GameState int

const ( //estados do jogo (pausado, game over, etc)
	StateTitle    GameState = iota // 0
	StatePlaying                   // 1
	StatePaused                    // 2
	StateGameOver                  // 3
)
const (
	screenWidth  = 320
	screenHeight = 240
	gridSize     = 10
)

// Point representa uma coordenada (x, y) na tela
type Point struct {
	x, y int
}

// Game armazena o estado do jogo
type Game struct {
	state       GameState
	snake       []Point
	dir         Point
	food        Point
	ticks       int
	score       int
	speed       int
	lastpressed Point
	inputFila   []Point
}

func (g *Game) prafila(p Point) {
	if len(g.inputFila) < 2 {
		g.inputFila = append(g.inputFila, p)
	}
}

// Update contém a lógica do jogo (roda 60 vezes por segundo)
func (g *Game) Update() error {
	// Update contém a lógica do jogo (roda 60 vezes por segundo)
	switch g.state {
	case StateTitle:
		if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
			g.state = StatePlaying
			g.reset()
		}
	case StatePaused:
		if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
			g.state = StatePlaying
		}
	case StateGameOver:
		if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
			g.reset()
			g.state = StateTitle
		}
	case StatePlaying:
		if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
			g.state = StatePaused
		}

		// 1. Controle de direção
		if ebiten.IsKeyPressed(ebiten.KeyW) {
			g.prafila(Point{0, -1})
		} else if ebiten.IsKeyPressed(ebiten.KeyS) {
			g.prafila(Point{0, 1})
		} else if ebiten.IsKeyPressed(ebiten.KeyA) {
			g.prafila(Point{-1, 0})
		} else if ebiten.IsKeyPressed(ebiten.KeyD) {
			g.prafila(Point{1, 0})
		}

		// 2. Controle de velocidade (a cobra move a cada x frames baseada na speed)
		g.ticks++
		if g.ticks < g.speed {
			return nil
		}
		g.ticks = 0

		if len(g.inputFila) > 0 {
			novadirec := g.inputFila[0]
			g.inputFila = g.inputFila[1:]
			if (novadirec.x != 0 && g.lastpressed.x == 0) || (novadirec.y != 0 && g.lastpressed.y == 0) {
				g.dir = novadirec
			}
		}
		g.lastpressed = g.dir

		// 3. Calcular a nova posição da cabeça da cobra
		head := g.snake[0]
		newHead := Point{head.x + g.dir.x, head.y + g.dir.y}

		// 4. Checar colisão com as paredes (Game Over)
		if newHead.x < 0 || newHead.y < 0 || newHead.x >= screenWidth/gridSize || newHead.y >= screenHeight/gridSize {
			g.state = StateGameOver
			return nil
		}

		// 5. Checar colisão com o próprio corpo
		for _, p := range g.snake {
			if p == newHead {
				g.state = StateGameOver
				return nil
			}
		}

		// 6. Lógica de Comer
		if newHead == g.food {
			g.score += 10
			if g.score%50 == 0 && g.speed > 3 {
				g.speed -= 2
			}
			// OTIMIZAÇÃO: A cobra comeu, então alocamos um espaço vazio no final.
			g.snake = append(g.snake, Point{})
			g.spawnFood()
		}

		// 7. OTIMIZAÇÃO DE MOVIMENTO: Desloca o corpo de trás pra frente
		for i := len(g.snake) - 1; i > 0; i-- {
			g.snake[i] = g.snake[i-1]
		}

		// 8. Define a nova cabeça na posição [0]
		g.snake[0] = newHead

		return nil
	}
	return nil
}

// Draw desenha os elementos na tela
func (g *Game) Draw(screen *ebiten.Image) {
	switch g.state {

	case StateTitle:
		ebitenutil.DebugPrint(screen, "=== JOGO DA COBRA COM AURA PORRA ===\n\nAperte Espaço para iniciar")

	case StatePlaying:
		// Aqui você desenha a cobra e a comida
		texto := fmt.Sprintf("\nPontos: %d\n\nAperte Esc para pausar", g.score)
		ebitenutil.DebugPrint(screen, texto)

	case StatePaused:
		// Como não colocamos um 'break' ou 'return' antes de desenhar o pause,
		// você pode chamar a função que desenha o jogo de fundo aqui, ou apenas
		// escrever por cima da última tela.
		ebitenutil.DebugPrint(screen, "--- JOGO PAUSADO ---\n\nAperte Esc para continuar")

	case StateGameOver:
		texto := fmt.Sprintf("GAME OVER!\n Pontos: %d\nAperte ENTER para tentar novamente", g.score)
		ebitenutil.DebugPrint(screen, texto)
	}
	if g.state == StatePlaying || g.state == StatePaused {
		// Desenha a comida (quadrado vermelho)
		ebitenutil.DrawRect(screen, float64(g.food.x*gridSize), float64(g.food.y*gridSize), float64(gridSize), float64(gridSize), color.RGBA{255, 0, 0, 255})

		// Desenha o corpo da cobra (quadrados verdes)
		for _, p := range g.snake {
			ebitenutil.DrawRect(screen, float64(p.x*gridSize), float64(p.y*gridSize), float64(gridSize), float64(gridSize), color.RGBA{0, 255, 0, 255})
		}
	}
}

// Layout define o tamanho interno da tela do jogo
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

// reset volta o jogo para o estado inicial
func (g *Game) reset() {
	g.score = 0
	g.speed = 10
	g.snake = []Point{{16, 12}} // Cobra começa no meio
	g.dir = Point{1, 0}         // Começa andando para a direita
	g.spawnFood()
}
func (g *Game) spawnFood() {
	var espacoslivres []Point //pontos possíveis de spawnar a maçã
	cols := screenWidth / gridSize
	lins := screenHeight / gridSize
	// loop mapeia todas as posições possíveis da tela
	for x := 0; x < cols; x++ {
		for y := 0; y < lins; y++ {
			p := Point{x, y}
			isFree := true

			// Checa se a coordenada está ocupada por alguma parte da cobra
			for _, sp := range g.snake {
				if p == sp {
					isFree = false
					break
				}
			}

			// Se o espaço estiver livre, adiciona na lista de opções
			if isFree {
				espacoslivres = append(espacoslivres, p)
			}
		}
	}
	// se houver espaços vagos, coloca a maçã em um deles aleatoriamente
	if len(espacoslivres) > 0 {
		g.food = espacoslivres[rand.Intn(len(espacoslivres))]
	}
}
func main() {
	game := &Game{}
	game.reset()

	// Configurações da janela do jogo
	ebiten.SetWindowSize(screenWidth*2, screenHeight*2)
	ebiten.SetWindowTitle("Jogo da Cobrinha com aura")

	// Inicia o loop principal
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
