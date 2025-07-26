-- +goose Up
CREATE TABLE order_transactions (
    id VARCHAR(255) PRIMARY KEY,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    order_id VARCHAR(255) NOT NULL,
    amount INT NOT NULL,
    method VARCHAR(100) NOT NULL,
    token VARCHAR(255) NOT NULL,
    session_iden VARCHAR(255) NOT NULL,
    transaction_iden VARCHAR(255) NOT NULL,
    payment_url TEXT NOT NULL,
    
    INDEX idx_order_transactions_created_at (created_at),
    INDEX idx_order_transactions_updated_at (updated_at),
    INDEX idx_order_transactions_deleted_at (deleted_at),
    INDEX idx_order_transactions_order_id (order_id),
    UNIQUE INDEX idx_order_transactions_session_iden (session_iden),
    UNIQUE INDEX idx_order_transactions_transaction_iden (transaction_iden),
    
    FOREIGN KEY (order_id) REFERENCES orders(id) ON UPDATE CASCADE ON DELETE RESTRICT
);

-- +goose Down
DROP INDEX idx_order_transactions_transaction_iden ON order_transactions;
DROP INDEX idx_order_transactions_session_iden ON order_transactions;
DROP INDEX idx_order_transactions_order_id ON order_transactions;
DROP INDEX idx_order_transactions_deleted_at ON order_transactions;
DROP INDEX idx_order_transactions_updated_at ON order_transactions;
DROP INDEX idx_order_transactions_created_at ON order_transactions;
DROP TABLE order_transactions;
