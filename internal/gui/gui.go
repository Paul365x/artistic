//
// package gui has all the ui elements
// this file has the top level ui funcs ~ one per personality
//

package gui

import (
	"github.com/artistic/internal/color_sets"
	//"github.com/artistic/internal/gizmo"
	"github.com/artistic/internal/preferences"
	"github.com/artistic/internal/state"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"

	//	"fyne.io/fyne/v2/theme"

	//"strconv"
	"sync"
)

// used to lock the top level function when drawing/redrawing
var Mu sync.Mutex

type PodRet struct {
	Content *fyne.Container
	View    *fyne.Container
	Rect    *canvas.Rectangle
}

// Pod displays the Pod struct giving an interface to load, edit and save data
func Pod(pod state.Pod_type) *PodRet {

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

	content := container.NewBorder(state.Error, nil, nil, nil,
		container.NewHSplit(left_pane, right_pane),
	)

	return &PodRet{
		Content: content,
		View:    view,
		Rect:    rect,
	}

}
