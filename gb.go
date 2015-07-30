package main

import "fmt"
import ui "github.com/nsf/termbox-go"
import "strings"

//import "strconv"

var orange = ui.Attribute(0x005f)
var f, d = orange, ui.ColorDefault
var selectColor = ui.ColorDefault
var debug []string
var selBoxIndex uint = 0

type Coord struct {
	X int
	Y int
}

type Box struct {
	Title       string
	LineColor   ui.Attribute
	RenderColor ui.Attribute
	TL          Coord // top left
	TR          Coord // top right
	BL          Coord // bottom left
	BR          Coord // bottom right
	W           int
	H           int
	Text        string
}

var boxList []Box

func makeBox(title, content string, x, y, height, width int, linecolor ui.Attribute) Box {
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
	box.LineColor = linecolor
	box.RenderColor = linecolor

	return box
}

func drawBoxCorners(box Box) {
	ui.SetCell(box.TL.X, box.TL.Y, '┌', box.RenderColor, d)
	ui.SetCell(box.TR.X, box.TR.Y, '┐', box.RenderColor, d)
	ui.SetCell(box.BL.X, box.BL.Y, '└', box.RenderColor, d)
	ui.SetCell(box.BR.X, box.BR.Y, '┘', box.RenderColor, d)
}

func drawBoxSides(box Box) {
	for i := box.TL.X + 1; i < box.TR.X; i++ {
		ui.SetCell(i, box.TL.Y, '─', box.RenderColor, d)
		ui.SetCell(i, box.BL.Y, '─', box.RenderColor, d)
	}
	for j := box.TL.Y + 1; j < box.BL.Y; j++ {
		ui.SetCell(box.TL.X, j, '│', box.RenderColor, d)
		ui.SetCell(box.BR.X, j, '│', box.RenderColor, d)
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
		ui.SetCell(x+i, y, c, d, d)
	}
}

func on_exit() {
	fmt.Printf("thanks for playing\n")
	fmt.Println(strings.Join(debug, "||"))
	fmt.Println(selBoxIndex)
}

func render(boxes []Box) {
	boxList[selBoxIndex].RenderColor = selectColor
	for _, box := range boxList {
		drawBox(box)
	}
	ui.Flush()
}

func main() {
	defer on_exit()
	err := ui.Init()
	if err != nil {
		panic(err)
	}

	defer ui.Close()

	ui.SetOutputMode(ui.Output256)
	// ui.SetInputMode(ui.InputEsc) //default
	w, h := ui.Size()
	ridonk := " Lorem ipsum dolor sit amet, consectetur adipiscing elit. Proin vitae mollis eros. Aenean at nisl laoreet, porta risus in, placerat dolor. Etiam massa augue, consequat efficitur augue pellentesque, viverra facilisis ex. Curabitur tristique nulla eget urna semper vestibulum. Cras eget eros ex. Vestibulum at ante leo. Quisque faucibus vulputate lacus vitae porta. Vivamus nulla erat, elementum id pulvinar sed, dapibus non diam. Nulla varius finibus est, ut laoreet enim venenatis id. Aliquam suscipit lectus neque, sit amet sollicitudin dui accumsan nec. Nulla enim felis, molestie at facilisis vitae, elementum et mauris. Aenean sit amet finibus ligula, ut luctus elit. Nunc sit amet massa varius est sodales eleifend sit amet et quam. Fusce ac scelerisque magna. Pellentesque porttitor rutrum neque vel fermentum.  Nam sem quam, ultrices quis aliquam in, varius a neque. Fusce elit leo, consectetur in eleifend eget, facilisis at mi. Etiam porttitor gravida ipsum et sollicitudin. Integer magna mauris, bibendum ac porttitor interdum, mattis non libero. Duis placerat felis in pellentesque dignissim. Sed vitae rhoncus purus, sit amet pulvinar magna. Donec aliquam efficitur eros, eu mattis purus porttitor sed. Fusce neque elit, imperdiet vel eleifend ut, ornare vitae lectus. Maecenas tempus tortor a ante scelerisque pellentesque. Suspendisse ac turpis at nisi dictum rutrum ac nec nisi. Sed non ipsum ac nunc venenatis placerat."
	boxList = append(boxList, makeBox("grimmwa.re", ridonk, 1, 2, h-3, w-3, orange))
	someMsg := "lots of words and stuff"
	boxList = append(boxList, makeBox("stuff", someMsg, 4, 5, 3, len(someMsg)+2, ui.ColorRed))
	writeln(":PRESS C-q TO EXIT", 0, 0)
	render(boxList)

loop:
	for {
		boxList[selBoxIndex].RenderColor = boxList[selBoxIndex].LineColor
		switch ev := ui.PollEvent(); ev.Type {
		case ui.EventKey:
			ui.Flush()
			if ev.Key == ui.KeyCtrlQ {
				break loop
			} else if ev.Ch == 'j' {
				selBoxIndex += 1
			} else if ev.Ch == 'k' {
				if 0 == selBoxIndex {
					selBoxIndex = uint(len(boxList)) - 1
				} else {
					selBoxIndex -= 1
				}
			}
			debug = append(debug, string(ev.Ch))
			debug = append(debug, string(selBoxIndex))
		case ui.EventError:
			panic(ev.Err)
		}
		render(boxList)
	}

	return
}
