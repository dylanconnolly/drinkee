package integration_tests

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dylanconnolly/drinkee/drinkee"
	drinkeehttp "github.com/dylanconnolly/drinkee/http"
	"github.com/dylanconnolly/drinkee/postgres"
	test_utils "github.com/dylanconnolly/drinkee/test/utils"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

var s = drinkeehttp.NewServer()

func TestGetDrinks(t *testing.T) {
	t.Parallel()
	db, p, resource := test_utils.SetupIntegrationTest(t, 5)
	defer test_utils.TeardownIntegrationTest(p, resource)

	s.DrinkService = postgres.NewDrinkService(db)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/drinks", nil)
	s.Router.ServeHTTP(w, req)

	var drinks []drinkee.DrinkResponse
	// read, _ := resp.Body.ReadBytes(0)
	// json.Unmarshal(read, &drinks)
	b, err := ioutil.ReadAll(w.Body)
	if err != nil {
		t.Errorf("Error reading drink response body: %s", err)
	}
	err = json.Unmarshal(b, &drinks)
	if err != nil {
		t.Errorf("Error unmarshalling drink response into drinks: %s", err)
	}

	assert.Equal(t, http.StatusOK, w.Code)
	assert.NotEmpty(t, drinks)
	assert.Equal(t, 5, len(drinks))
}

func TestGetIngredients(t *testing.T) {
	t.Parallel()
	db, p, resource := test_utils.SetupIntegrationTest(t, 5)
	defer test_utils.TeardownIntegrationTest(p, resource)

	s.DrinkService = postgres.NewDrinkService(db)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/ingredients", nil)
	s.Router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

// func TestCreateDrink(t *testing.T) {
// 	t.Parallel()

// 	reqBody := struct {
// 			name string
// 			displayName string
// 			description string
// 			instructions string
// 			drinkIngredients []struct {
// 				name string
// 				measurement string
// 			}
// 		}{
// 		"name": "moscow mule",
// 		"displayName": "Moscow Mule",
// 		"description": "Refreshing vodka based drink",
// 		"instructions": "combine ingredients and enjoy",
// 		"drinkIngredients": [
// 			{
// 				"name": "Vodka",
// 				"measurement": "1.5 fl oz"
// 			},
// 			{
// 				"name": "Ginger beer",
// 				"measurement": "3 fl oz"
// 			},
// 			{
// 				"name": "Lime",
// 				"measurement": "1 slice"
// 			}
// 		]
// 	}

// 	db, p, resource := startDatabase(t)
// 	s := drinkeehttp.NewServer()
// 	s.DrinkService = postgres.NewDrinkService(db)

// 	w := httptest.NewRecorder()
// 	req, err := http.NewRequest("POST", "/drinks")

// }
