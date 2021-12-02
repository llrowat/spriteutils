package spriteutils_test

import (
	"errors"
	"github.com/hajimehoshi/ebiten"
	"github.com/llrowat/spriteutils"
	"github.com/stretchr/testify/assert"
	"image/color"
	"os"
	"testing"
)

var regularTermination = errors.New("regular termination")

type mockGame struct {
	m    *testing.M
	code int
}

func (g *mockGame) Update(*ebiten.Image) error {
	g.code = g.m.Run()
	return regularTermination
}

func (*mockGame) Draw(*ebiten.Image) {
}

func (*mockGame) Layout(int, int) (int, int) {
	return 320, 240
}

func MainWithRunLoop(m *testing.M) {
	// Run an Ebiten process so that (*Image).At is available.
	g := &mockGame{
		m: m,
	}
	if err := ebiten.RunGame(g); err != nil && err != regularTermination {
		panic(err)
	}
	os.Exit(g.code)
}

func TestMain(m *testing.M) {
	MainWithRunLoop(m)
}

func TestSprite_Update(t *testing.T) {
	testImage, _ := ebiten.NewImage(5, 5, ebiten.FilterDefault)
	subject := &spriteutils.Sprite{
		Image:    testImage,
	}

	var testValues = []struct {
		xVel, yVel float64
		expX, expY int
	}{
		{0, 0, 0, 0},
		{1, 0, 1, 0},
		{0, 1, 0, 1},
		{0.6, 1.2, 0, 1},
		{1, 2, 1, 2},
		{-3, -5, -3, -5},
	}

	for _, testValue := range testValues {
		subject.X = 0
		subject.Y = 0
		subject.XVelocity = testValue.xVel
		subject.YVelocity = testValue.yVel

		subject.Update()

		assert.Equal(t, testValue.expX, subject.X)
		assert.Equal(t, testValue.expY, subject.Y)
	}

}

func TestSprite_Draw(t *testing.T) {
	testImage, _ := ebiten.NewImage(5, 5, ebiten.FilterDefault)
	testScreen, _ := ebiten.NewImage(10, 10, ebiten.FilterDefault)
	subject := &spriteutils.Sprite{
		Image:    testImage,
		X:        0,
		Y:        0,
		Rotation: 0,
	}

	var result = subject.Draw(testScreen)

	assert.NoError(t, result)
}

func TestSprite_ApplyImpulse(t *testing.T) {
	testImage, _ := ebiten.NewImage(5, 5, ebiten.FilterDefault)
	subject := &spriteutils.Sprite{
		Image:    testImage,
	}

	var testValues = []struct {
		x, y, expXVel, expYVel float64
	}{
		{0, 0, 0, 0},
		{1, 1, 1, 1},
		{-1, -1, -1, -1},
	}

	for _, testValue := range testValues {
		subject.XVelocity = 0
		subject.YVelocity = 0

		subject.ApplyImpulse(testValue.x, testValue.y)

		assert.Equal(t, testValue.x, subject.XVelocity)
		assert.Equal(t, testValue.y, subject.YVelocity)
	}

}

func TestIsColliding(t *testing.T) {
	testImage1, _ := ebiten.NewImage(2, 2, ebiten.FilterDefault)
	testImage1.Set(1, 1, color.RGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff})

	testImage2, _ := ebiten.NewImage(2, 2, ebiten.FilterDefault)
	testImage2.Set(0, 0, color.RGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff})

	subject := &spriteutils.Sprite{
		Image: testImage1,
	}

	otherSprite := &spriteutils.Sprite{
		Image: testImage2,
	}

	var testValues = []struct {
		x1, y1, x2, y2 int
		expected bool
	}{
		{0, 0, 0, 0, false},
		{0, 0, 5, 5, false},
		{0, 0, 1, 1, true},
	}

	for _, testValue := range testValues {
		subject.X = testValue.x1
		subject.Y = testValue.y1
		otherSprite.X = testValue.x2
		otherSprite.Y = testValue.y2

		var result = subject.IsColliding(otherSprite)

		assert.Equal(t, testValue.expected, result)
	}


}
