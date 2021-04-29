package types

type Shape struct {
	NodeID int
	Attributes Dictionary
	Models []int
}

const SGShape = "shape"

func (s *Shape) GetType() string {
	return SGShape
}

// Shape has no scene graph children, just models
func (s *Shape) GetChildren() []int {
	return []int{}
}

func (r *MagicaReader) GetShape() Shape {
	s := Shape{}
	s.NodeID = r.GetInt32()
	s.Attributes = r.GetDictionary()

	childNodeCount := r.GetInt32()
	childNodes := make([]int, childNodeCount)
	for i := 0; i < childNodeCount; i++ {
		childNodes[i] = r.GetInt32()
	}

	s.Models = childNodes
	return s
}
