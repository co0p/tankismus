package run

import (
	"encoding/json"
	"image/color"
	"os"
	"time"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/co0p/tankismus/game/assets"
	"github.com/co0p/tankismus/game/components"
	"github.com/co0p/tankismus/game/systems"
	"github.com/co0p/tankismus/pkg/ecs"
	"github.com/co0p/tankismus/pkg/input"
	mappkg "github.com/co0p/tankismus/pkg/map"
)

// Scene represents the main gameplay scene.
type Scene struct {
	world      *ecs.World
	player     ecs.EntityID
	tilemap    ecs.EntityID
	levelMap   *mappkg.Map
	lastUpdate time.Time
}

// New constructs a new run scene with a single player tank.
// If ctx is a *mappkg.Map, it is used as the level map (primarily for tests).
// Otherwise, the scene attempts to load game/assets/maps/map.json. If loading
// or validation fails, no level map or tilemap is created.
func New(ctx interface{}) *Scene {
	w := ecs.NewWorld()
	// Ensure core assets, including tile sprites, are loaded before composing
	// the level tilemap. Load is idempotent.
	_ = assets.Load()

	// Determine the level map to use.
	var levelMap *mappkg.Map
	if m, ok := ctx.(*mappkg.Map); ok && m != nil {
		levelMap = m
	} else {
		// Attempt to load the default JSON map.
		const defaultMapPath = "game/assets/maps/map.json"
		file, err := os.Open(defaultMapPath)
		if err == nil {
			defer file.Close()
			dec := json.NewDecoder(file)
			var loaded mappkg.Map
			if err := dec.Decode(&loaded); err == nil {
				// Ensure the loaded map satisfies basic invariants.
				if err := loaded.ValidateForGenerator(); err == nil {
					levelMap = &loaded
				}
			}
		}
	}

	var tilemapEntity ecs.EntityID
	if levelMap != nil {
		const tileSize = 16
		// Compose the tilemap image and register it in the assets registry.
		if img, err := assets.ComposeTilemap("tilemap_ground", levelMap, tileSize); err == nil {
			wImg, hImg := img.Size()
			tilemapEntity = w.NewEntity()
			// Position the tilemap so that its top-left corner aligns with the
			// world origin. RenderSystem draws sprites centered on their
			// Transform, so we offset by half the tilemap dimensions.
			w.AddComponent(tilemapEntity, &components.Transform{X: float64(wImg) / 2, Y: float64(hImg) / 2, Rotation: 0, Scale: 1})
			w.AddComponent(tilemapEntity, &components.Sprite{SpriteID: "tilemap_ground"})
			w.AddComponent(tilemapEntity, &components.RenderOrder{Z: 0})
		}
	}

	player := w.NewEntity()
	w.AddComponent(player, &components.Transform{X: 100, Y: 100, Rotation: 0, Scale: 1})
	w.AddComponent(player, &components.Velocity{})
	w.AddComponent(player, &components.ControlIntent{})
	w.AddComponent(player, &components.MovementParams{
		MaxForwardSpeed:     133.3333,
		MaxBackwardSpeed:    80,
		LinearAcceleration:  200,
		LinearDeceleration:  300,
		MaxTurnRate:         3,
		AngularAcceleration: 6,
		AngularDeceleration: 9,
	})
	w.AddComponent(player, &components.Sprite{SpriteID: "player_tank"})
	w.AddComponent(player, &components.RenderOrder{Z: 10})

	return &Scene{
		world:      w,
		player:     player,
		tilemap:    tilemapEntity,
		levelMap:   levelMap,
		lastUpdate: time.Now(),
	}
}

func (s *Scene) OnEnter() {}

func (s *Scene) OnExit() {}

func (s *Scene) Update(dt float64) {
	input.Poll()
	systems.InputMovementSystem(s.world, s.player)
	systems.MovementSystem(s.world, dt)
}

func (s *Scene) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{R: 10, G: 40, B: 10, A: 255})

	// ensure assets are loaded; Load is idempotent.
	_ = assets.Load()

	systems.RenderSystem(s.world, screen)
}

// World exposes the underlying ECS world for testing purposes.
func (s *Scene) World() *ecs.World {
	return s.world
}

// Player returns the player entity ID for testing purposes.
func (s *Scene) Player() ecs.EntityID {
	return s.player
}
