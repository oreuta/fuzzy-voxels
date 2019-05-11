package voxel

import (
	"encoding/binary"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"os"
)

// BlenderVoxelFormat represents a datacube.
// The order to write the full datacube after the header is
// 1. Frame by frame, where each frame has nz layers.
// 2. Layer by layer, where each layer has ny lines.
// 3. Line by line, where each line has nx values.
// 4. Value by value, until completing the nx values.
type BlenderVoxelFormat struct {
	Nx, Ny, Nz uint32          //the number of subdivisions of the domain in each axis
	Nf         uint32          //the total number of frames contained in the file
	V          [][][][]float32 //values, which describe data normalized to a [0, 1] interval
}

// Init creates an empty Voxel Model
func (f *BlenderVoxelFormat) Init(Nf, Nx, Ny, Nz uint32) {
	const zero = uint32(0)
	f.Nf = Nf
	f.Nx, f.Ny, f.Nz = Nx, Ny, Nz
	f.V = make([][][][]float32, f.Nf)
	for v := zero; v < f.Nf; v++ {
		f.V[v] = make([][][]float32, f.Nz)
		for z := zero; z < f.Nz; z++ {
			f.V[v][z] = make([][]float32, f.Ny)
			for y := zero; y < f.Ny; y++ {
				f.V[v][z][y] = make([]float32, f.Nx)
			}
		}
	}
}

// WriteToFile writes datacube to file in Blender Voxel format
func (f BlenderVoxelFormat) WriteToFile(fileName string) (err error) {
	file, err := os.Create(fileName)
	if err != nil {
		return
	}
	_, err = f.WriteTo(file)
	file.Close()
	return
}

// WriteTo writes datacube as a byte stream to the specified Writer
func (f BlenderVoxelFormat) WriteTo(w io.Writer) (n int64, err error) {
	var count int
	buf := make([]byte, 4)
	binary.LittleEndian.PutUint32(buf, f.Nx)
	count, err = w.Write(buf)
	if err != nil {
		n = int64(count)
		return
	}

	binary.LittleEndian.PutUint32(buf, f.Ny)
	count, err = w.Write(buf)
	if err != nil {
		n += int64(count)
		return
	}

	binary.LittleEndian.PutUint32(buf, f.Nz)
	count, err = w.Write(buf)
	if err != nil {
		n += int64(count)
		return
	}

	binary.LittleEndian.PutUint32(buf, f.Nf)
	count, err = w.Write(buf)
	if err != nil {
		n += int64(count)
		return
	}

	var q, x, y, z uint32
	for q = 0; q < f.Nf; q++ {
		for z = 0; z < f.Nz; z++ {
			for y = 0; y < f.Ny; y++ {
				for x = 0; x < f.Nx; x++ {
					binary.LittleEndian.PutUint32(buf, math.Float32bits(f.V[q][z][y][x]))
					count, err = w.Write(buf)
					if err != nil {
						n += int64(count)
						return
					}
				}
			}
		}
	}
	return
}

// WritePython ...
func (f BlenderVoxelFormat) WritePython(fileName string) error {
	var x, y, z uint32
	buf := []byte{}
	buf = append(buf, []byte(pythonHeader)...)
	buf = append(buf, []byte(fmt.Sprintf("VM = np.zeros( (%d, %d, %d), dtype=float)\n", f.Nx, f.Ny, f.Nz))...)
	for z = 0; z < f.Nz; z++ {
		for y = 0; y < f.Ny; y++ {
			for x = 0; x < f.Nx; x++ {
				s := fmt.Sprintf("VM[%d,%d,%d] = %.2f\n", x, y, z, f.V[uint32(0)][z][y][x])
				buf = append(buf, []byte(s)...)
			}
		}
	}
	buf = append(buf, []byte(fmt.Sprintf("N = %d\nM = %d\nK = %d\n", f.Nx, f.Ny, f.Nz))...)
	buf = append(buf, []byte("g = draw_voxel_model(VM, N, M, K)")...)

	return ioutil.WriteFile(fileName, buf, 0666)
}

const pythonHeader = `import numpy as np
import bpy

def draw_voxel_model(V, N, M, K, group_name='VM'):
  g = bpy.data.groups.new(group_name)
  mat_dict = dict()
  Nh = N/2
  Mh = M/2
  Kh = K/2
  for i in range(N):
    for j in range(M):
      for k in range(K):
        p = V[i,j,k]
        if p > 0:
          mat_name = 'm'+str(p)[2:]
          mat = bpy.data.materials.get(mat_name)
          if mat is None:  
            mat = bpy.data.materials.new(mat_name)
            mat.diffuse_color = (0.5,0.5,0.5)
            mat.alpha = p
            mat.use_transparency = True
          bpy.ops.mesh.primitive_cube_add(location=(i+1/2-Nh, j+1/2-Mh, k+1/2-Kh))
          v = bpy.context.active_object
          v.dimensions = (1,1,1)
          v.active_material = mat
          v.show_transparent = True
          g.objects.link(v)
  return g

`
