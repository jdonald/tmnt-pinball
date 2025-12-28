package main

import (
	"math"

	"github.com/veandco/go-sdl2/sdl"
)

const (
	FlipperLength     = 120.0 // Increased from 80 to close center gap
	FlipperWidth      = 15.0
	FlipperRestAngle  = 0.4  // ~23 degrees down
	FlipperActiveAngle = -0.6 // ~34 degrees up
	FlipperSpeed      = 0.3
)

type Flippers struct {
	LeftX, LeftY   float64
	RightX, RightY float64
	LeftAngle      float64
	RightAngle     float64
	LeftActive     bool
	RightActive    bool
}

func NewFlippers() *Flippers {
	return &Flippers{
		LeftX:      110, // Moved closer to left wall (was 188)
		LeftY:      750,
		RightX:     450, // Moved closer to right (was 413), but left of launcher at 488
		RightY:     750,
		LeftAngle:  FlipperRestAngle,
		RightAngle: math.Pi - FlipperRestAngle,
	}
}

func (f *Flippers) Update() {
	// Animate left flipper
	if f.LeftActive {
		if f.LeftAngle > FlipperActiveAngle {
			f.LeftAngle -= FlipperSpeed
			if f.LeftAngle < FlipperActiveAngle {
				f.LeftAngle = FlipperActiveAngle
			}
		}
	} else {
		if f.LeftAngle < FlipperRestAngle {
			f.LeftAngle += FlipperSpeed
			if f.LeftAngle > FlipperRestAngle {
				f.LeftAngle = FlipperRestAngle
			}
		}
	}

	// Animate right flipper
	if f.RightActive {
		targetAngle := math.Pi - FlipperActiveAngle
		if f.RightAngle < targetAngle {
			f.RightAngle += FlipperSpeed
			if f.RightAngle > targetAngle {
				f.RightAngle = targetAngle
			}
		}
	} else {
		targetAngle := math.Pi - FlipperRestAngle
		if f.RightAngle > targetAngle {
			f.RightAngle -= FlipperSpeed
			if f.RightAngle < targetAngle {
				f.RightAngle = targetAngle
			}
		}
	}
}

func (f *Flippers) Render(renderer *sdl.Renderer) {
	f.Update()

	// Render left flipper
	f.renderFlipper(renderer, f.LeftX, f.LeftY, f.LeftAngle, f.LeftActive)

	// Render right flipper
	f.renderFlipper(renderer, f.RightX, f.RightY, f.RightAngle, f.RightActive)
}

func (f *Flippers) renderFlipper(renderer *sdl.Renderer, x, y, angle float64, active bool) {
	// Color changes when active
	if active {
		renderer.SetDrawColor(255, 255, 0, 255) // Yellow when active
	} else {
		renderer.SetDrawColor(180, 180, 180, 255) // Gray when idle
	}

	// Calculate flipper endpoints
	endX := x + math.Cos(angle)*FlipperLength
	endY := y + math.Sin(angle)*FlipperLength

	// Draw flipper as a thick line (polygon would be better)
	perpX := -math.Sin(angle) * FlipperWidth / 2
	perpY := math.Cos(angle) * FlipperWidth / 2

	points := []sdl.Point{
		{X: int32(x + perpX), Y: int32(y + perpY)},
		{X: int32(x - perpX), Y: int32(y - perpY)},
		{X: int32(endX - perpX), Y: int32(endY - perpY)},
		{X: int32(endX + perpX), Y: int32(endY + perpY)},
	}

	// Draw filled polygon (simplified as lines)
	for i := 0; i < len(points); i++ {
		next := (i + 1) % len(points)
		renderer.DrawLine(points[i].X, points[i].Y, points[next].X, points[next].Y)
	}

	// Fill with horizontal lines (simple fill)
	for dy := -FlipperWidth / 2; dy <= FlipperWidth/2; dy += 1 {
		startX := x + math.Cos(angle)*0 + perpX*dy/(FlipperWidth/2)
		startY := y + math.Sin(angle)*0 + perpY*dy/(FlipperWidth/2)
		endFillX := endX + perpX*dy/(FlipperWidth/2)
		endFillY := endY + perpY*dy/(FlipperWidth/2)
		renderer.DrawLine(int32(startX), int32(startY), int32(endFillX), int32(endFillY))
	}

	// Draw pivot point
	renderer.SetDrawColor(100, 100, 100, 255)
	drawCircle(renderer, int32(x), int32(y), 8)
}
