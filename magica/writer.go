package magica

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"github.com/mattkimber/gandalf/magica/chunk"
	"github.com/mattkimber/gandalf/magica/types"
	"io"
	"os"
)

func(v *VoxelObject) GetData() []byte {
	if v.Size.X < 256 && v.Size.Y < 256 && v.Size.Z < 256 {
		// Simple case, fits in a single Magica object
		main := chunk.MainChunk{}
		size := chunk.Chunk{Item: &types.Size{
			X: v.Size.X,
			Y: v.Size.Y,
			Z: v.Size.Z,
		}}

		pd := v.GetPoints()
		xyzi := chunk.Chunk{Item: &pd}

		pal := types.Palette(v.PaletteData)
		rgba := chunk.Chunk{Item: &pal}

		chunks := chunk.Chunk{
			Item:     &main,
			Children: []chunk.Chunk{size,xyzi,rgba},
		}

		return chunks.GetBytes()
	}

	// Larger objects not yet implemented
	return nil
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