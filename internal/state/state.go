package state

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

// the following globals are used internally and not exposed to the user
var Window fyne.Window                                    // main window of running application
var FileMatch []string = []string{".*", ".json", ".JSON"} // file extensions we show in the file dialog
var CurrentFile fyne.URI = nil                            // file we are currently editing
var CWD string                                            // current working directory
var CurrentTreeid string
var Error *fyne.Container // notification container
// const AppId = "au.com.chubbpaul.artistic"                 // our app id
var Data Data_type = nil               // our json data
var Default_color string = "TEEPUBLIC" // default colorset
var Default_personality string = "POD" // default personality
var Default_size string = "100"        // default screen size
var Default_tree string = "12"         // default tree pane size
var Prefs_form *widget.Form            // form for preferences menu item
var Dirty bool = false                 // flag as to whether we have changes
var IndexName string = "/index.bleve"  // search index file name

// Prefs map is exposed to the user via the preferences menu item
var Prefs map[string]interface{}

// Personality types are all the types known so far
var Personality_types = []string{"POD", "EMB"}
