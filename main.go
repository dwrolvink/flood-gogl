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
)

var (
	WindowTitle = "Test GL Application"

	x        float32 = 0.0  // used to move the sprites around
	dir_x    float32 = 1    // used to change x
	tick     float32 = -1.0 // ticks up every game loop cycle
	delay_ms int64   = 20   // handles frame rate

	DrawMode      string  = "composite"                        // chooses whether to draw one dataset, or all of them ("composite" vs "single_set")
	ChosenDataset int     = 0                                  // Used only when DrawMode = "single_dataset"
	record        bool    = false                              // whether to record the screen.
	record_length float32 = float32(1.0 * (1000.0 / delay_ms)) // After how many ticks to stop recording (and close the program)

	// The texture that we can write to.
	PF_textureID gogl.TextureID
	//testTex      gogl.TextureID
)

func main() {
	// Init Window, OpenGL, and Data, get user input from commandline
	// -----------------------------------------------------------
	window := gogl.Init(WindowTitle, Width, Height)
	data, datalist := SetData()
	_ = datalist

	ParseCommandlineArgs()

	// spike bind front buffer
	// -----------------------
	// make home for the pixel data to be read after swap buffer

	upLeft := image.Point{0, 0}
	lowRight := image.Point{Width, Height}
	previous_frame := image.NewNRGBA(image.Rectangle{upLeft, lowRight})

	PF_textureID = gogl.GenTexture()

	//testTex = gogl.LoadImageToTexture("assets/img/texture.png")

	// Main loop
	// ===========================================================
	for !window.ShouldClose() && (!record || tick < record_length) {

		// Update game
		// ------------------------------------------------------
		// Naive way to manage FPS. See also bottom of this loop.
		start := time.Now()

		// Increment global clock
		tick += 1.0

		// x is used temporarily to move stuff around, will be removed when
		// there are actors with volition
		x += 0.01 * dir_x
		if x > 1.0 || x < -1.0 {
			dir_x *= -1.0
		}

		// Update DataObjects
		for i := range datalist {
			datalist[i].Update()
		}

		// Draw to screen
		// ------------------------------------------------------
		// Clear screen
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		// Draw new frame
		DrawDataset(data)

		// Put buffer that we painted on on the foreground
		window.SwapBuffers()

		// spike bind front buffer
		// -----------------------

		// update front buffer texture
		gogl.BindTexture(PF_textureID)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.REPEAT)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)

		// Bind front buffer to previous_frame
		gl.ReadBuffer(gl.FRONT)
		gl.ReadPixels(0, 0, Width, Height, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(previous_frame.Pix))

		previous_frame.Pix[0+4*Width*200] = 255.0
		previous_frame.Pix[3+4*Width*200] = 255.0

		// Bind to texture ID
		gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, int32(Width), int32(Height), 0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(previous_frame.Pix))

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
		if record {
			CreateImage(int(tick))
			fmt.Println(tick)
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

// Define the DataObjects that contain our Programs, Shaders, Sprites, etc
func SetData() (gogl.DataObject, []gogl.DataObject) {
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

	return data, datalist
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

		// DrawMode & ChosenDataset
		// -----------------------------------------------------------
		// Apply commandline choice for dataset

		if os.Args[i] == "--set" {
			if i+1 < len(os.Args) {

				// Print both datasets on top of eachother
				if os.Args[i+1] == "c" {
					DrawMode = "composite"
					continue
				}

				// Print only one dataset
				DrawMode = "single_set"
				choice, err := strconv.Atoi(os.Args[i+1])
				if err != nil {
					fmt.Println("ERROR: Dataset index not passed in. E.g. '-s 1'. Ignoring.")
					continue
				}
				ChosenDataset = choice

			} else {
				fmt.Println("ERROR: Dataset index not passed in. E.g. '-s 1'. Ignoring.")
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
func CreateImage(number int) {
	filename := fmt.Sprintf("image%03d.png", number)
	width := Width
	height := Height

	upLeft := image.Point{0, 0}
	lowRight := image.Point{width, height}

	img := image.NewNRGBA(image.Rectangle{upLeft, lowRight})

	gl.ReadBuffer(gl.FRONT)
	gl.ReadPixels(0, 0, Width, Height, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(img.Pix))
	img = imaging.FlipV(img)

	// Encode as PNG.
	f, _ := os.Create("recording/temp/" + filename)
	png.Encode(f, img)
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
