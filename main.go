package main

import (
	"github.com/artistic/internal/gui"
	"github.com/artistic/internal/notify"
	"github.com/artistic/internal/preferences"
	"github.com/artistic/internal/state"

	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"

	"github.com/kbinani/screenshot"
)

func main() {

	a := app.New()

	// this just shouldn't be necessary
	exe, _ := os.Executable()
	fp := filepath.Dir(exe)
	file := filepath.Join(fp, "Icon.png")
	res, _ := fyne.LoadResourceFromPath(file)
	app.SetMetadata(fyne.AppMetadata{
		Icon:    res,
		Name:    "Artistic",
		ID:      "au.com.chubbpaul.artistic",
		Version: "0.1.0",
		Build:   5,
	})

	state.Window = a.NewWindow("Artistic")

	preferences.Init_prefs()

	// I know - heresy according to fyne devs but this is a graphics program and by default we want
	// realestate. The user can both configure the default size and also resize.
	bounds := screenshot.GetDisplayBounds(0)
	factor64, err := strconv.ParseFloat(state.Prefs["scr_size"].(*preferences.Pref_single).Value, 32)
	if err != nil {
		factor64 = 100.0 // if the prefs corrupt, default to 100%
	}
	factor := float32(factor64 / 100.0)
	state.Window.Resize(fyne.NewSize(float32(bounds.Dx())*factor, float32(bounds.Dy())*factor))
	personality := state.Prefs["personality"].(*preferences.Pref_multi)
	switch personality.Value {
	case "POD":
		pod := state.Empty_pod()
		state.Data = &pod
		state.Error = notify.NewNotify("Started with Empty Artwork", "aok")
		gui.Mu.Lock()
		gui.Pod(pod)
		gui.Mu.Unlock()
		state.Window.ShowAndRun()
	default:
		fmt.Println("Unknown personality: What do I do now?")
	}

}
