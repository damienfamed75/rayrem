package common

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/markbates/pkger"

	// "github.com/gobuffalo/packr/v2"
	"github.com/spf13/viper"
)

// Global configuration files.
var (
	Config *configuration

	PublicConfig = viper.New()

	usingEncrypted bool
)

type configuration struct {
	Game struct {
		Gravity     float32 `json:"gravity"`
		EntityScale float32 `json:"entityScale"`
	} `json:"game"`
	Player struct {
		Spritesheet string  `json:"spritesheet"`
		Friction    float32 `json:"friction"`
		JumpHeight  float32 `json:"jumpHeight"`
		MaxSpeed    struct {
			X float32 `json:"X"`
			Y float32 `json:"Y"`
		} `json:"maxSpeed"`
	} `json:"player"`
	Camera struct {
		Zoom float32 `json:"zoom"`
		Lerp float32 `json:"lerp"`
	} `json:"camera"`
}

// LoadConfig loads in the debug and public configuration files.
func LoadConfig() error {
	if err := loadDebug(); err != nil {
		return fmt.Errorf("debug: %w", err)
	}

	if err := loadPublic(); err != nil {
		return fmt.Errorf("public: %w", err)
	}

	return nil
}

func loadDebug() error {
	var cfgRaw []byte

	cfgFile, err := pkger.Open("/config/settings.json")
	if err != nil {
		return fmt.Errorf("open debug config: %w", err)
	}

	cfgRaw, err = ioutil.ReadAll(cfgFile)
	if err != nil {
		return fmt.Errorf("reading debug config: %w", err)
	}

	if err := json.Unmarshal(cfgRaw, &Config); err != nil {
		return fmt.Errorf("unmarshal debug config: %w", err)
	}

	return nil
}

func loadPublic() error {
	setDefaults()

	PublicConfig.SetConfigType("json")

	if fileExists("config/game.config") {
		if err := loadEncryptedConfig(); err != nil {
			return err
		}

		usingEncrypted = true
	} else {
		PublicConfig.SetConfigFile("config/game.json")
		if err := PublicConfig.ReadInConfig(); err != nil {
			return err
		}
	}

	return nil
}

func setDefaults() {
	PublicConfig.SetDefault("screen.width", 800)
	PublicConfig.SetDefault("screen.height", 600)
	PublicConfig.SetDefault("volume.music", 1.0)
	PublicConfig.SetDefault("volume.sound", 1.0)
	PublicConfig.SetDefault("volume.master", 1.0)
}
