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

// Geometri untuk viewer
// Vec2, Vec3, Vec4, serta operasi matriks 4x4

type Mat4 [4][4] float64;

type Vec4 struct {
	W, X, Y, Z float64;
}

type ScreenVertex struct {
	X, Y, Z float64
}

func Add3(v1, v2 Vec3) Vec3 {
	return Vec3{v1.X + v2.X, v1.Y + v2.Y, v1.Z + v2.Z};
}

func Sub3(v1, v2 Vec3) Vec3 {
	return Vec3{v1.X - v2.X, v1.Y - v2.Y, v1.Z - v2.Z};
}

func MulN3(v1 Vec3, s float64) Vec3 {
	return Vec3{v1.X * s , v1.Y * s, v1.Z * s};
}

func Norm3(v1 Vec3) float64 {
	return math.Sqrt(math.Pow(v1.X, 2) + math.Pow(v1.Y, 2) + math.Pow(v1.Z, 2));
}

func Normalized3(v Vec3) Vec3 {
	if (Norm3(v) == 0) {
		return v
	}
	return Vec3{v.X / Norm3(v), v.Y / Norm3(v), v.Z / Norm3(v)};
}

func Mat4Identity() Mat4 {
	return Mat4 {
		{1, 0 , 0, 0},
		{0, 1 , 0, 0},
		{0, 0 , 1, 0},
		{0, 0 , 0, 1},
	}
}

func MulMat4Vec4 (m4 Mat4, v4 Vec4) Vec4 {
	var vRes Vec4
	var tempRes float64

	for i := range m4 {
		tempRes = 0.0
		tempRes += m4[i][0] * v4.X
		tempRes += m4[i][1] * v4.Y
		tempRes += m4[i][2] * v4.Z
		tempRes += m4[i][3] * v4.W
		
		switch i {
		case 0 : vRes.X = tempRes
		case 1 : vRes.Y = tempRes
		case 2 : vRes.Z = tempRes
		case 3 : vRes.W = tempRes
		}
	}

	return vRes
}

func MulMat4Mat4 (m0, m1 Mat4) Mat4 {
	var mRes Mat4
	var tempRes float64

	for r0 := 0; r0 < 4; r0++ {
		for c1 := 0; c1 < 4; c1++ {
			tempRes = 0.0
			for i := 0; i < 4; i++ {
				tempRes += m0[r0][i] * m1[i][c1]
			}
			mRes[r0][c1] = tempRes
		}
	}

	return mRes
}

func RotationX4 (rad float64) Mat4 {
	c := math.Cos(rad)
	s := math.Sin(rad)

	return Mat4 {
		{1, 0, 0, 0}, 
		{0, c, -s, 0},
		{0, s, c, 0},
		{0 ,0 ,0 ,1},
	}
}

func RotationY4(rad float64) Mat4 {
	c := math.Cos(rad)
	s := math.Sin(rad)

	return Mat4{
		{c, 0, s, 0},
		{0, 1, 0, 0},
		{-s, 0, c, 0},
		{0, 0, 0, 1},
	}
}

func RotationZ4(rad float64) Mat4 {
	c := math.Cos(rad)
	s := math.Sin(rad)

	return Mat4{
		{c, -s, 0, 0},
		{s, c, 0, 0},
		{0, 0, 1, 0},
		{0, 0, 0, 1},
	}
}

func RotationXYZ4(rX, rY, rZ float64) Mat4 {
	mX := RotationX4(rX);
	mY := RotationY4(rY);
	mZ := RotationZ4(rZ);
	return MulMat4Mat4(mZ, MulMat4Mat4(mY, mX))
}

func Perspective(opPresVAngle, windowRatio, near, far float64) Mat4 {
	f := 1.0 / math.Tan(opPresVAngle/2.0)

	return Mat4{
		{f / windowRatio, 0, 0, 0},
		{0, f, 0, 0},
		{0, 0, (far + near) / (near - far), (2 * far * near) / (near - far)},
		{0, 0, -1, 0},
	}
}

func LookAt(cam, target, up Vec3) Mat4 {
	forward := Normalized3(Sub3(target, cam))
	right := Normalized3(cross(forward, up))
	trueUp := cross(right, forward)

	return Mat4{
		{right.X, right.Y, right.Z, -dot(right, cam)},
		{trueUp.X, trueUp.Y, trueUp.Z, -dot(trueUp, cam)},
		{-forward.X, -forward.Y, -forward.Z, dot(forward, cam)},
		{0, 0, 0, 1},
	}
}

func ProjectVertex(v Vec3, mvp Mat4, width, height int) (ScreenVertex, bool) {
	v4 := Vec4{
		X: v.X,
		Y: v.Y,
		Z: v.Z,
		W: 1.0,
	}

	clip := MulMat4Vec4(mvp, v4)

	if clip.W <= 0.00001 {
		return ScreenVertex{}, false
	}

	ndcX := clip.X / clip.W
	ndcY := clip.Y / clip.W
	ndcZ := clip.Z / clip.W
	
	screenX := (ndcX + 1.0) * 0.5 * float64(width)
	screenY := (1.0 - (ndcY+1.0)*0.5) * float64(height)

	return ScreenVertex{
		X: screenX,
		Y: screenY,
		Z: ndcZ,
	}, true
}