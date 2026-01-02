package systems

import (
	"github.com/co0p/tankismus/game/components"
	"github.com/co0p/tankismus/pkg/ecs"
	"github.com/co0p/tankismus/pkg/input"
)

// InputMovementSystem updates the player's control intent based on input
// actions. It no longer writes velocity directly; MovementSystem interprets
// the intent and updates velocity and transform.
func InputMovementSystem(world *ecs.World, player ecs.EntityID) {
	cI, okI := world.GetComponent(player, components.TypeControlIntent)
	if !okI {
		return
	}

	intent, okIntent := cI.(*components.ControlIntent)
	if !okIntent {
		return
	}

	// Throttle: forward/backward along facing direction.
	throttle := 0.0
	if input.IsActionDown(input.ActionMoveForward) {
		throttle += 1
	}
	if input.IsActionDown(input.ActionMoveBackward) {
		throttle -= 1
	}
	if throttle > 1 {
		throttle = 1
	}
	if throttle < -1 {
		throttle = -1
	}

	// Turn: left/right.
	turn := 0.0
	if input.IsActionDown(input.ActionTurnLeft) {
		turn -= 1
	}
	if input.IsActionDown(input.ActionTurnRight) {
		turn += 1
	}
	if turn > 1 {
		turn = 1
	}
	if turn < -1 {
		turn = -1
	}

	intent.Throttle = throttle
	intent.Turn = turn
}
