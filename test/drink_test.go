package integration_tests

import (
	"bytes"
	"encoding/json"
	"fmt"
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

	var drinks []drinkee.Drink

	if w.Code != 200 {
		t.Errorf("Error with drink request: %s", w.Body)
		t.FailNow()
	}
	b, err := ioutil.ReadAll(w.Body)
	if err != nil {
		t.Errorf("Error reading drink response body: %s", err)
		t.FailNow()
	}
	err = json.Unmarshal(b, &drinks)
	if err != nil {
		t.Errorf("Error unmarshalling drink response into drinks: %s", err)
		t.FailNow()
	}

	assert.Equal(t, http.StatusOK, w.Code)
	assert.NotEmpty(t, drinks)
	assert.Equal(t, 5, len(drinks))
	assert.NotEmpty(t, drinks[0].ID)
	assert.NotEmpty(t, drinks[0].Name)
	assert.NotEmpty(t, drinks[0].DisplayName)
	assert.NotEmpty(t, drinks[0].Description)
	assert.NotEmpty(t, drinks[0].Instructions)
	assert.NotEmpty(t, drinks[0].DrinkIngredients)
}

func TestDrinkFilter(t *testing.T) {
	t.Parallel()
	db, p, resource := test_utils.SetupIntegrationTest(t, 10)
	defer test_utils.TeardownIntegrationTest(p, resource)

	s.DrinkService = postgres.NewDrinkService(db)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/drinks?limit=5", nil)
	s.Router.ServeHTTP(w, req)

	var drinks []drinkee.Drink

	if w.Code != 200 {
		t.Errorf("Error with drink request: %s", w.Body)
		t.FailNow()
	}
	b, err := ioutil.ReadAll(w.Body)
	if err != nil {
		t.Errorf("Error reading drink response body: %s", err)
		t.FailNow()
	}
	err = json.Unmarshal(b, &drinks)
	if err != nil {
		t.Errorf("Error unmarshalling drink response into drinks: %s", err)
		t.FailNow()
	}

	assert.Equal(t, http.StatusOK, w.Code)
	assert.NotEmpty(t, drinks)
	assert.Equal(t, 5, len(drinks))

	name := "test drink 1,test drink 2,test drink 3,test drink 4,test drink 5,test drink 6"
	filter := drinkee.DrinkFilter{
		Name: &name,
	}

	var buffer bytes.Buffer
	json.NewEncoder(&buffer).Encode(filter)

	req, _ = http.NewRequest("GET", "/api/v1/drinks", &buffer)
	s.Router.ServeHTTP(w, req)

	fmt.Printf("body response: %+v", w.Body)
}

func TestGetDrinkByID(t *testing.T) {
	t.Parallel()
	db, p, resource := test_utils.SetupIntegrationTest(t, 2)
	defer test_utils.TeardownIntegrationTest(p, resource)

	s.DrinkService = postgres.NewDrinkService(db)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/drinks/1", nil)
	s.Router.ServeHTTP(w, req)

	var drink drinkee.Drink

	if w.Code != 200 {
		t.Errorf("Error with drink request: %s", w.Body)
		t.FailNow()
	}

	b, err := ioutil.ReadAll(w.Body)
	if err != nil {
		t.Errorf("Error reading drink response body: %s", err)
		t.FailNow()
	}
	err = json.Unmarshal(b, &drink)
	if err != nil {
		t.Errorf("Error unmarshalling drink response into drinks: %s", err)
		t.FailNow()
	}

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "Test Drink 1", drink.DisplayName)
	assert.NotEmpty(t, drink.DrinkIngredients)
}

func TestGenerateDrinks(t *testing.T) {
	t.Parallel()
	db, p, resource := test_utils.SetupIntegrationTest(t, 5)
	defer test_utils.TeardownIntegrationTest(p, resource)

	s.DrinkService = postgres.NewDrinkService(db)

	ingredients := drinkeehttp.IngredientListRequest{
		Ingredients: []drinkee.Ingredient{
			{
				ID:          1,
				Name:        "test ingredient 1",
				DisplayName: "Test Ingredient 1",
			},
			{
				ID:          3,
				Name:        "test ingredient 3",
				DisplayName: "Test Ingredient 3",
			},
		},
	}

	var buffer bytes.Buffer
	err := json.NewEncoder(&buffer).Encode(ingredients)
	if err != nil {
		t.Errorf("Error encoding request body: %s", err)
		t.FailNow()
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/generateDrinks?strict=true", &buffer)

	s.Router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusAccepted, w.Code)
	fmt.Printf("response: %+v", w.Body)
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
