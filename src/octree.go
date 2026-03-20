package main

// VoxelResult holds everything produced by the voxelization pass.
type VoxelResult struct {
	Vertices          []Vec3
	Faces             [][4]int // quad faces (4 vertex indices, 0-based)
	VoxelCount        int
	NodeCountPerDepth []int // index 0 unused; index d = node count at depth d
	PrunedPerDepth    []int // same layout; nodes pruned (empty/fully outside) per depth
}

// OctreeNode represents a single node in the octree.
type OctreeNode struct {
	Box      AABB
	Children [8]*OctreeNode
	IsLeaf   bool
	HasVoxel bool // true when this leaf should generate a voxel cube
}

// Voxelize is the public entry-point. It builds an octree over the model and
// collects leaf voxels using a Divide-and-Conquer strategy.
func Voxelize(model *Model, maxDepth int) *VoxelResult {
	res := &VoxelResult{
		NodeCountPerDepth: make([]int, maxDepth+1),
		PrunedPerDepth:    make([]int, maxDepth+1),
	}

	// Pre-compute the root bounding box (a cube)
	rootBox := ComputeBounds(model.Vertices)

	// Build a list of raw triangles for fast intersection queries
	tris := make([]rawTriangle, len(model.Faces))
	for i, f := range model.Faces {
		tris[i] = rawTriangle{
			V0: model.Vertices[f.A],
			V1: model.Vertices[f.B],
			V2: model.Vertices[f.C],
		}
	}

	// Collect active leaf AABBs via divide-and-conquer
	var activeLeaves []AABB
	buildOctree(rootBox, tris, 1, maxDepth, res, &activeLeaves)

	// Generate OBJ geometry for each active leaf (a unit cube of 8 verts, 6 quad faces)
	res.VoxelCount = len(activeLeaves)
	for _, box := range activeLeaves {
		appendVoxelGeometry(box, res)
	}

	return res
}

// rawTriangle stores the three pre-resolved vertex positions of a face.
type rawTriangle struct {
	V0, V1, V2 Vec3
}

// buildOctree recursively subdivides `box` using divide-and-conquer.
//
//   - If no triangle from `tris` intersects this box → prune (empty node).
//   - If we are at max depth → this is an active leaf (becomes a voxel).
//   - Otherwise → subdivide into 8 children (octants) and recurse.
func buildOctree(
	box AABB,
	tris []rawTriangle,
	depth, maxDepth int,
	res *VoxelResult,
	leaves *[]AABB,
) {
	res.NodeCountPerDepth[depth]++

	// Intersection test: keep only triangles that overlap this box.
	overlapping := filterTriangles(tris, box)

	if len(overlapping) == 0 {
		// No geometry in this region – prune.
		res.PrunedPerDepth[depth]++
		return
	}

	if depth == maxDepth {
		// Leaf node with geometry → record as a voxel.
		*leaves = append(*leaves, box)
		return
	}

	// Divide: split box into 8 octants and conquer each.
	octants := subdivide(box)
	for _, child := range octants {
		buildOctree(child, overlapping, depth+1, maxDepth, res, leaves)
	}
}

// filterTriangles returns the subset of tris that intersect box.
func filterTriangles(tris []rawTriangle, box AABB) []rawTriangle {
	out := make([]rawTriangle, 0, len(tris))
	for _, t := range tris {
		if TriangleIntersectsAABB(t.V0, t.V1, t.V2, box) {
			out = append(out, t)
		}
	}
	return out
}

// subdivide splits an AABB into its 8 octant children.
func subdivide(box AABB) [8]AABB {
	c := box.Center()
	return [8]AABB{
		// Bottom (Y-) half
		{Min: Vec3{box.Min.X, box.Min.Y, box.Min.Z}, Max: Vec3{c.X, c.Y, c.Z}}, // 0: -x -y -z
		{Min: Vec3{c.X, box.Min.Y, box.Min.Z}, Max: Vec3{box.Max.X, c.Y, c.Z}}, // 1: +x -y -z
		{Min: Vec3{box.Min.X, box.Min.Y, c.Z}, Max: Vec3{c.X, c.Y, box.Max.Z}}, // 2: -x -y +z
		{Min: Vec3{c.X, box.Min.Y, c.Z}, Max: Vec3{box.Max.X, c.Y, box.Max.Z}}, // 3: +x -y +z
		// Top (Y+) half
		{Min: Vec3{box.Min.X, c.Y, box.Min.Z}, Max: Vec3{c.X, box.Max.Y, c.Z}}, // 4: -x +y -z
		{Min: Vec3{c.X, c.Y, box.Min.Z}, Max: Vec3{box.Max.X, box.Max.Y, c.Z}}, // 5: +x +y -z
		{Min: Vec3{box.Min.X, c.Y, c.Z}, Max: Vec3{c.X, box.Max.Y, box.Max.Z}}, // 6: -x +y +z
		{Min: Vec3{c.X, c.Y, c.Z}, Max: Vec3{box.Max.X, box.Max.Y, box.Max.Z}}, // 7: +x +y +z
	}
}

// appendVoxelGeometry appends 8 vertices and 6 quad faces for one voxel cube.
//
// The cube is axis-aligned with corners at box.Min and box.Max.
// Vertices are ordered so that face normals point outward.
//
//	  7----6
//	 /|   /|
//	4----5 |
//	| 3--|-2
//	|/   |/
//	0----1
//
// Vertex layout (0=MinXMinYMinZ corner):
//
//	0: Min
//	1: MaxX MinY MinZ
//	2: MaxX MinY MaxZ
//	3: MinX MinY MaxZ
//	4: MinX MaxY MinZ
//	5: MaxX MaxY MinZ
//	6: MaxX MaxY MaxZ
//	7: MinX MaxY MaxZ
func appendVoxelGeometry(box AABB, res *VoxelResult) {
	base := len(res.Vertices)
	mn := box.Min
	mx := box.Max

	res.Vertices = append(res.Vertices,
		Vec3{mn.X, mn.Y, mn.Z}, // 0
		Vec3{mx.X, mn.Y, mn.Z}, // 1
		Vec3{mx.X, mn.Y, mx.Z}, // 2
		Vec3{mn.X, mn.Y, mx.Z}, // 3
		Vec3{mn.X, mx.Y, mn.Z}, // 4
		Vec3{mx.X, mx.Y, mn.Z}, // 5
		Vec3{mx.X, mx.Y, mx.Z}, // 6
		Vec3{mn.X, mx.Y, mx.Z}, // 7
	)

	b := base
	res.Faces = append(res.Faces,
		[4]int{b + 0, b + 1, b + 2, b + 3}, // Bottom  (-Y)
		[4]int{b + 7, b + 6, b + 5, b + 4}, // Top     (+Y)
		[4]int{b + 0, b + 4, b + 5, b + 1}, // Front   (-Z)
		[4]int{b + 2, b + 6, b + 7, b + 3}, // Back    (+Z)
		[4]int{b + 0, b + 3, b + 7, b + 4}, // Left    (-X)
		[4]int{b + 1, b + 5, b + 6, b + 2}, // Right   (+X)
	)
}
