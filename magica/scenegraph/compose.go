package scenegraph

import (
	"github.com/mattkimber/gandalf/geometry"
	"github.com/mattkimber/gandalf/magica/types"
)

func (n *Node) Decompose() (graph Map, pointData []types.PointData, sizeData []types.Size) {
	id := 0
	shapeID := 0

	graph = make(Map)
	pointData = make([]types.PointData, 0)
	sizeData = make([]types.Size, 0)

	_ = n.decomposeWithIDs(&id, &shapeID, graph, &pointData, &sizeData)
	return
}

func (n *Node) decomposeWithIDs(id, shapeID *int, graph Map, pointData *[]types.PointData, sizeData *[]types.Size) (rootID int) {
	rootTranslation := types.Translation{
		NodeID:     *id,
		Attributes: types.Dictionary{},
		ReservedID: -1,
		LayerID: 0,
		Frames: []types.Frame{{X: 0, Y: 0, Z: 0}},
	}
	graph[*id] = &rootTranslation
	*id++

	rootGroup := types.Group{
		NodeID:     *id,
		Attributes: types.Dictionary{},
	}

	rootTranslation.ChildNodeID = rootGroup.NodeID

	graph[*id] = &rootGroup
	*id++

	// If this node has models, add them to the point and size data
	if len(n.Models) > 0 {
		translationNodeIDs := make([]int, 0)

		for _, model := range n.Models {
			*pointData = append(*pointData, model.Points)
			*sizeData = append(*sizeData, model.Size)

			shp := types.Shape{
				NodeID:     *id,
				Attributes: types.Dictionary{},
				Models:     []int{*shapeID},
			}
			graph[*id] = &shp
			*id++
			*shapeID++

			trn := types.Translation{
				NodeID:      *id,
				Attributes:  types.Dictionary{},
				ChildNodeID: shp.NodeID,
				ReservedID:  -1,
				LayerID:     0,
				Frames: []types.Frame{{
					X: n.Location.X,
					Y: n.Location.Y,
					Z: n.Location.Z,
				},
				},
			}
			graph[*id] = &trn
			translationNodeIDs = append(translationNodeIDs, *id)
			*id++
		}

		rootGroup.ChildNodes = translationNodeIDs
	} else {
		childNodeIDs := make([]int, 0)

		for _, node := range n.Children {
			var childID int
			childID = node.decomposeWithIDs(id, shapeID, graph, pointData, sizeData)
			childNodeIDs = append(childNodeIDs, childID)
		}

		rootGroup.ChildNodes = childNodeIDs
	}

	return rootTranslation.NodeID
}

func Compose(graph Map, current types.SceneGraphItem, x, y, z int, pointData []types.PointData, sizeData []types.Size) (result Node) {
	if current.GetType() == types.SGTranslation {
		tn := current.(*types.Translation)
		for _, frame := range tn.Frames {
			x += frame.X
			y += frame.Y
			z += frame.Z
		}
	}

	if current.GetType() == types.SGShape {
		shp := current.(*types.Shape)
		size := types.Size{}
		models := make([]Model, len(shp.Models))
		for idx, child := range shp.Models {
			models[idx] = Model{Points: pointData[child], Size: sizeData[child]}
			size = sizeData[child]
		}
		result.Models = models
		result.Location = geometry.Point{X: x - (size.X / 2), Y: y - (size.Y / 2), Z: z - (size.Z / 2)}
	}

	children := make([]Node, 0)
	for _, child := range current.GetChildren() {
		next, ok := graph[child]; if ok {
			children = append(children, Compose(graph, next, x, y, z, pointData, sizeData))
		}
	}

	result.Children = children

	return result
}
