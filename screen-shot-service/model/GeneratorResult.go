package model

// GeneratorResult : the result out of generator for each screenshot and will send to result service
type GeneratorResult struct {
	InputJSON []byte
	URL       string `json:"url"`
	URLHash   string `json:"url_hash"`
	Path      string `json:"image_path"`
	IsSuccess bool   `json:"is_image_created"`
	Err       error
}
