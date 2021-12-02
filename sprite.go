// Package spriteutils provides some utility functions for dealing with 2D sprites on top of the ebiten game library
package spriteutils

import (
	"github.com/hajimehoshi/ebiten"
	"image"
	"math"
)

// SimpleSprite is a sprite interface with methods to draw and update
type SimpleSprite interface {
	// Update the sprite
	Update()
	// Draw the sprite to screen
	Draw(screen *ebiten.Image) error
}

// Sprite represents an image with position, rotation, and velocity
type Sprite struct {
	// Image is the ebiten image to draw for this sprite
	Image *ebiten.Image
	// X is the sprites x-axis position in 2D space
	X int
	// Y is the sprites y-axis position in 2D space
	Y int
	// XVelocity is the sprite's velocity in the x-axis
	XVelocity float64
	// YVelocity is the sprite's velocity in the y-axis
	YVelocity float64
	// Rotation is the sprite's rotation in radians
	Rotation float64
}

// Update the sprite by applying its velocity to position
func (sprite *Sprite) Update() {
	sprite.Y += int(sprite.YVelocity)
	sprite.X += int(sprite.XVelocity)
}

// Draw the sprite to screen after applying rotation and translation transformations
func (sprite *Sprite) Draw(screen *ebiten.Image) error {
	op := &ebiten.DrawImageOptions{}

	// make sure rotation occurs around mid-point
	width, height := sprite.Image.Size()
	op.GeoM.Translate(-float64(width)/2.0, -float64(height)/2.0)
	op.GeoM.Rotate(sprite.Rotation)
	op.GeoM.Translate(float64(width)/2.0, float64(height)/2.0)

	op.GeoM.Translate(float64(sprite.X), float64(sprite.Y))
	return screen.DrawImage(sprite.Image, op)
}

// IsColliding determines whether there is a collision between this sprite and another.
// A collision is determined to have occurred if the non-transparent sprite images are touching
func (sprite *Sprite) IsColliding(otherSprite *Sprite) bool {
	// The "hitbox" is considered to be the real bounds of the image
	spriteHitbox := image.Rectangle{
		Min: image.Point{X: sprite.X, Y: sprite.Y},
		Max: image.Point{X: sprite.X + sprite.Image.Bounds().Dx(), Y: sprite.Y + sprite.Image.Bounds().Dy()},
	}
	otherSpriteHitbox := image.Rectangle{
		Min: image.Point{X: otherSprite.X, Y: otherSprite.Y},
		Max: image.Point{X: otherSprite.X + otherSprite.Image.Bounds().Dx(), Y: otherSprite.Y + otherSprite.Image.Bounds().Dy()},
	}

	// If the hitboxes don't overlap then there can be no collision
	if !spriteHitbox.Overlaps(otherSpriteHitbox) {
		return false
	}

	// Get the rectangle representing the overlap of the two sprites hitboxes
	intersection := spriteHitbox.Intersect(otherSpriteHitbox.Bounds())

	// Go through each pixel in the intersection rectangle.  If the corresponding pixel in both images is non-transparent
	// then we consider this to be a collision.
	for i := intersection.Min.X; i < intersection.Max.X; i++ {
		for y := intersection.Min.Y; y < intersection.Max.Y; y++ {
			var _, _, _, spritePixelAlpha = sprite.Image.At(rotatePoint(i-sprite.X, y-sprite.Y, sprite.Rotation, sprite.Image.Bounds().Dx()/2, sprite.Image.Bounds().Dy()/2)).RGBA()
			var _, _, _, otherSpritePixelAlpha = otherSprite.Image.At(rotatePoint(i-otherSprite.X, y-otherSprite.Y, otherSprite.Rotation, otherSprite.Image.Bounds().Dx()/2, otherSprite.Image.Bounds().Dy()/2)).RGBA()
			if spritePixelAlpha != 0 && otherSpritePixelAlpha != 0 {
				return true
			}
		}
	}

	return false
}

// ApplyImpulse applies a 2d vector force represented by (xVelocity, yVelocity) to the sprite
func (sprite *Sprite) ApplyImpulse(xVelocity, yVelocity float64) {
	sprite.XVelocity += xVelocity
	sprite.YVelocity += yVelocity
}

// rotatePoint takes a point (x,y) and rotates it by theta (radians) around the origin point (originX, originY)
// to get the resulting point post-rotation
func rotatePoint(x int, y int, theta float64, originX int, originY int) (int, int) {
	sinTheta := math.Sin(theta)
	cosTheta := math.Cos(theta)

	tx := x - originX
	ty := y - originY

	rx := int(float64(tx)*cosTheta-float64(ty)*sinTheta) + originX
	ry := int(float64(tx)*sinTheta+float64(ty)*cosTheta) + originY

	return rx, ry
}
