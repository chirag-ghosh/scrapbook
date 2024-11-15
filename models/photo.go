package models

type Photo struct {
	ID           int     `json:"id"`
	FileDir      string  `json:"file_dir"`
	Name         string  `json:"name"`
	CameraMake   string  `json:"camera_make"`
	CameraModel  string  `json:"camera_model"`
	LensID       string  `json:"lens_id"`
	Width        int     `json:"width"`
	Height       int     `json:"height"`
	FocalLength  float64 `json:"focal_length"`
	Aperture     float64 `json:"aperture"`
	ShutterSpeed string  `json:"shutter_speed"` // string to handle values in the form of a/b in database
	ISO          int     `json:"iso"`
	CapturedAt   string  `json:"captured_at"`
}
