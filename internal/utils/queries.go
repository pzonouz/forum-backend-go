package utils

const CreateQuestionTableQuery = `CREATE TABLE IF NOT EXISTS "questions" (
  "id" bigserial PRIMARY KEY,
  "title" varchar UNIQUE,
  "description" text,
  "created_at" timestamp DEFAULT (now())
);`

const CreateQuestionTableQueryTest = `CREATE TABLE IF NOT EXISTS "questions_test" (
  "id" bigserial PRIMARY KEY,
  "title" varchar UNIQUE,
  "description" text,
  "created_at" timestamp DEFAULT (now())
);`

const CreateUserTableQuery = `
  -- CREATE TYPE IF NOT EXISTS roleType as ENUM('admin','user','owner','modifier');
  CREATE TABLE IF NOT EXISTS users (
  id bigserial PRIMARY KEY,
  email varchar UNIQUE NOT NULL,
  password varchar NOT NULL,
  name varchar,
  address text,
  phone_number varchar, 
  role roleType NOT NULL DEFAULT 'user',
  created_at timestamp DEFAULT (now())
);`

const CreateUserTableQueryTest = `
  -- CREATE TYPE roleType as ENUM('admin','user','owner','modifier');
  CREATE TABLE IF NOT EXISTS "users_test" (
  "id" bigserial PRIMARY KEY,
  "email" varchar UNIQUE,
  "password" varchar,
  "name" varchar,
  "address" text,
  "phone_number" varchar,
  role roleType NOT NULL DEFAULT 'user',
  "created_at" timestamp DEFAULT (now())
);`

// const CreateUserQuery = `INSERT INTO "users" ("email","password","name","address","phone_number") VALUES ($1,$2,$3,$4,$5,$6) RETURNING "id";`
//
// const CreateUserQueryTest = `INSERT INTO "users_test" ("email","password",name","address","phone_number") VALUES ($1,$2,$3,$4,$5,$6) RETURNING "id";`
//
// const EditUserQueryTest = `UPDATE users_test SET first_name=$1,last_name=$2,address=$3,phone_number=$4 WHERE id=$5`
//
// const EditUserQuery = `UPDATE users SET first_name=$1,last_name=$2,address=$3,phone_number=$4 WHERE id=$5`
//
// const DeleteUserByIDQuery = `DELETE FROM "users" WHERE "id"=$1`
//
// const DeleteUserByIDQueryTest = `DELETE FROM "users_test" WHERE "id"=$1`
const DeleteTestTableQuery = `DROP TABLE IF EXISTS "%s_test";`
