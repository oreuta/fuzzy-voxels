package voxel

// BuildCone ...
func BuildCone(r, h float32) IndicatorFunc {
	return func(x, y, z float32) bool {
		return Cone(x, y, z, 0, 0, 0, r, h)
	}
}

//Cone ...
//Center is in the point (x0,y0,z0)
func Cone(x, y, z, x0, y0, z0, r, h float32) bool {
	xx, yy, zz := x-x0, y-y0, z-z0
	zrh := zz * r / h

	return xx*xx+yy*yy <= zrh*zrh && zz <= h/2 && zz >= -h/2
}

func BuildConeVM(R, H float32, N uint32) BlenderVoxelFormat {
	const (
		zero = uint32(0)
	)

	s := BuildCone(R, H)

	b := BlenderVoxelFormat{}

	b.Nf = 1
	b.Nx, b.Ny, b.Nz = N, N, N
	b.V = make([][][][]float32, b.Nf)
	for f := zero; f < b.Nf; f++ {
		b.V[f] = make([][][]float32, b.Nz)
		for z := zero; z < b.Nz; z++ {
			b.V[f][z] = make([][]float32, b.Ny)
			for y := zero; y < b.Ny; y++ {
				b.V[f][z][y] = make([]float32, b.Nx)
			}
		}
	}

	L := float32(2 * R)
	Lx, Ly, Lz := L, L, L
	dx, dy, dz := L/float32(b.Nx), L/float32(b.Ny), L/float32(b.Nz)
	nx, ny, nz := int32(b.Nz), int32(b.Ny), int32(b.Nx)
	for k := int32(0); k < nz; k++ {
		for j := int32(0); j < ny; j++ {
			for i := int32(0); i < nx; i++ {
				xx := float32(i-(nx-1)/2) * Lx / float32(nx)
				yy := float32(j-(ny-1)/2) * Ly / float32(ny)
				zz := float32(k-(nz-1)/2) * Lz / float32(nz)
				b.V[zero][k][j][i] = Ro(xx, yy, zz, dx, dy, dz, s)
			}
		}
	}

	return b
}
