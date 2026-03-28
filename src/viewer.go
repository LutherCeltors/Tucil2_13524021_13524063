package main

import (
	"fmt"
	"os"
	"unsafe"

	"github.com/veandco/go-sdl2/sdl"
)

const (
	WINDOW_WIDTH  = 800
	WINDOW_HEIGHT = 600
	WINDOWN_TITLE = "3D Viewer"
)

type Viewer3D struct {
	window *sdl.Window
	renderer *sdl.Renderer
	texture *sdl.Texture
}

func InitializeSdl2() error {
	var errSdl error
	var sdlFlags uint32 = sdl.INIT_EVERYTHING

	errSdl = sdl.Init(sdlFlags)

	if errSdl != nil {
		return fmt.Errorf("Gagal menginisiasi SDL2 : %v\n", errSdl)
	}

	return errSdl 
}

func QuitSdl2() {
	sdl.Quit();
	println("Keluar dari sdl2")
}

func newViewer3D() *Viewer3D {
	v := &Viewer3D{}
	return v
}

func (v *Viewer3D) init() error {
	var err error

	if v.window, err = sdl.CreateWindow(WINDOWN_TITLE, sdl.WINDOWPOS_CENTERED, sdl.WINDOWPOS_CENTERED, WINDOW_WIDTH, WINDOW_HEIGHT, sdl.WINDOW_SHOWN); err != nil {
		return fmt.Errorf("Gagal menginisiasi window : %v\n", err)
	}

	if v.renderer, err = sdl.CreateRenderer(v.window, -1, sdl.RENDERER_ACCELERATED); err != nil {
		return fmt.Errorf("Gagal menginisiasi renderer : %v\n", err)
	}
	
	if v.texture, err = v.renderer.CreateTexture(uint32(sdl.PIXELFORMAT_RGBA32), sdl.TEXTUREACCESS_STREAMING, int32(WINDOW_WIDTH), int32(WINDOW_HEIGHT)); err != nil {
		return fmt.Errorf("Gagal memuat texture sebagai sarana rendering: %v\n", err)
	}
	
	return err
}

func (v *Viewer3D) close() {
	if v != nil {
		v.texture.Destroy()
		println("Texture")
		println(v.texture)
		v.texture = nil
		println(v.texture)
		v.renderer.Destroy()
		println("Renderer")
		println(v.renderer)
		v.renderer = nil
		println(v.renderer)
		v.window.Destroy()
		println("Window")
		println(v.window)
		v.window = nil
		println(v.window)
	}
}

func (v *Viewer3D) run() error {
	var err error

	v.renderer.Clear()
	
	frame := NewFramebuff(WINDOW_WIDTH, WINDOW_HEIGHT)
	DrawLine(0,0, 700, 500, 255, 20, 20, 255, frame)
	DrawLine(600,100, 100, 500, 255, 255, 20, 255, frame)

	err = v.texture.Update(nil, unsafe.Pointer(&frame.colors[0]), frame.nByteInRow)
	if err != nil {
		return fmt.Errorf("Gagal melakukan randerisasi: %v\n", err)		
	}

	v.renderer.Clear()
	v.renderer.Copy(v.texture, nil, nil)
	v.renderer.Present()

	for true {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				return err
			}
		} 
	}

	return err
}

func main() {
	var err error

	defer QuitSdl2()
	if err = InitializeSdl2(); err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
		return
	}

	view := newViewer3D()
	defer view.close()
	if err = view.init(); err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
		return
	}

	if err = view.run(); err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
	}
}
