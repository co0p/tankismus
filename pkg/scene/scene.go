package scene

import "github.com/hajimehoshi/ebiten/v2"

// Scene represents a high level game state such as start menu, gameplay, or game over.
type Scene interface {
	OnEnter()
	OnExit()
	Update(dt float64)
	Draw(screen *ebiten.Image)
}

// Manager tracks the currently active scene.
type Manager struct {
	current Scene
}

// NewManager constructs a new scene manager.
func NewManager(initial Scene) *Manager {
	m := &Manager{}
	m.SetScene(initial)
	return m
}

// SetScene switches the active scene, calling lifecycle hooks.
func (m *Manager) SetScene(next Scene) {
	if m.current != nil {
		m.current.OnExit()
	}
	m.current = next
	if m.current != nil {
		m.current.OnEnter()
	}
}

// Update forwards the update call to the active scene.
func (m *Manager) Update(dt float64) {
	if m.current == nil {
		return
	}
	m.current.Update(dt)
}

// Draw forwards the draw call to the active scene.
func (m *Manager) Draw(screen *ebiten.Image) {
	if m.current == nil {
		return
	}
	m.current.Draw(screen)
}
