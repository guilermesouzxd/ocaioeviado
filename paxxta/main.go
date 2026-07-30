package main

import (
	"image/color"
	"log"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
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
	snake []Point
	dir   Point
	food  Point
	ticks int
}

// Update contém a lógica do jogo (roda 60 vezes por segundo)
func (g *Game) Update() error {
	// 1. Controle de direção pelas setas do teclado
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
	if g.ticks < 1 {
		return nil
	}
	g.ticks = 0

	// 3. Calcular a nova posição da cabeça da cobra
	head := g.snake[0]
	newHead := Point{head.x + g.dir.x, head.y + g.dir.y}

	// 4. Checar colisão com as paredes (Game Over: reseta o jogo)
	if newHead.x < 0 || newHead.y < 0 || newHead.x >= screenWidth/gridSize || newHead.y >= screenHeight/gridSize {
		g.reset()
		return nil
	}

	// 5. Checar colisão com o próprio corpo
	for _, p := range g.snake {
		if p == newHead {
			g.reset()
			return nil
		}
	}

	// 6. Mover a cobra (adiciona nova cabeça no início da lista)
	g.snake = append([]Point{newHead}, g.snake...)

	// 7. Lógica de comer a maçã
	if newHead == g.food {
		// Gera uma nova maçã em um local aleatório
		g.food = Point{rand.Intn(screenWidth / gridSize), rand.Intn(screenHeight / gridSize)}
	} else {
		// Remove a ponta da cauda se não comeu a maçã (para não crescer infinitamente)
		g.snake = g.snake[:len(g.snake)-1]
	}

	return nil
}

// Draw desenha os elementos na tela
func (g *Game) Draw(screen *ebiten.Image) {
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
	g.snake = []Point{{16, 12}} // Cobra começa no meio
	g.dir = Point{1, 0}         // Começa andando para a direita
	g.food = Point{rand.Intn(screenWidth / gridSize), rand.Intn(screenHeight / gridSize)}
}

func main() {
	game := &Game{}
	game.reset()

	// Configurações da janela do jogo
	ebiten.SetWindowSize(screenWidth*2, screenHeight*2)
	ebiten.SetWindowTitle("Jogo da Cobrinha em Go")

	// Inicia o loop principal
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
