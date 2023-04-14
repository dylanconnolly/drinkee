package test_utils

import (
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/dylanconnolly/drinkee/drinkee"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/ory/dockertest"
	"github.com/ory/dockertest/docker"
)

func SetupIntegrationTest(t *testing.T, mockDataCount int) (*sqlx.DB, *dockertest.Pool, *dockertest.Resource) {
	db, pool, resource := setupDatabase(t)

	if err := seedDatabase(db, mockDataCount); err != nil {
		log.Fatalf("Could not seed database: %s", err)
	}

	return db, pool, resource
}

func TeardownIntegrationTest(p *dockertest.Pool, r *dockertest.Resource) {
	p.Purge(r)
}

func setupDatabase(t *testing.T) (*sqlx.DB, *dockertest.Pool, *dockertest.Resource) {
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

	return db, pool, resource
}

func seedDatabase(db *sqlx.DB, n int) error {
	drinkStructs := generateMockDrinks(n)
	ingredientStructs := generateMockIngredients(n)
	drinkIngredientRows := generateMockDrinkIngredientRows(n)

	tx := db.MustBegin()

	tx.NamedExec(`
		INSERT INTO drinks (name, display_name, description, instructions) VALUES (:name, :display_name, :description, :instructions)
	`, drinkStructs)
	tx.NamedExec(`
		INSERT INTO ingredients (name, display_name) VALUES (:name, :display_name)
	`, ingredientStructs)
	tx.NamedExec(`
		INSERT INTO drink_ingredients (drink_id, ingredient_id, measurement) VALUES (:drink_id, :ingredient_id, :measurement)
	`, drinkIngredientRows)
	return tx.Commit()
}

type DrinkIngredientRow struct {
	DrinkID      int `db:"drink_id"`
	IngredientID int `db:"ingredient_id"`
	Measurement  string
}

func generateMockDrinks(n int) []drinkee.Drink {
	var drinks []drinkee.Drink

	for i := 1; i <= n; i++ {
		drink := drinkee.Drink{
			Name:         fmt.Sprintf("test drink %d", i),
			DisplayName:  fmt.Sprintf("Test Drink %d", i),
			Description:  fmt.Sprintf("Test drink description %d", i),
			Instructions: fmt.Sprintf("Instructions for Test Drink %d", i),
		}

		drinks = append(drinks, drink)
	}

	return drinks
}

func generateMockIngredients(n int) []drinkee.Ingredient {
	var ingredients []drinkee.Ingredient

	for i := 1; i <= n; i++ {
		ingredient := drinkee.Ingredient{
			Name:        fmt.Sprintf("test ingredient %d", i),
			DisplayName: fmt.Sprintf("Test Ingredient %d", i),
		}

		ingredients = append(ingredients, ingredient)
	}

	return ingredients
}

func generateMockDrinkIngredientRows(n int) []DrinkIngredientRow {
	var drinkIngredients []DrinkIngredientRow

	for i := 1; i <= n; i++ {
		di := DrinkIngredientRow{
			DrinkID:      i,
			IngredientID: i,
			Measurement:  fmt.Sprintf("Drink %d ingredient %d measurement", i, i),
		}

		drinkIngredients = append(drinkIngredients, di)
	}

	return drinkIngredients
}
