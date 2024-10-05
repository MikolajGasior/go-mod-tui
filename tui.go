package terminalui

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
)

// TUI is main interface definition. It has current terminal width and height, pointer to main pane,
// pointer to a function that is triggered when interface is being drawn (that happens when app is
// started and when terminal size is changed), and finally pointers to standard output and standard
// error File instances.
type TUI struct {
	stdout     *os.File
	stderr     *os.File
	h          int
	w          int
	pane       *TUIPane
	onDraw     func(*TUI) int
	onKeyPress func(*TUI, []byte)
	loopSleep  int
}


// NewTUI creates new instance of TUI and returns it
func NewTUI() *TUI {
	t := &TUI{}
	p := NewTUIPane("main", t)
	t.SetPane(p)
	t.SetLoopSleep(1000)
	return t
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
	fmt.Fprint(t.stdout, s)
}

// Exit closed the program
func (t *TUI) Exit(i int) {
	t.clear()
	cmd := exec.Command("stty", "sane")
	cmd.Stdin = os.Stdin
	cmd.Run()
	os.Exit(i)
}
