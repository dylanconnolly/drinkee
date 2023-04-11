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

type DrinkService interface {
	// FindDrinkByID(ctx *gin.Context, id int) (*Drink, error)
	FindDrinks(ctx *gin.Context) ([]*DrinkResponse, error)
	CreateDrink(ctx *gin.Context, cr *CreateDrink) error
	SimpleFind()
	GenerateDrinks(c *gin.Context, i []Ingredient) ([]DrinkResponse, error)
	FindIngredients(c *gin.Context) ([]*Ingredient, error)
}
