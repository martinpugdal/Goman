package main

import (
	"image/png"
	"math"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
)

// pacman struct
type Pacman struct {
	Image      *ebiten.Image // the image of Pacman
	X, Y       float64       // position in pixels
	Moving     bool          // is moving?
	DirX, DirY int           // direction, e.g. (1,0)=right, (0,1)=down
}

// load Pacman image from a file
func NewPacman(path string, x, y float64) (*Pacman, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	img, err := png.Decode(file)
	if err != nil {
		return nil, err
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
	return &Pacman{
		Image:  dst,
		X:      x,
		Y:      y,
		Moving: false,
	}, nil
}

func (p *Pacman) Draw(screen *ebiten.Image) {
	if p.Image == nil {
		return
	}
	op := &ebiten.DrawImageOptions{}

	// based on direction, set the image rotation
	angle := 0.0
	switch {
	case p.DirX == 1 && p.DirY == 0:
		angle = 0
	case p.DirX == 0 && p.DirY == 1:
		angle = 90
	case p.DirX == -1 && p.DirY == 0:
		angle = 180
	case p.DirX == 0 && p.DirY == -1:
		angle = -90
	}

	// rotate image
	op.GeoM.Translate(-float64(TileSize)/2, -float64(TileSize)/2)
	op.GeoM.Rotate(angle * (math.Pi / 180))
	op.GeoM.Translate(float64(TileSize)/2, float64(TileSize)/2)
	op.GeoM.Translate(p.X, p.Y)
	screen.DrawImage(p.Image, op)
}
