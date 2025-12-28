# tankismus

The followup AI supported iteration building a top down panzer game in golang, using the ebiten engine and ecs pattern.

## game design docs
### core concept

A fast paced top down action showcase game for the developer to tackle most problems in game design. 

### design pillars

The design revolves around the following design pillars
- **minimal ui** No real ui elements should be visible other than
- **immersion** through fullscreen, screen shake, sound and light effects
- **retro** its pixel graphics and retro sounds
- **action** there are lots of explosions
- **escalation** the longer you play the hard and more extreme the game becomes
- **endurance** for playing the game longer, the player gets rewards

### features & mechanics

The player controls a single tank from a camera that remains centered on the vehicle. The tank can move forward and backward and rotate freely along a full 360° axis.

The core game mode is a Survival / Horde mode. The player must withstand increasingly difficult waves of enemies for as long as possible, aiming to survive and achieve a high score or longest time survived.

The game map contains various terrain types and obstacles:

- Obstacles: Solid objects that block movement unless destroyed by the tank’s weapon fire.
- Terrain Effects on Movement Speed:
  - Grass: normal movement speed
  - Street/Road: increased movement speed
  - Sand: reduced movement speed
  - Water: no movement possible

Enemy forces consist of continuous waves of hostile tanks that spawn outside the player’s visible area. As the survival session progresses, enemy tanks scale in difficulty and behavior complexity, becoming more durable, more numerous, and more tactically intelligent.

Each enemy tank is defined by:

 - A field of view that determines when it detects the player
 - A health value
 - A shooting cooldown
 - A shooting accuracy parameter

The survival run ends when the player’s tank health is fully depleted.

### interface & controls
The tank is controlled using WASD keys

### visual / audio style
Using pixel art and retro style sound effects

## implementation details

### package layout
tries to follow https://www.ardanlabs.com/blog/2017/02/design-philosophy-on-packaging.html

### architecture
tried to follow the ECS pattern

## references
 - https://github.com/brotherhood-of-recursive-descent/tankism
 - https://ebitengine.org/
 - https://co0p.github.io/posts/ecs-animation/
 - https://gameprogrammingpatterns.com/contents.html
 - https://kenney.nl/assets/top-down-tanks-redux
