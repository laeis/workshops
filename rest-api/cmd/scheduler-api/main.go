package main

import (
	"database/sql"
	"fmt"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
	router "workshops/rest-api/internal/delivery/http"
	"workshops/rest-api/internal/repositories"
	"workshops/rest-api/internal/services"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}

}

func main() {
	//init db
	db := connectToDB()
	defer db.Close()

	taskRepo := repositories.NewTask(db)
	//taskRepo := repositories.NewInMemoryTask()
	taskService := services.NewTaskService(taskRepo)
	taskController := router.NewTaskHandler(taskService)
	log.Fatal(http.ListenAndServe(":8000", router.NewRouter(taskController)))
}

func connectToDB() *sql.DB {
	host := os.Getenv("POSTGRES_HOST")
	port := os.Getenv("POSTGRES_PORT")
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbname := os.Getenv("POSTGRES_DB")

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

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
