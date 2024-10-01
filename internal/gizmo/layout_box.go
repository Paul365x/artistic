package gizmo

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

type gizmo_box struct {
	Size float32
}

func NewGizmoBox( s float32 ) fyne.Layout {
	return  &gizmo_box{Size: s}
}

// MinSize gives a fixed size based on the size of the default textsize and the requested size
func (d *gizmo_box) MinSize(objects []fyne.CanvasObject) fyne.Size {
	if (d.Size >= 25.0) {
		panic("Tried to make gizmo_box too large")
	}
	txt_size := fyne.MeasureText( "M",
		theme.TextSize(), 
		fyne.TextStyle{Bold: false},
	)
	txt_size.Width = 0
	txt_size.Height *= d.Size
//txt_size = fyne.NewSize(60.0, 60.0)
	return txt_size
}

// Layout draws the object at the requested size with a 3px pad
func (d *gizmo_box) Layout(objects []fyne.CanvasObject, containerSize fyne.Size) {
 //   if (len(objects) != 1 ) {
        //only support one - initially to stop circles being squashed
 //   	panic("Tried to use gizmo_box layout for multiple objects")
 //   }    	
//	pos := fyne.NewPos(3,3)
	size := d.MinSize(objects)
//	size.Height -= 3
//	size.Width -= 3
   objects[0].Resize(size)
//    objects[0].Move(pos)
}


