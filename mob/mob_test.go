package mob

import ( 
	"reflect"
	"strings"
	"testing"
)

func TestInstance(t *testing.T) {
	var tests = []struct {
		in *MobClass
		out *Mob
	}{ 
		{
			&MobClass {
				Name:  "Test Class",
				Avatar: 1234,
				MaxHP: 1,
				Desc:  "Description!",
			},
			&Mob{
				Class: &MobClass {
					Name:  "Test Class",
					Avatar: 1234,
					MaxHP: 1,
					Desc:  "Description!",
				},
				CurHP: 1,
				Room: nil,
			},
		},
	}

	for _, test := range tests {
		out := test.in.Instance()

		if reflect.DeepEqual(out, test.out) == false {
			t.Errorf("%s.Instance() got %s, expected %s", test.in, out, test.out)
		}
	}
}

func TestLoadMobs(t *testing.T) {
	var tests = []struct {
		in  string
		out []*MobClass
		err error
	}{
		{"[]", []*MobClass{}, nil},
		{`
		[{
			"class": "Test Class",
			"avatar": 1234, 
			"max_hp": 1,
			"desc": "Description!"
		}]`,
		[]*MobClass{
			{
				Name:  "Test Class",
				Avatar: 1234,
				MaxHP: 1,
				Desc:  "Description!",
			},
		}, nil},
	}

	for _, test := range tests {
		r := strings.NewReader(test.in)
		out, err := loadMobClasses(r)
		if err != test.err {
			t.Errorf("loadMobClasses(\"%s\") got %v err, expected %+v err", test.in, err, test.err)
		}

		if reflect.DeepEqual(out, test.out) == false {
			t.Errorf("loadMobClasses(\"%s\") got %s, expected %s", test.in, out, test.out)
		}
	}
}
