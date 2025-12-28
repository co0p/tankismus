package components

import "github.com/co0p/tankismus/pkg/ecs"

// Type IDs used with the generic ECS world.
//
// Note: these values are used as bit positions in a uint64 mask,
// so the total number of distinct component types must stay below 64.
const (
	TypeTransform ecs.ComponentType = iota
	TypeVelocity
	TypePlayerTag
	TypeEnemyTag
	TypeHealth
	TypeSprite
	TypeCollider
	TypeProjectile
)

// Transform represents position, rotation and uniform scale.
type Transform struct {
	X        float64 `json:"x"`
	Y        float64 `json:"y"`
	Rotation float64 `json:"rotation"`
	Scale    float64 `json:"scale"`
}

func (Transform) Type() ecs.ComponentType { return TypeTransform }

// Velocity represents linear and angular velocity.
type Velocity struct {
	VX      float64 `json:"vx"`
	VY      float64 `json:"vy"`
	Angular float64 `json:"angular"`
}

func (Velocity) Type() ecs.ComponentType { return TypeVelocity }

// PlayerTag marks the player entity.
type PlayerTag struct {
	IsPlayer bool `json:"is_player"`
}

func (PlayerTag) Type() ecs.ComponentType { return TypePlayerTag }

// EnemyTag marks enemy entities.
type EnemyTag struct {
	IsEnemy bool `json:"is_enemy"`
}

func (EnemyTag) Type() ecs.ComponentType { return TypeEnemyTag }

// Health represents hit points.
type Health struct {
	Current float64 `json:"current"`
	Max     float64 `json:"max"`
}

func (Health) Type() ecs.ComponentType { return TypeHealth }

// Sprite identifies which sprite to render for an entity.
type Sprite struct {
	SpriteID string `json:"sprite_id"`
}

func (Sprite) Type() ecs.ComponentType { return TypeSprite }

// Collider is a simple axis-aligned bounding box.
type Collider struct {
	Width   float64 `json:"width"`
	Height  float64 `json:"height"`
	OffsetX float64 `json:"offset_x"`
	OffsetY float64 `json:"offset_y"`
}

func (Collider) Type() ecs.ComponentType { return TypeCollider }

// Projectile marks projectile entities.
type Projectile struct {
	Speed    float64 `json:"speed"`
	Lifetime float64 `json:"lifetime"`
	Damage   float64 `json:"damage"`
}

func (Projectile) Type() ecs.ComponentType { return TypeProjectile }
