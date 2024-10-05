package terminalui

import (
	"math"
)

const SPLIT_NONE = 0
const SPLIT_H = 1
const SPLIT_V = 2

const UNIT_PERCENT = 1
const UNIT_CHAR = 2

// TUIPane represent a pane within the terminal interface. It has a name.
// It can be split horizontally or vertically to create another 2 panes.
// Split can be described as percentage or fixed characters and only one
// of the panes created from split can have fixed size. Other one is calculated
// from total width.
// Pane also have min width, min height, style and have two events: onDraw
// and onIterate.
type TUIPane struct {
	name       string
	split      int
	splitValue int
	splitUnit  int
	tooSmall   bool
	tui        *TUI
	panes      [2]*TUIPane
	onDraw     func(p *TUIPane) int
	onIterate  func(p *TUIPane) int
	width      int
	height     int
	left       int
	top        int
	minWidth   int
	minHeight  int
	style      *TUIPaneStyle
}

// GetName returns name
func (p *TUIPane) GetName() string {
	return p.name
}

// GetSplit returns split type (horizontal or vertical)
func (p *TUIPane) GetSplit() int {
	return p.split
}

// GetTUI returns TUI instance that this pane is attached to
func (p *TUIPane) GetTUI() *TUI {
	return p.tui
}

// GetPanes returns pane instances created by split
func (p *TUIPane) GetPanes() [2]*TUIPane {
	return p.panes
}

// GetOnDraw returns onDraw event func
func (p *TUIPane) GetOnDraw() func(p *TUIPane) int {
	return p.onDraw
}

// GetOnIterate returns onIterate event func
func (p *TUIPane) GetOnIterate() func(p *TUIPane) int {
	return p.onIterate
}

// GetStyle returns style instance
func (p *TUIPane) GetStyle() *TUIPaneStyle {
	return p.style
}

// SetOnDraw sets onDraw event func
func (p *TUIPane) SetOnDraw(f func(p *TUIPane) int) {
	p.onDraw = f
}

// SetOnIterate sets onIterate event func
func (p *TUIPane) SetOnIterate(f func(p *TUIPane) int) {
	p.onIterate = f
}

// SetStyle sets style
func (p *TUIPane) SetStyle(s *TUIPaneStyle) {
	p.style = s
}

// Split creates new two panes by splitting this pane either
// horizontally or vertically.
// Type, size, size unit are func arguments.
// Function returns pointers to two new panes.
func (p *TUIPane) Split(t int, s int, u int) (*TUIPane, *TUIPane) {
	p.panes[0] = NewTUIPane("Nazwa", p.tui)
	p.panes[1] = NewTUIPane("Nazwa", p.tui)
	p.split = t
	p.splitValue = s
	p.splitUnit = u
	return p.panes[0], p.panes[1]
}

// SplitVertically splits pane vertically. It takes size and size unit
// as arguments. Only one of the two new panes gets the defined size. If the
// value is < 0 then it's the left one, when the value is > 0 then it is
// the right one.
func (p *TUIPane) SplitVertically(s int, u int) (*TUIPane, *TUIPane) {
	return p.Split(SPLIT_V, s, u)
}

// SplitHorizontally splits pane horizontally. It takes size and size unit
// as arguments. Only one of the two new panes gets the defined size. If the
// value is < 0 then it's the top one, when the value is > 0 then it is the
// right one.
func (p *TUIPane) SplitHorizontally(s int, u int) (*TUIPane, *TUIPane) {
	return p.Split(SPLIT_H, s, u)
}

// GetWidth returns pane width
func (p *TUIPane) GetWidth() int {
	return p.width
}

// GetHeight returns pane height
func (p *TUIPane) GetHeight() int {
	return p.height
}

// GetLeft returns x position of pane on terminal window (main pane)
func (p *TUIPane) GetLeft() int {
	return p.left
}

// GetTop returns y position of pane on terminal window (main pane)
func (p *TUIPane) GetTop() int {
	return p.top
}

// GetMinWidth returns minimal width necessary for pane content to work
func (p *TUIPane) GetMinWidth() int {
	return p.minWidth
}

// GetMinHeight returns minimal height necessary for pane content to work
func (p *TUIPane) GetMinHeight() int {
	return p.minHeight
}

// GetTotalMinWidth returns total minimal width necessary for pane to work
// It is GetMinWidth + width necessary for style
func (p *TUIPane) GetTotalMinWidth() int {
	if p.style != nil {
		return p.minWidth + p.style.H()
	}
	return p.minWidth
}

// GetTotalMinHeight returns total minimal height necessary for pane to work
// It is GetMinHeight + height necessary for style
func (p *TUIPane) GetTotalMinHeight() int {
	if p.style != nil {
		return p.minHeight + p.style.V()
	}
	return p.minHeight
}

// SetWidth sets width of pane, checks if it's not too small for the content
// (search for 'minimal width') and calls panes inside to set their width as
// well.
func (p *TUIPane) SetWidth(w int) {
	p.width = w
	if p.GetTotalMinWidth() > 0 && p.width < p.GetTotalMinWidth() {
		p.tooSmall = true
		return
	}
	p.tooSmall = false
	if p.split == SPLIT_H {
		p.panes[0].SetLeft(p.left)
		p.panes[1].SetLeft(p.left)
		p.panes[0].SetWidth(w)
		p.panes[1].SetWidth(w)
	} else if p.split == SPLIT_V {
		v1, v2, tooSmall := p.getSplitValues()
		if tooSmall {
			p.tooSmall = true
			return
		}
		p.tooSmall = false
		p.panes[0].SetLeft(p.left)
		p.panes[1].SetLeft(p.left + v1)
		p.panes[0].SetWidth(v1)
		p.panes[1].SetWidth(v2)
	}
}

// SetHeight sets height of pane, checks if it's not too small for the content
// (search for 'minimal height') and calls panes inside to set their height as
// well.
func (p *TUIPane) SetHeight(h int) {
	p.height = h
	if p.GetTotalMinHeight() > 0 && p.height < p.GetTotalMinHeight() {
		p.tooSmall = true
		return
	}
	if p.split == SPLIT_V {
		p.panes[0].SetTop(p.top)
		p.panes[1].SetTop(p.top)
		p.panes[0].SetHeight(h)
		p.panes[1].SetHeight(h)
	} else if p.split == SPLIT_H {
		v1, v2, tooSmall := p.getSplitValues()
		if tooSmall {
			p.tooSmall = true
			return
		}
		p.tooSmall = false
		p.panes[0].SetTop(p.top)
		p.panes[1].SetTop(p.top + v1)
		p.panes[0].SetHeight(v1)
		p.panes[1].SetHeight(v2)
	}
}

// SetMinWidth sets minimal width for pane content (without style)
func (p *TUIPane) SetMinWidth(w int) {
	p.minWidth = w
}

// SetMinHeight sets minimal height for pane content (without style)
func (p *TUIPane) SetMinHeight(h int) {
	p.minHeight = h
}

// getSplitValues is used by Split functions to calculate the width
// and height of panes. It takes the split type, split value (and its unit)
// and calculates the size in number of characters. It also checks if the size
// is not too small as well.
func (p *TUIPane) getSplitValues() (int, int, bool) {
	var baseVal int
	var calcVal int

	if p.split == SPLIT_V {
		baseVal = p.width
	} else if p.split == SPLIT_H {
		baseVal = p.height
	} else {
		return 0, 0, false
	}

	if p.splitUnit == UNIT_PERCENT {
		calcVal = int(math.Abs(float64(p.splitValue) / 100 * float64(baseVal)))
	} else {
		calcVal = int(math.Abs(float64(p.splitValue)))
	}
	if calcVal >= baseVal || calcVal < 1 {
		return 0, 0, true
	}

	if p.splitValue < 0 {
		return calcVal, baseVal - calcVal, false
	} else if p.splitValue > 0 {
		return baseVal - calcVal, calcVal, false
	}
	return 0, 0, false
}

// SetLeft sets the left value (x position on main pane)
func (p *TUIPane) SetLeft(l int) {
	p.left = l
}

// SetTop sets the top value (y position on main pane)
func (p *TUIPane) SetTop(t int) {
	p.top = t
}

// Write prints string on the pane
func (p *TUIPane) Write(x int, y int, s string, overwriteStyleFrame bool) {
	if p.split == SPLIT_NONE || p.tooSmall {
		if p.style != nil && !overwriteStyleFrame {
			p.tui.Write(p.left+x+p.style.L(), p.top+y+p.style.T(), s)
		} else {
			p.tui.Write(p.left+x, p.top+y, s)
		}
	}
}

// Draw prints the pane on terminal window
func (p *TUIPane) Draw() int {
	if p.tooSmall {
		if p.width > 0 && p.height > 0 {
			p.Write(0, 0, "!", false)
		}
		return 1
	}
	if p.split != SPLIT_NONE {
		p.panes[0].Draw()
		p.panes[1].Draw()
		return 1
	} else {
		if p.style != nil {
			p.style.Draw(p)
		}
		if p.onDraw != nil {
			return p.onDraw(p)
		}
		return 1
	}
	return 1
}

// Iterate is executed by TUI with every main loop iteration
func (p *TUIPane) Iterate() int {
	if p.tooSmall {
		if p.width > 0 && p.height > 0 {
			p.Write(0, 0, "!", false)
		}
		return 1
	}
	if p.split != SPLIT_NONE {
		p.panes[0].Iterate()
		p.panes[1].Iterate()
		return 1
	} else {
		if p.onIterate != nil {
			return p.onIterate(p)
		}
		return 1
	}
	return 1
}

// NewTUIPane returns new instance of TUIPane
func NewTUIPane(n string, t *TUI) *TUIPane {
	p := &TUIPane{name: n, split: SPLIT_NONE, tui: t}
	return p
}
