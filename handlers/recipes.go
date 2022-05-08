package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"recipe-api/data"
	"strconv"

	"github.com/gorilla/mux"
)

type Recipes struct {
	l *log.Logger
}

func NewRecipes(l *log.Logger) *Recipes {
	return &Recipes{l}
}

func (r *Recipes) GetRecipes(rw http.ResponseWriter, hr *http.Request) {
	r.l.Println("Handle GET Recipes")

	rows, err := data.GetRecipes()
	if err != nil {
		http.Error(rw, fmt.Errorf("something went wrong reading from DB: %s", err).Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	recipes := make([]*data.Recipe, 0)
	for rows.Next() {
		recipe := new(data.Recipe)

		err := rows.Scan(&recipe.ID, &recipe.Recipe, &recipe.CreatedAt, &recipe.UpdatedAt)
		if err != nil {
			http.Error(rw, fmt.Errorf("something went wrong reading data rows: %s", err).Error(), http.StatusInternalServerError)
			return
		}
		recipes = append(recipes, recipe)
	}
	for _, recipe := range recipes {
		recipe.ToJSON(r.l.Writer())
	}
}

func (r *Recipes) GetRecipe(rw http.ResponseWriter, hr *http.Request) {
	r.l.Println("Handle GET Recipe")
	vars := mux.Vars(hr)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(rw, "Unable to convert id", http.StatusBadRequest)
		return
	}

	recipe := new(data.Recipe)
	rows := data.GetRecipe(id).Scan(&recipe.ID, &recipe.Recipe, &recipe.CreatedAt, &recipe.UpdatedAt)
	if rows != nil {
		http.Error(rw, fmt.Errorf("something went wrong reading from DB: %s", rows).Error(), http.StatusInternalServerError)
		return
	}

	recipe.ToJSON(r.l.Writer())
}

func (r *Recipes) UpdateRecipes(rw http.ResponseWriter, hr *http.Request) {
	vars := mux.Vars(hr)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(rw, "Unable to convert id", http.StatusBadRequest)
		return
	}

	r.l.Println("Handle PUT Recipes", id)
	recipeAttrs := hr.Context().Value(KeyRecipe{}).(data.RecipeAttrs)

	_, err = data.UpdateRecipe(id, &recipeAttrs)
	if err != nil {
		http.Error(rw, fmt.Errorf("something went wrong updating the record: %s", err).Error(), http.StatusBadRequest)
		return
	}
}

func (r *Recipes) CreateRecipe(rw http.ResponseWriter, hr *http.Request) {
	r.l.Println("Handle POST Recipes")

	recipeAttrs := hr.Context().Value(KeyRecipe{}).(data.RecipeAttrs)

	_, err := data.CreateRecipe(&recipeAttrs)
	if err != nil {
		http.Error(rw, fmt.Errorf("something went wrong creating the record: %s", err).Error(), http.StatusBadRequest)
		return
	}
}

func (r *Recipes) DeleteRecipe(rw http.ResponseWriter, hr *http.Request) {
	r.l.Println("Handle DELETE Recipe")
	vars := mux.Vars(hr)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(rw, "Unable to convert id", http.StatusBadRequest)
		return
	}

	_, delErr := data.DeleteRecipe(id)
	if delErr != nil {
		http.Error(rw, fmt.Errorf("something went wrong deleting from DB: %s", delErr).Error(), http.StatusInternalServerError)
		return
	}
}

type KeyRecipe struct{}

func (r Recipes) MiddlewareRecipeValidation(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, hr *http.Request) {
		recipe := data.RecipeAttrs{}

		err := recipe.FromJSON(hr.Body)
		if err != nil {
			http.Error(rw, "Unable to unmarshal json", http.StatusBadRequest)
			return
		}

		ctx := context.WithValue(hr.Context(), KeyRecipe{}, recipe)
		req := hr.WithContext(ctx)

		next.ServeHTTP(rw, req)
	})
}
