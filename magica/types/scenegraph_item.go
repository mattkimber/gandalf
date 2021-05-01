package types

type SceneGraphItem interface {
	GetType() string
	GetChildren() []int
}
