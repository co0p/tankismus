package ecs

// EntityID is an opaque handle to an entity in the World.
type EntityID int

// ComponentType identifies a particular component kind.
//
// It is used both as a key into component storage and to compute
// bitmasks for fast "has components" queries.
type ComponentType int

// Component is the interface implemented by all components stored in the
// World. Concrete components live in higher-level packages (for example
// game/components) and return their ComponentType via Type().
type Component interface {
	Type() ComponentType
}

// Entity holds metadata about a single entity.
type Entity struct {
	id   EntityID
	mask uint64
}

// World owns entities and generic component storage.
type World struct {
	// nextID is incremented for every new entity.
	nextID EntityID

	entities   map[EntityID]*Entity
	components map[ComponentType]map[EntityID]Component
}

// NewWorld constructs an empty World.
func NewWorld() *World {
	return &World{
		nextID:     1,
		entities:   make(map[EntityID]*Entity),
		components: make(map[ComponentType]map[EntityID]Component),
	}
}

// NewEntity creates a new entity and returns its ID.
func (w *World) NewEntity() EntityID {
	id := w.nextID
	w.nextID++
	w.entities[id] = &Entity{id: id, mask: 0}
	return id
}

// DestroyEntity removes the entity and all its components.
func (w *World) DestroyEntity(id EntityID) {
	delete(w.entities, id)
	for t, store := range w.components {
		delete(store, id)
		w.components[t] = store
	}
}

// bitFor returns the bitmask corresponding to a component type.
func bitFor(t ComponentType) uint64 {
	return 1 << uint(t)
}

// ensureEntity makes sure the entity metadata exists.
func (w *World) ensureEntity(id EntityID) *Entity {
	if e, ok := w.entities[id]; ok {
		return e
	}
	e := &Entity{id: id, mask: 0}
	w.entities[id] = e
	return e
}

// AddComponent attaches a component to an entity.
func (w *World) AddComponent(id EntityID, c Component) {
	e := w.ensureEntity(id)
	t := c.Type()
	store, ok := w.components[t]
	if !ok {
		store = make(map[EntityID]Component)
		w.components[t] = store
	}
	store[id] = c
	e.mask |= bitFor(t)
}

// RemoveComponent detaches a component of the given type from an entity.
func (w *World) RemoveComponent(id EntityID, t ComponentType) {
	if store, ok := w.components[t]; ok {
		delete(store, id)
	}
	if e, ok := w.entities[id]; ok {
		e.mask &^= bitFor(t)
	}
}

// GetComponent returns the component of the given type for an entity, if any.
func (w *World) GetComponent(id EntityID, t ComponentType) (Component, bool) {
	store, ok := w.components[t]
	if !ok {
		return nil, false
	}
	c, ok := store[id]
	return c, ok
}

// HasComponent reports whether the entity has a component of the given type.
func (w *World) HasComponent(id EntityID, t ComponentType) bool {
	e, ok := w.entities[id]
	if !ok {
		return false
	}
	return e.mask&bitFor(t) != 0
}

// Mask returns the current component mask for an entity.
func (w *World) Mask(id EntityID) (uint64, bool) {
	e, ok := w.entities[id]
	if !ok {
		return 0, false
	}
	return e.mask, true
}

// MaskFor computes a mask for the given component types.
func MaskFor(types ...ComponentType) uint64 {
	var m uint64
	for _, t := range types {
		m |= bitFor(t)
	}
	return m
}

// Find returns all entities whose component mask contains all bits in required.
func (w *World) Find(required uint64) []EntityID {
	if required == 0 {
		return nil
	}

	result := make([]EntityID, 0)
	for id, e := range w.entities {
		if e.mask&required == required {
			result = append(result, id)
		}
	}
	return result
}
