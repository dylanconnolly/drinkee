CREATE TYPE ingredient_data AS (name VARCHAR(255), measurement VARCHAR(255));

CREATE TABLE IF NOT EXISTS drink_ingredients(
    id serial PRIMARY KEY,
    drink_id int REFERENCES drinks(id) ON DELETE CASCADE,
    ingredient_id int REFERENCES ingredients(id) ON DELETE CASCADE,
    measurement text NOT NULL,
    created_at timestamp NOT NULL DEFAULT current_timestamp,
    updated_at timestamp NOT NULL DEFAULT current_timestamp
);

CREATE TRIGGER set_timestamp
BEFORE UPDATE ON drink_ingredients
FOR EACH ROW
EXECUTE PROCEDURE set_updated_at();