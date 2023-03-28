CREATE TABLE IF NOT EXISTS public_user (
    id serial PRIMARY KEY,
    email text NOT NULL,
    pu_password text NOT NULL
);