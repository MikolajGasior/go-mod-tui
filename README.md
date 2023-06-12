# go-mod-tui

## This project is not maintained

Package `MikolajGasior/go-mod-tui` is meant to simplify printing on terminal window by
specifying boxes and adding static or dynamic content to it. These boxes here
are called panes and they are defined by vertical or horizontal split.
Terminal window is main pane which can be split into another panes, and these
panes can be split again into next ones, and so on... Just like in another
great tool which is tmux.

Pane size can be defined as a percentage or as number of characters. Pane
content can be dynamic and coming from an attached function (in sample code
below it's called a Widget).

Pane can have a border that can be styled by defining what characters should be
used for it (left side, top-left corder, top bar etc.).

Package tui implementation uses ANSI escape codes and so far has been tested on
MacOSX and Linux. It won't work on Windows (probably Cygwin as well).

### Install

Ensure you have your
[workspace directory](https://golang.org/doc/code.html#Workspaces) created and
run the following:

```
go get -u github.com/MikolajGasior/go-mod-tui
```

### Example

See below sample code and a screenshot of its execution.

```
package main

import (
    "os"
    "github.com/MikolajGasior/go-mod-tui"
)

// TUI has onDraw event and a function can be attached to it. onDraw is
// called when TUI is being drawn, eg. when first started or when
// terminal window is resized.
// getOnTUIDraw returns a func that will later be attached in main().
func getOnTUIDraw(n *NTree) func(*tui.TUI) int {
    // It does nothing actually.
    fn := func(c *tui.TUI) int {
        return 0
    }
    return fn
}

// TUIPane has onDraw event and a function can be attached to it. onDraw
// is called when TUI is being drawn, eg. when first started or when
// terminal window is resized.
// getOnTUIPaneDraw returns a func that will later be attached in main().
func getOnTUIPaneDraw(n *NTree, p *tui.TUIPane) func(*tui.TUIPane) int {
    // Func is defined separate in another struct which is called a Widget.
    // This Widget prints out current time. Check the source for more.
    t := tui.NewTUIWidgetSample()
    t.InitPane(p)
    fn := func(x *tui.TUIPane) int {
        return t.Run(x)
    }
    return fn
}

func main() {
    // Create TUI instance
    myTUI := tui.NewTUI("My Project", "Its description", "Author")
    // Attach func to onDraw event
    myTUI.SetOnDraw(getOnTUIDraw(n))

    // Get main pane which we are going to split
    p0 := myTUI.GetPane()

    // Create new panes by splitting the main pane. Split creates two
    // panes and we have to define size of one of them. If it's the
    // left (vertical) or top (horizontal) one then the value is lower than
    // 0 and if it's right (vertical) or bottom (horizontal) then the value
    // should be highter than 0. It can be a percentage of width/height or
    // number of characters, as it's shown below.
    p01, p02 := p0.SplitVertically(-50, tui.UNIT_PERCENT)
    p021, p022 := p02.SplitVertically(-40, tui.UNIT_CHAR)

    p11, p12 := p01.SplitHorizontally(20, tui.UNIT_CHAR)
    p21, p22 := p021.SplitHorizontally(50, tui.UNIT_PERCENT)
    p31, p32 := p022.SplitHorizontally(-35, tui.UNIT_CHAR)

    // Create style instances which will be attached to certain panes
    s1 := tui.NewTUIPaneStyleFrame()
    s2 := tui.NewTUIPaneStyleMargin()

    // Create custom TUIPaneStyle. Previous ones are predefined and come
    // with the package.
    s3 := &tui.TUIPaneStyle{
        NE: "/", NW: "\\", SE: " ", SW: " ", E: " ", W: " ", N: "_", S: " ",
    }

    // Set pane styles.
    p11.SetStyle(s1)
    p12.SetStyle(s1)
    p21.SetStyle(s2)
    p22.SetStyle(s2)
    p31.SetStyle(s3)
    p32.SetStyle(s1)

    // Attach previously defined func to panes' onDraw event. onDraw
    // handler is called whenever pane is being drawn: on start and
    // on terminal window resize.
    p11.SetOnDraw(getOnTUIPaneDraw(n, p11))
    p12.SetOnDraw(getOnTUIPaneDraw(n, p12))
    p21.SetOnDraw(getOnTUIPaneDraw(n, p21))
    p22.SetOnDraw(getOnTUIPaneDraw(n, p22))
    p31.SetOnDraw(getOnTUIPaneDraw(n, p31))
    p32.SetOnDraw(getOnTUIPaneDraw(n, p32))

    // Attach previously defined func to panes' onIterate event.
    // onIterate handler is called every iteration of TUI's main loop.
    // There is a one second delay between every iteration.
    p11.SetOnIterate(getOnTUIPaneDraw(n, p11))
    p12.SetOnIterate(getOnTUIPaneDraw(n, p12))
    p21.SetOnIterate(getOnTUIPaneDraw(n, p21))
    p22.SetOnIterate(getOnTUIPaneDraw(n, p22))
    p31.SetOnIterate(getOnTUIPaneDraw(n, p31))
    p32.SetOnIterate(getOnTUIPaneDraw(n, p32))

    // Run TUI
    myTUI.Run(os.Stdout, os.Stderr)
}
```

![Example](screenshot.png)
