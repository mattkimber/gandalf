package types

import (
	"strconv"
	"strings"
)

type Translation struct {
	NodeID int
	Dictionary Dictionary
	ChildNodeID int
	ReservedID int
	LayerID int
	Frames []Frame
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
		NodeID:         nodeID,
		Dictionary: 	dictionary,
		ChildNodeID:    childNodeID,
		ReservedID:     reservedID,
		LayerID:        layerID,
		Frames: 		frames,
	}
}