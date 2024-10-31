// fyne widget - label, entry and button
package gizmo

/*
**
 */
import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"

	//"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"

	//"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	//"image/color"

	"github.com/artistic/internal/notify"
	"github.com/artistic/internal/state"

	"golang.design/x/clipboard"
)

// Declare conformity with interfaces.
var _ fyne.Widget = (*EnhancedEntry)(nil)

type EnhancedEntry struct {
	widget.BaseWidget
	Label *widget.Label
	Input *widget.Entry
	clip  *widget.Button
	Dirty bool
}

func NewEnhancedEntry(label string, plc string, multi bool, on_chg func(string)) *EnhancedEntry {
	clippy := widget.NewButtonWithIcon("", theme.ContentCopyIcon(), nil)
	entry := &EnhancedEntry{
		Label: widget.NewLabel(label),
		Input: nil,
		clip:  clippy,
	}
	if multi {
		entry.Input = widget.NewMultiLineEntry()
	} else {
		entry.Input = widget.NewEntry()
	}

	entry.Input.OnChanged = on_chg
	entry.Input.SetPlaceHolder(plc)
	entry.ExtendBaseWidget(entry)

	clippy.OnTapped = func() {
		// need to add in notify
		clipboard.Write(clipboard.FmtText, []byte(entry.Input.Text))
		notify.Notify(string("Copied..."), "aok", state.Error)
	}
	return entry
}

func (e *EnhancedEntry) CreateRenderer() fyne.WidgetRenderer {
	sub_c := container.New(
		layout.NewFormLayout(),
		e.Label,
		e.Input,
	)
	c := container.NewVBox(
		container.NewBorder(
			nil,    //top
			nil,    //bottom
			nil,    //left
			e.clip, //right
			sub_c,  //body
		),
	)
	return widget.NewSimpleRenderer(c)
}

/*
func main() {
	// Init returns an error if the package is not ready for use.
	err := clipboard.Init()
	if err != nil {
		panic(err)
	}
	myApp := app.New()
	w := myApp.NewWindow("Lines")

	upd := func(value string) {
		fmt.Println("dirty input: ", value)
	}
	content := NewEnhancedEntry("Label", "plc holder", true, upd)
	w.SetContent(content)

	w.Resize(fyne.NewSize(1000, 1000))
	w.ShowAndRun()
}
*/
