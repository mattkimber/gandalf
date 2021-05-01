package scenegraph

import (
	"github.com/mattkimber/gandalf/geometry"
	"reflect"
	"testing"
)

func TestExtents_GetBounds(t *testing.T) {
	tests := []struct {
		name       string
		e          Extents
		wantResult Extent
	}{
		{
			name: "identity",
			e: Extents{{Min: geometry.Point{X: 1, Y: 2, Z: 3}, Max: geometry.Point{X: 4, Y: 5, Z: 6}}},
			wantResult: Extent{Min: geometry.Point{X: 1, Y: 2, Z: 3}, Max:geometry.Point{X: 4, Y: 5, Z: 6}},
		},
		{
			name: "two extents",
			e: Extents{
				{Min: geometry.Point{X: 1, Y: 2, Z: 3}, Max: geometry.Point{X: 4, Y: 5, Z: 6}},
				{Min: geometry.Point{X: -1, Y: 4, Z: -3}, Max: geometry.Point{X: 7, Y: 8, Z: 2}},
			},
			wantResult: Extent{Min: geometry.Point{X: -1, Y: 2, Z: -3}, Max:geometry.Point{X: 7, Y: 8, Z: 6}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotResult := tt.e.GetBounds(); !reflect.DeepEqual(gotResult, tt.wantResult) {
				t.Errorf("GetBounds() = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}