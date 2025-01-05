package main

import (
	"os"

	tui "github.com/go-phings/terminal-ui"
)

func getOnTUIDraw() func(*tui.TUI) int {
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

func main() {
	myTUI := tui.NewTUI()
	myTUI.SetOnDraw(getOnTUIDraw())

	p0 := myTUI.GetPane()

	p01, p02 := p0.SplitVertically(-50, tui.UNIT_PERCENT)
	p021, p022 := p02.SplitVertically(-40, tui.UNIT_CHAR)

	p11, p12 := p01.SplitHorizontally(20, tui.UNIT_CHAR)
	p21, p22 := p021.SplitHorizontally(50, tui.UNIT_PERCENT)
	p31, p32 := p022.SplitHorizontally(-35, tui.UNIT_CHAR)

	// Styles
	s1 := tui.NewTUIPaneStyleFrame()
	s2 := tui.NewTUIPaneStyleMargin()
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
