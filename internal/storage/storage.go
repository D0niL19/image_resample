package storage

import (
	"fmt"
	"os"
)

type DirectoryStorage struct {
	origPath string
	resPath  string
}

func NewDirectoryStorage(origPath, resPath string) *DirectoryStorage {
	return &DirectoryStorage{origPath: origPath, resPath: resPath}
}

func (s *DirectoryStorage) CheckAndRetrieveResized(hash string, width, height int) (string, bool) {
	filePath := fmt.Sprintf("%s%s_%dx%d.jpg", s.origPath, hash, width, height)
	if _, err := os.Stat(filePath); err == nil {
		return filePath, true
	}
	return "", false
}

func (s *DirectoryStorage) SaveOriginal(hash string, data []byte) error {
	filePath := fmt.Sprintf("%s/%s.jpg", s.origPath, hash)
	if _, err := os.Stat(filePath); err == nil {
		return nil
	}
	return os.WriteFile(filePath, data, 0644)
}

func (s *DirectoryStorage) SaveResized(hash string, width, height int, data []byte) error {
	filePath := fmt.Sprintf("%s/%s_%dx%d.jpg", s.origPath, hash, width, height)
	return os.WriteFile(filePath, data, 0644)
}
