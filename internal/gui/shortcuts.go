// package gui contains all ui elements
// keyboard shortcuts for menu items
package gui

import (
	"github.com/artistic/internal/state"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
)

// list of shortcuts
var (
	short_about = &desktop.CustomShortcut{
		KeyName:  fyne.KeyH,
		Modifier: fyne.KeyModifierControl,
	}

	short_new = &desktop.CustomShortcut{
		KeyName:  fyne.KeyF,
		Modifier: fyne.KeyModifierControl,
	}

	short_open = &desktop.CustomShortcut{
		KeyName:  fyne.KeyO,
		Modifier: fyne.KeyModifierControl,
	}

	short_prefs = &desktop.CustomShortcut{
		KeyName:  fyne.KeyP,
		Modifier: fyne.KeyModifierControl,
	}

	short_save = &desktop.CustomShortcut{
		KeyName:  fyne.KeyS,
		Modifier: fyne.KeyModifierControl,
	}
)

/*
**
** wrappers
**
 */

func wshort_about(shortcut fyne.Shortcut) {
	about_about()
}

func wshort_new(shortcut fyne.Shortcut) {
	file_new()
}

func wshort_open(shortcut fyne.Shortcut) {
	file_open()
}

func wshort_prefs(shortcut fyne.Shortcut) {
	about_prefs()
}

func wshort_save(shortcut fyne.Shortcut) {
	file_save()
}

func Pod_shorts() {
	c := state.Window.Canvas()
	c.AddShortcut(short_about, wshort_about)
	c.AddShortcut(short_new, wshort_new)
	c.AddShortcut(short_open, wshort_open)
	c.AddShortcut(short_prefs, wshort_prefs)
	c.AddShortcut(short_save, wshort_save)
}
