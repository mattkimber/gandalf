package chunk


type ChunkItem interface {
	GetBytes() ([]byte, error)
	IsChunk() bool
	GetChunkName() string
}

type Chunk struct {
	Item ChunkItem
	Children []Chunk
}

type MainChunk struct{}

func(m *MainChunk) GetBytes() ([]byte, error) {
	return []byte{}, nil
}

func(m *MainChunk) IsChunk() bool {
	return true
}

func(m *MainChunk) GetChunkName() string {
	return "MAIN"
}

