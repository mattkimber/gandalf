package types

type Palette []byte

func (r *MagicaReader) GetPalette() Palette {
	return r.buffer.Bytes()
}
