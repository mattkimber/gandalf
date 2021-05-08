package magica

import (
	"encoding/binary"
	"fmt"
	"github.com/mattkimber/gandalf/geometry"
	"github.com/mattkimber/gandalf/magica/scenegraph"
	"github.com/mattkimber/gandalf/magica/types"
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


func GetMagicaVoxelObject(handle io.Reader, layers []int) (VoxelObject, error) {
	if !isHeaderValid(handle) {
		return VoxelObject{}, fmt.Errorf("header not valid")
	}
	getChunkHeader(handle)

	sizeData := make([]types.Size, 0)
	pointData := make([]types.PointData, 0)

	scenegraphMap := make(scenegraph.Map)

	var palette types.Palette

	for {
		chunkType, err := getChunkHeader(handle)

		if err != nil {
			return VoxelObject{}, fmt.Errorf("error reading chunk header: %v", err)
		}

		if chunkType == "" {
			break
		}

		data, err := getChunkData(handle, 0)

		if err != nil {
			return VoxelObject{}, err
		}

		rd := types.GetReader(data)

		switch chunkType {
		case "SIZE":
			data := rd.GetSize()
			sizeData = append(sizeData, data)
		case "XYZI":
			data := rd.GetPointData()
			pointData = append(pointData, data)
		case "RGBA":
			palette = rd.GetPalette()
		case "nTRN":
			translation := rd.GetTranslation()
			scenegraphMap[translation.NodeID] = &translation
		case "nGRP":
			group := rd.GetGroup()
			scenegraphMap[group.NodeID] = &group
		case "nSHP":
			shape := rd.GetShape()
			scenegraphMap[shape.NodeID] = &shape
		default:
		}
	}


	graph := scenegraph.GetScenegraph(scenegraphMap, layers, pointData, sizeData)
	model := graph.GetCompositeModel()


	object := VoxelObject{}
	object.PaletteData = palette
	object.Size = geometry.Point{X: model.Size.X, Y: model.Size.Y, Z: model.Size.Z}
	object.Voxels = model.Data
	return object, nil
}


func GetFromReader(handle io.Reader, layers []int) (v VoxelObject, err error) {
	v, err = GetMagicaVoxelObject(handle, layers)
	return
}

func FromFileWithLayers(filename string, layers []int) (v VoxelObject, err error) {
	handle, err := os.Open(filename)
	if err != nil {
		return VoxelObject{}, err
	}

	v, err = GetFromReader(handle, layers)
	if err != nil {
		return v, err
	}

	if err := handle.Close(); err != nil {
		return v, err
	}

	return v, nil
}

func FromFile(filename string) (v VoxelObject, err error) {
	return FromFileWithLayers(filename, []int{})
}

