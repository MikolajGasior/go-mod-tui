package tui

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

// TUI is main interface definition. It has a name, description, author
// (which are not used anywhere yet), current terminal width and height,
// pointer to main pane, pointer to a function that is triggered when
// interface is being drawn (that happens when app is started and when
// terminal size is changed), and finally pointers to standard output
// and standard error File instances.
type TUI struct {
	name       string
	desc       string
	author     string
	stdout     *os.File
	stderr     *os.File
	h          int
	w          int
	pane       *TUIPane
	onDraw     func(*TUI) int
	onKeyPress func(*TUI, []byte)
	loopSleep  int
}

// GetName returns TUI name
func (t *TUI) GetName() string {
	return t.name
}

// GetDesc returns TUI description
func (t *TUI) GetDesc() string {
	return t.desc
}

// GetAuthor returns TUI author
func (t *TUI) GetAuthor() string {
	return t.author
}

// GetStdout returns stdout property
func (t *TUI) GetStdout() *os.File {
	return t.stdout
}

// GetStderr returns stderr property
func (t *TUI) GetStderr() *os.File {
	return t.stderr
}

// GetPane returns initial/first terminal pane
func (t *TUI) GetPane() *TUIPane {
	return t.pane
}

// GetWidth returns cached terminal width
func (t *TUI) GetWidth() int {
	return t.w
}

// GetHeight returns cached terminal height
func (t *TUI) GetHeight() int {
	return t.h
}

// GetLoopSleep returns delay between each iteration of main loop
func (t *TUI) GetLoopSleep() int {
	return t.loopSleep
}

// Run clears the terminal and starts program's main loop
func (t *TUI) Run(stdout *os.File, stderr *os.File) int {
	t.stdout = stdout
	t.stderr = stderr

	t.initTTY()
	t.clear()

	done := make(chan bool)
	go t.startMainLoop()
	go t.startStdioLoop()
	<-done

	return 0
}

// SetOnDraw attaches function that will be triggered when interface is being
// drawn (what happens on initialisation and terminal resize)
func (t *TUI) SetOnDraw(f func(*TUI) int) {
	t.onDraw = f
}

// SetOnKeyPress attaches function that will triggered when key is pressed (a byte is sent onto stdio)
func (t *TUI) SetOnKeyPress(f func(*TUI, []byte)) {
	t.onKeyPress = f
}

// SetPane sets the main terminal pane
func (t *TUI) SetPane(p *TUIPane) {
	t.pane = p
}

// SetLoopSleep sets the delay between each iteration of main loop
func (t *TUI) SetLoopSleep(s int) {
	t.loopSleep = s
}

// Write prints out on the terminal window at a specified position
func (t *TUI) Write(x int, y int, s string) {
	fmt.Fprintf(t.stdout, "\u001b[1000A\u001b[1000D")
	if x > 0 {
		fmt.Fprintf(t.stdout, "\u001b["+strconv.Itoa(x)+"C")
	}
	if y > 0 {
		fmt.Fprintf(t.stdout, "\u001b["+strconv.Itoa(y)+"B")
	}
	fmt.Fprintf(t.stdout, s)
}

// getSize gets terminal size by calling stty command
func (t *TUI) getSize() (int, int, error) {
	cmd := exec.Command("stty", "size")
	cmd.Stdin = os.Stdin
	out, err := cmd.Output()
	if err != nil {
		return 0, 0, err
	}

	nums := strings.Split(string(out), " ")
	h, err := strconv.Atoi(nums[0])
	if err != nil {
		return 0, 0, err
	}
	w, err := strconv.Atoi(strings.Replace(nums[1], "\n", "", 1))
	if err != nil {
		return 0, 0, err
	}
	return w, h, nil
}

// refreshSize gets terminal size and caches it
func (t *TUI) refreshSize() bool {
	w, h, err := t.getSize()
	if err != nil {
		return false
	}
	if t.w != w || t.h != h {
		t.w = w
		t.h = h

		t.pane.SetWidth(w)
		t.pane.SetHeight(h)
		return true
	}
	return false
}

func (t *TUI) initTTY() {
	cmd1 := exec.Command("stty", "cbreak", "min", "1")
	cmd1.Stdin = os.Stdin
	err := cmd1.Run()
	if err != nil {
		log.Fatal(err)
	}

	cmd2 := exec.Command("stty", "-echo")
	cmd2.Stdin = os.Stdin
	err = cmd2.Run()
	if err != nil {
		log.Fatal(err)
	}
}

// clear clears terminal window
func (t *TUI) clear() {
	fmt.Fprintf(t.stdout, "\u001b[2J\u001b[1000A\u001b[1000D")
}

// startMainLoop initialises program's main loop, controls the terminal size,
// ensures panes are correctly drawn and calls methods attached to their
// onIterate property
func (t *TUI) startMainLoop() {
	for {
		sizeChanged := t.refreshSize()
		if sizeChanged {
			t.clear()
			if t.onDraw != nil {
				t.onDraw(t)
			}
			t.pane.Draw()
		}
		t.pane.Iterate()
		time.Sleep(time.Millisecond * time.Duration(t.loopSleep))
	}
}

func (t *TUI) startStdioLoop() {
	var b []byte = make([]byte, 1)
	for {
		os.Stdin.Read(b)
		if t.onKeyPress != nil {
			t.onKeyPress(t, b)
		}
	}
}

func (t *TUI) Exit(i int) {
	t.clear()
	cmd := exec.Command("stty", "sane")
	cmd.Stdin = os.Stdin
	cmd.Run()
	os.Exit(i)
}

// NewTUI creates new instance of TUI and returns it
func NewTUI(n string, d string, a string) *TUI {
	t := &TUI{name: n, desc: d, author: a}
	p := NewTUIPane("main", t)
	t.SetPane(p)
	t.SetLoopSleep(1000)
	return t
}
