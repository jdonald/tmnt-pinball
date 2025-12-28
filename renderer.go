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

// drawText renders bitmap-style text with a simple 5x7 pixel font
func drawText(renderer *sdl.Renderer, text string, x, y int, r, g, b uint8) {
	renderer.SetDrawColor(r, g, b, 255)

	charWidth := int32(6)  // 5 pixels + 1 spacing
	charHeight := int32(7)
	scale := int32(1)

	for i, char := range text {
		charX := int32(x) + int32(i)*charWidth*scale
		charY := int32(y)
		drawChar(renderer, char, charX, charY, scale)
	}
}

// drawChar renders a single character using a 5x7 pixel bitmap font
func drawChar(renderer *sdl.Renderer, char rune, x, y, scale int32) {
	// Get the bitmap for this character
	bitmap := getCharBitmap(char)

	// Draw the bitmap
	for row := 0; row < 7; row++ {
		for col := 0; col < 5; col++ {
			if bitmap[row]&(1<<(4-col)) != 0 {
				if scale == 1 {
					renderer.DrawPoint(x+int32(col), y+int32(row))
				} else {
					renderer.FillRect(&sdl.Rect{
						X: x + int32(col)*scale,
						Y: y + int32(row)*scale,
						W: scale,
						H: scale,
					})
				}
			}
		}
	}
}

// getCharBitmap returns a 5x7 bitmap for a character
func getCharBitmap(char rune) [7]uint8 {
	// Each row is 5 bits (represented as uint8)
	// Bit 4 (MSB of 5 bits) is leftmost pixel

	switch char {
	case 'A':
		return [7]uint8{0x0E, 0x11, 0x11, 0x1F, 0x11, 0x11, 0x11}
	case 'B':
		return [7]uint8{0x1E, 0x11, 0x11, 0x1E, 0x11, 0x11, 0x1E}
	case 'C':
		return [7]uint8{0x0E, 0x11, 0x10, 0x10, 0x10, 0x11, 0x0E}
	case 'D':
		return [7]uint8{0x1E, 0x11, 0x11, 0x11, 0x11, 0x11, 0x1E}
	case 'E':
		return [7]uint8{0x1F, 0x10, 0x10, 0x1E, 0x10, 0x10, 0x1F}
	case 'F':
		return [7]uint8{0x1F, 0x10, 0x10, 0x1E, 0x10, 0x10, 0x10}
	case 'G':
		return [7]uint8{0x0E, 0x11, 0x10, 0x17, 0x11, 0x11, 0x0F}
	case 'H':
		return [7]uint8{0x11, 0x11, 0x11, 0x1F, 0x11, 0x11, 0x11}
	case 'I':
		return [7]uint8{0x0E, 0x04, 0x04, 0x04, 0x04, 0x04, 0x0E}
	case 'J':
		return [7]uint8{0x01, 0x01, 0x01, 0x01, 0x11, 0x11, 0x0E}
	case 'K':
		return [7]uint8{0x11, 0x12, 0x14, 0x18, 0x14, 0x12, 0x11}
	case 'L':
		return [7]uint8{0x10, 0x10, 0x10, 0x10, 0x10, 0x10, 0x1F}
	case 'M':
		return [7]uint8{0x11, 0x1B, 0x15, 0x15, 0x11, 0x11, 0x11}
	case 'N':
		return [7]uint8{0x11, 0x19, 0x15, 0x13, 0x11, 0x11, 0x11}
	case 'O':
		return [7]uint8{0x0E, 0x11, 0x11, 0x11, 0x11, 0x11, 0x0E}
	case 'P':
		return [7]uint8{0x1E, 0x11, 0x11, 0x1E, 0x10, 0x10, 0x10}
	case 'Q':
		return [7]uint8{0x0E, 0x11, 0x11, 0x11, 0x15, 0x12, 0x0D}
	case 'R':
		return [7]uint8{0x1E, 0x11, 0x11, 0x1E, 0x14, 0x12, 0x11}
	case 'S':
		return [7]uint8{0x0E, 0x11, 0x10, 0x0E, 0x01, 0x11, 0x0E}
	case 'T':
		return [7]uint8{0x1F, 0x04, 0x04, 0x04, 0x04, 0x04, 0x04}
	case 'U':
		return [7]uint8{0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x0E}
	case 'V':
		return [7]uint8{0x11, 0x11, 0x11, 0x11, 0x11, 0x0A, 0x04}
	case 'W':
		return [7]uint8{0x11, 0x11, 0x11, 0x15, 0x15, 0x1B, 0x11}
	case 'X':
		return [7]uint8{0x11, 0x11, 0x0A, 0x04, 0x0A, 0x11, 0x11}
	case 'Y':
		return [7]uint8{0x11, 0x11, 0x0A, 0x04, 0x04, 0x04, 0x04}
	case 'Z':
		return [7]uint8{0x1F, 0x01, 0x02, 0x04, 0x08, 0x10, 0x1F}
	case '0':
		return [7]uint8{0x0E, 0x11, 0x13, 0x15, 0x19, 0x11, 0x0E}
	case '1':
		return [7]uint8{0x04, 0x0C, 0x04, 0x04, 0x04, 0x04, 0x0E}
	case '2':
		return [7]uint8{0x0E, 0x11, 0x01, 0x02, 0x04, 0x08, 0x1F}
	case '3':
		return [7]uint8{0x0E, 0x11, 0x01, 0x0E, 0x01, 0x11, 0x0E}
	case '4':
		return [7]uint8{0x02, 0x06, 0x0A, 0x12, 0x1F, 0x02, 0x02}
	case '5':
		return [7]uint8{0x1F, 0x10, 0x1E, 0x01, 0x01, 0x11, 0x0E}
	case '6':
		return [7]uint8{0x06, 0x08, 0x10, 0x1E, 0x11, 0x11, 0x0E}
	case '7':
		return [7]uint8{0x1F, 0x01, 0x02, 0x04, 0x08, 0x08, 0x08}
	case '8':
		return [7]uint8{0x0E, 0x11, 0x11, 0x0E, 0x11, 0x11, 0x0E}
	case '9':
		return [7]uint8{0x0E, 0x11, 0x11, 0x0F, 0x01, 0x02, 0x0C}
	case ':':
		return [7]uint8{0x00, 0x04, 0x00, 0x00, 0x00, 0x04, 0x00}
	case '!':
		return [7]uint8{0x04, 0x04, 0x04, 0x04, 0x04, 0x00, 0x04}
	case '?':
		return [7]uint8{0x0E, 0x11, 0x01, 0x02, 0x04, 0x00, 0x04}
	case '.':
		return [7]uint8{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x04}
	case ',':
		return [7]uint8{0x00, 0x00, 0x00, 0x00, 0x00, 0x04, 0x08}
	case '-':
		return [7]uint8{0x00, 0x00, 0x00, 0x0E, 0x00, 0x00, 0x00}
	case '+':
		return [7]uint8{0x00, 0x04, 0x04, 0x1F, 0x04, 0x04, 0x00}
	case '/':
		return [7]uint8{0x01, 0x01, 0x02, 0x04, 0x08, 0x10, 0x10}
	case '(':
		return [7]uint8{0x02, 0x04, 0x08, 0x08, 0x08, 0x04, 0x02}
	case ')':
		return [7]uint8{0x08, 0x04, 0x02, 0x02, 0x02, 0x04, 0x08}
	case ' ':
		return [7]uint8{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	default:
		// For lowercase, use uppercase bitmap
		if char >= 'a' && char <= 'z' {
			return getCharBitmap(rune(char - 32))
		}
		// Unknown character - draw a small rectangle
		return [7]uint8{0x1F, 0x11, 0x11, 0x11, 0x11, 0x11, 0x1F}
	}
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
