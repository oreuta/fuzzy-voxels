package voxel

// BuildUnion returns a union of two models a and b translated to (dx,dy,dz)
func BuildUnion(a, b BlenderVoxelFormat, dx, dy, dz uint32) (g BlenderVoxelFormat) {
	return buildResultVM(a, b, dx, dy, dz, maxFloat32)
}

// BuildIntersection returns an intersection of two models a and b translated to (dx,dy,dz)
func BuildIntersection(a, b BlenderVoxelFormat, dx, dy, dz uint32) (g BlenderVoxelFormat) {
	return buildResultVM(a, b, dx, dy, dz, minFloat32)
}

// BuildSubstraction1minus2 returns a / b translated to (dx,dy,dz)
func BuildSubstraction1minus2(a, b BlenderVoxelFormat, dx, dy, dz uint32) (g BlenderVoxelFormat) {
	return buildResultVM(a, b, dx, dy, dz, sub12Float32)
}

// BuildSubstraction2minus1 returns b / a translated to (dx,dy,dz)
func BuildSubstraction2minus1(a, b BlenderVoxelFormat, dx, dy, dz uint32) (g BlenderVoxelFormat) {
	return buildResultVM(a, b, dx, dy, dz, sub21Float32)
}

type operation func(a, b float32) float32

func buildResultVM(a, b BlenderVoxelFormat, dx, dy, dz uint32, op operation) (g BlenderVoxelFormat) {
	nf := uint32(1)
	nx := dx + maxUint32(a.Nx-dx, b.Nx)
	ny := dy + maxUint32(a.Ny-dy, b.Ny)
	nz := dz + maxUint32(a.Nz-dz, b.Nz)
	g.Init(nf, nx, ny, nz)

	ro := buildRoFunc(dx, dy, dz, nx, ny, nz, op)

	var (
		roA, roB   float32
		xb, yb, zb uint32
	)
	for x := zeroUint32; x < nx; x++ {
		for y := zeroUint32; y < ny; y++ {
			for z := zeroUint32; z < nz; z++ {
				roA, roB = zeroFloat32, zeroFloat32
				if x < a.Nx && y < a.Ny && z < a.Nz {
					roA = a.V[nf-1][z][y][x]
				}
				xb = x - dx
				yb = y - dy
				zb = z - dz
				if xb < b.Nx && yb < b.Ny && zb < b.Nz {
					roB = b.V[nf-1][zb][yb][xb]
				}
				g.V[nf-1][x][y][z] = ro(x, y, z, roA, roB)
			}
		}
	}

	return
}

func maxFloat32(a, b float32) float32 {
	if a > b {
		return a
	}
	return b
}

func minFloat32(a, b float32) float32 {
	if a > b {
		return b
	}
	return a
}

func sub12Float32(a, b float32) float32 {
	if a > b {
		return a - b
	}
	return zeroFloat32
}

func sub21Float32(a, b float32) float32 {
	if b > a {
		return b - a
	}
	return zeroFloat32
}

func maxUint32(a, b uint32) uint32 {
	if a > b {
		return a
	}
	return b
}

func indexesToRo(condX, condY, condZ int, a, b, ab float32) float32 {
	index := condX*9 + condY*3 + condZ
	result := [27]float32{
		a, a, 0, a, a, 0, 0, 0, 0, a, a, 0, a, ab, b, 0, b, b, 0, 0, 0, 0, b, b, 0, b, b,
	}
	return result[index]
}

func cond(q, dq, nq uint32) int {
	switch {
	case 0 <= q && q < dq:
		return 0
	case dq <= q && q < nq:
		return 1
	case nq <= q && q < nq+dq:
		return 2
	}
	return -1
}

func buildRoFunc(dx, dy, dz, nx, ny, nz uint32, op operation) func(x, y, z uint32, a, b float32) float32 {
	const zero = float32(0)
	return func(x, y, z uint32, a, b float32) float32 {
		condX := cond(x, dx, nx)
		condY := cond(y, dy, ny)
		condZ := cond(z, dz, nz)

		roA := op(a, zero)
		roB := op(b, zero)
		roAB := op(a, b)
		return indexesToRo(condX, condY, condZ, roA, roB, roAB)
	}
}
