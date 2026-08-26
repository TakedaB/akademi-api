CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    email TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    role TEXT NOT NULL CHECK (role IN ('diretoria','financeiro','professor', 'aluno')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE students (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    enrollment_number TEXT NOT NULL UNIQUE,
    name TEXT NOT NULL,
    birth_date DATE NOT NULL,
    parent_name TEXT NOT NULL,
    city TEXT,
    phone TEXT NOT NULL,
    email TEXT,
    grade TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE enrollment_counters (
    year INT PRIMARY KEY,
    count INT NOT NULL DEFAULT 0
);

CREATE INDEX idx_students_enrollment_number ON students (enrollment_number);
CREATE INDEX idx_users_email ON users(email);    