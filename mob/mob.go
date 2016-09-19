package mob

import (
	"fmt"
	"sync"
)

/////////////////////////////////////////////////////////////////////////// Mob

// Mob (mobile) is an autonomous entity.
// TODO: the transport handler looks like it'll have a lot of duplicated message formatting
// routines for mobs, traps, etc... should there be some sort of "emotable" interface that
// each of those above types will implement, etc.?
type Mob struct {
	// The name of the mobile.
	Name string

	// Flavour text for the mobile.
	Desc string

	// CurHP is how many hit points the Mob currently has
	CurHP int

	// MaxHP is how many hit points the Mob can have
	MaxHP int

	// Room is the room that the Mob is currently in.
	// Room *room.Room

	mu sync.RWMutex
}

// healthiness produces an adjective describing the health of the Mob.
// TODO: maybe pull this into a json, idk
func (m *Mob) healthiness() string {
	r := float32(m.CurHP) / float32(m.MaxHP)
	if r > 0.99 {
		return "healthy"
	} else if r > 0.8 {
		return "slightly injured"
	} else if r > 0.6 {
		return "injured"
	} else if r > 0.4 {
		return "badly injured"
	} else {
		return "nearly dead"
	}
}

// Show produces a one-liner "current status" message for a given mob.
func (m *Mob) Show() string {
	return fmt.Sprintf("%s is here and looks %s.", m.Name, m.healthiness())
}
