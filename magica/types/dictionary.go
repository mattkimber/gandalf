package types

import (
	"bytes"
	"encoding/binary"
)

type Dictionary struct {
	Values map[string]string
}

func (r *MagicaReader) GetDictionary() Dictionary {
	items := r.GetInt32()
	values := make(map[string]string)
	for i := 0; i < items; i++ {
		key := r.GetString()
		value := r.GetString()
		values[key] = value
	}

	return Dictionary{Values: values}
}

func (d *Dictionary) GetBytes() ([]byte, error) {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, int32(len(d.Values)))
	if err != nil {
		return nil, err
	}

	for key, value := range d.Values {
		err = WriteString(buf, key)
		if err != nil {
			return nil, err
		}

		err = WriteString(buf, value)
		if err != nil {
			return nil, err
		}
	}

	return buf.Bytes(), nil
}

func(d *Dictionary) IsChunk() bool {
	return false
}

func(d *Dictionary) GetChunkName() string {
	return ""
}