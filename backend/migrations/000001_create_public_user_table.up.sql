-- Filename: BIOAFF/backend/migrations/000001_create_public_user_table.up.sql

--creating extension for the email type
create extension if not exists citext;

create table if not exists public_user(
    id serial primary key,
    name text not null,
    email citext unique not null,
    pu_password bytea not null,
    activated bool not null,
    created_at timestamp(0) with time zone not null default now(),
    version integer not null default 1
);