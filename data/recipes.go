package data

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var db *sql.DB
var err error

func init() {
	log.Println("Load ENV file")
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	port, err := strconv.Atoi(os.Getenv("POSTGRES_PORT"))
	if err != nil {
		log.Fatal("Error converting ENV POSTGRES_PORT to int")
	}

	log.Println("Connecting to Database")
	postgresInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("PI_ADDR"),
		port,
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DATABASE"))

	db, err = sql.Open("postgres", postgresInfo)
	if err != nil {
		log.Printf("Error connecting to rpi database: %s\n", err)
	}
}

type Recipe struct {
	ID        int         `json:"id"`
	Recipe    RecipeAttrs `json:"recipe"`
	CreatedAt time.Time   `json:"createdAt"`
	UpdatedAt time.Time   `json:"updatedAt"`
}

type RecipeAttrs struct {
	Name        string   `json:"name,omitempty"`
	Ingredients []string `json:"ingredients,omitempty"`
}

type Recipes []*Recipe

func (a RecipeAttrs) Value() (driver.Value, error) {
	return json.Marshal(a)
}

func (a *RecipeAttrs) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(b, &a)
}

func (r *Recipes) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(r)
}

func (r *Recipe) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(r)
}

func (r *RecipeAttrs) FromJSON(rdr io.Reader) error {
	e := json.NewDecoder(rdr)
	return e.Decode(r)
}

func GetRecipes() (*sql.Rows, error) {
	return db.Query("SELECT * FROM recipes")
}

func GetRecipe(id int) *sql.Row {
	return db.QueryRow("SELECT * FROM recipes where id = $1", id)
}

func CreateRecipe(rAttrs *RecipeAttrs) (int64, error) {
	stmt, err := db.Prepare("INSERT INTO recipes(recipe) VALUES($1)")
	if err != nil {
		return -1, err
	}

	result, err := stmt.Exec(rAttrs)
	if err != nil {
		return -1, err
	}

	rowCnt, err := result.RowsAffected()
	if err != nil {
		return -1, err
	}

	return rowCnt, nil
}

func UpdateRecipe(id int, rAttrs *RecipeAttrs) (int64, error) {
	stmt, err := db.Prepare("UPDATE recipes SET recipe = $1 WHERE id = $2")
	if err != nil {
		return -1, err
	}

	result, err := stmt.Exec(rAttrs, id)
	if err != nil {
		return -1, err
	}

	rowCnt, err := result.RowsAffected()
	if err != nil {
		return -1, err
	}

	return rowCnt, nil
}

func DeleteRecipe(id int) (int64, error) {
	stmt, err := db.Prepare("DELETE FROM recipes WHERE id = $1")
	if err != nil {
		return -1, err
	}

	result, err := stmt.Exec(id)
	if err != nil {
		return -1, err
	}

	rowCnt, err := result.RowsAffected()
	if err != nil {
		return -1, err
	}

	return rowCnt, nil
}
