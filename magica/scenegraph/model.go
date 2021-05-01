package scenegraph

import (
	"github.com/mattkimber/gandalf/geometry"
	"github.com/mattkimber/gandalf/magica/types"
	"github.com/mattkimber/gandalf/utils"
)

type OutputModel struct {
	Data [][][]byte
	Size types.Size
}

func (n *Node) GetCompositeModel() OutputModel {
	extents := n.GetExtents()
	offset := extents.Min

	size := types.Size{X: extents.Max.X - offset.X, Y: extents.Max.Y - offset.Y, Z: extents.Max.Z - offset.Z}
	data := utils.Make3DByteSlice(size)

	n.AppendVoxels(offset, size, &data)

	return OutputModel {
		Size: size,
		Data: data,
	}
}

func (n *Node) AppendVoxels(offset geometry.Point, size types.Size, data *[][][]byte) {
	for _, model := range n.Models {
		for _, p := range model.Points {
			x := p.Point.X + n.Location.X - offset.X
			y := p.Point.Y + n.Location.Y - offset.Y
			z := p.Point.Z + n.Location.Z - offset.Z

			if x < size.X && y < size.Y && z < size.Z && p.Colour != 0 {
				(*data)[x][y][z] = p.Colour
			}
		}
	}

	for _, child := range n.Children {
		child.AppendVoxels(offset, size, data)
	}
}
