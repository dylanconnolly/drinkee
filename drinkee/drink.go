package drinkee

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
)

type Drink struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	DisplayName  string `json:"displayName" db:"display_name"`
	Description  string `json:"description,omitempty"`
	Instructions string `json:"instructions"`
}

type DrinkResponse struct {
	Drink
	DrinkIngredients DrinkIngredientSlice `json:"drinkIngredients" db:"drink_ingredients"`
}

type DrinkIngredientSlice []DrinkIngredient

func (dis *DrinkIngredientSlice) Scan(src interface{}) error {
	var data []byte
	switch v := src.(type) {
	case string:
		data = []byte(v)
	case []byte:
		data = v
	default:
		return nil
	}
	return json.Unmarshal(data, dis)
}

type DrinkIngredient struct {
	Name        string `json:"name"`
	DisplayName string `json:"displayName" db:"display_name"`
	Measurement string `json:"measurement"`
}

type DrinkService interface {
	FindDrinkByID(ctx *gin.Context, id int) (*Drink, error)
	FindDrinks(ctx *gin.Context) ([]*DrinkResponse, error)
	CreateDrink(ctx *gin.Context) error
}
