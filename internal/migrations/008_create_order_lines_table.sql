-- +goose Up
CREATE TABLE order_lines (
    id VARCHAR(255) PRIMARY KEY,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    order_id VARCHAR(255) NOT NULL,
    video_composition_id VARCHAR(255) NOT NULL,
    product_id VARCHAR(255) NOT NULL,
    quantity INT NOT NULL,
    amount DECIMAL(20,8) NOT NULL,
    
    INDEX idx_order_lines_created_at (created_at),
    INDEX idx_order_lines_updated_at (updated_at),
    INDEX idx_order_lines_deleted_at (deleted_at),
    INDEX idx_order_lines_order_id (order_id),
    INDEX idx_order_lines_video_composition_id (video_composition_id),
    INDEX idx_order_lines_product_id (product_id),
    
    FOREIGN KEY (order_id) REFERENCES orders(id) ON UPDATE CASCADE ON DELETE RESTRICT,
    FOREIGN KEY (video_composition_id) REFERENCES video_compositions(id) ON UPDATE CASCADE ON DELETE RESTRICT,
    FOREIGN KEY (product_id) REFERENCES products(id) ON UPDATE CASCADE ON DELETE RESTRICT
);

-- +goose Down
DROP INDEX idx_order_lines_product_id ON order_lines;
DROP INDEX idx_order_lines_video_composition_id ON order_lines;
DROP INDEX idx_order_lines_order_id ON order_lines;
DROP INDEX idx_order_lines_deleted_at ON order_lines;
DROP INDEX idx_order_lines_updated_at ON order_lines;
DROP INDEX idx_order_lines_created_at ON order_lines;
DROP TABLE order_lines;
