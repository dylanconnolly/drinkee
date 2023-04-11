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

func (s *DrinkService) SimpleFind() {
	s.db.Queryx("SELECT * FROM drinks")
	s.db.Queryx(`\dt`)
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
