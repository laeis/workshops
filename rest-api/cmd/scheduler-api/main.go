package main

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sethvargo/go-envconfig"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
	"workshops/rest-api/internal/config"
	handler "workshops/rest-api/internal/delivery/http"
	"workshops/rest-api/internal/delivery/http/middlewares"
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

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		oscall := <-c
		log.Printf("system call:%+v", oscall)
		cancel()
	}()

	if err := serve(ctx); err != nil {
		log.Printf("failed to serve:+%v\n", err)
	}

}

func serve(ctx context.Context) (err error) {

	//Init App config
	var c config.AppConfig
	if err := envconfig.Process(ctx, &c); err != nil {
		return err
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
	r.Use(middlewares.Recover)
	r.Use(middlewares.Metrics)

	auhMiddleware := mux.MiddlewareFunc(middlewares.AuthAdapter(userService, &jwtWrapper))

	router.Task(r.PathPrefix("/tasks").Subrouter(), taskController, auhMiddleware)
	router.User(r.PathPrefix("/users").Subrouter(), userController, auhMiddleware)
	router.Auth(r, authController, auhMiddleware)

	srv := &http.Server{
		Addr:    ":8000",
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen:%+s\n", err)
		}
	}()

	log.Printf("server started")

	metrics := r.PathPrefix("/metrics").Subrouter()
	metrics.Handle("", promhttp.Handler())
	promSrv := &http.Server{
		Addr:    ":9000",
		Handler: metrics,
	}

	go func() {
		if err := promSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("promteus listen:%+s\n", err)
		}
	}()

	log.Printf("prometeus started")

	<-ctx.Done()

	log.Printf("server stopped")

	ctxShutDown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		cancel()
	}()

	if err = srv.Shutdown(ctxShutDown); err != nil {
		log.Fatalf("server Shutdown Failed:%+s", err)
	}

	if err == http.ErrServerClosed {
		err = nil
	}

	if err := promSrv.Shutdown(ctxShutDown); err != nil {
		log.Fatalf("promtetus server Shutdown Failed:%+s", err)
	}

	log.Printf("server exited properly")

	return
}

func connectToDB(ctx context.Context, c config.AppConfig) *sql.DB {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		c.DBConfig.Host, c.DBConfig.Port, c.DBConfig.User, c.DBConfig.Password, c.DBConfig.Database)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	err = db.PingContext(ctx)
	if err != nil {
		panic(err)
	}
	return db
}
