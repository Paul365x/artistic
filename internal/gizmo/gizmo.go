//
// the name widget was taken and means something specific in fyne
// these are little reusable artifacts for the gui
//

package gizmo

import (
	"errors"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"path/filepath"
	"slices"
	"strings"

	"golang.design/x/clipboard"
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

// SliceIndex finds something in a slice
func SliceIndex(limit int, predicate func(i int) bool) int {
	for i := 0; i < limit; i++ {
		if predicate(i) {
			return i
		}
	}
	return -1
}

// Pick_Radio creates a selectable, addable, deletable radio_group of items
type PickRadio struct {
	Sign *fyne.Container
	Window fyne.Window
	Rg *widget.RadioGroup
	S []string
	Root string
	Cwd *string
	Plc string
	Notify func(string, string, *fyne.Container) *fyne.Container
	F func(string)
}

func (p *PickRadio) Create () *fyne.Container {

	// setup the widgets and their bindings
	pick_shadow := binding.BindStringList(&p.S)
	selector := widget.NewEntry()
	selector.SetPlaceHolder(p.Plc)

	p.Rg = widget.NewRadioGroup(p.S, p.F)
	p.Rg.OnChanged = func(s string) {
		selector.SetText(s)
		selector.Refresh()
	}
	// add and delete also need to change the instances slice
	add_button := widget.NewButton("Add", func() {
		p.Rg.Append(selector.Text)
		p.Rg.SetSelected(selector.Text)
		p.S = append(p.S, selector.Text)
	})

	del_button := widget.NewButton("Del", func() {
		sel_id := SliceIndex(len(p.S), func(i int) bool { return p.S[i] == selector.Text })
		p.S = slices.Delete(p.S, sel_id, sel_id+1)
		p.Rg.Options = slices.Delete(p.Rg.Options, sel_id, sel_id+1)
		selector.SetText("")
		selector.Refresh()
		pick_shadow.Reload()
		p.Rg.Selected = ""
		p.Rg.Refresh()
	})

	copy_button := widget.NewButtonWithIcon("", theme.ContentCopyIcon(), func() {
		str := strings.Join(p.S, ",")
		clipboard.Write(clipboard.FmtText, []byte(str))
		p.Notify(string("Copied files"), "aok", p.Sign)
	})

	file_button := widget.NewButtonWithIcon("", theme.FileIcon(), func() {
		d := dialog.NewFileOpen(func(uc fyne.URIReadCloser, err error) {
		    if uc != nil {
		    	path := uc.URI().Path()
				if *p.Cwd != p.Root && filepath.Dir(path) != *p.Cwd {
					err = errors.ErrUnsupported
					p.Notify(string("Different directory to the other files"),
					"error",
					p.Sign)
					return
				}
				file,_ := filepath.Rel(p.Root, path)
		    	selector.SetText(file)	
				selector.Refresh()	
		    }
	    },	p.Window)
	    pathURI := storage.NewFileURI(*p.Cwd)
	    listURI, _ := storage.ListerForURI(pathURI)
	    d.SetLocation(listURI)
		d.Show()
	})

	// setup the layouts
	buttons := container.NewHBox(
		add_button,
		del_button,
		copy_button,
		file_button,
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
		p.Rg,
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
