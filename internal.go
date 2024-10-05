package terminalui

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

// initTTY initialises terminal window
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

// startMainLoop initialises program's main loop, controls the terminal size, ensures panes are correctly
// drawn and calls methods attached to their onIterate property
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

// startStdioLoop creates a loop that will get keyboard input
func (t *TUI) startStdioLoop() {
	var b []byte = make([]byte, 1)
	for {
		os.Stdin.Read(b)
		if t.onKeyPress != nil {
			t.onKeyPress(t, b)
		}
	}
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
