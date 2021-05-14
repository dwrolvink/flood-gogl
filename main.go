package main

import (
	//"fmt"
	"log"
	//"os"
	//"strconv"
	"time"

	"github.com/dwrolvink/gogl"
	"github.com/go-gl/gl/v4.5-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"

	"image"
)

const (
	Width  = 500 // Width of the main window
	Height = 500 // Height of the main window
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
	ActionReset = false
	ActionPrtsc = false
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
		if ActionPrtsc {
			CreateImage(int(tick), RECORDING_PRTSC)
			ActionPrtsc = false
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

// DRAWING
// ----------------------------------------------------------------------

// Draw a single dataset.
func DrawDataset(data gogl.DataObject) {
	data.Enable()

	data.Program.SetFloat("window_width", float32(Width))
	data.Program.SetFloat("window_height", float32(Height))
	gl.BindTexture(gl.TEXTURE_2D, uint32(PF_textureID))
	gl.DrawElements(gl.TRIANGLES, 6, gl.UNSIGNED_INT, gl.PtrOffset(0))

}
