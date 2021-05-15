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

	GL_DEFAULT_FBO = 0

	DRAW_MODE_ADD   = 1 // Just add both together
	DRAW_MODE_MERGE = 2 // uniformf A is used as mix, 0=game, 1=smell
	DRAW_MODE_SMELL = 3 // uniformf A is used as mix, 0=red, 0.5=both, 1.0=green

	StartImageSrc = "assets/img/text2.png"
)

var (
	WindowTitle           = "Test GL Application"
	tick          float32 = 0.0                                  // ticks up every game loop cycle
	delay_ms      int64   = 20                                   // handles frame rate
	record        bool    = false                                // whether to record the screen.
	record_length float32 = float32(1.0*(1000.0/delay_ms)) * 0.5 // After how many ticks to stop recording (and close the program)

	window = gogl.Init(WindowTitle, Width, Height) // Init Window, OpenGL

	datalist       = SetData() // Init data
	DataGame       = datalist[0]
	DataSmellRed   = datalist[1]
	DataSmellGreen = datalist[2]
	DataMix        = datalist[3]

	PfGameImg             *image.NRGBA   // intermediate storage for front buffer pixels - game state
	PfSmellImg            *image.NRGBA   // intermediate storage for front buffer pixels - smell (blurred gamestate)
	PfGameTextureID       gogl.TextureID // previous frame, the backbuffer gets bound to this texture after every draw.
	PfSmellRedTextureID   gogl.TextureID //
	PfSmellGreenTextureID gogl.TextureID //
	FBGame                uint32         // can be used as a draw target when we don't want to draw to the screen
	FBSmellRed            uint32         // can be used as a draw target when we don't want to draw to the screen
	FBSmellGreen          uint32         // can be used as a draw target when we don't want to draw to the screen

	// actions
	ActionReset              = false
	ActionPrtsc              = false
	ActionRecord             = false
	ActionPrintSmell         = true
	ActionPrintGame          = true
	ActionDrawMode   int32   = DRAW_MODE_ADD
	ActionDrawA      float32 = 0.5
)

func main() {
	// Get user input from commandline & link keypresses to actions
	SetKeyHandling(window)
	ParseCommandlineArgs()

	// Make an image to read the pixel data of the front buffer into
	// (I don't yet know how to put this directly into the texture)
	PfGameTextureID = NewDefaultTexture()
	PfSmellRedTextureID = NewDefaultTexture()
	PfSmellGreenTextureID = NewDefaultTexture()

	//Create render target for when gamestate is not outputted to the screen
	ResetFrame(PfGameTextureID, StartImageSrc) // texture needs to be initialized before linking to FBO
	CreateFramebuffer(&FBGame, PfGameTextureID)

	// Home for red smell
	ResetFrame(PfSmellRedTextureID, StartImageSrc) // texture needs to be initialized before linking to FBO
	CreateFramebuffer(&FBSmellRed, PfSmellRedTextureID)

	// Home for green smell
	ResetFrame(PfSmellGreenTextureID, StartImageSrc) // texture needs to be initialized before linking to FBO
	CreateFramebuffer(&FBSmellGreen, PfSmellGreenTextureID)

	// Main loop
	for !window.ShouldClose() {

		// Housekeeping
		start := time.Now()
		tick += 1.0

		if ActionReset || tick == 1.0 {
			ResetFrame(PfGameTextureID, StartImageSrc)
			ResetFrame(PfSmellRedTextureID, StartImageSrc)
			ResetFrame(PfSmellGreenTextureID, StartImageSrc)
			ActionReset = false
		}

		// Draw to screen / update gamestate
		UpdateGame()
		DrawDataset()

		// Handle window events
		glfw.PollEvents()

		// Check if shaders need to be recompiled
		gogl.HotloadShaders()

		// Sanity check
		if err := gl.GetError(); err != 0 {
			log.Println(err)
		}

		// Record output
		if ActionPrtsc {
			// Printscreens
			CreateImage(int(tick), RECORDING_PRTSC)
			ActionPrtsc = false
		}
		if ActionRecord {
			// Gif recording
			CreateImage(int(tick), RECORDING_GIF)
			if tick > record_length {
				ActionRecord = false
				CompileGif()
			}
		}

		// Sleep for a bit if the loop finished too quickly.
		// Note that this does nothing for if a loop takes too long.
		elapsed := time.Since(start)
		dif_ms := delay_ms - elapsed.Milliseconds()
		time.Sleep(time.Duration(dif_ms * int64(time.Millisecond)))
	}

	// useless here, but good to keep track of what needs to be deleted
	defer glfw.Terminate()
}

func UpdateGame() {
	// set to overwrite whatever is on the texture (no mixing)
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.ONE, gl.ZERO)

	// Update game
	// =======================================================
	// Select game buffer
	gl.BindFramebuffer(gl.FRAMEBUFFER, FBGame)

	// Enable game shader program
	data := DataGame
	data.Enable()

	// Set uniforms
	data.Program.SetFloat("window_width", float32(Width))
	data.Program.SetFloat("window_height", float32(Height))

	// Bind textures
	data.Program.SetInt("PfGameTexture", int32(0))
	data.Program.SetInt("PfSmellTexture", int32(1))
	gl.ActiveTexture(gl.TEXTURE0 + 0)                      // Game state
	gl.BindTexture(gl.TEXTURE_2D, uint32(PfGameTextureID)) //

	gl.ActiveTexture(gl.TEXTURE0 + 1)                          // Smell (blurred composite of game state)
	gl.BindTexture(gl.TEXTURE_2D, uint32(PfSmellRedTextureID)) //

	// Draw Game state
	gl.DrawElements(gl.TRIANGLES, 6, gl.UNSIGNED_INT, gl.PtrOffset(0))

	// Update red smell
	// =======================================================
	gl.BindFramebuffer(gl.FRAMEBUFFER, FBSmellRed)

	data = DataSmellRed
	data.Enable()

	// Set uniforms
	data.Program.SetFloat("window_width", float32(Width))
	data.Program.SetFloat("window_height", float32(Height))

	// Bind textures
	data.Program.SetInt("PfGameTexture", int32(0))
	data.Program.SetInt("PfSmellTexture", int32(1))
	gl.ActiveTexture(gl.TEXTURE0 + 0)                      // Game state
	gl.BindTexture(gl.TEXTURE_2D, uint32(PfGameTextureID)) //

	gl.ActiveTexture(gl.TEXTURE0 + 1)                          // Smell (blurred composite of game state)
	gl.BindTexture(gl.TEXTURE_2D, uint32(PfSmellRedTextureID)) //

	// Write
	gl.DrawElements(gl.TRIANGLES, 6, gl.UNSIGNED_INT, gl.PtrOffset(0))

	// Update green smell
	// =======================================================
	gl.BindFramebuffer(gl.FRAMEBUFFER, FBSmellGreen)

	data = DataSmellGreen
	data.Enable()

	// Set uniforms
	data.Program.SetFloat("window_width", float32(Width))
	data.Program.SetFloat("window_height", float32(Height))

	// Bind textures
	data.Program.SetInt("PfGameTexture", int32(0))
	data.Program.SetInt("PfSmellTexture", int32(1))
	gl.ActiveTexture(gl.TEXTURE0 + 0)                      // Game state
	gl.BindTexture(gl.TEXTURE_2D, uint32(PfGameTextureID)) //

	gl.ActiveTexture(gl.TEXTURE0 + 1)                            // Smell (blurred composite of game state)
	gl.BindTexture(gl.TEXTURE_2D, uint32(PfSmellGreenTextureID)) //

	// Write
	gl.DrawElements(gl.TRIANGLES, 6, gl.UNSIGNED_INT, gl.PtrOffset(0))
}

func DrawDataset() {

	// Clear screen
	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	// Draw game
	// ==============================================================
	// Enable Blit program
	data := DataMix
	data.Enable()

	// Set uniforms
	data.Program.SetFloat("window_width", float32(Width))
	data.Program.SetFloat("window_height", float32(Height))
	data.Program.SetInt("MODE", ActionDrawMode)
	data.Program.SetFloat("A", ActionDrawA)

	// Bind textures
	data.Program.SetInt("PfGameTexture", int32(0))
	data.Program.SetInt("PfSmellRedTexture", int32(1))
	data.Program.SetInt("PfSmellGreenTexture", int32(2))
	gl.ActiveTexture(gl.TEXTURE0 + 0)                            // Game state
	gl.BindTexture(gl.TEXTURE_2D, uint32(PfGameTextureID))       //
	gl.ActiveTexture(gl.TEXTURE0 + 1)                            // Smell (blurred composite of game state)
	gl.BindTexture(gl.TEXTURE_2D, uint32(PfSmellRedTextureID))   //
	gl.ActiveTexture(gl.TEXTURE0 + 2)                            // Smell (blurred composite of game state)
	gl.BindTexture(gl.TEXTURE_2D, uint32(PfSmellGreenTextureID)) //

	// Draw Game state
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.ONE, gl.ZERO)
	gl.DrawElements(gl.TRIANGLES, 6, gl.UNSIGNED_INT, gl.PtrOffset(0))

	// [bug magnet]
	/*
		if ActionPrintSmell && ActionPrintGame {
			gl.BlendEquation(gl.MAX)
		}
	*/

	// Put buffer that we painted on (gl.BACK) on the foreground.
	window.SwapBuffers()

}
