package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"golang.org/x/image/colornames"
	"image/color"
	"log"
	"math"
	"os"
)

const (
	screenWidth  = 800
	screenHeight = 400
	paddingX     = 10
	paddingY     = 10
)

var (
	targetDifficulty  = "medium"
	difficultyMapping = map[string]int{
		"easy":   5,
		"medium": 10,
		"hard":   15,
	}
	playerSpeed, cpuSpeed = difficultyMapping[targetDifficulty], difficultyMapping[targetDifficulty]
)

type Game struct {
	cpuPosOffset    int
	playerPosOffset int
	rects           []*Rect
	ball            *Ball
	p1Score         int
	cpuScore        int
}

func drawCircle(screen *ebiten.Image, x, y, radius int, clr color.Color) {
	radius64 := float64(radius)
	minAngle := math.Acos(1 - 1/radius64)

	for angle := float64(0); angle <= 360; angle += minAngle {
		xDelta := radius64 * math.Cos(angle)
		yDelta := radius64 * math.Sin(angle)

		x1 := int(math.Round(float64(x) + xDelta))
		y1 := int(math.Round(float64(y) + yDelta))

		screen.Set(x1, y1, clr)
	}
}

func (g *Game) Reset() {
	g.ball.pos.X = 2*paddingX + g.rects[0].width
	g.ball.pos.Y = screenHeight / 2
	g.ball.SpeedX, g.ball.SpeedY = 0, 0
	g.rects[0].pos.X, g.rects[0].pos.Y = paddingX, (screenHeight/2)-(g.rects[0].height/2)
	g.rects[1].pos.X, g.rects[1].pos.Y = screenWidth-(2*paddingX), (screenHeight/2)-(g.rects[1].height/2)
}

func (g *Game) Update() error {

	// Other Inputs
	if g.ball.SpeedX == 0 {
		if ebiten.IsKeyPressed(ebiten.KeyEscape) {
			os.Exit(0)
		}
		if ebiten.IsKeyPressed(ebiten.KeyE) {
			targetDifficulty = "easy"
		}
		if ebiten.IsKeyPressed(ebiten.KeyM) {
			targetDifficulty = "medium"
		}
		if ebiten.IsKeyPressed(ebiten.KeyH) {
			targetDifficulty = "hard"
		}
	}
	// -----------

	// P1 Rect Movement
	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		if g.rects[0].pos.Y+g.rects[0].height < screenHeight-paddingY {
			g.rects[0].pos.Y += playerSpeed
		}
	}
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		if g.playerPosOffset+g.rects[0].pos.Y > paddingY {
			g.rects[0].pos.Y -= playerSpeed
		}
	}
	if ebiten.IsKeyPressed(ebiten.KeySpace) {
		if g.ball.SpeedX == 0 {
			if g.ball.pos.Y < g.rects[0].pos.Y+g.rects[0].height && g.ball.pos.Y > g.rects[0].pos.Y {
				g.ball.SpeedX = difficultyMapping[targetDifficulty]
			}
		}
	}
	// -----------

	// CPU Rect movement
	if g.ball.SpeedX != 0 {
		if g.rects[1].pos.Y+g.rects[1].height >= screenHeight-paddingY {
			cpuSpeed = -difficultyMapping[targetDifficulty]
		}
		if g.rects[1].pos.Y <= paddingY {
			cpuSpeed = difficultyMapping[targetDifficulty]
		}
		g.rects[1].pos.Y += cpuSpeed
	}
	// -----------

	// Ball Movement
	if g.ball.SpeedX != 0 {
		// Ball hits P1
		if g.ball.pos.X <= (2*paddingX)+g.rects[0].width {
			// CPU Goal !
			if g.ball.pos.Y < g.rects[0].pos.Y || g.ball.pos.Y > (g.rects[0].pos.Y+g.rects[0].height) {
				g.cpuScore++
				g.Reset()
				return nil
			}
			g.ball.SpeedX = difficultyMapping[targetDifficulty]
			if g.ball.pos.Y < (g.rects[0].pos.Y + g.rects[0].height/2) {
				g.ball.SpeedY = -difficultyMapping[targetDifficulty]
			} else if g.ball.pos.Y > (g.rects[0].pos.Y + g.rects[0].height/2) {
				g.ball.SpeedY = difficultyMapping[targetDifficulty]
			} else {
				g.ball.SpeedY = 0
			}
		}
		// Ball hits CPU
		if g.ball.pos.X >= screenWidth-(2*paddingX)-g.rects[1].width {
			// P1 Goal !
			if g.ball.pos.Y < g.rects[1].pos.Y || g.ball.pos.Y > (g.rects[1].pos.Y+g.rects[1].height) {
				g.p1Score++
				g.Reset()
				return nil
			}
			g.ball.SpeedX = -difficultyMapping[targetDifficulty]
			if g.ball.pos.Y < (g.rects[1].pos.Y + g.rects[1].height/2) {
				g.ball.SpeedY = -difficultyMapping[targetDifficulty]
			} else if g.ball.pos.Y > (g.rects[1].pos.Y + g.rects[1].height/2) {
				g.ball.SpeedY = difficultyMapping[targetDifficulty]
			} else {
				g.ball.SpeedY = 0
			}
		}

		g.ball.pos.X += g.ball.SpeedX
		g.ball.pos.Y += g.ball.SpeedY
	}
	// Ball hits Wall Up/Down
	if g.ball.SpeedY != 0 {
		if g.ball.pos.Y <= paddingY {
			g.ball.SpeedY = difficultyMapping[targetDifficulty]
		}
		if g.ball.pos.Y >= screenHeight-paddingY {
			g.ball.SpeedY = -difficultyMapping[targetDifficulty]
		}
	}
	// -----------
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Draw Field
	screen.Fill(colornames.Black)
	// Middle vertical line && circle
	ebitenutil.DrawLine(screen, screenWidth/2, float64(paddingY-5), screenWidth/2, float64(screenHeight-paddingY+5), colornames.White)
	drawCircle(screen, screenWidth/2, screenHeight/2, 100, colornames.White)
	// Top Line
	ebitenutil.DrawLine(screen, float64(paddingX-5), float64(paddingY-5), float64(screenWidth-paddingX+5), float64(paddingY-5), colornames.White)
	// Bottom line
	ebitenutil.DrawLine(screen, float64(paddingX-5), float64(screenHeight-paddingY+5), float64(screenWidth-paddingX+5), float64(screenHeight-paddingY+5), colornames.White)
	// Left line
	ebitenutil.DrawLine(screen, float64(paddingX-5), float64(paddingY-5), float64(paddingX-5), float64(screenHeight-paddingY+5), colornames.White)
	// Right line
	ebitenutil.DrawLine(screen, float64(screenWidth-paddingX+5), float64(paddingY-5), float64(screenWidth-paddingX+5), float64(screenHeight-paddingY+5), colornames.White)
	// -----------

	// Draw Ball
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(g.ball.pos.X), float64(g.ball.pos.Y))
	screen.DrawImage(g.ball.form.Draw(), op)
	// -----------

	// Draw P1 Rect
	op = &ebiten.DrawImageOptions{}
	g.rects[0].pos.Y += g.playerPosOffset
	op.GeoM.Translate(float64(g.rects[0].pos.X), float64(g.rects[0].pos.Y))
	screen.DrawImage(g.rects[0].Draw(), op)
	// -----------

	// Draw CPU Rect
	op = &ebiten.DrawImageOptions{}
	g.rects[1].pos.Y += g.cpuPosOffset
	op.GeoM.Translate(float64(g.rects[1].pos.X), float64(g.rects[1].pos.Y))
	screen.DrawImage(g.rects[1].Draw(), op)
	// -----------

	// Draw Scores
	ebitenutil.DebugPrintAt(screen, "P1", screenWidth/2-paddingX-14, paddingY+5)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%d", g.p1Score), screenWidth/2-paddingX-13, paddingY+25)
	ebitenutil.DebugPrintAt(screen, "CPU", screenWidth/2+paddingX+5, paddingY+5)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%d", g.cpuScore), screenWidth/2+paddingX+12, paddingY+25)
	// -----------

	// Draw Instructions
	if g.ball.SpeedX == 0 {
		ebitenutil.DebugPrintAt(screen, "Press SPACE to start", paddingX, screenHeight-paddingY-30)
		ebitenutil.DebugPrintAt(screen, "Move with UP/DOWN Arrow Keys", paddingX, screenHeight-paddingY-15)
		ebitenutil.DebugPrintAt(screen, "Press for Difficulty: E[Easy], M[Medium], H[Hard]", paddingX, paddingY)
		ebitenutil.DebugPrintAt(screen, "Press ESC to Exit", paddingX, paddingY+15)
	}
	// -----------

}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}

func main() {
	var (
		rectHeight = 150
		rectWidth  = 10
	)
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Pong")
	if err := ebiten.RunGame(&Game{
		cpuPosOffset:    0,
		playerPosOffset: 0,
		ball:            NewBall(2*paddingX+rectWidth, screenHeight/2),
		rects: []*Rect{
			NewRect(rectHeight, rectWidth, paddingX, (screenHeight/2)-(rectHeight/2)),                 // P1
			NewRect(rectHeight, rectWidth, screenWidth-(2*paddingX), (screenHeight/2)-(rectHeight/2)), //CPU
		},
	}); err != nil {
		log.Fatal(err)
	}
}
