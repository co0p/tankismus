package systems

import (
	"math"

	"github.com/co0p/tankismus/game/components"
	"github.com/co0p/tankismus/pkg/ecs"
	"github.com/co0p/tankismus/pkg/input"
)

// InputMovementSystem updates the player's velocity based on input actions.
func InputMovementSystem(world *ecs.World, player ecs.EntityID) {
	cT, okT := world.GetComponent(player, components.TypeTransform)
	cV, okV := world.GetComponent(player, components.TypeVelocity)
	if !okT || !okV {
		return
	}

	p, okP := cT.(*components.Transform)
	v, okVel := cV.(*components.Velocity)
	if !okP || !okVel {
		return
	}

	const moveSpeed = 100.0
	const turnSpeed = 2.5 // radians per second

	// Turning
	if input.IsActionDown(input.ActionTurnLeft) {
		v.Angular = -turnSpeed
	} else if input.IsActionDown(input.ActionTurnRight) {
		v.Angular = turnSpeed
	} else {
		v.Angular = 0
	}

	// Forward/backward movement along facing direction
	forward := 0.0
	if input.IsActionDown(input.ActionMoveForward) {
		forward += 1
	}
	if input.IsActionDown(input.ActionMoveBackward) {
		forward -= 1
	}

	if forward != 0 {
		vx := math.Cos(p.Rotation) * moveSpeed * forward
		vy := math.Sin(p.Rotation) * moveSpeed * forward
		v.VX = vx
		v.VY = vy
	} else {
		v.VX = 0
		v.VY = 0
	}
}
