-- +goose Up
CREATE TABLE images (
    id VARCHAR(255) PRIMARY KEY,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    name VARCHAR(255) NOT NULL,
    preset_image_type VARCHAR(255) NOT NULL,
    video_composition_id VARCHAR(255) NOT NULL,
    
    INDEX idx_images_created_at (created_at),
    INDEX idx_images_updated_at (updated_at),
    INDEX idx_images_deleted_at (deleted_at),
    INDEX idx_images_video_composition_id (video_composition_id),
    INDEX idx_images_preset_image_type (preset_image_type),
    
    FOREIGN KEY (video_composition_id) REFERENCES video_compositions(id) ON UPDATE CASCADE ON DELETE RESTRICT
);

-- +goose Down
DROP INDEX idx_images_preset_image_type ON images;
DROP INDEX idx_images_video_composition_id ON images;
DROP INDEX idx_images_deleted_at ON images;
DROP INDEX idx_images_updated_at ON images;
DROP INDEX idx_images_created_at ON images;
DROP TABLE images;
