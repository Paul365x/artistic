package state

import (
	"encoding/json"
	"image/color"
	"os"
	"path/filepath"
	"strconv"
	"unicode"

	"github.com/blevesearch/bleve/v2"

	"github.com/artistic/internal/notify"
)

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
	it.BG.Color = *bg
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
	c.Color = color.RGBA{uint8(red), uint8(green), uint8(blue), uint8(alpha)}
	return nil
} // UnmarshalJson

// Data_type implementation for the POD personality
func (p Pod_type) What_am_i() string {
	return p.Personality
}
