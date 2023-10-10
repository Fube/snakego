package screen

import (
	"strings"
)


type screen struct {
	state [][]rune
	colored [][]string
	width int
	length int
}

func New(length, width int) Screen {
	s := screen{length: length, width: width}

	s.state = make([][]rune, length, length)
	for i := 0; i < length; i++ {
		s.state[i] = make([]rune, width, width)
	}

	s.colored = make([][]string, length, length)
	for i := 0; i < length; i++ {
		s.colored[i] = make([]string, width, width)
	}

	return &s
}

type Screen interface {
	DrawLimits()
	WriteAt(x, y int, c rune)
	WriteColoredAt(x, y int, c rune, clr string)
	Clear()
	Render() string
}

func (s *screen) DrawLimits() {
	for i := 0; i < s.length; i++ {
		s.WriteAt(0, i, '|')
	}

	for i := 0; i < s.length; i++ {
		s.WriteAt(s.width-1, i, '|')
	}


	for i := 0; i < s.width; i++ {
		s.WriteAt(i, 0, '-')
	}

	for i := 0; i < s.width; i++ {
		s.WriteAt(i, s.length - 1, '-')
	}
}

func (s *screen) WriteAt(x, y int, c rune) {
	s.state[y][x] = c
}

func (s *screen) WriteColoredAt(x, y int, c rune, clr string) {
	s.WriteAt(x, y, c)
	s.colored[y][x] = clr
}

func (s *screen) Clear() {
	for y := range s.state {
		for x := 0; x < len(s.state[y]); x++ {
			s.state[y][x] = 0
			s.colored[y][x] = ""
		}
	}
}

func (s *screen) Render() string {
	var bldr strings.Builder
	for y, row := range s.state {
		for x := 0; x < len(row); {
			if v := s.colored[y][x]; v != ""{
				bldr.WriteString(v)

				for x < len(row) && s.colored[y][x] == v {
					bldr.WriteRune(getRenderedRune(s.state[y][x]))
					x++
				}

				bldr.WriteString("\033[0m")
			} else {
				bldr.WriteRune(getRenderedRune(s.state[y][x]))
				x++
			}
		}
		bldr.WriteRune('\n')
	}

	return bldr.String()
}

func getRenderedRune(c rune) rune {
	if c == 0 {
		return ' '
	}

	return c
}
