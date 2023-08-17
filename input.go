package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/go-gl/glfw/v3.2/glfw"
)

func SetKeyHandling(window *glfw.Window) {
	window.SetKeyCallback(func(_ *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
		// Alias for readability
		Down := glfw.Press
		Up := glfw.Release
		Repeat := glfw.Repeat
		_ = Up

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
		translate_step_size := 0.01
		zoom_step_size := 0.01

		switch char {
		case "65": // space
			if action == Down || action == Repeat {
				ActionReset = true
			}
		case "110": // home
			if action == Down || action == Repeat {
				ZOOM -= zoom_step_size
			}
		case "115": // end
			if action == Down || action == Repeat {
				ZOOM += zoom_step_size
			}
		case "116": // arrow down
			if action == Down || action == Repeat {
				Y_TRANSLATE -= translate_step_size
				// ActionDrawA -= 0.1
				// if ActionDrawA < 0.0 {
				// 	ActionDrawA = 0.0
				// }
			}
		case "111": // arrow up
			if action == Down || action == Repeat {
				Y_TRANSLATE += translate_step_size
				// ActionDrawA += 0.1
				// if ActionDrawA > 1.0 {
				// 	ActionDrawA = 1.0
				// }
			}
		case "113": // arrow left
			if action == Down || action == Repeat {
				X_TRANSLATE -= translate_step_size
				// ActorDotRadius -= 0.002
				// if ActorDotRadius < 0.0 {
				// 	ActorDotRadius = 0.0
				// }
			}
		case "114": // arrow right
			if action == Down || action == Repeat {
				X_TRANSLATE += translate_step_size
				// ActorDotRadius += 0.002
				// if ActorDotRadius > 2.0 {
				// 	ActorDotRadius = 2.0
				// }
			}
		case "p":
			if action == Down {
				ActionPrtsc = true
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
