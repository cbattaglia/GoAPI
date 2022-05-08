module recipe-api

go 1.18

replace recipe-api/data => ./data

replace recipe-api/handlers => ./handlers

require (
	github.com/gorilla/mux v1.8.0 // indirect
	github.com/joho/godotenv v1.4.0 // indirect
	github.com/lib/pq v1.10.5 // indirect
)
