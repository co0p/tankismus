# Increment: Image-Driven Map Generator From Pixel Map

## User Story

As a developer working on Tankismus, I want to generate playable game maps from a color-coded source image so that I can design and iterate on levels visually (by editing the image) without relying on a random grass map generator.

## Acceptance Criteria

- When the map generator is run against a specified input pixel map (for example, `game/assets/maps/map.png`), the resulting JSON map layout is derived solely from the image pixels, not from any seeded random generator.
- Green pixels in the source image always become grass terrain tiles in the generated map; yellow pixels always become sand terrain tiles; black pixels always become road terrain tiles, and this mapping is clearly documented for contributors.
- For contiguous road pixels, the generator selects appropriate road subtiles (straight segments, corners, T-junctions, crossings, etc.) based on the configuration of adjacent road pixels, producing visually coherent road networks when rendered.
- At boundaries between grass and sand (and between terrain and road where appropriate), the generator uses available transition tiles so that visible edges look intentional rather than abrupt, while interior areas use suitable interior tiles (for example, `tileGrass1/2`, `tileSand1/2`).
- The generator outputs a JSON map in the existing map format (including width, height, and a 2D matrix of tile IDs) so that current map consumers in the game continue to function without changes.
- The existing seeded grass map generator and its CLI flow are fully retired from normal workflows, with the “genesis” or equivalent map generation entrypoint now using the image-driven generator instead.
- Automated checks (for example, tests run via `go test ./...`) validate color-to-terrain mapping and neighbor-aware subtile selection for representative small input images and fail if these behaviors regress.

## Use Case

**Actors**

- Developer working on Tankismus.
- Map generator tool (invoked via a CLI or equivalent entrypoint).
- Game runtime that consumes the generated JSON map.

**Preconditions**

- The Tankismus project is checked out and buildable according to the main README.
- A source pixel map exists at a known location (for example, `game/assets/maps/map.png`) and uses only the supported colors:
  - Green for grass terrain.
  - Yellow for sand terrain.
  - Black for road terrain.
- Tile images for grass, sand, and roads (including road variants and transitions) are present under the game assets images directory.
- The map generator tool is available to run from the project environment.

**Main Flow**

1. The developer opens the source pixel map in an image editor and designs or adjusts the level layout using only the supported colors to indicate grass, sand, and road.
2. From the project environment (for example, the repository root), the developer invokes the map generator tool, pointing it at the pixel map and specifying an output JSON file path.
3. The generator loads the pixel map, determines its dimensions, and iterates over each pixel in a well-defined order.
4. For each pixel, the generator:
   - Maps the pixel color to a logical terrain type (grass, sand, or road).
   - Examines neighboring pixels (at least north, south, east, and west, and possibly diagonals) to understand whether the current cell is interior, a boundary, a corner, a straight road segment, a T-junction, or a crossing.
   - Selects a specific tile ID from the existing tile set that matches the terrain type and neighborhood configuration.
5. After processing all pixels, the generator assembles a map object containing width, height, and a 2D matrix of tile IDs in the established JSON map format and writes it to the requested output path.
6. The developer runs the game (or a suitable demo) that loads this JSON map and renders the ground tilemap.
7. In the running game, the visible level layout (grass, sand, roads, and transitions) matches the intended design drawn in the pixel map, and roads and terrain boundaries appear visually coherent.
8. If the layout or appearance is not as desired, the developer edits the pixel map and re-runs the generator, iterating until satisfied.

**Alternate / Exception Flows**

- A1: Unsupported colors in the pixel map
  - The generator encounters a pixel color that does not map to grass, sand, or road.
  - It reports a clear error (for example, including pixel coordinates or counts of problematic pixels) and exits with a non-zero status so the developer can correct the image.

- A2: Missing or mismatched tile identifiers
  - The neighbor-aware logic chooses a tile ID that does not correspond to an available tile sprite.
  - The generator reports a clear error or warning that identifies the problematic tile ID and location so the developer can adjust the mapping logic or assets.

- A3: Input or output I/O issues
  - The generator cannot read the input pixel map or cannot write the output file.
  - It exits with a non-zero status and a descriptive error message so the developer can fix path or permission issues.

- A4: Game run with missing or stale generated map
  - The game is started without a freshly generated JSON map or cannot load the expected map file.
  - The behavior is clearly signaled (for example, via logs or an obvious fallback behavior) so developers understand that they need to run the map generator.

## Context

The current project already includes a JSON-based tile map format and a seeded grass-only map generator used via a CLI tool. This existing generator is useful for quickly producing deterministic test maps but does not support intentional level design based on terrain types like sand and roads. Designers and developers have no straightforward way to author a specific map layout without modifying code or working around the random generator.

At the same time, the game assets already include tile sprites for grass, sand, roads, and various transition and junction shapes, and the run scene is capable of rendering a composed tilemap as a ground layer beneath the player tank. There is also an embedded pixel map asset that can serve as the visual source of truth for level layout.

The desired direction is to empower developers to control map layout visually by editing a pixel map while retaining the existing JSON map format and downstream consumers. The constitution emphasizes small, testable increments, blackbox-style tests, and clear boundaries between game logic, tools, and engine-style packages. Retiring the old seeded grass generator in favor of an image-driven approach aligns with the principle of avoiding long-lived legacy paths.

## Goal

The goal of this increment is to make image-based map generation the primary way to create playable maps in Tankismus: given a color-coded pixel map, developers can reliably produce a JSON map in the existing format that renders into a coherent grass/sand/road tilemap using the current tile assets.

In scope:

- Providing a generator that accepts a pixel map as input and produces a JSON map matching the existing format.
- Defining and documenting a simple color-to-terrain mapping (green/grass, yellow/sand, black/road).
- Implementing neighbor-aware selection of road and transition subtiles so that roads and terrain boundaries render cleanly.
- Updating the main map generation workflow to use this new generator instead of the seeded grass-only generator.

Out of scope for this increment:

- Adding new terrain types beyond grass, sand, and road.
- Changing how the game runtime loads or renders maps beyond what is needed to consume the new generator’s output.
- Introducing advanced editing tools or UIs beyond the pixel map workflow.
- Implementing terrain-dependent movement, collisions, or AI behavior; these remain future increments that can build on top of the new maps.

This is a good increment because it is small and self-contained, replaces a narrow piece of functionality (map generation) with a more powerful but still straightforward approach, and can be validated quickly through both automated tests and visible in-game behavior.

## Tasks

- Task: Define and communicate the pixel color to terrain mapping and usage of the source pixel map.
  - User/Stakeholder Impact: Developers know exactly how to author and edit `map.png` (or equivalent) to represent grass, sand, and roads using specific colors.
  - Acceptance Clues: Project documentation or generator help text clearly describes which colors are supported and how they map to terrain types, and new contributors can successfully create a simple map using only that guidance.

- Task: Provide a generator experience that takes a pixel map as input and outputs a JSON map in the existing format.
  - User/Stakeholder Impact: Developers can run a single, well-defined command or tool to turn a pixel map into a game-ready JSON map without editing code.
  - Acceptance Clues: Given a valid input pixel map and output path, the tool completes successfully, writes a JSON file in the established map format, and can be invoked consistently as part of local workflows.

- Task: Implement neighbor-aware tile selection for roads and terrain boundaries.
  - User/Stakeholder Impact: When maps are rendered, roads and transitions between terrain types look visually coherent and intentional, improving the perceived quality of the level.
  - Acceptance Clues: For representative test images and real maps, observers can see that straight roads, corners, junctions, and boundaries align correctly, with appropriate subtile variants selected based on surrounding terrain.

- Task: Retire the seeded grass-only generator from normal workflows and align the main generation entrypoint with the image-driven approach.
  - User/Stakeholder Impact: Developers no longer rely on a random grass-only generator; instead, the image-based generator becomes the standard way to create or update maps, reducing confusion and duplicated paths.
  - Acceptance Clues: The primary documented way to generate maps uses the pixel map input; references to the old seeded grass generator are removed or clearly marked as historical, and day-to-day development no longer depends on it.

- Task: Establish a basic automated safety net for the new generator.
  - User/Stakeholder Impact: Future changes to the generator are less likely to silently break color mappings or neighbor logic, preserving confidence in map generation.
  - Acceptance Clues: Running the project’s standard test command executes checks that validate color-to-terrain mapping and neighbor-aware tile selection for small, controlled input images, and failures are clear when behavior is changed unintentionally.

## Risks and Assumptions

- Risk: The existing tile asset set may not cover every desired road or transition configuration, which could limit how complex road networks or terrain boundaries look in early versions.
- Risk: Developers might accidentally use unsupported colors in the pixel map, leading to confusing errors or incomplete maps if error handling is not clear.
- Risk: If the generator’s neighbor-aware logic is overly complex or brittle, it may be hard to reason about or adjust, increasing maintenance cost.

- Assumption: The current JSON map format and map-loading behavior are sufficient and should not be changed in this increment.
- Assumption: Using a single shared pixel map as the source of truth for level layout fits existing and near-term workflows for this project.
- Assumption: Developers are comfortable using an external image editor to manipulate the pixel map and do not need a bespoke in-game editor in this increment.

Mitigations include keeping the neighbor logic as simple as possible for the first version, providing clear documentation and error messages for unsupported colors, and treating any major expansion of tile patterns or terrain types as follow-up increments.

## Success Criteria and Observability

- Developers can, without touching code, modify the layout of grass, sand, and roads by editing the pixel map and then regenerating the JSON map, and the resulting changes are visible in the game.
- The image-driven generator becomes the documented, default way to generate maps, and the old seeded grass generator is no longer part of normal workflows.
- Automated tests covering color mapping and key neighbor-aware scenarios pass reliably in local runs and CI, and fail clearly if behavior changes unexpectedly.
- Informal visual inspection of rendered maps confirms that roads and terrain boundaries match the intended design in the pixel map and that there are no obvious broken joints or tile mismatches.

## Process Notes

- This increment should be implemented as a series of small, safe changes that introduce the new generator, wire it into existing workflows, and then remove reliance on the old seeded grass-only generator.
- All changes should pass through the normal build and test pipeline (including `go test ./...`) and avoid special, one-off deployment or release steps.
- The rollout should favor simplicity: once the image-driven generator is validated, it becomes the standard path for developers, and the older generator is removed rather than retained as a parallel fallback.
- Any adjustments needed to documentation or onboarding materials should be made alongside the generator change so that new contributors see the updated workflow.

## Follow-up Increments (Optional)

- Extend the color and terrain mapping to support additional terrain types such as water, obstacles, or destructible cover, with corresponding tile assets and visual rules.
- Introduce terrain-aware gameplay behaviors, such as movement speed changes or pathfinding weights based on terrain types derived from the generated map.
- Provide higher-level tools or visualizations (for example, an in-game debug overlay or a simple viewer) to inspect generated maps without running a full game session.
- Explore saving and loading multiple named maps to support different levels or arenas, building on the same image-driven generation approach.

## PRD Entry (for docs/PRD.md)

- Increment ID: image-driven-map-generator-from-pixel-map
- Title: Image-Driven Map Generator From Pixel Map
- Status: Proposed
- Increment Folder: docs/increments/image-driven-map-generator-from-pixel-map/
- User Story: As a developer working on Tankismus, I want to generate playable game maps from a color-coded source image so that I can design and iterate on levels visually (by editing the image) without relying on a random grass map generator.
- Acceptance Criteria:
  - Map generation for normal workflows is driven by an input pixel map instead of a seeded random generator.
  - Green, yellow, and black pixels map deterministically to grass, sand, and road tiles, and this behavior is documented.
  - Neighbor-aware subtile selection produces visually coherent roads and terrain boundaries when rendered.
  - The generator outputs the existing JSON map format and replaces the prior seeded grass-only generator in day-to-day use.
  - Basic automated checks validate the generator’s color mapping and neighbor logic under the standard test command.
- Use Case Summary: Developers edit a color-coded pixel map to design a level, run a generator tool that converts the image into a JSON map using neighbor-aware tile selection for grass, sand, and roads, and then run the game to see a rendered tilemap ground layer that matches the designed layout; errors for unsupported colors or missing tiles are clearly reported to support quick iteration.