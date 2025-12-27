package main

import (
	"math"

	"github.com/veandco/go-sdl2/sdl"
)

const (
	Gravity       = 0.3
	BallRadius    = 8.0
	LaunchPower   = -25.0
	Friction      = 0.99
	WallBounce    = 0.8
	FlipperBounce = 1.2
)

type Ball struct {
	X, Y         float64
	VelX, VelY   float64
	Radius       float64
	Active       bool
	InLauncher   bool
}

type PhysicsEngine struct {
	ball         *Ball
	extraBalls   []*Ball
	playWidth    float64
	playHeight   float64
	launcherX    float64
	launcherY    float64
}

func NewPhysicsEngine() *PhysicsEngine {
	pe := &PhysicsEngine{
		playWidth:  700,
		playHeight: 1100,
		launcherX:  650,
		launcherY:  1000,
	}
	pe.ball = &Ball{
		X:          pe.launcherX,
		Y:          pe.launcherY,
		Radius:     BallRadius,
		Active:     true,
		InLauncher: true,
	}
	pe.extraBalls = make([]*Ball, 0)
	return pe
}

func (pe *PhysicsEngine) Update(flippers *Flippers) {
	pe.updateBall(pe.ball, flippers)
	for _, ball := range pe.extraBalls {
		if ball.Active {
			pe.updateBall(ball, flippers)
		}
	}
}

func (pe *PhysicsEngine) updateBall(ball *Ball, flippers *Flippers) {
	if !ball.Active || ball.InLauncher {
		return
	}

	// Apply gravity
	ball.VelY += Gravity

	// Apply friction
	ball.VelX *= Friction
	ball.VelY *= Friction

	// Update position
	ball.X += ball.VelX
	ball.Y += ball.VelY

	// Wall collisions
	pe.handleWallCollisions(ball)

	// Flipper collisions
	pe.handleFlipperCollisions(ball, flippers)
}

func (pe *PhysicsEngine) handleWallCollisions(ball *Ball) {
	// Left wall
	if ball.X-ball.Radius < 50 {
		ball.X = 50 + ball.Radius
		ball.VelX = -ball.VelX * WallBounce
	}

	// Right wall
	if ball.X+ball.Radius > 50+pe.playWidth {
		ball.X = 50 + pe.playWidth - ball.Radius
		ball.VelX = -ball.VelX * WallBounce
	}

	// Top wall
	if ball.Y-ball.Radius < 50 {
		ball.Y = 50 + ball.Radius
		ball.VelY = -ball.VelY * WallBounce
	}

	// Bottom (ball lost if it goes past flippers)
	if ball.Y > 50+pe.playHeight {
		ball.Active = false
	}
}

func (pe *PhysicsEngine) handleFlipperCollisions(ball *Ball, flippers *Flippers) {
	// Left flipper
	if pe.checkFlipperCollision(ball, flippers.LeftX, flippers.LeftY, flippers.LeftAngle) {
		if flippers.LeftActive {
			// Strong upward bounce when flipper is active
			angle := flippers.LeftAngle - math.Pi/4
			speed := 20.0
			ball.VelX = math.Cos(angle) * speed
			ball.VelY = math.Sin(angle) * speed
		} else {
			ball.VelY = -ball.VelY * FlipperBounce
		}
	}

	// Right flipper
	if pe.checkFlipperCollision(ball, flippers.RightX, flippers.RightY, flippers.RightAngle) {
		if flippers.RightActive {
			// Strong upward bounce when flipper is active
			angle := flippers.RightAngle - math.Pi/4
			speed := 20.0
			ball.VelX = math.Cos(angle) * speed
			ball.VelY = math.Sin(angle) * speed
		} else {
			ball.VelY = -ball.VelY * FlipperBounce
		}
	}
}

func (pe *PhysicsEngine) checkFlipperCollision(ball *Ball, flipperX, flipperY, angle float64) bool {
	// Simplified collision check with flipper rectangle
	flipperLength := 80.0
	flipperWidth := 15.0

	// Calculate flipper endpoints
	endX := flipperX + math.Cos(angle)*flipperLength
	endY := flipperY + math.Sin(angle)*flipperLength

	// Distance from ball to line segment
	dist := pe.pointToLineDistance(ball.X, ball.Y, flipperX, flipperY, endX, endY)

	return dist < ball.Radius+flipperWidth/2
}

func (pe *PhysicsEngine) pointToLineDistance(px, py, x1, y1, x2, y2 float64) float64 {
	dx := x2 - x1
	dy := y2 - y1
	lenSq := dx*dx + dy*dy

	if lenSq == 0 {
		return math.Sqrt((px-x1)*(px-x1) + (py-y1)*(py-y1))
	}

	t := ((px-x1)*dx + (py-y1)*dy) / lenSq
	t = math.Max(0, math.Min(1, t))

	projX := x1 + t*dx
	projY := y1 + t*dy

	return math.Sqrt((px-projX)*(px-projX) + (py-projY)*(py-projY))
}

func (pe *PhysicsEngine) LaunchBall() {
	if pe.ball.InLauncher {
		pe.ball.InLauncher = false
		pe.ball.VelY = LaunchPower
		pe.ball.VelX = -3 // Slight angle
	}
}

func (pe *PhysicsEngine) ResetBall() {
	pe.ball.X = pe.launcherX
	pe.ball.Y = pe.launcherY
	pe.ball.VelX = 0
	pe.ball.VelY = 0
	pe.ball.Active = true
	pe.ball.InLauncher = true

	// Remove extra balls
	pe.extraBalls = make([]*Ball, 0)
}

func (pe *PhysicsEngine) IsBallLost() bool {
	if !pe.ball.Active && len(pe.extraBalls) == 0 {
		return true
	}
	// Check if all balls are inactive
	allInactive := !pe.ball.Active
	if allInactive {
		for _, ball := range pe.extraBalls {
			if ball.Active {
				allInactive = false
				break
			}
		}
	}
	return allInactive
}

func (pe *PhysicsEngine) AddBalls(count int) {
	for i := 0; i < count; i++ {
		newBall := &Ball{
			X:      400 + float64(i*50),
			Y:      300,
			VelX:   float64(i-count/2) * 3,
			VelY:   5,
			Radius: BallRadius,
			Active: true,
		}
		pe.extraBalls = append(pe.extraBalls, newBall)
	}
}

func (pe *PhysicsEngine) ActiveBalls() int {
	count := 0
	if pe.ball.Active {
		count++
	}
	for _, ball := range pe.extraBalls {
		if ball.Active {
			count++
		}
	}
	return count
}

func (pe *PhysicsEngine) PlayBumperSound() {
	// Sound will be played here when audio is implemented
}

func (pe *PhysicsEngine) Render(renderer *sdl.Renderer) {
	// Render main ball
	if pe.ball.Active || pe.ball.InLauncher {
		renderer.SetDrawColor(255, 255, 255, 255)
		drawCircle(renderer, int32(pe.ball.X), int32(pe.ball.Y), int32(pe.ball.Radius))

		// Add shine effect
		renderer.SetDrawColor(200, 200, 255, 255)
		drawCircle(renderer, int32(pe.ball.X-2), int32(pe.ball.Y-2), int32(pe.ball.Radius/3))
	}

	// Render extra balls (multiball)
	for _, ball := range pe.extraBalls {
		if ball.Active {
			renderer.SetDrawColor(255, 255, 255, 255)
			drawCircle(renderer, int32(ball.X), int32(ball.Y), int32(ball.Radius))

			renderer.SetDrawColor(200, 200, 255, 255)
			drawCircle(renderer, int32(ball.X-2), int32(ball.Y-2), int32(ball.Radius/3))
		}
	}

	// Draw launcher
	renderer.SetDrawColor(150, 150, 150, 255)
	renderer.FillRect(&sdl.Rect{X: int32(pe.launcherX) - 15, Y: int32(pe.launcherY) - 100, W: 30, H: 100})
}
