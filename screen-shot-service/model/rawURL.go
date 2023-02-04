package model

// RawURL : the stander schema will recived from kafka topic (raw_url)
type RawURL struct {
	URL string `json:"url"`
}
