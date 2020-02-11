package common

import (
	"crypto/rand"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"os"

	"golang.org/x/crypto/nacl/secretbox"
)

// SavePublicConfig unmarshals the viper config into a map then remarshals it
// into bytes, encrypts the data using salt encryption and then overwrites the
// current game.config or if you're not using game.config then it'll overwrite
// the game.json file instead.
func SavePublicConfig() {
	// if the game config is encrypted then save it encrypted.
	if usingEncrypted {
		tmp := make(map[string]interface{})
		PublicConfig.Unmarshal(&tmp)

		writeEncrypted(tmp)
	} else {
		// non encrypted save.
		PublicConfig.WriteConfig()
	}

	// Update the controls structure.
	loadControls()
}

func loadTemp(path string) ([]byte, error) {
	f, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	return ioutil.ReadAll(f)
}

func writeEncrypted(tmp map[string]interface{}) {
	// You must use a different nonce for each message you encrypt with the
	// same key. Since the nonce here is 192 bits long, a random value
	// provides a sufficiently small probability of repeats.
	var nonce [24]byte
	if _, err := io.ReadFull(rand.Reader, nonce[:]); err != nil {
		panic(err)
	}

	cfgFile, err := json.Marshal(tmp)
	if err != nil {
		panic(err)
	}

	// This encrypts "hello world" and appends the result to the nonce.
	encrypted := secretbox.Seal(nonce[:], cfgFile, &nonce, &secretKey)

	if err := os.Remove("config/game.config"); err != nil {
		panic(err)
	}

	f, err := os.OpenFile("config/game.config", os.O_CREATE|os.O_RDWR, 0777)
	if err != nil {
		panic(err)
	}

	f.Write(encrypted)

	f.Sync()

	f.Close()
}
