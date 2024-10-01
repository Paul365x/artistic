package state

import (
	"fmt"
	"testing"
)

func TestInstEmpty(t *testing.T) {
	it := Empty_instance()
	if it.Image != "" {
		t.Fatalf(`Instance failed empty Image`)
	}
	st := fmt.Sprintf("%T", it.BG)
	fmt.Println(st)
	if st != "state.Color_type" {
		t.Fatalf(`Instance failed empty background`)
	}
}

func TestMetaEmpty(t *testing.T) {
	mt := Empty_meta()
	at := fmt.Sprintf("%T", mt.About)
	sr := fmt.Sprintf("%T", mt.Search_data)

	if at != "state.About_type" {
		t.Fatalf(`Meta failed empty About`)
	}
	if sr != "state.Search_type" {
		t.Fatalf(`Meta failed empty Search_data`)
	}
}

func TestPodEmpty(t *testing.T) {
	art := Empty_pod()
	at := fmt.Sprintf("%T", art.Artwork)
	mt := fmt.Sprintf("%T", art.Metadata)

	if at != "state.Artwork_type" {
		t.Fatalf(`Art failed empty Artwork`)
	}
	if mt != "state.Meta_type" {
		t.Fatalf(`Art failed empty Metadata`)
	}
}
