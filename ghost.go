package main

import (
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	Red    = "red"
	Blue   = "blue"
	Pink   = "pink"
	Orange = "orange"
)

type Ghost struct {
	Image      *ebiten.Image // image of the ghost
	X, Y       float64       // position in pixels
	DirX, DirY int           // current direction
	Color      string        // color of the ghost
	moveTick   int           // tick for movement delay
}

func (g *Ghost) Update(pacmanX, pacmanY float64) {
	switch g.Color {

	case Red:
		g.chasePacman(pacmanX, pacmanY)
	case Pink:
		g.tileRandomMove()

	default:
		g.randomMove()
	}
}

// show and random move
func (g *Ghost) tileRandomMove() {
	// 3 frames delay for movement
	g.moveTick = (g.moveTick + 1) % 3
	if g.moveTick != 0 {
		return
	}
	tileX := int((g.X + float64(TileSize)/2) / float64(TileSize))
	tileY := int((g.Y + float64(TileSize)/2) / float64(TileSize))

	if int(g.X)%TileSize == 0 && int(g.Y)%TileSize == 0 {
		dirs := [][2]int{{1, 0}, {-1, 0}, {0, 1}, {0, -1}}
		valid := [][2]int{}
		for _, d := range dirs {
			nx, ny := tileX+d[0], tileY+d[1]
			if nx >= 0 && ny >= 0 && ny < len(maze) && nx < len(maze[0]) && maze[ny][nx] != 0 {
				// exclude to go back to the previous tile
				if !(g.DirX == -d[0] && g.DirY == -d[1]) {
					valid = append(valid, d)
				}
			}
		}
		// if no valid directions, stay still
		if len(valid) > 0 {
			choice := valid[rand.Intn(len(valid))]
			g.DirX, g.DirY = choice[0], choice[1]
		} else {
			g.DirX, g.DirY = 0, 0
		}
	}

	// apply the direction
	g.X += float64(g.DirX)
	g.Y += float64(g.DirY)
}

func (g *Ghost) Draw(screen *ebiten.Image) {
	if g.Image == nil {
		return
	}
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(g.X, g.Y)
	screen.DrawImage(g.Image, op)
}

// just chasing the pacman with tile-based movement
func (g *Ghost) chasePacman(px, py float64) {
	g.moveTick = (g.moveTick + 1) % 2 // 2 frames delay for movement
	if g.moveTick != 0 {
		return
	}
	tileX := int((g.X + float64(TileSize)/2) / float64(TileSize))
	tileY := int((g.Y + float64(TileSize)/2) / float64(TileSize))
	pacTileX := int((px + float64(TileSize)/2) / float64(TileSize))
	pacTileY := int((py + float64(TileSize)/2) / float64(TileSize))

	if int(g.X)%TileSize == 0 && int(g.Y)%TileSize == 0 {
		dirs := [][2]int{{1, 0}, {-1, 0}, {0, 1}, {0, -1}}
		bestDir := [2]int{0, 0}
		minDist := 999999.0
		for _, d := range dirs {
			nx, ny := tileX+d[0], tileY+d[1]
			if nx >= 0 && ny >= 0 && ny < len(maze) && nx < len(maze[0]) && maze[ny][nx] != 0 {
				// exclude to go back to the previous tile
				if !(g.DirX == -d[0] && g.DirY == -d[1]) {
					dist := math.Abs(float64(nx-pacTileX)) + math.Abs(float64(ny-pacTileY))
					if dist < minDist {
						minDist = dist
						bestDir = d
					}
				}
			}
		}
		g.DirX, g.DirY = bestDir[0], bestDir[1]
	}
	g.X += float64(g.DirX)
	g.Y += float64(g.DirY)
}

// fall back to random movement, if no specific behavior is defined for the ghost
func (g *Ghost) randomMove() {
	g.X += float64((rand.Intn(3) - 1) * 2)
	g.Y += float64((rand.Intn(3) - 1) * 2)
}
