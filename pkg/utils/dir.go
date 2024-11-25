package utils

import (
	"log"
	"os"
)

func EnsureDir(path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := os.MkdirAll(path, 0755); err != nil {
			log.Fatalf("Не удалось создать директорию %s: %v", path, err)
		}
	}
}
