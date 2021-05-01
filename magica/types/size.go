package types

import (
	"encoding/binary"
	"github.com/mattkimber/gandalf/geometry"
)

type Size geometry.Point

func (r *MagicaReader) GetSize() Size {
	data := r.buffer.Bytes()

	return Size{
		X: int(binary.LittleEndian.Uint32(data[0:4])),
		Y: int(binary.LittleEndian.Uint32(data[4:8])),
		Z: int(binary.LittleEndian.Uint32(data[8:12])),
	}
}
