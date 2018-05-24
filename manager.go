package entity

import (
	"sync/atomic"
)

// ID represents an identifier for a unique entity in a manager. It is aliased to uint64 to provide compile-time optimization and to reduce the number of type casts necessary when, for example, serializing it.
type ID = uint64

// Manager ties together entities and systems.
type Manager struct {
	systems []System
	// id should only be modified atomically to prevent race conditions.
	id uint64
}

// AddSystem adds a system to the manager to be executed on subsequent calls to update. Systems are executed in the order in which they are added.
func (m *Manager) AddSystem(system System) {
	m.systems = append(m.systems, system)
}

// RemoveSystem removes a system from the manager. It will not execute during the next call to update.
func (m *Manager) RemoveSystem(system System) {
	for i, s := range m.systems {
		if s == system {
			m.systems = append(m.systems[:i], m.systems[i+1:]...)
			// don't break because a system could theoretically be added multiple times
		}
	}
}

// Systems returns all added systems in the order in which they are executed.
func (m *Manager) Systems() (systems []System) {
	return m.systems
}

// Update synchronously executes system updates and removes marked entities.
func (m *Manager) Update(delta float64) {
	for _, system := range m.systems {
		system.Update(delta)
	}
}

// Remove marks an ID for removal from all systems.
func (m *Manager) Remove(entity ID) {
	for _, system := range m.systems {
		system.Remove(entity)
	}
}

// NewEntity creates a new entity that can be later removed. It provides a synchronized way to add an entity to
func (m *Manager) NewEntity() ID {
	return ID(atomic.AddUint64(&m.id, 1) - 1)
}
