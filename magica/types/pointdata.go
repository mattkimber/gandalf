package types

import "github.com/mattkimber/gandalf/geometry"

type PointData []geometry.PointWithColour

func (r *MagicaReader) GetPointData() PointData {
	data := r.buffer.Bytes()
	result := make([]geometry.PointWithColour, len(data)/4)


	for i := 0; i + 4 <= len(data); i += 4 {
		point := geometry.PointWithColour{
			Point: geometry.Point{X: int(data[i]), Y: int(data[i+1]), Z: int(data[i+2])}, Colour: data[i+3],
		}

		result[i/4] = point
	}

	return result
}
