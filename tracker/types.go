package tracker

// An Event represents a musical event to occur at a point in time.
type Event struct {
	NoteNum  int
	Velocity int
	Generator
}

// A Generator maps to a musical device that can play an Event.
type Generator interface {
	Play(e Event)
	String() string
}

// A Track is a series of Events meant to be played sequentially
// through time on one or more Generators.
type Track []*Event

// A Line is a set of Events meant to be played concurrently
// at a single moment in time on one or more Generators.
type Line []*Event

// A Pattern is a set of Tracks meant to be played
// concurrently through time.
type Pattern []Track

// A PatternTable is a set of patterns to play in sequence.
// Playing back a PatternTable in entiretly may be thought
// of as playing an entire song.
type PatternTable []*Pattern

type Player struct {
	BPM int
	PatternTable
}

type Tracker struct {
	screen    *screen
	Player    *Player
	stop      chan bool
	isPlaying bool
}

// TODO(aoeu): A PatternTable isn't really a "table."  Rename it?
