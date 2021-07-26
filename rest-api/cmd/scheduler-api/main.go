package main

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/sethvargo/go-envconfig"
	"log"
	"net/http"
	"workshops/rest-api/internal/config"
	handler "workshops/rest-api/internal/delivery/http"
	"workshops/rest-api/internal/delivery/http/router"
	"workshops/rest-api/internal/entities"
	"workshops/rest-api/internal/repositories/postgre_sql"
	"workshops/rest-api/internal/services"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func main() {

	ctx := context.Background()
	//Init App config
	var c config.AppConfig
	if err := envconfig.Process(ctx, &c); err != nil {
		log.Fatal(err)
	}
	//init db
	db := connectToDB(ctx, c)
	defer db.Close()

	jwtWrapper := entities.NewJwtWrapper(c)
	taskRepo := postgre_sql.NewTask(db)
	userRepo := postgre_sql.NewUser(db)

	taskService := services.NewTask(taskRepo, userRepo)
	userService := services.NewUser(userRepo, &jwtWrapper)

	taskController := handler.NewTask(taskService)
	userController := handler.NewUser(userService)
	authController := handler.NewAuth(userService)

	r := mux.NewRouter()
	r.Use(handler.RecoverMiddleware)
	auhMiddleware := mux.MiddlewareFunc(handler.AuthMiddlewareAdapter(userService, &jwtWrapper))

	router.Task(r.PathPrefix("/tasks").Subrouter(), taskController, auhMiddleware)
	router.User(r.PathPrefix("/users").Subrouter(), userController, auhMiddleware)
	router.Auth(r, authController, auhMiddleware)

	log.Fatal(http.ListenAndServe(":8000", r))
}

func connectToDB(ctx context.Context, c config.AppConfig) *sql.DB {
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
