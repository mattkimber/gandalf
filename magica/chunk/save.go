package chunk

import (
	"bytes"
	"encoding/binary"
)

func (c *Chunk) GetBytes() []byte {
	allChildBytes := make([]byte, 0)

	for _, child := range c.Children {
		childBytes := child.GetBytes()
		allChildBytes = append(allChildBytes, childBytes...)
	}

	thisBytes, _ := c.Item.GetBytes()

	buf := new(bytes.Buffer)
	if c.Item.IsChunk() {
		buf.WriteString(c.Item.GetChunkName())
		binary.Write(buf, binary.LittleEndian, int32(len(thisBytes)))
		binary.Write(buf, binary.LittleEndian, int32(len(allChildBytes)))
	}

	buf.Write(thisBytes)
	buf.Write(allChildBytes)

	return buf.Bytes()
}