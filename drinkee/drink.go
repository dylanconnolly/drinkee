package drinkee

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
)

type DrinkService interface {
	FindDrinkByID(ctx *gin.Context, id int) (*Drink, error)
	FindDrinks(ctx *gin.Context, f DrinkFilter) ([]*Drink, error)
	CreateDrink(ctx *gin.Context, cr *CreateDrink) error
	GenerateDrinks(c *gin.Context, i []Ingredient) ([]*Drink, error)
	FindIngredients(c *gin.Context) ([]*Ingredient, error)
}

type Drink struct {
	ID               int                  `json:"id"`
	Name             string               `json:"name"`
	DisplayName      string               `json:"displayName" db:"display_name"`
	Description      string               `json:"description,omitempty"`
	Instructions     string               `json:"instructions"`
	DrinkIngredients DrinkIngredientSlice `json:"drinkIngredients" db:"drink_ingredients"`
}

type DrinkResponse struct {
	Drink      `json:"drink"`
	TotalCount int `json:"total_count"`
}

type CreateDrink struct {
	Name             string            `json:"name" binding:"required"`
	DisplayName      string            `json:"displayName" binding:"required"`
	Description      string            `json:"description"`
	Instructions     string            `json:"instructions" binding:"required"`
	DrinkIngredients []DrinkIngredient `json:"drinkIngredients" binding:"required"`
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

type DrinkFilter struct {
	Limit int
	Skip  int
	Name  *string `json:"name,omitempty"`
	ID    *int    `json:"id,omitempty"`
}
