package dal

import (
	"screen-shot-api/db"
)

// ScreenshotDAL : data access layer  (screenshot) table.
type ScreenshotDAL struct {
	ID             int32  `json:"id" gorm:"column:id;primary_key:true"`
	URL            string `json:"url" gorm:"column:url"`
	URLHash        string `json:"url_hash" gorm:"column:url_hash"`
	IsImageCreated bool   `json:"is_image_created" gorm:"column:is_image_created"`
	ImagePath      string `json:"image_path" gorm:"column:image_path"`
	CreatedAt      int64  `json:"created_at" gorm:"column:created_at"`
}

// TableName sets the insert table name for this struct type
func (s *ScreenshotDAL) TableName() string {
	return "screenshot"
}

// GetAllScreenshots : get all screenshots.
func GetAllScreenshots() []*ScreenshotDAL {
	screenshots := []*ScreenshotDAL{}
	db.DB().Find(&screenshots)
	return screenshots
}

// GetScreenshot : get one screenshot by id.
func GetScreenshot(id int32) (*ScreenshotDAL, error) {
	s := &ScreenshotDAL{}
	result := db.DB().First(s, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return s, nil
}

// CreateScreenshot : create new screenshot.
func CreateScreenshot(s *ScreenshotDAL) (*ScreenshotDAL, error) {
	result := db.DB().Create(s)
	if result.Error != nil {
		return nil, result.Error
	}
	return s, nil
}

// UpdateScreenshot : update exist screenshot.
func UpdateScreenshot(s *ScreenshotDAL) (*ScreenshotDAL, error) {
	_, err := GetScreenshot(s.ID)
	if err != nil {
		return nil, err
	}
	result := db.DB().Save(s)
	if result.Error != nil {
		return nil, result.Error
	}
	return s, nil
}

// DeleteScreenshot : delete screenshot by id.
func DeleteScreenshot(id int32) error {
	s, err := GetScreenshot(id)
	if err != nil {
		return err
	}
	result := db.DB().Delete(s)
	return result.Error
}
