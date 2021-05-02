package types

type Palette []byte

func (r *MagicaReader) GetPalette() Palette {
	return r.buffer.Bytes()
}

func (p *Palette) GetBytes() ([]byte, error) {
	return *p, nil
}

func(p *Palette) IsChunk() bool {
	return true
}

func(p *Palette) GetChunkName() string {
	return "RGBA"
}