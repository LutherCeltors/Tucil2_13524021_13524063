package main

import (
	"fmt"
	"math"
)

type Framebufer struct {
	width int
	height int
	nByteInRow int
	depth []float64
	colors []byte
}

func NewFramebuff(width int, height int) *Framebufer {
	frame := &Framebufer{
		width: width,
		height: height,
		nByteInRow: width * 4,
		colors: make([]byte, width * height * 4),
		depth: make([]float64, width * height),
	}
	frame.ClearFrameBuff(0,0,0,255);
	frame.ClearDepth(math.Inf(1))

	return frame
}

func (frame *Framebufer) ClearFrameBuff(r, g, b, a byte) {
	for i := 0; i < frame.width; i++ {
		for j := 0; j < frame.height; j++{
			idx := i * 4 + frame.nByteInRow * j
			frame.colors[idx + 0] = r
			frame.colors[idx + 1] = g
			frame.colors[idx + 2] = b
			frame.colors[idx + 3] = a  
		}
	}	
}

func (frame *Framebufer) ClearDepth(n float64) {
	for i := range frame.depth {
		frame.depth[i] = n
	}
}

func (frame *Framebufer) SetPixelColors(x, y int, r, g, b, a byte) error {
	var err error
	if x < 0 || x >= frame.width || y < 0 || y >= frame.height {
		return fmt.Errorf("Gagal mengubah warna: x = %d dan y = %d, sementara width, height = %d, %d\n", x, y, frame.width, frame.height)
	}

	idx := x * 4 + frame.nByteInRow * y
	frame.colors[idx + 0] = r
	frame.colors[idx + 1] = g
	frame.colors[idx + 2] = b
	frame.colors[idx + 3] = a  

	return err
}

func (frame *Framebufer) SetDepth(x, y int, val float64) error {
	var err error
	if x < 0 || x >= frame.width || y < 0 || y >= frame.height {
		return fmt.Errorf("Gagal mengubah warna: x = %d dan y = %d, sementara width, height = %d, %d\n", x, y, frame.width, frame.height)
	}
	
	idx := x + y * frame.width
	if frame.depth[idx] > val {
		frame.depth[idx] = val
	}

	return err
}

func (frame *Framebufer) SetPixelDepthColor(x, y int, z float64, r, g, b, a byte) {
	if x < 0 || x >= frame.width || y < 0 || y >= frame.height {
		return
	}

	depthIdx := x + y*frame.width
	if z >= frame.depth[depthIdx] {
		return
	}

	frame.depth[depthIdx] = z

	colorIdx := x*4 + frame.nByteInRow*y
	frame.colors[colorIdx+0] = r
	frame.colors[colorIdx+1] = g
	frame.colors[colorIdx+2] = b
	frame.colors[colorIdx+3] = a
}