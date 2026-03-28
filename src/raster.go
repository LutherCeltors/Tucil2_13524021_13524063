package main

import "math"

func absVal(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// masih belum menggunakan anti aliasing, masih dengna bresenham line algorithm
func DrawLine(x0, y0 int, x1, y1 int, r, g, b, a byte, frame *Framebufer) {
	dx := absVal(x1 - x0)
	dy := absVal(y1 - y0)

	sx := -1
	if x0 < x1 {
		sx = 1
	}

	sy := -1
	if y0 < y1 {
		sy = 1
	}

	s := dx - dy

	for {
		_ = frame.SetPixelColors(x0, y0, r, g, b, a)

		if x0 == x1 && y0 == y1 {
			break
		}

		e2 := 2 * s

		if e2 > -dy {
			s -= dy
			x0 += sx
		}
		if e2 < dx {
			s += dx
			y0 += sy
		}
	}	
}

func RoundToInt (val float64) int {
	return int(math.Round(val))
}

func DrawTriangle (x0, y0, x1, y1, x2, y2 int, r, g, b, a byte, frame *Framebufer) {
	DrawLine(x0, y0, x1, y1, r, g, b, a, frame)
	DrawLine(x0, y0, x2, y2, r, g, b, a, frame)
	DrawLine(x1, y1, x2, y2, r, g, b, a, frame)
}

func DrawTriangle3D(v0, v1, v2 Vec3, mView ,mvp Mat4, r, g, b, a byte, frame *Framebufer) {
	if IsBackFace(v0, v1, v2, mView) {return}

	p0, ok0 := ProjectVertex(v0, mvp, frame.width, frame.height)
	p1, ok1 := ProjectVertex(v1, mvp, frame.width, frame.height)
	p2, ok2 := ProjectVertex(v2, mvp, frame.width, frame.height)

	if !ok0 || !ok1 || !ok2 {
		return
	}

	DrawTriangle( RoundToInt(p0.X), RoundToInt(p0.Y), RoundToInt(p1.X), RoundToInt(p1.Y), RoundToInt(p2.X), RoundToInt(p2.Y), r, g, b, a, frame)
}


func DrawMeshWireframe(mesh Model, mView, mvp Mat4, r, g, b, a byte, frame *Framebufer) {
	for _, f := range mesh.Faces {
		DrawTriangle3D(
			mesh.Vertices[f.A],
			mesh.Vertices[f.B],
			mesh.Vertices[f.C],
			mView,
			mvp,
			r, g, b, a,
			frame,
		)
	}
}

func TransformPoint(m Mat4, v Vec3) Vec3 {
	p := Vec4{X: v.X, Y: v.Y, Z: v.Z, W: 1.0}
	res := MulMat4Vec4(m, p)

	return Vec3{
		X: res.X,
		Y: res.Y,
		Z: res.Z,
	}
}

func IsBackFace(v0, v1, v2 Vec3, modelView Mat4) bool {
	p0 := TransformPoint(modelView, v0)
	p1 := TransformPoint(modelView, v1)
	p2 := TransformPoint(modelView, v2)

	e1 := Sub3(p1, p0)
	e2 := Sub3(p2, p0)

	normal := cross(e1, e2)

	toCamera := Vec3{
		X: -p0.X,
		Y: -p0.Y,
		Z: -p0.Z,
	}

	return dot(normal, toCamera) <= 0
}

func minInt(a, b int) int {
	if a < b {return a}
	return b
}

func maxInt(a, b int) int {
	if a > b {return a}
	return b
}

func clampInt(val, min, max int) int {
	if val < min {return min}
	if val > max {return max}
	return val
}

func EdgeFunction(ax, ay, bx, by, px, py float64) float64 {
	return (px-ax)*(by-ay) - (py-ay)*(bx-ax)
}

func FillTriangle2D(x0, y0, x1, y1, x2, y2 int, r, g, b, a byte, frame *Framebufer) {
	minX := clampInt(minInt(x0, minInt(x1, x2)), 0, frame.width-1)
	maxX := clampInt(maxInt(x0, maxInt(x1, x2)), 0, frame.width-1)
	minY := clampInt(minInt(y0, minInt(y1, y2)), 0, frame.height-1)
	maxY := clampInt(maxInt(y0, maxInt(y1, y2)), 0, frame.height-1)

	ax, ay := float64(x0), float64(y0)
	bx, by := float64(x1), float64(y1)
	cx, cy := float64(x2), float64(y2)

	area := EdgeFunction(ax, ay, bx, by, cx, cy)
	if area == 0 {
		return
	}

	for y := minY; y <= maxY; y++ {
		for x := minX; x <= maxX; x++ {
			px := float64(x) + 0.5
			py := float64(y) + 0.5

			w0 := EdgeFunction(bx, by, cx, cy, px, py)
			w1 := EdgeFunction(cx, cy, ax, ay, px, py)
			w2 := EdgeFunction(ax, ay, bx, by, px, py)

			inside := (w0 >= 0 && w1 >= 0 && w2 >= 0) || (w0 <= 0 && w1 <= 0 && w2 <= 0)
			if !inside {
				continue
			}

			alpha := w0 / area
			beta := w1 / area
			gamma := w2 / area

			_ = alpha
			_ = beta
			_ = gamma

			_ = frame.SetPixelColors(x, y, r, g, b, a)
		}
	}
}

func DrawFilledTriangle2D(x0, y0, x1, y1, x2, y2 int, fillR, fillG, fillB, fillA byte, lineR, lineG, lineB, lineA byte, frame *Framebufer) {
	FillTriangle2D(x0, y0, x1, y1, x2, y2, fillR, fillG, fillB, fillA, frame)
	DrawTriangle(x0, y0, x1, y1, x2, y2, lineR, lineG, lineB, lineA, frame)
}

func FillTriangleProjected(p0, p1, p2 ScreenVertex, r, g, b, a byte, frame *Framebufer) {
	x0 := RoundToInt(p0.X)
	y0 := RoundToInt(p0.Y)
	x1 := RoundToInt(p1.X)
	y1 := RoundToInt(p1.Y)
	x2 := RoundToInt(p2.X)
	y2 := RoundToInt(p2.Y)

	minX := clampInt(minInt(x0, minInt(x1, x2)), 0, frame.width-1)
	maxX := clampInt(maxInt(x0, maxInt(x1, x2)), 0, frame.width-1)
	minY := clampInt(minInt(y0, minInt(y1, y2)), 0, frame.height-1)
	maxY := clampInt(maxInt(y0, maxInt(y1, y2)), 0, frame.height-1)

	ax, ay := p0.X, p0.Y
	bx, by := p1.X, p1.Y
	cx, cy := p2.X, p2.Y

	area := EdgeFunction(ax, ay, bx, by, cx, cy)
	if area == 0 {
		return
	}

	for y := minY; y <= maxY; y++ {
		for x := minX; x <= maxX; x++ {
			px := float64(x) + 0.5
			py := float64(y) + 0.5

			w0 := EdgeFunction(bx, by, cx, cy, px, py)
			w1 := EdgeFunction(cx, cy, ax, ay, px, py)
			w2 := EdgeFunction(ax, ay, bx, by, px, py)

			inside := (w0 >= 0 && w1 >= 0 && w2 >= 0) || (w0 <= 0 && w1 <= 0 && w2 <= 0)
			if !inside {
				continue
			}

			alpha := w0 / area
			beta  := w1 / area
			gamma := w2 / area

			z := alpha*p0.Z + beta*p1.Z + gamma*p2.Z

			frame.SetPixelDepthColor(x, y, z, r, g, b, a)
		}
	}
}

func DrawFilledTriangle3D(v0, v1, v2 Vec3, modelView, mvp Mat4, r, g, b, a byte, frame *Framebufer) {
	if IsBackFace(v0, v1, v2, modelView) {
		return
	}

	p0, ok0 := ProjectVertex(v0, mvp, frame.width, frame.height)
	p1, ok1 := ProjectVertex(v1, mvp, frame.width, frame.height)
	p2, ok2 := ProjectVertex(v2, mvp, frame.width, frame.height)

	if !ok0 || !ok1 || !ok2 {
		return
	}

	FillTriangleProjected(p0, p1, p2, r, g, b, a, frame)
}

func DrawMeshFilled(model *Model, modelView, mvp Mat4, r, g, b, a byte, frame *Framebufer) {
	if model == nil {
		return
	}

	for _, face := range model.Faces {
		DrawFilledTriangle3D(
			model.Vertices[face.A],
			model.Vertices[face.B],
			model.Vertices[face.C],
			modelView,
			mvp,
			r, g, b, a,
			frame,
		)
	}
}

func ClampByte(x float64) byte {
	if x < 0 {
		return 0
	}
	if x > 255 {
		return 255
	}
	return byte(x)
}

func ComputeFaceNormalView(v0, v1, v2 Vec3, modelView Mat4) Vec3 {
	p0 := TransformPoint(modelView, v0)
	p1 := TransformPoint(modelView, v1)
	p2 := TransformPoint(modelView, v2)

	e1 := Sub3(p1, p0)
	e2 := Sub3(p2, p0)

	return Normalized3(cross(e1, e2))
}

func ApplyFlatShading(baseR, baseG, baseB byte, normal, lightDir Vec3) (byte, byte, byte) {
	intensity := dot(normal, Normalized3(lightDir))
	if intensity < 0 {
		intensity = 0
	}

	ambient := 0.15
	brightness := ambient + (1.0-ambient)*intensity

	r := ClampByte(float64(baseR) * brightness)
	g := ClampByte(float64(baseG) * brightness)
	b := ClampByte(float64(baseB) * brightness)

	return r, g, b
}

func DrawFlatShadedTriangle3D(
	v0, v1, v2 Vec3,
	modelView, mvp Mat4,
	baseR, baseG, baseB, a byte,
	lightDir Vec3,
	frame *Framebufer,
) {
	if IsBackFace(v0, v1, v2, modelView) {
		return
	}

	p0, ok0 := ProjectVertex(v0, mvp, frame.width, frame.height)
	p1, ok1 := ProjectVertex(v1, mvp, frame.width, frame.height)
	p2, ok2 := ProjectVertex(v2, mvp, frame.width, frame.height)

	if !ok0 || !ok1 || !ok2 {
		return
	}

	normal := ComputeFaceNormalView(v0, v1, v2, modelView)
	r, g, b := ApplyFlatShading(baseR, baseG, baseB, normal, lightDir)

	FillTriangleProjected(p0, p1, p2, r, g, b, a, frame)
}

func DrawModelFlatShaded(
	model *Model,
	modelView, mvp Mat4,
	baseR, baseG, baseB, a byte,
	lightDir Vec3,
	frame *Framebufer,
) {
	if model == nil {
		return
	}

	for _, face := range model.Faces {
		DrawFlatShadedTriangle3D(
			model.Vertices[face.A],
			model.Vertices[face.B],
			model.Vertices[face.C],
			modelView, mvp,
			baseR, baseG, baseB, a,
			lightDir,
			frame,
		)
	}
}