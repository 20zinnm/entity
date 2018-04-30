package entity_test

import (
	"testing"
	"github.com/20zinnm/entity"
)

type testSystem struct {
	updateCb func(dt float64)
	removeCb func(id entity.ID)
}

func (s *testSystem) Update(delta float64) {
	if s.updateCb != nil {
		s.updateCb(delta)
	}
}

func (s *testSystem) Remove(e entity.ID) {
	if s.removeCb != nil {
		s.removeCb(e)
	}
}

func TestManager_AddSystem(t *testing.T) {
	manager := new(entity.Manager)
	system := &testSystem{
	}
	manager.AddSystem(system)
	systems := manager.Systems()
	if len(systems) != 1 {
		t.Fail()
	}
	if systems[0] != system {
		t.Fail()
	}
}

func TestManager_RemoveSystem(t *testing.T) {
	manager := new(entity.Manager)
	system := &testSystem{}
	manager.AddSystem(system)
	manager.RemoveSystem(system)
	if len(manager.Systems()) > 0 {
		t.Fail()
	}
}

func TestManager_Update(t *testing.T) {
	manager := new(entity.Manager)
	updated := false
	system := &testSystem{
		updateCb: func(dt float64) {
			if dt != 1 {
				t.Fail()
			}
			updated = true
		},
	}
	manager.AddSystem(system)
	manager.Update(1)
	if !updated {
		t.Fail()
	}
}

func TestManager_Remove(t *testing.T) {
	manager := new(entity.Manager)
	id := manager.NewEntity()
	removed := false
	system := &testSystem{
		removeCb: func(rm entity.ID) {
			if rm != id {
				t.Fail()
			}
			removed = true
		},
		updateCb: func(_ float64) {}, // todo: check that entity is removed before update
	}
	manager.AddSystem(system)
	manager.Remove(id)
	manager.Update(1)
	if !removed {
		t.Fail()
	}
}
