package main

import (
	"fmt"
	"math"
	"math/rand"

	"github.com/veandco/go-sdl2/sdl"
)

type GameState int

const (
	StateMenu GameState = iota
	StateTurtleSelection
	StatePlaying
	StateMultiball
	StatePizzaTime
	StateGameOver
)

type Turtle int

const (
	Leonardo Turtle = iota
	Donatello
	Raphael
	Michelangelo
)

var TurtleNames = []string{"Leonardo", "Donatello", "Raphael", "Michelangelo"}
var TurtleColors = []sdl.Color{
	{0, 0, 255, 255},   // Leonardo - Blue
	{128, 0, 128, 255}, // Donatello - Purple
	{255, 0, 0, 255},   // Raphael - Red
	{255, 165, 0, 255}, // Michelangelo - Orange
}

type Game struct {
	state           GameState
	selectedTurtle  Turtle
	score           int
	balls           int
	physics         *PhysicsEngine
	flippers        *Flippers
	bumpers         []Bumper
	targets         []Target
	pizzaCount      int
	pizzaMode       bool
	pizzaModeTimer  int
	multiballActive bool
	multiballCount  int
	comboProgress   int
	episodeMode     bool
	renderer        *sdl.Renderer
}

type Bumper struct {
	X, Y   float64
	Radius float64
	Active bool
}

type Target struct {
	X, Y      float64
	Width     float64
	Height    float64
	Hit       bool
	Points    int
	TargetType string // "pizza", "foot", "episode"
}

func NewGame(renderer *sdl.Renderer) *Game {
	g := &Game{
		state:    StateMenu,
		balls:    3,
		renderer: renderer,
		physics:  NewPhysicsEngine(),
	}
	g.setupPlayfield()
	return g
}

func (g *Game) setupPlayfield() {
	// Create flippers
	g.flippers = NewFlippers()

	// Setup bumpers (like in the real table)
	g.bumpers = []Bumper{
		{X: 300, Y: 300, Radius: 30, Active: true},
		{X: 500, Y: 300, Radius: 30, Active: true},
		{X: 400, Y: 250, Radius: 30, Active: true},
	}

	// Setup targets
	g.targets = []Target{
		// Pizza targets
		{X: 100, Y: 200, Width: 40, Height: 20, Points: 1000, TargetType: "pizza"},
		{X: 150, Y: 200, Width: 40, Height: 20, Points: 1000, TargetType: "pizza"},
		{X: 200, Y: 200, Width: 40, Height: 20, Points: 1000, TargetType: "pizza"},
		// Foot clan targets (1-2-3 combo)
		{X: 600, Y: 400, Width: 30, Height: 30, Points: 500, TargetType: "foot1"},
		{X: 650, Y: 400, Width: 30, Height: 30, Points: 500, TargetType: "foot2"},
		{X: 700, Y: 400, Width: 30, Height: 30, Points: 500, TargetType: "foot3"},
		// Episode targets
		{X: 350, Y: 150, Width: 40, Height: 30, Points: 2000, TargetType: "episode"},
	}
}

func (g *Game) Update(input *InputManager) {
	switch g.state {
	case StateMenu:
		if input.StartPressed {
			g.state = StateTurtleSelection
			input.StartPressed = false
		}
	case StateTurtleSelection:
		g.handleTurtleSelection(input)
	case StatePlaying:
		g.updatePlaying(input)
	case StatePizzaTime:
		g.updatePizzaMode(input)
	case StateMultiball:
		g.updateMultiball(input)
	case StateGameOver:
		if input.StartPressed {
			g.resetGame()
			input.StartPressed = false
		}
	}
}

func (g *Game) handleTurtleSelection(input *InputManager) {
	// Navigate turtles with left/right or d-pad
	if input.LeftPressed {
		g.selectedTurtle = (g.selectedTurtle - 1 + 4) % 4
		input.LeftPressed = false
	}
	if input.RightPressed {
		g.selectedTurtle = (g.selectedTurtle + 1) % 4
		input.RightPressed = false
	}

	// Confirm selection with flipper button or A button
	if input.LeftFlipper || input.RightFlipper || input.StartPressed {
		g.state = StatePlaying
		g.physics.LaunchBall()
		input.LeftFlipper = false
		input.RightFlipper = false
		input.StartPressed = false
	}
}

func (g *Game) updatePlaying(input *InputManager) {
	// Handle flippers
	g.flippers.LeftActive = input.LeftFlipper
	g.flippers.RightActive = input.RightFlipper

	// Handle pizza button
	if input.PizzaButton && g.pizzaCount > 0 {
		g.activatePizzaMode()
		input.PizzaButton = false
	}

	// Launch ball
	if input.LaunchButton {
		g.physics.LaunchBall()
		input.LaunchButton = false
	}

	// Update physics
	g.physics.Update(g.flippers)

	// Check collisions
	g.checkCollisions()

	// Check if ball is lost
	if g.physics.IsBallLost() {
		g.balls--
		if g.balls <= 0 {
			g.state = StateGameOver
		} else {
			g.physics.ResetBall()
		}
	}

	// Update pizza mode timer
	if g.pizzaMode {
		g.pizzaModeTimer--
		if g.pizzaModeTimer <= 0 {
			g.pizzaMode = false
		}
	}
}

func (g *Game) updatePizzaMode(input *InputManager) {
	g.pizzaModeTimer--
	if g.pizzaModeTimer <= 0 {
		g.pizzaMode = false
		g.state = StatePlaying
	}
	// Pizza mode gives double points and slow motion
	g.updatePlaying(input)
}

func (g *Game) updateMultiball(input *InputManager) {
	// Similar to playing but with multiple balls
	g.updatePlaying(input)
	// Check if multiball should end
	if g.physics.ActiveBalls() <= 1 {
		g.multiballActive = false
		g.state = StatePlaying
	}
}

func (g *Game) checkCollisions() {
	ball := g.physics.ball

	// Check bumper collisions
	for i := range g.bumpers {
		bumper := &g.bumpers[i]
		dx := ball.X - bumper.X
		dy := ball.Y - bumper.Y
		dist := math.Sqrt(dx*dx + dy*dy)

		if dist < bumper.Radius+ball.Radius {
			// Bounce ball
			angle := math.Atan2(dy, dx)
			ball.VelX = math.Cos(angle) * 15
			ball.VelY = math.Sin(angle) * 15
			g.score += 100
			g.physics.PlayBumperSound()
		}
	}

	// Check target collisions
	for i := range g.targets {
		target := &g.targets[i]
		if !target.Hit &&
			ball.X > target.X && ball.X < target.X+target.Width &&
			ball.Y > target.Y && ball.Y < target.Y+target.Height {
			target.Hit = true
			g.score += target.Points
			g.handleTargetHit(target)
		}
	}
}

func (g *Game) handleTargetHit(target *Target) {
	switch target.TargetType {
	case "pizza":
		g.pizzaCount++
	case "foot1":
		if g.comboProgress == 0 {
			g.comboProgress = 1
		}
	case "foot2":
		if g.comboProgress == 1 {
			g.comboProgress = 2
		}
	case "foot3":
		if g.comboProgress == 2 {
			// Complete foot combo!
			g.score += 5000
			g.comboProgress = 0
			g.activateMultiball()
		}
	case "episode":
		g.activateEpisodeMode()
	}
}

func (g *Game) activatePizzaMode() {
	g.pizzaMode = true
	g.pizzaModeTimer = 300 // 5 seconds at 60 FPS
	g.pizzaCount--
	g.state = StatePizzaTime
}

func (g *Game) activateMultiball() {
	g.multiballActive = true
	g.multiballCount = 3
	g.state = StateMultiball
	g.physics.AddBalls(2) // Add 2 more balls
}

func (g *Game) activateEpisodeMode() {
	g.episodeMode = true
	g.score += 3000
}

func (g *Game) resetGame() {
	g.score = 0
	g.balls = 3
	g.pizzaCount = 0
	g.pizzaMode = false
	g.multiballActive = false
	g.comboProgress = 0
	g.episodeMode = false
	g.state = StateMenu
	g.physics.ResetBall()

	// Reset targets
	for i := range g.targets {
		g.targets[i].Hit = false
	}
}

func (g *Game) Render(renderer *sdl.Renderer) {
	switch g.state {
	case StateMenu:
		g.renderMenu(renderer)
	case StateTurtleSelection:
		g.renderTurtleSelection(renderer)
	case StatePlaying, StatePizzaTime, StateMultiball:
		g.renderPlayfield(renderer)
	case StateGameOver:
		g.renderGameOver(renderer)
	}
}

func (g *Game) renderMenu(renderer *sdl.Renderer) {
	// Green background
	renderer.SetDrawColor(0, 100, 0, 255)
	renderer.Clear()

	// Title
	drawText(renderer, "TEENAGE MUTANT NINJA TURTLES", 100, 300, 255, 255, 255)
	drawText(renderer, "PINBALL", 320, 350, 255, 255, 0)
	drawText(renderer, "Press START or SPACE to begin", 200, 600, 255, 255, 255)
}

func (g *Game) renderTurtleSelection(renderer *sdl.Renderer) {
	renderer.SetDrawColor(20, 20, 40, 255)
	renderer.Clear()

	drawText(renderer, "CHOOSE YOUR TURTLE", 220, 100, 255, 255, 0)

	// Draw all turtles
	for i := 0; i < 4; i++ {
		x := int32(150 + i*150)
		y := int32(400)
		color := TurtleColors[i]

		// Highlight selected turtle
		if Turtle(i) == g.selectedTurtle {
			renderer.SetDrawColor(255, 255, 0, 255)
			renderer.FillRect(&sdl.Rect{X: x - 10, Y: y - 10, W: 80, H: 120})
		}

		// Draw turtle representation (circle)
		renderer.SetDrawColor(0, 200, 0, 255) // Green body
		drawCircle(renderer, x+25, y+25, 30)

		// Draw bandana color
		renderer.SetDrawColor(color.R, color.G, color.B, 255)
		renderer.FillRect(&sdl.Rect{X: x, Y: y - 5, W: 60, H: 15})

		// Draw name
		drawText(renderer, TurtleNames[i], int(x)-20, int(y)+80, 255, 255, 255)
	}

	drawText(renderer, "Use LEFT/RIGHT to select, FLIPPER to confirm", 120, 700, 200, 200, 200)
}

func (g *Game) renderPlayfield(renderer *sdl.Renderer) {
	// Dark playfield background
	renderer.SetDrawColor(20, 20, 30, 255)
	renderer.Clear()

	// Draw playfield border
	renderer.SetDrawColor(100, 100, 100, 255)
	renderer.DrawRect(&sdl.Rect{X: 50, Y: 50, W: 700, H: 1100})

	// Draw bumpers
	for _, bumper := range g.bumpers {
		renderer.SetDrawColor(255, 0, 0, 255)
		drawCircle(renderer, int32(bumper.X), int32(bumper.Y), int32(bumper.Radius))
	}

	// Draw targets
	for _, target := range g.targets {
		if !target.Hit {
			switch target.TargetType {
			case "pizza":
				renderer.SetDrawColor(255, 165, 0, 255) // Orange
			case "foot1", "foot2", "foot3":
				renderer.SetDrawColor(128, 0, 128, 255) // Purple
			case "episode":
				renderer.SetDrawColor(0, 255, 0, 255) // Green
			}
			renderer.FillRect(&sdl.Rect{
				X: int32(target.X),
				Y: int32(target.Y),
				W: int32(target.Width),
				H: int32(target.Height),
			})
		}
	}

	// Draw flippers
	g.flippers.Render(renderer)

	// Draw ball(s)
	g.physics.Render(renderer)

	// Draw HUD
	g.renderHUD(renderer)
}

func (g *Game) renderHUD(renderer *sdl.Renderer) {
	// Turtle indicator
	color := TurtleColors[g.selectedTurtle]
	renderer.SetDrawColor(color.R, color.G, color.B, 255)
	renderer.FillRect(&sdl.Rect{X: 20, Y: 20, W: 30, H: 30})
	drawText(renderer, TurtleNames[g.selectedTurtle], 60, 25, 255, 255, 255)

	// Score
	drawText(renderer, fmt.Sprintf("SCORE: %d", g.score), 20, 60, 255, 255, 0)

	// Balls remaining
	drawText(renderer, fmt.Sprintf("BALLS: %d", g.balls), 20, 90, 255, 255, 255)

	// Pizza count
	if g.pizzaCount > 0 {
		drawText(renderer, fmt.Sprintf("PIZZA: %d", g.pizzaCount), 20, 120, 255, 165, 0)
		drawText(renderer, "Press MIDDLE BUTTON to EAT PIZZA!", 20, 145, 255, 200, 0)
	}

	// Pizza mode indicator
	if g.pizzaMode {
		drawText(renderer, "PIZZA POWER!", 300, 100, 255, 255, 0)
	}

	// Multiball indicator
	if g.multiballActive {
		drawText(renderer, "MULTIBALL!", 300, 150, 255, 0, 0)
	}

	// Foot combo progress
	if g.comboProgress > 0 {
		drawText(renderer, fmt.Sprintf("FOOT COMBO: %d/3", g.comboProgress), 20, 180, 128, 0, 128)
	}
}

func (g *Game) renderGameOver(renderer *sdl.Renderer) {
	renderer.SetDrawColor(0, 0, 0, 255)
	renderer.Clear()

	drawText(renderer, "GAME OVER", 300, 400, 255, 0, 0)
	drawText(renderer, fmt.Sprintf("Final Score: %d", g.score), 280, 450, 255, 255, 255)
	drawText(renderer, "Press START to play again", 240, 550, 200, 200, 200)
}
