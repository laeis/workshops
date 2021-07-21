package main

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/sethvargo/go-envconfig"
	"log"
	"net/http"
	router "workshops/rest-api/internal/delivery/http"
	"workshops/rest-api/internal/repositories/postgre_sql"
	"workshops/rest-api/internal/services"
)

type AppConfig struct {
	*DBConfig
}

type DBConfig struct {
	Host     string `env:"POSTGRES_HOST,default=localhost"`
	Port     string `env:"POSTGRES_PORT,default=5432"`
	User     string `env:"POSTGRES_USER"`
	Password string `env:"POSTGRES_PASSWORD"`
	Database string `env:"POSTGRES_DB"`
}

func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func main() {

	ctx := context.Background()
	//Init App config
	var c AppConfig
	if err := envconfig.Process(ctx, &c); err != nil {
		log.Fatal(err)
	}
	//init db
	db := connectToDB(ctx, c)
	defer db.Close()

	taskRepo := postgre_sql.NewTask(db)

	taskService := services.NewTask(taskRepo)
	taskController := router.NewTask(taskService)
	log.Fatal(http.ListenAndServe(":8000", router.NewRouter(taskController)))
}

func connectToDB(ctx context.Context, c AppConfig) *sql.DB {
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
	return db
}
