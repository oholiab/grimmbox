package main

import "fmt"
import ui "github.com/nsf/termbox-go"
import "strings"
import "net/http"
import "log"
import "io/ioutil"
import "encoding/json"
import "os"

//import "strconv"

var orange = ui.Attribute(150)
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

func emptyBox(box Box) {
	for i := box.TL.X + 1; i < box.TR.X; i++ {
		for j := box.TL.Y + 1; j < box.BL.Y; j++ {
			ui.SetCell(i, j, ' ', d, d)
		}
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
	rem := len(lines) - height + 2
	if rem > 0 {
		return lines[:height-2], rem
	} else {
		return lines, rem
	}
}

func drawBox(box Box) {
	drawBoxCorners(box)
	drawBoxSides(box)
	emptyBox(box)
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

func render(boxes []Box) {
	boxList[selBoxIndex].RenderColor = selectColor
	for _, box := range boxes {
		drawBox(box)
	}
	ui.Flush()
}

// grimmbox functions above here for separation later
// jira crap

type Ticket struct {
	ID          int
	Key         string
	Summary     string
	Description string
}

type IssueWrapper struct {
	Fields Ticket
}

type Issues struct {
	Issues []IssueWrapper
}

func queryJIRA(jira, user, passwd string) []Ticket {
	var issues Issues
	var tickets []Ticket
	proto := "https://"
	req, err := http.NewRequest("GET", strings.Join([]string{proto, jira, "/rest/api/2/search?jql=assignee=", user, "+AND+status=Open&fields=id,key,summary,description"}, ""), nil)
	req.Header = map[string][]string{
		"Host":         {jira},
		"Content-Type": {"application/json"},
	}
	req.SetBasicAuth(user, passwd)
	req.ProtoMinor = 1
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	json.Unmarshal(body, &issues)
	for _, ticket := range issues.Issues {
		tickets = append(tickets, ticket.Fields)
	}
	return tickets
}

var ticketWidth = 20
var ticketHeight = 5
var tooSmall = false

func get_viewport_width() int {
	termWidth, _ := ui.Size()
	if termWidth < ticketWidth {
		tooSmall = true
		return 0
	}
	return termWidth - ticketWidth
}

func on_exit() {
	fmt.Println("Debug:")
	fmt.Println(strings.Join(debug, "||"))
	fmt.Println(selBoxIndex)
}

func main() {
	defer on_exit()
	err := ui.Init()
	if err != nil {
		panic(err)
	}

	defer ui.Close()

	ui.SetOutputMode(ui.Output256)
	_, h := ui.Size()
	offset := 1
	viewPort := makeBox("grimmwa.re", "", ticketWidth, offset, h-3, get_viewport_width(), orange)
	for i, ticket := range queryJIRA(os.Getenv("JIRA"), os.Getenv("USER"), os.Getenv("JIRAPASS")) {
		boxList = append(boxList, makeBox(ticket.Key, ticket.Summary, 0, i*ticketHeight+offset, ticketHeight, ticketWidth, orange))
	}
	writeln(":PRESS C-c TO EXIT", 0, 0)
	render(append(boxList, viewPort))

loop:
	for {
		boxList[selBoxIndex].RenderColor = boxList[selBoxIndex].LineColor
		switch ev := ui.PollEvent(); ev.Type {
		case ui.EventKey:
			ui.Flush()
			if ev.Key == ui.KeyCtrlC {
				break loop
			} else if ev.Ch == 'j' {
				if uint(len(boxList))-1 == selBoxIndex {
					selBoxIndex = 0
				} else {
					selBoxIndex += 1
				}
			} else if ev.Ch == 'k' {
				if 0 == selBoxIndex {
					selBoxIndex = uint(len(boxList)) - 1
				} else {
					selBoxIndex -= 1
				}
			}
		case ui.EventError:
			panic(ev.Err)
		}
		viewPort.Text = boxList[selBoxIndex].Text
		render(append(boxList, viewPort))
	}

	return
}
