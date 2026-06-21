// fyne widget - label, entry and button
package gizmo

/*
**
 */
import (
	"errors"
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/theme"

	//"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"

	//"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	//"image/color"

	"golang.design/x/clipboard"
)

// Declare conformity with interfaces.
var _ fyne.Widget = (*EnhancedEntry)(nil)

type EnhancedEntry struct {
	widget.BaseWidget
	Label    *widget.Label
	Input    *widget.Entry
	clip     *widget.Button
	selector * widget.Button
	file_pick bool
	Dirty    bool
	sign *fyne.Container
	window fyne.Window
	root string
	cwd *string
	notify func(string, string, *fyne.Container) *fyne.Container
	
}

func NewEnhancedEntry(label string,
					  wd *string,
					  r string, 
					  plc string, 
					  multi bool, 
					  file_pick bool,
					  on_chg func(string),
					  n func(string, string, *fyne.Container) *fyne.Container,
					  w fyne.Window,
					  error *fyne.Container) *EnhancedEntry {

	clippy := widget.NewButtonWithIcon("", theme.ContentCopyIcon(), nil)
	selector := widget.NewButtonWithIcon("", theme.FileIcon(), nil)

	entry := &EnhancedEntry{
		Label:    widget.NewLabel(label),
		Input:    nil,
		clip:     clippy,
		selector: selector,
		file_pick: file_pick,
		sign:     error,
		window:   w,
		root:     r,
		cwd:      wd,
		notify:   n,
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
		n(string("Copied..."), "aok", error)
	}
	selector.OnTapped = entry.fileHandler
	return entry
}

func (e *EnhancedEntry) CreateRenderer() fyne.WidgetRenderer {
	sub_c := container.New(
		layout.NewFormLayout(),
		e.Label,
		e.Input,
	)
	btn_c := container.NewHBox(
		e.clip,
	)
	if e.file_pick {
		btn_c.Add(e.selector)
		btn_c.Refresh()
	}

	c := container.NewVBox(
		container.NewBorder(
			nil,    //top
			nil,    //bottom
			nil,    //left
			btn_c, //right
			sub_c,  //body
		),
	)
	return widget.NewSimpleRenderer(c)
}

func (e *EnhancedEntry) fileHandler() {
	d := dialog.NewFileOpen(func(uc fyne.URIReadCloser, err error) {
		if uc != nil {
			path := uc.URI().Path()
			if *e.cwd != e.root && filepath.Dir(path) != *e.cwd {
					err = errors.ErrUnsupported
					e.notify(string("Different directory to the other files"),
					"error",
					e.sign)
					return
			}
			e.Input.Text,_ = filepath.Rel(e.root, path)
			e.Input.OnChanged(uc.URI().Path())
			e.Input.Refresh()
			*e.cwd = filepath.Dir(path)
		}
	},	e.window)
	
	LocationURI, err := storage.ListerForURI(storage.NewFileURI(*e.cwd))
	if err != nil {
		switch err.Error() {
		case "uri is not listable":
			e.notify("Failed to open folder", "error", e.sign)
		default:
			e.notify(err.Error(), "error", e.sign)
		}
		return
	}
	d.SetLocation(LocationURI)
	d.Show()
}

