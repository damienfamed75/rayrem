package common

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os"

	"golang.org/x/crypto/nacl/secretbox"
)

var secretKey [32]byte

func init() {
	secretKeyBytes, err := hex.DecodeString("87f328c60b3d80c9cdd921b8e85ac2bbb760020ef77e9ab3f8ba629193626133")
	if err != nil {
		panic(err)
	}

	copy(secretKey[:], secretKeyBytes)
}

func loadEncryptedConfig() error {
	// Open the encrypted config file.
	f, err := os.Open("config/game.config")
	if err != nil {
		return err
	}
	defer f.Close()

	// Read the bytes from the encrypted config file.
	fBytes, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}

	// When you decrypt, you must use the same nonce and key you used to
	// encrypt the message. One way to achieve this is to store the nonce
	// alongside the encrypted message. Above, we stored the nonce in the first
	// 24 bytes of the encrypted text.
	var decryptNonce [24]byte
	copy(decryptNonce[:], fBytes[:24])

	decrypted, ok := secretbox.Open(nil, fBytes[24:], &decryptNonce, &secretKey)
	if !ok {
		return fmt.Errorf("decryption error: %w", err)
	}

	// Create a reader for the config to use.
	reader := bytes.NewReader(decrypted)

	// Read the bytes from the bytes.Reader and load it into the PublicConfig.
	if err := PublicConfig.ReadConfig(reader); err != nil {
		return fmt.Errorf("viper: %w", err)
	}

	return nil
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
