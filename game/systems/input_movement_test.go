package systems

import (
	"testing"

	"github.com/co0p/tankismus/game/components"
	"github.com/co0p/tankismus/pkg/ecs"
	"github.com/co0p/tankismus/pkg/input"
)

func newInputTestWorld() (*ecs.World, ecs.EntityID) {
	w := ecs.NewWorld()
	id := w.NewEntity()
	w.AddComponent(id, &components.Transform{})
	w.AddComponent(id, &components.Velocity{})
	w.AddComponent(id, &components.ControlIntent{})
	w.AddComponent(id, &components.MovementParams{})
	return w, id
}

// helper to clear all input state between tests
func clearInputState() {
	// Poll with no keys pressed will clear all states; this relies on
	// Ebiten state being initially empty when tests run headless.
	input.Poll()
}

func TestInputMovementSystem_SetsThrottleFromMoveKeys(t *testing.T) {
	w, id := newInputTestWorld()
	clearInputState()

	// With no keys pressed, intent should remain neutral.
	InputMovementSystem(w, id)
	cI, _ := w.GetComponent(id, components.TypeControlIntent)
	intent := cI.(*components.ControlIntent)
	if intent.Throttle != 0 {
		t.Fatalf("expected neutral throttle with no input, got %v", intent.Throttle)
	}

	// Simulate pressing forward.
	input.SetActionStateForTest(input.ActionMoveForward, true)
	input.SetActionStateForTest(input.ActionMoveBackward, false)
	InputMovementSystem(w, id)
	if intent.Throttle != 1 {
		t.Fatalf("expected throttle=1 when moving forward, got %v", intent.Throttle)
	}

	// Simulate pressing backward.
	input.SetActionStateForTest(input.ActionMoveForward, false)
	input.SetActionStateForTest(input.ActionMoveBackward, true)
	InputMovementSystem(w, id)
	if intent.Throttle != -1 {
		t.Fatalf("expected throttle=-1 when moving backward, got %v", intent.Throttle)
	}

	// No move keys.
	input.SetActionStateForTest(input.ActionMoveForward, false)
	input.SetActionStateForTest(input.ActionMoveBackward, false)
	InputMovementSystem(w, id)
	if intent.Throttle != 0 {
		t.Fatalf("expected throttle=0 when no move keys pressed, got %v", intent.Throttle)
	}
}

func TestInputMovementSystem_SetsTurnFromTurnKeys(t *testing.T) {
	w, id := newInputTestWorld()
	clearInputState()

	InputMovementSystem(w, id)
	cI, _ := w.GetComponent(id, components.TypeControlIntent)
	intent := cI.(*components.ControlIntent)
	if intent.Turn != 0 {
		t.Fatalf("expected neutral turn with no input, got %v", intent.Turn)
	}

	// Turn left.
	input.SetActionStateForTest(input.ActionTurnLeft, true)
	input.SetActionStateForTest(input.ActionTurnRight, false)
	InputMovementSystem(w, id)
	if intent.Turn != -1 {
		t.Fatalf("expected turn=-1 when turning left, got %v", intent.Turn)
	}

	// Turn right.
	input.SetActionStateForTest(input.ActionTurnLeft, false)
	input.SetActionStateForTest(input.ActionTurnRight, true)
	InputMovementSystem(w, id)
	if intent.Turn != 1 {
		t.Fatalf("expected turn=1 when turning right, got %v", intent.Turn)
	}

	// No turn keys.
	input.SetActionStateForTest(input.ActionTurnLeft, false)
	input.SetActionStateForTest(input.ActionTurnRight, false)
	InputMovementSystem(w, id)
	if intent.Turn != 0 {
		t.Fatalf("expected turn=0 when no turn keys pressed, got %v", intent.Turn)
	}
}

func TestInputMovementSystem_DoesNotModifyVelocityDirectly(t *testing.T) {
	w, id := newInputTestWorld()
	clearInputState()

	// Pre-set some non-zero velocity and ensure it is not overwritten.
	cV, _ := w.GetComponent(id, components.TypeVelocity)
	v := cV.(*components.Velocity)
	v.VX = 10
	v.VY = -5
	v.Angular = 1.5

	input.SetActionStateForTest(input.ActionMoveForward, true)
	input.SetActionStateForTest(input.ActionTurnRight, true)
	InputMovementSystem(w, id)

	if v.VX != 10 || v.VY != -5 || v.Angular != 1.5 {
		t.Fatalf("expected velocity unchanged by input system, got vx=%v vy=%v ang=%v", v.VX, v.VY, v.Angular)
	}
}
