package types

import (
	"bytes"
	"encoding/binary"
)

type MagicaReader struct {
	buffer *bytes.Buffer
}

func GetReader(data []byte) MagicaReader {
	buf := bytes.NewBuffer(data)
	return MagicaReader{
		buffer: buf,
	}
}

func (r *MagicaReader) GetString() string {
	size := r.GetInt32()
	bytes := r.buffer.Next(size)
	return string(bytes)
}

func  (r *MagicaReader) GetInt32() int {
	bytes := r.buffer.Next(4)
	return int(binary.LittleEndian.Uint32(bytes[0:4]))
}