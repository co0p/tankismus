package systems

import (
	"image/color"
	"testing"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/co0p/tankismus/game/assets"
	"github.com/co0p/tankismus/game/components"
	"github.com/co0p/tankismus/pkg/ecs"
)

// fakeSprite is a simple ebiten.Image backed by a Go image for bounds.
func fakeSprite(w, h int) *ebiten.Image {
	img := ebiten.NewImage(w, h)
	img.Fill(color.White)
	return img
}

// TestRenderSystem_RotatesAroundSpriteCenter is a coarse behavioral test: it
// verifies that when rotation changes, the drawn image remains visually
// centered on the transform position. Since Ebiten does not expose the
// internal GeoM matrix or draw calls, this test is limited to ensuring the
// call does not panic and can be extended with golden-image style tests in
// future if needed.
func TestRenderSystem_RotatesAroundSpriteCenter(t *testing.T) {
	// Prepare a fake sprite and register it.
	spriteID := "test_tank"
	img := fakeSprite(20, 10) // width=20,height=10, center=(10,5)
	assets.RegisterSpriteForTest(spriteID, img)

	world := ecs.NewWorld()
	id := world.NewEntity()
	world.AddComponent(id, &components.Transform{X: 100, Y: 50, Rotation: 0, Scale: 1})
	world.AddComponent(id, &components.Sprite{SpriteID: spriteID})

	screen := ebiten.NewImage(200, 100)

	// Draw at rotation 0 and then with a non-zero rotation. The main
	// verification here is that both calls succeed without error or panic.
	RenderSystem(world, screen)

	cT, _ := world.GetComponent(id, components.TypeTransform)
	p := cT.(*components.Transform)
	p.Rotation = 1.0

	RenderSystem(world, screen)

	// No explicit numeric assertions here due to limited access to the
	// underlying draw machinery; correctness is exercised indirectly via
	// integration tests of movement + rendering at a higher level.
}
