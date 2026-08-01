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
	state GameState
	snake []Point
	dir   Point
	food  Point
	ticks int
	score int
	speed int
}

// Update contém a lógica do jogo (roda 60 vezes por segundo)
func (g *Game) Update() error {
	switch g.state {
	case StateTitle: //caso esteja na tela de titulo, aperta espaço p começar
		if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
			g.state = StatePlaying
			g.reset()
		}
	case StatePaused: // caso esteja pausado, aperte esc p jogar
		if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
			g.state = StatePlaying
		}
	case StateGameOver: // se game over, apaert enter p jogar
		if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
			g.reset()
			g.state = StateTitle
		}
	case StatePlaying: // o jogo roda todo aí enquanto estiver nesse estado, aperta esc p pausar
		if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
			g.state = StatePaused
		}
		if ebiten.IsKeyPressed(ebiten.KeyW) && g.dir.y != 1 {
			g.dir = Point{0, -1}
		} else if ebiten.IsKeyPressed(ebiten.KeyS) && g.dir.y != -1 {
			g.dir = Point{0, 1}
		} else if ebiten.IsKeyPressed(ebiten.KeyA) && g.dir.x != 1 {
			g.dir = Point{-1, 0}
		} else if ebiten.IsKeyPressed(ebiten.KeyD) && g.dir.x != -1 {
			g.dir = Point{1, 0}
		}

		// 2. Controle de velocidade (a cobra move a cada 10 frames)
		g.ticks++
		if g.ticks < g.speed {
			return nil
		}
		g.ticks = 0
		// 3. Calcular a nova posição da cabeça da cobra
		head := g.snake[0]
		newHead := Point{head.x + g.dir.x, head.y + g.dir.y}

		// 4. Checar colisão com as paredes (Game Over: reseta o jogo)
		if newHead.x < 0 || newHead.y < 0 || newHead.x >= screenWidth/gridSize || newHead.y >= screenHeight/gridSize {
			g.state = StateGameOver
			return nil
		}

		// 5. Checar colisão com o próprio corpo
		for _, p := range g.snake {
			if p == newHead {
				g.state = StateGameOver
			}
		}
		g.snake = append([]Point{newHead}, g.snake...)
		//7. funcao de comer
		if newHead == g.food {
			g.score += 10
			if g.score%50 == 0 && g.speed > 3 {
				g.speed--
			}
			g.spawnFood()
		} else {
			// Remove a ponta da cauda se não comeu a maçã (para não crescer infinitamente)
			g.snake = g.snake[:len(g.snake)-1]
		}
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
	// Desenha a comida (quadrado vermelho)
	ebitenutil.DrawRect(screen, float64(g.food.x*gridSize), float64(g.food.y*gridSize), float64(gridSize), float64(gridSize), color.RGBA{255, 0, 0, 255})

	// Desenha o corpo da cobra (quadrados verdes)
	for _, p := range g.snake {
		ebitenutil.DrawRect(screen, float64(p.x*gridSize), float64(p.y*gridSize), float64(gridSize), float64(gridSize), color.RGBA{0, 255, 0, 255})
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
	g.food = Point{rand.Intn(screenWidth / gridSize), rand.Intn(screenHeight / gridSize)}
	g.spawnFood()
}
func (g *Game) spawnFood() {
	for { // coordenadas de spawn da maçã
		g.food = Point{
			x: rand.Intn(screenWidth / gridSize),
			y: rand.Intn(screenHeight / gridSize),
		}

		// Checa se a nova posição caiu em cima da cobra
		emCimaDaCobra := false
		for _, p := range g.snake {
			if p == g.food {
				emCimaDaCobra = true
				break // Para a verificação, já sabemos que deu ruim
			}
		}

		// Se NÃO estiver em cima da cobra, a posição é válida!
		if !emCimaDaCobra {
			break // Sai do loop infinito
		}
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
