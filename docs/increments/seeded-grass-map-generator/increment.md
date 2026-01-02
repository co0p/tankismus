# Increment: Seeded Grass Map Generator

## User Story

As a developer, I want to generate a grass-only tile map as a JSON matrix from a given seed and tile grid size, so that I can quickly create and reuse deterministic test maps using the tileGrass1 and tileGrass2 tile IDs.

## Acceptance Criteria

1. **Developer can generate a grass-only map from seed and dimensions**  
   - Given a seed, tile width, and tile height, a developer can trigger map generation in a straightforward way.  
   - A JSON file is produced without needing to run the game itself.

2. **Output is a JSON matrix of tile IDs**  
   - The main content of the output is a 2D matrix structure representing the tile grid.  
   - Each cell in the matrix contains a tile ID string.  
   - The structure makes it clear how many rows and columns are present.

3. **Only tileGrass1 and tileGrass2 are used**  
   - Every entry in the matrix is either the string "tileGrass1" or "tileGrass2".  
   - No other terrain types or objects appear in the output for this increment.

4. **Deterministic randomness based on seed**  
   - For the same seed, tile width, and tile height, the generator always produces exactly the same JSON matrix.  
   - Changing the seed (while keeping the dimensions the same) results in a visibly different grass pattern.

5. **Basic validity and inspectability of the JSON file**  
   - The produced JSON can be parsed successfully by standard JSON tools.  
   - A developer can inspect the file and easily understand the grid layout and which tiles were chosen.

6. **Clear behavior for invalid inputs**  
   - If tile dimensions are non-positive or the output path is unusable, a clear error is surfaced.  
   - No partial or corrupt JSON file is left behind on failure.

## Use Case

**Actors**

- Developer working on the tankismus project.  
- Map generation tool that produces JSON tile maps.

**Preconditions**

- The project can be built and run in the developer’s environment.  
- The developer has decided on:  
  - A seed value.  
  - Tile grid width (number of columns).  
  - Tile grid height (number of rows).  
  - An output file location and name (or is comfortable using a default).
- No game runtime map loading is required for this increment.

**Main Flow**

1. The developer decides to create a grass-only test map for experimentation or future integration.  
2. The developer invokes the map generator, providing:  
   - Seed.  
   - Tile width and tile height.  
   - Output file path (or accepting a sensible default).  
3. The generator validates the provided inputs and reports any obvious issues.  
4. The generator initializes its randomness based on the seed so that its behavior is deterministic.  
5. The generator constructs an internal grid with the requested width and height.  
6. For each cell in the grid, the generator selects either "tileGrass1" or "tileGrass2" using seed-based randomness.  
7. The generator encodes the grid as a JSON structure representing a 2D matrix of tile ID strings, with any basic metadata needed to understand the grid.  
8. The generator writes the JSON to the output file path without errors.  
9. The developer opens or parses the JSON file and verifies that:  
   - The matrix dimensions match the requested width and height.  
   - All entries are either "tileGrass1" or "tileGrass2".  
   - The JSON is valid and understandable.

**Alternate / Exception Flows**

- **Invalid dimensions**  
  - The developer provides a zero or negative tile width or height.  
  - The generator reports a clear error and does not create a JSON file.

- **Unusable output path**  
  - The specified output path is unwritable or invalid.  
  - The generator reports a clear error and does not leave a partial JSON file behind.

- **Determinism check**  
  - The developer runs the generator twice with the same seed and dimensions.  
  - The resulting JSON matrices are identical.  
  - The developer runs it again with a different seed (same dimensions) and observes a different grass pattern.

## Context

The project is a top-down, arcade-style tank game where maps and terrain will influence movement, tactics, and overall feel. Over time, the game will need varied maps with different terrain types and obstacles, but currently there is no simple, deterministic way for developers to generate and reuse map layouts.

For this early stage, the priority is to enable developers to create basic, repeatable maps that can serve as test data and as a foundation for future terrain and obstacle work. A simple grass-only generator using two tile variants (tileGrass1 and tileGrass2) provides enough variation to exercise map-related workflows without introducing complexity from additional terrain types.

The project’s constitution emphasizes small, testable increments and blackbox-style validation of behavior. A seed-based map generator fits well within this approach: it is narrow in scope, easy to demonstrate, and produces artifacts (JSON files) that can be inspected and compared.

This increment is explicitly developer-facing and does not require integration with the game’s runtime or rendering yet. Later increments can build on this generator to introduce richer terrain, obstacles, and in-game loading of generated maps.

## Goal

The goal of this increment is to provide a simple, deterministic map generation capability that produces a JSON matrix of grass tiles from a seed and tile grid size, enabling developers to quickly create, reuse, and share basic test maps.

**Scope**

- A developer-facing way to generate maps without running the full game.  
- Grass-only terrain using exactly two tile IDs: "tileGrass1" and "tileGrass2".  
- JSON output that clearly represents a 2D grid of tile IDs.  
- Deterministic behavior based on seed and grid dimensions.

**Non-Goals**

- Loading generated maps into the running game or rendering them on screen.  
- Introducing additional terrain types, obstacles, or interactive elements.  
- Defining collision, movement modifiers, or gameplay behavior tied to terrain.  
- Designing the long-term map file format in full detail; the format may evolve in later increments.

**Why this is a good increment**

- It is small and self-contained, focused solely on producing a simple, repeatable JSON map artifact.  
- It is easy to evaluate: developers can generate maps, inspect the JSON, and verify determinism quickly.  
- It fits naturally into the existing build and run workflow without special processes.  
- It lays groundwork for future, more complex map and terrain systems without committing to them prematurely.

## Tasks

- **Task:** Provide a simple way for developers to invoke map generation with seed and grid dimensions  
  - **User/Stakeholder Impact:** Developers can create test maps on demand without manual editing, improving iteration speed.  
  - **Acceptance Clues:** A developer can specify a seed, tile width, tile height, and output location, and receive a JSON file in response.

- **Task:** Ensure the JSON output clearly represents a 2D grid of tile IDs  
  - **User/Stakeholder Impact:** Developers and future tools can understand and consume the map structure easily.  
  - **Acceptance Clues:** The JSON contains an obvious matrix-like structure with rows and columns of tile ID strings, and it is straightforward to map indices to positions.

- **Task:** Restrict tile values to tileGrass1 and tileGrass2 for all grid cells  
  - **User/Stakeholder Impact:** The generated maps are predictable in terms of terrain type, simplifying early testing and debugging.  
  - **Acceptance Clues:** Spot-checks and simple scripts confirm that every cell value in the matrix is either "tileGrass1" or "tileGrass2".

- **Task:** Make generation deterministic based on seed and dimensions  
  - **User/Stakeholder Impact:** Developers can regenerate the same map later (or on another machine) by reusing the same inputs, aiding reproducible testing.  
  - **Acceptance Clues:** Re-running generation with the same seed and dimensions yields identical JSON; changing only the seed yields a different pattern.

- **Task:** Handle invalid inputs and output issues with clear feedback  
  - **User/Stakeholder Impact:** Developers receive understandable messages when something goes wrong, avoiding silent failures or confusing partial outputs.  
  - **Acceptance Clues:** Invalid dimensions or unwritable output paths result in clear error reporting and no corrupted or half-written files.

- **Task:** Document the intended JSON structure and usage at a high level  
  - **User/Stakeholder Impact:** Future work on map loading, terrain, and tools can build on a shared understanding of what the generator produces.  
  - **Acceptance Clues:** There is a short description of the JSON structure and how seed and dimensions relate to the matrix, accessible to team members.

## Risks and Assumptions

- **Risk:** The initial JSON structure might need adjustments once integrated with rendering and gameplay.  
  - **Assumption:** The structure will remain simple enough that minor changes will not invalidate the value of this increment.

- **Risk:** If determinism is incorrectly implemented, developers may see inconsistent maps for the same seed and dimensions.  
  - **Assumption:** Simple manual checks (re-running with the same inputs) will be part of validation, making such issues easy to spot and correct.

- **Risk:** Overly simplistic maps might not reveal issues that appear with more complex terrain later.  
  - **Assumption:** This increment is a foundation; future increments will introduce more terrain types and structures for richer testing.

- **Assumption:** Developers are comfortable using a command or similar mechanism to generate maps and inspect JSON files directly.  
- **Assumption:** The seed concept is acceptable to use as the primary handle for reproducible maps at this stage.

## Success Criteria and Observability

- **Behavioral Success Criteria**  
  - Developers can reliably generate grass-only maps using seed and grid dimensions and obtain valid JSON files.  
  - Generated maps are reproducible: the same inputs always lead to the same output.  
  - Different seeds with the same dimensions lead to noticeably different grass patterns.

- **Observability**  
  - Developers can compare JSON files across runs to verify determinism and variation.  
  - Simple JSON parsing tools or scripts can confirm that matrix dimensions and tile values match expectations.  
  - Any failures (invalid inputs or output issues) are clearly visible through error messages.

## Process Notes

- This increment should be implemented in small, safe steps that can be tested frequently, in line with the project’s preference for small, reversible changes.  
- Normal build and test workflows should be sufficient to validate the generator’s behavior.  
- The map generator should be introduced in a way that does not disrupt existing game flows; it is an additive capability.  
- Future increments can extend or refine the generator without needing special rollout or coordination.

## Follow-up Increments (Optional)

- Extend the generator to support additional terrain types (e.g., roads, sand, water) and basic obstacle placement while keeping seed-based determinism.  
- Introduce map loading into the game runtime so that generated JSON maps can be rendered and played.  
- Add support for multiple layers or richer metadata (such as spawn points, obstacles, and objectives) in the map format.  
- Provide simple visualization or editing tools to view and tweak generated maps.  
- Explore parameterization of terrain distribution (e.g., clusters of certain tiles) for more interesting layouts.

## PRD Entry (for docs/PRD.md)

- **Increment ID:** seeded-grass-map-generator  
- **Title:** Seeded Grass Map Generator  
- **Status:** Proposed  
- **Increment Folder:** docs/increments/seeded-grass-map-generator/  
- **User Story:** As a developer, I want to generate a grass-only tile map as a JSON matrix from a given seed and tile grid size, so that I can quickly create and reuse deterministic test maps using the tileGrass1 and tileGrass2 tile IDs.  
- **Acceptance Criteria:**  
  - Developers can generate a JSON file representing a grass-only map from seed and grid dimensions without running the game.  
  - The JSON contains a 2D matrix of tile ID strings where every entry is either "tileGrass1" or "tileGrass2".  
  - The generator is deterministic with respect to seed and dimensions, and the JSON is valid and easy to inspect.  
- **Use Case Summary:** A developer invokes a map generator with a seed, tile width, tile height, and an output path; the tool validates inputs, fills a width×height grid with tileGrass1 and tileGrass2 using seed-based randomness, writes a valid JSON matrix to disk, and allows the developer to regenerate or vary maps by reusing or changing the seed.