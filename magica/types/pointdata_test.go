package types

import (
	"github.com/mattkimber/gandalf/geometry"
	"reflect"
	"testing"
)

func TestPointData_GetPointData(t *testing.T) {
	pointData := PointData{
		{
			Point:  geometry.Point{X:1, Y: 2, Z: 3},
			Colour: 4,
		},
		{
			Point:  geometry.Point{X:3, Y: 2, Z: 1},
			Colour: 5,
		},
	}

	bytes, err := pointData.GetBytes()
	if err != nil {
		t.Errorf("error saving point data: %s", err)
	}

	// Something weird goes on with the point data... this shouldn't be here?
	rd := GetReader(bytes[4:])
	pointData2 := rd.GetPointData()

	if !reflect.DeepEqual(pointData, pointData2) {
		t.Errorf("point data did not survive round trip. In: %+v, Out: %+v", pointData, pointData2)
	}
}