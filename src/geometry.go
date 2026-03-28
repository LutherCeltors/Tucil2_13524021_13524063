package main

import "math"

type AABB struct {
	Min, Max Vec3
}

func (a AABB) Center() Vec3 {
	return Vec3{
		(a.Min.X + a.Max.X) * 0.5,
		(a.Min.Y + a.Max.Y) * 0.5,
		(a.Min.Z + a.Max.Z) * 0.5,
	}
}

func (a AABB) HalfSize() Vec3 {
	return Vec3{
		(a.Max.X - a.Min.X) * 0.5,
		(a.Max.Y - a.Min.Y) * 0.5,
		(a.Max.Z - a.Min.Z) * 0.5,
	}
}

func (a AABB) ContainsPoint(p Vec3) bool {
	return p.X >= a.Min.X && p.X <= a.Max.X &&
		p.Y >= a.Min.Y && p.Y <= a.Max.Y &&
		p.Z >= a.Min.Z && p.Z <= a.Max.Z
}

func ComputeBounds(verts []Vec3) AABB {
	const eps = 1e-6
	b := AABB{
		Min: Vec3{math.MaxFloat64, math.MaxFloat64, math.MaxFloat64},
		Max: Vec3{-math.MaxFloat64, -math.MaxFloat64, -math.MaxFloat64},
	}
	for _, v := range verts {
		if v.X < b.Min.X { b.Min.X = v.X }
		if v.Y < b.Min.Y { b.Min.Y = v.Y }
		if v.Z < b.Min.Z { b.Min.Z = v.Z }
		if v.X > b.Max.X { b.Max.X = v.X }
		if v.Y > b.Max.Y { b.Max.Y = v.Y }
		if v.Z > b.Max.Z { b.Max.Z = v.Z }
	}
	dx := b.Max.X - b.Min.X
	dy := b.Max.Y - b.Min.Y
	dz := b.Max.Z - b.Min.Z
	maxD := math.Max(dx, math.Max(dy, dz))
	cx := (b.Min.X + b.Max.X) * 0.5
	cy := (b.Min.Y + b.Max.Y) * 0.5
	cz := (b.Min.Z + b.Max.Z) * 0.5
	half := maxD*0.5 + eps
	return AABB{
		Min: Vec3{cx - half, cy - half, cz - half},
		Max: Vec3{cx + half, cy + half, cz + half},
	}
}

func TriangleIntersectsAABB(v0, v1, v2 Vec3, box AABB) bool {
	c := box.Center()
	h := box.HalfSize()

	a0 := Vec3{v0.X - c.X, v0.Y - c.Y, v0.Z - c.Z}
	a1 := Vec3{v1.X - c.X, v1.Y - c.Y, v1.Z - c.Z}
	a2 := Vec3{v2.X - c.X, v2.Y - c.Y, v2.Z - c.Z}

	minX := math.Min(a0.X, math.Min(a1.X, a2.X))
	maxX := math.Max(a0.X, math.Max(a1.X, a2.X))
	if minX > h.X || maxX < -h.X { return false }

	minY := math.Min(a0.Y, math.Min(a1.Y, a2.Y))
	maxY := math.Max(a0.Y, math.Max(a1.Y, a2.Y))
	if minY > h.Y || maxY < -h.Y { return false }

	minZ := math.Min(a0.Z, math.Min(a1.Z, a2.Z))
	maxZ := math.Max(a0.Z, math.Max(a1.Z, a2.Z))
	if minZ > h.Z || maxZ < -h.Z { return false }

	e0 := Vec3{a1.X - a0.X, a1.Y - a0.Y, a1.Z - a0.Z}
	e1 := Vec3{a2.X - a1.X, a2.Y - a1.Y, a2.Z - a1.Z}
	normal := cross(e0, e1)
	if !planeAABBOverlap(normal, a0, h) { return false }

	edges := [3]Vec3{
		e0,
		e1,
		{a0.X - a2.X, a0.Y - a2.Y, a0.Z - a2.Z},
	}
	verts := [3]Vec3{a0, a1, a2}
	axes := [3]Vec3{{1, 0, 0}, {0, 1, 0}, {0, 0, 1}}

	for _, e := range edges {
		for _, ax := range axes {
			sep := cross(e, ax)
			p0 := dot(sep, verts[0])
			p1 := dot(sep, verts[1])
			p2 := dot(sep, verts[2])
			pMin := math.Min(p0, math.Min(p1, p2))
			pMax := math.Max(p0, math.Max(p1, p2))
			r := h.X*math.Abs(sep.X) + h.Y*math.Abs(sep.Y) + h.Z*math.Abs(sep.Z)
			if pMin > r || pMax < -r {
				return false
			}
		}
	}
	return true
}

func cross(a, b Vec3) Vec3 {
	return Vec3{
		a.Y*b.Z - a.Z*b.Y,
		a.Z*b.X - a.X*b.Z,
		a.X*b.Y - a.Y*b.X,
	}
}

func dot(a, b Vec3) float64 {
	return a.X*b.X + a.Y*b.Y + a.Z*b.Z
}

func planeAABBOverlap(normal, point Vec3, h Vec3) bool {
	r := h.X*math.Abs(normal.X) + h.Y*math.Abs(normal.Y) + h.Z*math.Abs(normal.Z)
	s := dot(normal, point)
	return math.Abs(s) <= r
}
