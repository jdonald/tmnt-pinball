package main

import (
	"github.com/veandco/go-sdl2/sdl"
)

type InputManager struct {
	LeftFlipper  bool
	RightFlipper bool
	PizzaButton  bool
	LaunchButton bool
	StartPressed bool
	LeftPressed  bool
	RightPressed bool
	controller   *sdl.GameController
}

func NewInputManager() *InputManager {
	im := &InputManager{}
	im.initController()
	return im
}

func (im *InputManager) initController() {
	// Try to open the first available controller
	for i := 0; i < sdl.NumJoysticks(); i++ {
		if sdl.IsGameController(i) {
			im.controller = sdl.GameControllerOpen(i)
			if im.controller != nil {
				break
			}
		}
	}
}

func (im *InputManager) HandleKeyboard(event *sdl.KeyboardEvent) {
	pressed := event.Type == sdl.KEYDOWN

	switch event.Keysym.Sym {
	case sdl.K_LEFT, sdl.K_a:
		if pressed {
			im.LeftPressed = true
		}
	case sdl.K_RIGHT, sdl.K_d:
		if pressed {
			im.RightPressed = true
		}
	case sdl.K_LSHIFT, sdl.K_z:
		im.LeftFlipper = pressed
	case sdl.K_RSHIFT, sdl.K_SLASH:
		im.RightFlipper = pressed
	case sdl.K_SPACE:
		im.PizzaButton = pressed
	case sdl.K_RETURN, sdl.K_1:
		if pressed {
			im.StartPressed = true
		}
	case sdl.K_DOWN, sdl.K_s:
		im.LaunchButton = pressed
	}
}

func (im *InputManager) HandleControllerButton(event *sdl.ControllerButtonEvent) {
	pressed := event.Type == sdl.CONTROLLERBUTTONDOWN

	switch event.Button {
	case sdl.CONTROLLER_BUTTON_A:
		if pressed {
			im.StartPressed = true
		}
	case sdl.CONTROLLER_BUTTON_START:
		if pressed {
			im.StartPressed = true
		}
	case sdl.CONTROLLER_BUTTON_LEFTSHOULDER:
		im.LeftFlipper = pressed
	case sdl.CONTROLLER_BUTTON_RIGHTSHOULDER:
		im.RightFlipper = pressed
	case sdl.CONTROLLER_BUTTON_X:
		im.PizzaButton = pressed
	case sdl.CONTROLLER_BUTTON_Y:
		im.LaunchButton = pressed
	case sdl.CONTROLLER_BUTTON_DPAD_LEFT:
		if pressed {
			im.LeftPressed = true
		}
	case sdl.CONTROLLER_BUTTON_DPAD_RIGHT:
		if pressed {
			im.RightPressed = true
		}
	case sdl.CONTROLLER_BUTTON_DPAD_DOWN:
		im.LaunchButton = pressed
	}
}

func (im *InputManager) HandleControllerAxis(event *sdl.ControllerAxisEvent) {
	// Handle analog triggers
	const triggerThreshold = 16000

	switch event.Axis {
	case sdl.CONTROLLER_AXIS_TRIGGERLEFT:
		im.LeftFlipper = event.Value > triggerThreshold
	case sdl.CONTROLLER_AXIS_TRIGGERRIGHT:
		im.RightFlipper = event.Value > triggerThreshold
	}
}

func (im *InputManager) HandleControllerDevice(event *sdl.ControllerDeviceEvent) {
	if event.Type == sdl.CONTROLLERDEVICEADDED {
		if im.controller == nil {
			im.controller = sdl.GameControllerOpen(int(event.Which))
		}
	} else if event.Type == sdl.CONTROLLERDEVICEREMOVED {
		if im.controller != nil {
			im.controller.Close()
			im.controller = nil
		}
	}
}

func (im *InputManager) Cleanup() {
	if im.controller != nil {
		im.controller.Close()
	}
}
