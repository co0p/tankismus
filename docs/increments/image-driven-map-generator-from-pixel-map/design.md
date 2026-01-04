# Design: Image-Driven Map Generator From Pixel Map

## Context and Problem

This design describes how to implement the increment defined in
`increment.md` in this folder: replacing the seeded grass-only map generator
with an image-driven generator that derives tilemaps from a pixel image.

Today:

- The map data model lives in an engine-style package as `pkg/map.Map`
  (Width, Height, Seed, Tiles[][]string) and is JSON-serializable.
- The `cmd/genesis` CLI takes a seed, width, and height, calls
  `pkg/map.NewGrassMap(seed, w, h)`, and writes the resulting map as JSON.
- The run scene and assets packages can already compose and render a tilemap
  from a `Map` instance, but level layout is random grass only.

Problem and goal (from the increment):

- Developers cannot intentionally author a specific level layout without
  touching code; they are limited to grass-only random maps.
- The goal is to make a pixel map (for example, `game/assets/maps/map.png`)
  the source of truth: fixed RGB colors (green, yellow, black) are mapped to
  terrain types (grass, sand, road), then to tile IDs, producing a JSON map in
  the existing `Map` format. The seeded grass generator is removed from
  normal workflows.

This design must respect the project’s `constitution-mode: lite` (focused,
pragmatic design and safety net), the layering rules in `ARCHITECTURE.md`, and
keep engine-style packages free from game-specific concerns.

References:

- Project constitution: `CONSTITUTION.md`.
- Increment definition: `increment.md` in this folder.
- Architecture overview: `ARCHITECTURE.md`.
- Existing map data model and generator: `pkg/map` and `cmd/genesis`.

## Proposed Solution (Technical Overview)

At a high level, the solution introduces an image-based map generator in the
**game layer** and reorients the `cmd/genesis` CLI around it.

Key ideas:

- **Image-driven generation**
  - The generator loads a PNG pixel map using the Go standard library
    (`image`, `image/png`, `image/color`).
  - Pixels with exact, fixed RGB values are treated as terrain hints:
    - Green → grass.
    - Yellow → sand.
    - Black → road.
  - The generator first converts the image into a grid of terrain kinds,
    then into tile IDs.

- **Neighbor-aware tile selection**
  - For each cell, the generator examines its neighbors (north, south, east,
    west, and optionally diagonals) to decide which subtile to assign:
    - Grass and sand interiors use base tiles such as `tileGrass1/2` and
      `tileSand1/2`.
    - Roads use road variants (`*roadNorth`, `*roadEast`, `*roadCorner*`,
      `*roadSplit*`, `*roadCrossing*`) based on the presence of road
      neighbors.
    - Terrain boundaries between grass and sand (and where appropriate,
      between terrain and road) use existing transition tiles such as
      `tileGrass_transition*` and `*_roadTransition*` when available.

- **Reuse of existing map model**
  - The generator constructs a `*pkg/map.Map` with:
    - `Width` and `Height` derived from the image dimensions.
    - `Seed` set to a conventional, non-random value (for example, `0`) and
      not used for behavior.
    - `Tiles` populated with tile IDs that correspond to sprites under
      `game/assets/images`.
  - JSON encoding and downstream map consumers (e.g., the run scene) remain
    unchanged.

- **CLI: updated `cmd/genesis`**
  - The CLI interface is changed to accept an input image path and output
    JSON path, instead of seed/width/height.
  - It delegates to the new generator module:
    - Loads the PNG.
    - Passes it to the generator to obtain a `Map`.
    - Encodes the map as JSON at the requested path.
  - The seeded grass generator is no longer part of the normal `genesis`
    flow.

After this change, the typical flow is:

1. Developer edits `game/assets/maps/map.png` using a pixel editor.
2. Developer runs `genesis` with the input map and output JSON path.
3. Genesis calls the generator, which returns a `Map` derived from the image.
4. Genesis writes the JSON to disk.
5. The game loads the JSON and renders a tilemap that visually matches the
   image-driven layout.

## Scope and Non-Scope (Technical)

In scope:

- Parsing a PNG pixel map and converting it into a terrain grid using fixed
  RGB mappings for grass, sand, and road.
- Neighbor-aware tile selection for road shapes and terrain boundaries,
  using the existing tile set.
- Constructing a `*pkg/map.Map` populated with tile IDs and encoding it as
  JSON in the current schema.
- Updating `cmd/genesis` to:
  - Accept an input image and output JSON path.
  - Use the new generator instead of `NewGrassMap`.
- Keeping automated tests and CI in line with `go test ./...`.

Out of scope:

- Introducing new terrain types (for example, water or obstacles) or new tile
  sets.
- Changing the `Map` struct’s public fields or JSON schema beyond potentially
  redefining the semantics of `Seed` as metadata.
- Altering how scenes load or render maps, beyond pointing them at the
  generated JSON.
- Adding in-game map editors or advanced tooling.
- Terrain-aware gameplay changes (movement speeds, collision rules, AI) that
  depend on the new terrain information.

This design is a targeted change to where maps come from, not to how they are
used at runtime.

## Architecture and Boundaries

This section situates the new generator in the existing layered architecture.

- **Engine-style layer (`pkg/map`)**
  - Remains responsible for the `Map` data model and functions that operate
    on it in an engine-agnostic way (validation, tile lookups, world-space
    queries).
  - Continues to have no knowledge of Ebiten, assets, or tile naming.

- **Game layer (new generator module)**
  - A new game-level module owns:
    - Mapping from fixed RGB pixel values to logical terrain kinds
      (grass, sand, road).
    - Neighbor-aware classification rules for roads and terrain boundaries.
    - Mapping from terrain and neighbor configuration to tile ID strings
      that match sprite IDs under `game/assets/images`.
    - Construction of a `*pkg/map.Map` from an `image.Image`.
  - This module depends on:
    - The Go standard library for image decoding.
    - `pkg/map` for the data model.
  - It does not depend on Ebiten or on scene/system packages.

- **Applications layer (`cmd/genesis`)**
  - Continues to act as a thin CLI boundary:
    - Responsible for argument parsing and input/output paths.
    - Invokes the generator module to transform an image into a `Map`.
    - Encodes and writes JSON to disk.
  - No direct knowledge of tile IDs or terrain rules.

The layering continues to follow the rules from `ARCHITECTURE.md`: binaries
may depend on game and engine packages, game packages may depend on engine
packages, and engine packages remain game-agnostic.

A conceptual component-level view for this increment is:

```mermaid
%% C4 Component-level view focused on map generation
graph TD
    subgraph Applications
        GenesisCLI[cmd/genesis
        CLI]
    end

    subgraph GameLayer[Game Layer]
        ImgGen[Image-Driven Map Generator
        (game module)]
    end

    subgraph Engine[Engine-Style Packages]
        MapModel[pkg/map
        Map Model]
    end

    subgraph StdLib[Standard Library]
        ImagePkg[image/png, image/color]
        IOPkg[os, filepath, encoding/json]
    end

    GenesisCLI --> ImgGen
    ImgGen --> MapModel
    ImgGen --> ImagePkg
    GenesisCLI --> IOPkg
```

## Contracts and Data

### Map JSON Contract

The JSON map contract remains that of `pkg/map.Map`. For reference, a
simplified JSON Schema fragment for maps after this change is:

```json
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "$id": "https://example.com/schemas/tankismus/map.json",
  "title": "Tankismus Tile Map",
  "type": "object",
  "required": ["width", "height", "seed", "tiles"],
  "properties": {
    "width": { "type": "integer", "minimum": 1 },
    "height": { "type": "integer", "minimum": 1 },
    "seed": {
      "type": "integer",
      "description": "Metadata for map generation; not used for randomness in image-driven maps"
    },
    "tiles": {
      "type": "array",
      "minItems": 1,
      "items": {
        "type": "array",
        "minItems": 1,
        "items": { "type": "string" }
      }
    }
  }
}
```

The semantics of `seed` change for image-driven maps: it no longer drives
random generation and may be set to a conventional value (for example, `0` or
an image-derived checksum) to preserve the field while leaving behavior
unchanged.

### Generator API (Internal Contract)

The generator module exposes an internal, Go-level API that can be used by
`cmd/genesis` and potentially by tests or other tools. Conceptually:

- **Purpose**: convert an `image.Image` into a populated `*pkg/map.Map` using
  fixed RGB and neighbor-aware tile rules.
- **Behavior**:
  - Validates the image (non-zero dimensions, rectangular grid).
  - Converts pixel colors to terrain kinds, failing with an error if any
    unsupported color is encountered.
  - Classifies each cell’s neighbors and selects a tile ID from the existing
    tile naming scheme.
  - Constructs and returns a `Map` that passes `Map.validate()`.

### Color Mapping Contract

- Fixed RGB values (one per terrain type) are used for this increment. For
  example (illustrative):
  - Grass: `#00FF00`.
  - Sand: `#FFFF00`.
  - Road: `#000000`.
- Any pixel not matching one of these exact RGB values is treated as an
  unsupported color and causes the generator to return an error with a clear
  message.

### Neighbor and Tile Selection

- The generator maintains an internal grid of terrain kinds that mirrors the
  image dimensions.
- For each cell, it determines a neighbor bitmask or pattern, using at least
  the four cardinal directions:
  - A road cell with neighbors north and south but not east/west maps to a
    vertical road tile (e.g., `*_roadNorth`).
  - A road cell with neighbors east and west but not north/south maps to a
    horizontal road tile (e.g., `*_roadEast`).
  - A road cell with neighbors north and east only maps to a corner tile
    (e.g., `*_roadCornerUR`), and so on.
  - Cells with three road neighbors become T-junctions (e.g., `*_roadSplit*`).
  - Cells with four road neighbors become crossings.
- Terrain boundary cells (grass next to sand, or terrain next to road where a
  transition tile exists) use available transition tiles; interior terrain
  falls back to grass/sand base tiles.
- Where the tile set does not support a specific pattern, the generator can
  fall back to a reasonable approximation (e.g., treat certain diagonal-only
  connections as no connection) but must do so consistently.

These rules remain internal to the generator module and may be refined in
follow-up increments without changing the external JSON contract.

## Testing and Safety Net

Given `constitution-mode: lite`, the safety net focuses on essential
blackbox tests around the generator behavior and CLI.

### Unit-Level Tests (Generator Module)

- **Color Mapping**
  - Create tiny in-memory images (for example, 2×2 or 3×3) with known pixel
    colors and verify that:
    - Grass, sand, and road pixels are converted to the correct terrain kind.
    - Unsupported colors produce clear errors.

- **Neighbor-Aware Road Shapes**
  - Build small terrain grids that represent:
    - Straight vertical and horizontal roads.
    - Corners (e.g., north+east, east+south, south+west, west+north).
    - T-junctions (three neighbors) and crossings (four neighbors).
  - Assert that the generator produces the expected tile IDs for these
    configurations.

- **Terrain Boundaries**
  - Simple grass–sand and terrain–road boundaries are tested to ensure
    transition tiles are used where expected and interior tiles are used
    elsewhere.

- **Map Structure Validity**
  - Generated maps should pass the existing `Map.validate()` constraints:
    - Width/Height > 0.
    - Tiles has the correct dimensions.
    - Seed is within allowed range, even if its semantics are different.

### CLI-Level Tests (`cmd/genesis`)

- **Happy Path**
  - Invoke the `run`-style entrypoint used by `main` with arguments pointing
    to a small test PNG and an output path.
  - Verify that it returns no error and writes a JSON file that decodes into
    a valid `Map` with expected dimensions and representative tile IDs.

- **Error Paths**
  - Invalid input path: verify that a descriptive error is returned.
  - Unsupported pixel color in the input image: verify that the error
    message explains the issue in human-readable terms.

All tests are invoked via `go test ./...` and should execute quickly.

## CI/CD and Rollout

- **CI Integration**
  - No changes are required to the CI pipeline beyond adding the new tests;
    `go test ./...` remains the primary check.
  - The generator module and CLI tests are simply new packages included in
    the existing test run.

- **Rollout**
  - The updated `cmd/genesis` usage becomes the new documented way to
    generate maps:
    - Input: path to a PNG pixel map.
    - Output: path to the JSON map file.
  - Developer documentation (for example, README snippets or tooling docs)
    should be aligned with the new behavior so that new contributors adopt
    the image-driven workflow.
  - The old seed/width/height interface is removed rather than preserved as
    an alternate path.

This is a forward-only change: once the new generator and CLI semantics are
in place and validated, development proceeds under the new behavior without a
built-in rollback path.

## Observability and Operations

The main observability needs for this increment are around generator usage
and failures.

- **Logging in `cmd/genesis`**
  - Log at least the following at an informational level:
    - Input image path and derived map dimensions.
    - Output path.
  - On errors, log clear messages including:
    - Whether the failure was due to I/O, unsupported colors, or
      classification issues.
    - Counts or sample coordinates for unsupported pixels.

- **Developer Experience**
  - CLI exit codes clearly indicate success (0) or failure (non-zero).
  - Error messages are actionable enough that developers can correct invalid
    images without needing to inspect the internals of the generator.

Given the project’s lite constitution mode, additional observability (metrics,
alerts, dashboards) is not required for this increment.

## Risks, Trade-offs, and Alternatives

- **Risks**
  - The tile set may not perfectly cover every neighbor configuration,
    leading to visually imperfect but still playable roads.
  - If the neighbor-classification logic becomes too clever, it may be
    difficult to reason about or extend; small misclassifications could be
    hard to debug.
  - Strict fixed RGB matching may surprise users if they slightly change
    colors in the pixel editor; they will see errors instead of graceful
    degradation.

- **Trade-offs**
  - The design deliberately keeps the generator logic in a game-layer module
    rather than in `pkg/map` to maintain a clean separation between
    engine-agnostic data structures and game-specific asset knowledge.
  - Neighbor matching focuses on cardinal directions and a manageable number
    of patterns, accepting that some corner cases may not look perfect in
    exchange for code simplicity.
  - Fixed RGB values make the system predictable and easy to test but place
    a small burden on map authors to adhere to exact colors.

- **Alternatives Considered**
  - **Generator in `pkg/map`**: rejected because it would couple the
    engine-style package to specific tile IDs and image formats, violating
    the architecture’s boundary.
  - **Configurable rules file for tile selection**: deferred as unnecessary
    complexity for a first increment; a hard-coded ruleset is sufficient and
    easier to refactor later.
  - **Retaining seeded grass generation as an alternate path**: rejected to
    align with the constitution’s "no long-lived fallbacks" principle.

## Follow-up Work

This design enables several natural follow-ups that are intentionally out of
scope for this increment:

- Extend the color and terrain mapping to support additional terrain types,
  such as water, obstacles, or destructible cover, with corresponding tile
  assets.
- Use the terrain grid for gameplay-affecting behavior (movement speeds,
  collision rules, AI pathfinding) based on terrain kind.
- Add lightweight developer tools to inspect or preview generated maps:
  - In-game debug overlays that show tile kinds under the cursor or tank.
  - A standalone viewer for map JSONs.
- Support multiple named maps or levels (e.g., by mapping multiple PNGs to
  distinct JSON files) and a way to select which map to load for a session.
- Revisit the `Map` data model and the role of `Seed` once more map-related
  features are in place, possibly simplifying the schema or adding metadata
  fields.

## References

- `CONSTITUTION.md` — project mode, layout, and delivery principles.
- `ARCHITECTURE.md` — layered architecture and dependency rules.
- `docs/increments/image-driven-map-generator-from-pixel-map/increment.md` —
  product-level definition of this increment.
- Existing map engine: `pkg/map`.
- Existing generator CLI: `cmd/genesis`.
- Tile assets and source map: `game/assets/images`, `game/assets/maps`.
