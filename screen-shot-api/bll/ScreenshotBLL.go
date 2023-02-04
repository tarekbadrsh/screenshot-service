package bll

import (
	"strconv"

	"screen-shot-api/dal"
	"screen-shot-api/dto"
)

// ConvertID : covnert ID string to ID int32.
func ConvertID(str string) (int32, error) {
	pram, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return 0, err
	}
	id := int32(pram)
	return id, nil
}

// GetAllScreenshots : get All screenshots.
func GetAllScreenshots() ([]*dto.ScreenshotDTO, error) {
	screenshots := dal.GetAllScreenshots()
	return dto.ScreenshotDALToDTOArr(screenshots)
}

// GetScreenshot : get one screenshot by id.
func GetScreenshot(id int32) (*dto.ScreenshotDTO, error) {
	s, err := dal.GetScreenshot(id)
	if err != nil {
		return nil, err
	}
	return dto.ScreenshotDALToDTO(s)
}

// CreateScreenshot : create new screenshot.
func CreateScreenshot(s *dto.ScreenshotDTO) (*dto.ScreenshotDTO, error) {
	screenshot, err := s.ScreenshotDTOToDAL()
	if err != nil {
		return nil, err
	}
	newscreenshot, err := dal.CreateScreenshot(screenshot)
	if err != nil {
		return nil, err
	}
	return dto.ScreenshotDALToDTO(newscreenshot)
}

// UpdateScreenshot : update exist screenshot.
func UpdateScreenshot(s *dto.ScreenshotDTO) (*dto.ScreenshotDTO, error) {
	screenshot, err := s.ScreenshotDTOToDAL()
	if err != nil {
		return nil, err
	}
	updatescreenshot, err := dal.UpdateScreenshot(screenshot)
	if err != nil {
		return nil, err
	}
	return dto.ScreenshotDALToDTO(updatescreenshot)
}

// DeleteScreenshot : delete screenshot by id.
func DeleteScreenshot(id int32) error {
	return dal.DeleteScreenshot(id)
}
