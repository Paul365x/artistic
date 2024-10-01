package gui

import (

	// go imports
	"os"
	"path/filepath"

	// temp imports
	//"fmt"

	// third party imports

	"github.com/prplecake/go-thumbnail"

	// internal
	"github.com/artistic/internal/notify"
	"github.com/artistic/internal/state"
)

// get_thumb, because some of these files are huge we use thumbnails. This function checks for a
// thumbnail and if it is not there creates it and then returns the path
// all internal paths are relative to the artwork root
func get_thumb(path string) string {
	notify.Progress.Start(state.Window)
	dir, file := filepath.Split(path)
	//dir = state.Prefs["root"] + dir
	thb := dir + "thb_" + file

	if !file_exists(thb) {
		var config = thumbnail.Generator{
			DestinationPath:   "",
			DestinationPrefix: "thb_",
			Scaler:            "CatmullRom",
		}

		gen := thumbnail.NewGenerator(config)

		i, err := gen.NewImageFromFile(path)
		if err != nil {
			notify.Notify("Failed to open image", "error", state.Error)
			notify.Progress.Stop()
			return ""
		}

		thumbBytes, err := gen.CreateThumbnail(i)
		if err != nil {
			//notify.Notify("Failed to create thumbnail", "error", state.Error)
			notify.Progress.Stop()
			return ""
		}

		err = os.WriteFile(thb, thumbBytes, 0644)
		if err != nil {
			notify.Notify("Failed to write thumbnail to disk", "error", state.Error)
		}
	}
	notify.Progress.Stop()
	return thb
} // get_thumb()

// file_exists checks whether the given path exists or not
func file_exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
} // file_exists()
