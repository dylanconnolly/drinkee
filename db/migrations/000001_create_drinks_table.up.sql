CREATE TABLE IF NOT EXISTS drinks(
    id serial PRIMARY KEY,
    name text NOT NULL,
    description text,
    instructions text
);