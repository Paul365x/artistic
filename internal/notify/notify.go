// package notify contains elements that tell the user what is happening
// mechanism to notify user of success or failure of an action
package notify

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"

	"image/color"
)

// status level colors
var (
	// background colours
	error   = color.RGBA{0xc6, 0x2b, 0x29, 255} // red
	warning = color.RGBA{0xe8, 0x81, 0x14, 255} // orange
	notify  = color.RGBA{0xd8, 0xe8, 0x14, 255} // yellow
	aok     = color.RGBA{0x0f, 0x8c, 0x37, 255} // green

	// text colors
	t_error   = color.RGBA{0xff, 0xff, 0xff, 255} // red
	t_warning = color.RGBA{0xff, 0xff, 0xff, 255} // orange e88114
	t_notify  = color.RGBA{0x00, 0x00, 0x00, 255} // yellow
	t_aok     = color.RGBA{0xff, 0xff, 0xff, 255} // green

)

// Notify takes a message and error level and returns a stack container showing text over
// color with the color indicating the error level.
func assemble_notify(msg string, error_level string) (*canvas.Rectangle, *canvas.Text) {
	var bg_color color.Color
	var fg_color color.Color
	switch error_level {
	case "aok":
		msg = "   OK: " + msg
		bg_color = aok
		fg_color = t_aok
	case "notify":
		msg = "   Notify: " + msg
		bg_color = notify
		fg_color = t_notify
	case "warning":
		msg = "   Warning: " + msg
		bg_color = warning
		fg_color = t_warning
	case "error":
		msg = "   Error: " + msg
		bg_color = error
		fg_color = t_error
	default:
		panic("unknown error level")
	}
	lbl := canvas.NewText(msg, fg_color)
	lbl.TextStyle.Bold = true
	bg := canvas.NewRectangle(bg_color)
	return bg, lbl
}

// NewNotify constructor
func NewNotify(s string, level string) *fyne.Container {
	bg, lbl := assemble_notify(s, level)
	return container.NewStack(bg, lbl)
}

// Notify puts message in notify container
func Notify(s string, level string, c *fyne.Container) *fyne.Container {
	bg, lbl := assemble_notify(s, level)
	c.Objects[0] = bg
	c.Objects[1] = lbl
	c.Refresh()
	return c
}
