package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"time"

	"github.com/dwrolvink/gogl"
	"github.com/go-gl/gl/v4.5-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"

	"image"
	"image/png"

	"github.com/disintegration/imaging"
)

const (
	Width  = 500 // Width of the main window
	Height = 500 // Height of the main window

	RECORDING_GIF   = 0
	RECORDING_PRTSC = 1
)

var (
	WindowTitle                   = "Test GL Application"
	tick           float32        = -1.0                               // ticks up every game loop cycle
	delay_ms       int64          = 5                                  // handles frame rate
	record         bool           = false                              // whether to record the screen.
	record_length  float32        = float32(1.0 * (1000.0 / delay_ms)) // After how many ticks to stop recording (and close the program)
	previous_frame *image.NRGBA                                        // intermediate storage for front buffer pixels
	PF_textureID   gogl.TextureID                                      // previous frame, the backbuffer gets bound to this texture after every draw.

	// actions
	action_reset = false
	action_prtsc = false
)

func main() {
	// Init Window, OpenGL, and Data, get user input from commandline
	window := gogl.Init(WindowTitle, Width, Height)
	data := SetData()
	SetKeyHandling(window)
	ParseCommandlineArgs()

	// Make an image to read the pixel data of the front buffer into
	// (I don't yet know how to put this directly into the texture)
	previous_frame = image.NewNRGBA(image.Rectangle{image.Point{0, 0}, image.Point{Width, Height}})
	PF_textureID = gogl.GenTexture()

	// Main loop
	// ===========================================================
	for !window.ShouldClose() {

		// Update game
		// ---------------------
		start := time.Now()
		tick += 1.0

		// Draw to screen
		// ---------------------
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT) // Clear screen
		DrawDataset(data)                                   // Draw new frame
		window.SwapBuffers()                                // Put buffer that we painted on on the foreground

		// Load current frontbuffer into texture for next round
		// ----------------------------------------------------
		UpdateFrontBufferTexture()

		// Event handling
		// ------------------------------------------------------
		// Handle window events
		glfw.PollEvents()

		// Check if shaders need to be recompiled
		gogl.HotloadShaders()

		if err := gl.GetError(); err != 0 {
			log.Println(err)
		}

		// Record output
		// ------------------------------------------------------
		if action_prtsc {
			CreateImage(int(tick), RECORDING_PRTSC)
			action_prtsc = false
		}

		// FPS management
		// ------------------------------------------------------
		// Sleep for a bit if the loop finished too quickly.
		// A better way would be to update actor positions based on
		// elapsed time, but the neccessary code isn't present yet for
		// that (i.e. volition).
		elapsed := time.Since(start)
		dif_ms := delay_ms - elapsed.Milliseconds()
		time.Sleep(time.Duration(dif_ms * int64(time.Millisecond)))
	}

	// Compile gif
	if record {
		CompileGif()
	}

	// useless here, but good to keep track of what needs to be deleted
	defer glfw.Terminate()
}

func SetKeyHandling(window *glfw.Window) {
	window.SetKeyCallback(func(_ *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
		// Get key as defined by the locale (qwerty, dvorak, etc)
		char := glfw.GetKeyName(key, scancode)
		// Alias for readability
		Down := glfw.Press
		Up := glfw.Release
		_ = Up

		// Handle keystrokes
		switch char {
		case "e":
			if action == Down {
				action_reset = true
			}
		case "p":
			if action == Down {
				action_prtsc = true
			}
		}
	})
}

// Used to reset the FP_texture to initial conditions
func ResetFrame(img *image.NRGBA) {
	pixels, _ := gogl.LoadPixelDataFromImage("assets/img/start.png")

	for i := 0; i < len(img.Pix) && i < len(*pixels); i++ {
		img.Pix[i] = (*pixels)[i]
	}
}

func ReadFrontBuffer() {
	// Get pixels from front buffer, and put them in previous_frame
	gl.ReadBuffer(gl.FRONT)
	gl.ReadPixels(0, 0, Width, Height, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(previous_frame.Pix))
}

func UpdateFrontBufferTexture() {
	// Either reset image, or load image data from frontbuffer
	if action_reset || tick == 1 {
		ResetFrame(previous_frame)
		action_reset = false
	} else {
		ReadFrontBuffer()
	}

	// Bind texture
	gogl.BindTexture(PF_textureID)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)

	// Write image data to the PF_texture
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, int32(Width), int32(Height), 0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(previous_frame.Pix))
}

// Define the DataObjects that contain our Programs, Shaders, Sprites, etc
func SetData() gogl.DataObject {
	/*
		   Multiple datasets can be defined.
		   Each set contains all that it needs to draw to the screen,
		   think of: Program, VOA, VBO, EBO, Textures, Sprites, etc

		   Below, each dataset is defined, added to datalist, and at
		   the end the commandline args are checked what the choice is.
		   Choices include:
		     - Print either dataset 0, or 1 --> -s 0, -s 1
			 - Print both as a composite --> -s c
	*/

	// List of datasets
	datalist := make([]gogl.DataObject, 2)

	// Fist dataset: Vertex type: Quad, uses Sprites
	// -----------------------------------------------------------
	datalist[0] = gogl.DataObject{
		ProgramName: "FrontBufferLoop",
		Type:        gogl.GOGL_QUADS,
		Vertices:    CreateQuadVertexMatrix(1.0, 0.0, 0.0),
		Indices: []uint32{
			1, 0, 3, // triangle 1
			0, 2, 3, // triangle 2
		},
		VertexShaderSource:   "shaders/quad.vert",
		FragmentShaderSource: "shaders/loop.frag",
	}

	datalist[0].ProcessData()

	// Pick one or the other data set
	// -----------------------------------------------------------
	data := datalist[0]

	return data
}

// DRAWING
// ----------------------------------------------------------------------

// Draw a single dataset.
func DrawDataset(data gogl.DataObject) {
	data.Enable()

	// load uniforms

	// load sprite
	/*
		sprite := data.SelectSprite(0)
		sprite.SetUniforms(&data)

		// Draw pepe 1
		data.Program.SetFloat("scale", 0.5)
		gl.DrawElements(gl.TRIANGLES, 6, gl.UNSIGNED_INT, gl.PtrOffset(0))

		// Draw pepe 2
		data.Program.SetFloat("scale", 0.25)
		data.Program.SetFloat("x", -x)
		gl.DrawElements(gl.TRIANGLES, 6, gl.UNSIGNED_INT, gl.PtrOffset(0))

		// Draw pepe 3
		data.Program.SetFloat("scale", 0.16)
		data.Program.SetFloat("x", -x)
		data.Program.SetFloat("y", -x)
		gl.DrawElements(gl.TRIANGLES, 6, gl.UNSIGNED_INT, gl.PtrOffset(0))
	*/
	data.Program.SetFloat("window_width", float32(Width))
	data.Program.SetFloat("window_height", float32(Height))
	gl.BindTexture(gl.TEXTURE_2D, uint32(PF_textureID))
	gl.DrawElements(gl.TRIANGLES, 6, gl.UNSIGNED_INT, gl.PtrOffset(0))

}

// HELPER FUNCTIONS
// ----------------------------------------------------------------------

// Applies commandline args: --fps <N>, --record <N>, --set <'c', N>
func ParseCommandlineArgs() {

	for i := range os.Args {

		// FPS
		// -----------------------------------------------------------
		// Apply commandline choice for fps, if present.
		// Note that for recording, 50 fps is the max.

		if os.Args[i] == "--fps" {
			// check if fps value has been passed directly after --fps
			if i+1 < len(os.Args) {
				choice, err := strconv.Atoi(os.Args[i+1])
				if err == nil {
					delay_ms = int64(1000 / choice)
				} else {
					fmt.Println("ERROR: Could not parse input after --fps as an int.")
				}
			}
		}

		// Record
		// -----------------------------------------------------------
		// Apply commandline choice for recording settings, if present

		if os.Args[i] == "--record" {
			// Enable recording
			record = true
			InitRecording()

			// Check if record_length has been passed
			if i+1 < len(os.Args) {
				choice, err := strconv.Atoi(os.Args[i+1]) // convert string input to int
				if err == nil {
					record_length = float32(int64(choice) * (1000.0 / delay_ms)) // input is in seconds, convert to ticks
				}
			}
		}
	}
}

// Easy way to create a quad with a certain size and offset
func CreateQuadVertexMatrix(size float32, x_offset float32, y_offset float32) []float32 {
	screen_left := -size + x_offset
	screen_bottom := -size + y_offset
	screen_right := size + x_offset
	screen_top := size + y_offset
	texture_top := float32(1.0)
	texture_bottom := float32(0.0)
	texture_left := float32(0.0)
	texture_right := float32(1.0)
	z := float32(0.0)

	vertices := []float32{
		// x, y, z, texcoordx, texcoordy
		screen_left, screen_top, z, texture_left, texture_top,
		screen_right, screen_top, z, texture_right, texture_top,
		screen_left, screen_bottom, z, texture_left, texture_bottom,
		screen_right, screen_bottom, z, texture_right, texture_bottom,
	}

	return vertices
}

// RECORDING
// ----------------------------------------------------------------------

// Ensures that the recording folder (recording/temp) is present.
func InitRecording() {
	err := os.Mkdir("recording/temp/", 0755)
	if err != nil {
		panic(err)
	}
}

// Reads out the pixel data in gl.FRONT, and saves it to recording/temp/image<Tick>.png
func CreateImage(number int, mode int) {
	filename := fmt.Sprintf("image%03d.png", number)
	width := Width
	height := Height

	var folder string
	if mode == RECORDING_PRTSC {
		folder = "recording/printscreens/"

		currentTime := time.Now()
		filename = currentTime.Format("2006-01-02 15:04:05.000000") + ".png"
	} else {
		folder = "recording/temp/"
	}

	img := image.NewNRGBA(image.Rectangle{image.Point{0, 0}, image.Point{width, height}})

	gl.ReadBuffer(gl.FRONT)
	gl.ReadPixels(0, 0, Width, Height, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(img.Pix))
	img = imaging.FlipV(img)

	// Encode as PNG.
	f, _ := os.Create(folder + filename)
	png.Encode(f, img)

	if mode == RECORDING_PRTSC {
		fmt.Println("Created " + folder + filename)
	}

}

// Takes all the frame images in recording/temp and makes a palletted gif out of it using ffmpeg.
func CompileGif() {
	filename := time.Now().Unix()

	cmd, err := exec.Command("/bin/sh", "scripts/make_gif.sh", fmt.Sprint(filename)).Output()
	if err != nil {
		fmt.Printf("error %s", err)
	}
	fmt.Println(cmd)
}
