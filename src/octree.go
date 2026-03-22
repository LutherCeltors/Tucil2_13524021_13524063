package main

import ("sync")

type VoxelResult struct {
	Vertices          []Vec3
	Faces             [][4]int 
	VoxelCount        int
	NodeCountPerDepth []int 
	PrunedPerDepth    []int 
}
type rawTriangle struct {
	V0, V1, V2 Vec3
}
// Buat ngitung statsny
type concurrentCollector struct {
	mu     sync.Mutex
	leaves []AABB
	nodes  []int
	pruned []int 
}

func newCollector(maxDepth int) *concurrentCollector {
	return &concurrentCollector{
		nodes:  make([]int, maxDepth+1),
		pruned: make([]int, maxDepth+1),
	}
}
func (c *concurrentCollector) addLeaf(box AABB) {
	c.mu.Lock()
	c.leaves = append(c.leaves, box)
	c.mu.Unlock()
}
func (c *concurrentCollector) addNode(depth int) {
	c.mu.Lock()
	c.nodes[depth]++
	c.mu.Unlock()
}
func (c *concurrentCollector) addPruned(depth int) {
	c.mu.Lock()
	c.pruned[depth]++
	c.mu.Unlock()
}

// Dibawah ini masih pake rekursi normal, minimalisir overhead
const LimitDepth = 4

// Main
func Voxelize(model *Model, maxDepth int) *VoxelResult {
	rootBox := ComputeBounds(model.Vertices)
	tris := make([]rawTriangle, len(model.Faces))
	for i, f := range model.Faces {
		tris[i] = rawTriangle{
			V0: model.Vertices[f.A],
			V1: model.Vertices[f.B],
			V2: model.Vertices[f.C],
		}
	}
	collector := newCollector(maxDepth)
	// Cap semaphore goroutine (Buat ngontrol Thrashing)
	sem := make(chan struct{}, 512)
	var wg sync.WaitGroup
	wg.Add(1)
	go buildOctreeConcurrent(rootBox, tris, 1, maxDepth, collector, sem, &wg)
	wg.Wait()
	res := &VoxelResult{
		NodeCountPerDepth: collector.nodes,
		PrunedPerDepth:    collector.pruned,
	}
	res.VoxelCount = len(collector.leaves)
	for _, box := range collector.leaves {
		appendVoxelGeometry(box, res)
	}
	return res
}

func buildOctreeConcurrent(
	box AABB,
	tris []rawTriangle,
	depth, maxDepth int,
	col *concurrentCollector,
	sem chan struct{},
	wg *sync.WaitGroup,
) {
	defer wg.Done()
	col.addNode(depth)
	overlapping := filterTriangles(tris, box)
	if len(overlapping) == 0 {
		col.addPruned(depth)
		return
	}
	if depth == maxDepth {
		col.addLeaf(box)
		return
	}
	octants := subdivide(box)
	if depth <= LimitDepth {
		var childWg sync.WaitGroup
		for _, child := range octants {
			child := child          // capture
			sem <- struct{}{}       // ambil semaphore slot
			childWg.Add(1)
			wg.Add(1)
			go func() {
				defer func() { <-sem }() // lepas slot
				buildOctreeConcurrent(child, overlapping, depth+1, maxDepth, col, sem, wg)
				childWg.Done()
			}()
		}
		childWg.Wait()
	} else {
		for _, child := range octants {
			wg.Add(1)
			buildOctreeConcurrent(child, overlapping, depth+1, maxDepth, col, sem, wg)
		}
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