package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/go-gl/glfw/v3.2/glfw"
)

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
				ActionReset = true
			}
		case "p":
			if action == Down {
				ActionPrtsc = true
			}
		case "s":
			if action == Down {
				ActionPrintSmell = !ActionPrintSmell
			}
		}

	})
}

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
