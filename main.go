package main

import (
	"fmt"
	"os"
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

const (
	WindowWidth  = 800
	WindowHeight = 1200
	FPS          = 60
)

func main() {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize SDL: %s\n", err)
		os.Exit(1)
	}
	defer sdl.Quit()

	window, err := sdl.CreateWindow(
		"TMNT Pinball",
		sdl.WINDOWPOS_UNDEFINED,
		sdl.WINDOWPOS_UNDEFINED,
		WindowWidth,
		WindowHeight,
		sdl.WINDOW_SHOWN,
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create window: %s\n", err)
		os.Exit(1)
	}
	defer window.Destroy()

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create renderer: %s\n", err)
		os.Exit(1)
	}
	defer renderer.Destroy()

	// Initialize game controller support
	// Note: SDL2 has built-in controller mappings for common controllers
	// Custom mappings from gamecontrollerdb.txt can be loaded manually if needed

	game := NewGame(renderer)
	inputManager := NewInputManager()

	running := true
	frameDelay := time.Second / FPS

	for running {
		frameStart := time.Now()

		// Handle events
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch e := event.(type) {
			case *sdl.QuitEvent:
				running = false
			case *sdl.KeyboardEvent:
				inputManager.HandleKeyboard(e)
			case *sdl.ControllerButtonEvent:
				inputManager.HandleControllerButton(e)
			case *sdl.ControllerAxisEvent:
				inputManager.HandleControllerAxis(e)
			case *sdl.ControllerDeviceEvent:
				inputManager.HandleControllerDevice(e)
			}
		}

		// Update game state
		game.Update(inputManager)

		// Render
		renderer.SetDrawColor(0, 0, 0, 255)
		renderer.Clear()
		game.Render(renderer)
		renderer.Present()

		// Frame rate limiting
		frameTime := time.Since(frameStart)
		if frameTime < frameDelay {
			time.Sleep(frameDelay - frameTime)
		}
	}
}
