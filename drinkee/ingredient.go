package drinkee

type Ingredient struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	DisplayName string `json:"displayName" db:"display_name"`
}
