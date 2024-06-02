package utils

const CreateUserTableQuery = `CREATE TABLE IF NOT EXISTS "users" (
  "id" bigserial PRIMARY KEY,
  "email" varchar,
  "password" varchar,
  "first_name" varchar,
  "last_name" varchar,
  "phone_number" varchar,
  "created_at" timestamp DEFAULT (now())
);`

const CreateUserTableQueryTest = `CREATE TABLE IF NOT EXISTS "users_test" (
  "id" bigserial PRIMARY KEY,
  "email" varchar,
  "password" varchar,
  "first_name" varchar,
  "last_name" varchar,
  "phone_number" varchar,
  "created_at" timestamp DEFAULT (now())
);`

const DeleteUserTableQueryTest = `DROP TABLE IF EXISTS "users_test";`
