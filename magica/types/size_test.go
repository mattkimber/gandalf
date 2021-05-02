package types

import (
	"reflect"
	"testing"
)

func TestMagicaReader_GetSize(t *testing.T) {
	size := Size{
		X: 10,
		Y: 20,
		Z: 15,
	}

	bytes, err := size.GetBytes()
	if err != nil {
		t.Errorf("error saving size: %s", err)
	}

	rd := GetReader(bytes)
	size2 := rd.GetSize()

	if !reflect.DeepEqual(size, size2) {
		t.Errorf("size did not survive round trip. In: %+v, Out: %+v", size, size2)
	}
}