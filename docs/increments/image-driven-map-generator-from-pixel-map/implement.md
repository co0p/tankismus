# Implement: Image-Driven Map Generator From Pixel Map

## 1. Context

- Goal: Replace the seeded grass-only generator with an image-driven generator that reads a fixed-color PNG, maps pixels to grass/sand/road, chooses neighbor-aware subtiles, and outputs the existing `pkg/map.Map` JSON format.
- Design: Introduce a generator module in the game layer (`game/maps`) that converts `image.Image` → terrain grid → tile IDs → `*mappkg.Map`, and rewire `cmd/genesis` to use it (PNG in, JSON out) with no seed/width/height arguments.
- Constraints: `constitution-mode: lite` (keep plan and tests pragmatic), preserve `pkg/map.Map` JSON schema, keep `pkg/*` game-agnostic, and avoid long-lived fallbacks (forward-only change, no rollback path).
- References: `increment.md`, `design.md`, `CONSTITUTION.md`.

Status: Not started  
Next step: Step 1 – Scaffold game/maps generator and strict color mapping

## 2. Workstreams

- **Workstream A – Game-layer image generator (game/maps)**
- **Workstream B – Genesis CLI rewire (cmd/genesis)**
- **Workstream C – Safety net, docs, and cleanup**

## 3. Steps

- [ ] Step 1: Scaffold game/maps generator and strict color mapping
- [ ] Step 2: Implement straight-road neighbor classification
- [ ] Step 3: Implement corners, T-junctions, crossings, and basic transitions
- [ ] Step 4: Rewire cmd/genesis to use the image-driven generator
- [ ] Step 5: Align docs and usage with image-driven generation
- [ ] Step 6: Add a small integration-style tilemap composition check

---

### Step 1: Scaffold game/maps generator and strict color mapping

- **Workstream:** A – Game-layer image generator (game/maps)
- **Based on Design:**
  - `design.md` – Proposed Solution (Technical Overview: Image-driven generation)
  - `design.md` – Architecture and Boundaries (generator in game layer, depends on `pkg/map` and stdlib image packages)
  - `design.md` – Contracts and Data (Color Mapping Contract, Generator API)
- **Files:**
  - `game/maps/generator.go` (new)
  - `game/maps/generator_test.go` (new)
- **TDD Cycle:**
  - **Red – Failing test first:**
    - In `game/maps/generator_test.go`, create small in-memory images (e.g., 2×2) with `image.NewRGBA` and set pixels to the exact RGB values chosen for grass, sand, and road.
    - Add tests that call a new function `GenerateFromImage(img image.Image) (*mappkg.Map, error)` and assert:
      - It returns a non-nil `*mappkg.Map` with `Width`/`Height` matching the image bounds.
      - `Tiles` is a `[][]string` with dimensions `Height`×`Width`.
      - For a simple all-grass/all-sand/all-road image, each tile cell contains a non-empty tile ID string (e.g., placeholder IDs like `"tileGrass1"`, `"tileSand1"`, `"tileGrass_roadNorth"`).
    - Add a test where one pixel is set to an unsupported color (not one of the fixed RGB values) and assert that `GenerateFromImage` returns a non-nil error containing an indication of an unsupported color.
  - **Green – Make the test(s) pass:**
    - Implement `game/maps/generator.go` with:
      - A `TerrainKind` internal enum (e.g., grass, sand, road, unknown).
      - Fixed RGB constants for grass, sand, and road (exact values agreed for `map.png`).
      - A helper `colorToTerrain(c color.Color) (TerrainKind, error)` that strictly matches RGB to terrain kinds and returns an error for unsupported colors.
      - `GenerateFromImage(img image.Image) (*mappkg.Map, error)` that:
        - Validates non-zero image bounds.
        - Builds a `[][]TerrainKind` grid from the image using `colorToTerrain`.
        - Constructs a `mappkg.Map` with `Width`/`Height` from the image, `Seed` set to a fixed valid value (e.g., `0`), and a `Tiles` matrix sized `Height`×`Width`.
        - Initially maps each terrain kind to a simple placeholder tile ID (e.g., grass → `"tileGrass1"`, sand → `"tileSand1"`, road → a simple road tile ID) without neighbor-based refinement yet.
        - Ensures `Map.validate()` passes for the generated map.
  - **Refactor – Clean up with tests green:**
    - Extract small helper functions (`colorToTerrain`, `newMapFromTerrainGrid`) to keep `GenerateFromImage` readable.
    - Clarify comments on fixed RGB matching and temporary placeholder tile IDs.
- **CI / Checks:**
  - Run `go test ./game/maps` to verify new tests and implementation.
  - Run `go test ./...` to ensure no regressions elsewhere.

---

### Step 2: Implement straight-road neighbor classification

- **Workstream:** A – Game-layer image generator (game/maps)
- **Based on Design:**
  - `design.md` – Contracts and Data (Neighbor and Tile Selection: straight segments)
  - `design.md` – Testing and Safety Net (Neighbor-Aware Road Shapes)
- **Files:**
  - `game/maps/generator.go`
  - `game/maps/generator_test.go`
- **TDD Cycle:**
  - **Red – Failing test first:**
    - In `generator_test.go`, add tests that build small images representing:
      - A vertical road: a column of road pixels with grass/sand around.
      - A horizontal road: a row of road pixels with grass/sand around.
    - For each scenario, call `GenerateFromImage` and assert that:
      - Road cells in vertical segments use the intended vertical road tile ID (e.g., `tileGrass_roadNorth` or `tileSand_roadNorth` depending on terrain underneath).
      - Road cells in horizontal segments use the intended horizontal road tile ID (e.g., `tileGrass_roadEast` or `tileSand_roadEast`).
      - Non-road cells still use appropriate base grass/sand tile IDs.
  - **Green – Make the test(s) pass:**
    - Extend `GenerateFromImage` (or a helper) to:
      - Compute, for each road cell, a neighbor bitmask based on N/E/S/W road presence.
      - Use this bitmask to distinguish between:
        - Vertical road segments (neighbors in N and S, not E/W).
        - Horizontal road segments (neighbors in E and W, not N/S).
      - Map these patterns to specific straight-road tile IDs for the appropriate terrain type; keep other road patterns mapped to a fallback ID for now.
  - **Refactor – Clean up with tests green:**
    - Extract a `selectRoadTileForStraight` helper that takes the terrain kind and neighbor bitmask and returns the tile ID, to avoid cluttering `GenerateFromImage`.
    - Keep mapping logic data-driven where possible (e.g., small lookup tables keyed by bitmask).
- **CI / Checks:**
  - `go test ./game/maps` then `go test ./...`.

---

### Step 3: Implement corners, T-junctions, crossings, and basic transitions

- **Workstream:** A – Game-layer image generator (game/maps)
- **Based on Design:**
  - `design.md` – Contracts and Data (Neighbor and Tile Selection: corners, T-junctions, crossings, terrain boundaries)
  - `design.md` – Testing and Safety Net (Neighbor-Aware Road Shapes, Terrain Boundaries)
- **Files:**
  - `game/maps/generator.go`
  - `game/maps/generator_test.go`
- **TDD Cycle:**
  - **Red – Failing test first:**
    - In `generator_test.go`, add targeted tests for more complex patterns:
      - Road corners: 2×2 or 3×3 images where a road cell has road neighbors in N+E, E+S, S+W, or W+N and grass/sand elsewhere; assert that the center/bend cell uses the correct corner tile ID (e.g., `*_roadCornerUR`, `*_roadCornerLR`, etc.).
      - T-junctions: cells with three road neighbors; assert that they map to `*_roadSplit*` tile IDs based on the missing direction.
      - Crossings: cells with four road neighbors; assert they map to `*_roadCrossing*` or `*_roadCrossingRound*` as chosen.
      - Simple grass–sand boundaries: horizontal or vertical stripes where one side is grass and the other is sand; assert that boundary cells use transition tiles (e.g., `tileGrass_transitionN`, `tileGrass_transitionE`) or, where tiles are not available, a documented fallback.
  - **Green – Make the test(s) pass:**
    - Extend neighbor-classification logic to compute full bitmasks and use them to:
      - Identify corner, T-junction, and crossing patterns and map them to the appropriate `*_roadCorner*`, `*_roadSplit*`, and `*_roadCrossing*` tile IDs, for both grass and sand terrains as needed.
      - Detect grass–sand (and where appropriate terrain–road) boundaries and map those edge cells to transition tiles defined in the existing assets; when no matching transition tile exists, fall back to a consistent base terrain tile.
  - **Refactor – Clean up with tests green:**
    - Consolidate neighbor-to-tile mapping into clear helper functions or tables (e.g., `selectRoadTile` and `selectTransitionTile`).
    - Add brief comments documenting which patterns are supported and any approximations due to tile set limitations.
- **CI / Checks:**
  - `go test ./game/maps` then `go test ./...`.

---

### Step 4: Rewire cmd/genesis to use the image-driven generator

- **Workstream:** B – Genesis CLI rewire (cmd/genesis)
- **Based on Design:**
  - `design.md` – Proposed Solution (Technical Overview: CLI), Architecture and Boundaries (binaries call game-layer generator)
  - `design.md` – Testing and Safety Net (CLI-Level Tests)
  - `design.md` – CI/CD and Rollout (new CLI contract, forward-only)
- **Files:**
  - `cmd/genesis/main.go`
  - `cmd/genesis/main_test.go`
- **TDD Cycle:**
  - **Red – Failing test first:**
    - Update `cmd/genesis/main_test.go` to reflect the new CLI contract:
      - Replace existing tests that call `run([]string{"seed", "width", "height", out})` with tests that call `run([]string{inputPNG, out})`.
      - In the happy-path test:
        - Create a small PNG file in a temp dir using `image/png` and `os.Create`, with pixels set to valid fixed RGB colors.
        - Call `run` and assert it returns `nil`.
        - Read the output JSON file into a `pkg/map.Map` and assert that `Width`/`Height` match the PNG dimensions and that `Tiles` has the correct size.
      - Add error-path tests:
        - Non-existent input PNG path → `run` returns an error.
        - A PNG with an unsupported color pixel → `run` returns an error that mentions invalid/unsupported colors.
      - Let these tests fail until `main.go` is adjusted.
  - **Green – Make the test(s) pass:**
    - Modify `cmd/genesis/main.go`:
      - Update `run(args []string) error` to:
        - Expect exactly two arguments: `inputPath` and `outputPath`; if not provided, return a usage error.
        - Use `os.Open` and `image/png.Decode` (or `image.Decode`) to load the input PNG.
        - Call `game/maps.GenerateFromImage(img)` to obtain a `*mappkg.Map`.
        - Ensure the directory for `outputPath` exists, then write JSON using the existing safe pattern (temp file + `os.Rename`).
      - Remove seed/width/height parsing and `NewGrassMap` calls from `run`.
      - Adjust `main()` to keep delegating to `run(os.Args[1:])` unchanged.
  - **Refactor – Clean up with tests green:**
    - Simplify error messages and usage text (e.g., `usage: genesis <input-png> <output-json>`), making them clear and concise.
    - Keep `run` focused and small; consider extracting an internal helper if argument parsing or I/O becomes large.
- **CI / Checks:**
  - `go test ./cmd/genesis` to validate CLI tests.
  - `go test ./...` for full suite.

---

### Step 5: Align docs and usage with image-driven generation

- **Workstream:** C – Safety net, docs, and cleanup
- **Based on Design:**
  - `design.md` – CI/CD and Rollout (document new workflow)
  - `design.md` – Observability and Operations (developer-facing behavior and clear errors)
  - `increment.md` – Tasks (document pixel color mapping and generator usage)
- **Files:**
  - Root `README.md` (or other top-level docs that mention map generation / genesis).
  - Increment docs already describe the intent but may be referenced.
- **TDD Cycle:**
  - **Red – Failing test first:**
    - Conceptually identify outdated references:
      - Search for `genesis` usage that still mentions seed/width/height.
      - Note any documentation that does not mention `game/assets/maps/map.png` and fixed RGB mapping.
  - **Green – Make the test(s) pass:**
    - Update `README.md` and any other relevant docs to:
      - Explain that `genesis` is invoked as `genesis <input-png> <output-json>`.
      - Describe the fixed RGB mapping (including exact RGB values) for grass, sand, and road.
      - Mention that unsupported colors cause the generator to fail with a clear error.
    - Optionally cross-link to this increment’s docs for more detail.
  - **Refactor – Clean up with tests green:**
    - Ensure docs are concise, consistent, and do not reference the old seeded generator as the main path.
    - Remove or clearly mark any references to the seed-based workflow as historical.
- **CI / Checks:**
  - `go test ./...` as a smoke check; consider a quick manual run of `genesis` following the updated docs to validate instructions.

---

### Step 6: Add a small integration-style tilemap composition check

- **Workstream:** C – Safety net, docs, and cleanup
- **Based on Design:**
  - `design.md` – Testing and Safety Net (Map Structure Validity, integration behavior)
  - `design.md` – Observability and Operations (ensure generated maps work with existing systems)
  - `increment.md` – Success Criteria (visible, coherent tilemap based on image)
- **Files:**
  - A test file in an appropriate package, for example:
    - `game/maps/integration_test.go` or
    - An additional test in `game/assets/tilemap_test.go` that exercises the generator + tilemap composition.
- **TDD Cycle:**
  - **Red – Failing test first:**
    - Add an integration-style test that:
      - Creates a small PNG image on disk or in memory with a simple pattern of grass, sand, and road that exercises at least one corner or junction.
      - Runs the generator (`GenerateFromImage`) or invokes `run` from `cmd/genesis` on that PNG to produce a `*mappkg.Map`.
      - Uses `assets.RegisterSpriteForTest` to register dummy sprites for all tile IDs that are expected to appear in the generated map.
      - Calls `assets.ComposeTilemap` and asserts that:
        - It returns a non-nil image.
        - The resulting image has the expected dimensions (`Width` * tileSize, `Height` * tileSize).
        - No `ErrTileSpriteNotFound` error is returned.
  - **Green – Make the test(s) pass:**
    - Ensure tile IDs emitted by the generator match the registered ones in the test; update the generator or test’s set of registered IDs if mismatches arise.
    - Fix any bugs revealed by the composed tilemap (e.g., uninitialized tiles or incorrect IDs at boundaries) while keeping unit tests green.
  - **Refactor – Clean up with tests green:**
    - Factor out helper functions inside the test to build sample images and register sprites to keep the test readable.
    - Document briefly in comments what scenario the integration test covers.
- **CI / Checks:**
  - `go test ./game/...` to verify integration with assets and maps.
  - `go test ./...` for full suite.

## 4. Rollout & Validation Notes

- **Suggested PR Grouping:**
  - PR 1: Steps 1–3 (introduce `game/maps` generator with color and neighbor logic + unit tests).
  - PR 2: Step 4 (rewire `cmd/genesis` to use the generator, update CLI tests).
  - PR 3: Steps 5–6 (docs alignment and integration-style tilemap composition test).

- **Validation Checkpoints:**
  - After Step 3:
    - `go test ./game/maps` and `go test ./...` pass.
    - Inspect generated `Map` for a few small synthetic PNGs to confirm tile IDs look reasonable for grass/sand/road and basic road shapes.
  - After Step 4:
    - `go test ./cmd/genesis` passes with the new CLI contract.
    - Running `genesis game/assets/maps/map.png out.json` produces a JSON map that the game can load without changes.
  - After Step 6:
    - Integration test confirms that a generated map can be composed into a tilemap image using `assets.ComposeTilemap` without missing tile sprites.
    - Manually running the game with the generated map shows a tilemap layout that visually matches the input PNG’s design (within the approximations documented in the design).
