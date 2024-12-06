package state

// ThreadChg encapsulates a change line tag eg "CC1", order eg 1, thread of color_type
type ThreadChg struct {
	ChgTag string
	Order  uint
	Thread Color_type
}

// These are the design params from an Embroidery Library PDF
type Design_params struct {
	Stitches      uint
	Size          string
	ColChngs      uint
	ColUsed       uint
	ThreadChanges []ThreadChg
}

type Emb_type struct {
	About   About_type
	Design  Design_params
	Threads []Color_type
}

func Empty_emb() Emb_type {
	e := new(Emb_type)
	e.Threads = []Color_type{}
	e.Design.ThreadChanges = []ThreadChg{}
	return *e
}
