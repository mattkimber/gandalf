package types

type SceneGraphItem interface {
	GetType() string
	GetChildren() []int
	GetBytes() ([]byte, error)
	IsChunk() bool
	GetChunkName() string
}
