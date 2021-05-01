package magica

import (
	"github.com/mattkimber/gandalf/geometry"
	"github.com/mattkimber/gandalf/magica/types"
	"github.com/mattkimber/gandalf/utils"
)

type VoxelData [][][]byte

type VoxelObject struct {
	Voxels      [][][]byte
	PaletteData []byte
	Size        geometry.Point
}

func (v *VoxelObject) GetPoints() (result []geometry.PointWithColour) {
	ct := 0

	v.Iterate(func(x, y, z int) {
		if v.Voxels[x][y][z] != 0 {
			ct++
		}
	})

	result = make([]geometry.PointWithColour, ct)
	ct = 0

	v.Iterate(func(x, y, z int) {
		if v.Voxels[x][y][z] != 0 {
			result[ct] = geometry.PointWithColour{
				Point:  geometry.Point{X: x, Y: y, Z: z},
				Colour: v.Voxels[x][y][z],
			}
			ct++
		}
	})

	return result
}

// NewVoxelObject returns an empty voxel object of the specified size and palette
func NewVoxelObject(size geometry.Point, palette []byte) VoxelObject {
	voxelData := make([][][]byte, size.X)
	for x := 0; x < size.X; x++ {
		voxelData[x] = make([][]byte, size.Y)
		for y := 0; y < size.Y; y++ {
			voxelData[x][y] = make([]byte, size.Z)
		}
	}

	v := VoxelObject{
		Voxels:      voxelData,
		PaletteData: palette,
		Size:        size,
	}

	return v
}

func (v *VoxelObject) Copy() (result VoxelObject) {
	result = VoxelObject{}

	result.Size = v.Size
	// We don't do anything with the palette data, so a shallow copy is okay
	result.PaletteData = v.PaletteData

	result.Voxels = utils.Make3DByteSlice(types.Size{X: v.Size.X, Y: v.Size.Y, Z: v.Size.Z})
	v.Iterate(func(x, y, z int) { result.Voxels[x][y][z] = v.Voxels[x][y][z] })

	return
}

// Set sets the voxel at loc to index i
func (v *VoxelObject) Set(loc geometry.Point, i byte) {
	v.Voxels[loc.X][loc.Y][loc.Z] = i
}

// Get gets the voxel index at loc
func (v *VoxelObject) Get(loc geometry.Point) byte {
	return v.Voxels[loc.X][loc.Y][loc.Z]
}

// SafeSet sets the voxel at loc to index i, if loc is in bounds
func (v *VoxelObject) SafeSet(loc geometry.Point, i byte) {
	if loc.IsInBounds(geometry.Bounds{Max: v.Size}) {
		v.Voxels[loc.X][loc.Y][loc.Z] = i
	}
}

// SafeGet gets the voxel at loc if loc is in bounds, or 0 if it isn't
func (v *VoxelObject) SafeGet(loc geometry.Point) byte {
	if loc.IsInBounds(geometry.Bounds{Max: v.Size}) {
		return v.Voxels[loc.X][loc.Y][loc.Z]
	}

	return 0
}

func (v *VoxelObject) Iterate(iterator func(int, int, int)) {
	for x := 0; x < v.Size.X; x++ {
		for y := 0; y < v.Size.Y; y++ {
			for z := 0; z < v.Size.Z; z++ {
				iterator(x, y, z)
			}
		}
	}
}
