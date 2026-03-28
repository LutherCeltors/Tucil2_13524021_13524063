package main

import (
	"fmt"
	"math"
	"os"
	"unsafe"
	"strings"

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
	mesh2 [2]*Model
	activeMesh int
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

func (v *Viewer3D) init(mesh0, mesh1 *Model) error {
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
	
	v.mesh2[0] = mesh0
	v.mesh2[1] = mesh1

	v.activeMesh = 0
	
	return err
}

func (v *Viewer3D) close() {
	if v != nil {
		v.texture.Destroy()
		v.texture = nil
		v.renderer.Destroy()
		v.renderer = nil
		v.window.Destroy()
		v.window = nil
	}
}

func (v *Viewer3D) run() error {
	var err error
	cam := 	NewOrbitCamera()
	
	frame := NewFramebuff(WINDOW_WIDTH, WINDOW_HEIGHT)
	
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
					case sdl.SCANCODE_TAB:
						v.activeMesh = absVal(v.activeMesh - 1)
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

		mView := MulMat4Mat4(view, model)
		mvp := MulMat4Mat4(projection, mView)

		lightDir := Normalized3(Vec3{1, 1, 1})

		NormalizeModel(v.mesh2[v.activeMesh])

		DrawModelFlatShaded(
			v.mesh2[v.activeMesh],
			mView,
			mvp,
			0, 180, 255, 255,
			lightDir,
			frame,
		)

		if err = v.texture.Update(nil, unsafe.Pointer(&frame.colors[0]), frame.nByteInRow); err != nil {
			return fmt.Errorf("Gagal melakukan randerisasi: %v\n", err)	
		}
		
		v.renderer.Copy(v.texture, nil, nil)
		v.renderer.Present()

		sdl.Delay(20)
	}

	return err
}

type ControlItem struct {
	Key      string
	Function string
}

func PrintViewerControlTable() {
	controls := []ControlItem{
		{"W / S", "Pitch atas / bawah"},
		{"A / D", "Yaw kiri / kanan"},
		{"Q / E", "Zoom out / in"},
		{"Arrow Keys", "Pan target kamera"},
		{"R", "Reset kamera"},
		{"TAB", "Swap model"},
		{"ESC", "Keluar"},
	}

	header1 := "Tombol"
	header2 := "Fungsi"

	col1Width := len(header1)
	col2Width := len(header2)

	for _, c := range controls {
		if len(c.Key) > col1Width {
			col1Width = len(c.Key)
		}
		if len(c.Function) > col2Width {
			col2Width = len(c.Function)
		}
	}

	line := "+" + strings.Repeat("-", col1Width+2) + "+" + strings.Repeat("-", col2Width+2) + "+"

	fmt.Println(line)
	fmt.Printf("| %-*s | %-*s |\n", col1Width, header1, col2Width, header2)
	fmt.Println(line)

	for _, c := range controls {
		fmt.Printf("| %-*s | %-*s |\n", col1Width, c.Key, col2Width, c.Function)
	}

	fmt.Println(line)
}

func Render(realObjPath, voxelObjPath string) {
	var err error
	var mesh0, mesh1 *Model
	
	if mesh0 ,err = ParseOBJ(realObjPath); err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
	} 

	if mesh1 ,err = ParseOBJ(voxelObjPath); err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
	} 

	defer QuitSdl2()
	if err = InitializeSdl2(); err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
		return
	}

	view := newViewer3D()
	
	defer view.close()
	if err = view.init(mesh0, mesh1); err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
		return
	}

	if err = view.run(); err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
	}
		
}
