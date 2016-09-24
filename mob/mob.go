package mob

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

/////////////////////////////////////////////////////////////////////// MobClass

// MobClass holds all information for a particular mob class.
// TODO: the transport handler looks like it'll have a lot of duplicated message formatting
// routines for mobs, traps, etc... should there be some sort of "emotable" interface that
// each of those above types will implement, etc.?
type MobClass struct {
	// The class name of the mobile.
	Name string `json:"class"`

	// A utf32 character code point for an emoji visual display of the mob,
	Avatar int32 `json:"avatar"`

	// Flavour text for the mobile.
	Desc string `json:"desc"`

	// MaxHP is how many hit points the Mob can have
	MaxHP int `json:"max_hp"`
}

//Instance returns a new mob of the given type.
func (mc *MobClass) Instance(name string) *Mob {
	return &Mob{
		Name:  name,
		Class: mc,
		CurHP: mc.MaxHP,
	}
}

/////////////////////////////////////////////////////////////////////////// Mob

// Mob is an instance of a MobClass.
type Mob struct {
	// The name of the individual mob.
	Name string

	// Class is the mob's class.
	Class *MobClass

	// CurrHP is how many hit points the Mob currently has.
	CurHP int
}

// healthiness produces an adjective describing the health of the Mob.
// TODO: maybe pull this into a json, idk
func (m *Mob) healthiness() string {
	r := float32(m.CurHP) / float32(m.Class.MaxHP)
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
	return fmt.Sprintf("%q %s is here and looks %s.", rune(m.Class.Avatar), m.Class.Name, m.healthiness())
}

///////////////////////////////////////////////////////////////////////////////

// loadMobClasses reads a Readable to consume a JSON array of mob class objects and returns
// them in unmarshalled form.
func loadMobClasses(r io.Reader) ([]*MobClass, error) {
	var mobClasses []*MobClass

	decoder := json.NewDecoder(r)
	if err := decoder.Decode(&mobClasses); err != nil {
		return nil, err
	}

	return mobClasses, nil
}

// Load parses a JSON file containing mob structure and produces all the mobile classes, keyed
// on the class name.
func Load(path string) (map[string]*MobClass, error) {
	fd, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer fd.Close()

	// decode the supplied JSON mob data...
	mobClasses, err := loadMobClasses(fd)
	if err != nil {
		return nil, fmt.Errorf("Error parsing mob class description file %v: %v", path, err)
	}

	// ...and construct the mapping of room name to room structures.
	mobClassMap := make(map[string]*MobClass)
	for _, c := range mobClasses {
		mobClassMap[c.Name] = c
	}
	return mobClassMap, nil
}
