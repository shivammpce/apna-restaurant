CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TABLE users (
    "id" UUID NOT NULL DEFAULT (uuid_generate_v4()),
    "email" VARCHAR NOT NULL,
    "name" VARCHAR NOT NULL,
    "phone_number" VARCHAR NOT NULL,
    "password" VARCHAR NOT NULL,
    "created_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP(3) NOT NULL,
    CONSTRAINT "users_pkey" PRIMARY KEY ("id")
);
CREATE UNIQUE INDEX "users_email_key" ON "users"("email");

