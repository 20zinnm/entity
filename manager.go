package entity

import (
	"sync/atomic"
	"sync"
)

// ID represents an identifier for a unique entity in a manager. It is aliased to uint64 to provide compile-time optimization and to reduce the number of type casts necessary when, for example, serializing it.
type ID = uint64

// Manager ties together entities and systems. It is thread-safe; all methods can be called concurrently.
type Manager struct {
	systemsMu sync.Mutex
	systems   []System
	removeMu  sync.Mutex
	remove    []ID
	// id should only be modified atomically to prevent race conditions.
	id uint64
}

// AddSystem adds a system to the manager to be executed on subsequent calls to update. Systems are executed in the order in which they are added.
func (m *Manager) AddSystem(system System) {
	m.systemsMu.Lock()
	m.systems = append(m.systems, system)
	m.systemsMu.Unlock()
}

// RemoveSystem removes a system from the manager. It will not execute during the next call to update.
func (m *Manager) RemoveSystem(system System) {
	m.systemsMu.Lock()
	for i, s := range m.systems {
		if s == system {
			m.systems = append(m.systems[:i], m.systems[i+1:]...)
			// don't break because a system could theoretically be added multiple times
		}
	}
	m.systemsMu.Unlock()
}

// Systems returns all added systems in the order in which they are executed.
func (m *Manager) Systems() []System {
	return m.systems
}

// Update synchronously executes system updates and removes marked entities.
func (m *Manager) Update(delta float64) {
	m.removeMu.Lock()
	var remove = make([]ID, len(m.remove))
	copy(remove, m.remove)
	m.remove = nil
	m.removeMu.Unlock()
	m.systemsMu.Lock()
	for _, system := range m.systems {
		for _, entity := range remove {
			system.Remove(entity)
		}
		system.Update(delta)
	}
	m.systemsMu.Unlock()
}

// Remove marks an ID for removal from all systems. The entity will actually be removed before the next update.
func (m *Manager) Remove(entity ID) {
	m.removeMu.Lock()
	m.remove = append(m.remove, entity)
	m.removeMu.Unlock()
}

// NewEntity creates a new entity that can be later removed. It provides a synchronized way to add an entity to
func (m *Manager) NewEntity() ID {
	return ID(atomic.AddUint64(&m.id, 1) - 1)
}
