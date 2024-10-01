package color_sets

import (
	"image/color"
	"maps"
	"math"

	"github.com/artistic/internal/notify"
	"github.com/artistic/internal/state"
)

// Exhaustive list of colors
var (
	ArmyR             = color.RGBA{0x3e, 0x52, 0x2b, 255}
	BlackT            = color.Gray16{0}
	BlackR            = color.RGBA{0x10, 0x10, 0x10, 255}
	BlueR             = color.RGBA{0x2f, 0x4c, 0xb5, 255}
	ButterYellowR     = color.RGBA{0xff, 0xcf, 0x6e, 255}
	CremeT            = color.RGBA{0xea, 0xe0, 0xc7, 255}
	CremeR            = color.RGBA{0xe5, 0xd6, 0xc5, 255}
	DenimHeatherR     = color.RGBA{0x3c, 0x3b, 0x43, 255}
	DarkGreyR         = color.RGBA{0x4a, 0x44, 0x40, 255}
	DarkRedR          = color.RGBA{0x5a, 0x1f, 0x32, 255}
	ForestGreenR      = color.RGBA{0x13, 0x29, 0x0c, 255}
	GoldR             = color.RGBA{0xf8, 0x9f, 0x2b, 255}
	GreenR            = color.RGBA{0x0f, 0x8c, 0x37, 255}
	HeatherT          = color.RGBA{0xd1, 0xd1, 0xd1, 255}
	HeatherGreyR      = color.RGBA{0xb6, 0xb6, 0xb6, 255}
	KiwiR             = color.RGBA{0xa1, 0xc7, 0x40, 255}
	LightBlueT        = color.RGBA{0xc8, 0xe0, 0xec, 255}
	LightBlueR        = color.RGBA{0xbe, 0xd4, 0xe8, 255}
	LightPinkR        = color.RGBA{0xff, 0xc7, 0xca, 255}
	NavyR             = color.RGBA{0x0d, 0x16, 0x2e, 255}
	OrangeT           = color.RGBA{0xdc, 0x44, 0x05, 255}
	PurpleR           = color.RGBA{0x53, 0x19, 0x63, 255}
	RedT              = color.RGBA{0xc6, 0x2b, 0x29, 255}
	RedR              = color.RGBA{0xdd, 0x21, 0x21, 255}
	RedHeatherT       = color.RGBA{0xcf, 0x42, 0x53, 255}
	RoyalBlueT        = color.RGBA{0x36, 0x53, 0x8b, 255}
	RoyalHeatherT     = color.RGBA{0x37, 0x67, 0xa0, 255}
	SlateT            = color.RGBA{0x76, 0x8e, 0x9a, 255}
	SoftPinkT         = color.RGBA{0xfa, 0xc2, 0xcd, 255}
	TealT             = color.RGBA{0x01, 0x95, 0xc3, 255}
	TurquoiseHeatherT = color.RGBA{0x33, 0x8f, 0xb1, 255}
	WhiteR            = color.RGBA{0xfa, 0xfa, 0xfa, 255}
	WhiteT            = color.Gray16{0xffff}
	YellowT           = color.RGBA{0xff, 0xb8, 0x1c, 255}
)

// Color_sets is a jumptable to the colorset funcions
var Color_sets = make(map[string]func() *map[string]color.Color)

// build_sets initialises the Color_sets jumptable
func Build_sets() {
	Color_sets["TEEPUBLIC"] = Teepee_set
	Color_sets["REDBUBBLE"] = CommieGlobe_set
	Color_sets["All"] = Complete_set
}

// Load_set takes a key and returns a function that in turn creates a colorset
func Load_set(set_name string) func() *map[string]color.Color {
	set, ok := Color_sets[set_name]
	if ok {
		return set
	} else {
		notify.Notify(string("Bad Color_set: ")+
			set_name+
			string(" defaulting to ")+
			state.Default_color,
			"error",
			state.Error)
		return Color_sets[state.Default_color]
	}
}

// teepee_set returns the Teepublic set of colors
func Teepee_set() *map[string]color.Color {
	return &map[string]color.Color{
		"TP_Black":             BlackT,
		"TP_White":             WhiteT,
		"TP_Red":               RedT,
		"TP_Soft Pink":         SoftPinkT,
		"TP_Orange":            OrangeT,
		"TP_Yellow":            YellowT,
		"TP_Creme":             CremeT,
		"TP_Royal Blue":        RoyalBlueT,
		"TP_Teal":              TealT,
		"TP_Slate":             SlateT,
		"TP_Light Blue":        LightBlueT,
		"TP_Heather":           HeatherT,
		"TP_Red Heather":       RedHeatherT,
		"TP_Turquoise Heather": TurquoiseHeatherT,
		"TP_Royal Heather":     RoyalHeatherT,
	}
} //teepee_set()

// commieglobe_set returns the RedBubble set of colors
func CommieGlobe_set() *map[string]color.Color {
	return &map[string]color.Color{
		"CB_Black":         BlackR,
		"CB_White":         WhiteR,
		"CB_Red":           RedR,
		"CB_Heather Grey":  HeatherGreyR,
		"CB_Denim Heather": DenimHeatherR,
		"CB_Navy":          NavyR,
		"CB_Blue":          BlueR,
		"CB_Creme":         CremeR,
		"CB_Light Blue":    LightBlueR,
		"CB_Dark Grey":     DarkGreyR,
		"CB_Kiwi":          KiwiR,
		"CB_Green":         GreenR,
		"CB_Army":          ArmyR,
		"CB_Forest Green":  ForestGreenR,
		"CB_Light Pink":    LightPinkR,
		"CB_Purple":        PurpleR,
		"CB_Dark Red":      DarkRedR,
		"CB_Butter Yellow": ButterYellowR,
		"CB_Gold":          GoldR,
	}
} //CommieGlobe_set()

// Complete_set returns all the colours known/hooked into a set
func Complete_set() *map[string]color.Color {
	c := make(map[string]color.Color)
	for k, v := range Color_sets {
		if k == "All" {
			continue
		}
		s := v()
		maps.Copy(c, *s)
	}
	return &c
} //Complete_set()

// compare_color calculates the distance between two colors and their opacity
func Compare_color(ref, test color.Color) (int64, int64) {
	//rgba := color.RGBAModel.Convert(test).(color.RGBA)
	//mt.Printf("%d %d %d\n",rgba.R, rgba.G, rgba.B)
	a := color.RGBAModel.Convert(ref).(color.RGBA)
	b := color.RGBAModel.Convert(test).(color.RGBA)
	var r_diff = float64(a.R) - float64(b.R)
	var g_diff = float64(a.G) - float64(b.G)
	var b_diff = float64(a.B) - float64(b.B)
	diff := math.Round(math.Sqrt((r_diff * r_diff) + (g_diff * g_diff) + (b_diff * b_diff)))
	alpha := int64(a.A - b.A)
	return int64(diff), alpha

} // compare_color()
