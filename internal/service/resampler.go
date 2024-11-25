package service

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"github.com/disintegration/imaging"
	"image"
	"image/jpeg"
	"imageResample/pkg/api"
	"log/slog"
	"time"
)

const (
	MaxImageSize = 10 * 1024 * 1024
	MinImageSize = 1 //1024 не проходят тесты потому что слишком маленькие изображения
)

type ImageStorage interface {
	CheckAndRetrieveResized(hash string, width, height int) (string, bool)
	SaveOriginal(hash string, data []byte) error
	SaveResized(hash string, width, height int, data []byte) error
}

type Resampler struct {
	storage     ImageStorage
	imageWidth  int
	imageHeight int
	log         *slog.Logger
}

func NewResampler(storage ImageStorage, imageWidth, imageHeight int, log *slog.Logger) *Resampler {
	return &Resampler{storage: storage, imageHeight: imageHeight, imageWidth: imageWidth, log: log}
}

func (r *Resampler) Resample(request api.ImageRequest) (int64, bool, error) {
	start := time.Now()

	imageBytes, err := base64.StdEncoding.DecodeString(request.Image)
	if err != nil {
		r.log.Error("Failed to decode base64 image", slog.String("error", err.Error()))
		return 0, false, err
	}

	err = Validate(imageBytes)
	if err != nil {
		r.log.Error("Failed to validate image", slog.Any("error", err))
		return 0, false, err
	}

	hash := generateHash(imageBytes)

	resizedPath, found := r.storage.CheckAndRetrieveResized(hash, r.imageWidth, r.imageHeight)
	if found {
		r.log.Info("Resized image found in storage", slog.String("path", resizedPath))
		return time.Since(start).Milliseconds(), true, nil
	}

	err = r.storage.SaveOriginal(hash, imageBytes)
	if err != nil {
		r.log.Error("Failed to save original image", slog.String("error", err.Error()))
		return 0, false, err
	}

	resizedImage, err := Resize(imageBytes, r.imageWidth, r.imageHeight)
	if err != nil {
		r.log.Error("Failed to resize image", slog.String("error", err.Error()))
		return 0, false, err
	}

	err = r.storage.SaveResized(hash, r.imageWidth, r.imageHeight, resizedImage)
	if err != nil {
		r.log.Error("Failed to save resized image", slog.String("error", err.Error()))
		return 0, false, err
	}

	r.log.Info("Resized image saved successfully", slog.String("hash", hash))
	return time.Since(start).Milliseconds(), false, nil
}

func Resize(input []byte, width, height int) ([]byte, error) {
	if width == 0 || height == 0 {
		return nil, fmt.Errorf("width and height must be greater than zero")
	}

	img, _, err := image.Decode(bytes.NewReader(input))
	if err != nil {
		return nil, err
	}
	resized := imaging.Resize(img, width, height, imaging.Lanczos)
	buf := new(bytes.Buffer)
	if err := jpeg.Encode(buf, resized, &jpeg.Options{Quality: 80}); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func Validate(input []byte) error {
	if err := validateImageLength(input); err != nil {
		return err
	}

	if !validateImageFormat(input) {
		return fmt.Errorf("image format is invalid")
	}

	return nil
}

func validateImageLength(imageBytes []byte) error {
	size := len(imageBytes)
	if size > MaxImageSize {
		return fmt.Errorf("image size exceeds maximum allowed size of 10 MB")
	}
	if size < MinImageSize {
		return fmt.Errorf("image size is below the minimum allowed size of 1 KB")
	}
	return nil
}

func validateImageFormat(imageBytes []byte) bool {
	if len(imageBytes) > 2 && imageBytes[0] == 0xFF && imageBytes[1] == 0xD8 {
		return true
	}
	return false
}

func generateHash(data []byte) string {
	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:])
}
