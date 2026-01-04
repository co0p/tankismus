package assets

import (
	"embed"
	"errors"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"

	_ "image/png"

	mappkg "github.com/co0p/tankismus/pkg/map"
)

//go:embed images/*
var imagesFS embed.FS

// Registry maps sprite IDs to loaded Ebiten images.
var Registry = map[string]*ebiten.Image{}

// ErrTileSpriteNotFound is returned by ComposeTilemap when a tile ID in the
// map does not have a corresponding sprite registered in the assets registry.
var ErrTileSpriteNotFound = errors.New("assets: tile sprite not found")

// Load loads all core assets into the registry.
// It is safe to call multiple times; later calls will simply overwrite
// existing entries with the same IDs.
func Load() error {
	entries, err := imagesFS.ReadDir("images")
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if !strings.HasSuffix(name, ".png") {
			continue
		}

		img, _, err := ebitenutil.NewImageFromFileSystem(imagesFS, "images/"+name)
		if err != nil {
			// Skip images that fail to load; the game can still run,
			// and ComposeTilemap will surface missing sprites as needed.
			continue
		}

		id := strings.TrimSuffix(name, ".png")
		if id == "tank" {
			id = "player_tank"
		}
		Registry[id] = img
	}

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

// ComposeTilemap creates a new Ebiten image by drawing all tiles from the
// provided map into a single image and registers it under the given spriteID.
//
// It expects that each tile ID used in the map (for example, "tileGrass1",
// "tileGrass2") is already present in the Registry and that all tile sprites
// share the same square dimensions given by tileSize. If any tile sprite is
// missing, ComposeTilemap returns an error.
func ComposeTilemap(spriteID string, m *mappkg.Map, tileSize int) (*ebiten.Image, error) {
	if m == nil {
		return nil, nil
	}
	if tileSize <= 0 {
		return nil, nil
	}

	width := m.Width * tileSize
	height := m.Height * tileSize
	img := ebiten.NewImage(width, height)

	for y := 0; y < m.Height; y++ {
		for x := 0; x < m.Width; x++ {
			tileID, ok := m.TileAt(x, y)
			if !ok {
				continue
			}
			src := GetSprite(tileID)
			if src == nil {
				return nil, ErrTileSpriteNotFound
			}

			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(x*tileSize), float64(y*tileSize))
			img.DrawImage(src, op)
		}
	}

	Registry[spriteID] = img
	return img, nil
}
