// fyne widget - label, entry and button
package search

/*
**
 */
import (
	"os"
	"path/filepath"

	"github.com/artistic/internal/state"
	"github.com/blevesearch/bleve/v2"
)

// search
func Build_dex(path string, personality string) bleve.Index {
	var str state.StringsFunc
	var bleveIdx bleve.Index
	var err error

	idx_path := path + "/index.bleve"
	err = os.RemoveAll(idx_path)
	if err != nil {
		panic(err)
	}
	mapping := bleve.NewIndexMapping()
	bleveIdx, err = bleve.New(idx_path, mapping)
	if err != nil {
		panic(err)
	}

	str = state.Get_search_strings(personality)
	err = filepath.Walk(path,
		func(path string, info os.FileInfo, ret error) error {
			if ret != nil {
				return ret
			}
			ext := filepath.Ext(path)
			if ext == ".json" {
				data := str(path)
				bleveIdx.Index(path, data)
			}
			return nil
		},
	)
	if err != nil {
		panic(err)
	}
	return bleveIdx
}
