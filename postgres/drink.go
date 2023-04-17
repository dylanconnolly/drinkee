package postgres

import (
	"encoding/json"

	"github.com/dylanconnolly/drinkee/drinkee"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type DrinkService struct {
	db *sqlx.DB
}

func NewDrinkService(db *sqlx.DB) *DrinkService {
	return &DrinkService{db: db}
}

func (s *DrinkService) FindDrinkByID(c *gin.Context, id int) (*drinkee.DrinkResponse, error) {
	tx, err := s.db.BeginTxx(c, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	drinks, err := findDrinkByID(c, tx, id)
	if err != nil {
		return nil, err
	}

	return drinks, nil
}

func (s *DrinkService) FindDrinks(ctx *gin.Context) ([]*drinkee.DrinkResponse, error) {
	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	drinks, err := findDrinks(ctx, tx)
	if err != nil {
		return nil, err
	}

	return drinks, nil
}

func (s *DrinkService) CreateDrink(c *gin.Context, cd *drinkee.CreateDrink) error {
	tx, err := s.db.BeginTxx(c, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err := createDrink(c, tx, cd); err != nil {
		return err
	}

	return tx.Commit()
}

func (s *DrinkService) GenerateDrinks(c *gin.Context, i []drinkee.Ingredient) ([]drinkee.DrinkResponse, error) {
	var ingredientIDs []int

	tx, err := s.db.BeginTxx(c, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	for _, ingredient := range i {
		ingredientIDs = append(ingredientIDs, ingredient.ID)
	}

	drinks, err := generateDrinks(c, tx, ingredientIDs)
	if err != nil {
		return nil, err
	}

	return drinks, nil
}

func createDrink(c *gin.Context, tx *sqlx.Tx, cd *drinkee.CreateDrink) error {
	var ingredientNames []string
	for _, di := range cd.DrinkIngredients {
		ingredientNames = append(ingredientNames, di.Name)
	}

	diJSON, err := json.Marshal(cd.DrinkIngredients)

	_, err = tx.Exec(`
		WITH drink AS (
			INSERT INTO drinks (name, display_name, description, instructions)
			VALUES ($1, $2, $3, $4)
			RETURNING id
		),
		ingredient_ids AS (
			SELECT id, name FROM ingredients WHERE name = ANY($5)
		),
		ingredient_data AS (
			SELECT * FROM json_populate_recordset(null::ingredient_data, $6)
		)
		INSERT INTO drink_ingredients (drink_id, ingredient_id, measurement)
		SELECT drink.id, ingredient_ids.id, ingredient_data.measurement
		FROM drink, ingredient_ids, ingredient_data
		WHERE ingredient_ids.name = ingredient_data.name
	`, cd.Name, cd.DisplayName, cd.Description, cd.Instructions, pq.Array(ingredientNames), string(diJSON))

	if err != nil {
		return err
	}

	return nil
}

func findDrinks(ctx *gin.Context, tx *sqlx.Tx) ([]*drinkee.DrinkResponse, error) {
	var drinks []*drinkee.DrinkResponse

	queryStr := `
	SELECT d.id, d.name, d.display_name, d.description, d.instructions, json_agg(json_build_object('name', i.name, 'displayName', i.display_name, 'measurement', di.measurement)) as drink_ingredients 
	FROM drinks d 
	JOIN drink_ingredients di ON di.drink_id=d.id
	JOIN ingredients i ON di.ingredient_id=i.id 
	GROUP BY d.id, d.name ORDER BY d.name;
	`

	err := tx.Select(&drinks, queryStr)

	if err != nil {
		return nil, err
	}

	return drinks, nil
}

func generateDrinks(ctx *gin.Context, tx *sqlx.Tx, ingredientIDs []int) ([]drinkee.DrinkResponse, error) {
	var drinks []drinkee.DrinkResponse

	queryStr := `SELECT md.id,md.name,md.display_name,md.description,md.instructions, ij.drink_ingredients
		FROM 
			(SELECT d.*, COUNT(*) AS ingredients_present,
			(SELECT COUNT(*) FROM drink_ingredients WHERE drink_ingredients.drink_id=d.id) AS total_ingredients 
			FROM drinks d JOIN drink_ingredients di ON di.drink_id=d.id WHERE di.ingredient_id = ANY($1) GROUP BY d.id) AS md 
      JOIN (SELECT d.id, json_agg(json_build_object('name', i.name, 'displayName', i.display_name, 'measurement', di.measurement)) as drink_ingredients 
            FROM drinks d 
            JOIN drink_ingredients di ON di.drink_id=d.id
            JOIN ingredients i ON di.ingredient_id=i.id 
            GROUP BY d.id, d.name ) AS ij ON ij.id=md.id
		WHERE ingredients_present=total_ingredients
		ORDER BY md.name;`

	err := tx.Select(&drinks, queryStr, pq.Array(ingredientIDs))
	if err != nil {
		return nil, err
	}

	return drinks, nil
}

func findDrinkByID(c *gin.Context, tx *sqlx.Tx, id int) (*drinkee.DrinkResponse, error) {
	var drink drinkee.DrinkResponse

	err := tx.Get(&drink, `
	SELECT d.id, d.name, d.display_name, d.description, d.instructions, json_agg(json_build_object('name', i.name, 'displayName', i.display_name, 'measurement', di.measurement)) as drink_ingredients 
	FROM drinks d 
	JOIN drink_ingredients di ON di.drink_id=d.id
	JOIN ingredients i ON di.ingredient_id=i.id 
	WHERE d.id = $1
	GROUP BY d.id, d.name ORDER BY d.name
	`, id)

	if err != nil {
		return nil, err
	}

	return &drink, nil
}

func (s *DrinkService) FindIngredients(c *gin.Context) ([]*drinkee.Ingredient, error) {
	tx, err := s.db.BeginTxx(c, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	ingredients, err := findIngredients(c, tx)
	if err != nil {
		return nil, err
	}

	return ingredients, nil
}

func findIngredients(c *gin.Context, tx *sqlx.Tx) ([]*drinkee.Ingredient, error) {
	var ingredients []*drinkee.Ingredient

	err := tx.Select(&ingredients, "SELECT id, name, display_name FROM ingredients ORDER BY name")
	if err != nil {
		return nil, err
	}

	return ingredients, nil
}
