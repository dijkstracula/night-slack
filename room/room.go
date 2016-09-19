package room

import (
	"encoding"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
)

////////////////////////////////////////////////////////////////////// Direction

// Direction is an enum comprising all the ways a mob can move.
// TODO: better suited to be in package mob?
type Direction int

const (
	north Direction = iota
	south
	east
	west
	up
	down
)

// UnmarshalJSON takes a byte string and transforms it into a Direction.
func (d *Direction) UnmarshalJSON(value []byte) error {
	s := strings.Trim(string(value), "\"")

	switch strings.ToLower(s) {
	case "north":
		*d = north
	case "south":
		*d = south
	case "east":
		*d = east
	case "west":
		*d = west
	case "up":
		*d = up
	case "down":
		*d = down
	default:
		return fmt.Errorf("Unknown Direction \"%s\"", s)
	}
	return nil
}

// MarshalText transforms a Direction into its byte representation.
func (d Direction) MarshalText() ([]byte, error) {
	return []byte(d.String()), nil
}

func (d *Direction) String() string {
	switch *d {
	case north:
		return "north"
	case south:
		return "south"
	case east:
		return "east"
	case west:
		return "west"
	case up:
		return "up"
	case down:
		return "down"
	default:
		return "<unknown Direction>"
	}
}

/////////////////////////////////////////////////////////////////////////// Room

var (
	// ErrNoExit is what is returned when a player tries to go in a Direction
	// that doesn't exist in the given room.
	ErrNoExit = fmt.Errorf("There's no obvious exit in that Direction.")

	// ErrInvalidExit is what is returned when a player tries to go in a non-existent Direction.
	ErrInvalidExit = fmt.Errorf("That isn't a Direction you can go.")
)

// Room holds all information about a room in the house.
type Room struct {
	encoding.TextMarshaler

	// The name of a room - the name will form part of the channel name.
	Name string `json:"name"`

	// A textual description of the room - will be returned to "look" commands
	Desc string `json:"desc"`

	// All the connecting rooms to this room.
	// TODO: https://golang.org/pkg/encoding/#TextMarshaler - key should be a `Direction`
	Exits map[string]string `json:"exits"`

	// TODO: items, mobs, traps.

	// Ensures mutual exclusion for mutating operations on the Room.
	mu sync.RWMutex
}

// ExitByDirection produces the name of a room in the supplied Direction, or an error if there
// is no room that way.
func (r *Room) ExitByDirection(s string) (string, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var d Direction
	if err := d.UnmarshalJSON([]byte(s)); err != nil {
		// `s` wasn't a Direction at all
		return "", ErrInvalidExit
	}

	room, ok := r.Exits[s] //TODO: key should be a Direction
	if !ok {
		// `s` wasn't in our list of exits
		return "", ErrNoExit
	}
	return room, nil
}

// Show produces a textual representation of the current state of the room.
// TODO: This ought to be in transport.
func (r *Room) Show() string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var strs []string
	strs = append(strs, r.Name)
	strs = append(strs, r.Desc)

	//TODO
	//for m := range(r.Mobs) {
	// strs = append(strs, m.Show())
	//}

	//TODO
	//for i := range(r.Item) {
	// strs = append(strs, fmt.Printf("There is an %s here.", s)
	//}

	for dir := range r.Exits {
		strs = append(strs, fmt.Sprintf("An exit lies to the %s.", dir))
	}

	return strings.Join(strs, "\n")
}

///////////////////////////////////////////////////////////////////////////////

// loadRoomArray reads a Readable to consume a JSON array of room objects and returns
// them in unmarshalled form.
func loadRoomArray(r io.Reader) ([]*Room, error) {
	var rooms []*Room

	decoder := json.NewDecoder(r)
	if err := decoder.Decode(&rooms); err != nil {
		return nil, err
	}

	return rooms, nil
}

// Load parses a json file containing room information and produces a map
// mapping room names to *Room structures.
func Load(path string) (map[string]*Room, error) {

	fd, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer fd.Close()

	// decode the supplied JSON room data...
	roomArray, err := loadRoomArray(fd)
	if err != nil {
		return nil, fmt.Errorf("Error parsing room description file %v: %v", path, err)
	}

	// ...and construct the mapping of room name to room structures.
	roomMap := make(map[string]*Room)
	for _, room := range roomArray {
		roomMap[room.Name] = room
	}
	return roomMap, nil
}
