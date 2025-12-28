package start

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"

	"github.com/co0p/tankismus/game/scenes/run"
	"github.com/co0p/tankismus/pkg/input"
	"github.com/co0p/tankismus/pkg/scene"
)

// Scene is the start scene showing a simple prompt.
type Scene struct {
	manager *scene.Manager
}

// New constructs a new start scene.
func New(manager *scene.Manager) *Scene {
	return &Scene{manager: manager}
}

func (s *Scene) OnEnter() {}

func (s *Scene) OnExit() {}

func (s *Scene) Update(dt float64) {
	_ = dt
	// Any key press starts the game.
	if input.AnyKeyPressed() {
		s.manager.SetScene(run.New(s.manager))
	}
}

func (s *Scene) Draw(screen *ebiten.Image) {
	screen.Fill(color.Black)
	ebitenutil.DebugPrint(screen, "tankismus\nPress any key to start")
}
