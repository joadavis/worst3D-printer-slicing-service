-- initial database schema for worst3dprintservice
-- joadavis Nov 2020

CREATE EXTENSION pgcrypto;

CREATE TABLE projects
(
    id SERIAL,
    project_name TEXT NOT NULL
);


CREATE TABLE users
(
    id uuid NOT NULL DEFAULT gen_random_uuid() PRIMARY KEY,
    username TEXT NOT NULL,
    email TEXT NOT NULL,
    project_id INT,
    password TEXT NOT NULL,
    FOREIGN KEY (project_id) REFERENCES projects (project_id)
);

CREATE TABLE jobs
(
    id uuid NOT NULL DEFAULT gen_random_uuid() PRIMARY KEY,
    status TEXT,
    requesting_user_id uuid,
    input_file_path TEXT,
    output_path TEXT
    FOREIGN KEY (requesting_user_id) REFERENCES users (id)
);

