package state

import (
	"image/color"
)

// internal_color_struct contains colors, name and which is selected
type Internal_color struct {
	Colors   map[string]color.Color
	Names    []string
	Selected string
}

// settings contains the personality and an array of colors
/*
type Settings struct {
	Personality string
	Colors      []color.Color
}
*/

// about_type contains descriptive metadata such as title and description
type About_type struct {
	Title       string
	Description string
	Id          string
}

// color type_struct links a color to a name
type Color_type struct {
	Color color.Color
	Name  string
	Code  string
}

// interface to allow personality identification
type Data_type interface {
	What_am_i() string
}
