// fyne widget - search results list, entry and button
package gizmo

/*
**
 */
import (
	//"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/blevesearch/bleve/v2"

	"github.com/artistic/internal/state"
)

// Declare conformity with interfaces.
var _ fyne.Widget = (*SearchState)(nil)

// our widget state
type SearchState struct {
	widget.BaseWidget
	Label         *widget.Label
	Input         *widget.Entry
	Search        *widget.Button
	List          *widget.List
	Results       []string
	Our_container *fyne.Container
	boxtype       bool
	idx_path      string
}

// NewSearchBox creates a search box
// root is the media tree root used to find the index
// filesY switches between IDs ie the filename and Values ie the metadata
func NewSearchBox(root string, filesY bool) *SearchState {

	search := widget.NewButtonWithIcon("", theme.SearchIcon(), nil)

	data := &SearchState{
		Label:         widget.NewLabel("Search:"),
		Input:         widget.NewEntry(),
		Search:        search,
		Results:       []string{},
		List:          nil,
		Our_container: nil,
		boxtype:       filesY,
		idx_path:      root + state.IndexName,
	}

	data.Input.SetPlaceHolder("Enter search term...")

	search.OnTapped = data.SearchTap

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

	return data
} // NewSearchBox

// CreateRenderer returns the widget renderer
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
	s.Our_container = c
	return widget.NewSimpleRenderer(c)
}

// MinSize hooks in the minsize processing otherwise the parent container does strange things
func (s *SearchState) MinSize() fyne.Size {
	s.ExtendBaseWidget(s)
	return s.BaseWidget.MinSize()
}

// SearchTap is the callback for the search button: lookup and populates the results
func (s *SearchState) SearchTap() {

	idx, err := bleve.Open(s.idx_path)
	if err != nil {
		panic(err)
	}

	item := s.Input.Text
	query := bleve.NewMatchQuery(item)
	searchRequest := bleve.NewSearchRequest(query)
	searchRequest.Fields = []string{"*"}
	searchResult, _ := idx.Search(searchRequest)
	s.Results = []string{}
	idx.Close()

	for _, result := range searchResult.Hits {
		if s.boxtype { // return the file paths ie id
			s.Results = append(s.Results, result.ID)

		} else {

			//for _, value := range result.Fields {
			title := "N: " + result.Fields["Metadata.About.Title"].(string)
			description := "D: " + result.Fields["Metadata.About.Title"].(string)
			maintag := "M: " + result.Fields["Metadata.Search_data.Maintag"].(string)
			tags := result.Fields["Metadata.Search_data.Tags"].(interface{})
			s.Results = append(s.Results, []string{
				title,
				description,
				maintag,
			}...)
			value := "T: "
			for _, tag := range tags.([]interface{}) {
				value = value + " " + tag.(string)

			}
			s.Results = append(s.Results, value)
		}
	}

	s.List.Refresh()
	s.Our_container.Refresh()
} // SearchTap
