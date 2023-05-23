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
JOIN UNNEST($5::ingredient_data[]) AS ingredient_data ON ingredient_data.name = ingredient_ids.name, dr.Name, dr.Description, dr.Instructions, pq.Array(ingredientNames), pq.Array(ingredientMeasurements))


-- from list of ingredient IDs return drinks containing at least 1 of those ingredients and append a column (n) to show number of ingredients present in list
SELECT d.*, COUNT(*) AS N FROM drinks d JOIN drink_ingredients di ON di.drink_id=d.id JOIN ingredients i ON di.ingredient_id=i.id WHERE i.id IN (2,7,8) GROUP BY d.id;

-- select drinks only that have the exact ingredients in list (can't have extra ingredients)
SELECT d.* FROM drinks d WHERE Not Exists (SELECT 1 FROM ingredients i WHERE id IN (2,7,8) AND Not Exists (SELECT 1 FROM drink_ingredients di WHERE di.drink_id=d.id AND di.ingredient_id=i.id));

-- get total ingredient count alongside number of ingredients we have available
SELECT d.*, COUNT(*) AS ingredient_count, (SELECT COUNT(*) FROM drink_ingredients WHERE drink_ingredients.drink_id=d.id) FROM drinks d JOIN drink_ingredients di ON di.drink_id=d.id WHERE di.ingredient_id IN (2,7,8) GROUP BY d.id;

-- return only drinks that we can make
SELECT * FROM (SELECT d.*, COUNT(*) AS ic, (SELECT COUNT(*) FROM drink_ingredients WHERE drink_ingredients.drink_id=d.id) AS ti FROM drinks d JOIN drink_ingredients di ON di.drink_id=d.id WHERE di.ingredient_id IN (2,7,8) GROUP BY d.id) AS joiny WHERE ic=ti;




SELECT md.id,md.name,md.display_name,md.description,md.instructions, ij.drink_ingredients
		FROM 
			(SELECT d.*, COUNT(*) AS ingredients_present,
			(SELECT COUNT(*) FROM drink_ingredients WHERE drink_ingredients.drink_id=d.id) AS total_ingredients 
			FROM drinks d JOIN drink_ingredients di ON di.drink_id=d.id WHERE di.ingredient_id IN (15,24) GROUP BY d.id) AS md 
      JOIN (SELECT d.id, json_agg(json_build_object('name', i.name, 'displayName', i.display_name, 'measurement', di.measurement)) as drink_ingredients 
            FROM drinks d 
            JOIN drink_ingredients di ON di.drink_id=d.id
            JOIN ingredients i ON di.ingredient_id=i.id 
            GROUP BY d.id, d.name ) AS ij ON ij.id=md.id
		WHERE ingredients_present=total_ingredients;



SELECT md.id,md.name,md.display_name,md.description,md.instructions, ij.drink_ingredients, ingredients_present, total_ingredients - ingredients_present AS missing_ingredients
		FROM 
			(SELECT d.*, COUNT(*) AS ingredients_present,
			(SELECT COUNT(*) FROM drink_ingredients WHERE drink_ingredients.drink_id=d.id) AS total_ingredients 
			FROM drinks d JOIN drink_ingredients di ON di.drink_id=d.id WHERE di.ingredient_id IN (15,24) GROUP BY d.id) AS md 
      JOIN (SELECT d.id, json_agg(json_build_object('name', i.name, 'displayName', i.display_name, 'measurement', di.measurement)) as drink_ingredients 
            FROM drinks d 
            JOIN drink_ingredients di ON di.drink_id=d.id
            JOIN ingredients i ON di.ingredient_id=i.id 
            GROUP BY d.id, d.name ) AS ij ON ij.id=md.id
		WHERE ingredients_present>=1;