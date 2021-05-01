package scenegraph

import (
	"github.com/mattkimber/gandalf/geometry"
	"math"
)

type Extent struct {
	Min geometry.Point
	Max geometry.Point
}

type Extents []Extent

func min(a int, b int) int {
	if a < b {
		return a
	}

	return b
}

func max(a int, b int) int {
	if a > b {
		return a
	}

	return b
}

func (e *Extents) GetBounds() Extent {
	result := Extent{
		Min: geometry.Point{X: math.MaxInt32, Y: math.MaxInt32, Z: math.MaxInt32 },
		Max: geometry.Point{X: math.MinInt32, Y: math.MinInt32, Z: math.MinInt32 },
	}

	for _, ex := range *e {
		result.Min.X = min(result.Min.X, ex.Min.X)
		result.Min.Y = min(result.Min.Y, ex.Min.Y)
		result.Min.Z = min(result.Min.Z, ex.Min.Z)

		result.Max.X = max(result.Max.X, ex.Max.X)
		result.Max.Y = max(result.Max.Y, ex.Max.Y)
		result.Max.Z = max(result.Max.Z, ex.Max.Z)
	}

	return result
}
