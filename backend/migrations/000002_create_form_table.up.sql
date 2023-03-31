-- Filename: BIOAFF/backend/migrations/000002_create_form_table.up.sql

create table if not exists form(
    user_id serial primary key references public_user(id) not null,
    form_id serial not null,
    form_status text not null,
    archive_status boolean not null,
    full_name text not null,
    other_names text,
    name_change_status text,
    social_security_num int not null,
    social_security_date date not null,
    social_security_country text not null,
    passport_number text not null,
    passport_date date not null,
    passport_country text not null,
    dob date not null,
    place_of_birth text not null,
    nationality text not null,
    acquired_nationality text not null,
    spouse_name text,
    address text not null,
    residential_phone_number text not null,
    residential_fax_number text,
    residentia_email citext,
    created_on timestamp(0) with time zone not null default now()
);