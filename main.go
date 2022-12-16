package main

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	screenWidth   = 1920
	screenHeight  = 1080
	maxVel        = 1.0
	stdDevAcc     = 0.1
	maxSprites    = 2000
	searchRad     = 0.05 * screenHeight
	maxNeighbours = 3
	maxBatch      = 500
)

var (
	whiteImage    = ebiten.NewImage(5, 5)
	blankSubImage = whiteImage.SubImage(image.Rect(1, 1, 2, 2)).(*ebiten.Image)
)

func init() {
	rand.Seed(time.Now().UnixNano())
	whiteImage.Fill(color.White)
}

type Game struct {
	particles  []*Sprite
	showPoints bool
	showLines  bool
	bgColor    bool
}

func NewGame() *Game {
	particles := []*Sprite{}
	for i := 0; i < maxSprites; i++ {
		particles = append(particles, newSprite(whiteImage, screenHeight, screenWidth))
	}
	return &Game{
		particles:  particles,
		showPoints: true,
		showLines:  true,
		bgColor:    true,
	}
}

func (g *Game) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyP) {
		g.showPoints = !g.showPoints
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyL) {
		g.showLines = !g.showLines
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyB) {
		g.bgColor = !g.bgColor
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyQ) {
		return ebiten.Termination
	}
	for _, s := range g.particles {
		ax := rand.NormFloat64() * stdDevAcc
		ay := rand.NormFloat64() * stdDevAcc
		s.update(ax, ay)
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	if g.bgColor {
		screen.Fill(color.RGBA{245, 245, 245, 245})
	}
	if g.showPoints {
		for _, s := range g.particles {
			s.draw(screen)
		}
	}
	if g.showLines {
		drawAllDetectionLines(g.particles, screen)
	}

	// Parrallelization on naive approach: Not a huge improvement
	// drawParrallel(g.particles, screen)

	msg := fmt.Sprintf("FPS: %0.0f, TPS: %0.0f\nToggle KEYS: B(bg)\n    P(points)\n    L(lines)",
		ebiten.ActualFPS(), ebiten.ActualTPS())
	ebitenutil.DebugPrint(screen, msg)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	// ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetFullscreen(true)
	ebiten.SetWindowTitle("Sprite")
	if err := ebiten.RunGame(NewGame()); err != nil {
		log.Fatal(err)
	}
}
