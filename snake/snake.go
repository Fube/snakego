package snake

import (
	"sync"
	"math/rand"
	"errors"
)

type keys struct {
	sync.Mutex
	up bool
	right bool
	down bool
	left bool
}

type Keys interface {
	IsUp() bool
	IsDown() bool
	IsLeft() bool
	IsRight() bool
	Register(b []byte)
}

func (k *keys) IsUp() bool {
	return k.up
}
func (k *keys) IsRight() bool {
	return k.right
}
func (k *keys) IsDown() bool {
	return k.down
}

func (k *keys) IsLeft() bool {
	return k.left
}

func (k *keys) Register(b []byte) {
	k.Lock()
	defer k.Unlock()

	if b == nil || len(b) < 3 || b[0] != 27 || b[1] != 91 {
		return
	}

	if b[2] == 65 && !k.down {
		k.up = true
		k.down = false
		k.left = false
		k.right = false
	} else if b[2] == 67 && !k.left {
		k.right = true
		k.down = false
		k.left = false
		k.up = false
	} else if b[2] == 66 && !k.up {
		k.down = true
		k.up = false
		k.left = false
		k.right = false
	} else if b[2] == 68 && !k.right {
		k.left = true
		k.down = false
		k.up = false
		k.right = false
	}

}

type Coord struct {
	X int
	Y int
}

type snake struct {
	*keys
	body []*Coord
	grow int
}

func New() Snake {
	k := &keys {down: true}
	s := snake{
		keys: k,
	}

	s.body = append(s.body, &Coord{X: 3, Y: 3})
	s.body = append(s.body, &Coord{X: 4, Y: 3})
	s.body = append(s.body, &Coord{X: 5, Y: 3})

	return &s
}

type Snake interface {
	Keys
	Move() error
	IsOnApple(*Coord) bool
	Grow()
	IsOOB(length, width int) bool
	SpawnApple(length, width int) *Coord
	WalkBody(func (*Coord))
}

func (s *snake) WalkBody(w func (*Coord)) {
	for _, c := range s.body {
		w(c)
	}
}

func (s *snake) Move() error {
	s.keys.Lock()
	defer s.keys.Unlock()

	last := s.body[len(s.body) - 1]
	s.body[0].X = last.X
	s.body[0].Y = last.Y

	g := false

	if s.grow <= 0 {
		s.grow = 0
	} else {
		cp := *(s.body[0])
		s.body = append([]*Coord{&cp}, s.body...)
		s.grow--
		g = true
	}


	if s.keys.IsUp() {
		s.body[0].Y = max(s.body[0].Y - 1, 0)
	}

	if s.keys.IsRight() {
		s.body[0].X = max(s.body[0].X + 1, 0)
	}

	if s.keys.IsDown() {
		s.body[0].Y = max(s.body[0].Y + 1, 0)
	}

	if s.keys.IsLeft() {
		s.body[0].X = max(s.body[0].X - 1, 0)
	}

	t := s.body[0]
	s.body = s.body[1:]
	s.body = append(s.body, t)

	if g {
		return nil
	}

	for _, c := range s.body {
		for _, c2 := range s.body {
			if c == c2 {
				continue
			}

			if c.X == c2.X && c.Y == c2.Y {
				return errors.New("collided")
			}
		}
	}

	return nil
}

func (s *snake) IsOnApple(a *Coord) bool {
	last := s.body[len(s.body) - 1]
	return last.X == a.X && last.Y == a.Y
}

func (s *snake) Grow() {
	s.grow++
}

func (s *snake) IsOOB(length, width int) bool {
	for _, c := range s.body {
		if c.X <= 0 || c.Y <= 0 || c.Y >= length || c.X >= width {
			return true
		}
	}

	return false
}

func (s *snake) SpawnApple(length, width int) *Coord {
	c := &Coord{Y: rng(3, length), X: rng(3, width)}
	i := 0
	for s.IsOnApple(c) && i <= 500 {
		c = &Coord{Y: rng(3, length), X: rng(3, width)}
		i++
	}

	if i <= 500 {
		return c
	}

	panic("Cannot spawn apple")
}

func max(a, b int) int {
	if a < b {
		return b
	}

	return a
}

func rng(min, max int) int {
	return rand.Intn(max-min) + min
}
