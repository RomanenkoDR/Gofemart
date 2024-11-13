package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
)

type Config struct {
	Port string `default:"8080"`
}

type Application struct {
	Config Config
	//Model
}

func (app *Application) Serve() error {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	port := os.Getenv("PORT")
	fmt.Printf("Listening on port %s\n", port)

	srv := &http.Server{
		Addr: fmt.Sprintf(":%s" + port),
		// TODO: add router
	}
	return srv.ListenAndServe()
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	cfg := Config{
		Port: os.Getenv("PORT"),
	}

	// TODO: connection to db

	app := &Application{
		Config: cfg,
		// TODO: add models later
	}

	err = app.Serve()
	if err != nil {
		log.Fatal(err)
	}

}
