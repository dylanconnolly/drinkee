package router_test

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/dylanconnolly/drinkee/router"
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
	p.Purge(r)
}

// func TestMain(m *testing.M) {
// 	// uses a sensible default on windows (tcp/http) and linux/osx (socket)
// 	pool, err := dockertest.NewPool("")
// 	if err != nil {
// 		log.Fatalf("Could not construct pool: %s", err)
// 	}

// 	err = pool.Client.Ping()
// 	if err != nil {
// 		log.Fatalf("Could not connect to Docker: %s", err)
// 	}

// 	// pulls an image, creates a container based on it and runs it
// 	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
// 		Repository: "postgres",
// 		Tag:        "13",
// 		Env: []string{
// 			"POSTGRES_PASSWORD=password",
// 			"POSTGRES_USER=testuser",
// 			"POSTGRES_DB=drinkee",
// 			"listen_addresses = '*'",
// 		},
// 	}, func(config *docker.HostConfig) {
// 		// set AutoRemove to true so that stopped container goes away by itself
// 		config.AutoRemove = true
// 		config.RestartPolicy = docker.RestartPolicy{Name: "no"}
// 	})
// 	if err != nil {
// 		log.Fatalf("Could not start resource: %s", err)
// 	}

// 	hostAndPort := resource.GetHostPort("5432/tcp")
// 	databaseUrl := fmt.Sprintf("postgres://testuser:password@%s/drinkee?sslmode=disable", hostAndPort)

// 	log.Println("Connecting to database on url: ", databaseUrl)

// 	resource.Expire(120) // Tell docker to hard kill the container in 120 seconds

// 	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
// 	pool.MaxWait = 120 * time.Second
// 	if err = pool.Retry(func() error {
// 		db, err = sqlx.Open("postgres", databaseUrl)
// 		if err != nil {
// 			return err
// 		}
// 		return db.Ping()
// 	}); err != nil {
// 		log.Fatalf("Could not connect to docker: %s", err)
// 	}
// 	//Run tests
// 	code := m.Run()

// 	// You can't defer this because os.Exit doesn't care for defer
// 	if err := pool.Purge(resource); err != nil {
// 		log.Fatalf("Could not purge resource: %s", err)
// 	}

// 	os.Exit(code)
// }

func TestGetDrinks(t *testing.T) {
	t.Parallel()
	db, p, resource := startDatabase(t)
	r := router.CreateNewRouter(db)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/drinks", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	cleanupDatabase(p, resource)
}

func TestGetIngredients(t *testing.T) {
	t.Parallel()
	db, p, resource := startDatabase(t)
	r := router.CreateNewRouter(db)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ingredients", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	cleanupDatabase(p, resource)
}
