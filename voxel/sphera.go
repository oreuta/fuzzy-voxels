package voxel

// BuildSphera ...
func BuildSphera(r float32) IndicatorFunc {
	return func(x, y, z float32) bool {
		return Sphera(x, y, z, 0, 0, 0, r)
	}
}

//Sphera ...
//Center is in the point (x0,y0,z0)
func Sphera(x, y, z, x0, y0, z0, r float32) bool {
	xx, yy, zz := x-x0, y-y0, z-z0
	return xx*xx+yy*yy+zz*zz <= r*r
}

func BuildSpheraVM(R float32, N uint32) BlenderVoxelFormat {
	const (
		zero = uint32(0)
	)

	s := BuildSphera(R)

	var b BlenderVoxelFormat
	b.Init(1, N, N, N)

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
