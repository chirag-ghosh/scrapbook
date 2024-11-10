package indexer

import (
	"fmt"
	"log"
	"mime"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/chirag-ghosh/scrapbook/db"
	"github.com/chirag-ghosh/scrapbook/models"
	"github.com/rwcarlsen/goexif/exif"
)

func checkRootDirectoryIndexState(dirPath string, lastModifiedTime time.Time) (bool, error) {
	dirPathAbs, err := filepath.Abs(dirPath)
	if err != nil {
		return false, fmt.Errorf("failed to get absolute path: %v", err)
	}

	db := db.GetDB()
	rows, err := db.Query("SELECT indexed_at FROM index_directories WHERE path = ?", dirPathAbs)
	if err != nil {
		return false, fmt.Errorf("failed to query database: %v", err)
	}
	defer rows.Close()

	var indexedAt time.Time
	if rows.Next() {
		err = rows.Scan(&indexedAt)
		if err != nil {
			return false, fmt.Errorf("failed to scan database row: %v", err)
		}
	}

	if indexedAt.After(lastModifiedTime) {
		return true, nil
	}

	return false, nil
}

func createRootDirectoryIndex(dirName string, dirPath string) (int, error) {
	dirPathAbs, err := filepath.Abs(dirPath)
	if err != nil {
		return -1, fmt.Errorf("failed to get absolute path: %v", err)
	}

	db := db.GetDB()
	_, err = db.Exec("INSERT INTO index_directories (name, path, indexed_at) VALUES (?, ?, ?) ON CONFLICT(path) DO UPDATE SET name = excluded.name, indexed_at = excluded.indexed_at", dirName, dirPathAbs, time.Now())
	if err != nil {
		return -1, fmt.Errorf("failed to insert into database: %v", err)
	}

	rows, err := db.Query("SELECT id FROM index_directories WHERE path = ?", dirPathAbs)
	if err != nil {
		return -1, fmt.Errorf("failed to query database: %v", err)
	}
	defer rows.Close()

	var id int
	if rows.Next() {
		err = rows.Scan(&id)
		if err != nil {
			return -1, fmt.Errorf("failed to scan database row: %v", err)
		}
	}

	return id, nil
}

func reduceRational(num int64, den int64) (int64, int64) {
	if num == 0 {
		return 0, den
	}

	if den == 0 {
		return num, 0
	}

	for i := num; i > 0; i-- {
		if num%i == 0 && den%i == 0 {
			num /= i
			den /= i
			break
		}
	}

	return num, den
}

func addPhoto(directoryId int, photoPath string) error {
	var photo models.Photo

	photo.FileDir = filepath.Dir(photoPath)
	photo.Name = filepath.Base(photoPath)

	mimeType := mime.TypeByExtension(filepath.Ext(photoPath))
	if !strings.HasPrefix(mimeType, "image/") {
		return nil
	}

	file, err := os.Open(photoPath)
	if err != nil {
		return fmt.Errorf("failed to open file %s: %v", photoPath, err)
	}
	defer file.Close()

	exifData, err := exif.Decode(file)
	if err == nil {
		tag, err := exifData.Get(exif.Make)
		if err == nil {
			photo.CameraMake, _ = tag.StringVal()
		}

		tag, err = exifData.Get(exif.Model)
		if err == nil {
			photo.CameraModel, _ = tag.StringVal()
		}

		tag, err = exifData.Get(exif.LensModel)
		if err == nil {
			photo.LensID, _ = tag.StringVal()
		}

		tag, err = exifData.Get(exif.PixelXDimension)
		if err == nil {
			photo.Width, _ = tag.Int(0)
		}

		tag, err = exifData.Get(exif.PixelYDimension)
		if err == nil {
			photo.Height, _ = tag.Int(0)
		}

		tag, err = exifData.Get(exif.FocalLength)
		if err == nil {
			num, den, _ := tag.Rat2(0)
			if den != 0 {
				photo.FocalLength = float64(num) / float64(den)
			}
		}

		tag, err = exifData.Get(exif.FNumber)
		if err == nil {
			num, den, _ := tag.Rat2(0)
			if den != 0 {
				photo.Aperture = float64(num) / float64(den)
			}
		}

		tag, err = exifData.Get(exif.ExposureTime)
		if err == nil {
			num, den, _ := tag.Rat2(0)
			if den != 0 {
				num, den = reduceRational(num, den)
				photo.ShutterSpeed = fmt.Sprintf("%d/%d", num, den)
			}
		}

		tag, err = exifData.Get(exif.ISOSpeedRatings)
		if err == nil {
			photo.ISO, _ = tag.Int(0)
		}

		tag, err = exifData.Get(exif.DateTime)
		if err == nil {
			photo.CapturedAt, _ = tag.StringVal()
		}
	}

	db := db.GetDB()
	_, err = db.Exec(`
		INSERT INTO photos (directory_id, file_dir, name, camera_make, camera_model, lens_id, width, height, focal_length, aperture, shutter_speed, iso, captured_at)
		VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?)
	`, directoryId, photo.FileDir, photo.Name, photo.CameraMake, photo.CameraModel, photo.LensID, photo.Width, photo.Height, photo.FocalLength, photo.Aperture, photo.ShutterSpeed, photo.ISO, photo.CapturedAt)
	if err != nil {
		return fmt.Errorf("failed to insert into database: %v", err)
	}

	return nil
}

func IndexRootDirectory(dirName string, dirPath string) error {
	fileInfo, err := os.Stat(dirPath)
	if os.IsNotExist(err) {
		return fmt.Errorf("directory does not exist")
	} else if err != nil {
		return fmt.Errorf("failed to get file info: %v", err)
	}

	if !fileInfo.IsDir() {
		return fmt.Errorf("path is not a directory")
	}

	isIndexed, err := checkRootDirectoryIndexState(dirPath, fileInfo.ModTime())
	if err != nil {
		return fmt.Errorf("failed to check index state: %v", err)
	}

	if isIndexed {
		log.Println("Directory is already indexed")
		return nil
	}

	directoryId, err := createRootDirectoryIndex(dirName, dirPath)
	if err != nil {
		return fmt.Errorf("failed to create index: %v", err)
	}

	err = filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("failed to walk directory: %v", err)
		}

		if !info.IsDir() {
			err := addPhoto(directoryId, path)
			if err != nil {
				log.Println(err)
			}
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("failed to walk directory: %v", err)
	}

	return nil
}
