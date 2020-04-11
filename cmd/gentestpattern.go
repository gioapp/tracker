package main

import (
	"encoding/gob"
	"flag"
	"fmt"
	"github.com/gioapp/tracker/tracker"
	"log"
	"os"
)

func makeTestPattern() tracker.Pattern {
	var gen1 = tracker.MockGenerator{}
	var gen2 = tracker.MockGenerator{}
	var gen3 = tracker.MockGenerator{}

	return tracker.Pattern{
		tracker.Track{
			&tracker.Event{1, 127, gen1},
			&tracker.Event{4, 127, gen1},
			&tracker.Event{4, 127, gen1},
			&tracker.Event{1, 127, gen1},
			&tracker.Event{2, 127, gen1},
			&tracker.Event{4, 127, gen1},
			&tracker.Event{4, 127, gen1},
			&tracker.Event{1, 127, gen1},
		},
		tracker.Track{
			&tracker.Event{0, 127, gen2},
			&tracker.Event{2, 127, gen2},
			&tracker.Event{3, 127, gen2},
		},
		tracker.Track{
			&tracker.Event{127, 127, gen3},
		},
	}
}

func main() {
	var filePath string
	flag.StringVar(&filePath, "out", "testpattern.trkr", "The file to write the encoded test pattern to.")
	flag.Parse()

	// Use OpenFile so we can clobber existing files by setting the relevant syscall flag.
	out, err := os.OpenFile(filePath, os.O_WRONLY, 0600)
	if err != nil {
		fmt.Println("here")
		log.Fatal(err)
	}
	defer out.Close()
	gob.Register(tracker.Pattern{})
	gob.Register(tracker.MockGenerator{})

	enc := gob.NewEncoder(out)
	pat := makeTestPattern()
	err = enc.Encode(pat)
	if err != nil {
		fmt.Println("there")
		log.Fatal(err)
	}

}
