package types

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strconv"
	"strings"
)

type Translation struct {
	NodeID      int
	Attributes  Dictionary
	ChildNodeID int
	ReservedID  int
	LayerID     int
	Frames      []Frame
}

const SGTranslation = "translation"

func (t *Translation) GetType() string {
	return SGTranslation
}

func (t *Translation) GetChildren() []int {
	return []int{t.ChildNodeID}
}

type Frame struct {
	X int
	Y int
	Z int
}

func (r *MagicaReader) GetTranslation() Translation {
	nodeID := r.GetInt32()
	dictionary := r.GetDictionary()
	childNodeID := r.GetInt32()
	reservedID := r.GetInt32()
	layerID := r.GetInt32()
	numFrames := r.GetInt32()

	frames := make([]Frame, 0)

	for i := 0; i < numFrames; i++ {
		dict := r.GetDictionary()

		if t, ok := dict.Values["_t"]; ok {
			values := strings.Split(t, " ")
			if len(values) >= 3 {
				x, _ := strconv.Atoi(values[0])
				y, _ := strconv.Atoi(values[1])
				z, _ := strconv.Atoi(values[2])

				frames = append(frames, Frame{X: x, Y: y, Z:z})
			}
		}
	}

	return Translation{
		NodeID:      nodeID,
		Attributes:  dictionary,
		ChildNodeID: childNodeID,
		ReservedID:  reservedID,
		LayerID:     layerID,
		Frames:      frames,
	}
}


func (t *Translation) GetBytes() ([]byte, error) {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, int32(t.NodeID))
	if err != nil {
		return nil, err
	}

	dictBytes, err := t.Attributes.GetBytes()
	if err != nil {
		return nil, err
	}

	_, err = buf.Write(dictBytes)
	if err != nil {
		return nil, err
	}

	err = binary.Write(buf, binary.LittleEndian, int32(t.ChildNodeID))
	if err != nil {
		return nil, err
	}

	err = binary.Write(buf, binary.LittleEndian, int32(t.ReservedID))
	if err != nil {
		return nil, err
	}

	err = binary.Write(buf, binary.LittleEndian, int32(t.LayerID))
	if err != nil {
		return nil, err
	}

	err = binary.Write(buf, binary.LittleEndian, int32(len(t.Frames)))
	if err != nil {
		return nil, err
	}

	for _, frame := range t.Frames {
		dictionary := Dictionary{}

		if frame.X != 0 || frame.Y != 0 || frame.Z != 0 {
			dictionary.Values = map[string]string{
				"_t": fmt.Sprintf("%d %d %d", frame.X, frame.Y, frame.Z),
			}
		}

		dictBytes, err = dictionary.GetBytes()
		if err != nil {
			return nil, err
		}

		_, err = buf.Write(dictBytes)
		if err != nil {
			return nil, err
		}
	}

	return buf.Bytes(), nil
}

func(t *Translation) IsChunk() bool {
	return true
}

func(t *Translation) GetChunkName() string {
	return "nTRN"
}