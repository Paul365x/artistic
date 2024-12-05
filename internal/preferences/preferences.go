// user configuration of look and feel
package preferences

import (
	"github.com/artistic/internal/color_sets"
	"github.com/artistic/internal/state"

	"os"
	"sort"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

var Window fyne.Window = nil // preferences dialog window

/*
** Form items production functions
 */

type Pref_unit_type interface {
	Init()
	Flavor() string
}

type Pref_element struct {
	Populate func() []string
	Label    string
	Hint     string
	Item     *widget.FormItem
	Value    string
}

type Pref_multi struct {
	Values []string
	Pref_element
}

type Pref_single struct {
	Pref_element
}

// Init_prefs sets up the prefs map and other globals
func Init_prefs() {

	state.Prefs = make(map[string]interface{})
	color_sets.Build_sets()

	m := &Pref_multi{}
	m.Populate = Populate_personality
	m.Init()
	state.Prefs["personality"] = m

	m = &Pref_multi{}
	m.Populate = Populate_color
	m.Init()
	state.Prefs["color_set"] = m

	s := &Pref_single{}
	s.Populate = Populate_root
	s.Init()
	state.Prefs["root"] = s
	state.CWD = s.Value
	state.CurrentTreeid = "file://" + state.CWD

	s = &Pref_single{}
	s.Populate = Populate_scr
	s.Init()
	state.Prefs["scr_size"] = s

	s = &Pref_single{}
	s.Populate = Populate_tree
	s.Init()
	state.Prefs["nav_size"] = s

	var items []*widget.FormItem

	keys := make([]string, len(state.Prefs))
	i := 0
	for k := range state.Prefs {
		keys[i] = k
		i++
	}
	sort.Strings(keys)

	for i := range keys {
		pref := keys[i]
		flavor := state.Prefs[pref].(Pref_unit_type).Flavor()
		switch flavor {
		case "multi":
			items = append(items, state.Prefs[pref].(*Pref_multi).Item)
		case "single":
			items = append(items, state.Prefs[pref].(*Pref_single).Item)
		}
	}

	help := widget.NewLabel("Restart is required to use any changes.")
	help.TextStyle = fyne.TextStyle{Italic: true}
	save := widget.NewButton("Save", save_button)
	cancel := widget.NewButton("Cancel", cancel_button)
	reset := widget.NewButton("Reset", reset_button)
	factory := widget.NewButton("Factory", factory_button)
	buttons := widget.FormItem{
		HintText: "Restart is required for any changes to take effect.",
		Widget:   container.NewHBox(save, cancel, reset, factory),
	}

	items = append(items, &buttons)

	state.Prefs_form = &widget.Form{
		Items: items,
	}

}

// This Init creates a single entry element FormItem
func (p *Pref_single) Init() {

	var key string
	var val string
	s := p.Populate()
	if len(s) < 4 {
		panic("Preferences: Populate function returns too few")
	}
	p.Label = s[0]
	p.Hint = s[1]
	key = s[2]
	val = s[3]
	p.Value = my_string_with_fallback(key, val)

	input := widget.NewEntry()
	input.OnChanged = func(v string) {
		p.Value = v
	}
	input.SetText(p.Value)

	p.Item = &widget.FormItem{
		Text:     p.Label,
		HintText: p.Hint,
		Widget:   input,
	}
} // Pref_single()

// This Init creates a select/combo box FormItem
func (p *Pref_multi) Init() {

	var key string
	var val string
	s := p.Populate()
	if len(s) < 5 {
		panic("Preferences: Populate function returns too few")
	}
	p.Label = s[0]
	p.Hint = s[1]
	key = s[2]
	val = s[3]
	p.Values = s[4:]
	p.Value = my_string_with_fallback(key, val)

	combo := widget.NewSelect(p.Values, func(value string) {
		p.Value = value
	})
	combo.SetSelected(p.Value)
	p.Item = &widget.FormItem{
		Text:     p.Label,
		HintText: p.Hint,
		Widget:   combo,
	}
} // Pref_multi()

// Flavor returns what type of element this is - required to cast correctly
func (p Pref_multi) Flavor() string {
	return "multi"
}

// Flavor returns what type of element this is - required to cast correctly
func (p Pref_single) Flavor() string {
	return "single"
}

/*
** button handlers
 */

func cancel_button() {
	Window.Close()
	Window = nil
}

func reset_button() {
	Init_prefs()
	Window.SetContent(state.Prefs_form)
	Window.Show()
}

func factory_button() {
	RemoveAll_prefs()
	reset_button()
}

func save_button() {
	SaveAll()
	cancel_button()
}

/*
** Populate or data supply functions
 */

// Populate_personality returns data to create a personality FormItem
// returns a splice of strings: Label, hint, key to prefs, default value
func Populate_personality() []string {
	p := my_string_with_fallback("personality", state.Default_personality)
	c := []string{"Personality",
		"This indicates what use the app is configured for",
		"personality",
		"POD",
		p,
	}
	return c
} // Populate_personality()

// Populate_root returns data to create a root path FormItem
// returns a splice of strings: Label, hint, key to prefs, default value
func Populate_root() []string {
	hd, err := os.UserHomeDir()
	if err != nil {
		hd = "."
	}
	root := my_string_with_fallback("root", hd)
	c := []string{"Artwork Root",
		"This is the top of the folder tree that contains your artwork",
		"root",
		root,
	}
	return c
} // Populate_root()

// Populate_color returns data to create a multi line FormItem for color sets
// returns a splice of strings: Label, hint, key to prefs, list of color_sets
func Populate_color() []string {
	var c []string
	var result []string
	for v := range color_sets.Color_sets {
		c = append(c, v)
	}
	cs := my_string_with_fallback("color_set", state.Default_color)
	result = append(result, string("Color Sets: "),
		string("This is the default color set in use"),
		string("color_set"),
		cs)
	result = append(result, c...)
	return result
} // Populate_color()

// Populate_scr returns data to create a screen size FormItem
// returns a splice of strings: Label, hint, key to prefs, default value
func Populate_scr() []string {
	sz := my_string_with_fallback("scr_size", state.Default_size)
	c := []string{"Window Size (%)",
		"This is the size of the window on startup in percent of the screen",
		"scr_size",
		sz,
	}
	return c
} // Populate_scr()

// Populate_tree returns data to create a tree pane size FormItem
// returns a splice of strings: Label, hint, key to prefs, default value
func Populate_tree() []string {
	sz := my_string_with_fallback("tree_size", state.Default_tree)
	c := []string{"File Pane Size (%)",
		"This is the size of the File Pane on startup in percent of the screen",
		"tree_size",
		sz,
	}
	return c
} // Populate_tree()

/*
** Utility Functions
 */

func Get_value(key string) string {
	flavor := state.Prefs[key].(Pref_unit_type).Flavor()
	var value string
	switch flavor {
	case "multi":
		value = state.Prefs[key].(*Pref_multi).Value
	case "single":
		value = state.Prefs[key].(*Pref_single).Value
	}
	return value
}

// my_string_with_fallback is a nasty hack that overcomes the ios limitation on remove value
func my_string_with_fallback(key string, value string) string {
	stored := fyne.CurrentApp().Preferences().StringWithFallback(key, value)
	if stored == "" {
		stored = value
	}
	return stored
}

// my_remove_value is a nasty hack that overcomes the ios limitation on remove value
func my_remove_value(key string) {
	pref := fyne.CurrentApp().Preferences()
	pref.RemoveValue(key)
	test := pref.String(key)
	if test != "" {
		pref.SetString(key, "")
	}
}

// RemoveAll_prefs() remove our standard preferences
func RemoveAll_prefs() {
	for key := range state.Prefs {
		my_remove_value(key)
	}
}

// SaveAll_prefs() save our standard preferences
func SaveAll() {
	pref := fyne.CurrentApp().Preferences()
	var value string
	for key := range state.Prefs {
		value = Get_value(key)
		pref.SetString(key, value)
	}
	Init_prefs()
}
