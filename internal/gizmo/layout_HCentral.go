//
// a fyne layout
//

package gizmo

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

// layout state
type gizmo_HCentral struct {
	Size float32
}

// NewGizmoHCentral constructor
func NewGizmoHCentral(s float32) fyne.Layout {
	return &gizmo_HCentral{Size: s}
}

// MinSize gives a fixed size based on the size of the default textsize and the requested size
func (d *gizmo_HCentral) MinSize(objects []fyne.CanvasObject) fyne.Size {
	tst_string := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	if d.Size >= 25.0 {
		panic("Tried to make gizmo_box too large")
	}
	txt_size := fyne.MeasureText(tst_string[0:uint(d.Size-1.0)],
		theme.TextSize(),
		fyne.TextStyle{Bold: false},
	)
	return txt_size
}

// Layout draws the object centered on the width and at the top
func (d *gizmo_HCentral) Layout(objects []fyne.CanvasObject, containerSize fyne.Size) {
	if len(objects) != 1 {
		//only support one - initially to stop circles being squashed
		panic("Tried to use gizmo_box layout for multiple objects")
	}
	size := d.MinSize(objects)
	o := objects[0]
	childmin := o.MinSize()
	o.Resize(childmin)
	o.Move(fyne.NewPos(float32(size.Width-childmin.Width)/2, 0.0))
}
