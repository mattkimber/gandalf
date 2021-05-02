package types

import (
	"bytes"
	"encoding/binary"
)

type Shape struct {
	NodeID int
	Attributes Dictionary
	Models []int
}

const SGShape = "shape"

func (s *Shape) GetType() string {
	return SGShape
}

// Shape has no scene graph children, just models
func (s *Shape) GetChildren() []int {
	return []int{}
}

func (r *MagicaReader) GetShape() Shape {
	s := Shape{}
	s.NodeID = r.GetInt32()
	s.Attributes = r.GetDictionary()

	childNodeCount := r.GetInt32()
	childNodes := make([]int, childNodeCount)
	for i := 0; i < childNodeCount; i++ {
		childNodes[i] = r.GetInt32()
	}

	s.Models = childNodes
	return s
}


func (s *Shape) GetBytes() ([]byte, error) {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, int32(s.NodeID))
	if err != nil {
		return nil, err
	}

	dictBytes, err := s.Attributes.GetBytes()
	if err != nil {
		return nil, err
	}

	_, err = buf.Write(dictBytes)
	if err != nil {
		return nil, err
	}

	err = binary.Write(buf, binary.LittleEndian, int32(len(s.Models)))
	if err != nil {
		return nil, err
	}

	for _, model := range s.Models {
		err = binary.Write(buf, binary.LittleEndian, int32(model))
		if err != nil {
			return nil, err
		}
	}

	return buf.Bytes(), nil
}


func(s *Shape) IsChunk() bool {
	return true
}

func(s *Shape) GetChunkName() string {
	return "nSHP"
}