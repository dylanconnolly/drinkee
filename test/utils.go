package test_utils

import (
	"fmt"

	"github.com/dylanconnolly/drinkee/drinkee"
	"github.com/jmoiron/sqlx"
)

func SeedDb(db *sqlx.DB) error {
	drinkStructs := generateMockDrinks(20)
	ingredientStructs := generateMockIngredients(20)

	tx := db.MustBegin()

	tx.NamedExec(`
		INSERT INTO drinks (name, display_name, description, instructions) VALUES (:name, :display_name, :description, :instructions)
	`, drinkStructs)
	tx.NamedExec(`
		INSERT INTO ingredients (name, display_name) VALUES (:name, :display_name)
	`, ingredientStructs)
	return tx.Commit()
}

func generateMockDrinks(n int) []drinkee.Drink {
	var drinks []drinkee.Drink

	for i := 1; i <= n; i++ {
		drink := drinkee.Drink{
			Name:         "test drink " + fmt.Sprint(i),
			DisplayName:  "Test Drink " + fmt.Sprint(i),
			Description:  "Test drink " + fmt.Sprint(i) + " description",
			Instructions: "Instructions for Test Drink " + fmt.Sprint(i),
		}

		drinks = append(drinks, drink)
	}

	return drinks
}

func generateMockIngredients(n int) []drinkee.Ingredient {
	var ingredients []drinkee.Ingredient

	for i := 1; i <= n; i++ {
		ingredient := drinkee.Ingredient{
			Name:        "test ingredient " + fmt.Sprint(i),
			DisplayName: "Test Ingredient" + fmt.Sprint(i),
		}

		ingredients = append(ingredients, ingredient)
	}

	return ingredients
}
