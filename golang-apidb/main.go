package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"

	"github.com/lib/pq"
	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	host = "postgres.postgres"
	port = 5432
	//host   = "localhost"
	//port   = 63711
	user   = "postgres"
	dbname = "postgres"
	table  = "movies"
)

var password = os.Getenv("DB_PASSWORD")

// Regex to validate API URL
var (
	getDataRe    = regexp.MustCompile(`^\/movies\/*$`)
	createDataRe = regexp.MustCompile(`^\/movies\/*$`)
)
var listGetMatches = []*regexp.Regexp{getDataRe}
var listPostMatches = []*regexp.Regexp{createDataRe}

// Pass DB object via handlers
type dataHandler struct {
	DB *gorm.DB
}

// Agreed data format: https://github.com/prust/wikipedia-movie-data
type Movie struct {
	Title            string         `json:"title,omitempty"`
	Year             int            `json:"year,omitempty"`
	Cast             pq.StringArray `json:"cast,omitempty" gorm:"type:string[]"`
	Genres           pq.StringArray `json:"genres,omitempty" gorm:"type:string[]"`
	Href             string         `json:"href,omitempty"`
	Extract          string         `json:"extract,omitempty"`
	Thumbnail        string         `json:"thumbnail,omitempty"`
	Thumbnail_width  int            `json:"thumbnail_width,omitempty"`
	Thumbnail_height int            `json:"thumbnail_height,omitempty"`
}

// Connect to DB
func Connect() *gorm.DB {
	dsn := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal(err)
	}
	log.Println("Successfully connected to db at " + host)
	return db
}

// Create movies table if not exist
func CreateTable(db *sql.DB) {
	var exists bool
	queryStmt := `SELECT EXISTS (SELECT FROM pg_tables WHERE  schemaname = 'public' AND tablename = $1);`
	if err := db.QueryRow(queryStmt, table).Scan(&exists); err != nil {
		log.Println("failed to execute query", err)
		return
	}
	if !exists {
		queryStmt := `CREATE TABLE ` + table + ` (title varchar,year integer,"cast" varchar[],genres varchar[],href varchar,extract TEXT,thumbnail varchar,thumbnail_width integer,thumbnail_height integer);`
		results, err := db.Query(queryStmt)
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
			if result := h.DB.Create(&movie); result.Error != nil {
				log.Println(result.Error)
			}
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
			//var err error
			var movies []Movie
			if len(r.URL.Query()) == 0 {
				if result := h.DB.Find(&movies); result.Error != nil {
					fmt.Println(result.Error)
				}
			} else {
				for _, key := range []string{"title", "year", "cast", "genre"} {
					val, exists := r.URL.Query()[key]
					if exists {
						switch true {
						case (key == "title" || key == "year"):
							if result := h.DB.Where(key+" = ?", val).Find(&movies); result.Error != nil {
								log.Println("failed to execute query", result.Error)
								notFound(w, r)
								return
							}
						case (key == "cast" || key == "genre"):
							if key == "genre" {
								key = "genres"
							}
							if result := h.DB.Where("?=ANY(\""+key+"\")", val).Find(&movies); result.Error != nil {
								log.Println("failed to execute query", result.Error)
								notFound(w, r)
								return
							}
						}
					}
				}
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
func handleRequests(DB *gorm.DB) {
	mux := http.NewServeMux()
	mux.Handle("/", &dataHandler{DB})
	log.Fatal(http.ListenAndServe(":8081", mux))
}

func main() {
	_, exists := os.LookupEnv("DB_PASSWORD")
	if !exists {
		panic("DB_PASSWORD not set")
	}
	DB := Connect()
	//CreateTable(DB)
	handleRequests(DB)
}
