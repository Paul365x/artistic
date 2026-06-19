// package gui contains all ui elements
// this file all the menu code
package gui

import (
	"log"
	"path/filepath"
	"strings"

	"github.com/artistic/internal/color_sets"
	"github.com/artistic/internal/notify"
	"github.com/artistic/internal/preferences"
	"github.com/artistic/internal/search"
	"github.com/artistic/internal/state"

	"cmp"
	"net/url"
	"slices"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"

	"fyne.io/fyne/v2/widget"

	x_dialog "fyne.io/x/fyne/dialog"
)

// menu_palette creates a menu allowing selection of a color_set for a wrap_colors wrap
func Menu_palette(
	rect *canvas.Rectangle,
	view *fyne.Container,
	main *fyne.MainMenu,
	instances []state.Instance_type) *fyne.Menu {

	var palette fyne.Menu

	palette.Label = "Palette"
	for k, v := range color_sets.Color_sets {
		colors := v()
		menu_item := fyne.NewMenuItem(k, nil)
		palette.Items = append(palette.Items, menu_item)
		menu_item.Action = func() {
			view.Objects[1] = wrap_colors(rect, "White", colors, instances)
			view.Refresh()
			for _, item := range palette.Items {
				item.Checked = false
			}
			menu_item.Checked = true
			main.Refresh()
		}

		// sort them for display
		ItemCmp := func(a, b *fyne.MenuItem) int {
			return cmp.Compare(a.Label, b.Label)
		}
		slices.SortFunc(palette.Items, ItemCmp)

		if k == state.Prefs["color_set"].(*preferences.Pref_multi).Value {
			menu_item.Checked = true
		}
	}
	return &palette
} // menu_palette()

// menu_file creates a menu with the usual type of stuff for file menus ala microsoft
func Menu_file() *fyne.Menu {
	file := fyne.NewMenu("File",
		fyne.NewMenuItem("New", file_new),
		fyne.NewMenuItem("Open", file_open),
		fyne.NewMenuItem("SaveAs", file_save_as),
		fyne.NewMenuItem("Save", file_save),
	)
	return file
} // menu_file()

// menu_pref returns a menu allowing changing of preferences
func Menu_about() *fyne.Menu {

	var about_menu fyne.Menu
	about_menu.Label = "About"

	menu_item := fyne.NewMenuItem("Preferences", about_prefs)
	about_menu.Items = append(about_menu.Items, menu_item)

	menu_item = fyne.NewMenuItem("Reindex", about_index)
	about_menu.Items = append(about_menu.Items, menu_item)

	menu_item = fyne.NewMenuItem("About", about_about)
	about_menu.Items = append(about_menu.Items, menu_item)

	return &about_menu
} // menu_about

/*
** callbacks for menu items
 */

// file_new is the callback for the file new menu item
func file_new() {
	// so we populate the front end with an empty whatsit
	// we save nil to current file
	pod := state.Empty_pod()
	state.Data = &pod
	state.CurrentFile = nil
	state.CWD = state.Prefs["root"].(*preferences.Pref_single).Value
	state.CurrentTreeid = "file://" + state.CWD
	tmp := Pod(*state.Data.(*state.Pod_type))
	var content *container.Split
	content = state.Window.Content().(*container.Split)
	content.Trailing = tmp.Content
	notify.Notify(string("Loaded: Empty"), "aok", state.Error)
	content.Refresh()
} // file_new()

// file_open is the callback for the file open menu item
func file_open() {
	d := dialog.NewFileOpen(func(uc fyne.URIReadCloser, err error) {
		if uc != nil {
			file := uc.URI()
			p := state.Empty_pod()
			err := p.Unserialise(file.Path())
			if err != nil {
				str := file.Name()
				notify.Notify(err.Error()+str, "error", state.Error)
				return
			}
			state.Data = &p
			state.CurrentFile = file
			state.CWD, _ = filepath.Split(uc.URI().Path())
			state.CurrentTreeid = "file://" + state.CWD
			tmp := Pod(*state.Data.(*state.Pod_type))
			var content *container.Split
			content = state.Window.Content().(*container.Split)
			content.Trailing = tmp.Content
			notify.Notify(string("Loaded: ")+file.Name(), "aok", state.Error)
			//state.Window.Content().Refresh()
			content.Refresh()
		}
	},
		state.Window)

	// calls the function (shown nil) with the uri of the file. need to update state.CWD
	// need to update current file to this one
	// load it and refresh
	LocationURI, err := storage.ListerForURI(storage.NewFileURI(state.CWD))
	if err != nil {
		log.Println(err.Error())
		switch err.Error() {
		case "uri is not listable":
			notify.Notify("Failed to open folder", "error", state.Error)
		default:
			notify.Notify(err.Error(), "error", state.Error)
		}
		return
	}
	d.SetLocation(LocationURI)
	d.SetFilter(storage.NewExtensionFileFilter(state.FileMatch))
	d.Show()
} // file_open()

// file_save is the callback for the file save menu item. Also used by the save-as menu item
func file_save() {
	// check we have data to save
	if state.Data == nil {
		notify.Notify(string("No Data to save"), "warning", state.Error)
		return
	}
	// check we have a file
	if state.CurrentFile == nil {
		notify.Notify(string("No file to save into - use Save As"), "notify", state.Error)
		return
	}
	// serialise
	err := state.Data.(*state.Pod_type).Serialise(state.CurrentFile.Path(),
		state.Prefs["root"].(*preferences.Pref_single).Value)
	if err != nil {
		notify.Notify(err.Error()+state.CurrentFile.Name(), "error", state.Error)
		return
	}
	notify.Notify(string("Saved ")+state.CurrentFile.Name(), "aok", state.Error)
} // file_save()

// file_save_as is the callback for the file save as menu item
func file_save_as() {
	f := dialog.NewFileSave(func(uc fyne.URIWriteCloser, err error) {
		if uc != nil {
			file := uc.URI()
			// add ext if missing
			if file.Extension() == "" {
				uc.Close()
				correctURI, err := AddExtension(file, ".json")
        		if err != nil {
           			notify.Notify(string("failed to add JSON extension ")+state.CurrentFile.Name(), "error", state.Error)
					return       
				}
				file = correctURI
			}
			state.CurrentFile = file
			file_save()
			notify.Notify(string("Saved ")+state.CurrentFile.Name(), "aok", state.Error)
		}

	}, state.Window)
	loc, err := storage.ListerForURI(storage.NewFileURI(state.CWD))
	if err != nil {
		log.Println(err.Error())
		switch err.Error() {
		case "uri is not listable":
			notify.Notify("Failed to open folder", "error", state.Error)
		default:
			notify.Notify(err.Error(), "error", state.Error)
		}
		return
	}
	f.SetLocation(loc)
	f.SetFilter(storage.NewExtensionFileFilter(state.FileMatch))
	f.Show()
	
} // save_as()

// about_prefs is the callback for the preferences menu item.
func about_prefs() {
	a := fyne.CurrentApp()
	w := a.NewWindow("Preferences")
	preferences.Window = w
	w.SetContent(state.Prefs_form)
	w.Show()
} // about_prefs()

func about_index() {
	personality := state.Prefs["personality"].(*preferences.Pref_multi).Value
	path := state.Prefs["root"].(*preferences.Pref_single).Value
	notify.Notify(string("Search Reindex: Please wait... "), "aok", state.Error)
	notify.Progress.Start(state.Window)
	search.Build_dex(path, personality)
	notify.Progress.Stop()
	notify.Notify(string("Search Index Ready"), "aok", state.Error)
}

// about_about is the callback for the about menu item
func about_about() {
	docURL, _ := url.Parse("https://docs.fyne.io")
	iconURL, _ := url.Parse("https://www.flaticon.com/free-icons/artist")
	links := []*widget.Hyperlink{
		widget.NewHyperlink("Docs", docURL),
		widget.NewHyperlink("Icon created by Freepik - Flaticon", iconURL),
	}
	res, _ := fyne.LoadResourceFromPath("./artistic-icon.png")
	//	log.Println(err)
	fyne.CurrentApp().SetIcon(res)
	x_dialog.ShowAboutWindow("Managing artistic works and their metadata.", links, fyne.CurrentApp())
} // about_about()

/*
** Utility/helper
*/
func AddExtension(u fyne.URI, ext string) (fyne.URI, error) {
	// Format extension to ensure it starts with a dot
	if !strings.HasPrefix(ext, ".") {
		ext = "." + ext
	}

	// If the URI already ends with this extension, return it as-is
	if strings.HasSuffix(strings.ToLower(u.String()), strings.ToLower(ext)) {
		return u, nil
	}

	// Append extension to the raw string and re-parse
	newURIStr := u.String() + ext
	return storage.ParseURI(newURIStr)
}
