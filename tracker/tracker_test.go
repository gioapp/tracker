package tracker

import (
	"fmt"
	"testing"
)

var gen1 = MockGenerator{}
var gen2 = MockGenerator{}
var gen3 = MockGenerator{}

var testPattern = Pattern{
	Track{
		Event{1, 1, gen1},
		Event{2, 1, gen1},
		Event{3, 1, gen1},
		Event{4, 1, gen1},
		Event{5, 1, gen1},
	},
	Track{
		Event{64, 127, gen2},
		Event{60, 127, gen2},
		Event{67, 127, gen2},
	},
	Track{
		Event{127, 127, gen3},
	},
}

func TestGetLines(t *testing.T) {
	lines := testPattern.GetLines()
	if len(lines) != 5 {
		t.Error("Expected 5 Lines but actual was %v", len(lines))
	}
}

func TestGetLine(t *testing.T) {
	for i := 0; i < 5; i++ {
		line := testPattern.GetLine(i)
		if len(line) != len(testPattern) {
			t.Error("Expeceted %v Events in Line but actual was %v",
				len(testPattern), len(line))
		}
	}
}

func TestNothing(t *testing.T) {
	// TODO(aoeu): Remove after actually implementing tests.
	fmt.Println(testPattern)
}

func TestNewTrack(t *testing.T) {
	g := MockGenerator{}
	v := 127
	notes := []int{3, 32, 64, 68, 91}
	tr := NewTrack(g, v, notes...)
	for i, e := range tr {
		if e.NoteNum != notes[i] {
			t.Error("Expected note %v but actual was %v", notes[i], e.NoteNum)
		}
		if e.Velocity != v {
			t.Error("Expected velocity %v but actual was %v", v, e.Velocity)
		}
		if e.Generator != g {
			t.Error("Expceted Generator %v but actual was %v", g, e.Generator)
		}
	}
}
