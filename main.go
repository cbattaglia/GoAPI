package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"recipe-api/handlers"
	"time"

	"github.com/gorilla/mux"
)

func main() {

	l := log.New(os.Stdout, "recipe-api", log.LstdFlags)

	recipeHandler := handlers.NewRecipes(l)

	sm := mux.NewRouter()

	getRouter := sm.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/", recipeHandler.GetRecipes)           //GET ALL
	getRouter.HandleFunc("/{id:[0-9]+}", recipeHandler.GetRecipe) //GET ONE

	putRouter := sm.Methods(http.MethodPut).Subrouter()
	putRouter.HandleFunc("/{id:[0-9]+}", recipeHandler.UpdateRecipes)
	putRouter.Use(recipeHandler.MiddlewareRecipeValidation) //UPDATE

	postRouter := sm.Methods(http.MethodPost).Subrouter()
	postRouter.HandleFunc("/", recipeHandler.CreateRecipe)
	postRouter.Use(recipeHandler.MiddlewareRecipeValidation) //CREATE

	deleteRouter := sm.Methods(http.MethodDelete).Subrouter()
	deleteRouter.HandleFunc("/{id:[0-9]+}", recipeHandler.DeleteRecipe)

	s := http.Server{
		Addr:         ":9090",
		Handler:      sm,
		ErrorLog:     l,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	go func() {
		l.Println("Starting server on port 9090")

		err := s.ListenAndServe()
		if err != nil {
			l.Printf("Error starting server: %s\n", err)
			os.Exit(1)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, os.Kill)

	sig := <-c
	log.Println("Got signal:", sig)

	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	s.Shutdown(ctx)
}
