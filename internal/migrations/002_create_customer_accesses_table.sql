-- +goose Up
CREATE TABLE customer_accesses (
    id VARCHAR(255) PRIMARY KEY,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    customer_id VARCHAR(255) NOT NULL,
    access_token VARCHAR(255) NOT NULL,
    token_expire_date DATETIME NOT NULL,
    
    INDEX idx_customer_accesses_created_at (created_at),
    INDEX idx_customer_accesses_updated_at (updated_at),
    INDEX idx_customer_accesses_deleted_at (deleted_at),
    INDEX idx_customer_accesses_customer_id (customer_id),
    UNIQUE INDEX idx_customer_accesses_access_token (access_token),
    
    FOREIGN KEY (customer_id) REFERENCES customers(id) ON UPDATE CASCADE ON DELETE RESTRICT
);

-- +goose Down
DROP INDEX idx_customer_accesses_access_token ON customer_accesses;
DROP INDEX idx_customer_accesses_customer_id ON customer_accesses;
DROP INDEX idx_customer_accesses_deleted_at ON customer_accesses;
DROP INDEX idx_customer_accesses_updated_at ON customer_accesses;
DROP INDEX idx_customer_accesses_created_at ON customer_accesses;
DROP TABLE customer_accesses;
