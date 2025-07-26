-- +goose Up
CREATE TABLE orders (
    id VARCHAR(255) PRIMARY KEY,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    customer_id VARCHAR(255) NOT NULL,
    gross_amount DECIMAL(10,2) NOT NULL,
    net_amount DECIMAL(10,2) NOT NULL,
    tax_amount DECIMAL(10,2) NOT NULL,
    tax_breakdown TEXT NOT NULL,
    status VARCHAR(255) NOT NULL,
    payment_status VARCHAR(255) NOT NULL,
    idempotency_key VARCHAR(255) NOT NULL,
    
    INDEX idx_orders_created_at (created_at),
    INDEX idx_orders_updated_at (updated_at),
    INDEX idx_orders_deleted_at (deleted_at),
    INDEX idx_orders_customer_id (customer_id),
    INDEX idx_orders_status (status),
    INDEX idx_orders_payment_status (payment_status),
    UNIQUE INDEX idx_orders_idempotency_key (idempotency_key),
    
    FOREIGN KEY (customer_id) REFERENCES customers(id) ON UPDATE CASCADE ON DELETE RESTRICT
);

-- +goose Down
DROP INDEX idx_orders_idempotency_key ON orders;
DROP INDEX idx_orders_payment_status ON orders;
DROP INDEX idx_orders_status ON orders;
DROP INDEX idx_orders_customer_id ON orders;
DROP INDEX idx_orders_deleted_at ON orders;
DROP INDEX idx_orders_updated_at ON orders;
DROP INDEX idx_orders_created_at ON orders;
DROP TABLE orders;
