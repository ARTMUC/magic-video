-- +goose Up
CREATE TABLE order_payments (
    id VARCHAR(255) PRIMARY KEY,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    order_transaction_id VARCHAR(255) NOT NULL,
    order_id VARCHAR(255) NOT NULL,
    session_id VARCHAR(255) NOT NULL,
    
    INDEX idx_order_payments_created_at (created_at),
    INDEX idx_order_payments_updated_at (updated_at),
    INDEX idx_order_payments_deleted_at (deleted_at),
    INDEX idx_order_payments_order_transaction_id (order_transaction_id),
    INDEX idx_order_payments_order_id (order_id),
    UNIQUE INDEX idx_order_payments_session_id (session_id),
    
    FOREIGN KEY (order_transaction_id) REFERENCES order_transactions(id) ON UPDATE CASCADE ON DELETE RESTRICT,
    FOREIGN KEY (order_id) REFERENCES orders(id) ON UPDATE CASCADE ON DELETE RESTRICT
);

-- +goose Down
DROP INDEX idx_order_payments_session_id ON order_payments;
DROP INDEX idx_order_payments_order_id ON order_payments;
DROP INDEX idx_order_payments_order_transaction_id ON order_payments;
DROP INDEX idx_order_payments_deleted_at ON order_payments;
DROP INDEX idx_order_payments_updated_at ON order_payments;
DROP INDEX idx_order_payments_created_at ON order_payments;
DROP TABLE order_payments;
