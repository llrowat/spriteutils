package spriteutils_test

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/llrowat/spriteutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

type MockSprite struct{
	mock.Mock
}

func(m *MockSprite) Draw(screen *ebiten.Image) error {
	m.Called(screen)
	return nil
}

func(m *MockSprite) Update() {
	m.Called()
}

func TestTransientSprite_Update(t *testing.T) {
	var mockSprite = new(MockSprite)
	mockSprite.On("Update").Return()

	var subject = &spriteutils.TransientSprite{
		CreatedAtGameTime: time.Second*0,
		LifetimeDuration: time.Second*1,
		Sprite: mockSprite,
	}

	subject.Update(time.Second*0)
	mockSprite.AssertNumberOfCalls(t, "Update", 1)

	subject.Update(time.Second*2)
	mockSprite.AssertNumberOfCalls(t, "Update", 1)
	assert.Empty(t, subject.Sprite)
}

func TestTransientSprite_Draw(t *testing.T) {
	testScreen, _ := ebiten.NewImage(5, 5, ebiten.FilterDefault)

	var mockSprite = new(MockSprite)
	mockSprite.On("Draw", testScreen).Return()

	var subject = &spriteutils.TransientSprite{
		CreatedAtGameTime: time.Second*0,
		LifetimeDuration: time.Second*1,
		Sprite: mockSprite,
	}

	subject.Draw(testScreen)
	mockSprite.AssertNumberOfCalls(t, "Draw", 1)

	subject.IsExpired = true
	subject.Draw(testScreen)
	mockSprite.AssertNumberOfCalls(t, "Draw", 1)

	subject.IsExpired = false
	subject.Sprite = nil
	subject.Draw(testScreen)
	mockSprite.AssertNumberOfCalls(t, "Draw", 1)
}