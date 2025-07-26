-- +goose Up
CREATE TABLE video_compositions (
    id VARCHAR(255) PRIMARY KEY,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    customer_id VARCHAR(255) NOT NULL,
    video_template VARCHAR(255) NOT NULL,
    status VARCHAR(255) NOT NULL,
    
    INDEX idx_video_compositions_created_at (created_at),
    INDEX idx_video_compositions_updated_at (updated_at),
    INDEX idx_video_compositions_deleted_at (deleted_at),
    INDEX idx_video_compositions_customer_id (customer_id),
    INDEX idx_video_compositions_status (status),
    
    FOREIGN KEY (customer_id) REFERENCES customers(id) ON UPDATE CASCADE ON DELETE SET NULL
);

-- +goose Down
DROP INDEX idx_video_compositions_status ON video_compositions;
DROP INDEX idx_video_compositions_customer_id ON video_compositions;
DROP INDEX idx_video_compositions_deleted_at ON video_compositions;
DROP INDEX idx_video_compositions_updated_at ON video_compositions;
DROP INDEX idx_video_compositions_created_at ON video_compositions;
DROP TABLE video_compositions;
