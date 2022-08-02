package tui

import (
	"strings"
)

// TUIPaneStyle defined pane style
type TUIPaneStyle struct {
	NE string
	N  string
	NW string
	W  string
	SW string
	S  string
	SE string
	E  string
}

// H (horizontal) returns minimal width for borders
func (s *TUIPaneStyle) H() int {
	return s.L() + s.R()
}

// V (vertical) returns minimal height for borders
func (s *TUIPaneStyle) V() int {
	return s.T() + s.B()
}

// L (left) returns width of left border
func (s *TUIPaneStyle) L() int {
	if len(s.NW) > 0 || len(s.W) > 0 || len(s.SW) > 0 {
		return 1
	}
	return 0
}

// R (right) returns width of right border
func (s *TUIPaneStyle) R() int {
	if len(s.NE) > 0 || len(s.E) > 0 || len(s.SE) > 0 {
		return 1
	}
	return 0
}

// T (top) returns height of top border
func (s *TUIPaneStyle) T() int {
	if len(s.NE) > 0 || len(s.N) > 0 || len(s.NW) > 0 {
		return 1
	}
	return 0
}

// B (bottom) returns height of bottom border
func (s *TUIPaneStyle) B() int {
	if len(s.SE) > 0 || len(s.S) > 0 || len(s.SW) > 0 {
		return 1
	}
	return 0
}

// Draw prints border around the pane
func (s *TUIPaneStyle) Draw(p *TUIPane) {
	if s.L() > 0 && s.T() > 0 {
		p.Write(0, 0, s.NW, true)
	}
	if s.L() > 0 && s.B() > 0 {
		p.Write(0, p.GetHeight()-1, s.SW, true)
	}
	if s.R() > 0 && s.T() > 0 {
		p.Write(p.GetWidth()-1, 0, s.NE, true)
	}
	if s.R() > 0 && s.B() > 0 {
		p.Write(p.GetWidth()-1, p.GetHeight()-1, s.SE, true)
	}
	if s.T() > 0 || s.B() > 0 {
		st := 0
		en := p.GetWidth() - 1
		if s.L() > 0 {
			st++
		}
		if s.R() > 0 {
			en--
		}
		if s.T() > 0 {
			p.Write(st, 0, strings.Repeat(s.N, en), true)
		}
		if s.B() > 0 {
			p.Write(st, p.GetHeight()-1, strings.Repeat(s.S, en), true)
		}
	}
	if s.L() > 0 || s.B() > 0 {
		st := 0
		en := p.GetHeight() - 1
		if s.T() > 0 {
			st++
		}
		if s.B() > 0 {
			en--
		}
		if s.L() > 0 {
			for i := st; i <= en; i++ {
				p.Write(0, i, s.W, true)
			}
		}
		if s.R() > 0 {
			for i := st; i <= en; i++ {
				p.Write(p.GetWidth()-1, i, s.E, true)
			}
		}
	}
}

// NewTUIPaneStyleFrame returns TUIPaneStyle instance with a nice frame around
func NewTUIPaneStyleFrame() *TUIPaneStyle {
	w := &TUIPaneStyle{
		NE: "┐", NW: "┌", SE: "┘", SW: "└", E: "│", W: "│", N: "─", S: "─",
	}
	return w
}

// NewTUIPaneStyleMargin returns TUIPaneStyle instance without a frame but with
// an internal margin
func NewTUIPaneStyleMargin() *TUIPaneStyle {
	w := &TUIPaneStyle{
		NE: " ", NW: " ", SE: " ", SW: " ", E: " ", W: " ", N: " ", S: " ",
	}
	return w
}

// NewTUIPaneStyleNone returns TUIPaneStyle instance without any frame or margin
func NewTUIPaneStyleNone() *TUIPaneStyle {
	w := &TUIPaneStyle{
		NE: "", NW: "", SE: "", SW: "", E: "", W: "", N: "", S: "",
	}
	return w
}
