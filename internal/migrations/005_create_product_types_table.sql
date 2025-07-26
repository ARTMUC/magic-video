-- +goose Up
CREATE TABLE product_types (
    id VARCHAR(255) PRIMARY KEY,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    product_id VARCHAR(255) NOT NULL,
    
    INDEX idx_product_types_created_at (created_at),
    INDEX idx_product_types_updated_at (updated_at),
    INDEX idx_product_types_deleted_at (deleted_at),
    INDEX idx_product_types_product_id (product_id),
    
    FOREIGN KEY (product_id) REFERENCES products(id) ON UPDATE CASCADE ON DELETE RESTRICT
);

-- +goose Down
DROP INDEX idx_product_types_product_id ON product_types;
DROP INDEX idx_product_types_deleted_at ON product_types;
DROP INDEX idx_product_types_updated_at ON product_types;
DROP INDEX idx_product_types_created_at ON product_types;
DROP TABLE product_types;
