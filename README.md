# terminal-ui

[![Go Reference](https://pkg.go.dev/badge/github.com/go-phings/terminal-ui.svg)](https://pkg.go.dev/github.com/go-phings/terminal-ui) [![Go Report Card](https://goreportcard.com/badge/github.com/go-phings/terminal-ui)](https://goreportcard.com/report/github.com/go-phings/terminal-ui)

The `terminalui` package is designed to simplify output to a terminal window by allowing the specification of panes with static or dynamic content. These panes, defined by either vertical or horizontal splits, structure the terminal window. The main pane, which represents the entire terminal window, can be split into additional panes, which in turn can be further subdivided, much like the functionality found in the popular tool, tmux.

Pane sizes can be specified either as a percentage or by a fixed number of characters. The content within a pane can be dynamic, sourced from an attached function (referred to as a Widget in the example code below).

Panes can also feature borders, which are customisable by defining the characters to be used for each side (e.g., left edge, top-left corner, top bar, etc.).

The package utilises ANSI escape codes and has been tested on macOS and Linux.

### Live examples
Check out two command-line games using this library:

* [Ortotris clone](https://github.com/cli-games/ortotris)
* [Snakey Letters](https://github.com/cli-games/snakey-letters)

### Example

See below sample code and a screenshot of its execution.

```
package main

import (
    "os"
    tui "github.com/go-phings/terminal-ui"
)

func main() {
    myTUI := tui.NewTUI()
    myTUI.SetOnDraw(getOnTUIDraw())

    p0 := myTUI.GetPane()

    p01, p02 := p0.SplitVertically(-50, tui.UNIT_PERCENT)
    p021, p022 := p02.SplitVertically(-40, tui.UNIT_CHAR)

    p11, p12 := p01.SplitHorizontally(20, tui.UNIT_CHAR)
    p21, p22 := p021.SplitHorizontally(50, tui.UNIT_PERCENT)
    p31, p32 := p022.SplitHorizontally(-35, tui.UNIT_CHAR)

    s1 := tui.NewTUIPaneStyleFrame()
    s2 := tui.NewTUIPaneStyleMargin()

    // with the package.
    s3 := &tui.TUIPaneStyle{
        NE: "/", NW: "\\", SE: " ", SW: " ", E: " ", W: " ", N: "_", S: " ",
    }

    p11.SetStyle(s1)
    p12.SetStyle(s1)
    p21.SetStyle(s2)
    p22.SetStyle(s2)
    p31.SetStyle(s3)
    p32.SetStyle(s1)

    p11.SetOnDraw(getOnTUIPaneDraw(p11))
    p12.SetOnDraw(getOnTUIPaneDraw(p12))
    p21.SetOnDraw(getOnTUIPaneDraw(p21))
    p22.SetOnDraw(getOnTUIPaneDraw(p22))
    p31.SetOnDraw(getOnTUIPaneDraw(p31))
    p32.SetOnDraw(getOnTUIPaneDraw(p32))

    p11.SetOnIterate(getOnTUIPaneDraw(p11))
    p12.SetOnIterate(getOnTUIPaneDraw(p12))
    p21.SetOnIterate(getOnTUIPaneDraw(p21))
    p22.SetOnIterate(getOnTUIPaneDraw(p22))
    p31.SetOnIterate(getOnTUIPaneDraw(p31))
    p32.SetOnIterate(getOnTUIPaneDraw(p32))

    myTUI.Run(os.Stdout, os.Stderr)
}

func getOnTUIDraw() func(*tui.TUI) int {
    // It does nothing actually.
    fn := func(c *tui.TUI) int {
        return 0
    }
    return fn
}

func getOnTUIPaneDraw(p *tui.TUIPane) func(*tui.TUIPane) int {
    t := tui.NewTUIWidgetSample()
    t.InitPane(p)
    fn := func(x *tui.TUIPane) int {
        return t.Run(x)
    }
    return fn
}

```

![Example](screenshot.png)
