-- initial database schema for worst3dprintservice
-- joadavis Nov 2020
-- NOTE: THIS MUST match internal/db_model.go!!!

CREATE EXTENSION pgcrypto;

CREATE TABLE IF NOT EXISTS projects
(
    id SERIAL PRIMARY KEY,
    project_name TEXT NOT NULL
);


CREATE TABLE IF NOT EXISTS users
(
    id uuid NOT NULL DEFAULT gen_random_uuid() PRIMARY KEY,
    username TEXT NOT NULL,
    email TEXT NOT NULL,
    password TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS jobs
(
    id uuid NOT NULL DEFAULT gen_random_uuid() PRIMARY KEY,
    status TEXT,
    requesting_user_id uuid,
    project_id int,
    input_file_path TEXT,
    output_path TEXT,
    FOREIGN KEY (requesting_user_id) REFERENCES users (id),
    FOREIGN KEY (project_id) references projects (id)
);

CREATE TABLE IF NOT EXISTS roles
(
    id int PRIMARY KEY,
    role_name TEXT NOT NULL
);


-- Specify the relation for an RBAC
CREATE TABLE IF NOT EXISTS user_project_roles
(
    id uuid NOT NULL DEFAULT gen_random_uuid() PRIMARY KEY,
    user_id uuid,
    project_id int,
    role_id int,
    FOREIGN KEY (user_id) REFERENCES users (id),
    FOREIGN KEY (project_id) REFERENCES projects (id),
    FOREIGN KEY (role_id) REFERENCES roles (id)
);


INSERT INTO roles (id, role_name) VALUES (13, 'admin');
INSERT INTO roles (id, role_name) VALUES (50, 'user');