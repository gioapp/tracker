package main

import (
	"github.com/gioapp/tracker/tracker"
	"log"
	"os"
	"runtime"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	f, err := os.OpenFile("errors.log", os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	l := log.New(f, "", log.LstdFlags|log.Llongfile)
	if t, err := tracker.New(); err != nil {
		l.Println(err)
	} else {
		defer t.Exit()
		defer func() {
			if r := recover(); r != nil {
				l.Println(r)
			}
		}()
		//t.ApplySampler("/home/tasm/ir/src/tracker/cmd/config/waves.json")
		t.Run()
	}
}
