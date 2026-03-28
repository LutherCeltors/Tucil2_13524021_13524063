package main

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: voxelizer <input.obj> <max_depth>")
		fmt.Println("  input.obj  : path to the input .obj file")
		fmt.Println("  max_depth  : maximum octree depth (positive integer)")
		os.Exit(1)
	}

	inputPath := os.Args[1]
	depthStr := os.Args[2]

	maxDepth, err := strconv.Atoi(depthStr)
	if err != nil || maxDepth <= 0 {
		fmt.Fprintf(os.Stderr, "Error: max_depth must be a positive integer, got: %s\n", depthStr)
		os.Exit(1)
	}

	fmt.Printf("Loading model: %s\n", inputPath)
	model, err := ParseOBJ(inputPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading OBJ: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Loaded %d vertices and %d faces.\n", len(model.Vertices), len(model.Faces))

	fmt.Printf("Voxelizing with max depth %d...\n", maxDepth)
	startTime := time.Now()

	result := Voxelize(model, maxDepth)

	elapsed := time.Since(startTime)

	outputPath := DeriveOutputPath(inputPath)

	err = WriteOBJ(outputPath, result.Vertices, result.Faces)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error writing output OBJ: %v\n", err)
		os.Exit(1)
	}

	fmt.Println()
	fmt.Println("=== Voxelization Report ===")
	fmt.Printf("Voxel count        : %d\n", result.VoxelCount)
	fmt.Printf("Vertex count       : %d\n", len(result.Vertices))
	fmt.Printf("Face count         : %d\n", len(result.Faces))
	fmt.Printf("Octree depth       : %d\n", maxDepth)
	fmt.Println()

	fmt.Println("Octree nodes per depth:")
	for d := 1; d <= maxDepth; d++ {
		fmt.Printf("  %d : %d\n", d, result.NodeCountPerDepth[d])
	}
	fmt.Println()

	fmt.Println("Pruned (skipped) nodes per depth:")
	for d := 1; d <= maxDepth; d++ {
		fmt.Printf("  %d : %d\n", d, result.PrunedPerDepth[d])
	}
	fmt.Println()

	fmt.Printf("Elapsed time       : %v\n", elapsed)
	fmt.Printf("Output saved to    : %s\n", outputPath)

	var ans string
	fmt.Printf("View the result (y/n) ? ")
	fmt.Scanln(&ans)
	if ans == "y" {
		PrintViewerControlTable()
		fmt.Printf("Tekan enter untuk melanjutkan...")
		fmt.Scanln(&ans)
		Render(inputPath, outputPath)
	}
}
