package main

import (
	"fmt"

	"github.com/oreuta/fuzzy-voxels/voxel"
)

func main() {
	var err error

	R := float32(5)
	//H := float32(10) //cylinder's height
	N := uint32(21)
	//vm := voxel.BuildCylinderVM(R, H, N)
	//vm := voxel.BuildConeVM(R, H, N)
	//vm := voxel.BuildSpheraVM(R, N)
	a := voxel.BuildSpheraVM(R, N)
	//b := voxel.BuildSpheraVM(R, N)
	//delta := uint32(5)
	//g := voxel.BuildUnion(a, b, delta, delta, delta)
	//g := voxel.BuildIntersection(a, b, delta, delta, delta)
	//g := voxel.BuildSubstraction2minus1(a, b, delta, delta, delta)
	//err = a.WritePython("C:\\tmp\\sub10-5-2.py")
	//err = a.WriteToFile("C:\\tmp\\sub10-5-2.blender")
	err = a.WritePython("C:\\tmp\\sphera21-5.py")
	fmt.Printf("Finish. Err = %v", err)
}
