package entity

// System represents a process that executes every tick and deals with game state. Systems are responsible for keeping track of entity data.
type System interface {
	Update(delta float64)
	Remove(entity ID)
}
