package magica

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"github.com/mattkimber/gandalf/magica/chunk"
	"github.com/mattkimber/gandalf/magica/scenegraph"
	"github.com/mattkimber/gandalf/magica/types"
	"io"
	"os"
	"sort"
)

const maxVoxelSize = 256

func(v *VoxelObject) GetData() []byte {
	var graph scenegraph.Map
	var pointData []types.PointData
	var sizeData []types.Size

	if v.Size.X < maxVoxelSize && v.Size.Y < maxVoxelSize && v.Size.Z < maxVoxelSize {
		// Simple case, fits in a single Magica object
		sizeData = []types.Size{{
			X: v.Size.X,
			Y: v.Size.Y,
			Z: v.Size.Z,
		}}

		pointData = []types.PointData{v.GetPoints()}
	} else {
		// Complex case, object needs to be split
		scenegraphNode := v.Split(maxVoxelSize)
		graph, pointData, sizeData = scenegraphNode.Decompose()
	}

	chunks := make([]chunk.Chunk, 0)

	for idx, _ := range pointData {
		size := chunk.Chunk{Item: &sizeData[idx]}
		chunks = append(chunks, size)
		xyzi := chunk.Chunk{Item: &pointData[idx]}
		chunks = append(chunks, xyzi)
	}

	graphKeys := make([]int, 0, len(graph))
	for k, _ := range graph {
		graphKeys = append(graphKeys, k)
	}
	sort.Ints(graphKeys)

	for _, k := range graphKeys {
		mapItem, _ := graph[k]

		mi := chunk.Chunk{Item: mapItem}
		chunks = append(chunks, mi)
	}

	pal := types.Palette(v.PaletteData)
	chunks = append(chunks, chunk.Chunk{Item: &pal})

	main := chunk.MainChunk{}

	mainChunk := chunk.Chunk{
		Item:     &main,
		Children: chunks,
	}

	return mainChunk.GetBytes()
}

func (v *VoxelObject) writeHeader(handle io.Writer) (err error) {
	if _, err = handle.Write([]byte("VOX ")); err != nil {
		return err
	}

	err = binary.Write(handle, binary.LittleEndian, int32(150))
	return
}

func (v *VoxelObject) Save(handle io.Writer) (err error) {
	bw := bufio.NewWriter(handle)
	v.writeHeader(bw)
	data := v.GetData()
	bw.Write(data)
	bw.Flush()
	return
}

// SaveToFile saves the voxel object to the specified file
func (v *VoxelObject) SaveToFile(filename string) (err error) {
	handle, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("could not create output file %s: %v", filename, err)
	}

	err = v.Save(handle)
	if err != nil {
		handle.Close()
		return fmt.Errorf("could not open output file: %v", err)
	}

	err = handle.Close()
	return err
}