// fyne widget - clickable coloured circle
package gizmo

/*
** Splatch is a widget consisting of a clickable coloured circle
 */
import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"image/color"
)

// Declare conformity with interfaces.
var _ fyne.Widget = (*Splatch)(nil)

// id of the splatch
type SplatchItemID = int

// splatch state
type Splatch struct {
	widget.DisableableWidget
	widget.BaseWidget
	Title    string              // title of the color
	Bg       color.Color         // color itself
	Circle   *canvas.Circle      // color of splatch
	Label    *widget.Label       // title of splatch
	OnTapped func()              `json:"-"` // mouseclick call back
	Span     float32             // number of chars width
	Selected bool                // whether the splatch is selected
	Renderer fyne.WidgetRenderer // ref to our renderer - absurd.
}

// NewSplatch returns a splatch widget
func NewSplatch(label string, bg color.Color, size float32, tapped func()) *Splatch {

	item := &Splatch{
		Title:    label,
		Bg:       bg,
		Span:     size,
		OnTapped: tapped,
		Selected: false,
		Renderer: nil,
		Circle:   nil,
		Label:    nil,
	}

	item.ExtendBaseWidget(item)
	return item
}

// CreateRenderer sets up the splatch renderer
func (item *Splatch) CreateRenderer() fyne.WidgetRenderer {
	// create the circle/dot
	bg := item.Bg
	circle := canvas.NewCircle(bg)
	nrgba := color.NRGBAModel.Convert(theme.BackgroundColor()).(color.NRGBA)

	// colors 0..255, choose on basis of midpoint
	lightness := nrgba.R + nrgba.G + nrgba.B
	if lightness <= 128 {
		circle.StrokeColor = color.White
	} else {
		circle.StrokeColor = color.Black
	}
	circle.StrokeWidth = 2
	circle.Position2.X = 25.0
	circle.Position2.Y = 25.0
	//c_cont := container.NewCenter(circle)
	//c_cont.Resize(fyne.NewSize(25.0, 25.0))
	lbl := widget.NewLabel(item.Title)
	lbl.Wrapping = fyne.TextWrapWord
	lbl.Alignment = fyne.TextAlignCenter
	c := container.New(
		layout.NewGridLayout(1), layout.NewSpacer(), circle, lbl,
	)
	item.Circle = circle
	item.Label = lbl
	item.Renderer = widget.NewSimpleRenderer(c)
	return item.Renderer
}

// update called when the splatch changes
func (b *Splatch) Update(name string, color color.Color) {
	b.Title = name
	nbsp := string('\u2060')

	// setup label - we want to cause two lines
	length := b.Span * 2
	b.Title += nbsp

	for len(b.Title) < int(length) {
		b.Title = b.Title + nbsp
	}

	b.Bg = color
	b.Label.SetText(b.Title)
	b.Circle.FillColor = color
	b.Renderer.Refresh()
}

// MinSize returns the size that this widget should not shrink below
func (b *Splatch) MinSize() fyne.Size {
	txt_size := fyne.MeasureText("M",
		theme.TextSize(),
		fyne.TextStyle{Bold: false},
	)
	txt_size.Width *= b.Span
	txt_size.Height += txt_size.Width + (2 * txt_size.Height) + 15
	return txt_size
}

// Tapped is called when a pointer tapped event is captured and triggers any tap handler
func (b *Splatch) Tapped(*fyne.PointEvent) {
	if b.Disabled() {
		return
	}

	if b.OnTapped != nil {
		b.OnTapped()
	}
}
