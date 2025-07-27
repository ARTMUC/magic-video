-- +goose Up
CREATE TABLE products (
    id VARCHAR(255) PRIMARY KEY,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    name VARCHAR(100) NOT NULL,
    unit_price DECIMAL(10,2) NOT NULL,
    tax_rate DECIMAL(5,2) NOT NULL,
    
    INDEX idx_products_created_at (created_at),
    INDEX idx_products_updated_at (updated_at),
    INDEX idx_products_deleted_at (deleted_at),
    INDEX idx_products_name (name)

);

-- +goose Down
DROP INDEX idx_products_name ON products;
DROP INDEX idx_products_deleted_at ON products;
DROP INDEX idx_products_updated_at ON products;
DROP INDEX idx_products_created_at ON products;
DROP TABLE products;
