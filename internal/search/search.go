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
	"github.com/blevesearch/bleve/v2/mapping"
)

func buildMappings() mapping.IndexMapping {

	mapping := bleve.NewIndexMapping()
	mapping.TypeField = "Personality"

	// mappings for POD personality
	//	podMapping := bleve.NewDocumentMapping()
	//	mapping.AddDocumentMapping("POD", podMapping)

	// link in the subsection mappings for pod
	metaMapping := bleve.NewDocumentMapping()
	mapping.AddDocumentMapping("Metadata", metaMapping)
	aboutMapping := bleve.NewDocumentMapping()
	searchMapping := bleve.NewDocumentMapping()
	metaMapping.AddSubDocumentMapping("About", aboutMapping)
	metaMapping.AddSubDocumentMapping("Search_data", searchMapping)

	// title - actually indexed
	titleFieldMapping := bleve.NewTextFieldMapping()
	titleFieldMapping.Analyzer = "en"
	aboutMapping.AddFieldMappingsAt("Title", titleFieldMapping)

	// Description - actually indexed
	descriptionFieldMapping := bleve.NewTextFieldMapping()
	descriptionFieldMapping.Analyzer = "en"
	aboutMapping.AddFieldMappingsAt("Description", descriptionFieldMapping)

	// Maintag - actually indexed
	mainFieldMapping := bleve.NewTextFieldMapping()
	mainFieldMapping.Analyzer = "en"
	searchMapping.AddFieldMappingsAt("Maintag", mainFieldMapping)

	// Tags - actually indexed
	tagFieldMapping := bleve.NewTextFieldMapping()
	tagFieldMapping.Analyzer = "en"
	searchMapping.AddFieldMappingsAt("Tags", tagFieldMapping)

	return mapping
}
func Build_dex(path string, p string) bleve.Index {

	//	var str state.StringsFunc
	var bleveIdx bleve.Index
	var err error

	mapping := buildMappings()

	idx_path := path + "/index.bleve"
	err = os.RemoveAll(idx_path)
	if err != nil {
		panic(err)
	}

	bleveIdx, err = bleve.New(idx_path, mapping)
	if err != nil {
		panic(err)
	}

	//	str = state.Get_search_strings(personality)
	err = filepath.Walk(path,
		func(path string, info os.FileInfo, ret error) error {
			if ret != nil {
				return ret
			}
			ext := filepath.Ext(path)
			if ext == ".json" {
				pod := state.Empty_pod()
				pod.Unserialise(path)
				bleveIdx.Index(path, pod)
			}
			return nil
		},
	)
	if err != nil {
		panic(err)
	}
	return bleveIdx
}

/*
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
*/
