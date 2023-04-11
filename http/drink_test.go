package http_test

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	drinkeehttp "github.com/dylanconnolly/drinkee/http"
	"github.com/dylanconnolly/drinkee/postgres"
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
	pool.MaxWait = 120 * time.Second
	if err = pool.Retry(func() error {
		db, err = sqlx.Open("postgres", databaseUrl)
		if err != nil {
			return err
		}
		return db.Ping()
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	return db, pool, resource
}

func cleanupDatabase(p *dockertest.Pool, r *dockertest.Resource) {
	fmt.Print(p, r)
	// p.Purge(r)
}

func TestGetDrinks(t *testing.T) {
	t.Parallel()
	db, p, resource := startDatabase(t)
	s := drinkeehttp.NewServer()
	s.DrinkService = postgres.NewDrinkService(db)

	s.DrinkService.SimpleFind()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/drinks", nil)
	s.Router.ServeHTTP(w, req)

	fmt.Print(w.Body)
	assert.Equal(t, http.StatusOK, w.Code)
	cleanupDatabase(p, resource)
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
