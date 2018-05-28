package entity

// System represents a process that executes every tick and deals with game state.
// Systems are responsible for keeping track of entity data and ensuring that Remove and Update can be called concurrently. This is usually accomplished through the use of a mutex.
type System interface {
	Update(delta float64)
	Remove(entity ID)
}
