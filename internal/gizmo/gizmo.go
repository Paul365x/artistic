package gizmo

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"slices"
)

// gizmo_title returns a container surrounding text that has been styled as a heading/title
func Title(s string) *fyne.Container {
	color := theme.ForegroundColor()
	txt := canvas.NewText(s, color)
	//txt.TextSize = theme.TextHeadingSize()
	// could simply bold them instead
	txt.TextStyle = fyne.TextStyle{Bold: true}
	return container.NewBorder(
		txt,
		nil, nil, nil, nil,
	)
} // gizmo_title()

// Pick_box creates a selectable, addable, deletable list of items
func Pick_box(s []string, placeholder string) *fyne.Container {

	sel_id := -1 // index of selected list item - -1 is no selection

	// setup the widgets and their bindings
	pick_shadow := binding.BindStringList(&s)
	selector := widget.NewEntry()
	selector.SetPlaceHolder(placeholder)

	list := widget.NewListWithData(
		pick_shadow,
		func() fyne.CanvasObject {
			return widget.NewLabel("template")
		},
		func(i binding.DataItem, o fyne.CanvasObject) {
			o.(*widget.Label).Bind(i.(binding.String))
		},
	)
	list.OnSelected = func(id int) {
		sel_id = id
		selector.Text = s[id]
		selector.Refresh()
	}

	add_button := widget.NewButton("Add", func() {
		pick_shadow.Append(selector.Text)
	})
	del_button := widget.NewButton("Del", func() {
		if sel_id >= 0 {
			s = slices.Delete(s, sel_id, sel_id+1)
			sel_id = -1
			selector.SetText("")
			selector.Refresh()
			pick_shadow.Reload()
			list.UnselectAll()
			list.Refresh()
		}
	})

	// setup the layouts
	buttons := container.NewHBox(
		add_button,
		del_button,
	)
	select_bar := container.NewBorder(
		nil, nil,
		buttons,
		nil,
		selector,
	)

	return container.NewBorder(
		select_bar,
		nil, nil, nil,
		list,
	)
}

func SliceIndex(limit int, predicate func(i int) bool) int {
	for i := 0; i < limit; i++ {
		if predicate(i) {
			return i
		}
	}
	return -1
}

// Pick_Radio creates a selectable, addable, deletable radio_group of items
func Pick_Radio(s []string, placeholder string, f func(value string)) *fyne.Container {

	// setup the widgets and their bindings
	pick_shadow := binding.BindStringList(&s)
	selector := widget.NewEntry()
	selector.SetPlaceHolder(placeholder)

	radio_group := widget.NewRadioGroup(s, f)
	radio_group.OnChanged = func(s string) {
		selector.SetText(s)
		selector.Refresh()
	}
	// add and delete also need to change the instances slice
	add_button := widget.NewButton("Add", func() {
		radio_group.Append(selector.Text)
		s = append(s, selector.Text)
	})
	del_button := widget.NewButton("Del", func() {
		sel_id := SliceIndex(len(s), func(i int) bool { return s[i] == selector.Text })
		s = slices.Delete(s, sel_id, sel_id+1)
		radio_group.Options = slices.Delete(radio_group.Options, sel_id, sel_id+1)
		selector.SetText("")
		selector.Refresh()
		pick_shadow.Reload()
		radio_group.Selected = ""
		radio_group.Refresh()
	})

	// setup the layouts
	buttons := container.NewHBox(
		add_button,
		del_button,
	)
	select_bar := container.NewBorder(
		nil, nil,
		buttons,
		nil,
		selector,
	)

	return container.NewBorder(
		select_bar,
		nil, nil, nil,
		radio_group,
	)
}

// Labeled_input creates and entry widget with label and data binding
func Labeled_input(shadow_var *string, plc_holder string, display string) (binding.ExternalString, *fyne.Container) {
	shadow := binding.BindString(shadow_var)
	input := widget.NewEntryWithData(shadow)
	input.SetPlaceHolder(plc_holder)
	return shadow, container.NewBorder(
		nil,
		nil,
		widget.NewLabel(display),
		nil,
		input,
	)

} // Labeled_input
