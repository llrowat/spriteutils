// Package spriteutils provides some utility functions for dealing with 2D sprites on top of the ebiten game library
package spriteutils

import (
	"github.com/hajimehoshi/ebiten"
	"time"
)

// TransientSprite is a sprite that exists for a period of time.
// Once it is expired it will no longer Draw or Update, but it is up to the user to remove any references so that it is garbage collected
type TransientSprite struct {
	// CreatedAtGameTime is the duration of time that has past since the game started.
	// We cannot use real time since if the game window is unfocused, it will pause.
	CreatedAtGameTime time.Duration
	// LifetimeDuration is the duration of time to keep this sprite alive
	LifetimeDuration time.Duration
	// Sprite is the limited lifetime sprite
	Sprite SimpleSprite
	// IsExpired represents whether the sprite lifetime has surpassed or not
	IsExpired bool
}

// Draw the sprite to screen if it has not expired
func (t *TransientSprite) Draw(screen *ebiten.Image) error {
	if !t.IsExpired && t.Sprite != nil {
		return t.Sprite.Draw(screen)
	}

	return nil
}

// Update the sprite is it has not expired
func (t *TransientSprite) Update(currentGameTime time.Duration) {
	t.checkExpired(currentGameTime)

	if !t.IsExpired && t.Sprite != nil {
		t.Sprite.Update()
	}
}

// checkExpired will set IsExpired and nil the Sprite if it's lifetime has elapsed
func (t *TransientSprite) checkExpired(currentGameTime time.Duration) {
	if currentGameTime-t.CreatedAtGameTime > t.LifetimeDuration {
		t.Sprite = nil
		t.IsExpired = true
	}

	t.IsExpired = false
}
