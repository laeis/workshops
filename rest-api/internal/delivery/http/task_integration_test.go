//+build integration
package http

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/sethvargo/go-envconfig"
	"github.com/stretchr/testify/assert"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"workshops/rest-api/internal/entities"
	"workshops/rest-api/internal/repositories/postgre_sql"
	"workshops/rest-api/internal/services"
	"workshops/rest-api/internal/validators"
)

type TestConfig struct {
	*DBConfig
}

type DBConfig struct {
	Host     string `env:"TEST_POSTGRES_HOST,default=localhost"`
	Port     string `env:"TEST_POSTGRES_PORT,default=5432"`
	User     string `env:"TEST_POSTGRES_USER"`
	Password string `env:"TEST_POSTGRES_PASSWORD"`
	Database string `env:"TEST_POSTGRES_DB"`
}

var c TestConfig

var taskController TaskHandler

func TestMain(m *testing.M) {
	ctx := context.Background()
	if err := godotenv.Load("../../../.env-test"); err != nil {
		log.Fatal("No .env file found", err)
	}

	if err := envconfig.Process(ctx, &c); err != nil {
		log.Fatal(err)
	}

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		c.DBConfig.Host, c.DBConfig.Port, c.DBConfig.User, c.DBConfig.Password, c.DBConfig.Database)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}

	truncateTables(db, []string{"tasks"})
	taskRepo := postgre_sql.NewTask(db)
	userRepo := postgre_sql.NewUser(db)
	taskService := services.NewTask(taskRepo, userRepo)
	taskController = NewTask(taskService)

	os.Exit(m.Run())
}

func TestTaskHandler_Create(t *testing.T) {
	cases := []struct {
		name         string
		statusCode   int
		bodyResponse string
	}{
		{
			name:         "Test",
			statusCode:   200,
			bodyResponse: `{"payload":{"Id":1,"Title":"Test event","Description":"","Category":"event","Date":"0001-01-01T00:00:00Z"}}`,
		},
	}
	for i, c := range cases {
		t.Run(fmt.Sprintf("Fetch Case %d: %s ", i, c.name), func(t *testing.T) {
			r := mux.NewRouter()
			r.HandleFunc("/tasks", taskController.Create).Methods(http.MethodPost)

			data := new(bytes.Buffer)
			json.NewEncoder(data).Encode(entities.Task{
				Category: validators.EVENT,
				Title:    "Test event",
			})
			req := httptest.NewRequest("POST", "/tasks", data)
			res := httptest.NewRecorder()

			r.ServeHTTP(res, req)

			body := res.Body.String()

			assert.Equal(t, c.statusCode, res.Code)
			assert.Equal(t, c.bodyResponse, body)
		})
	}

}

func TestTaskHandler_Fetch(t *testing.T) {
	cases := []struct {
		name         string
		statusCode   int
		bodyResponse string
	}{
		{
			name:         "Test",
			statusCode:   200,
			bodyResponse: `{"payload":{"Id":1,"Title":"Test event","Description":"","Category":"event","Date":"0001-01-01T00:00:00Z"}}`,
		},
	}
	for i, c := range cases {
		t.Run(fmt.Sprintf("Fetch Case %d: %s ", i, c.name), func(t *testing.T) {
			r := mux.NewRouter()
			r.HandleFunc("/tasks/{id}", taskController.Get).Methods(http.MethodGet)

			req := httptest.NewRequest("GET", "/tasks/1", nil)
			res := httptest.NewRecorder()

			r.ServeHTTP(res, req)

			body := res.Body.String()

			assert.Equal(t, c.statusCode, res.Code)
			assert.Equal(t, c.bodyResponse, body)
		})
	}
}

func truncateTables(db *sql.DB, tables []string) {
	_, _ = db.Exec("SET FOREIGN_KEY_CHECKS=0;")

	for _, v := range tables {
		_, _ = db.Exec(fmt.Sprintf("TRUNCATE TABLE %s RESTART IDENTITY;", v))
	}

	_, _ = db.Exec("SET FOREIGN_KEY_CHECKS=1;")
}
