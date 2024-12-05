// bleve search routines apart from those embedded in gizmo::searchbox.go
package search

import (
	"os"
	"path/filepath"

	"github.com/artistic/internal/state"

	"github.com/blevesearch/bleve/v2"
	"github.com/blevesearch/bleve/v2/mapping"
)

// buildMappings builds a mapping path for the index to find the text needed for search
func buildMappings() mapping.IndexMapping {

	mapping := bleve.NewIndexMapping()
	mapping.TypeField = "Personality"

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
} // buildMapping

// Build_dex deletes the old index, if any, and indexes all files
func Build_dex(path string, p string) {

	//	var str state.StringsFunc
	var bleveIdx bleve.Index
	var err error

	mapping := buildMappings()

	idx_path := path + state.IndexName
	err = os.RemoveAll(idx_path)
	if err != nil {
		panic(err)
	}

	bleveIdx, err = bleve.New(idx_path, mapping)
	if err != nil {
		panic(err)
	}

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
	bleveIdx.Close()
} // Build_dex
