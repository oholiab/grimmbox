package main

import "fmt"
import "os"
import ioutil "io/ioutil"
import ui "github.com/nsf/termbox-go"
import "strings"
import "strconv"

var f, b = ui.ColorDefault, ui.ColorDefault
var debug []string

type Coord struct {
	X int
	Y int
}

type Box struct {
	Title string
	Color ui.Attribute
	TL    Coord // top left
	TR    Coord // top right
	BL    Coord // bottom left
	BR    Coord // bottom right
	W     int
	H     int
	Text  string
}

func on_exit() {
	fmt.Printf("thanks for playing\n")
	ls()
	fmt.Println("%v\n", strings.Join(debug, "||"))
}

func makeBox(title, content string, x, y, height, width int) Box {
	var box Box
	box.Title = title
	box.TL.X = x
	box.TL.Y = y
	box.TR.X = x + width - 1
	box.TR.Y = y
	box.BL.X = x
	box.BL.Y = y + height - 1
	box.BR.X = x + width - 1
	box.BR.Y = y + height - 1
	box.W = width
	box.H = height
	box.Text = content

	return box
}
func ls() {
	wd, _ := os.Getwd()
	list, _ := ioutil.ReadDir(wd)
	for _, file := range list {
		fmt.Println("%v\n", file.Name())
	}
}

func drawBoxCorners(box Box) {
	ui.SetCell(box.TL.X, box.TL.Y, '┌', f, b)
	ui.SetCell(box.TR.X, box.TR.Y, '┐', f, b)
	ui.SetCell(box.BL.X, box.BL.Y, '└', f, b)
	ui.SetCell(box.BR.X, box.BR.Y, '┘', f, b)
}

func drawBoxSides(box Box) {
	for i := box.TL.X + 1; i < box.TR.X; i++ {
		ui.SetCell(i, box.TL.Y, '─', f, b)
		ui.SetCell(i, box.BL.Y, '─', f, b)
	}
	for j := box.TL.Y + 1; j < box.BL.Y; j++ {
		ui.SetCell(box.TL.X, j, '│', f, b)
		ui.SetCell(box.BR.X, j, '│', f, b)
	}
}

func boxText(text string, width, height int) ([]string, int) {
	var lines []string
	var sliceLength int
	buf := []byte(text)
	for len(buf) > 0 {
		if width-2 < len(buf) {
			sliceLength = width - 2
		} else {
			sliceLength = len(buf)
		}
		lines = append(lines, string(buf[0:sliceLength]))
		if sliceLength == len(buf) {
			break
		} else {
			buf = buf[sliceLength:]
		}
	}
	rem := len(lines) - height
	debug = append(debug, strings.Join(lines, " "))
	debug = append(debug, strconv.FormatInt(int64(len(lines)), 10))
	if rem > 0 {
		return lines[:height], rem
	} else {
		return lines, rem
	}
}

func drawBox(box Box) {
	drawBoxCorners(box)
	drawBoxSides(box)
	writeln(box.Title, box.TL.X+1, box.TL.Y)
	lines, _ := boxText(box.Text, box.W, box.H)
	for i, line := range lines {
		writeln(line, box.TL.X+1, box.TL.Y+1+i)
	}
}

func writeln(s string, x, y int) {
	for i, c := range s {
		ui.SetCell(x+i, y, c, f, b)
	}
}

func main() {
	defer on_exit()
	err := ui.Init()
	if err != nil {
		panic(err)
	}

	defer ui.Close()

	// ui.SetInputMode(ui.InputEsc) //default
	w, h := ui.Size()
	ridonk := "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
	box := makeBox("grimmwa.re", ridonk, 1, 2, h-3, w-3)
	writeln(":PRESS C-q TO EXIT", 0, 0)
	drawBox(box)
	ui.Flush()

loop:
	for {
		switch ev := ui.PollEvent(); ev.Type {
		case ui.EventKey:
			ui.Flush()
			if ev.Key == ui.KeyCtrlQ {
				break loop
			}
		case ui.EventError:
			panic(ev.Err)
		}
	}

	return
}
