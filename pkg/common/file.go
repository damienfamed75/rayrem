package common

import (
	"fmt"
	"image"
	"io/ioutil"
	"path/filepath"

	"github.com/damienfamed75/aseprite"
	"github.com/markbates/pkger"
	"github.com/markbates/pkger/pkging"
)

func open(directory string) (pkging.File, error) {
	return pkger.Open(directory)
}

// OpenAsset returns an io.Reader file of an asset file.
func OpenAsset(fileName string) (pkging.File, error) {
	return open(filepath.Join("/assets", fileName))
}

// ReadAsset reads the bytes from an asset file and returns if there were errors.
func ReadAsset(fileName string) ([]byte, error) {
	// Open the asset file as an io.Reader
	f, err := OpenAsset(fileName)
	if err != nil {
		return nil, fmt.Errorf("open asset: %w", err)
	}
	defer f.Close()

	// Read all the bytes from the file.
	raw, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, fmt.Errorf("reading asset: %w", err)
	}

	return raw, nil
}

// LoadPNG loads an asset file and decodes it as a PNG image.
func LoadPNG(fileName string) (image.Image, error) {
	// Opens the asset as a file.
	f, err := OpenAsset(fileName)
	if err != nil {
		return nil, fmt.Errorf("png asset: %w", err)
	}
	defer f.Close()

	// Decode the file as an image.
	img, _, err := image.Decode(f)
	if err != nil {
		return nil, fmt.Errorf("decoding png: %w", err)
	}

	return img, nil
}

// LoadSpritesheet returns an aseprite sheet based on the path given.
func LoadSpritesheet(fileName string) (*aseprite.File, error) {
	// Open and read the asset file for its bytes.
	aseRaw, err := ReadAsset(fileName)
	if err != nil {
		return nil, err
	}

	// Try to open the spritesheet and return the error if not successful.
	ase, err := aseprite.NewFile(aseRaw)
	if err != nil {
		return nil, fmt.Errorf("open aseprite: %w", err)
	}

	return ase, nil
}
