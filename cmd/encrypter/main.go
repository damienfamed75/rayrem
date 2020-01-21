// encrypter simply encrypts the configuration files for packaging of the game.

package main

import (
	"crypto/rand"
	"encoding/hex"
	"io"
	"io/ioutil"
	"log"
	"os"

	"golang.org/x/crypto/nacl/secretbox"
)

var (
	secretKey [32]byte
)

func init() {
	secretKeyBytes, err := hex.DecodeString("87f328c60b3d80c9cdd921b8e85ac2bbb760020ef77e9ab3f8ba629193626133")
	if err != nil {
		panic(err)
	}

	copy(secretKey[:], secretKeyBytes)
}

func loadConfig(path string) ([]byte, error) {
	f, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	return ioutil.ReadAll(f)
}

func encrypt(fBytes []byte, targetPath string) error {
	// You must use a different nonce for each message you encrypt with the
	// same key. Since the nonce here is 192 bits long, a random value
	// provides a sufficiently small probability of repeats.
	var nonce [24]byte
	if _, err := io.ReadFull(rand.Reader, nonce[:]); err != nil {
		panic(err)
	}

	// This encrypts "hello world" and appends the result to the nonce.
	encrypted := secretbox.Seal(nonce[:], fBytes, &nonce, &secretKey)

	f, err := os.OpenFile(targetPath, os.O_CREATE|os.O_RDWR, 0777)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	f.Write(encrypted)

	f.Sync()

	return nil
}

func main() {
	cfgFile, err := loadConfig("config/game.json")
	if err != nil {
		panic(err)
	}

	encrypt(cfgFile, "config/game.config")

	settingsFile, err := loadConfig("config/settings.json")
	if err != nil {
		panic(err)
	}

	encrypt(settingsFile, "config/settings.config")
}
