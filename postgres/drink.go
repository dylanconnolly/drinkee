package postgres

import (
	"github.com/dylanconnolly/drinkee/drinkee"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
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
