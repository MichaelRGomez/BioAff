CREATE TABLE IF NOT EXISTS admin_users(
  id serial PRIMARY KEY,
  email text NOT NULL,
  au_password text NOT NULL,
  created_at TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT NOW()
);