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
)

const (
	GL_DEFAULT_FBO = 0

	DRAW_MODE_ADD   = 1 // Just add both textures together
	DRAW_MODE_MERGE = 2 // uniformf A is used as mix, 0=game, 1=smell
	DRAW_MODE_SMELL = 3 // uniformf A is used as mix, 0=red, 0.5=both, 1.0=green

	StartImageSrc = "assets/img/image.png" // Needs to be a png
)

var (
	Width  = Config.Screen.Width  // Width of the main window, needs to be the same as the image that's loaded in StartImageSrc
	Height = Config.Screen.Height // Height of the main window, needs to be the same as the image that's loaded in StartImageSrc

	WindowTitle = "Test GL Application"
	window      = gogl.Init(WindowTitle, Width, Height) // Init Window, OpenGL

	// Helper variables
	tick          float32 = 0.0   // ticks up every game loop cycle
	fps           int     = 50    //
	delay_ms      int64   = 5     // handles frame rate
	record        bool    = false // whether to record the screen.
	record_length float32 = 0.0   //float32(2.0*(fps)) * 0.5 // After how many ticks to stop recording (and close the program)
	record_start  float32 = 0.0
	record_stop   float32 = 0.0

	// Shader programs
	datalist       = SetData()   // Init data
	DataGame       = datalist[0] // The actual gamestate.
	DataSmellRed   = datalist[1] // Just the red channel. Blurred to it extends the actual red regions
	DataSmellGreen = datalist[2] // Just the green channel. Blurred to it extends the actual green regions
	DataMix        = datalist[3] // Used to display a mix of the framebuffers on the screen
	DataSwapGame   = datalist[4] // Used to display a mix of the framebuffers on the screen

	// Framebuffers and their attached textures
	PfGameTextureID       gogl.TextureID // previous frame, the backbuffer gets bound to this texture after every draw.
	PfGameTextureWriteID  gogl.TextureID // [test]
	PfSmellRedTextureID   gogl.TextureID //
	PfSmellGreenTextureID gogl.TextureID //
	FBGame                uint32         // can be used as a draw target when we don't want to draw to the screen
	FBGameWrite           uint32         // [test]
	FBSmellRed            uint32         // can be used as a draw target when we don't want to draw to the screen
	FBSmellGreen          uint32         // can be used as a draw target when we don't want to draw to the screen

	// Actions. These are influenced by keystrokes
	ActionReset              = false
	ActionPrtsc              = false
	ActionRecord             = false
	ActionPrintSmell         = true
	ActionPrintGame          = true
	ActionDrawMode   int32   = DRAW_MODE_ADD
	ActionDrawA      float32 = 0.0

	ZOOM              = 1.0
	X_TRANSLATE       = 0.0
	Y_TRANSLATE       = 0.0
	SHOW_HUD          = 1
	TIMESTAMP   int32 = 0

	KeyWActive bool = false
	KeyAActive bool = false
	KeyRActive bool = false
	KeySActive bool = false

	// Actors
	ActorDotPos    [2]float32 = [2]float32{0.5, 0.5}
	ActorDotRadius float32    = 0.002
)

func main() {
	// disable player
	if Config.Player.Enabled == false {
		ActorDotRadius = 0
	}

	// Get user input from commandline & say what each keypress should do
	SetKeyHandling(window)
	SetMouseHandling(window)
	ParseCommandlineArgs()

	// Create textures for the different frame buffers
	PfGameTextureWriteID = NewDefaultTexture() // [test]
	PfGameTextureID = NewDefaultTexture()

	PfSmellRedTextureID = NewDefaultTexture()
	PfSmellGreenTextureID = NewDefaultTexture()

	ResetFrame(PfGameTextureWriteID, StartImageSrc) // texture needs to be initialized before linking to FBO
	CreateFramebuffer(&FBGameWrite, PfGameTextureWriteID)

	// Create frame buffer for the game state
	ResetFrame(PfGameTextureID, StartImageSrc) // texture needs to be initialized before linking to FBO
	CreateFramebuffer(&FBGame, PfGameTextureID)

	// Create frame buffer for the red "smell" (used for green cells to know where the red cells are)
	ResetFrame(PfSmellRedTextureID, StartImageSrc) // texture needs to be initialized before linking to FBO
	CreateFramebuffer(&FBSmellRed, PfSmellRedTextureID)

	// Create frame buffer for the green "smell" (used for red cells to know where the green cells are)
	ResetFrame(PfSmellGreenTextureID, StartImageSrc) // texture needs to be initialized before linking to FBO
	CreateFramebuffer(&FBSmellGreen, PfSmellGreenTextureID)

	// Main loop
	for !window.ShouldClose() {

		// Housekeeping
		start := time.Now()
		tick += 1.0

		// Reset the textures that are bound to the frame buffers when we want to restart
		if ActionReset {
			ActionReset = false
			ResetFrame(PfGameTextureWriteID, StartImageSrc)
			ResetFrame(PfGameTextureID, StartImageSrc)
			ResetFrame(PfSmellRedTextureID, StartImageSrc)
			ResetFrame(PfSmellGreenTextureID, StartImageSrc)
		}

		// key actions
		if KeyAActive {
			ActorDotPos[0] -= 0.01
			if ActorDotPos[0] < 0.0 {
				ActorDotPos[0] += 1.0
			}
		}
		if KeySActive {
			ActorDotPos[0] += 0.01
			if ActorDotPos[0] > 1.0 {
				ActorDotPos[0] -= 1.0
			}
		}
		if KeyRActive {
			ActorDotPos[1] -= 0.01
			if ActorDotPos[1] < 0.0 {
				ActorDotPos[1] += 1.0
			}
		}
		if KeyWActive {
			ActorDotPos[1] += 0.01
			if ActorDotPos[1] > 1.0 {
				ActorDotPos[1] -= 1.0
			}
		}

		// get new timestamp
		TIMESTAMP = int32(makeTimestamp())

		// Each framebuffer gets updated in a recursive manner
		// (Each fbo is the input of a calculation and the output of that same calculation)
		UpdateGame()

		// What ends up on the screen is a combination of the framebuffers.
		// What is drawn is influenced by the draw mode.
		Draw()

		// Fetch window events.
		glfw.PollEvents()

		// Check if shaders need to be recompiled, and recompile them if so.
		// This allows us to change the shaders mid-run.
		gogl.HotloadShaders()

		// Sanity check
		if err := gl.GetError(); err != 0 {
			log.Println("test", err)
		}

		// Record output (prtsc and gif/mp4 recording)
		if ActionPrtsc {
			CreateImage(int(tick), RECORDING_PRTSC)
			ActionPrtsc = false
		}
		if ActionRecord {
			CreateImage(int(tick), RECORDING_GIF)
			if record_stop != 0.0 && tick > record_stop {
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
	// Set to overwrite whatever is on the texture (no mixing)
	gl.Disable(gl.BLEND)

	// Update game
	// ---------------------------------------------
	// Select game buffer/shader program
	//  gl.BindFramebuffer(gl.FRAMEBUFFER, FBGameWrite)
	gl.BindFramebuffer(gl.DRAW_FRAMEBUFFER, FBGameWrite) // [test]
	gl.BindFramebuffer(gl.READ_FRAMEBUFFER, FBGame)      // [test]
	data := DataGame
	data.Enable()

	// Bind correct textures
	data.Program.SetInt("PfGameTexture", int32(0))
	data.Program.SetInt("PfSmellRedTexture", int32(1))
	data.Program.SetInt("PfSmellGreenTexture", int32(2))
	gl.ActiveTexture(gl.TEXTURE0 + 0)                            // Game state
	gl.BindTexture(gl.TEXTURE_2D, uint32(PfGameTextureID))       //
	gl.ActiveTexture(gl.TEXTURE0 + 1)                            // Smell (blurred composite of game state)
	gl.BindTexture(gl.TEXTURE_2D, uint32(PfSmellRedTextureID))   //
	gl.ActiveTexture(gl.TEXTURE0 + 2)                            // Smell (blurred composite of game state)
	gl.BindTexture(gl.TEXTURE_2D, uint32(PfSmellGreenTextureID)) //

	// Set uniforms
	// fmt.Println(int32(now.Unix()%100) / 2)
	data.Program.SetFloat("window_width", float32(Width))
	data.Program.SetFloat("window_height", float32(Height))
	data.Program.SetFloatVector2("Actor1", &ActorDotPos)
	data.Program.SetFloat("Actor1Radius", ActorDotRadius)
	data.Program.SetInt("TIMESTAMP", int32(TIMESTAMP))

	// Iterate Game state
	gl.DrawElements(gl.TRIANGLES, 6, gl.UNSIGNED_INT, gl.PtrOffset(0))

	// copy write PBO to read FBO
	// ---------------------------------------------
	// Select game buffer/shader program
	gl.BindFramebuffer(gl.FRAMEBUFFER, FBGame)
	data = DataSwapGame
	data.Enable()

	// Bind correct textures
	data.Program.SetInt("PfGameTexture", int32(0))
	gl.ActiveTexture(gl.TEXTURE0 + 0)                           // Game state
	gl.BindTexture(gl.TEXTURE_2D, uint32(PfGameTextureWriteID)) //

	// Iterate Game state
	gl.DrawElements(gl.TRIANGLES, 6, gl.UNSIGNED_INT, gl.PtrOffset(0))

	// Update red smell
	// ---------------------------------------------
	gl.BindFramebuffer(gl.FRAMEBUFFER, FBSmellRed)
	data = DataSmellRed
	data.Enable()

	// Bind textures
	data.Program.SetInt("PfGameTexture", int32(0))
	data.Program.SetInt("PfSmellTexture", int32(1))
	gl.ActiveTexture(gl.TEXTURE0 + 0)                          // Game state
	gl.BindTexture(gl.TEXTURE_2D, uint32(PfGameTextureID))     //
	gl.ActiveTexture(gl.TEXTURE0 + 1)                          // Smell (blurred composite of game state)
	gl.BindTexture(gl.TEXTURE_2D, uint32(PfSmellRedTextureID)) //

	// Set uniforms
	data.Program.SetFloat("window_width", float32(Width))
	data.Program.SetFloat("window_height", float32(Height))

	// Iterate red smell
	gl.DrawElements(gl.TRIANGLES, 6, gl.UNSIGNED_INT, gl.PtrOffset(0))

	// Update green smell
	// ---------------------------------------------
	gl.BindFramebuffer(gl.FRAMEBUFFER, FBSmellGreen)
	data = DataSmellGreen
	data.Enable()

	// Bind textures
	data.Program.SetInt("PfGameTexture", int32(0))
	data.Program.SetInt("PfSmellTexture", int32(1))
	gl.ActiveTexture(gl.TEXTURE0 + 0)                            // Game state
	gl.BindTexture(gl.TEXTURE_2D, uint32(PfGameTextureID))       //
	gl.ActiveTexture(gl.TEXTURE0 + 1)                            // Smell (blurred composite of game state)
	gl.BindTexture(gl.TEXTURE_2D, uint32(PfSmellGreenTextureID)) //

	// Set uniforms
	data.Program.SetFloat("window_width", float32(Width))
	data.Program.SetFloat("window_height", float32(Height))

	// Iterate green smell
	gl.DrawElements(gl.TRIANGLES, 6, gl.UNSIGNED_INT, gl.PtrOffset(0))
}

func Draw() {

	// Clear screen
	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	// Set to overwrite whatever is on the texture (no mixing)
	gl.Disable(gl.BLEND)

	// Enable Mix program
	data := DataMix
	data.Enable()

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

	// Set uniforms
	data.Program.SetInt("TIMESTAMP", int32(TIMESTAMP))
	data.Program.SetInt("MODE", ActionDrawMode)
	data.Program.SetFloat("A", ActionDrawA)
	data.Program.SetFloat("window_width", float32(Width))
	data.Program.SetFloat("window_height", float32(Height))
	data.Program.SetFloat("ZOOM", float32(ZOOM))
	data.Program.SetFloat("Y_TRANSLATE", float32(Y_TRANSLATE))
	data.Program.SetFloat("X_TRANSLATE", float32(X_TRANSLATE))
	data.Program.SetInt("SHOW_HUD", int32(SHOW_HUD))

	// Draw Game state to screen
	gl.DrawElements(gl.TRIANGLES, 6, gl.UNSIGNED_INT, gl.PtrOffset(0))

	// Put buffer that we painted on (gl.BACK) on the foreground.
	window.SwapBuffers()
}
