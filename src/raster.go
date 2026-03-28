package main

import (
	"fmt"
	"math"
)

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

func DrawTriangle3D(v0, v1, v2 Vec3, mvp Mat4, r, g, b, a byte, frame *Framebufer) {
	p0, ok0 := ProjectVertex(v0, mvp, frame.width, frame.height)
	p1, ok1 := ProjectVertex(v1, mvp, frame.width, frame.height)
	p2, ok2 := ProjectVertex(v2, mvp, frame.width, frame.height)

	if !ok0 || !ok1 || !ok2 {
		fmt.Printf("Ditemukan vertex yang tidak valid\n")
		return
	}

	DrawTriangle( RoundToInt(p0.X), RoundToInt(p0.Y), RoundToInt(p1.X), RoundToInt(p1.Y), RoundToInt(p2.X), RoundToInt(p2.Y), r, g, b, a, frame)
}


func DrawMeshWireframe(mesh Model, mvp Mat4, r, g, b, a byte, frame *Framebufer) {
	for _, f := range mesh.Faces {
		DrawTriangle3D(
			mesh.Vertices[f.A],
			mesh.Vertices[f.B],
			mesh.Vertices[f.C],
			mvp,
			r, g, b, a,
			frame,
		)
	}
}