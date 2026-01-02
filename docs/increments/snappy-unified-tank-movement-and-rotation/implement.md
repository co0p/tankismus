Implement: Snappy Unified Tank Movement and Rotation
====================================================

2. Context
----------

- Goal: Implement a unified, snappy tank movement and rotation model with acceleration toward bounded linear speeds, bounded rotational speed, and rotation around each tank's visual center for all tanks, as described in `increment.md` and `design.md`.
- Non-goals: Terrain-dependent movement, complex physics (sliding, friction beyond what is needed), AI decision-making, and detailed per-tank balancing remain out of scope for this increment.
- Design approach: Control systems (player now, AI later) write normalized movement intents; the movement system applies acceleration/deceleration and caps to compute velocities and integrate position/rotation; the render system draws tanks rotated around their centers.
- Constitution: `CONSTITUTION.md` defines `constitution-mode: lite`, favoring small, pragmatic steps, blackbox-style tests for non-trivial logic, and simple `go test ./...`-driven CI.

Status: Not started
Next step: Step 1 – Introduce tank control intent and movement parameter components


1. Workstreams
---------------

- Workstream A – Components & Data Contracts
- Workstream B – Movement System Behavior & Tests
- Workstream C – Control System Refactor (Player Input → Intents) & Tests
- Workstream D – Rendering Pivot & Rotation
- Workstream E – Debug Observability & Feel Checks (Optional)


2. Steps
--------

- [ ] Step 1: Introduce tank control intent and movement parameter components
- [ ] Step 2: Add core movement model behavior and tests
- [ ] Step 3: Ensure position and rotation integration matches movement direction
- [ ] Step 4: Refactor InputMovementSystem to produce intents, not raw velocities
- [ ] Step 5: Wire new components and systems in the run scene
- [ ] Step 6: Update the render system for rotation around tank center
- [ ] Step 7: Optional debug overlay/logging for movement inspection


### Step 1: Introduce tank control intent and movement parameter components

- Workstream: A – Components & Data Contracts
- Based on Design: `design.md` §4–5 (Architecture and Boundaries; Contracts and Data)
- Files: `game/components/components.go`, `game/components/components_test.go`, `game/scenes/run/run.go`
- TDD Cycle:
  - Red – Failing test first:
    - In `game/components/components_test.go`, add tests that construct new tank control intent and tank movement parameter components and assert:
      - They implement the ECS `Component` interface (their `Type()` methods return distinct `ComponentType` values).
      - Default values conform to the design contracts (intent starts neutral with `throttle=0` and `turn=0`; parameters have non-negative, reasonable defaults for max speeds and accel/decel values).
  - Green – Make the test(s) pass:
    - In `game/components/components.go`, add `ControlIntent` and `MovementParams` component types with unique `ComponentType` constants and struct fields matching the conceptual schemas from `design.md`.
    - In `game/scenes/run/run.go`, when creating the player entity, attach default instances of `ControlIntent` and `MovementParams` alongside existing `Transform`, `Velocity`, and `Sprite` components.
  - Refactor – Clean up with tests green:
    - Tidy component organization and comments to clearly group movement-related components while preserving the existing ECS boundaries and naming style.
- CI / Checks:
  - Run `go test ./game/components/...` and then `go test ./game/...` to ensure components compile and tests pass across the game layer.


### Step 2: Add core movement model behavior and tests

- Workstream: B – Movement System Behavior & Tests
- Based on Design: `design.md` §4–6 (Movement system responsibilities; Movement Behavior Contract)
- Files: `game/systems/movement.go`, `game/systems/movement_test.go`
- TDD Cycle:
  - Red – Failing test first:
    - Create `game/systems/movement_test.go` with tests that:
      - Construct a small ECS world and an entity with `Transform`, `Velocity`, `ControlIntent`, and `MovementParams`.
      - Set a constant forward throttle in the intent and call the movement logic repeatedly with a fixed `dt`, asserting that:
        - Linear speed increases smoothly over updates and approaches but does not exceed `maxForwardSpeed`.
      - Set a constant backward throttle and assert that linear speed approaches but does not exceed `-maxBackwardSpeed`.
      - Set throttle from non-zero to zero and assert that speed decays toward zero within a reasonable number of steps based on the configured deceleration.
      - Set a constant non-zero turn value and assert that angular velocity approaches but does not exceed `maxTurnRate` in magnitude.
    - Initially, target an extended behavior of `MovementSystem` (or a small internal helper it uses) so tests compile but fail due to missing acceleration/cap logic.
  - Green – Make the test(s) pass:
    - In `game/systems/movement.go`, extend `MovementSystem` to:
      - For entities with `Transform`, `Velocity`, `ControlIntent`, and `MovementParams`, compute target linear and angular speeds from the intent (`throttle`/`turn`) and parameters (`maxForwardSpeed`, `maxBackwardSpeed`, `maxTurnRate`).
      - Adjust the entity's current linear and angular velocity toward these targets at rates bounded by the configured `linearAcceleration`/`linearDeceleration` and `angularAcceleration`/`angularDeceleration`, clamping to the max values.
      - Preserve existing integration of velocity to position and rotation, but now driven by the updated velocities.
  - Refactor – Clean up with tests green:
    - If needed, extract pure helper functions inside `game/systems` for "step velocity toward target" to simplify logic while keeping tests focused on observable behavior of `MovementSystem`.
- CI / Checks:
  - Run `go test ./game/systems/...` to validate movement behavior, then `go test ./game/...`.


### Step 3: Ensure position and rotation integration matches movement direction

- Workstream: B – Movement System Behavior & Tests
- Based on Design: `design.md` §5–6 (Movement Behavior Contract; Visual alignment)
- Files: `game/systems/movement.go`, `game/systems/movement_test.go`
- TDD Cycle:
  - Red – Failing test first:
    - Extend `game/systems/movement_test.go` with tests that verify:
      - With non-zero linear speed and zero angular velocity, repeated updates move the entity along a straight line consistent with its rotation (e.g. X/Y displacements match `cos(rotation)`/`sin(rotation)` scaled by speed and `dt`).
      - With non-zero linear and angular velocities, position changes in both axes and rotation changes monotonically such that the entity follows a smooth arc.
    - These tests should fail if position/rotation integration does not fully match the contract (for example, if VX/VY are inconsistent with rotation or if rotation is not applied correctly).
  - Green – Make the test(s) pass:
    - Adjust or clarify position integration in `game/systems/movement.go` so that:
      - Movement direction is always consistent with the entity's rotation (either by treating linear speed as a scalar derived from VX/VY or by ensuring VX/VY are always derived from rotation and linear speed).
      - Rotation advancement strictly follows `Angular * dt`.
  - Refactor – Clean up with tests green:
    - Simplify variable naming and comments in `MovementSystem` to make the relationship between control intent, velocity, and transform explicit and easy to reason about.
- CI / Checks:
  - Run `go test ./game/systems/...` and `go test ./game/...`.


### Step 4: Refactor InputMovementSystem to produce intents, not raw velocities

- Workstream: C – Control System Refactor (Player Input → Intents) & Tests
- Based on Design: `design.md` §2, §4–5 (Control systems express intent only; separation of concerns)
- Files: `game/systems/input_movement.go`, `game/systems/input_movement_test.go`
- TDD Cycle:
  - Red – Failing test first:
    - Create `game/systems/input_movement_test.go` with tests that:
      - Build a minimal ECS world with a player entity having `Transform`, `Velocity`, `ControlIntent`, and `MovementParams`.
      - Stub or control the input state via the `pkg/input` API (for example, by configuring actions or by structuring the system to accept an input query function) so that specific combinations of W/S/A/D are treated as pressed.
      - Assert that `InputMovementSystem` updates the entity's `ControlIntent` as follows:
        - W pressed → `throttle` = +1.
        - S pressed → `throttle` = -1.
        - W and S both released → `throttle` = 0.
        - A pressed → `turn` = -1.
        - D pressed → `turn` = +1.
        - No turn keys → `turn` = 0.
      - Assert that `Velocity` is not directly overwritten by input.
    - Initially, these tests will fail because `InputMovementSystem` still writes to `Velocity` directly.
  - Green – Make the test(s) pass:
    - In `game/systems/input_movement.go`, change `InputMovementSystem` so that it:
      - Reads action states from `pkg/input` as before.
      - Maps actions into normalized `throttle` and `turn` values on the player's `ControlIntent` component.
      - Leaves `Velocity` untouched, allowing `MovementSystem` to handle actual speed and rotation.
  - Refactor – Clean up with tests green:
    - Remove or simplify any leftover logic that directly manipulates `Velocity` in response to input, ensuring that control systems now exclusively write intents.
- CI / Checks:
  - Run `go test ./game/systems/...` and `go test ./game/...`.


### Step 5: Wire new components and systems in the run scene

- Workstream: A/B/C – Components, Movement, and Control Wiring
- Based on Design: `design.md` §4 (Architecture and Boundaries) and §6 (Consistency across tanks)
- Files: `game/scenes/run/run.go`, `game/scenes/run/run_test.go`
- TDD Cycle:
  - Red – Failing test first:
    - Create `game/scenes/run/run_test.go` with tests that:
      - Construct a new run `Scene` and verify that the created player entity in its `world` has the following components: `Transform`, `Velocity`, `ControlIntent`, `MovementParams`, and `Sprite`.
      - Call `Scene.Update` with a small `dt` and assert that it performs, in order: input polling, `InputMovementSystem`, and `MovementSystem` (for example, by observing resulting changes in control intent and transform over multiple updates).
    - These tests will initially fail until components and system calls are fully wired.
  - Green – Make the test(s) pass:
    - Ensure `run.New` adds all required components to the player entity and that the scene's `Update` method:
      - Polls input via `pkg/input.Poll()`.
      - Invokes `InputMovementSystem` and then `MovementSystem` with consistent `dt` handling.
  - Refactor – Clean up with tests green:
    - If helpful for testability and clarity, simplify or centralize frame time handling while preserving actual gameplay behavior.
- CI / Checks:
  - Run `go test ./game/scenes/run/...` and `go test ./game/...`.


### Step 6: Update the render system for rotation around tank center

- Workstream: D – Rendering Pivot & Rotation
- Based on Design: `design.md` §4–5 (Render system responsibilities; Visual Pivot and Alignment)
- Files: `game/systems/render.go`, `game/systems/render_test.go`
- TDD Cycle:
  - Red – Failing test first:
    - Create `game/systems/render_test.go` with tests that:
      - Use a small fake or stub for the screen (or wrap `ebiten.Image` where practical) to capture the `DrawImageOptions` used by `RenderSystem` when drawing a tank.
      - Given a `Transform` with different rotation angles and a `Sprite`, assert that:
        - The transformation matrix (`GeoM`) applies rotation about the sprite center before translation.
        - The translation places the sprite so its visual center coincides with the tank's logical position (within an acceptable tolerance).
    - Initially, these tests will fail because the current render system only translates and does not rotate or center.
  - Green – Make the test(s) pass:
    - In `game/systems/render.go`, change `RenderSystem` so that it:
      - Retrieves the sprite image and computes its center based on bounds.
      - Applies rotation from `Transform.Rotation` around the sprite center, then translates so that the center aligns with `Transform.X`/`Transform.Y`.
  - Refactor – Clean up with tests green:
    - If needed, factor setup of `DrawImageOptions` into a small helper function to avoid repetition and keep the main loop over entities straightforward.
- CI / Checks:
  - Run `go test ./game/systems/...` and `go test ./game/...`.


### Step 7: Optional debug overlay/logging for movement inspection

- Workstream: E – Debug Observability & Feel Checks (Optional)
- Based on Design: `design.md` §8 (Observability and Operations)
- Files: `game/scenes/run/run.go`, optionally a small debug helper in `game/`
- TDD Cycle:
  - Red – Failing test first (lightweight/optional):
    - If implementing a debug overlay or helper, add a simple test (e.g. in `game/scenes/run/run_test.go` or a debug-specific test file) that:
      - Given a tank entity with known speed and angular velocity, a debug helper function returns a formatted string or structure summarizing current and target speeds consistent with movement parameters and control intents.
  - Green – Make the test(s) pass:
    - Implement the debug helper or overlay so that, when a debug flag is active, it reads movement state (speed, angular velocity, possibly target speeds) and exposes it for on-screen drawing or logging.
  - Refactor – Clean up with tests green:
    - Ensure debug functionality is clearly separated from core game logic and easy to disable in normal builds.
- CI / Checks:
  - Run `go test ./game/...` to confirm all tests, including optional debug tests, pass.


3. Rollout & Validation Notes
-----------------------------

- Suggested grouping into PRs:
  - PR 1: Steps 1–2 (components and core movement model + tests).
  - PR 2: Steps 3–5 (movement integration checks, input refactor to intents, and run scene wiring + tests).
  - PR 3: Step 6–7 (render rotation/pivot and optional debug observability).

- Suggested validation checkpoints:
  - After Step 2: `go test ./game/systems/...` should pass, and inspection of movement tests should show linear and angular speeds respecting acceleration and caps.
  - After Step 5: Running the game locally should show the player tank moving and rotating using the new movement model, with consistent feel and no obvious regressions in controllability.
  - After Step 6: Visual inspection during play should confirm that tanks rotate around their centers and move in alignment with their visual facing direction.
