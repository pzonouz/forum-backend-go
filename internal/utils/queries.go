package utils

const CreateRoleTableQuery = `CREATE TABLE IF NOT EXISTS "roles" (
  "id" bigserial PRIMARY KEY,
  "name" varchar UNIQUE,
  "created_at" timestamp DEFAULT (now())
);`

const CreateRoleTableQueryTest = `CREATE TABLE IF NOT EXISTS "roles_test" (
  "id" bigserial PRIMARY KEY,
  "name" varchar UNIQUE,
  "created_at" timestamp DEFAULT (now())
);`

const CreateUserTableQuery = `CREATE TABLE IF NOT EXISTS "users" (
  "id" bigserial PRIMARY KEY,
  "email" varchar UNIQUE,
  "password" varchar,
  "name" varchar,
  "address" text,
  "phone_number" varchar,
  "created_at" timestamp DEFAULT (now()),
  "roler_id" integer REFERENCES roles(id)
);`

const CreateUserTableQueryTest = `CREATE TABLE IF NOT EXISTS "users_test" (
  "id" bigserial PRIMARY KEY,
  "email" varchar UNIQUE,
  "password" varchar,
  "name" varchar,
  "address" text,
  "phone_number" varchar,
  "created_at" timestamp DEFAULT (now()),
  "roler_id" integer REFERENCES roles_test(id)
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
