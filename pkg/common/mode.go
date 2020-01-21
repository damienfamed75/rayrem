package common

// Mode is used to tell what scene to draw in the game.
type Mode uint

// All these are iota+0 so then if something goes wrong then it'll default
// to the main menu.
const (
	ModeMainMenu Mode = iota
	ModeGame
	ModeTesting
	ModeGameOver
)
