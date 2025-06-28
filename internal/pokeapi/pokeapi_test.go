package pokeapi

import (
	"testing"
	"time"
)

func TestNewClientLocs(t *testing.T) {
	client := NewClient(5 * time.Second)

	if len(client.Config.NextUrl) == 0 {
		t.Errorf("NextUrl expected to be non empty")
	}

	locs, _ := client.GetLocationAreas("back")

	if loclen := len(locs); loclen > 0 {
		t.Errorf("len of locations should be 0, %v found", loclen)
	}

	locs, _ = client.GetLocationAreas("forward")

	if loclen := len(locs); loclen != 20 {
		t.Errorf("len of locations should be 20, %v found", loclen)
	}

}
