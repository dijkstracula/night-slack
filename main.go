package main

import (
	"flag"
	"fmt"
	"path"

	"github.com/dijkstracula/night-slack/mob"
	"github.com/dijkstracula/night-slack/room"
	"github.com/golang/glog"
)

////////////////////////////////////////////////////////////////////// Game

// World holds all data for the current game world, as well as global handlers to
// interface with players behind some certain protocol.
type World struct {
	// rooms is a collection of all the rooms in the game, keyed on their name.
	rooms map[string]*room.Room

	// mobs is a collection of all mobs in the game, keyed on their name.
	mobs map[string]*mob.Mob

	// traps is a collection of all items in the game, keyed on their name.
	// traps map[string]*trap.Trap
}

var (
	dataDirPath string
)

func init() {
	flag.StringVar(&dataDirPath, "data_dir_path", "./data/", "Path to game data files")
	flag.Parse()
}

func main() {
	// TODO
	rooms, err := room.Load(path.Join(dataDirPath, "rooms.json"))
	if err != nil {
		glog.Fatalf("Could not load room data: %v", err)
	} else {
		glog.Infof("%d room(s) loaded.", len(rooms))
	}

	mobs, err := mob.Load(path.Join(dataDirPath, "mobs.json"))
	if err != nil {
		glog.Fatalf("Could not load mob data: %v", err)
	} else {
		glog.Infof("%d mob(s) loaded.", len(mobs))
	}

	// TODO
	fmt.Printf("%v\n", rooms["Living Room"].Show())
	fmt.Printf("%v\n", mobs["Dastardly Auger"].Instance("larry").Show())
}
