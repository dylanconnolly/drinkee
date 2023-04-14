package integration_tests

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/dylanconnolly/drinkee/drinkee"
	drinkeehttp "github.com/dylanconnolly/drinkee/http"
	drinkeepg "github.com/dylanconnolly/drinkee/postgres"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/ory/dockertest"
	"github.com/ory/dockertest/docker"
	"github.com/stretchr/testify/assert"
)

func startDatabase(t *testing.T) (*sqlx.DB, *dockertest.Pool, *dockertest.Resource) {
	var db *sqlx.DB
	t.Helper()

	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not construct pool: %s", err)
	}

	err = pool.Client.Ping()
	if err != nil {
		log.Fatalf("Could not connect to Docker: %s", err)
	}

	// pulls an image, creates a container based on it and runs it
	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "13",
		Env: []string{
			"POSTGRES_PASSWORD=password",
			"POSTGRES_USER=testuser",
			"POSTGRES_DB=drinkee",
			"listen_addresses = '*'",
		},
	}, func(config *docker.HostConfig) {
		// set AutoRemove to true so that stopped container goes away by itself
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{Name: "no"}
	})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	hostAndPort := resource.GetHostPort("5432/tcp")
	databaseUrl := fmt.Sprintf("postgres://testuser:password@%s/drinkee?sslmode=disable", hostAndPort)

	log.Println("Connecting to database on url: ", databaseUrl)

	resource.Expire(120) // Tell docker to hard kill the container in 120 seconds

	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	pool.MaxWait = 300 * time.Second
	if err = pool.Retry(func() error {
		db, err = sqlx.Open("postgres", databaseUrl)
		if err != nil {
			return err
		}
		return db.Ping()
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	// run migrations
	log.Println("running db migrations...")
	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
	if err != nil {
		fmt.Println("error: ", err)
	}

	m, err := migrate.NewWithDatabaseInstance("file:///Users/dconnolly/repos/drinkee/db/migrations", "postgres", driver)
	if err != nil {
		fmt.Println("error: ", err)
	}
	err = m.Up()
	if err != nil {
		fmt.Println("error: ", err)
	}

	log.Println("db migrations complete!")

	// seeding
	log.Println("starting seeding db...")
	// test_utils.SeedDb(db, 5)
	log.Println("done seeding")

	return db, pool, resource
}

// func runMigrations(migrationsPath string, db *sqlx.DB) error {
// 	if migrationsPath == "" {
// 		return errors.New("missing migrations path")
// 	}

// 	fmt.Println("im in the migrate thing")

// 	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
// 	if err != nil {
// 		return err
// 	}
// 	m, err := migrate.NewWithDatabaseInstance("file://"+migrationsPath, "postgres", driver)
// 	if err != nil {
// 		return err
// 	}
// 	err = m.Up()
// 	if err != nil && err != migrate.ErrNoChange {
// 		return err
// 	}
// 	return nil
// }

func cleanupDatabase(p *dockertest.Pool, r *dockertest.Resource) {
	// fmt.Print(p, r)
	p.Purge(r)
}

var s = drinkeehttp.NewServer()

func TestGetDrinks(t *testing.T) {
	t.Parallel()
	db, p, resource := startDatabase(t)
	defer cleanupDatabase(p, resource)
	// s := drinkeehttp.NewServer()
	s.DrinkService = drinkeepg.NewDrinkService(db)

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
}

func TestGetIngredients(t *testing.T) {
	t.Parallel()
	db, p, resource := startDatabase(t)
	defer cleanupDatabase(p, resource)
	// s := drinkeehttp.NewServer()
	s.DrinkService = drinkeepg.NewDrinkService(db)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/ingredients", nil)
	s.Router.ServeHTTP(w, req)

	// log.Printf("body response: %+v", w.Body)
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
