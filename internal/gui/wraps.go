/*
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
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"

	x_widget "fyne.io/x/fyne/widget"
)

// wrap_about displays an about struct and keeps the backing store in sync
// returns some nested containers
func wrap_about(about *state.About_type) *fyne.Container {

	title_input := widget.NewEntry()
	title_input.Text = about.Title
	title_input.OnChanged = func(value string) {
		p := state.Data.(*state.Pod_type)
		p.Metadata.About.Title = value
		state.Dirty = true
	}
	title_input.SetPlaceHolder("Enter Title...")
	title_label := widget.NewLabel("Title:")

	desc_input := widget.NewMultiLineEntry()
	desc_input.Text = about.Description
	desc_input.Wrapping = fyne.TextWrapWord
	desc_input.OnChanged = func(value string) {
		p := state.Data.(*state.Pod_type)
		p.Metadata.About.Description = value
		state.Dirty = true
	}
	desc_input.SetPlaceHolder("Enter Description...")
	desc_label := widget.NewLabel("Description: ")

	form := container.New(
		layout.NewFormLayout(),
		title_label,
		title_input,
		desc_label,
		desc_input,
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

	// setup maintag box with binding
	main_shadow := binding.BindString(&search.Maintag)
	main_input := widget.NewEntryWithData(main_shadow)
	main_input.SetPlaceHolder("Enter Main Tag...")
	main_input.OnChanged = func(v string) {
		state.Dirty = true
	}

	//main tag
	main_box := container.New(
		layout.NewVBoxLayout(),
		container.NewBorder(
			nil, nil,
			widget.NewLabel("Main Tag:  "),
			nil,
			main_input,
		),
		layout.NewSpacer(),
	)

	tags := gizmo.Pick_box(search.Tags, "Enter tag here...")
	top_bar := container.NewVBox(
		container.NewVBox(
			gizmo.Title("Search"),
			layout.NewSpacer(),
			main_box,
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

type Disp_type struct {
	Instance state.Instance_type
	Index    int
}

var Img *fyne.Container
var Instances map[string]Disp_type
var Instance_idx int

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

func wrap_files(artwork *state.Artwork_type, img *fyne.Container) *fyne.Container {

	// set up the parent file name
	parent_shadow := binding.BindString(&artwork.Parent)
	parent_input := widget.NewEntryWithData(parent_shadow)
	parent_input.SetPlaceHolder("Enter Parent File...")
	parent_input.OnChanged = func(v string) {
		state.Dirty = true
	}

	title := gizmo.Title("Files:")
	row := container.NewBorder(
		title,
		nil,
		widget.NewLabel("Parent File: "),
		nil,
		parent_input,
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

func NewFileDirFilter(ext []string) storage.FileFilter {
	return &FileDirFilter{Exts: ext}
}

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

func open_down_to(u string, t *x_widget.FileTree) {
	path := strings.Replace(u, "file:///", "", 1)
	parts := strings.Split(path, "/")
	u_walk := "file://"
	for _, part := range parts {
		u_walk = u_walk + "/" + part
		t.OpenBranch(u_walk)
	}
}

func wrap_file_tree() *x_widget.FileTree {
	// create the file tree
	rootUri := storage.NewFileURI(state.Prefs["root"].(*preferences.Pref_single).Value)
	tree := x_widget.NewFileTree(rootUri)
	tree.Filter = NewFileDirFilter(state.FileMatch) // Filter files and dirs
	tree.Sorter = func(u1, u2 fyne.URI) bool {
		return u1.String() < u2.String() // Sort alphabetically
	}
	tree.OnSelected = func(u string) {
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
		Pod(*state.Data.(*state.Pod_type))
		notify.Notify(string("Loaded: ")+file_name, "aok", state.Error)
		state.Window.Content().Refresh()
	}
	tree.Show()
	open_down_to(state.CurrentTreeid, tree)
	return tree
}
