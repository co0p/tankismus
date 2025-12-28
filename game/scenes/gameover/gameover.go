package gameover

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"

	"github.com/co0p/tankismus/game/scenes/start"
	"github.com/co0p/tankismus/pkg/scene"
)

// Scene represents the game over scene.
type Scene struct {
	manager *scene.Manager
}

// New constructs a new game over scene.
func New(manager *scene.Manager) *Scene {
	return &Scene{manager: manager}
}

func (s *Scene) OnEnter() {}

func (s *Scene) OnExit() {}

func (s *Scene) Update(dt float64) {
	_ = dt
	if len(inpututil.PressedKeys()) > 0 {
		s.manager.SetScene(start.New(s.manager))
	}
}

func (s *Scene) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{R: 30, G: 0, B: 0, A: 255})
	ebitenutil.DebugPrint(screen, "Game Over\nPress any key to return to start")
}
