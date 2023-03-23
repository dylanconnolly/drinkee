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