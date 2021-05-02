package types

import (
	"reflect"
	"testing"
)

func TestMagicaReader_GetShape(t *testing.T) {
	shape := Shape{
		NodeID:     1,
		Attributes: Dictionary{
			Values: map[string]string{
				"foo": "bar",
			},
		},
		Models: []int{2,3,4},
	}

	bytes, err := shape.GetBytes()
	if err != nil {
		t.Errorf("error saving shape: %s", err)
	}

	rd := GetReader(bytes)
	shape2 := rd.GetShape()

	if !reflect.DeepEqual(shape, shape2) {
		t.Errorf("shape did not survive round trip. In: %+v, Out: %+v", shape, shape2)
	}
}