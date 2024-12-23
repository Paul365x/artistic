// fyne widget - label, entry and button
package gizmo

/*
**
 */
import (
	"unicode"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"

	//"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"

	//"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	//"image/color"

	"slices"
	"strings"

	"golang.design/x/clipboard"
)

// Declare conformity with interfaces.
var _ fyne.Widget = (*PickBox)(nil)

type PickBox struct {
	widget.BaseWidget
	Label  *widget.Label
	Input  *widget.Entry
	clip   *widget.Button
	add    *widget.Button
	del    *widget.Button
	List   *widget.List
	Data   []string
	sel_id int
}

func NewPickBox(label string, plc string, on_chg func([]string)) *PickBox {

	clippy := widget.NewButtonWithIcon("", theme.ContentCopyIcon(), nil)
	add_b := widget.NewButtonWithIcon("", theme.ContentAddIcon(), nil)
	del_b := widget.NewButtonWithIcon("", theme.ContentRemoveIcon(), nil)

	entry := &PickBox{
		Label:  widget.NewLabel(label),
		Input:  widget.NewEntry(),
		clip:   clippy,
		add:    add_b,
		del:    del_b,
		List:   nil,
		Data:   []string{},
		sel_id: -1,
	}

	entry.List = widget.NewList(
		func() int {
			return len(entry.Data)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("template")
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			o.(*widget.Label).SetText(entry.Data[i])
		},
	)

	entry.List.OnSelected = func(id int) {
		entry.sel_id = id
		entry.Input.Text = entry.Data[id]
		entry.Input.Refresh()
	}

	entry.Input.SetPlaceHolder(plc)
	entry.ExtendBaseWidget(entry)

	clippy.OnTapped = func() {
		// need to add in notify
		str := strings.Join(entry.Data, ",")
		clipboard.Write(clipboard.FmtText, []byte(str))
	}
	// maybe add focus to add_b on change of input to handle typing then enter
	add_b.OnTapped = func() {
		str := strings.FieldsFunc(entry.Input.Text, func(r rune) bool {
			if unicode.IsSpace(r) {
				return true
			}
			if unicode.IsPunct(r) {
				return true
			}
			return false
		})
		entry.Data = append(entry.Data, str...)
		entry.List.Refresh()
		on_chg(entry.Data)
	}

	del_b.OnTapped = func() {
		if entry.sel_id >= 0 {
			entry.Data = slices.Delete(entry.Data, entry.sel_id, entry.sel_id+1)
			entry.sel_id = -1
			entry.Input.SetText("")
			entry.Input.Refresh()
			entry.List.UnselectAll()
			entry.List.Refresh()
			on_chg(entry.Data)
		}
	}
	return entry
}

func (e *PickBox) CreateRenderer() fyne.WidgetRenderer {
	sub_c := container.New(
		layout.NewFormLayout(),
		e.Label,
		e.Input,
	)
	buttons := container.NewHBox(
		e.add,
		e.del,
		e.clip,
	)
	selector := container.NewVBox(
		container.NewBorder(
			nil,     //top
			nil,     //bottom
			nil,     //left
			buttons, //right
			sub_c,   //body
		),
	)

	c := container.NewBorder(
		selector,
		nil, nil, nil,
		e.List,
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
	content := NewPickBox("Label", "plc holder", upd)
	w.SetContent(content)

	w.Resize(fyne.NewSize(1000, 1000))
	w.ShowAndRun()
}
*/
