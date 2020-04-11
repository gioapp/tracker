package tracker

import (
	"errors"
	"fmt"
	"github.com/nsf/termbox-go"
	"regexp"
	"strconv"
	"strings"
)

type dir int

const (
	UP dir = iota
	DOWN
	LEFT
	RIGHT
)

func New() (*Tracker, error) {
	if err := termboxInit(); err != nil {
		return &Tracker{}, err
	}
	// TODO(aoeu): Don't hardcode the pattern name, use flags.
	t, err := NewTracker("cmd/testpattern.trkr")
	if err != nil {
		return &Tracker{}, err
	}
	t.screen.printThings()
	return t, nil
}

type screen struct {
	fg, bg         termbox.Attribute
	editMode       bool
	cX, cY         int
	redraw         chan bool
	currentPattern *Pattern // TODO(aoeu): Do we need the full PatternTable?
	lineOffset     int
}

func NewScreen() *screen {
	return &screen{
		fg:         termbox.ColorDefault,
		bg:         termbox.ColorDefault,
		editMode:   false,
		redraw:     make(chan bool),
		lineOffset: -1,
	}
}

func (t *Tracker) Run() {
	t.screen.refresh()
	t.UserIn()
}

func (t *Tracker) Exit() {
	termbox.Close()
}

func (s *screen) printThings() {
	s.printEditMode()
	s.drawTable()
}

func (s *screen) drawTable() {
	switch {
	case s.editMode:
		s.drawPattern(5, 5, s.cY, s.cX, *s.currentPattern)
	default:
		s.drawPattern(5, 5, s.cY, -1, *s.currentPattern)
	}
}

func (s *screen) moveC(d dir) {
	maxX, maxY := len(*s.currentPattern)-1, s.currentPattern.maxTrackLen()-1
	if s.editMode {
		switch d {
		case UP:
			if s.cY > 0 {
				s.cY--
			}
		case DOWN:
			if s.cY < maxY {
				s.cY++
			}
		case RIGHT:
			if s.cX < maxX {
				s.cX++
			}
		case LEFT:
			if s.cX > 0 {
				s.cX--
			}
		}
	}
}

const (
	editMode     = "EDIT MODE"
	playbackMode = "Press 'e' to edit."
)

func (s *screen) printEditMode() {
	if s.editMode {
		s.prints(1, 1, editMode)
	} else {
		s.prints(1, 1, playbackMode)
	}
}

func termboxInit() error {
	if err := termbox.Init(); err != nil {
		return err
	}
	termbox.SetInputMode(termbox.InputEsc)
	return nil
}

//	Prints text to the screen
func (s screen) prints(x, y int, n interface{}) {
	fmtString := fmt.Sprintf("%%%vd", numWidth)
	switch n.(type) {
	case int:
		s.drawString(x, y, fmt.Sprintf(fmtString, n))
	default:
		s.drawString(x, y, fmt.Sprint(n))
	}
}

func (s *screen) refresh() {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	s.printThings()
	termbox.Flush()
}

var spacesRe = regexp.MustCompile(`\s+`)

func (s *screen) editCell() error {
	if !s.editMode {
		return nil
	}
	input := s.NewEditBox(20, 20, "NoteNumber Velocity")
	var params [2]int
	input = spacesRe.ReplaceAllString(input, " ")
	vals := strings.Split(strings.TrimSpace(input), " ")
	if len(vals) < 2 {
		return errors.New(fmt.Sprintf("Invalid input when editing cell: '%v'", input))
	}
	i := 0
	for _, val := range vals {
		if n, err := strconv.Atoi(string(val)); err != nil {
			return errors.New(fmt.Sprintf("Could not set use value '%v' from input '%v' to edit cell.", val, input))
		} else {
			params[i] = n
			i++
		}
	}
	e := &Event{NoteNum: params[0], Velocity: params[1]}
	(*s.currentPattern).InsertAt(s.cX, s.cY, e)
	return nil
}

func (t *Tracker) UserIn() {
	keyEvents := make(chan termbox.Event)

	go func() {
		for {
			e := termbox.PollEvent()
			switch e.Type {
			case termbox.EventKey:
				keyEvents <- e
			}
		}
	}()

	for {
		select {
		case <-t.screen.redraw:
			t.screen.refresh()
		case e := <-keyEvents:
			switch e.Key {
			case termbox.KeyEsc:
				return
			case termbox.KeyArrowUp:
				t.screen.moveC(UP)
			case termbox.KeyArrowDown:
				t.screen.moveC(DOWN)
			case termbox.KeyArrowRight:
				t.screen.moveC(RIGHT)
			case termbox.KeyArrowLeft:
				t.screen.moveC(LEFT)
			case termbox.KeyEnter:
				t.screen.editCell() //TODO(aoeu): Log any error returned from func editCell.
			}
			switch e.Ch {
			case 'e':
				if t.isPlaying {
					t.TogglePlayback()
				}
				if t.screen.editMode {
					t.screen.editMode = false
					t.screen.refresh()
				} else {
					t.screen.editMode = true
					// TODO(aoeu): t.Stop() ?
				}
			case 'p':
				t.screen.editMode = false
				t.screen.refresh()
				t.TogglePlayback()
			}
			t.screen.refresh()
		}
	}
}

func (s screen) drawString(x, y int, str string) {
	for i, r := range str {
		s.drawChar(x+i, y, r)
	}
}

func (s screen) drawChar(x, y int, r rune) {
	termbox.SetCell(x, y, r, s.fg, s.bg)
}

const (
	numWidth   = 4
	trackWidth = 8
)

func (s *screen) drawPattern(x, y, hr, hc int, p Pattern) {
	defer func() { //TODO(Brad) - Don't make this a deferred func
		s.fg = termbox.ColorBlue
		s.bg = termbox.ColorDefault
	}()
	for i, l := range p.GetLines() {
		s.fg = termbox.ColorBlue
		s.bg = termbox.ColorDefault
		s.prints(x, y+numWidth+i, i)
		for z, e := range l {
			if i == hr && z == hc || i == s.lineOffset {
				s.bg = termbox.ColorRed
			} else {
				s.bg = termbox.ColorDefault
			}
			x1 := x + (z * trackWidth) + 1
			y1 := y + numWidth + i
			s.fg = termbox.ColorDefault
			s.prints(x1+numWidth, y1, e.NoteNum)
			s.fg = termbox.ColorGreen
			s.prints(x1+(numWidth*2), y1, e.Velocity)
		}

	}
	//	Function that calls function drawEvent to draw an event next to the previous event (from left to write)
	//	Note, generator, effect, parameter
}

//	Windows
//	-Draw to windows
//	-Dynamically adds coordinate offsets based on window position
//	-Record table consisting of window IDs
//	-Draw in window like normally drawing in termbox buffer

//	NOTESTSEOTNSET
//	Struct variable indicating "theme:
//	Functions use current "Theme" when printing/drawing

//TODO(Brad) - Add a parameter for entering in current data.
//	This way when the escape key is pressed, no values have been changed.....
func (s screen) NewEditBox(x, y int, title string) (out string) {
	s.drawEditBox(x, y, title)
	termbox.Flush()
	for inChar := s.getInput(); inChar != '\n'; inChar = s.getInput() {
		switch inChar {
		case -1:
			return ""
		case -2:
			if len(out) > 0 {
				out = out[:len(out)-1]
			}
		default:
			out += string(inChar)
		}
		//TODO - Add clearing to whitespace characters deleted
		s.printToEditBox(x, y, out)
		termbox.Flush()
	}
	return
}

func (s screen) printToEditBox(x, y int, in string) {
	s.drawString(x+2, y+3, in)
	s.drawString(x+2+len(in), y+3, " ")
}

func (s screen) drawEditBox(x, y int, title string) {
	s.drawString(x+3, y+1, title)
	for col := x; col < x+10+len(title); col++ {
		s.drawChar(col, y, '-')
		s.drawChar(col, y+4, '-')
	}
	width := 9 + len(title)
	height := 5
	for row := 0; row < height; row++ {
		s.drawChar(x, y+row, '|')
		if row > 2 && row < height-1 {
			for i := 1; i < width; i++ {
				s.drawChar(x+i, y+row, ' ')
			}
		}
		s.drawChar(x+width, y+row, '|')
	}
}

func (s screen) getInput() rune {
	switch e := termbox.PollEvent(); e.Type {
	case termbox.EventKey:
		switch e.Key {
		case termbox.KeyEsc:
			return rune(-1)
		case termbox.KeySpace:
			return ' '
		case termbox.KeyEnter:
			return '\n'
		case termbox.KeyDelete, termbox.KeyBackspace, termbox.KeyBackspace2:
			return -2
		}
		return e.Ch
	}
	return 0
}
