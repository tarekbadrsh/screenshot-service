package generator

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"net/url"
	"screen-shot-service/config"
	"screen-shot-service/logger"
	"screen-shot-service/model"
	"screen-shot-service/storage"
	"sync"
)

var (
	internalGenerator *IGenerator
	configuration     config.Config
	initGeneratorOnce sync.Once
)

// IGenerator :
type IGenerator interface {
	ScreenshotURL(targetURL string, destination string) error
}

// InitializeGenerator : initialize the main Generator
func InitializeGenerator(gen *IGenerator) {
	initGeneratorOnce.Do(func() {
		internalGenerator = gen
		configuration = config.Configuration()
		logger.Info("Generator has been initialized")
	})
}

// ScreenshotURL : generate screen-shot of url
func ScreenshotURL(targetURL string, destination string) error {
	return (*internalGenerator).ScreenshotURL(targetURL, destination)
}

func handler(result *model.GeneratorResult) {
	raw := model.RawURL{}
	result.Err = json.Unmarshal(result.InputJSON, &raw)
	if result.Err != nil {
		logger.Error(result.Err)
		result.URL = string(result.InputJSON)
		return
	}

	u := &url.URL{}
	u, result.Err = url.ParseRequestURI(raw.URL)
	if result.Err != nil {
		logger.Errorf("Parse Url Error %v", result.Err)
		result.URL = string(result.InputJSON)
		return
	}
	result.URL = u.String()
	// get url hash
	urlHash := md5.Sum([]byte(result.URL))
	result.URLHash = fmt.Sprintf("%x", urlHash)

	result.Err = storage.ImagePath(configuration.StoragePath, result)
	if result.Err != nil {
		logger.Errorf("Error while creating storage directory %v", result.Err)
		return
	}

	result.Err = ScreenshotURL(result.URL, result.Path)
	if result.Err != nil {
		logger.Errorf("Generate Screenshot error : %v", result.Err)
		return
	}
}

// GetGeneratorHandler :
func GetGeneratorHandler() func([]byte) error {
	return func(input []byte) error {
		result := &model.GeneratorResult{InputJSON: input, IsSuccess: false}

		// retry with retry_count until success or end trials.
		for i := 1; i <= configuration.RetryCount; i++ {
			handler(result)
			if result.Err == nil {
				result.IsSuccess = true
				break
			}
			logger.Warnf("Handle message failed massage(%s) err(%v); trial(%v)", input, result.Err, i)
		}

		// store the result
		storErr := storage.SendGeneratedResult(result)
		if storErr != nil {
			logger.Error("Store result error : %v", storErr)
			return storErr
		}
		if result.Err != nil || storErr != nil {
			return fmt.Errorf("Result Error (%v), Storage Error (%v)", result.Err, storErr)
		}
		logger.Infof("Screenshot of (%v) successfully generated and the result has been stored in Result service", result.URL)
		return nil
	}
}
