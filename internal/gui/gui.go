package gui

import (
	"internal/color_sets"
	"internal/preferences"
	"internal/state"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
)

// Pod displays the Pod struct giving an interface to load, edit and save data
func Pod(pod state.Pod_type) {

	w := state.Window

	tree := wrap_file_tree()

	// create the lefthand interior
	left_pane := container.NewVSplit(
		wrap_about(&pod.Metadata.About),
		wrap_search(&pod.Metadata.Search_data),
	)
	left_pane.SetOffset(0.2)

	// create image and background
	rect := canvas.NewRectangle(pod.Artwork.Instances[0].BG.BG)
	image, img_wrap := wrap_image(rect, &pod.Artwork.Instances[0])

	// setup the color picker
	col := wrap_colors(
		rect,
		"White",
		color_sets.Load_set(preferences.Get_value("color_set"))(),
		pod.Artwork.Instances,
	)

	// setup the parent and child files selector
	files := wrap_files(&pod.Artwork, image)
	files.Refresh()
	view := container.New(
		layout.NewGridLayout(1),
		img_wrap,
		col,
		files,
	)

	right_pane := container.NewBorder(
		nil, nil, nil, nil, view,
	)

	// setup our menus
	var menu fyne.MainMenu
	menu.Items = append(menu.Items,
		menu_file(),
		menu_palette(rect, view, &menu, pod.Artwork.Instances),
		menu_about(),
	)

	content := container.NewBorder(state.Error, nil, nil, nil,

		container.NewHSplit(left_pane, right_pane),
	)

	w_layout := container.NewHSplit(tree, content)
	w_layout.SetOffset(0.12)
	w_layout.Refresh()

	// put everything in place and kick it off
	Pod_shorts()
	w.SetMainMenu(&menu)
	w.SetMaster()
	w.SetContent(w_layout)
}
