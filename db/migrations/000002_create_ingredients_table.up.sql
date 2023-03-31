CREATE TABLE IF NOT EXISTS ingredients(
    id serial PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    display_name VARCHAR(255) NOT NULL,
    created_at timestamp NOT NULL DEFAULT current_timestamp,
    updated_at timestamp NOT NULL DEFAULT current_timestamp
);

CREATE TRIGGER set_timestamp
BEFORE UPDATE ON ingredients
FOR EACH ROW
EXECUTE PROCEDURE set_updated_at();