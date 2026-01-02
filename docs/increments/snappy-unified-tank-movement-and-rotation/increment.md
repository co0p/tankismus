# Increment: Snappy Unified Tank Movement and Rotation

## User Story

As a player, I want all tanks in the game (player and enemies) to use the same smooth, snappy tank-style movement and rotation, so that the controls feel fast, responsive, and consistent across the game.

## Acceptance Criteria

1. **Smooth linear acceleration with bounded speeds for all tanks**
   - When any tank starts moving forward or backward, its speed increases over time from rest toward a defined maximum speed instead of instantly jumping to full speed.
   - Maximum forward and backward speeds are clearly bounded and remain stable during normal play.
   - Observing multiple tanks (player and enemies), they all respect this “accelerate toward max speed” behavior.

2. **Bounded rotational speed with smooth, snappy turning**
   - When a tank is commanded to rotate (via player input or AI), it turns at a consistent rotational speed rather than snapping to a new angle.
   - There is a clear maximum turn rate so that tanks cannot spin arbitrarily fast.
   - Reversing or stopping rotation input results in a visibly smooth and quick change in turning behavior, without jarring overshoot.

3. **Rotation around each tank’s visual center**
   - When any tank rotates, the visual sprite or representation clearly pivots around the tank’s center point (no noticeable “orbiting” or off-center wobble).
   - Forward and backward movement follow the tank’s facing direction as seen on screen, so visual orientation and movement direction always match.

4. **Consistent movement model across tank types**
   - The same movement model (acceleration, maximum linear speeds, rotational speed limits) applies to all tank entities.
   - It is possible to use different parameter values per tank type later, but no tank uses a legacy “instant speed” or “snap rotation” behavior once this increment is complete.

5. **Terrain-independent movement for this increment**
   - For this increment, tank movement behavior (acceleration, maximum speed, turning) is not modified by terrain type; tanks behave the same on all surfaces.
   - Any future terrain-based movement effects, if introduced, are explicitly out of scope and handled in separate increments.

6. **No regressions in basic controllability and feel**
   - The player can reliably move, rotate, and navigate around obstacles using the keyboard without new obvious glitches (for example, getting stuck, uncontrolled spinning, or unresponsive turns).
   - Short test sessions show that movement feels at least as controllable as before, with noticeably snappier and more “tank-like” behavior.

## Use Case

**Actors**

- Player controlling a tank via keyboard.
- AI-controlled tanks using movement commands from their behavior logic.
- Game loop updating all tanks each frame.

**Preconditions**

- The game is running and at least one tank entity (player and/or enemies) is active on screen.
- Input is available for the player (keyboard).
- Each tank has an associated facing direction and a visible sprite or representation on screen.

**Main Flow**

1. The player presses and holds the move-forward key while their tank is idle or moving slowly.
2. The player tank’s forward speed increases smoothly from its current speed toward a defined maximum forward speed, without instantly jumping to full speed.
3. The player releases the move-forward key; the tank stops accelerating and comes to rest or slows quickly enough to feel snappy, without long sliding.
4. The player presses and holds rotate-left or rotate-right.
5. The player tank rotates smoothly around its center point at a bounded rotational speed, clearly visible as continuous turning rather than snapping.
6. While rotating, if the player also holds forward or backward, the tank moves in the direction it is currently facing, so motion and visual orientation always match.
7. Enemy tanks, when instructed by their AI to move or rotate, exhibit the same style of acceleration and turning: they speed up toward a maximum speed, respect a maximum turn rate, and visually rotate around their center.
8. During normal play, all tanks feel responsive and snappy: they reach useful speeds quickly, can change direction without sluggish delays, and stop or change rotation predictably when input or AI commands change.

**Alternate / Exception Flows**

- **Short key taps or brief movement commands**
  - The player quickly taps move or rotate keys; the tank responds with a small, noticeable movement or turn without overshooting or feeling unresponsive.

- **Rapid direction changes**
  - The player switches from moving forward to backward or from rotating left to rotating right; the tank transitions smoothly and quickly, without jitter or unrealistic instant reversal that breaks the sense of tank weight.

- **AI path adjustments**
  - An AI-controlled tank changes its desired direction (for example, to avoid an obstacle or reposition); the tank’s rotation and movement adjust using the same smooth, bounded turning and acceleration model rather than snapping to a new heading.

- **Speed cap behavior**
  - When any tank is commanded to move in one direction for an extended period, it does not exceed its defined maximum speed, and movement appears stable and predictable at that top speed.

## Context

The game is a top-down, fast-paced tank survival shooter with an ECS-style architecture. The design pillars emphasize minimal UI, immersion, retro aesthetics, and high-intensity action with responsive controls. Currently, tank movement and rotation need to better match the intended “tank feel” described in the game design: small, agile tanks that can move and turn quickly while still feeling like tracked vehicles.

The project constitution encourages small, safe increments, clear separation of concerns, and blackbox-style testing of key behaviors such as movement rules and timing. A unified movement model for all tanks fits well within this approach: it affects a focused aspect of gameplay, is easy to demonstrate, and can be verified through observable behavior in short play sessions.

This increment focuses specifically on:

- Introducing a shared, snappy movement and rotation model for all tanks.
- Ensuring that visual rotation and movement direction align around each tank’s center.
- Keeping changes independent of terrain and other gameplay systems for now.

All tank visual assets are initially authored to face upwards in their neutral orientation; this increment assumes that orientation when defining movement and rotation behavior.

Terrain-dependent movement, per-model balancing, and more advanced behaviors are intentionally deferred to future increments so that this change remains small, testable, and low-risk.

## Goal

The goal of this increment is to establish a unified, snappy tank movement and rotation model used by all tanks in the game, featuring linear acceleration toward bounded speeds and smooth, bounded rotation around each tank’s center.

**Scope**

- All tanks (player and enemies) share the same conceptual movement model:
  - Linear acceleration toward maximum forward and backward speeds.
  - Bounded rotational speed for smooth, continuous turning.
  - Movement direction always aligned with on-screen facing.
- The visual representation of tanks clearly rotates around their center, so the pivot point feels natural and consistent.
- Player input and AI movement commands both drive tanks through this shared movement model, resulting in a fast-paced but predictable feel.

**Non-Goals**

- Introducing terrain-dependent movement or traction differences.
- Implementing per-tank-type speed or rotation differences beyond using the shared model (fine-tuning per model is reserved for later increments).
- Changing weapons, shooting, collision rules, or enemy decision-making logic beyond how tanks physically move and rotate.

**Why this is a good increment**

- It is small and self-contained, focusing solely on how tanks move and rotate.
- It has clear, observable outcomes that can be evaluated quickly in short play sessions.
- It flows through the normal build, test, and release process without requiring special coordination.
- It lays a solid foundation for later balancing and terrain effects without locking in specific parameter choices too early.

## Tasks

- **Task:** Establish a unified tank movement behavior for all tanks
  - **User/Stakeholder Impact:** Players experience consistent, predictable movement for player and enemy tanks, enhancing the sense of fairness and control.
  - **Acceptance Clues:** Observing any tank over time shows it accelerating toward a clear maximum speed and rotating at a bounded rate, with no tank using an obviously different movement logic.

- **Task:** Ensure tanks visually rotate around their center with aligned movement direction
  - **User/Stakeholder Impact:** Players perceive tanks as small, agile tracked vehicles whose movement and visual orientation are tightly coupled, improving immersion and readability in fast-paced action.
  - **Acceptance Clues:** When a tank rotates in place, its sprite clearly pivots around its center; when moving forward or backward, its motion follows the on-screen facing direction without visual drift or offset.

- **Task:** Align player input and AI commands with the shared movement model
  - **User/Stakeholder Impact:** Players and observers see that both controlled and AI tanks behave according to the same physical rules, avoiding surprises when fighting enemies.
  - **Acceptance Clues:** Player-controlled and AI-controlled tanks both accelerate, decelerate, and turn in comparable ways when given similar commands, with no special-case movement that breaks consistency.

- **Task:** Validate controllability and feel in short play sessions
  - **User/Stakeholder Impact:** Players can quickly get a feel for how tanks respond, supporting the game’s fast-paced, arcade-like experience without frustration.
  - **Acceptance Clues:** Short internal playtests confirm that tanks reach useful speeds quickly, can navigate around obstacles reliably, and feel at least as controllable as before while noticeably more “tank-like.”

- **Task:** Document the movement model and tunable parameters at a product level
  - **User/Stakeholder Impact:** Future balancing work and new tank types can reuse the same movement concept while adjusting speeds and turn rates safely.
  - **Acceptance Clues:** There is a short, accessible description of the movement model and which parameters (for example, acceleration, max speed, turn rate) can vary per tank type in future increments.

## Risks and Assumptions

- **Risk:** Movement becomes too twitchy or too sluggish if the chosen acceleration and speed caps are not well tuned.
  - **Assumption:** Iterative tuning during implementation and playtesting will keep the feel snappy but controllable.

- **Risk:** Changes to movement behavior may subtly affect collision interactions, enemy pathing, or difficulty.
  - **Assumption:** Existing ECS boundaries and simple test scenarios will help detect regressions early, and adjustments can be made without broad architectural changes.

- **Risk:** Players used to the previous feel might need a brief adjustment period.
  - **Assumption:** The new model will be clearly better aligned with the game’s design pillars (fast, responsive tanks), making the change net-positive after a short adaptation.

- **Assumption:** All tanks in the game are intended to follow the same basic movement rules, with parameter differences rather than entirely separate movement concepts.
- **Assumption:** Terrain-based movement effects and more advanced behaviors will be introduced later, building on, not replacing, this model.

## Success Criteria and Observability

- **Behavioral Success Criteria**
  - In test runs, observers report that tanks feel noticeably more “tank-like,” snappy, and consistent in how they move and turn.
  - Player-controlled and AI-controlled tanks exhibit similar acceleration and turning behavior when given similar commands.
  - No obvious visual pivot issues or mismatches between facing direction and movement direction are observed.

- **Observability**
  - Short manual test scenarios (for example, moving in straight lines, rotating in place, circling obstacles) visibly demonstrate acceleration curves, speed caps, and bounded rotation.
  - Simple logs or debug overlays, if used, can show tank speed and rotation values over time to confirm that speed and turn rates remain within expected bounds.
  - After release, informal feedback from play sessions focuses on feel and control rather than confusion or glitches in movement.

## Process Notes

- This increment should be implemented via small, safe changes that can be tested frequently, in line with the project’s preference for small, reversible steps.
- Existing automated tests and new tests for movement behavior should run under the normal `go test` pipeline.
- The change should be rolled out through the standard build and run workflow, with the ability to quickly adjust parameters or revert the movement model if major issues are discovered.
- Any larger architectural shifts (for example, a new global physics model) are out of scope for this increment and should be considered separately if needed.

## Follow-up Increments (Optional)

- Introduce terrain-dependent movement and traction, where different surfaces affect speed and turning behavior.
- Define per-tank-type movement profiles (for example, heavier but tougher tanks vs lighter, faster scouts) using the same underlying movement model.
- Explore advanced movement behaviors such as drift, recoil influence from firing, or terrain-based turning penalties.
- Refine enemy behaviors and pathing to take advantage of the new movement model, improving challenge and variety.

## PRD Entry (for docs/PRD.md)

- **Increment ID:** snappy-unified-tank-movement-and-rotation
- **Title:** Snappy Unified Tank Movement and Rotation
- **Status:** Proposed
- **Increment Folder:** docs/increments/snappy-unified-tank-movement-and-rotation/
- **User Story:** As a player, I want all tanks in the game (player and enemies) to use the same smooth, snappy tank-style movement and rotation, so that the controls feel fast, responsive, and consistent across the game.
- **Acceptance Criteria:**
  - All tanks use linear acceleration toward bounded maximum speeds instead of instant full speed.
  - All tanks rotate with a bounded, smooth turn rate around their visual center.
  - Movement is independent of terrain for this increment, and basic controllability is at least as good as before.
- **Use Case Summary:** Player and AI tanks share a unified movement model with acceleration toward capped speeds and bounded rotation around each tank’s center; both respond quickly and predictably to input and AI commands, producing a fast-paced but consistent tank feel without terrain-dependent behavior yet.
