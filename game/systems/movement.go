package systems

import (
	"github.com/co0p/tankismus/game/components"
	"github.com/co0p/tankismus/pkg/ecs"
)

// MovementSystem applies velocity to transform for all entities
// that have both Transform and Velocity components.
func MovementSystem(world *ecs.World, dt float64) {
	required := ecs.MaskFor(components.TypeTransform, components.TypeVelocity)
	entities := world.Find(required)
	for _, id := range entities {
		cT, okT := world.GetComponent(id, components.TypeTransform)
		cV, okV := world.GetComponent(id, components.TypeVelocity)
		if !okT || !okV {
			continue
		}

		p, okP := cT.(*components.Transform)
		v, okVel := cV.(*components.Velocity)
		if !okP || !okVel {
			continue
		}

		p.X += v.VX * dt
		p.Y += v.VY * dt
		p.Rotation += v.Angular * dt
	}
}
