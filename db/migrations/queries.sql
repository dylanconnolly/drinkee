-- insert drink then get id of new record
-- lookup ingredient IDs off list of names
-- add drink ID and ingredient IDs to drink_ingredients table

WITH drink AS (
  INSERT INTO drinks (name, description, instructions)
  VALUES (
    'Margarita',
    'A classic tequila-based cocktail',
    'Shake all ingredients with ice and strain into a chilled cocktail glass.'
  )
  RETURNING id
),
ingredient_ids AS (
  SELECT id FROM ingredients WHERE name IN ('Tequila', 'Triple sec', 'Lime juice', 'Ice')
)
INSERT INTO drink_ingredients (drink_id, ingredient_id)
SELECT drink.id, ingredient_ids.id
FROM drink, ingredient_ids;


-- keep incoming drink ingredients as json and create a table using the array of ingredients and measurements
WITH drink AS (
    INSERT INTO drinks (name, description, instructions)
    VALUES ($1, $2, $3)
    RETURNING id
),
ingredient_ids AS (
    SELECT id, name FROM ingredients WHERE name = ANY($4)
),
ingredient_data AS (
    SELECT name, measurement FROM json_populate_recordset(null::ingredient, $4)
)
INSERT INTO drink_ingredients (drink_id, ingredient_id, measurement)
SELECT drink.id, ingredient_ids.id, ingredient_data.measurement
FROM drink, ingredient_ids, ingredient_data
WHERE ingredient_ids.name = ingredient_data.name


-- using postgres arrays
WITH drink AS (
    INSERT INTO drinks (name, description, instructions)
    VALUES ($1, $2, $3)
    RETURNING id
),
ingredient_ids AS (
    SELECT id, name FROM ingredients WHERE name = ANY($4)
)
INSERT INTO drink_ingredients (drink_id, ingredient_id, measurement)
SELECT drink.id, ingredient_ids.id, ingredient_data.measurement
FROM drink, ingredient_ids
JOIN UNNEST($5::ingredient_data[]) AS ingredient_data ON ingredient_data.name = ingredient_ids.name`, dr.Name, dr.Description, dr.Instructions, pq.Array(ingredientNames), pq.Array(ingredientMeasurements))