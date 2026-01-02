package systems

import (
	"github.com/hajimehoshi/ebiten/v2"

	"github.com/co0p/tankismus/game/assets"
	"github.com/co0p/tankismus/game/components"
	"github.com/co0p/tankismus/pkg/ecs"
)

// RenderSystem draws all entities that have a Transform and a Sprite component.
func RenderSystem(world *ecs.World, screen *ebiten.Image) {
	required := ecs.MaskFor(components.TypeTransform, components.TypeSprite)
	entities := world.Find(required)
	for _, id := range entities {
		cT, okT := world.GetComponent(id, components.TypeTransform)
		cS, okS := world.GetComponent(id, components.TypeSprite)
		if !okT || !okS {
			continue
		}

		p, okP := cT.(*components.Transform)
		s, okSprite := cS.(*components.Sprite)
		if !okP || !okSprite {
			continue
		}

		img := assets.GetSprite(s.SpriteID)
		if img == nil {
			continue
		}

		op := &ebiten.DrawImageOptions{}
		w, h := img.Size()
		cx := float64(w) / 2
		cy := float64(h) / 2

		// Move origin to sprite center, rotate by logical rotation, then move
		// to world position. Sprites are authored facing +X (to the right).
		op.GeoM.Translate(-cx, -cy)
		op.GeoM.Rotate(p.Rotation)
		op.GeoM.Translate(p.X, p.Y)
		screen.DrawImage(img, op)
	}
}
