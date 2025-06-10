package main

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	// INITIALIZE AES IN GCM MODE
	key := []byte("thisisthesecretkeythatwillbeused") // 32 bytes key for AES-256
	if len(key) != 32 {
		panic("key length must be 32 bytes for AES-256")
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		panic("error while setting up aes: " + err.Error())
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		panic("error while setting up gcm: " + err.Error())
	}

	// Check current working directory
	cwd, err := os.Getwd()
	if err != nil {
		panic("error getting current directory: " + err.Error())
	}
	fmt.Println("Current working directory:", cwd)

	// Define the directory path
	dirPath := "C:\\Users\\hieng\\Desktop\\Ransomware\\Experiment"

	// Ensure the directory exists
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		panic("directory does not exist: " + dirPath)
	}

	// LOOPING THROUGH TARGET FILES
	err = filepath.WalkDir(dirPath, func(path string, info os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// SKIP IF DIRECTORY
		if info.IsDir() {
			return nil
		}

		// Process only encrypted files (those ending in .enc)
		if filepath.Ext(path) == ".enc" {
			fmt.Println("Decrypting", path, "...")

			// READ ENCRYPTED FILE CONTENTS
			encrypted, err := os.ReadFile(path)
			if err != nil {
				fmt.Println("error while reading encrypted file contents:", err)
				return err
			}

			// Extract nonce from the start of the encrypted content
			nonceSize := gcm.NonceSize()
			if len(encrypted) < nonceSize {
				fmt.Println("error: encrypted file too short")
				return fmt.Errorf("encrypted file too short")
			}
			nonce, ciphertext := encrypted[:nonceSize], encrypted[nonceSize:]

			// DECRYPT BYTES
			plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
			if err != nil {
				fmt.Println("error while decrypting file contents:", err)
				return err
			}

			// WRITE DECRYPTED CONTENTS
			originalFilePath := path[:len(path)-len(".enc")]
			if err := os.WriteFile(originalFilePath, plaintext, 0666); err != nil {
				fmt.Println("error while writing decrypted contents:", err)
				return err
			}

			// DELETE THE ENCRYPTED FILE
			if err := os.Remove(path); err != nil {
				fmt.Println("error while deleting the encrypted file:", err)
				return err
			}
		}

		return nil
	})

	if err != nil {
		fmt.Println("error walking the path:", err)
	}
}
