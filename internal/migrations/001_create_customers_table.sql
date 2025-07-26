-- +goose Up
CREATE TABLE customers (
    id VARCHAR(255) PRIMARY KEY,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    email VARCHAR(255) NOT NULL,
    name VARCHAR(255) NULL,
    address_home_number VARCHAR(255) NULL,
    address_street_name VARCHAR(255) NULL,
    address_city VARCHAR(255) NULL,
    address_zip_code VARCHAR(255) NULL,
    address_state VARCHAR(255) NULL,
    address_country_code VARCHAR(255) NULL,
    
    INDEX idx_customers_created_at (created_at),
    INDEX idx_customers_updated_at (updated_at),
    INDEX idx_customers_deleted_at (deleted_at),
    UNIQUE INDEX idx_customers_email (email)
);

-- +goose Down
DROP INDEX idx_customers_email ON customers;
DROP INDEX idx_customers_deleted_at ON customers;
DROP INDEX idx_customers_updated_at ON customers;
DROP INDEX idx_customers_created_at ON customers;
DROP TABLE customers;
