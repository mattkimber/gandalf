package types

import (
	"reflect"
	"testing"
)

func TestPalette_GetPalette(t *testing.T) {
	palette := Palette{1,2,3,4,5,6,7,8,9,10,11,12}

	bytes, err := palette.GetBytes()
	if err != nil {
		t.Errorf("error saving palette: %s", err)
	}

	rd := GetReader(bytes)
	palette2 := rd.GetPalette()

	if !reflect.DeepEqual(palette, palette2) {
		t.Errorf("palette did not survive round trip. In: %+v, Out: %+v", palette, palette2)
	}
}