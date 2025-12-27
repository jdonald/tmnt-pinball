package main

import (
	"math"

	"github.com/veandco/go-sdl2/sdl"
)

// drawCircle draws a filled circle using midpoint circle algorithm
func drawCircle(renderer *sdl.Renderer, centerX, centerY, radius int32) {
	for y := -radius; y <= radius; y++ {
		for x := -radius; x <= radius; x++ {
			if x*x+y*y <= radius*radius {
				renderer.DrawPoint(centerX+x, centerY+y)
			}
		}
	}
}

// drawText renders simple text (bitmap style)
// This is a simplified version - in production you'd use SDL_ttf
func drawText(renderer *sdl.Renderer, text string, x, y int, r, g, b uint8) {
	renderer.SetDrawColor(r, g, b, 255)

	charWidth := int32(8)
	charHeight := int32(12)

	for i, char := range text {
		charX := int32(x) + int32(i)*charWidth
		charY := int32(y)

		// Draw simple character representation
		switch char {
		case ' ':
			// Space - do nothing
		default:
			// Draw a simple rectangle for each character
			drawSimpleChar(renderer, char, charX, charY, charWidth, charHeight)
		}
	}
}

func drawSimpleChar(renderer *sdl.Renderer, char rune, x, y, width, height int32) {
	// This is a very simplified character renderer
	// For each character, we'll draw a pattern

	// Draw character outline
	renderer.DrawRect(&sdl.Rect{X: x, Y: y, W: width, H: height})

	// Add some vertical lines for readability
	renderer.DrawLine(x+width/3, y, x+width/3, y+height)
	renderer.DrawLine(x+2*width/3, y, x+2*width/3, y+height)
}

// drawLine draws a line between two points
func drawLine(renderer *sdl.Renderer, x1, y1, x2, y2 int32) {
	renderer.DrawLine(x1, y1, x2, y2)
}

// drawPolygon draws a polygon
func drawPolygon(renderer *sdl.Renderer, points []sdl.Point) {
	for i := 0; i < len(points); i++ {
		next := (i + 1) % len(points)
		renderer.DrawLine(points[i].X, points[i].Y, points[next].X, points[next].Y)
	}
}

// fillPolygon fills a polygon (simplified)
func fillPolygon(renderer *sdl.Renderer, points []sdl.Point) {
	if len(points) < 3 {
		return
	}

	// Find bounding box
	minX, maxX := points[0].X, points[0].X
	minY, maxY := points[0].Y, points[0].Y

	for _, p := range points {
		if p.X < minX {
			minX = p.X
		}
		if p.X > maxX {
			maxX = p.X
		}
		if p.Y < minY {
			minY = p.Y
		}
		if p.Y > maxY {
			maxY = p.Y
		}
	}

	// Scan line fill
	for y := minY; y <= maxY; y++ {
		for x := minX; x <= maxX; x++ {
			if pointInPolygon(x, y, points) {
				renderer.DrawPoint(x, y)
			}
		}
	}
}

// pointInPolygon checks if a point is inside a polygon using ray casting
func pointInPolygon(x, y int32, points []sdl.Point) bool {
	n := len(points)
	inside := false

	p1 := points[0]
	for i := 1; i <= n; i++ {
		p2 := points[i%n]

		if y > min(p1.Y, p2.Y) {
			if y <= max(p1.Y, p2.Y) {
				if x <= max(p1.X, p2.X) {
					if p1.Y != p2.Y {
						xinters := int32((int64(y-p1.Y))*int64(p2.X-p1.X)/int64(p2.Y-p1.Y) + int64(p1.X))
						if p1.X == p2.X || x <= xinters {
							inside = !inside
						}
					}
				}
			}
		}
		p1 = p2
	}

	return inside
}

func min(a, b int32) int32 {
	if a < b {
		return a
	}
	return b
}

func max(a, b int32) int32 {
	if a > b {
		return a
	}
	return b
}

// drawRamp draws a pinball ramp
func drawRamp(renderer *sdl.Renderer, x1, y1, x2, y2 int32, width int32) {
	// Calculate perpendicular vector
	dx := float64(x2 - x1)
	dy := float64(y2 - y1)
	length := math.Sqrt(dx*dx + dy*dy)

	if length == 0 {
		return
	}

	perpX := -dy / length * float64(width) / 2
	perpY := dx / length * float64(width) / 2

	points := []sdl.Point{
		{X: x1 + int32(perpX), Y: y1 + int32(perpY)},
		{X: x1 - int32(perpX), Y: y1 - int32(perpY)},
		{X: x2 - int32(perpX), Y: y2 - int32(perpY)},
		{X: x2 + int32(perpX), Y: y2 + int32(perpY)},
	}

	fillPolygon(renderer, points)
}
