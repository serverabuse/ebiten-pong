package main

type Ball struct {
	form   *Rect
	pos    Position
	SpeedX int
	SpeedY int
}

func NewBall(x, y int) *Ball {
	return &Ball{
		form: NewRect(10, 10, x, y),
		pos: Position{
			X: x,
			Y: y,
		},
		SpeedX: 0,
		SpeedY: 0,
	}
}
