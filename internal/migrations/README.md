# Database Migrations

This directory contains MySQL database migrations for the magic-video project.

## Requirements

Before running migrations, you need to install the `goose` tool:

```bash
go install github.com/pressly/goose/v3/cmd/goose@latest
```

## Configuration

Migrations use the following environment variables (with default values):

- `DB_HOST` - database host (default: localhost)
- `DB_PORT` - database port (default: 3306)
- `DB_USER` - database user (default: root)
- `DB_PASSWORD` - database password (default: password)
- `DB_NAME` - database name (default: magic_video)

## Usage

### Run all migrations (compose up)

```bash
cd internal/migrations
./compose_up.sh
```

Or with custom parameters:

```bash
DB_HOST=localhost DB_PORT=3306 DB_USER=myuser DB_PASSWORD=mypass DB_NAME=mydb ./compose_up.sh
```

### Rollback migrations (compose down)

Rollback all migrations:

```bash
cd internal/migrations
./compose_down.sh
```

Rollback specific number of migrations:

```bash
./compose_down.sh 3  # rollback last 3 migrations
```

## Migration Structure

Migrations are numbered and contain the following tables:

1. `001_create_customers_table.sql` - customers table
2. `002_create_customer_accesses_table.sql` - customer access tokens table
3. `003_create_video_compositions_table.sql` - video compositions table
4. `004_create_images_table.sql` - images table
5. `005_create_product_types_table.sql` - product types table
6. `006_create_products_table.sql` - products table
7. `007_create_orders_table.sql` - orders table
8. `008_create_order_lines_table.sql` - order lines table
9. `009_create_order_transactions_table.sql` - order transactions table
10. `010_create_order_payments_table.sql` - order payments table
11. `011_create_video_composition_jobs_table.sql` - video composition jobs table
12. `012_create_mail_logs_table.sql` - mail logs table

## Indexes

Each table contains appropriate indexes for:
- `created_at`, `updated_at`, `deleted_at` fields (soft delete)
- Foreign keys
- Unique fields (e.g., email, access_token)
- Frequently queried fields (status, template, etc.)

## Foreign Keys

Migrations contain appropriate foreign keys with constraints:
- `ON UPDATE CASCADE` - cascading updates
- `ON DELETE RESTRICT` - prevent deletion when references exist
- `ON DELETE SET NULL` - set NULL on deletion (for optional relationships)
