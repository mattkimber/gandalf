package types

import (
	"bytes"
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

func (s *Size) GetBytes() ([]byte, error) {
	buf := new(bytes.Buffer)

	err := binary.Write(buf, binary.LittleEndian, int32(s.X))
	if err != nil {
		return nil, err
	}

	err = binary.Write(buf, binary.LittleEndian, int32(s.Y))
	if err != nil {
		return nil, err
	}

	err = binary.Write(buf, binary.LittleEndian, int32(s.Z))
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}


func(s *Size) IsChunk() bool {
	return true
}

func(s *Size) GetChunkName() string {
	return "SIZE"
}