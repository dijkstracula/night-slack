package room

import (
	"reflect"
	"strings"
	"testing"
)

func TestExitByDirection(t *testing.T) {
	var tests = []struct {
		in  string
		out string
		err error
	}{
		{"north", "Room2", nil},
		{"south", "", ErrNoExit},
		{"sdlkfjsd", "", ErrInvalidExit},
	}

	r := Room{
		Name:  "Room1",
		Desc:  "TODO",
		Exits: map[string]string{"north": "Room2"},
	}

	for _, test := range tests {
		out, err := r.ExitByDirection(test.in)
		if out != test.out {
			t.Errorf("Room1.ExitByDirection(\"%s\") got \"%+v\", expected \"%+v\"", test.in, out, test.out)
		}
		if err != test.err {
			t.Errorf("Room1.ExitByDirection(\"%s\") got err \"%+v\", expected err \"%+v\"", test.in, err, test.err)
		}
	}
}

func TestUnmarshalDirection(t *testing.T) {
	var tests = []struct {
		in  []byte
		out Direction
	}{
		{[]byte("north"), north},
		{[]byte("south"), south},
		{[]byte("east"), east},
		{[]byte("west"), west},
		{[]byte("up"), up},
		{[]byte("down"), down},
		{[]byte("NoRtH"), north}, // case insensitivity
	}

	for _, test := range tests {
		var out Direction
		if err := out.UnmarshalJSON(test.in); err != nil {
			t.Errorf("Direction.UnmarshalJSON(\"%s\") got non=nil \"%v\" error", test.in, err)
		}
		if out != test.out {
			t.Errorf("Direction.UnmarshalJSON(\"%s\") got %d, expected %d", test.in, out, test.out)
		}
	}
}

func TestShow(t *testing.T) {
	r := Room{
		Name:  "Room1",
		Desc:  "TODO",
		Exits: map[string]string{"north": "Room2"},
	}

	txt := r.Show()

	if !strings.Contains(txt, "TODO") {
		t.Errorf("Missing desription in room formatted output")
	}

	if !strings.Contains(txt, "north") {
		t.Errorf("Missing exit in room formatted output")
	}

	if strings.Contains(txt, "%%s") {
		t.Errorf("\"%%s\" in room formatted output")
	}
}

func TestLoadRoomArray(t *testing.T) {
	var tests = []struct {
		in  string
		out []*Room
		err error
	}{
		{"[]", []*Room{}, nil},
		{`[ { "Name": "Room1",
		"Desc": "TODO",
		"Exits": { "north": "Room2" }
	}
	]`,
			[]*Room{
				{
					Name:  "Room1",
					Desc:  "TODO",
					Exits: map[string]string{"north": "Room2"},
				},
			}, nil},
	}

	for _, test := range tests {
		r := strings.NewReader(test.in)
		out, err := loadRoomArray(r)
		if err != test.err {
			t.Errorf("loadRoomArray(\"%s\") got %v err, expected %+v err", test.in, err, test.err)
		}

		if reflect.DeepEqual(out, test.out) == false {
			t.Errorf("loadRoomArray(\"%s\") got %s, expected %s", test.in, out, test.out)
		}
	}
}
