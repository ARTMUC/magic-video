-- +goose Up
CREATE TABLE active_products (
    id VARCHAR(255) PRIMARY KEY,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    product_id VARCHAR(255) NOT NULL,

    INDEX idx_active_products_created_at (created_at),
    INDEX idx_active_products_updated_at (updated_at),
    INDEX idx_active_products_deleted_at (deleted_at),
    INDEX idx_active_products_product_id (product_id),

    FOREIGN KEY (product_id) REFERENCES products(id) ON UPDATE CASCADE ON DELETE RESTRICT
);

-- +goose Down
DROP INDEX idx_active_products_created_at ON active_products;
DROP INDEX idx_active_products_updated_at ON active_products;
DROP INDEX idx_active_products_deleted_at ON active_products;
DROP INDEX idx_active_products_product_id ON active_products;
DROP TABLE active_products;



