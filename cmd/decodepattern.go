package main

import (
	"bytes"
	"encoding/gob"
	"flag"
	"fmt"
	"github.com/gioapp/tracker/tracker"
	"io"
	"io/ioutil"
	"log"
	"strings"
)

func main() {
	var filePath string
	flag.StringVar(&filePath, "in", "testpattern.trkr", "The file to decode and print as a test pattern")
	flag.Parse()
	b, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatal()
	}
	gob.Register(tracker.Pattern{})
	gob.Register(tracker.MockGenerator{})
	dec := gob.NewDecoder(bytes.NewReader(b))
	pat := tracker.Pattern{}
	err = dec.Decode(&pat)
	if err != nil && err != io.EOF {
		log.Fatal()
	}
	s := fmt.Sprintf("%#v", pat)
	// TODO(aoeu): Is there something like json.Ident for dumping go-formatted values?
	s = strings.Replace(s, ",", ",\n", -1)
	s = strings.Replace(s, "{", "{\n ", -1)
	fmt.Println(s)
}
