package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

type AppleManager struct {
	Apples []*Apple
}

func (am *AppleManager) InitApples(maze [][]int) {
	am.Apples = nil
	for y, row := range maze {
		for x, tile := range row {
			if tile == 2 {
				am.Apples = append(am.Apples, &Apple{X: x, Y: y})
			}
		}
	}
}

func (am *AppleManager) Draw(screen *ebiten.Image) {
	for _, a := range am.Apples {
		a.Draw(screen)
	}
}

type Apple struct {
	X, Y  int // Position in the maze
	Eaten bool
}

func (a *Apple) Collected() bool {
	return a.Eaten
}

func (a *Apple) Collect() {
	a.Eaten = true
}

func (a *Apple) Draw(screen *ebiten.Image) {
	if a.Eaten {
		return
	}
	radius := float32(TileSize) * 0.2
	centerX := float32(a.X*TileSize) + float32(TileSize)/2
	centerY := float32(a.Y*TileSize) + float32(TileSize)/2

	// Draw a simple red circle for the apple
	ebitenutilDrawCircle(screen, centerX, centerY, radius, color.RGBA{255, 0, 0, 255})
}

func ebitenutilDrawCircle(screen *ebiten.Image, cx, cy, r float32, clr color.Color) {
	for y := -r; y <= r; y++ {
		for x := -r; x <= r; x++ {
			if x*x+y*y <= r*r {
				screen.Set(int(cx+x), int(cy+y), clr)
			}
		}
	}
}
