package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/lib/pq"
	_ "github.com/lib/pq"
)

const (
	host = "postgres.postgres"
	port = 5432
	//host     = "localhost"
	//port     = 51591
	user     = "postgres"
	//Use kube secret environment variable
	password = <PASSWORD>
	dbname   = "postgres"
	table    = "movies"
)

// Regex to validate API URL
var (
	getDataRe    = regexp.MustCompile(`^\/movies\/*$`)
	createDataRe = regexp.MustCompile(`^\/movies\/*$`)
)
var listGetMatches = []*regexp.Regexp{getDataRe}
var listPostMatches = []*regexp.Regexp{createDataRe}

// Pass DB object via handlers
type dataHandler struct {
	DB *sql.DB
}

// Agreed data format: https://github.com/prust/wikipedia-movie-data
type Movie struct {
	Title            string   `json:"title,omitempty"`
	Year             int      `json:"year,omitempty"`
	Cast             []string `json:"cast,omitempty"`
	Genres           []string `json:"genres,omitempty"`
	Href             string   `json:"href,omitempty"`
	Extract          string   `json:"extract,omitempty"`
	Thumbnail        string   `json:"thumbnail,omitempty"`
	Thumbnail_width  int      `json:"thumbnail_width,omitempty"`
	Thumbnail_height int      `json:"thumbnail_height,omitempty"`
}

// Connect to DB
func Connect() *sql.DB {
	connInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", connInfo)
	if err != nil {
		log.Fatal(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	log.Println("Successfully connected to db at " + host)
	return db
}

// Close DB connection
func CloseConnection(db *sql.DB) {
	defer db.Close()
}

// Create movies table if not exist
func CreateTable(db *sql.DB) {
	var exists bool
	if err := db.QueryRow("SELECT EXISTS (SELECT FROM pg_tables WHERE  schemaname = 'public' AND tablename = '" + table + "' );").Scan(&exists); err != nil {
		log.Println("failed to execute query", err)
		return
	}
	if !exists {
		results, err := db.Query(`CREATE TABLE " + table + " (title varchar,year integer,"cast" varchar[],genres varchar[],href varchar,extract TEXT,thumbnail varchar,thumbnail_width integer,thumbnail_height integer);`)
		if err != nil {
			log.Println("failed to execute query", err)
			return
		}
		log.Println("Table created successfully", results)
	} else {
		log.Println("Table '" + table + "' already exists ")
	}

}

// Create movies DB record (not required for implementation)
func (h *dataHandler) Create(w http.ResponseWriter, r *http.Request) {
	for _, match := range listPostMatches {
		matches := match.FindStringSubmatch(r.URL.Path)
		log.Println(matches, string(r.Method))
		if len(matches) < 1 {
			continue
		}
		switch {
		case match == createDataRe:
			defer r.Body.Close()
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				log.Fatalln(err)
				w.WriteHeader(500)
				return
			}
			var movie Movie
			json.Unmarshal(body, &movie)
			queryStmt := `INSERT INTO ` + table + ` (title,year,"cast",genres,href,extract,thumbnail,thumbnail_width,thumbnail_height) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING title;`
			err = h.DB.QueryRow(queryStmt, &movie.Title, &movie.Year, pq.Array(&movie.Cast), pq.Array(&movie.Genres), &movie.Href, &movie.Extract, &movie.Thumbnail, &movie.Thumbnail_width, &movie.Thumbnail_height).Scan(&movie.Title)
			if err != nil {
				log.Println("failed to execute query", err)
				w.WriteHeader(500)
				return
			}

			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode("Created")
		}
	}
}

// Get movies by title, year, cast member, or genre
func (h *dataHandler) Get(w http.ResponseWriter, r *http.Request) {
	for _, match := range listGetMatches {
		matches := match.FindStringSubmatch(r.URL.Path)
		log.Println(matches, string(r.Method))
		if len(matches) < 1 {
			continue
		}

		switch {
		case match == getDataRe:
			var results *sql.Rows
			var err error
			queryStmt := `SELECT * FROM ` + table
			if len(r.URL.Query()) == 0 {
				results, err = h.DB.Query(queryStmt)
			} else {
				count := 0
				var val string
				for k, v := range r.URL.Query() {
					count += 1
					val = strings.Join(v, " ")
					switch {
					case k == "title":
						queryStmt += ` WHERE title = $1;`
					case k == "year":
						queryStmt += ` WHERE year = $1;`
					case k == "cast":
						queryStmt += ` WHERE $1=ANY("cast");`
					case k == "genre":
						queryStmt += ` WHERE $1=ANY(genres);`
					}
					if count > 0 {
						break
					}
				}
				results, err = h.DB.Query(queryStmt, val)
			}
			if err != nil {
				log.Println("failed to execute query", err)
				w.WriteHeader(500)
				return
			}

			var movies = make([]Movie, 0)
			for results.Next() {
				var movie Movie
				err = results.Scan(&movie.Title, &movie.Year, pq.Array(&movie.Cast), pq.Array(&movie.Genres), &movie.Href, &movie.Extract, &movie.Thumbnail, &movie.Thumbnail_width, &movie.Thumbnail_height)
				if err != nil {
					log.Println("failed to scan", err)
					w.WriteHeader(500)
					return
				}

				movies = append(movies, movie)
			}

			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(movies)
		}
	}
}

// Route user requests based on HTTP method
func (h *dataHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	switch {
	case r.Method == http.MethodPost && (createDataRe.MatchString(r.URL.Path)):
		h.Create(w, r)
		return
	case r.Method == http.MethodGet && (getDataRe.MatchString(r.URL.Path)):
		h.Get(w, r)
		return
	default:
		notFound(w, r)
		return
	}
}

// Not found HTTP response
func notFound(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode("not found")
}

// HTTP request handler
func handleRequests(DB *sql.DB) {
	mux := http.NewServeMux()
	mux.Handle("/", &dataHandler{DB})
	log.Fatal(http.ListenAndServe(":8081", mux))
}

func main() {
	DB := Connect()
	CreateTable(DB)
	handleRequests(DB)
	CloseConnection(DB)
}
