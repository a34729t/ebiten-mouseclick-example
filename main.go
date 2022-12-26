package main

import (
	"log"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

var USE_INPUT_MANAGER = true // Change to demonstrate mouse click weirdness
const BUTTON0 = ebiten.MouseButton0

type Game struct {
	state        *MouseState
	inputManager *InputManager
}

type MouseState struct {
	FirstPressedTs int64
}

// main starts the program.
// If USE_INPUT_MANAGER is enabled, start the input manager listener in a goroutine.
// Else, listen for input in main thread.
func main() {
	game := &Game{
		state:        &MouseState{},
		inputManager: NewInputManager(),
	}

	if USE_INPUT_MANAGER {
		go game.inputManager.Listen()
	}

	if err := ebiten.RunGame(game); err != nil {
		log.Fatalf("Error while running game: %s", err)
	}
}

// Update checks for mouse input.
// If USE_INPUT_MANAGER is enabled, check channel for mouse clicks.
// Else, call detectMouseClick directly!
func (g *Game) Update() error {
	if USE_INPUT_MANAGER {
		g.inputManager.HandleInput()
	} else {
		detectMouseClick(g.state)
	}
	return nil
}
func (g *Game) Draw(screen *ebiten.Image)                       {}
func (g *Game) Layout(int, int) (screenWidth, screenHeight int) { return 1024, 768 }

// detectMouseClick checks for mouse press/release events and tracks the first press time
// and logs duration of press on release
func detectMouseClick(s *MouseState) bool {
	if inpututil.IsMouseButtonJustPressed(BUTTON0) {
		s.FirstPressedTs = time.Now().UnixNano()
		log.Print("Mouse just pressed")
	}
	if inpututil.IsMouseButtonJustReleased(BUTTON0) {
		durationMillis := (time.Now().UnixNano() - s.FirstPressedTs) / 1000000
		x, y := ebiten.CursorPosition()
		log.Printf("Mouse released at <%d, %d> with %d", x, y, durationMillis)
		return true
	}
	return false
}

// InputManager listens for mouse input in a goroutine, and pushes a notification on a channel
// when a mouse click is detected
type InputManager struct {
	inputChannel chan bool
	state        *MouseState
}

func NewInputManager() *InputManager {
	inputManager := InputManager{
		inputChannel: make(chan bool),
		state:        &MouseState{},
	}
	return &inputManager
}

// HandleInput performs a nonblocking read of the channel and checks for mouse clicks.
func (inputManager *InputManager) HandleInput() {
	// nonblocking read of channel
	select {
	case <-inputManager.inputChannel:
		log.Println("Mouse click!")
	default:
	}
}

// Listen starts the input manager loop
func (inputManager *InputManager) Listen() {
	for {
		click := detectMouseClick(inputManager.state)
		if click {
			inputManager.inputChannel <- true
		}
	}
}
