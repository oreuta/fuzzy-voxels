package voxel

import "fmt"

const (
	zeroUint32  = uint32(0)
	zeroFloat32 = float32(0)
)

// SpaceObject represents 3D-object behaviour as a Voxel Model.
type SpaceObject interface {
	GetIndicatorFunc() IndicatorFunc
	GetVoxelModel() BlenderVoxelFormat
}

// IndicatorFunc describes shape of a 3D-object.
// Returns 'ture' if a point (x,y,z) belongs to the object, and 'false' otherwise.
type IndicatorFunc func(x, y, z float32) bool

// Ro is an optical density of a voxel of size (Lx,Ly,Lz) in (xc,yc,yz).
// The object is defined by its indicator function f.
func Ro(xc, yc, zc, Lx, Ly, Lz float32, f IndicatorFunc) float32 {
	n := 5
	nx, ny, nz := n, n, n
	count := 0
	for x := 0; x < nx; x++ {
		for y := 0; y < ny; y++ {
			for z := 0; z < nz; z++ {
				xx := xc + float32(x-(nx-1)/2)*Lx/float32(nx)
				yy := yc + float32(y-(ny-1)/2)*Ly/float32(ny)
				zz := zc + float32(z-(nz-1)/2)*Lz/float32(nz)
				if f(xx, yy, zz) {
					count++
				}
			}
		}
	}
	b := float32(count) / float32(nx*ny*nz)
	fmt.Printf("%.2f\n", b)
	return b
}
