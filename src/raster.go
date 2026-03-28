package main

import (
	"fmt"
	"math"
	"os"
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

	err := dx - dy

	for {
		if pixErr := frame.SetPixelColors(x0, y0, r, g, b, a); pixErr != nil {
			fmt.Fprintf(os.Stderr, "%v", pixErr)
			return
		}

		if x0 == x1 && y0 == y1 {
			break
		}

		e2 := 2 * err

		if e2 > -dy {
			err -= dy
			x0 += sx
		}
		if e2 < dx {
			err += dx
			y0 += sy
		}
	}	
}

func RoundToInt (val float64) int {
	return int(math.Round(val))
}