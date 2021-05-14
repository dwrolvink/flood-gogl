package main

import (
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
		}
	})
}
