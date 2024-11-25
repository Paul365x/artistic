/*
** package gui contains all ui elements
** Wraps are segments of ui that can be (hopefully) reused across the different personalities.
 */
package gui

import (
	"image/color"
	"os"
	"path/filepath"
	"strings"

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
	title := gizmo.NewEnhancedEntry("Title:\u2007\u2007\u2007\u2007\u2007\u2007\u2007",
		"Enter Title...",
		false,
		title_chg)
	title.Input.Text = about.Title

	desc_chg := func(value string) {
		p := state.Data.(*state.Pod_type)
		p.Metadata.About.Description = value
		state.Dirty = true
	}
	desc := gizmo.NewEnhancedEntry("Description:\u2007",
		"Enter Description...",
		true,
		desc_chg)
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
	main := gizmo.NewEnhancedEntry("Main\u2007Tag:\u2007\u2007\u2007",
		"Enter tag here...",
		false,
		main_chg)
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
	default_colors *map[string]color.Color,
	instances []state.Instance_type,
) *fyne.Container {

	// setup our colors
	var color_state state.Internal_color
	color_state.Colors = *default_colors
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
	Instance state.Instance_type
	Index    int
}

var Img *fyne.Container
var Instances map[string]Disp_type
var Instance_idx int

// file_radio_callback is the callback for the file radio button selector
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

// wrap_files contains the file selector and other files
func wrap_files(artwork *state.Artwork_type, img *fyne.Container) *fyne.Container {

	// set up the parent file name
	parent_chg := func(value string) {
		artwork.Parent = value
		state.Dirty = true
		notify.Notify(string("Copied parent file"), "aok", state.Error)
	}

	// want the strings the same length so that they line up - 13 char
	parent := gizmo.NewEnhancedEntry("Parent File: ",
		"Enter Parent File...",
		false,
		parent_chg)
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
	Img = img
	Instances = make(map[string]Disp_type)
	var images []string

	// setup the files
	i := 0
	for _, instance := range artwork.Instances {
		file_name := filepath.Base(instance.Image)
		if file_name == "." {
			continue
		}
		images = append(images, file_name)
		Instances[file_name] = Disp_type{
			Instance: instance,
			Index:    i,
		}
		i++
	}

	radio_cont := gizmo.Pick_Radio(images, "Enter Child File...", file_radio_callback)
	file_radio := radio_cont.Objects[0]
	radio_cont.Hide()
	if artwork.Instances[0].Image != "" {

		radio_cont.Show()
		file_radio.Show()
		file_radio.(*widget.RadioGroup).SetSelected(filepath.Base(artwork.Instances[0].Image))
		file_radio.Refresh()
		radio_cont.Refresh()

	}

	file_container := container.NewBorder(
		row,
		nil, nil, nil,
		radio_cont,
	)

	return file_container
}

// need to create our own file filter for file tree since the standard doesn't also do dirs
type FileDirFilter struct {
	Exts []string
}

// NewFileDirFilter constructs the file extension struct
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

func Wrap_nav() *container.AppTabs {
	tree := wrap_file_tree()
	root := state.Prefs["root"].(*preferences.Pref_single).Value
	search := gizmo.NewSearchBox(root)
	search.List.OnSelected = func(id int) {
		load_file(search.Results[id])
	}
	return container.NewAppTabs(
		container.NewTabItemWithIcon("", theme.FolderIcon(), tree),
		container.NewTabItemWithIcon("", theme.SearchIcon(), search),
	)
}

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
	state.Data = &p
	state.CurrentFile = storage.NewFileURI(path)
	var file_name string
	state.CWD, file_name = filepath.Split(path)
	state.CurrentTreeid = "file://" + state.CWD
	tmp := Pod(*state.Data.(*state.Pod_type))
	var content *container.Split
	content = state.Window.Content().(*container.Split)
	content.Trailing = tmp.Content
	notify.Notify(string("Loaded: ")+file_name, "aok", state.Error)
	content.Refresh()
}
