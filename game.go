package main

import (
	"fmt"

	"github.com/dijkstracula/night-slack/mob"
	"github.com/dijkstracula/night-slack/room"
)

// idk some kind of event loop that parses actions, updates the house, and sends
// messages back out to the transport.
//
// TODO: add Game as an arg to all actions
type Game struct {
	House House
}

////////////////////////////////////////////////////////////////////////////////

// The House is the entire state of game of night-slack. The House contains
// the current state of the world (rooms, mobs, traps, etc).
//
// The rest of the world "plays" night-slack by manipulating a single house with
// Actions. Both scripted/timed events and player actions are turned into
// Actions and then applied to a House.
//
// In fact, because night-slack is weird, there's no concept of players, just
// what mobs are alive in the house.
//
// The game is structured such that it expects to be wrapped in an event loop
// that handles parsing input into Actions and sends them to Tick at an
// appropriate interval.
//
// A House is not safe for access from multiple goroutines.
//
// TODO: Figure out what "winning" and "losing" mean and have that checked after
// every call to Tick.
type House struct {
	// The rooms in the house, by name.
	rooms map[string]*room.Room

	// The mobs in the house, by name.
	mobs map[string]*mob.Mob

	// The current locations of every mob in the house.
	locations map[*mob.Mob]*room.Room
}

// Tick accepts a batch of action and applies them all in order.
//
// TODO: Check to see if the game is still going
func (g *House) Tick(actions []Action) {
	for _, action := range actions {
		action.DoIt(g)
	}
}

////////////////////////////////////////////////////////////////////////////////

// An action is a function on the current state of the game. It may modify the
// game, and can safely assume that no other Action is happening while it does
// so.
type Action interface {
	DoIt(*House) error
}

//  Move a mob around
//
// Since there is no structure to rooms, a mob can move instantly from one room
// to another. This doesn't make much sense in our reality, but in the reality
// of night-slack it's perfectly normal.
type Move struct {
	MobName         string
	DestinationName string
}

func (m *Move) DoIt(g *House) error {
	mob, mobExists := g.mobs[m.MobName]
	if !mobExists {
		// FIXME: there should be some difference between an error in something
		// a player does and something we messed up while scripting mob movements.
		return fmt.Errorf("hey! %s isn't a mob", m.MobName)
	}

	room, roomExists := g.rooms[m.DestinationName]
	if !roomExists {
		// FIXME: there should be some difference between an error in something
		// a player does and something we messed up while scripting mob movements.
		return fmt.Errorf("hey! %s isn't a room", m.DestinationName)
	}

	g.locations[mob] = room
	return nil
}

// Describe a mob or a room by name
//
// TODO: are there things that can change when they're described?!?
type Describe struct {
	Thingy string
}

func (d Describe) DoIt(g *House) {
	if room, ok := g.rooms[d.Thingy]; ok {
		g.ShowMessage(room.Show())
		return
	}

	if mob, ok := g.mobs[d.Thingy]; ok {
		g.ShowMessage(mob.Show())
		return
	}
}

// type AttackMob struct { ... }
// type DisarmTrap struct { ... }
// etc
