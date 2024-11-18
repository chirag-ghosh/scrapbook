package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/chirag-ghosh/scrapbook/db"
	"github.com/chirag-ghosh/scrapbook/models"
)

func handleHealthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Scrapbook API is running.")
}

func handleTimeline(w http.ResponseWriter, r *http.Request) {
	pageParam := r.URL.Query().Get("page")
	limitParam := r.URL.Query().Get("limit")

	page, err := strconv.Atoi(pageParam)
	if err != nil {
		page = 1
	}

	limit, err := strconv.Atoi(limitParam)
	if err != nil {
		limit = 10
	}

	offset := (page - 1) * limit

	query := `SELECT id, file_dir, name, camera_make, camera_model, lens_id, width, height, focal_length, aperture, shutter_speed, iso, captured_at 
			  FROM photos
			  WHERE name LIKE '%.jpg' OR name LIKE '%.png' OR name LIKE '%.jpeg'
			  ORDER BY captured_at DESC 
			  LIMIT $1 OFFSET $2`

	db := db.GetDB()

	rows, err := db.Query(query, limit, offset)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error querying database: %v", err), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var photos []models.Photo
	for rows.Next() {
		var photo models.Photo
		err := rows.Scan(&photo.ID, &photo.FileDir, &photo.Name, &photo.CameraMake, &photo.CameraModel, &photo.LensID, &photo.Width, &photo.Height, &photo.FocalLength, &photo.Aperture, &photo.ShutterSpeed, &photo.ISO, &photo.CapturedAt)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error scanning row: %v", err), http.StatusInternalServerError)
			return
		}

		photos = append(photos, photo)
	}

	if err := rows.Err(); err != nil {
		http.Error(w, fmt.Sprintf("Error scanning row: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(photos); err != nil {
		http.Error(w, fmt.Sprintf("Error encoding JSON: %v", err), http.StatusInternalServerError)
	}
}

func handlePhotoServe(w http.ResponseWriter, r *http.Request) {
	imageIdParam := r.PathValue("id")

	imageId, err := strconv.Atoi(imageIdParam)
	if err != nil {
		http.Error(w, "Invalid image ID", http.StatusBadRequest)
		return
	}

	db := db.GetDB()

	var photo models.Photo
	err = db.QueryRow(`SELECT id, file_dir, name FROM photos WHERE id = $1`, imageId).Scan(&photo.ID, &photo.FileDir, &photo.Name)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error querying database: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Cache-Control", "public, max-age=3600")
	http.ServeFile(w, r, fmt.Sprintf("%s/%s", photo.FileDir, photo.Name))
}
