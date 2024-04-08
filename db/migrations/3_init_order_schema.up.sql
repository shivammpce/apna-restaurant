CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS tables (
    "id" UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    "table_number" INT NOT NULL,
    "order_ids" UUID[],
    CONSTRAINT "unique_table_number" UNIQUE ("table_number") 
);

CREATE TABLE IF NOT EXISTS orders (
    "id" UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    "table_id" UUID NOT NULL,
    "amount" INT,
    "order_items" JSONB, 
    "created_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "delivered_at" TIMESTAMP(3),
    FOREIGN KEY ("table_id") REFERENCES tables ("id")
);
