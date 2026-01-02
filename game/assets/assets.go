package assets

import (
	"embed"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

//go:embed images/*
var imagesFS embed.FS

// Registry maps sprite IDs to loaded Ebiten images.
var Registry = map[string]*ebiten.Image{}

// Load loads all core assets into the registry.
// For now we load a single placeholder tank sprite if present.
func Load() error {
	// Example: try to load images/tank.png if it exists.
	img, _, err := ebitenutil.NewImageFromFileSystem(imagesFS, "images/tank.png")
	if err == nil {
		Registry["player_tank"] = img
	}
	// It's okay if the image is missing; the game can still run.
	return nil
}

// GetSprite returns the Ebiten image for a sprite ID, if loaded.
func GetSprite(id string) *ebiten.Image {
	return Registry[id]
}

// RegisterSpriteForTest allows tests to inject sprites into the registry
// without loading from disk.
func RegisterSpriteForTest(id string, img *ebiten.Image) {
	Registry[id] = img
}
