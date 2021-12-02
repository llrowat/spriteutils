package spriteutils_test

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/llrowat/spriteutils"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSpriteFactory_GenerateSprite(t *testing.T) {
	testImage, _ := ebiten.NewImage(5, 5, ebiten.FilterDefault)

	var subject = &spriteutils.SpriteFactory{
		Images: []*ebiten.Image{ testImage },
	}

	var testValues = []struct {
		minX, maxX, minY, maxY int
	}{
		{0, 0, 0, 0},
		{-10, 10, -10, 10},
		{1, 10, 1, 10},
		{-10, -5, -10, -5},
	}

	for _, testValue := range testValues {
		subject.MinX = testValue.minX
		subject.MaxX = testValue.maxX
		subject.MinY = testValue.minY
		subject.MaxY = testValue.maxY

		// Generate a bunch of sprites from factory and ensure the assertions remain true for all.
		for i := 0; i < 50; i++ {
			result := subject.GenerateSprite()

			assert.True(t, result.X >= testValue.minX)
			assert.True(t, result.X <= testValue.maxX)
			assert.True(t, result.Y >= testValue.minY)
			assert.True(t, result.Y <= testValue.maxY)
		}
	}
}
