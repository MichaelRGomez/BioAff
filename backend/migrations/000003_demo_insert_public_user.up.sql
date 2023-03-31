--Filename: BIOAFF/backend/internal/migrations/000003_demo_insert_public_user.up.sql

insert into public_user
(name, email, pu_password, activated)
values ('john doe', 'jd@gmail.com', 'null',true);