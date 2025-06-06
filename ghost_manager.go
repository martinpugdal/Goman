package main

import (
	"image/png"
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
)

type GhostManager struct {
	Ghosts []*Ghost
}

func (gm *GhostManager) InitGhosts() {
	gm.Ghosts = []*Ghost{
		loadGhost(Red, 8, 6),
		loadGhost(Blue, 8, 7),
		loadGhost(Orange, 8, 8),
		loadGhost(Pink, 9, 8),
	}
}

func loadGhost(color string, tileX, tileY int) *Ghost {
	file, err := os.Open("assets/" + color + "-ghost.png")
	if err != nil {
		log.Printf("Could not load ghost image for %s: %v", color, err)
		return &Ghost{Color: color, X: float64(tileX * TileSize), Y: float64(tileY * TileSize)}
	}
	defer file.Close()
	img, err := png.Decode(file)
	if err != nil {
		log.Printf("Could not decode ghost image for %s: %v", color, err)
		return &Ghost{Color: color, X: float64(tileX * TileSize), Y: float64(tileY * TileSize)}
	}
	w, h := img.Bounds().Dx(), img.Bounds().Dy()
	var scale float64
	if w > h {
		scale = float64(TileSize) / float64(w)
	} else {
		scale = float64(TileSize) / float64(h)
	}
	scaledW := float64(w) * scale
	scaledH := float64(h) * scale
	offsetX := (float64(TileSize) - scaledW) / 2
	offsetY := (float64(TileSize) - scaledH) / 2
	dst := ebiten.NewImage(TileSize, TileSize)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(scale, scale)
	op.GeoM.Translate(offsetX, offsetY)
	dst.DrawImage(ebiten.NewImageFromImage(img), op)
	return &Ghost{
		Image: dst,
		X:     float64(tileX * TileSize),
		Y:     float64(tileY * TileSize),
		Color: color,
	}
}

func (gm *GhostManager) Update(pacmanX, pacmanY float64) {
	for _, g := range gm.Ghosts {
		g.Update(pacmanX, pacmanY)
	}
}

func (gm *GhostManager) Draw(screen *ebiten.Image) {
	for _, g := range gm.Ghosts {
		g.Draw(screen)
	}
}
