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

CREATE TABLE teachers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,
    subject TEXT NOT NULL,
    phone TEXT NOT NULL,
    hire_date DATE NOT NULL,
    class_assigned TEXT NOT NULL,
    workload_hours INT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE finance (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    student_id UUID NOT NULL REFERENCES students(id) ON DELETE CASCADE,
    description TEXT NOT NULL,
    amount NUMERIC(10, 2) NOT NULL,
    payment_method TEXT NOT NULL CHECK (payment_method IN ('cash', 'pix', 'cartao', 'boleto', 'transferencia')),
    status TEXT NOT NULL DEFAULT 'pendente' CHECK (status IN ( 'pendente', 'pago', 'atrasado', 'isento')),
    due_date DATE NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_finance_student_id ON finance(student_id);