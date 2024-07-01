package utils

const CreateUserTableQuery = `
  -- CREATE TYPE roleType as ENUM('admin','user','owner','modifier');
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

const CreateQuestionTableQuery = `CREATE TABLE IF NOT EXISTS "questions" (
  "id" bigserial PRIMARY KEY,
  "title" varchar UNIQUE,
  "description" text,
  "created_at" timestamp DEFAULT (now()),
  "user_name" varchar,
  "user_id" bigint REFERENCES users(id)
);`

const CreateQuestionTableQueryTest = `CREATE TABLE IF NOT EXISTS "questions_test" (
  "id" bigserial PRIMARY KEY,
  "title" varchar UNIQUE,
  "description" text,
  "created_at" timestamp DEFAULT (now()),
  "user_name" varchar,
  "user_id" bigint REFERENCES users(id)
);`

const CreateAnswerTableQuery = `CREATE TABLE IF NOT EXISTS "answers" (
  "id" bigserial PRIMARY KEY,
  "description" text,
  "created_at" timestamp DEFAULT (now()),
  "user_name" varchar,
  "user_id" bigint REFERENCES users(id),
  "question_id" bigint REFERENCES questions(id)
);`

const CreateAnswerTableQueryTest = `CREATE TABLE IF NOT EXISTS "answers_test" (
  "id" bigserial PRIMARY KEY,
  "description" text,
  "created_at" timestamp DEFAULT (now()),
  "user_name" varchar,
  "user_id" bigint REFERENCES users(id),
  "question_id" bigint REFERENCES questions(id)
);`

const CreateScoreTableQuery = `
  --CREATE TYPE Operator as ENUM('plus','minus');
  CREATE TABLE IF NOT EXISTS "scores" (
  "id" bigserial PRIMARY KEY,
  "operator" Operator NOT NULL,
  "created_at" timestamp DEFAULT (now()),
  "user_id" bigint REFERENCES users(id),
  "question_id" bigint REFERENCES questions(id),
  "answer_id" bigint REFERENCES answers(id)
);`

const CreateScoreTableQueryTest = `
  --CREATE TYPE Operator as ENUM('plus','minus');
  CREATE TABLE IF NOT EXISTS "scores_test" (
  "id" bigserial PRIMARY KEY,
  "operator" Operator NOT NULL,
  "created_at" timestamp DEFAULT (now()),
  "user_id" bigint REFERENCES users(id),
  "question_id" bigint REFERENCES questions(id),
  "answer_id" bigint REFERENCES answers(id)
  );`

const DeleteTestTableQuery = `DROP TABLE IF EXISTS "%s_test";`
