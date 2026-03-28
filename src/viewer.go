package main

import (
	"fmt"
	"math"
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
	objModPath string
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

func (v *Viewer3D) init(p string) error {
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

	v.objModPath = p
	
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
	rotatAngle := 0.0 
	cam := 	NewOrbitCamera()
	
	frame := NewFramebuff(WINDOW_WIDTH, WINDOW_HEIGHT)
	
	//Vektor segitiga 3D
	// v0 := Vec3{X: -1, Y: -1, Z: 0}
	// v1 := Vec3{X: 1, Y: -1, Z: 0}
	// v2 := Vec3{X: 0, Y: 1, Z: 0}

	//Mesh model kubus
	// cube := Model {
	// 	Vertices : []Vec3{
	// 		{-1.0,-1.0,-1.0}, 
	// 		{ 1.0,-1.0,-1.0}, 
	// 		{ 1.0, 1.0,-1.0},
	// 		{-1.0, 1.0, -1.0},
	// 		{-1.0, -1.0, 1.0},
	// 		{ 1.0, -1.0, 1.0},
	// 		{ 1.0, 1.0, 1.0},
	// 		{-1.0, 1.0, 1.0},
	// 	},
	// 	Faces: []Triangle{
	// 		{0, 1, 2}, {0, 2, 3}, 
	// 		{4, 5, 6}, {4, 6, 7},
	// 		{0, 1, 5}, {0, 5, 4}, 
	// 		{2, 3, 7}, {2, 7, 6},
	// 		{1, 2, 6}, {1, 6, 5},
	// 		{0, 3, 7}, {0, 7, 4},
	// 	},
	// }

	mesh, err := ParseOBJ("cow.obj")
	if err != nil {return err}
	NormalizeModel(mesh)
	
	for true {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch e := event.(type) {
			case *sdl.QuitEvent:
				return err
			case *sdl.KeyboardEvent:
				if e.Type == sdl.KEYDOWN {
					switch e.Keysym.Scancode {
					case sdl.SCANCODE_ESCAPE:
						return err
					case sdl.SCANCODE_W:
						cam.MoveCamera(0,-1,0,0,0)
					case sdl.SCANCODE_S:
						cam.MoveCamera(0,1,0,0,0)
					case sdl.SCANCODE_D:
						cam.MoveCamera(1,0,0,0,0)
					case sdl.SCANCODE_A:
						cam.MoveCamera(-1,0,0,0,0)
					case sdl.SCANCODE_Q:
						cam.MoveCamera(0,0,1,0,0)
					case sdl.SCANCODE_E:
						cam.MoveCamera(0,0,-1,0,0)
					case sdl.SCANCODE_UP:
						cam.MoveCamera(0,0,0,1,0)
					case sdl.SCANCODE_DOWN:
						cam.MoveCamera(0,0,0,-1,0)
					case sdl.SCANCODE_RIGHT:
						cam.MoveCamera(0,0,0,0,1)
					case sdl.SCANCODE_LEFT:
						cam.MoveCamera(0,0,0,0,-1)
					case sdl.SCANCODE_R:
						cam.Reset()
					}
				}
			}
		} 

		v.renderer.Clear()
		frame.ClearFrameBuff(0, 0, 0, 255)
		frame.ClearDepth(math.Inf(1))


		model := Mat4Identity()
		view := cam.ViewMatrix()
		projection := cam.ProjectionMatrix(WINDOW_WIDTH, WINDOW_HEIGHT)

		mvp := MulMat4Mat4(projection, MulMat4Mat4(view, model))

		DrawMeshWireframe(*mesh, mvp, 255, 0, 0, 255, frame)

		if err = v.texture.Update(nil, unsafe.Pointer(&frame.colors[0]), frame.nByteInRow); err != nil {
			return fmt.Errorf("Gagal melakukan randerisasi: %v\n", err)		
		}
		
		v.renderer.Copy(v.texture, nil, nil)
		v.renderer.Present()

		rotatAngle += 0.1
		sdl.Delay(20)
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
