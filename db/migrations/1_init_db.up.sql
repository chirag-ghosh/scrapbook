CREATE TABLE IF NOT EXISTS photos (
    id INTEGER PRIMARY KEY,
    directory_id INT NOT NULL,
    file_dir VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    camera_make VARCHAR(255),
    camera_model VARCHAR(255),
    lens_id VARCHAR(255),
    width INT NOT NULL,
    height INT NOT NULL,
    focal_length FLOAT,
    aperture FLOAT,
    shutter_speed VARCHAR(255),
    iso INT,
    captured_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (directory_id, file_dir, name),
    FOREIGN KEY (directory_id) REFERENCES index_directories(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS index_directories (
    id INTEGER PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    path VARCHAR(255) NOT NULL UNIQUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    indexed_at TIMESTAMP
);
