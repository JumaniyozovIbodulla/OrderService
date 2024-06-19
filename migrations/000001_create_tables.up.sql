CREATE TYPE "order_enum" AS ENUM('self_pickup','delivery');
CREATE TYPE "payment_type" AS ENUM('uzum','cash', 'terminal');
CREATE TYPE "payment_enum" AS ENUM('waiting_for_payment', 'collecting', 'shipping', 'waiting_on_branch', 'finished', 'cancelled');



CREATE TABLE IF NOT EXISTS "orders" (
    "id" UUID PRIMARY KEY,
    "external_id" VARCHAR NOT NULL,
    "type" order_enum,
    "customer_phone" VARCHAR(20) NOT NULL,
    "customer_name" VARCHAR(20) NOT NULL,
    "customer_id" UUID,
    "payment_type" payment_type,
    "status" payment_enum,
    "to_address" VARCHAR NOT NULL,
    "to_location" GEOMETRY(POINT, 4326) NOT NULL,
    "discount_amount" DECIMAL(10, 2) DEFAULT 0.0,
    "amount" DECIMAL(10, 2) DEFAULT 0.0,
    "delivery_price" DECIMAL(10, 2) DEFAULT 0.0,
    "paid" BOOLEAN DEFAULT FALSE,
    "courier_id" UUID,
    "courier_phone" VARCHAR(20) NOT NULL,
    "courier_name" VARCHAR(20) NOT NULL,
    "created_at" TIMESTAMPTZ DEFAULT NOW(),
    "updated_at" TIMESTAMPTZ DEFAULT NOW(),
    "deleted_at" BIGINT DEFAULT 0,
    CONSTRAINT "unique_user_phone_deleted_at" UNIQUE ("customer_phone", "deleted_at")
);

UPDATE "orders"
SET "to_location" = ST_GeomFromText("to_location"::text, 4326);


CREATE TABLE IF NOT EXISTS "order_products" (
    "id" UUID PRIMARY KEY,
    "product_id" UUID NOT NULL,
    "count" INT DEFAULT 1,
    "discount_price" DECIMAL(10, 2) DEFAULT 0.0,
    "price" DECIMAL(10, 2) DEFAULT 0.0,
    "order_id" UUID REFERENCES "orders"("id"),
    "created_at" TIMESTAMPTZ DEFAULT NOW(),
    "updated_at" TIMESTAMPTZ DEFAULT NOW(),
    "deleted_at" BIGINT DEFAULT 0
);

CREATE TABLE IF NOT EXISTS "order_status_notes" (
    "id" UUID PRIMARY KEY,
    "order_id" UUID REFERENCES "orders"("id"),
    "status" payment_enum,
    "user_id" UUID,
    "reason" VARCHAR(100) NOT NULL,
    "created_at" TIMESTAMPTZ DEFAULT NOW()
);