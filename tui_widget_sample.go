package terminalui

import (
	"time"
)

type TUIWidgetSample struct {
}

// InitPane sets pane minimal width and height that's necessary for the pane
// to work.
func (w *TUIWidgetSample) InitPane(p *TUIPane) {
	p.SetMinWidth(5)
	p.SetMinHeight(3)
}

// Run is main function which just prints out the current time.
func (w *TUIWidgetSample) Run(p *TUIPane) int {
	t := time.Now()
	p.Write(0, 0, t.Format("15:04:05"), false)
	return 1
}

// NewTUIWidgetSample returns instance of TUIWidgetSample struct
func NewTUIWidgetSample() *TUIWidgetSample {
	w := &TUIWidgetSample{}
	return w
}
