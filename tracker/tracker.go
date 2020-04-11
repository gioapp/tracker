package tracker

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"github.com/aoeu/audio"
	"io"
	"io/ioutil"
	"time"
)

// A mockGenerator is only intended for testing or debugging.
type MockGenerator struct{}

func (m MockGenerator) Play(e Event)   {}
func (m MockGenerator) String() string { return "Mock generator." }

func (e Event) String() string {
	return fmt.Sprintf("%v %v", e.NoteNum, e.Velocity)
}

type AudioGenerator struct {
	Sampler *audio.Sampler
}

func NewAudioGenerator(filepath string) (*AudioGenerator, error) {
	a := AudioGenerator{}
	if s, err := audio.NewLoadedSampler(filepath); err != nil {
		return &AudioGenerator{}, err
	} else {
		a.Sampler = s
	}
	return &a, nil
}

func (a AudioGenerator) Play(e Event) {
	go a.Sampler.Play(e.NoteNum, float32(e.Velocity)/127.0)
}

func (a AudioGenerator) String() string {
	return "Audio generator." // TODO(aoeu): Return something useful.
}

func NewTrack(g Generator, velocity int, notes ...int) *Track {
	t := make(Track, len(notes))
	for i := 0; i < len(t); i++ {
		t[i] = &Event{Generator: g, NoteNum: notes[i], Velocity: velocity}
	}
	return &t
}

func NewPattern(filePath string) (*Pattern, error) {
	b, err := ioutil.ReadFile(filePath)
	if err != nil {
		return &Pattern{}, err
	}
	// TODO(aoeu): Is there a programmatic way of registering the types?
	gob.Register(Pattern{})
	gob.Register(MockGenerator{})
	dec := gob.NewDecoder(bytes.NewReader(b))
	pat := Pattern{}
	err = dec.Decode(&pat)
	if err != nil && err != io.EOF {
		return &Pattern{}, err
	}
	return &pat, nil
}

func (p Pattern) maxTrackLen() int {
	maxLen := 0
	for _, track := range p {
		if len(track) > maxLen {
			maxLen = len(track)
		}
	}
	return maxLen
}

func (p Pattern) minTrackLen() int {
	// TODO(aoeu): Is this actually needed?
	minLen := int(^uint(0) >> 1)
	for _, track := range p {
		if len(track) < minLen {
			minLen = len(track)
		}
	}
	return minLen
}

// GetLine returns a Line containing the Events
// associated with all of the Tracks at a given
// offest in a Pattern.
//
// If a Pattern is thought of as a "table", and
// Tracks are thought of as "columns", GetLine
// returns "row" of the table.
func (p Pattern) GetLine(offset int) Line {
	l := make(Line, len(p))
	for i, track := range p {
		switch {
		case len(track) > offset:
			l[i] = track[offset]
		default:
			l[i] = &Event{}
		}
	}
	return l
}

// GetLines returns a series of Line types
// containing the Events associated with all of
// the Tracks in a Pattern.
//
// Any Track that is shorter in length than others
// in the pattern is still still represented in a
// respective Line with empty Event values (as padding).
func (p Pattern) GetLines() []Line {
	maxTrackLen := p.maxTrackLen()
	l := make([]Line, maxTrackLen)
	for i := range l {
		l[i] = p.GetLine(i)
	}
	return l
}

func NewPlayer(filepath string) (*Player, error) {
	p := Player{}
	p.PatternTable = make(PatternTable, 1)
	if pattern, err := NewPattern(filepath); err != nil {
		return &Player{}, err
	} else {
		p.PatternTable[0] = pattern
	}
	p.BPM = 120 // TODO(aoeu): Don't hardcode the BPM.
	return &p, nil

}

func (t *Tracker) TogglePlayback() {
	if t.isPlaying {
		t.Stop()
	} else {
		go t.Play()
	}
}

func (t *Tracker) Stop() {
	t.stop <- true
}

func (t *Tracker) Play() {
	t.isPlaying = true
	defer func() {
		t.isPlaying = false
		t.screen.lineOffset = -1
		t.screen.redraw <- true
	}()
	// TODO(aoeu): Can pieces of this be decoupled from Tracker?
	nsPerBeat := 60000000000 / t.Player.BPM
	for _, pattern := range t.Player.PatternTable {
		for _, line := range pattern.GetLines() {
			t.screen.lineOffset += 1
			t.screen.redraw <- true
			for _, e := range line {
				if e.Generator != nil {
					// TODO(aoeu): Reconsider ownership of Events and Generators.
					go e.Generator.Play(*e)
				}
			}
			select {
			case <-time.After(time.Duration(nsPerBeat) * time.Nanosecond):
			case <-t.stop:
				return
			}
		}
	}

}

func (t *Tracker) ApplySampler(samplerConfig string) error {
	g, err := NewAudioGenerator(samplerConfig)
	if err != nil {
		return err
	}
	if err := g.Sampler.Run(); err != nil {
		return err
	}
	for _, pattern := range t.Player.PatternTable {
		pattern.ApplyGenerator(g)
	}
	return nil
}

func (p *Pattern) ApplyGenerator(g Generator) {
	for _, track := range *p {
		track.ApplyGenerator(g)
	}
}

func (t *Track) ApplyGenerator(g Generator) {
	for _, e := range *t {
		e.Generator = g
	}
}

func NewTracker(trkrFilepath string) (*Tracker, error) {
	t := Tracker{}
	if p, err := NewPlayer(trkrFilepath); err != nil {
		return &Tracker{}, err
	} else {
		t.Player = p
	}
	t.screen = NewScreen()
	t.screen.currentPattern = t.Player.PatternTable[0]
	t.stop = make(chan bool)
	return &t, nil
}

func (p *Pattern) InsertAt(x, y int, e *Event) {
	if len(*p) < x {
		return //TODO - Return error ???????
	}

	t := (*p)[x]
	l := len((*p)[x])
	switch {
	case l > y:
		t[y].NoteNum, t[y].Velocity = e.NoteNum, e.Velocity
	case l == y:
		(*p)[x] = append(t, e)
		fmt.Println("yea")
	default:
		k := make([]*Event, y-l)
		for i := range k {
			k[i] = &Event{}
		}
		(*p)[x] = append((*p)[x], k...)
		(*p)[x] = append((*p)[x], e)
	}

}
