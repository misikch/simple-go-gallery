package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	_ "fmt"
	"github.com/nfnt/resize"
	"image/jpeg"
	"io"
	"math/rand"
	"os"
)

var (
	sizes = []uint{80, 160, 320}
	letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
)

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}

	return string(b)
}

func SaveFile(in io.Reader) (md5FileSum string, err error) {
	tmpName := RandStringRunes(32)

	tmpFilePath := fmt.Sprintf("./images/%s.jpg", tmpName)
	newFile, err := os.Create(tmpFilePath)
	if err != nil {
		return "", fmt.Errorf("failed to create file %q: %w", tmpFilePath, err)
	}

	hasher := md5.New()
	_, err = io.Copy(newFile, io.TeeReader(in, hasher))
	if err != nil {
		return "", fmt.Errorf("copy file failed: %w", err)
	}
	err = newFile.Sync()
	if err != nil {
		return "", fmt.Errorf("failed to store file to disk: %w", err)
	}
	err = newFile.Close()
	if err != nil {
		return "", fmt.Errorf("failed to close file: %w", err)
	}

	md5Sum := hex.EncodeToString(hasher.Sum(nil))

	realFilePath := fmt.Sprintf("./images/%s.jpg", md5Sum)
	err = os.Rename(tmpFilePath, realFilePath)
	if err != nil {
		return "", fmt.Errorf("failed to rename file from %q to %q: %w", tmpFilePath, realFilePath, err)
	}

	return md5Sum, nil
}

func MakeThumbnail(realFilePath string, md5FileSum string) error {
	for _, size := range sizes {
		resizedFilePath := fmt.Sprintf("./images/%s_%d.jpg", md5FileSum, size)
		err := ResizeImage(realFilePath, resizedFilePath, size)
		if err != nil {
			return fmt.Errorf("failed to resize file %q: %w", realFilePath, err)
		}
	}

	return nil
}

func ResizeImage(realFilePath string, resizedFilePath string, size uint) error {
	originalFile, err := os.Open(realFilePath)
	if err != nil {
		return fmt.Errorf("failed to open file %s: %w", realFilePath, err)
	}

	img, err := jpeg.Decode(originalFile)
	if err != nil {
		return fmt.Errorf("failed to decode jpeg file %s: %w", realFilePath, err)
	}
	originalFile.Close()

	resizedImg := resize.Resize(size, 0, img, resize.Lanczos2)

	out, err := os.Create(resizedFilePath)
	if err != nil {
		return fmt.Errorf("failed to crete file %q: %w", resizedFilePath, err)
	}
	defer out.Close()

	err = jpeg.Encode(out, resizedImg, nil)
	if err != nil {
		return fmt.Errorf("failed to jpeg.Encode resized image: %w", err)
	}

	return nil
}



