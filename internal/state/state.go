package state

import (
	"encoding/json"
	"image/color"
	"os"
	"path/filepath"
	"strconv"
	"unicode"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"

	"github.com/blevesearch/bleve/v2"

	"github.com/artistic/internal/notify"
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

// internal_color_struct contains colors, name and which is selected
type Internal_color struct {
	Colors   map[string]color.Color
	Names    []string
	Selected string
}

// settings contains the personality and an array of colors
type Settings struct {
	Personality string
	Colors      []color.Color
}

// color type_struct links a color to a name
type Color_type struct {
	BG   color.Color
	Name string
}

// instance_type struct links an image path to a background color
type Instance_type struct {
	Image string
	BG    Color_type
}

// artwork_type links a parent file to a series of children
type Artwork_type struct {
	Parent    string
	Instances []Instance_type
}

// search_type contains tags used on the various platforms to search
type Search_type struct {
	Maintag string
	Tags    []string
}

// about_type contains descriptive metadata such as title and description
type About_type struct {
	Title       string
	Description string
}

// meta_type contains about data and search data
type Meta_type struct {
	About       About_type
	Search_data Search_type
}

// pod type is the top level struct for the POD personality
type Pod_type struct {
	Personality string
	Artwork     Artwork_type
	Metadata    Meta_type
}

// Empty_instance creates and returns an empty instance_type
func Empty_instance() Instance_type {
	bg := new(color.Color)
	it := new(Instance_type)
	it.BG.BG = *bg
	it.BG.Name = ""
	it.Image = ""
	return *it
}

// Empty_meta creates and returns an empty emta_type
func Empty_meta() Meta_type {
	mt := new(Meta_type)
	at := new(About_type)
	sd := new(Search_type)
	mt.About = *at
	mt.Search_data = *sd
	return *mt
}

// Empty_pod creates and returns an empty pod_type
func Empty_pod() Pod_type {
	at := new(Artwork_type)
	at.Instances = append(at.Instances, Empty_instance())
	mt := Empty_meta()
	art := new(Pod_type)
	art.Personality = "POD"
	art.Artwork = *at
	art.Metadata = mt
	return *art
}

// Serialise  parses and writes a meta file to disk
func (p *Pod_type) Serialise(file_name string, root string) error {
	str, _ := json.MarshalIndent(p, "", "    ")
	err := os.WriteFile(file_name, str, 0644)
	if err != nil {
		notify.Notify(err.Error()+file_name,
			"error", Error)
	} else {
		// add to index
		idx_path := filepath.Join(root, IndexName)
		idx, err := bleve.Open(idx_path)
		if err != nil {
			notify.Notify(err.Error()+idx_path+": Try reindexing...",
				"error", Error)
		}
		idx.Index(file_name, p)
		idx.Close()
	}

	return err
}

// Unserialise reads and parses a meta file off disk
func (p *Pod_type) Unserialise(file_name string) error {
	data, err := os.ReadFile(file_name)
	if err == nil {
		err = json.Unmarshal(data, p)
	}
	return err
}

// UnmarshalJSON is a helper function to transform JSON color to go image/color
func (c *Color_type) UnmarshalJSON(b []byte) error {
	str := ""
	key := ""
	var red, green, blue, alpha uint64
	for _, ch := range b {
		if unicode.IsSpace(rune(ch)) {
			continue
		}
		switch ch {
		case '{', '"':
		case '}', ',':
			switch key {
			case "Name":
				c.Name = str
			case "R":
				red, _ = strconv.ParseUint(str, 10, 0)
			case "G":
				green, _ = strconv.ParseUint(str, 10, 0)
			case "B":
				blue, _ = strconv.ParseUint(str, 10, 0)
			case "A":
				alpha, _ = strconv.ParseUint(str, 10, 0)
			}
			key = ""
			str = ""
		case ':':
			if str != "BG" {
				key = str
			}
			str = ""
		default:
			str += string([]byte{ch})

		}
	}
	c.BG = color.RGBA{uint8(red), uint8(green), uint8(blue), uint8(alpha)}
	return nil
} // UnmarshalJson

// interface to allow personality identification
type Data_type interface {
	What_am_i() string
}

// Data_type implementation for the POD personality
func (p Pod_type) What_am_i() string {
	return p.Personality
}
