package main

type VoxelResult struct {
	Vertices          []Vec3
	Faces             [][4]int
	VoxelCount        int
	NodeCountPerDepth []int
	PrunedPerDepth    []int
}

type OctreeNode struct {
	Box      AABB
	Children [8]*OctreeNode
	IsLeaf   bool
	HasVoxel bool
}

func Voxelize(model *Model, maxDepth int) *VoxelResult {
	res := &VoxelResult{
		NodeCountPerDepth: make([]int, maxDepth+1),
		PrunedPerDepth:    make([]int, maxDepth+1),
	}

	rootBox := ComputeBounds(model.Vertices)
	tris := make([]rawTriangle, len(model.Faces))
	for i, f := range model.Faces {
		tris[i] = rawTriangle{
			V0: model.Vertices[f.A],
			V1: model.Vertices[f.B],
			V2: model.Vertices[f.C],
		}
	}
	var activeLeaves []AABB
	buildOctree(rootBox, tris, 1, maxDepth, res, &activeLeaves)

	res.VoxelCount = len(activeLeaves)
	for _, box := range activeLeaves {
		appendVoxelGeometry(box, res)
	}

	return res
}

type rawTriangle struct {
	V0, V1, V2 Vec3
}

func buildOctree(
	box AABB,
	tris []rawTriangle,
	depth, maxDepth int,
	res *VoxelResult,
	leaves *[]AABB,
) {
	res.NodeCountPerDepth[depth]++
	overlapping := filterTriangles(tris, box)

	if len(overlapping) == 0 {
		res.PrunedPerDepth[depth]++
		return
	}

	if depth == maxDepth {
		*leaves = append(*leaves, box)
		return
	}

	octants := subdivide(box)
	for _, child := range octants {
		buildOctree(child, overlapping, depth+1, maxDepth, res, leaves)
	}
}

func filterTriangles(tris []rawTriangle, box AABB) []rawTriangle {
	out := make([]rawTriangle, 0, len(tris))
	for _, t := range tris {
		if TriangleIntersectsAABB(t.V0, t.V1, t.V2, box) {
			out = append(out, t)
		}
	}
	return out
}

func subdivide(box AABB) [8]AABB {
	c := box.Center()
	return [8]AABB{
		{Min: Vec3{box.Min.X, box.Min.Y, box.Min.Z}, Max: Vec3{c.X, c.Y, c.Z}},
		{Min: Vec3{c.X, box.Min.Y, box.Min.Z}, Max: Vec3{box.Max.X, c.Y, c.Z}},
		{Min: Vec3{box.Min.X, box.Min.Y, c.Z}, Max: Vec3{c.X, c.Y, box.Max.Z}},
		{Min: Vec3{c.X, box.Min.Y, c.Z}, Max: Vec3{box.Max.X, c.Y, box.Max.Z}},
		{Min: Vec3{box.Min.X, c.Y, box.Min.Z}, Max: Vec3{c.X, box.Max.Y, c.Z}},
		{Min: Vec3{c.X, c.Y, box.Min.Z}, Max: Vec3{box.Max.X, box.Max.Y, c.Z}},
		{Min: Vec3{box.Min.X, c.Y, c.Z}, Max: Vec3{c.X, box.Max.Y, box.Max.Z}},
		{Min: Vec3{c.X, c.Y, c.Z}, Max: Vec3{box.Max.X, box.Max.Y, box.Max.Z}},
	}
}

func appendVoxelGeometry(box AABB, res *VoxelResult) {
	base := len(res.Vertices)
	mn := box.Min
	mx := box.Max

	res.Vertices = append(res.Vertices,
		Vec3{mn.X, mn.Y, mn.Z},
		Vec3{mx.X, mn.Y, mn.Z},
		Vec3{mx.X, mn.Y, mx.Z},
		Vec3{mn.X, mn.Y, mx.Z},
		Vec3{mn.X, mx.Y, mn.Z},
		Vec3{mx.X, mx.Y, mn.Z},
		Vec3{mx.X, mx.Y, mx.Z},
		Vec3{mn.X, mx.Y, mx.Z},
	)

	b := base
	res.Faces = append(res.Faces,
		[4]int{b + 0, b + 1, b + 2, b + 3},
		[4]int{b + 7, b + 6, b + 5, b + 4},
		[4]int{b + 0, b + 4, b + 5, b + 1},
		[4]int{b + 2, b + 6, b + 7, b + 3},
		[4]int{b + 0, b + 3, b + 7, b + 4},
		[4]int{b + 1, b + 5, b + 6, b + 2},
	)
}
