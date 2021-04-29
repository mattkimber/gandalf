package types

type Group struct {
	NodeID int
	Attributes Dictionary
	ChildNodes []int
}

const SGGroup = "group"

func (g *Group) GetType() string {
	return SGGroup
}

func (g *Group) GetChildren() []int {
	return g.ChildNodes
}

func (r *MagicaReader) GetGroup() Group {
	g := Group{}
	g.NodeID = r.GetInt32()
	g.Attributes = r.GetDictionary()

	childNodeCount := r.GetInt32()
	childNodes := make([]int, childNodeCount)
	for i := 0; i < childNodeCount; i++ {
		childNodes[i] = r.GetInt32()
	}

	g.ChildNodes = childNodes
	return g
}
