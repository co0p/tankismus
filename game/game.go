package game

import (
	"time"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/co0p/tankismus/game/scenes/start"
	"github.com/co0p/tankismus/pkg/scene"
)

// Game implements ebiten.Game and delegates to a scene manager.
type Game struct {
	manager  *scene.Manager
	lastTime time.Time
}

// NewGame constructs a new Game wired to the start scene.
func NewGame() *Game {
	g := &Game{}
	// manager is initialized with nil, then StartScene will set itself.
	m := scene.NewManager(nil)
	startScene := start.New(m)
	m.SetScene(startScene)
	g.manager = m
	g.lastTime = time.Now()
	return g
}

func (g *Game) Update() error {
	now := time.Now()
	dt := now.Sub(g.lastTime).Seconds()
	g.lastTime = now

	g.manager.Update(dt)
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.manager.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	// For now use the window size directly.
	return outsideWidth, outsideHeight
}
