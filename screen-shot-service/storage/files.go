package storage

import (
	"fmt"
	"os"
	"screen-shot-service/logger"
	"screen-shot-service/model"
	"time"
)

const (
	dateStorelayout = "2006/01/02/15/04"
)

// ImagePath : get new image path
func ImagePath(storagePath string, result *model.GeneratorResult) error {

	// create directory to store the images
	timePath := time.Now().Format(dateStorelayout)
	dirPath := fmt.Sprintf("%v/%v", storagePath, timePath)
	err := createDirectoryRecursively(dirPath)
	if err != nil {
		logger.Error(err)
		return err
	}

	result.Path = fmt.Sprintf("%v/%v.png", dirPath, result.URLHash)
	return nil
}

func createDirectoryRecursively(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		Mkdirerr := os.MkdirAll(path, os.ModePerm)
		if Mkdirerr != nil {
			logger.Errorf("Error while creating directory %v", logger.WithFields(map[string]interface{}{"Path": path, "error": Mkdirerr}))
			return err
		}
	}
	return nil
}
