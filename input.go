package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/go-gl/glfw/v3.2/glfw"
)

var (
	CURSOR_XPOS_SCREEN float64
	CURSOR_YPOS_SCREEN float64

	DRAG_XPOS_ORIGIN float64
	DRAG_YPOS_ORIGIN float64

	DRAG_X float64 // difference in uv since last drag()
	DRAG_Y float64 // difference in uv since last drag()

	MOUSE_LB_PRESSED bool
)

const (
	MOUSE_LEFT  = 0
	MOUSE_RIGHT = 1

	MOUSE_BUTTON_DOWN = 1
	MOUSE_BUTTON_UP   = 0

	MIN_PAN_DISTANCE = 0.002
)

// reset cursor related state when the cursor exits the window
func ExitScreen() {
	fmt.Println("screen exit")
	MOUSE_LB_PRESSED = false
}

func SetMouseHandling(window *glfw.Window) {
	window.SetCursorPosCallback(func(_ *glfw.Window, xpos float64, ypos float64) {
		CURSOR_XPOS_SCREEN = xpos
		CURSOR_YPOS_SCREEN = ypos

		if MOUSE_LB_PRESSED == true {
			update_drag_delta()
			drag()
		}
	})

	window.SetMouseButtonCallback(func(_ *glfw.Window, button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
		fmt.Println("button, action, mods", button, action, mods)

		if button == MOUSE_LEFT {
			if action == MOUSE_BUTTON_DOWN {
				MOUSE_LB_PRESSED = true
				reset_drag_origin()
			}
			if action == MOUSE_BUTTON_UP {
				MOUSE_LB_PRESSED = false
			}
		}
	})

	window.SetScrollCallback(func(_ *glfw.Window, xoffset float64, yoffset float64) {
		if yoffset == -1.0 {
			zoom(ZOOM_OUT, ZOOM_STEP_SIZE*4.0*ZOOM)
		}
		if yoffset == 1.0 {
			zoom(ZOOM_IN, ZOOM_STEP_SIZE*4.0*ZOOM)
		}
	})

	window.SetCursorEnterCallback(func(_ *glfw.Window, entered bool) {
		if entered == false {
			ExitScreen()
		} else {
			fmt.Println("screen enter")
		}
	})
}

func SetKeyHandling(window *glfw.Window) {
	window.SetKeyCallback(func(_ *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
		// Alias for readability
		Down := glfw.Press
		Up := glfw.Release
		Repeat := glfw.Repeat

		// Get key as defined by the locale (qwerty, dvorak, etc)
		char := glfw.GetKeyName(key, scancode)
		if char == "" {
			char = fmt.Sprint(scancode)
		}

		// Handle keystrokes by keyboard position (locale independent)
		switch key {
		case glfw.KeyA:
			if action == Down {
				KeyAActive = true
			} else if action == Up {
				KeyAActive = false
			}
		case glfw.KeyR:
			if action == Down {
				KeyRActive = true
			} else if action == Up {
				KeyRActive = false
			}

		case glfw.KeyS:
			if action == Down {
				KeySActive = true
			} else if action == Up {
				KeySActive = false
			}
		case glfw.KeyW:
			if action == Down {
				KeyWActive = true
			} else if action == Up {
				KeyWActive = false
			}
		}

		// Handle keystrokes
		translate_step_size := 0.02

		switch char {
		case "65": // space
			if action == Down || action == Repeat {
				ActionReset = true
			}
		// ZOOM
		// -----------------------------------------------
		case "110": // home
			if action == Down || action == Repeat {
				zoom(ZOOM_IN, ZOOM_STEP_SIZE)
			}
		case "115": // end
			if action == Down || action == Repeat {
				zoom(ZOOM_OUT, ZOOM_STEP_SIZE)
			}
		// PAN
		// -----------------------------------------------
		case "116": // arrow down
			if action == Down || action == Repeat {
				// Y_TRANSLATE -= max(0.002, ZOOM*translate_step_size)
				pan(0.0, -translate_step_size)
				if action == Repeat {
					pan(0.0, -translate_step_size)
				}
			}
		case "111": // arrow up
			if action == Down || action == Repeat {
				pan(0.0, translate_step_size)
				if action == Repeat {
					pan(0.0, translate_step_size)
				}
			}
		case "113": // arrow left
			if action == Down || action == Repeat {
				pan(-translate_step_size, 0.0)
				if action == Repeat {
					pan(-translate_step_size, 0.0)
				}
			}
		case "114": // arrow right
			if action == Down || action == Repeat {
				pan(translate_step_size, 0.0)
				if action == Repeat {
					pan(translate_step_size, 0.0)
				}
			}
		// TOGGLES
		// -----------------------------------------------
		case "h":
			if action == Down {
				SHOW_HUD = 1 - SHOW_HUD
			}
		// PRINT/RECORD
		// -----------------------------------------------
		case "p":
			if action == Down {
				ActionPrtsc = true
			}
		case "r":
			if action == Down {
				if ActionRecord {
					// end recording
					record_stop = tick
				} else {
					ActionRecord = true
					record_start = tick
					if record_length > 0 {
						record_stop = tick + record_length
					} else {
						record_stop = 0.0
					}
				}
				fmt.Println("record, record_start, record_stop", ActionRecord, record_start, record_stop)

			}
		case "s":
			if action == Down || action == Repeat {
				ActionPrintSmell = !ActionPrintSmell
			}
		case "g":
			if action == Down || action == Repeat {
				ActionPrintGame = !ActionPrintGame
			}
		case "0":
			digit, _ := strconv.Atoi(char)
			ActionDrawMode = int32(digit)
		case "1":
			digit, _ := strconv.Atoi(char)
			ActionDrawMode = int32(digit)
		case "2":
			digit, _ := strconv.Atoi(char)
			ActionDrawMode = int32(digit)
		case "3":
			digit, _ := strconv.Atoi(char)
			ActionDrawMode = int32(digit)
		case "4":
			digit, _ := strconv.Atoi(char)
			ActionDrawMode = int32(digit)
		case "5":
			digit, _ := strconv.Atoi(char)
			ActionDrawMode = int32(digit)
		case "6":
			digit, _ := strconv.Atoi(char)
			ActionDrawMode = int32(digit)
		case "7":
			digit, _ := strconv.Atoi(char)
			ActionDrawMode = int32(digit)
		case "8":
			digit, _ := strconv.Atoi(char)
			ActionDrawMode = int32(digit)
		case "9":
			digit, _ := strconv.Atoi(char)
			ActionDrawMode = int32(digit)

		default:
			// get keychars if key is repeated and not matched
			if action == Repeat {
				//fmt.Println("action, char, key, scancode:", action, char, key, scancode)
			}
		}

		fmt.Println("action, char, key, scancode:", action, char, key, scancode)

	})
}

// Applies commandline args: --fps <N>, --record <N>
func ParseCommandlineArgs() {

	for i := range os.Args {

		// FPS
		// -----------------------------------------------------------
		// Apply commandline choice for fps, if present.
		// Note that for gif recording, 50 fps is the max.

		if os.Args[i] == "--fps" {
			// check if fps value has been passed directly after --fps
			if i+1 < len(os.Args) {
				choice, err := strconv.Atoi(os.Args[i+1])
				if err == nil {
					fps = choice
					delay_ms = int64(1000 / choice)
				} else {
					fmt.Println("ERROR: Could not parse input after --fps as an int.")
				}
			}
		}
	}

	for i := range os.Args {

		// Record
		// -----------------------------------------------------------
		// Apply commandline choice for recording settings, if present

		if os.Args[i] == "--record" {
			// Enable recording
			ActionRecord = true
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
