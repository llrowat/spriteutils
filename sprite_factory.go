// Package spriteutils provides some utility functions for dealing with 2D sprites on top of the ebiten game library
package spriteutils

import (
	"github.com/hajimehoshi/ebiten"
	"math/rand"
)

// SpriteFactory is a helper for creating sprites in random x, y positions
type SpriteFactory struct {
	// Images represent the possible images for the generated sprite.  Will be chosen randomly when generating
	Images []*ebiten.Image
	// MinX represents the generated sprite's minimum x-axis position in 2D space
	MinX int
	// MaxX represents the generated sprite's maximum x-axis position in 2D space
	MaxX int
	// MinY represents the generated sprite's minimum y-axis position in 2D space
	MinY int
	// MaxY represents the generated sprite's maximum y-axis position in 2D space
	MaxY int
}

// GenerateSprite generates a sprite using the factories settings
func (factory *SpriteFactory) GenerateSprite() *Sprite {
	x := rand.Intn(factory.MaxX-factory.MinX+1) + factory.MinX
	y := rand.Intn(factory.MaxY-factory.MinY+1) + factory.MinY
	image := rand.Intn(len(factory.Images))

	return &Sprite{
		Image:     factory.Images[image],
		X:         x,
		Y:         y,
		XVelocity: 0,
		YVelocity: 0,
		Rotation:  0,
	}
}
