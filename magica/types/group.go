package types

import (
	"bytes"
	"encoding/binary"
)

type Group struct {
	NodeID int
	Attributes Dictionary
	ChildNodes []int
}

const SGGroup = "group"

func (g *Group) GetType() string {
	return SGGroup
}

func (g *Group) GetChildren() []int {
	return g.ChildNodes
}

func (r *MagicaReader) GetGroup() Group {
	g := Group{}
	g.NodeID = r.GetInt32()
	g.Attributes = r.GetDictionary()

	childNodeCount := r.GetInt32()
	childNodes := make([]int, childNodeCount)
	for i := 0; i < childNodeCount; i++ {
		childNodes[i] = r.GetInt32()
	}

	g.ChildNodes = childNodes
	return g
}


func (g *Group) GetBytes() ([]byte, error) {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, int32(g.NodeID))
	if err != nil {
		return nil, err
	}

	dictBytes, err := g.Attributes.GetBytes()
	if err != nil {
		return nil, err
	}

	_, err = buf.Write(dictBytes)
	if err != nil {
		return nil, err
	}

	err = binary.Write(buf, binary.LittleEndian, int32(len(g.ChildNodes)))
	if err != nil {
		return nil, err
	}

	for _, node := range g.ChildNodes {
		err = binary.Write(buf, binary.LittleEndian, int32(node))
		if err != nil {
			return nil, err
		}
	}

	return buf.Bytes(), nil
}

func(g *Group) IsChunk() bool {
	return true
}

func(g *Group) GetChunkName() string {
	return "nGRP"
}