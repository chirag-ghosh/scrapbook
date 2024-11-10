package models

type Photo struct {
	ID           int
	FileDir      string
	Name         string
	CameraMake   string
	CameraModel  string
	LensID       string
	Width        int
	Height       int
	FocalLength  float64
	Aperture     float64
	ShutterSpeed string // string to handle values in the form of a/b in database
	ISO          int
	CapturedAt   string
}
