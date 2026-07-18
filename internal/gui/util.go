// package gui contains all ui elements
// this file a selection of helper functions
package gui

import (

	// go imports
	"image"
	"os"
	"path/filepath"

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
			notify.Notify("Failed to configure thumb", "error", state.Error)
			notify.Progress.Stop()
			return ""
		}

		// need to get and set our image dimensions since this is not setting them in the lib
		file, _ := os.Open(path)
		defer file.Close()

		img, _, err := image.Decode(file)
		if err != nil {
			notify.Notify("Failed to dimension image", "error", state.Error)
			notify.Progress.Stop()
			return ""
		}
		bounds := img.Bounds()
		i.Current.Width = bounds.Dx()  // Returns r.Max.X - r.Min.X
		i.Current.Height = bounds.Dy() // Returns r.Max.Y - r.Min.Y
		i.Current.X = 0
		i.Current.Y = 0

		thumbBytes, err := gen.CreateThumbnail(i)
		if err != nil {
			notify.Notify("Failed to create thumbnail", "error", state.Error)
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
