package magica

import (
	"encoding/binary"
	"fmt"
	"github.com/mattkimber/gandalf/geometry"
	"github.com/mattkimber/gandalf/magica/types"
	"github.com/mattkimber/gandalf/utils"
	"io"
	"io/ioutil"
	"os"
)

const magic = "VOX "

func isHeaderValid(handle io.Reader) bool {
	result, err := getChunkHeader(handle)
	return err == nil && result == magic
}

func getChunkHeader(handle io.Reader) (string, error) {
	limitedReader := io.LimitReader(handle, 4)
	result, err := ioutil.ReadAll(limitedReader)
	return string(result), err
}

func getSizeFromChunk(handle io.Reader) (geometry.Point, error) {
	data, err := getChunkData(handle, 12)

	if err != nil {
		return geometry.Point{}, err
	}

	return geometry.Point{
		X: int(binary.LittleEndian.Uint32(data[0:4])),
		Y: int(binary.LittleEndian.Uint32(data[4:8])),
		Z: int(binary.LittleEndian.Uint32(data[8:12])),
	}, nil
}

func getPointDataFromChunk(handle io.Reader) ([]geometry.PointWithColour, error) {
	data, err := getChunkData(handle, 4)

	if err != nil {
		return getNilValueForPointDataFromChunk(), err
	}

	result := make([]geometry.PointWithColour, len(data)/4)

	for i := 0; i < len(data); i += 4 {
		point := geometry.PointWithColour{
			Point: geometry.Point{X: int(data[i]), Y: int(data[i+1]), Z: int(data[i+2])}, Colour: data[i+3],
		}

		result[i/4] = point
	}

	return result, nil
}

func getPaletteDataFromChunk(handle io.Reader) (data []byte, err error) {
	data, err = getChunkData(handle, 0)
	return
}

func getTranslationDataFromChunk(handle io.Reader) (tx types.Translation, err error) {
	data, err := getChunkData(handle, 4)

	if err != nil {
		return
	}

	rd := types.GetReader(data)

	return rd.GetTranslation(), nil
}

func getGroupDataFromChunk(handle io.Reader) (g types.Group, err error) {
	data, err := getChunkData(handle, 4)

	if err != nil {
		return
	}

	rd := types.GetReader(data)

	return rd.GetGroup(), nil
}

func getShapeDataFromChunk(handle io.Reader) (s types.Shape, err error) {
	data, err := getChunkData(handle, 4)

	if err != nil {
		return
	}

	rd := types.GetReader(data)

	return rd.GetShape(), nil
}

func getVoxelObjectFromPointData(size geometry.Point, data []geometry.PointWithColour) VoxelData {
	result := utils.Make3DByteSlice(size)

	for _, p := range data {
		if p.Point.X < size.X && p.Point.Y < size.Y && p.Point.Z < size.Z && p.Colour != 0 {
			result[p.Point.X][p.Point.Y][p.Point.Z] = p.Colour
		}
	}

	return result
}

func skipUnhandledChunk(handle io.Reader) {
	_, _ = getChunkData(handle, 0)
}

func getChunkData(handle io.Reader, minSize int64) ([]byte, error) {
	parsedSize := getSize(handle)

	// Still need to read to the end even if the size
	// is invalid
	limitedReader := io.LimitReader(handle, parsedSize)
	data, err := ioutil.ReadAll(limitedReader)

	if parsedSize < minSize {
		return nil, fmt.Errorf("invalid chunk size")
	}

	if int64(len(data)) < parsedSize {
		return nil, fmt.Errorf("chunk size declared %d but was %d", parsedSize, len(data))
	}

	return data, err
}

func getSize(handle io.Reader) int64 {
	limitedReader := io.LimitReader(handle, 8)
	size, err := ioutil.ReadAll(limitedReader)

	if err != nil {
		return 0
	}

	parsedSize := int64(binary.LittleEndian.Uint32(size[0:4]))
	return parsedSize
}

func getNilValueForPointDataFromChunk() []geometry.PointWithColour {
	return []geometry.PointWithColour{}
}

func GetMagicaVoxelObject(handle io.Reader) (VoxelObject, error) {
	if !isHeaderValid(handle) {
		return VoxelObject{}, fmt.Errorf("header not valid")
	}
	getChunkHeader(handle)

	size := geometry.Point{}

	sizes := make([]geometry.Point, 0)
	allPointData := make([][]geometry.PointWithColour, 0)

	scenegraph := make(map[int]types.SceneGraphItem)

	var paletteData []byte

	for {
		chunkType, err := getChunkHeader(handle)

		if err != nil {
			return VoxelObject{}, fmt.Errorf("error reading chunk header: %v", err)
		}

		if chunkType == "" {
			break
		}

		switch chunkType {
		case "SIZE":
			// We only expect one SIZE chunk, but use the last value
			size, err = getSizeFromChunk(handle)
			if err != nil {
				return VoxelObject{}, fmt.Errorf("error reading size chunk: %v", err)
			}
			sizes = append(sizes, size)
		case "XYZI":
			data, err := getPointDataFromChunk(handle)
			if err != nil {
				return VoxelObject{}, fmt.Errorf("error reading size chunk: %v", err)
			}

			allPointData = append(allPointData, data)
		case "RGBA":
			paletteData, err = getPaletteDataFromChunk(handle)
			if err != nil {
				return VoxelObject{}, fmt.Errorf("Error reading palette chunk: %v", err)
			}
		case "nTRN":
			translation, err := getTranslationDataFromChunk(handle)
			if err != nil {
				return VoxelObject{}, fmt.Errorf("error reading nTRN chunk: %v", err)
			}
			scenegraph[translation.NodeID] = &translation
		case "nGRP":
			group, err := getGroupDataFromChunk(handle)
			if err != nil {
				return VoxelObject{}, fmt.Errorf("error reading nGRP chunk: %v", err)
			}
			scenegraph[group.NodeID] = &group
		case "nSHP":
			shape, err := getShapeDataFromChunk(handle)
			if err != nil {
				return VoxelObject{}, fmt.Errorf("error reading nSHP chunk: %v", err)
			}
			scenegraph[shape.NodeID] = &shape
		default:
			skipUnhandledChunk(handle)
		}
	}

	if size.X == 0 || size.Y == 0 || size.Z == 0 {
		return VoxelObject{}, fmt.Errorf("invalid size %v", size)
	}

	pointData := make([]geometry.PointWithColour, 0)
	rootNode, ok := scenegraph[0]
	if ok {
		TraverseScenegraph(scenegraph, rootNode, 0, 0, 0, &pointData, allPointData, sizes, &size)
		RebaseShape(pointData, &size)
	} else {
		for _, pd := range allPointData {
			pointData = append(pointData, pd...)
		}
	}

	object := VoxelObject{}
	object.Voxels = getVoxelObjectFromPointData(size, pointData)
	object.PaletteData = paletteData
	object.Size = size
	return object, nil
}

func RebaseShape(pointData []geometry.PointWithColour, size *geometry.Point) {
	minX, minY, minZ := 0, 0, 0

	// Establish the minimum point of a voxel
	for _, pd := range pointData {
		minX = min(minX, pd.Point.X)
		minY = min(minY, pd.Point.Y)
		minZ = min(minZ, pd.Point.Z)
	}

	// Offset if those are greater than zero
	if minX < 0 || minY < 0 || minZ < 0 {
		for idx, _ := range pointData {
			pointData[idx].Point.X -= minX
			pointData[idx].Point.Y -= minY
			pointData[idx].Point.Z -= minZ
		}
	}

	size.X = max(size.X, size.X - minX)
	size.Y = max(size.Y, size.Y - minY)
	size.Z = max(size.Z, size.Z - minZ)
}

func min(a int, b int) int {
	if a < b {
		return a
	}

	return b
}

func max(a int, b int) int {
	if a > b {
		return a
	}

	return b
}

func TraverseScenegraph(
	graph map[int]types.SceneGraphItem,
	current types.SceneGraphItem,
	x, y, z int,
	pointData *[]geometry.PointWithColour,
	allPointData [][]geometry.PointWithColour,
	sizes []geometry.Point,
	size *geometry.Point,
) {
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
		for _, child := range shp.Models {
			size.X = max(size.X, x + sizes[child].X)
			size.Y = max(size.Y, y + sizes[child].Y)
			size.Z = max(size.Z, z + sizes[child].Z)

			for _, pd := range allPointData[child] {
				point := geometry.PointWithColour{
					Point:  geometry.Point{
						X: pd.Point.X + x,
						Y: pd.Point.Y + y,
						Z: pd.Point.Z + z,
					},
					Colour: pd.Colour,
				}

				*pointData = append(*pointData, point)
			}
		}
	}

	for _, child := range current.GetChildren() {
		next, ok := graph[child]; if ok {
			TraverseScenegraph(graph, next, x, y, z, pointData, allPointData, sizes, size)
		}
	}
}

func GetFromReader(handle io.Reader) (v VoxelObject, err error) {
	v, err = GetMagicaVoxelObject(handle)
	return
}

func FromFile(filename string) (v VoxelObject, err error) {
	handle, err := os.Open(filename)
	if err != nil {
		return VoxelObject{}, err
	}

	v, err = GetFromReader(handle)
	if err != nil {
		return v, err
	}

	if err := handle.Close(); err != nil {
		return v, err
	}

	return v, nil
}
