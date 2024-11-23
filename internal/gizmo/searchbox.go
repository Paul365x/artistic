// fyne widget - search results list, entry and button
package gizmo

/*
**
 */
import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
	"github.com/blevesearch/bleve/v2"

	//"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"

	//"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	//"image/color"
)

// Declare conformity with interfaces.
var _ fyne.Widget = (*SearchState)(nil)

type SearchState struct {
	widget.BaseWidget
	Label         *widget.Label
	Input         *widget.Entry
	Search        *widget.Button
	List          *widget.List
	Results       []string
	idx           bleve.Index
	our_container *fyne.Container
}

func NewSearchBox(root string) *SearchState {
	search := widget.NewButtonWithIcon("", theme.SearchIcon(), nil)

	idx_path := root + "/index.bleve"
	idx_db, err := bleve.Open(idx_path)
	if err != nil {
		panic(err)
	}

	data := &SearchState{
		Label:   widget.NewLabel("Search:"),
		Input:   widget.NewEntry(),
		Search:  search,
		Results: []string{},
		List:    nil,
		idx:     idx_db,
	}

	data.Input.SetPlaceHolder("Enter search term...")
	search.OnTapped = func() {
		item := data.Input.Text
		query := bleve.NewMatchQuery(item)
		searchRequest := bleve.NewSearchRequest(query)
		searchResult, _ := data.idx.Search(searchRequest)
		data.Results = []string{}
		for _, result := range searchResult.Hits {
			data.Results = append(data.Results, result.ID)
		}
		data.List = widget.NewList(
			func() int {
				return len(data.Results)
			},
			func() fyne.CanvasObject {
				return widget.NewLabel("template")
			},
			func(i widget.ListItemID, o fyne.CanvasObject) {
				o.(*widget.Label).SetText(data.Results[i])
			})
		data.our_container.Refresh()
	}

	data.List = widget.NewList(
		func() int {
			return len(data.Results)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("template")
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			o.(*widget.Label).SetText(data.Results[i])
		})
	data.List.OnSelected = func(id widget.ListItemID) {
		fmt.Println(id, data.Results[id])
	}

	return data
}

func (s *SearchState) CreateRenderer() fyne.WidgetRenderer {
	selector := container.NewVBox(
		container.NewBorder(
			nil,      //top
			nil,      //bottom
			s.Label,  //left
			s.Search, //right
			s.Input,  //body
		),
	)
	c := container.NewBorder(
		selector,
		nil,
		nil,
		nil,
		s.List,
	)
	s.our_container = c
	return widget.NewSimpleRenderer(c)
}

/*
func main() {
	// Init returns an error if the package is not ready for use.
	err := clipboard.Init()
	if err != nil {
		panic(err)
	}
	myApp := app.New()
	w := myApp.NewWindow("Lines")

	upd := func(value string) {
		fmt.Println("dirty input: ", value)
	}
	content := NewEnhancedEntry("Label", "plc holder", true, upd)
	w.SetContent(content)

	w.Resize(fyne.NewSize(1000, 1000))
	w.ShowAndRun()
}
*/
