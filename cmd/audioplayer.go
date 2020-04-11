package main

import (
	"github.com/aoeu/audio"
	"github.com/gioapp/tracker/tracker"
)

func main() {
	configPath := "/home/tasm/ir/src/tracker/cmd/config/waves.json"
	sampler, err := audio.NewLoadedSampler(configPath)
	g := tracker.MockGenerator{}
	p := tracker.PatternTable{
		tracker.Pattern{
			tracker.NewTrack(g, 127, []int{64, 60, 67}),
			tracker.NewTrack(g, 127, []int{52, 48, 55}),
			tracker.NewTrack(g, 127, []int{40, 36, 42}),
		},
	}
	t := tracker.Tracker{BPM: 120, PatternTable: p}
	t.Play()
}
