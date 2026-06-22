/*
** package gui contains all ui elements
** Wraps are segments of ui that can be (hopefully) reused across the different personalities.
 */
package gui

import (
	//"fmt"
	"image/color"
	"os"
	"path/filepath"
	"strings"

	"golang.design/x/clipboard"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"github.com/artistic/internal/gizmo"
	"github.com/artistic/internal/notify"
	"github.com/artistic/internal/preferences"
	"github.com/artistic/internal/state"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	x_widget "fyne.io/x/fyne/widget"
)

// wrap_about displays an about struct and keeps the backing store in sync
// returns some nested containers
func wrap_about(about *state.About_type) *fyne.Container {

	title_chg := func(value string) {
		p := state.Data.(*state.Pod_type)
		p.Metadata.About.Title = value
		state.Dirty = true
	}

	// want the strings the same length so that they line up - 13 char
	title := gizmo.NewEnhancedEntry(
		"Title:\u2007\u2007\u2007\u2007\u2007\u2007\u2007",
		&state.CWD,
		state.Prefs["root"].(*preferences.Pref_single).Value,
		"Enter Title...",
		false,
		false,
		title_chg,
		notify.Notify,
		state.Window,
		state.Error,
	)
	title.Input.Text = about.Title

	desc_chg := func(value string) {
		p := state.Data.(*state.Pod_type)
		p.Metadata.About.Description = value
		state.Dirty = true
	}

	desc := gizmo.NewEnhancedEntry(
		"Description:\u2007",
		&state.CWD,
		state.Prefs["root"].(*preferences.Pref_single).Value,
		"Enter Description...",
		true,
		false,
		desc_chg,
		notify.Notify,
		state.Window,
		state.Error,
	)
	desc.Input.Text = about.Description
	desc.Input.Wrapping = fyne.TextWrapWord

	form := container.NewVBox(
		title,
		desc,
	)

	// pass back the module
	return container.New(
		layout.NewVBoxLayout(),
		gizmo.Title("About"),
		form,
		layout.NewSpacer(),
	)

} //wrap_about()

// wrap_search displays a search type struct and keeps the backing store in sync
func wrap_search(search *state.Search_type) *fyne.Container {

	// setup main
	main_chg := func(value string) {
		search.Maintag = value
		state.Dirty = true
	}

	main := gizmo.NewEnhancedEntry(
		"Main\u2007Tag:\u2007\u2007\u2007",
		&state.CWD,
		state.Prefs["root"].(*preferences.Pref_single).Value,
		"Enter tag here...",
		false,
		false,
		main_chg,
		notify.Notify,
		state.Window,
		state.Error,
	)
	main.Input.SetText(search.Maintag)

	// setup tags
	tags_chg := func(values []string) {
		search.Tags = values
		state.Dirty = true
	}
	tags := gizmo.NewPickBox("Tags:\u2007\u2007\u2007\u2007\u2007\u2007\u2007",
		"Enter tag here...",
		tags_chg)
	tags.Data = append(tags.Data, search.Tags...)
	// lay it out
	top_bar := container.NewVBox(
		container.NewVBox(
			gizmo.Title("Search"),
			layout.NewSpacer(),
			main,
		),
		tags,
	)

	return container.NewBorder(top_bar, nil, nil, nil, tags)

} // wrap_search()

// wrap_image displays an instance  type struct and keeps the backing store in sync
// it allows the changing of the image background.
func wrap_image(rect *canvas.Rectangle, instance *state.Instance_type) (
	*fyne.Container, *fyne.Container) {

	var img *canvas.Image
	if instance.Image != "" {
		file := instance.Image
		d := filepath.Join(state.Prefs["root"].(*preferences.Pref_single).Value, file)
		thb := get_thumb(d)
		img = canvas.NewImageFromFile(thb)
		img.FillMode = canvas.ImageFillOriginal
	} else {
		img = new(canvas.Image)
	}
	// not sure why this is necessary but get it wrong and the color change in rect doesn't display
	img_border := container.NewBorder(nil, nil, nil, nil, img)
	title := gizmo.Title("Artwork")
	img_stack := container.NewStack(rect, img_border)
	img_main := container.NewBorder(title, nil, nil, nil, img_stack)
	return img_border, img_main
} // wrap_image()

// wrap_colours creates a container holding the colour palette
func wrap_colors(
	rect *canvas.Rectangle,
	selected string,
	Default_colorsets *map[string]color.Color,
	instances []state.Instance_type,
) *fyne.Container {

	// setup our colors
	var color_state state.Internal_color
	color_state.Colors = *Default_colorsets
	color_state.Selected = selected
	for k := range color_state.Colors {
		color_state.Names = append(color_state.Names, k)
	}

	// calculate the size - start with biggest string
	var width = 0
	for name := range color_state.Colors {
		words := strings.Fields(name)
		for _, w := range words {
			if len(w) > width {
				width = len(w)
			}
		}
	}

	grid := widget.NewGridWrap(
		// length callback
		func() int {
			return len(color_state.Colors)
		},
		// CreateItem callback
		func() fyne.CanvasObject {
			name := color_state.Selected
			col := color_state.Colors[name]
			return gizmo.NewSplatch(name, col, float32(width-3), nil)
		},
		// UpdateItem callback
		func(id widget.GridWrapItemID, item fyne.CanvasObject) {
			item.(*gizmo.Splatch).Update(color_state.Names[id],
				color_state.Colors[color_state.Names[id]])
			item.(*gizmo.Splatch).OnTapped = func() {
				rect.FillColor = color_state.Colors[color_state.Names[id]]
				rect.Refresh()
				color_state.Selected = color_state.Names[id]
				state.Dirty = true
				instances[Instance_idx].BG.BG = color_state.Colors[color_state.Names[id]]
				instances[Instance_idx].BG.Name = color_state.Names[id]
			}
		},
	)
	caser := cases.Title(language.English)
	str := string("Palette: ") + caser.String(state.Prefs["color_set"].(*preferences.Pref_multi).Value)
	title := gizmo.Title(str)
	return container.NewBorder(title, nil, nil, nil, grid)
}

// art object instance store
type Disp_type struct {
	Instance *state.Instance_type
	Index    int
}

var Img *fyne.Container              // container for image display
var Instances map[string]Disp_type   // map file base name to a disp_type
var Instance_idx int                 // index into state.Data.Artwork.Instances

// file_radio_callback is the callback for the file radio button selection
func file_radio_callback(value string) {
	file_path := Instances[value].Instance.Image
	if file_path != "" {
		thb := get_thumb(preferences.Get_value("root") + file_path)
		new_img := canvas.NewImageFromFile(thb)
		new_img.FillMode = canvas.ImageFillOriginal
		Img.RemoveAll()
		Img.Add(new_img)
		Img.Refresh()
		Instance_idx = Instances[value].Index
	}
}

// file_radio_add is the callback for the pickRadio add button
// string is the relative path of a file
func file_radio_add(value string ) bool {
	base_name := filepath.Base(value)
	if _, exists := Instances[base_name]; exists {
		notify.Notify(string("This file has already been added"), "error", state.Error)
		return false
	}

	// setup the instance
	instance := state.Empty_instance()
	instance.Image = value
	instance.BG.Name = state.Default_color_name
	instance.BG.BG = state.Default_color

	// add to data and get the index
	insts := state.Data.(*state.Pod_type).Artwork.Instances
	idx := len(insts) 
	insts = append(insts, instance)

	// add to the Instances
	d := Disp_type {
		Instance: &instance,
		Index: idx,
	}
	Instances[base_name] = d

	return true
}

// file_radio_del is the callback for the pickRadio del button
// string is the relative path of a file
func file_radio_del(value string ) bool {
	base_name := filepath.Base(value)
	if _, exists := Instances[base_name]; !exists {
		notify.Notify(string("This file hasn't been added"), "error", state.Error)
		return false
	}

	// remove from the backing store
	insts := state.Data.(*state.Pod_type).Artwork.Instances
	idx := Instances[base_name].Index 
	insts = append(insts[:idx], insts[idx+1:]...)

	// remove from display store
	delete(Instances, base_name)

	return true

}

// wrap_files contains the file selector and other files
func wrap_files(artwork *state.Artwork_type, img *fyne.Container) *fyne.Container {

	// need this in the parent_chg call back
	pr := gizmo.PickRadio {
			Sign :   state.Error,
			Window : state.Window,
			Rg :     nil,
	        S :      []string{},
	        Root :   state.Prefs["root"].(*preferences.Pref_single).Value,
			Cwd :    &state.CWD,     
	        Plc :    "Enter Child File...",
			Notify : notify.Notify,
	        Change : file_radio_callback,
			Add :    file_radio_add,
			Del :    file_radio_del,
	}
	radio_cont := pr.Create()

	// set up the parent file name widget - call back for changes to that field
	parent_chg := func(value string) {
		if len(value) > 0 {
			radio_cont.Show()
			radio_cont.Objects[0].Show()
		} else {
			radio_cont.Hide()
		//	radio_cont.Objects[0].Hide()
		}
		artwork.Parent = value
		state.Dirty = true
		notify.Notify(string("Copied parent file"), "aok", state.Error)
	}

	// want the strings the same length so that they line up - 13 char
	parent := gizmo.NewEnhancedEntry("Parent File: ",
		&state.CWD,
		state.Prefs["root"].(*preferences.Pref_single).Value,
		"Enter Parent File...",
		false,
		true,
		parent_chg,
		notify.Notify,
		state.Window,
		state.Error,
	)
	parent.Input.Text = artwork.Parent

	title := gizmo.Title("Files:")
	row := container.NewBorder(
		title,
		nil,
		nil,
		nil,
		parent,
	)

	// setup globals and locals required for Pick_Radio
	Instances = make(map[string]Disp_type)	

	// setup the files
	i := 0
	for _, instance := range artwork.Instances {
		file_name := filepath.Base(instance.Image)
		if file_name == "." {
			continue
		}
		pr.S = append(pr.S, file_name)
		Instances[file_name] = Disp_type{
			Instance: &instance,
			Index:    i,
		}
		pr.Rg.Append(file_name)
		i++
	}


	file_radio := radio_cont.Objects[0]
	//radio_cont.
	radio_cont.Hide()
	if artwork.Instances[0].Image != "" {		
		radio_cont.Show()
		file_radio.Show()
		pr.Rg.SetSelected(filepath.Base(artwork.Instances[0].Image))
		file_radio.Refresh()
		radio_cont.Refresh()
	}

	file_container := container.NewBorder(
		row,
		nil, nil, nil,
		radio_cont,
	)

	return file_container
} // wrap_files()

// need to create our own file filter for file tree since the standard doesn't also do dirs
type FileDirFilter struct {
	Exts []string
}

// NewFileDirFilter constructs the file extension struct
// used for the tree
func NewFileDirFilter(ext []string) storage.FileFilter {
	return &FileDirFilter{Exts: ext}
}

// Matches indicates whether a path matches the filter
func (fd *FileDirFilter) Matches(uri fyne.URI) bool {
	path := uri.Path()

	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	if info.IsDir() {
		return true
	}
	extension := uri.Extension()
	for _, ext := range fd.Exts {
		if strings.EqualFold(extension, ext) {
			return true
		}
	}
	return false

}

// open_down_to opens all branches in a tree to the given node
func open_down_to(u string, t *x_widget.FileTree) {
	path := strings.Replace(u, "file:///", "", 1)
	parts := strings.Split(path, "/")
	u_walk := "file://"
	for _, part := range parts {
		u_walk = u_walk + "/" + part
		t.OpenBranch(u_walk)
	}
}

// wrap_file_tree contains the file tree
func wrap_file_tree() *x_widget.FileTree {
	// create the file tree
	rootUri := storage.NewFileURI(state.Prefs["root"].(*preferences.Pref_single).Value)
	tree := x_widget.NewFileTree(rootUri)
	tree.Filter = NewFileDirFilter(state.FileMatch) // Filter files and dirs
	tree.Sorter = func(u1, u2 fyne.URI) bool {
		return u1.String() < u2.String() // Sort alphabetically
	}
	tree.OnSelected = load_file
	tree.Show()
	//open_down_to(state.CurrentTreeid, tree)
	return tree
} // wrap_file_tree()

// Wrap_nav creates the lefthand tabbed navigation pane
func Wrap_nav() *container.AppTabs {
	tree := wrap_file_tree()
	root := state.Prefs["root"].(*preferences.Pref_single).Value
	search := gizmo.NewSearchBox(root, true)
	search.List.OnSelected = func(id int) {
		load_file(search.Results[id])
	}
	meta := gizmo.NewSearchBox(root, false)
	meta.List.OnSelected = func(id int) {
		clipboard.Write(clipboard.FmtText, []byte(meta.Results[id][3:]))
		notify.Notify(string("Copied..."), "aok", state.Error)
	}
	return container.NewAppTabs(
		container.NewTabItemWithIcon("", theme.FolderIcon(), tree),
		container.NewTabItemWithIcon("", theme.SearchIcon(), search),
		container.NewTabItemWithIcon("", theme.HistoryIcon(), meta),
	)
} // wrap_nav

// load_file loads the selected file in the callbacks for nav
func load_file(u string) {

	path := strings.Replace(u, "file://", "", 1)
	info, err := os.Stat(path)
	if err != nil {
		return
	}
	if info.IsDir() {
		return
	}

	p := state.Empty_pod()
	p.Unserialise(path)
	var file_name string
	state.CWD, file_name = filepath.Split(path)
	state.CWD = gizmo.AddTrailingSlash(state.CWD)
	state.Data = &p
	state.CurrentFile = storage.NewFileURI(path)
		
	state.CurrentTreeid = "file://" + state.CWD
	tmp := Pod(state.Data.(*state.Pod_type))
	var content *container.Split
	content = state.Window.Content().(*container.Split)
	content.Trailing = tmp.Content
	notify.Notify(string("Loaded: ")+file_name, "aok", state.Error)
	content.Refresh()
} // load_file
