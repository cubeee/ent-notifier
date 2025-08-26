package lib_test

import (
	"github.com/cubeee/ent-notifier/lib"
	"testing"
)

func TestGetMappedLocation(t *testing.T) {
	mappedLocations := []*lib.MappedLocation{
		{
			Name:   "Lumbridge",
			X:      3222,
			Y:      3222,
			Radius: 50,
		},
		{
			Name:   "Edgeville",
			X:      3093,
			Y:      3494,
			Radius: 50,
		},
	}
	mapped := lib.GetMappedLocation(3200, 3200, mappedLocations)
	if mapped != mappedLocations[0] {
		t.Errorf("got %v, want %v", mapped, mappedLocations[0])
	}
}

func TestGetClosestOverlappingLocation(t *testing.T) {
	mappedLocations := []*lib.MappedLocation{
		{
			Name:   "Pmahog",
			X:      3303,
			Y:      6125,
			Radius: 20,
		},
		{
			Name:   "Pteak",
			X:      3310,
			Y:      6122,
			Radius: 20,
		},
	}
	mapped := lib.GetMappedLocation(3309, 6121, mappedLocations)
	if mapped != mappedLocations[1] {
		t.Errorf("got %v, want %v", mapped, mappedLocations[1])
	}
}
