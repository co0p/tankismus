package run

import (
	"image/color"
	"time"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/co0p/tankismus/game/assets"
	"github.com/co0p/tankismus/game/components"
	"github.com/co0p/tankismus/game/systems"
	"github.com/co0p/tankismus/pkg/ecs"
	"github.com/co0p/tankismus/pkg/input"
)

// Scene represents the main gameplay scene.
type Scene struct {
	world      *ecs.World
	player     ecs.EntityID
	lastUpdate time.Time
}

// New constructs a new run scene with a single player tank.
func New(_ interface{}) *Scene {
	w := ecs.NewWorld()
	player := w.NewEntity()
	w.AddComponent(player, &components.Transform{X: 100, Y: 100, Rotation: 0, Scale: 1})
	w.AddComponent(player, &components.Velocity{})
	w.AddComponent(player, &components.Sprite{SpriteID: "player_tank"})

	return &Scene{
		world:      w,
		player:     player,
		lastUpdate: time.Now(),
	}
}

func (s *Scene) OnEnter() {}

func (s *Scene) OnExit() {}

func (s *Scene) Update(dt float64) {
	_ = dt
	// Compute dt based on real time for now.
	now := time.Now()
	elapsed := now.Sub(s.lastUpdate).Seconds()
	s.lastUpdate = now

	input.Poll()
	systems.InputMovementSystem(s.world, s.player)
	systems.MovementSystem(s.world, elapsed)
}

func (s *Scene) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{R: 10, G: 40, B: 10, A: 255})

	// ensure assets are loaded; Load is idempotent.
	_ = assets.Load()

	systems.RenderSystem(s.world, screen)
}
