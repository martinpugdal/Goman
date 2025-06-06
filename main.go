package main

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// Constants for screen and tile sizes
const (
	ScreenWidth  = 640 * 2
	ScreenHeight = 480 * 2
	TileSize     = 32 * 2
)

// 0 = wall, 1 = path, 2 = apple
var maze = [][]int{
	{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	{0, 2, 2, 2, 0, 2, 2, 2, 2, 2, 2, 2, 2, 2, 0, 2, 2, 2, 2, 0},
	{0, 2, 0, 2, 0, 2, 0, 0, 0, 2, 0, 0, 0, 2, 0, 2, 0, 2, 2, 0},
	{0, 2, 0, 2, 2, 2, 2, 0, 0, 2, 0, 0, 2, 2, 2, 2, 0, 2, 2, 0},
	{0, 2, 0, 0, 0, 0, 2, 0, 0, 2, 0, 0, 2, 0, 0, 0, 0, 0, 2, 0},
	{0, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 0},
	{0, 2, 0, 0, 0, 0, 2, 0, 0, 1, 0, 0, 2, 2, 0, 0, 0, 0, 2, 0},
	{0, 2, 2, 2, 2, 2, 2, 0, 1, 1, 1, 0, 2, 2, 2, 2, 2, 2, 2, 0},
	{0, 2, 0, 0, 0, 0, 2, 0, 0, 1, 0, 0, 2, 0, 0, 0, 0, 0, 2, 0},
	{0, 2, 2, 2, 2, 2, 2, 0, 0, 0, 0, 0, 2, 2, 2, 2, 2, 2, 2, 0},
	{0, 2, 0, 2, 0, 0, 0, 0, 2, 0, 0, 2, 0, 0, 0, 0, 2, 0, 2, 0},
	{0, 2, 2, 2, 2, 2, 2, 0, 2, 2, 2, 2, 0, 2, 2, 2, 2, 2, 2, 0},
	{0, 2, 0, 0, 0, 0, 2, 0, 0, 0, 0, 0, 0, 2, 0, 0, 0, 0, 2, 0},
	{0, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 0},
	{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
}

type Game struct {
	pacmanX, pacmanY   int          // placement for Pacman in the maze
	pacman             *Pacman      // Pacman instance
	apples             AppleManager // Apple manager for collecting apples
	ghosts             GhostManager // Ghost manager for handling ghosts
	moveDelay          int          // frames between movements
	moveCounter        int          // frame counter
	lastDirX, lastDirY int          // last direction
}

func (g *Game) Update() error {
	g.ghosts.Update(g.pacman.X, g.pacman.Y)

	dx, dy := 0, 0
	moving := false

	// input movement check
	switch {
	case ebiten.IsKeyPressed(ebiten.KeyArrowUp) || ebiten.IsKeyPressed(ebiten.KeyW):
		dx, dy = 0, -1
		moving = true
	case ebiten.IsKeyPressed(ebiten.KeyArrowDown) || ebiten.IsKeyPressed(ebiten.KeyS):
		dx, dy = 0, 1
		moving = true
	case ebiten.IsKeyPressed(ebiten.KeyArrowLeft) || ebiten.IsKeyPressed(ebiten.KeyA):
		dx, dy = -1, 0
		moving = true
	case ebiten.IsKeyPressed(ebiten.KeyArrowRight) || ebiten.IsKeyPressed(ebiten.KeyD):
		dx, dy = 1, 0
		moving = true
	}

	g.moveCounter++

	// only move if the delay has passed and we are moving
	if moving && g.moveCounter >= g.moveDelay {
		g.moveCounter = 0
		newX, newY := g.pacmanX+dx, g.pacmanY+dy
		// if the new position is within bounds and not a wall (0)
		if newX >= 0 && newX < len(maze[0]) && newY >= 0 && newY < len(maze) && maze[newY][newX] != 0 {
			g.pacmanX, g.pacmanY = newX, newY
		}
	}

	// if we are not moving, reset the move counter
	if !moving {
		g.moveCounter = g.moveDelay
	}

	// update pacman info based on movement
	if g.pacman != nil {
		g.pacman.X = float64(g.pacmanX * TileSize)
		g.pacman.Y = float64(g.pacmanY * TileSize)
		g.pacman.Moving = moving

		if moving && (dx != 0 || dy != 0) {
			g.pacman.DirX = dx
			g.pacman.DirY = dy
		}
	}

	if g.apples.Apples != nil {
		for _, apple := range g.apples.Apples {
			// Check if Pacman collects an apple
			if apple.X == g.pacmanX && apple.Y == g.pacmanY && !apple.Collected() {
				apple.Collect()
			}
		}
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	for y, row := range maze {
		for x, tile := range row {
			switch tile {
			case 0:
				vector.DrawFilledRect(screen, float32(x*TileSize), float32(y*TileSize), float32(TileSize), float32(TileSize), color.RGBA{200, 200, 200, 255}, false)
			case 1, 2:
				vector.DrawFilledRect(screen, float32(x*TileSize), float32(y*TileSize), float32(TileSize), float32(TileSize), color.RGBA{100, 100, 100, 255}, false)
			}
		}
	}

	g.apples.Draw(screen)
	g.ghosts.Draw(screen)
	g.pacman.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return ScreenWidth, ScreenHeight
}

func main() {
	pacman, err := NewPacman("assets/yellow-pacman.png", float64(1*TileSize), float64(1*TileSize))
	if err != nil {
		log.Fatalf("Failed to load Pacman image: %v", err)
	}

	game := &Game{
		pacmanX:     1,
		pacmanY:     1,
		pacman:      pacman,
		moveDelay:   10,
		moveCounter: 0,
		lastDirX:    0,
		lastDirY:    0,
	}

	game.apples.InitApples(maze)
	game.ghosts.InitGhosts()

	// Create the Ebiten game
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetWindowSize(ScreenWidth, ScreenHeight)
	ebiten.SetWindowTitle("Goman (Pacman in Go)")
	err = ebiten.RunGame(game)
	if err != nil {
		log.Fatalf("Failed to run game: %v", err)
	}
}
