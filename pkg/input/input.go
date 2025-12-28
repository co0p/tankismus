package input

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// Action represents a high-level input action like moving or shooting.
type Action string

const (
	ActionMoveForward  Action = "move_forward"
	ActionMoveBackward Action = "move_backward"
	ActionTurnLeft     Action = "turn_left"
	ActionTurnRight    Action = "turn_right"
	ActionFire         Action = "fire"
)

// mapping from actions to ebiten keys.
var actionBindings = map[Action][]ebiten.Key{
	ActionMoveForward:  {ebiten.KeyW},
	ActionMoveBackward: {ebiten.KeyS},
	ActionTurnLeft:     {ebiten.KeyA},
	ActionTurnRight:    {ebiten.KeyD},
	ActionFire:         {ebiten.KeySpace},
}

// state caches the current frame's action state.
var state = map[Action]bool{}

// Poll updates the cached action state from Ebiten's keyboard state.
func Poll() {
	for action, keys := range actionBindings {
		pressed := false
		for _, k := range keys {
			if ebiten.IsKeyPressed(k) {
				pressed = true
				break
			}
		}
		state[action] = pressed
	}
}

// IsActionDown reports whether the given action is currently active.
func IsActionDown(a Action) bool {
	return state[a]
}

// AnyKeyPressed reports whether any key was pressed in the current frame.
// This is useful for simple "press any key" screens while still keeping
// Ebiten-specific details inside the input package.
func AnyKeyPressed() bool {
	return len(inpututil.PressedKeys()) > 0
}
