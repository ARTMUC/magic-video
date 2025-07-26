-- +goose Up
CREATE TABLE mail_logs (
    id VARCHAR(255) PRIMARY KEY,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    recipient_name VARCHAR(255) NOT NULL,
    recipient_email VARCHAR(255) NOT NULL,
    template VARCHAR(255) NOT NULL,
    status VARCHAR(255) NOT NULL,
    error TEXT NULL,
    reference VARCHAR(255) NOT NULL,
    
    INDEX idx_mail_logs_created_at (created_at),
    INDEX idx_mail_logs_updated_at (updated_at),
    INDEX idx_mail_logs_deleted_at (deleted_at),
    INDEX idx_mail_logs_recipient_email (recipient_email),
    INDEX idx_mail_logs_template (template),
    INDEX idx_mail_logs_status (status),
    INDEX idx_mail_logs_reference (reference)
);

-- +goose Down
DROP INDEX idx_mail_logs_reference ON mail_logs;
DROP INDEX idx_mail_logs_status ON mail_logs;
DROP INDEX idx_mail_logs_template ON mail_logs;
DROP INDEX idx_mail_logs_recipient_email ON mail_logs;
DROP INDEX idx_mail_logs_deleted_at ON mail_logs;
DROP INDEX idx_mail_logs_updated_at ON mail_logs;
DROP INDEX idx_mail_logs_created_at ON mail_logs;
DROP TABLE mail_logs;
