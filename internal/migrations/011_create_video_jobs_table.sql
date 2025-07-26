-- +goose Up
CREATE TABLE video_composition_jobs (
    id VARCHAR(255) PRIMARY KEY,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    video_composition_id VARCHAR(255) NOT NULL,
    order_line_id VARCHAR(255) NOT NULL,
    status VARCHAR(255) NOT NULL,
    
    INDEX idx_video_composition_jobs_created_at (created_at),
    INDEX idx_video_composition_jobs_updated_at (updated_at),
    INDEX idx_video_composition_jobs_deleted_at (deleted_at),
    INDEX idx_video_composition_jobs_video_composition_id (video_composition_id),
    INDEX idx_video_composition_jobs_order_line_id (order_line_id),
    INDEX idx_video_composition_jobs_status (status),
    
    FOREIGN KEY (video_composition_id) REFERENCES video_compositions(id) ON UPDATE CASCADE ON DELETE RESTRICT,
    FOREIGN KEY (order_line_id) REFERENCES order_lines(id) ON UPDATE CASCADE ON DELETE RESTRICT
);

-- +goose Down
DROP INDEX idx_video_composition_jobs_status ON video_composition_jobs;
DROP INDEX idx_video_composition_jobs_order_line_id ON video_composition_jobs;
DROP INDEX idx_video_composition_jobs_video_composition_id ON video_composition_jobs;
DROP INDEX idx_video_composition_jobs_deleted_at ON video_composition_jobs;
DROP INDEX idx_video_composition_jobs_updated_at ON video_composition_jobs;
DROP INDEX idx_video_composition_jobs_created_at ON video_composition_jobs;
DROP TABLE video_composition_jobs;
