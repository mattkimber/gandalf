package types

import (
	"reflect"
	"testing"
)

func TestMagicaReader_GetTranslation(t *testing.T) {
	translation := Translation{
		NodeID: 1,
		Attributes: Dictionary{
			Values: map[string]string{
				"foo": "bar",
			},
		},
		ChildNodeID: 2,
		ReservedID:  1,
		LayerID:     5,
		Frames:      []Frame{
			{
				X: 12,
				Y: -3,
				Z: 54,
			},
		},
	}

	bytes, err := translation.GetBytes()
	if err != nil {
		t.Errorf("error saving translation: %s", err)
	}

	rd := GetReader(bytes)
	translation2 := rd.GetTranslation()

	if !reflect.DeepEqual(translation, translation2) {
		t.Errorf("translation did not survive round trip. In: %+v, Out: %+v", translation, translation2)
	}
}