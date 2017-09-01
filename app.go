package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

const (
	MYSQL_USERNAME = "mlabouardy"
	MYSQL_PASSWORD = "12345678"
	MYSQL_HOST     = ""
	MYSQL_PORT     = 3306
	MYSQL_DB       = "mydb"
	DATASOURCE     = "%s:%s@tcp(%s:%d)/%s"
)

type Movie struct {
	ID   string
	Name string
}

var db *sql.DB
var err error

func GetMoviesEndpoint(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT * FROM movies")
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
	}
	movies := make([]Movie, 0)
	for rows.Next() {
		var movie Movie
		err := rows.Scan(&movie.ID, &movie.Name)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		movies = append(movies, movie)
	}
	respondWithJson(w, http.StatusOK, movies)
}

func PostMoviesEndpoint(w http.ResponseWriter, r *http.Request) {
	var movie Movie
	if err := json.NewDecoder(r.Body).Decode(&movie); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
	}
	_, err := db.Exec("INSERT INTO movies(name) VALUES(?)", movie.Name)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
	}
	respondWithJson(w, http.StatusCreated, map[string]string{"result": "success"})
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s : %s", msg, err)
	}
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	respondWithJson(w, code, map[string]string{"error": msg})
}

func respondWithJson(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func main() {
	ds := fmt.Sprintf(DATASOURCE, MYSQL_USERNAME, MYSQL_PASSWORD, MYSQL_HOST, MYSQL_PORT, MYSQL_DB)
	fmt.Println(ds)
	db, err = sql.Open("mysql", ds)
	if err != nil {
		failOnError(err, "Failed to connect to database")
	}
	//defer db.Close()

	r := mux.NewRouter()
	r.HandleFunc("/api/movies", GetMoviesEndpoint).Methods("GET")
	r.HandleFunc("/api/movies", PostMoviesEndpoint).Methods("POST")
	if err := http.ListenAndServe(":3000", r); err != nil {
		failOnError(err, "Cannot start an HTTP server")
	}
}
