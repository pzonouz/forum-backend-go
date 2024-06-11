package utils

const CreateUserTableQuery = `CREATE TABLE IF NOT EXISTS "users" (
  "id" bigserial PRIMARY KEY,
  "email" varchar UNIQUE,
  "password" varchar,
  "first_name" varchar,
  "last_name" varchar,
  "address" text,
  "phone_number" varchar,
  "created_at" timestamp DEFAULT (now())
);`

const CreateUserTableQueryTest = `CREATE TABLE IF NOT EXISTS "users_test" (
  "id" bigserial PRIMARY KEY,
  "email" varchar UNIQUE,
  "password" varchar,
  "first_name" varchar,
  "last_name" varchar,
  "address" text,
  "phone_number" varchar,
  "created_at" timestamp DEFAULT (now())
);`

const CreateUserQuery = `INSERT INTO "users" ("email","password","first_name","last_name","address","phone_number") VALUES ($1,$2,$3,$4,$5,$6) RETURNING "id";`

const CreateUserQueryTest = `INSERT INTO "users_test" ("email","password","first_name","last_name","address","phone_number") VALUES ($1,$2,$3,$4,$5,$6) RETURNING "id";`

const EditUserQueryTest = `UPDATE users_test SET first_name=$1,last_name=$2,address=$3,phone_number=$4 WHERE id=$5`

const EditUserQuery = `UPDATE users SET first_name=$1,last_name=$2,address=$3,phone_number=$4 WHERE id=$5`

const DeleteUserByIDQuery = `DELETE FROM "users" WHERE "id"=$1`

const DeleteUserByIDQueryTest = `DELETE FROM "users_test" WHERE "id"=$1`

const DeleteTestTableQuery = `DROP TABLE IF EXISTS "%s_test";`
