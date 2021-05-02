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

func WriteString(buf *bytes.Buffer, value string) error {
	err := binary.Write(buf, binary.LittleEndian, int32(len(value)))
	if err != nil {
		return  err
	}

	_, err = buf.WriteString(value)
	return err
}