package scene

import (
	"fmt"
	"strconv"

	"github.com/damienfamed75/rayrem/pkg/camera"
	"github.com/damienfamed75/rayrem/pkg/common"

	r "github.com/lachee/raylib-goplus/raylib"
)

var _ common.Scene = &Menu{}

// Menu us the main menu of the game.
type Menu struct {
	camera       *camera.StaticCamera
	sceneManager common.SceneManager

	states         map[string]bool
	actives        map[string]int
	resolutions    []r.Vector2
	lastFullscreen bool
	fullscreen     bool
	vol            float32
}

// NewMenu creates and sets up settings in the menu.
func NewMenu(sceneManager common.SceneManager) *Menu {
	m := &Menu{
		sceneManager: sceneManager,
		states:       make(map[string]bool),
		actives:      make(map[string]int),
		vol:          float32(common.PublicConfig.GetFloat64("volume.master")),
		resolutions: []r.Vector2{
			r.NewVector2(640, 480),
			r.NewVector2(800, 600),
			r.NewVector2(960, 720),
			r.NewVector2(1024, 768),
			r.NewVector2(1280, 960),
			r.NewVector2(1400, 1050),
			r.NewVector2(1440, 1080),
			r.NewVector2(1600, 1200),
			r.NewVector2(1920, 1080),
			r.NewVector2(1920, 1440),
			r.NewVector2(2048, 1536),
			r.NewVector2(3840, 2160),
		},
	}

	m.camera = camera.NewStatic(r.NewVector2(0, 0), 1)

	currentResX := float32(common.PublicConfig.GetFloat64("screen.width"))
	currentResY := float32(common.PublicConfig.GetFloat64("screen.height"))
	for i, res := range m.resolutions {
		if res.X == currentResX && res.Y == currentResY {
			m.actives["drop"] = i
		}
	}

	m.fullscreen = common.PublicConfig.GetBool("screen.fullscreen")

	return m
}

// Update doesn't really do much in the main menu.
func (m *Menu) Update(dt float32) {}

// Draw draws all the buttons and stuff.
func (m *Menu) Draw() {
	r.BeginMode2D(m.camera.Camera2D)
	r.ClearBackground(r.White)

	// If the settings window is not showing.
	if !m.states["settings"] {
		// Play button
		if r.GuiButton(r.NewRectangle(float32(r.GetScreenWidth()/2)-150, 100, 300, 100), "play") {
			m.sceneManager.SetScene(common.ModeTesting) // Update scene.
		}

		// Settings button
		if r.GuiButton(r.NewRectangle(float32(r.GetScreenWidth()/2)-150, 220, 300, 100), "settings") {
			m.states["settings"] = true // Show settings.
		}
	}

	// If settings should show.
	if m.states["settings"] {
		// global settings position.
		pos := r.NewRectangle(
			float32(r.GetScreenWidth()/2),
			float32(r.GetScreenHeight()/2-r.GetScreenHeight()/4),
			0, 0,
		)

		// Draw the window.
		if r.GuiWindowBox(r.NewRectangle(
			float32(r.GetScreenWidth()/8), float32(r.GetScreenHeight()/8),
			float32(r.GetScreenWidth()-(r.GetScreenWidth()/4)),
			float32(r.GetScreenHeight()-(r.GetScreenHeight()/4)),
		), "settings") {
			m.states["settings"] = false
		}

		// Apply button
		if r.GuiButton(r.NewRectangle(pos.X+pos.Width/2-50, float32(r.GetScreenHeight()-175), 100, 20), "Apply") {
			resolution := m.resolutions[m.actives["drop"]]
			fmt.Println(resolution)

			// If the fullscreen value has changed since last time then toggle.
			if m.fullscreen != m.lastFullscreen {
				m.lastFullscreen = m.fullscreen
				r.ToggleFullscreen()
			}

			r.SetWindowSize(int(resolution.X), int(resolution.Y))
			// Make sure that the window doesn't get stuck in the corner of the monitor screen.
			if r.GetWindowPosition().X == 0 && r.GetWindowPosition().Y == 0 {
				r.SetWindowPosition(50, 50)
			}

			// Update variables in the config file.
			common.PublicConfig.Set("screen.width", resolution.X)
			common.PublicConfig.Set("screen.height", resolution.Y)
			common.PublicConfig.Set("volume.master", m.vol)
			common.PublicConfig.Set("screen.fullscreen", m.fullscreen)

			// Save the configuration file on disk.
			common.SavePublicConfig()
		}

		// Volume slider.
		newVol := r.GuiSlider(r.NewRectangle(pos.X+pos.Width/2-50, pos.Y+150, 100, 20), "Volume - "+strconv.Itoa(int(m.vol*100)), "100", m.vol, 0, 1)
		if newVol != m.vol {
			m.vol = newVol
			r.SetMasterVolume(m.vol) // Update realtime.
		}

		// Fullscreen checkbox.
		newFullscreen := r.GuiCheckBox(r.NewRectangle(pos.X+20, pos.Y+100, 20, 20), "fullscreen", m.fullscreen)
		if newFullscreen != m.fullscreen {
			m.fullscreen = newFullscreen
		}

		// Resolution drop down menu. (THIS SHOULD ALWAYS BE AT THE BOTTOM OF SETTINGS)
		if choose, num := r.GuiDropdownBox(r.NewRectangle(pos.X-100, pos.Y+100, 100, 20), m.getScreenResolutions(), m.actives["drop"], m.states["drop"]); choose {
			m.states["drop"] = !m.states["drop"]
			m.actives["drop"] = num
		}

	}

	r.EndMode2D()
}

// Unload doesn't do much.
func (m *Menu) Unload() {

}

// returns the string for the dropdown menu.
// Should look something like this:
// 640x480;800x600;960x720 etc..etc...
// semicolons represent when there is a new entry in the dropdown.
func (m *Menu) getScreenResolutions() string {
	var result string

	for i, res := range m.resolutions {
		result += strconv.Itoa(int(res.X)) + "x" + strconv.Itoa(int(res.Y))

		if len(m.resolutions) > i+1 {
			result += ";"
		}
	}

	return result
}
