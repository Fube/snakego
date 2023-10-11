package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"os/signal"
	"snakego/screen"
	"snakego/snake"
	"time"
)

const (
	cls = "\033[2J"
)

func main() {
	colorCodes := map[string]string{
		"red":     "\033[31m",
		"green":   "\033[32m",
		"yellow":  "\033[33m",
		"magenta": "\033[35m",
		"cyan":    "\033[36m",
		"white":   "\033[37m",
	}

	snakeColors := make([]string, 0)
	appleColors := make([]string, 0)

	snakeColorFlag := flag.String("snake-color", "white", "Choose a color: red, green, yellow, magenta, cyan, white, rainbow")
	appleColorFlag := flag.String("apple-color", "white", "Choose a color: red, green, yellow, magenta, cyan, white, rainbow")
	flag.Parse()

	selectedSnakeColor, ok := colorCodes[*snakeColorFlag]
	if !ok && *snakeColorFlag != "rainbow" {
		panic("Cannot find color")
	}

	if *snakeColorFlag == "rainbow" {
		for _, v := range colorCodes {
			snakeColors = append(snakeColors, v)
		}
	} else {
		snakeColors = append(snakeColors, selectedSnakeColor)
	}

	selectedAppleColor, ok := colorCodes[*appleColorFlag]
	if !ok && *appleColorFlag != "rainbow" {
		panic("Cannot find color")
	}

	if *appleColorFlag == "rainbow" {
		for _, v := range colorCodes {
			appleColors = append(appleColors, v)
		}
	} else {
		appleColors = append(appleColors, selectedAppleColor)
	}

	width := 80
	length := 20

	rand.Seed(time.Now().UnixNano())

	exec.Command("stty", "-F", "/dev/tty", "cbreak").Run()
	exec.Command("stty", "-F", "/dev/tty", "-echo").Run()

	s := snake.New()
	apple := s.SpawnApple(length-2, width-2)

	go in(s)

	last := ""

	sigintc := make(chan os.Signal, 1)
	signal.Notify(sigintc, os.Interrupt)

	screenReset()
	scr := screen.New(length, width)
Draw:
	for {
		select {
		case <-sigintc:
			{
				break Draw
			}
		default:
			{
			}
		}

		scr.DrawLimits()

		s.WalkBody(func(c *snake.Coord) {
			scr.WriteColoredAt(c.X, c.Y, 'x', snakeColors[(c.Y+c.X)%len(snakeColors)])
		})

		scr.WriteColoredAt(apple.X, apple.Y, 'o', appleColors[(apple.Y+apple.X)%len(appleColors)])

		cur := scr.Render()
		if last == "" {
			fmt.Print(cur)
		} else {
			writeDiff(last, cur)
		}

		last = cur

		fmt.Printf("\033[%d;%dH", 0, 0)
		time.Sleep(time.Millisecond * 100)
		scr.Clear()

		if s.IsOnApple(apple) {
			s.Grow()
			apple = s.SpawnApple(length-2, width-2)
		}

		err := s.Move()

		if err != nil || s.IsOOB(length-1, width-1) {
			break
		}
	}
	cleanUp()
}

func screenReset() {
	fmt.Print(cls)
	fmt.Printf("\r")
	fmt.Printf("\033[%d;%dH", 0, 0)
}

func cleanUp() {
	screenReset()
	exec.Command("stty", "-F", "/dev/tty", "-cbreak").Run()
	exec.Command("stty", "-F", "/dev/tty", "echo").Run()
}

func writeDiff(l, c string) {
	row := 1
	col := 1

	lgt := len(l)
	if len(c) < lgt {
		lgt = len(c)
	}

	li := 0
	ci := 0
	for li < lgt && ci < lgt {
		if l[li] == '\033' {
			for {
				if l[li] == 'm' {
					li++
					break
				}
				li++
			}
			continue
		}

		if c[ci] == '\033' {
			a := ci
			for {
				if c[ci] == 'm' {
					ci++
					break
				}
				ci++
			}

			fmt.Printf("\033[%d;%dH", row, col)
			fmt.Print(c[a:ci])
			continue
		}

		if l[li] != c[ci] {
			fmt.Printf("\033[%d;%dH", row, col)
			fmt.Print(string(c[ci]))
		}

		if c[ci] == '\n' {
			row++
			col = 1
		} else {
			col++
		}

		ci++
		li++
	}
}

func in(k snake.Keys) {
	b := make([]byte, 16)
	nanos := time.Now().UnixNano() / 1e6
	for {
		os.Stdin.Read(b)

		cnanos := time.Now().UnixNano() / 1e6
		if cnanos-nanos < 100 {
			continue
		}

		nanos = cnanos

		k.Register(b)
		b[0] = 0
		b[1] = 0
		b[2] = 0
	}
}
