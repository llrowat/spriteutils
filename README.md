# spriteutils

## Overview
Some simple utilities for working with 2D sprites ontop of the [Ebiten](https://ebiten.org/) game library.  Core features:
- Simplified handling of sprite movement and rotation around image center point
- Pixel-perfect collision detection
- Ability to apply impulse to sprite
- Easily create sprites that exist on screen for a short period of time
- A factory that can create sprites in random locations within its set constaints

I created this module for use in [github.com/llrowat/galactic-asteroid-belt](https://github.com/llrowat/galactic-asteroid-belt) which was a simple game project I worked on to teach myself Go.  If anyone else gets some usefulness from it, cool!  

## Installing
To install the library, you must have Go 1.17 installed.  From the project that will use this library just run:

```go get github.com/llrowat/spriteutils```

## Usage Examples

### Sprite

Create a Sprite from an (empty in this example) ebiten image at position (100, 20)
```
exampleSprite = &spriteutils.Sprite{
  Image:     ebiten.NewImage(2, 2, ebiten.FilterDefault),
  X:         100,
  Y:         20,
  XVelocity: 0,
  YVelocity: 0,
  Rotation:  0,
}
```

To update the sprite movement based on velocity, you would add the following to your Game update loop
```
exampleSprite.Update()
```

To draw the sprite image based on position and rotation, you would add the following to your Game draw method, where ```screen``` is the ebiten image representing the game canvas
```
err := exampleSprite.Draw(screen)
if err != nil {
    log.Fatal(err)
}
```

To apply an impulse force to the sprite, for example to increase x velocity by 5 pixels/frame and y velocity by 10 pixels/frame
```
exampleSprite.ApplyImpulse(5, 10)
```

To determine whether there has been a collision between two sprites.  A collision is occuring if the sprites are overlapping and there exists a non-transparent pixel at one of the overlapping pixels for both sprites.
```
isColliding := exampleSprite.IsColliding(someOtherSprite)
```

### SpriteFactory

To Create a new sprite factory with one of two (empty for this example) images in a 0 to 150 pixel square
```
exampleFactory = &spriteutils.SpriteFactory{
  Images: []*ebiten.Image{ ebiten.NewImage(1, 1, ebiten.FilterDefault), 
                           ebiten.NewImage(2, 2, ebiten.FilterDefault)},
  MinX:   0,
  MaxX:   150,
  MinY:   0,
  MaxY:   150,
}
```

To generate a new sprite with one of the images, and within the bounds (0, 0) to (150, 150)
```
newSprite := exampleFactory.GenerateSprite()
```

### TransientSprite

Transient sprites can be useful for single frame effects, such as an object death/destruction frame.  It is best if they are stored in slices so that you are not referencing them directly and can more easily remove refrence for garbage collection.  The ```CreatedAtGameTime``` should always be the game time (duration since game started).  You cannot use real world time reliably because when you unfocus the game window the game will pause.  You can get game time by incrementing a frame counter in the Game update loop.  Since Ebiten runs at 60 TPS we can calculate the game time using ```time.Duration(gameFrameCount) * time.Second / 60```

To create a transient sprite from some sprite object ```exampleSprite``` that will last for 5 seconds
```
transientSprites []*spriteutils.TransientSprite
newTransientSprite := &spriteutils.TransientSprite {
  CreatedAtGameTime: time.Duration(gameFrameCount) * time.Second / 60,
  LifetimeDuration:  time.Second * 5,
  Sprite: exampleSprite 
}

transientSprites = append(transientSprites, newTransientSprite)
```

Like a normal sprite, you would put the Update call in the Game update loop.  The only difference is you pass in the current game time
```
temp := transientSprites[:0]

for _, transientSprite := range transientSprites {
  transientSprite.Update(time.Duration(gameFrameCount) * time.Second / 60)
  
  // Only keep sprites that are not expired.  The expired sprites will be dereferenced for garbage collection
  if !transientSprite.IsExpired {
    temp = append(temp, transientSprite)
  }
  
}

transientSprites = temp
```

Like a normal sprite, you would put the Draw call in the Game draw method
```
for _, transientSprite := range transientSprites {
  transientSprite.Draw(screen)
}
```
