CREATE TABLE IF NOT EXISTS drink_ingredients(
    id serial PRIMARY KEY,
    drink_id int REFERENCES drinks(id) ON DELETE CASCADE,
    ingredient_id int REFERENCES ingredients(id) ON DELETE CASCADE,
    measurement text NOT NULL
);