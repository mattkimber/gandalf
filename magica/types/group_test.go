package types

import (
	"reflect"
	"testing"
)

func TestGroup_GetGroup(t *testing.T) {
	group := Group{
		NodeID:     1,
		Attributes: Dictionary{
			Values: map[string]string{
				"foo": "bar",
			},
		},
		ChildNodes: []int{2,3,4},
	}

	bytes, err := group.GetBytes()
	if err != nil {
		t.Errorf("error saving group: %s", err)
	}

	rd := GetReader(bytes)
	group2 := rd.GetGroup()

	if !reflect.DeepEqual(group, group2) {
		t.Errorf("group did not survive round trip. In: %+v, Out: %+v", group, group2)
	}
}