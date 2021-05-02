package types

import (
	"reflect"
	"testing"
)

func TestMagicaReader_GetDictionary(t *testing.T) {
	dict := Dictionary{Values: map[string]string{
		"a": "foo",
		"b": "bar",
	}}

	bytes, err := dict.GetBytes()
	if err != nil {
		t.Errorf("error saving dictionary: %s", err)
	}

	rd := GetReader(bytes)
	dict2 := rd.GetDictionary()

	if !reflect.DeepEqual(dict, dict2) {
		t.Errorf("dictionary did not survive round trip. In: %+v, Out: %+v", dict, dict2)
	}
}