package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"image/color"
)

type Rect struct {
	height int
	width  int
	pos    Position
}

func (r Rect) Draw() *ebiten.Image {
	img := ebiten.NewImage(r.width, r.height)
	img.Fill(color.White)
	return img
}

func NewRect(height, width, x, y int) *Rect {
	return &Rect{
		height: height,
		width:  width,
		pos: Position{
			X: x,
			Y: y,
		},
	}
}
