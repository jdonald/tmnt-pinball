# TMNT Pinball - Teenage Mutant Ninja Turtles Pinball Game

A pinball video game inspired by the Stern TMNT pinball machine, built in Go for macOS with game controller support.

## Features

### Game Modes

- **Turtle Selection**: Choose your favorite turtle (Leonardo, Donatello, Raphael, or Michelangelo)
- **Pizza Power Mode**: Collect pizza slices and activate pizza mode for double points and slow motion
- **Ninja Pizza Multiball**: Lock 3 balls to start multiball madness
- **Foot Clan Combo**: Complete the 1-2-3 combo to earn big points
- **Episode Mode**: Hit special targets to activate episode bonuses

### Ninja Instructions (from the real machine)

- **Choose Your Turtle**: At game start, push FLIPPER BUTTON to choose or plunge ball to SELECT
- **Ninja Pizza Multiball**: Lock 3 balls in the NINJA PIZZA PARLOR to start multiball
- **Turtle Power Multiball**: Advance TURTLE POWER to light RIGHT RAMP to start multiball
- **Training**: Shoot the ball BEHIND the UPPER FLIPPER to start TRAINING MISSIONS
- **Foot Combo**: Complete 1-2-3 WHEN LIT to collect TURTLE COMBO AWARDS
- **Episodes**: Shoot the Television Shot when flashing to start an EPISODE
- **Pizza Power**: USE BUTTON TO EAT PIZZA for special powers!

## Prerequisites

### macOS Requirements

1. **Homebrew** (if not installed):
   ```bash
   /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
   ```

2. **SDL2** libraries:
   ```bash
   brew install sdl2
   ```

3. **Go** (1.21 or later):
   ```bash
   brew install go
   ```

## Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/jdonald/tmnt-pinball.git
   cd tmnt-pinball
   ```

2. Download dependencies:
   ```bash
   go mod download
   ```

3. Build the game:
   ```bash
   go build -o tmnt-pinball
   ```

## Running the Game

```bash
./tmnt-pinball
```

Or run directly with Go:
```bash
go run .
```

## Controls

### Keyboard Controls

| Key | Action |
|-----|--------|
| **Left Shift** or **Z** | Left Flipper |
| **Right Shift** or **/** | Right Flipper |
| **Space** | Pizza Button (Eat Pizza) |
| **Down Arrow** or **S** | Launch Ball |
| **Left/Right Arrows** or **A/D** | Navigate Menus |
| **Enter** or **1** | Start Game / Confirm Selection |

### Game Controller Controls

| Button | Action |
|--------|--------|
| **Left Shoulder (L1/LB)** | Left Flipper |
| **Right Shoulder (R1/RB)** | Right Flipper |
| **Left Trigger (L2/LT)** | Left Flipper (analog) |
| **Right Trigger (R2/RT)** | Right Flipper (analog) |
| **X Button** | Pizza Button |
| **Y Button** | Launch Ball |
| **D-Pad Left/Right** | Navigate Menus |
| **D-Pad Down** | Launch Ball |
| **A Button** or **Start** | Start Game / Confirm |

## Gameplay Tips

1. **Choose Your Turtle Wisely**: Each turtle has a unique color scheme that affects your visual experience
2. **Collect Pizza**: Hit the orange pizza targets to collect pizza slices
3. **Pizza Power**: Press the pizza button when you have pizza to activate special mode
4. **Multiball Strategy**: Try to keep all balls in play during multiball for maximum points
5. **Combo System**: Hit the purple Foot Clan targets in sequence (1-2-3) for combo awards
6. **Aim for Bumpers**: The red bumpers give quick points and exciting action
7. **Use Flippers Strategically**: Timing is everything - flip at the right moment!

## Game Structure

```
tmnt-pinball/
├── main.go           # Main game loop and SDL initialization
├── game.go           # Game state management and logic
├── physics.go        # Physics engine for ball movement
├── flippers.go       # Flipper mechanics
├── input.go          # Keyboard and controller input handling
├── renderer.go       # Rendering utilities
├── go.mod            # Go module file
└── README.md         # This file
```

## Scoring

- **Bumper Hit**: 100 points
- **Pizza Target**: 1,000 points
- **Foot Clan Target**: 500 points
- **Foot Clan Combo**: 5,000 points
- **Episode Target**: 2,000 points
- **Episode Mode**: 3,000 points
- **Pizza Mode**: Double all points!

## Troubleshooting

### SDL2 Not Found

If you get an error about SDL2 not being found:

```bash
brew install sdl2
export CGO_CFLAGS="-I/opt/homebrew/include"
export CGO_LDFLAGS="-L/opt/homebrew/lib"
go build
```

### Controller Not Detected

1. Make sure your controller is connected before starting the game
2. The game uses SDL2's game controller database
3. Most modern controllers (Xbox, PlayStation, Switch Pro) are supported
4. Check that `gamecontrollerdb.txt` exists in the game directory

### Performance Issues

- Close other applications to free up system resources
- The game runs at 60 FPS by default
- Reduce window size if needed (edit `WindowWidth` and `WindowHeight` in `main.go`)

## Future Enhancements

Potential features for future versions:

- Sound effects and music
- More detailed graphics and animations
- Additional game modes (Boss battles, Sewer missions, etc.)
- Network multiplayer
- Tournament mode
- Customizable table layouts
- Visual effects (particle systems, lighting)

## Credits

Inspired by the Stern Teenage Mutant Ninja Turtles pinball machine.

TMNT and all related characters are © Nickelodeon.

This is a fan-made tribute project for educational and entertainment purposes.

## License

MIT License - See LICENSE file for details

## Contributing

Contributions are welcome! Please feel free to submit pull requests or open issues for bugs and feature requests.

---

**Cowabunga! Enjoy the game!** 🐢🍕
