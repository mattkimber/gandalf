package types

import (
	"bytes"
	"encoding/binary"
	"github.com/mattkimber/gandalf/geometry"
)

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

func (p *PointData) GetBytes() ([]byte, error) {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, int32(len(*p)))
	for _, pt := range *p {
		_, err := buf.Write([]byte{
			byte(pt.Point.X),
			byte(pt.Point.Y),
			byte(pt.Point.Z),
			pt.Colour,
		})
		if err != nil {
			return nil, err
		}
	}

	return buf.Bytes(), nil
}


func(p *PointData) IsChunk() bool {
	return true
}

func(p *PointData) GetChunkName() string {
	return "XYZI"
}