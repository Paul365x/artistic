package gizmo

import (
	"fyne.io/fyne/v2"
	//"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	//"fyne.io/fyne/v2/data/binding"
	//"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"slices"
)

// Declare conformity with interfaces.
var _ fyne.Widget = (*RadioBox)(nil)

type RadioBox struct {
	widget.BaseWidget
	Input *widget.Entry
	List  *widget.RadioGroup
	Items []string
	Add   *widget.Button
	Del   *widget.Button
}

// Pick_box creates a selectable, addable, deletable list of items
func NewRadioBox(placeholder string, callback func(value string)) *RadioBox {
	box := &RadioBox{
		Input: widget.NewEntry(),
		List:  nil,
		Items: []string{},
	}
	box.Input.SetPlaceHolder(placeholder)

	box.List = widget.NewRadioGroup(box.Items, nil)

	box.List.OnChanged = func(s string) {
		box.Input.SetText(s)
		box.Input.Refresh()
		callback(s)
	}

	box.Add = widget.NewButton("Add", func() {
		box.List.Append(box.Input.Text)
		box.Items = append(box.Items, box.Input.Text)
		box.Input.SetText("")
	})

	box.Del = widget.NewButton("Del", func() {
		sel_id := sliceIndex(len(box.Items), func(i int) bool {
			return box.Items[i] == box.Input.Text
		})
		box.Items = slices.Delete(box.Items, sel_id, sel_id+1)
		box.List.Options = slices.Delete(box.List.Options, sel_id, sel_id+1)
		box.Input.SetText("")
		box.Input.Refresh()
		box.List.Selected = ""
		box.List.Refresh()
	})
	return box
}

func (p *RadioBox) CreateRenderer() fyne.WidgetRenderer {
	// setup the layouts
	buttons := container.NewHBox(
		p.Add,
		p.Del,
	)
	select_bar := container.NewBorder(
		nil, nil,
		buttons,
		nil,
		p.Input,
	)

	c := container.NewBorder(
		select_bar,
		nil, nil, nil,
		p.List,
	)
	return widget.NewSimpleRenderer(c)
}

// SliceIndex finds something in a slice
func sliceIndex(limit int, predicate func(i int) bool) int {
	for i := 0; i < limit; i++ {
		if predicate(i) {
			return i
		}
	}
	return -1
}
