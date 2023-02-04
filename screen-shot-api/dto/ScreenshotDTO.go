package dto

import (
	"fmt"
	"screen-shot-api/config"
	"screen-shot-api/dal"
)

// ScreenshotDTO : data transfer object  (screenshot) table.
type ScreenshotDTO struct {
	ID             int32  `json:"id"`
	URL            string `json:"url"`
	URLHash        string `json:"url_hash"`
	IsImageCreated bool   `json:"is_image_created"`
	ImagePath      string `json:"image_path"`
	CreatedAt      int64  `json:"created_at"`
	ScreenShotURL  string `json:"screen_shot_url"`
}

// ScreenshotDTOToDAL : convert ScreenshotDTO to ScreenshotDAL
func (a *ScreenshotDTO) ScreenshotDTOToDAL() (*dal.ScreenshotDAL, error) {
	screenshot := &dal.ScreenshotDAL{
		ID:             a.ID,
		URL:            a.URL,
		URLHash:        a.URLHash,
		IsImageCreated: a.IsImageCreated,
		ImagePath:      a.ImagePath,
		CreatedAt:      a.CreatedAt,
	}
	return screenshot, nil
}

// ScreenshotDALToDTO : convert ScreenshotDAL to ScreenshotDTO
func ScreenshotDALToDTO(a *dal.ScreenshotDAL) (*ScreenshotDTO, error) {
	c := config.Configuration()
	screenshot := &ScreenshotDTO{
		ID:             a.ID,
		URL:            a.URL,
		URLHash:        a.URLHash,
		IsImageCreated: a.IsImageCreated,
		ImagePath:      a.ImagePath,
		CreatedAt:      a.CreatedAt,
		ScreenShotURL:  fmt.Sprintf("%v%v", c.ScreenShotServer, a.ImagePath),
	}
	return screenshot, nil
}

// ScreenshotDALToDTOArr : convert Array of ScreenshotDAL to Array of ScreenshotDTO
func ScreenshotDALToDTOArr(screenshots []*dal.ScreenshotDAL) ([]*ScreenshotDTO, error) {
	var err error
	res := make([]*ScreenshotDTO, len(screenshots))
	for i, screenshot := range screenshots {
		res[i], err = ScreenshotDALToDTO(screenshot)
		if err != nil {
			return res, err
		}
	}
	return res, nil
}
