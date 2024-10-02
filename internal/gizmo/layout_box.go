// a fyne layout
package gizmo

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

// layout state
type gizmo_box struct {
	Size float32
}

// NewGizmoBox construstor
func NewGizmoBox(s float32) fyne.Layout {
	return &gizmo_box{Size: s}
}

// MinSize gives a fixed size based on the size of the default textsize and the requested size
func (d *gizmo_box) MinSize(objects []fyne.CanvasObject) fyne.Size {
	if d.Size >= 25.0 {
		panic("Tried to make gizmo_box too large")
	}
	txt_size := fyne.MeasureText("M",
		theme.TextSize(),
		fyne.TextStyle{Bold: false},
	)
	txt_size.Width = 0
	txt_size.Height *= d.Size
	return txt_size
}

// Layout draws the object at the requested size with a 3px pad
func (d *gizmo_box) Layout(objects []fyne.CanvasObject, containerSize fyne.Size) {
	size := d.MinSize(objects)
	objects[0].Resize(size)
}
