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
	Width  = 480 // Width of the main window
	Height = 480 // Height of the main window
)

var (
	WindowTitle           = "Test GL Application"
	tick          float32 = -1.0                                 // ticks up every game loop cycle
	delay_ms      int64   = 30                                   // handles frame rate
	record        bool    = false                                // whether to record the screen.
	record_length float32 = float32(1.0*(1000.0/delay_ms)) * 0.5 // After how many ticks to stop recording (and close the program)

	window   = gogl.Init(WindowTitle, Width, Height) // Init Window, OpenGL
	datalist = SetData()                             // Init data

	PfGameImg        *image.NRGBA   // intermediate storage for front buffer pixels - game state
	PfSmellImg       *image.NRGBA   // intermediate storage for front buffer pixels - smell (blurred gamestate)
	PfGameTextureID  gogl.TextureID // previous frame, the backbuffer gets bound to this texture after every draw.
	PfSmellTextureID gogl.TextureID //
	FrameBuffer1     uint32         // can be used as a draw target when we don't want to draw to the screen

	// actions
	ActionReset      = false
	ActionPrtsc      = false
	ActionRecord     = false
	ActionPrintSmell = true
)

func main() {
	// Get user input from commandline & link keypresses to actions
	SetKeyHandling(window)
	ParseCommandlineArgs()

	// Make an image to read the pixel data of the front buffer into
	// (I don't yet know how to put this directly into the texture)
	PfGameImg = image.NewNRGBA(image.Rectangle{image.Point{0, 0}, image.Point{Width, Height}})
	PfSmellImg = image.NewNRGBA(image.Rectangle{image.Point{0, 0}, image.Point{Width, Height}})
	PfGameTextureID = gogl.GenTexture()
	PfSmellTextureID = gogl.GenTexture()

	// draw target for smell
	// ---------------------

	// generate a framebuffer
	gl.GenFramebuffers(1, &FrameBuffer1)

	// Load textures
	ResetFrame(PfGameTextureID, PfGameImg, "assets/img/texture.png", 0)
	ResetFrame(PfSmellTextureID, PfSmellImg, "assets/img/texture.png", 1)

	// attach texture to fbo
	gl.BindFramebuffer(gl.DRAW_FRAMEBUFFER, FrameBuffer1)
	gl.FramebufferTexture2D(gl.FRAMEBUFFER, gl.COLOR_ATTACHMENT0, gl.TEXTURE_2D, uint32(PfSmellTextureID), 0)

	// check state of framebuffer
	if gl.CheckFramebufferStatus(gl.FRAMEBUFFER) != gl.FRAMEBUFFER_COMPLETE {
		panic("framebuffer not complete!")
	}

	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)

	//gl.FramebufferParameteri(gl.DRAW_FRAMEBUFFER, gl.FRAMEBUFFER_DEFAULT_WIDTH, Width);
	//gl.FramebufferParameteri(gl.DRAW_FRAMEBUFFER, gl.FRAMEBUFFER_DEFAULT_HEIGHT, Height);
	//gl.FramebufferParameteri(gl.DRAW_FRAMEBUFFER, gl.FRAMEBUFFER_DEFAULT_SAMPLES, 4);

	// Main loop
	// ===========================================================
	for !window.ShouldClose() {

		// Update game
		// ---------------------
		start := time.Now()
		tick += 1.0

		if ActionReset {
			ResetFrame(PfGameTextureID, PfGameImg, "assets/img/texture.png", 0)
			ResetFrame(PfSmellTextureID, PfSmellImg, "assets/img/texture.png", 1)
			ActionReset = false
		}

		// Draw to screen
		// ---------------------
		DrawDataset(datalist) // Draw new frame

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
		if ActionRecord {
			CreateImage(int(tick), RECORDING_GIF)
			if tick > record_length {
				ActionRecord = false
				CompileGif()
			}
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

	// useless here, but good to keep track of what needs to be deleted
	defer glfw.Terminate()
}

// DRAWING
// ----------------------------------------------------------------------

// Draw a single dataset.
func DrawDataset(datalist []gogl.DataObject) {

	// Clear screen
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	// Update game state
	// ==============================================================

	// Enable program, VAO, etc
	data := datalist[0]
	data.Enable()

	// Set uniforms
	data.Program.SetFloat("window_width", float32(Width))
	data.Program.SetFloat("window_height", float32(Height))

	// Bind textures
	data.Program.SetInt("PfGameTexture", int32(0))
	data.Program.SetInt("PfSmellTexture", int32(1))
	gl.ActiveTexture(gl.TEXTURE0 + 0)                      // Game state
	gl.BindTexture(gl.TEXTURE_2D, uint32(PfGameTextureID)) //

	gl.ActiveTexture(gl.TEXTURE0 + 1)                       // Smell (blurred composite of game state)
	gl.BindTexture(gl.TEXTURE_2D, uint32(PfSmellTextureID)) //

	// Draw Game state
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.ONE, gl.ZERO)
	gl.DrawElements(gl.TRIANGLES, 6, gl.UNSIGNED_INT, gl.PtrOffset(0))

	TakeBufferSnapshot(gl.BACK, PfGameTextureID, PfGameImg, 0) // overwrites PfGameTexture

	// Calc smell
	// =======================================================

	// bind it as the target for rendering commands
	//gl.BindFramebuffer(gl.DRAW_FRAMEBUFFER, FrameBuffer1)

	data = datalist[1]
	data.Enable()

	// Set uniforms
	data.Program.SetFloat("window_width", float32(Width))
	data.Program.SetFloat("window_height", float32(Height))

	// Bind textures
	data.Program.SetInt("PfGameTexture", int32(0))
	data.Program.SetInt("PfSmellTexture", int32(1))
	gl.ActiveTexture(gl.TEXTURE0 + 0)                      // Game state
	gl.BindTexture(gl.TEXTURE_2D, uint32(PfGameTextureID)) //

	gl.ActiveTexture(gl.TEXTURE0 + 1)                       // Smell (blurred composite of game state)
	gl.BindTexture(gl.TEXTURE_2D, uint32(PfSmellTextureID)) //

	// Draw Smell
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	gl.DrawElements(gl.TRIANGLES, 6, gl.UNSIGNED_INT, gl.PtrOffset(0))

	TakeBufferSnapshot(gl.BACK, PfSmellTextureID, PfSmellImg, 1)

	// Overwrite backbuffer with gamestate
	// =============================================================

	if ActionPrintSmell == false {
		// Enable program, VAO, etc
		data := datalist[2]
		data.Enable()

		// Set uniforms
		data.Program.SetFloat("window_width", float32(Width))
		data.Program.SetFloat("window_height", float32(Height))

		// Bind textures
		data.Program.SetInt("PfGameTexture", int32(0))
		data.Program.SetInt("PfSmellTexture", int32(1))
		gl.ActiveTexture(gl.TEXTURE0 + 0)                      // Game state
		gl.BindTexture(gl.TEXTURE_2D, uint32(PfGameTextureID)) //

		gl.ActiveTexture(gl.TEXTURE0 + 1)                       // Smell (blurred composite of game state)
		gl.BindTexture(gl.TEXTURE_2D, uint32(PfSmellTextureID)) //

		// Draw Game state
		gl.Enable(gl.BLEND)
		gl.BlendFunc(gl.ONE, gl.ZERO)
		gl.DrawElements(gl.TRIANGLES, 6, gl.UNSIGNED_INT, gl.PtrOffset(0))
	}

	// =========================================================

	// Put buffer that we painted on on the foreground
	window.SwapBuffers()

}
