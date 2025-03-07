package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"
)

// Define MSDSCourse struct
type MSDSCourse struct {
	CID     string `json:"course_id"`
	CNAME   string `json:"course_name"`
	CPREREQ string `json:"prerequisite"`
}

// CSV file path for storing courses
var CSVFILE = "./courses.csv"

type CourseCatalog []MSDSCourse

var data = CourseCatalog{}
var index map[string]int
var db *sql.DB

func init() {
	var err error

	fmt.Println("Initializing the DB connection")

	db_connection := "user=postgres dbname=msds password=root host=/cloudsql/module8-452821:us-central1:mypostgres sslmode=disable port = 5432"
	db, err = sql.Open("postgres", db_connection)
	if err != nil {
		log.Fatal(fmt.Println("Couldn't Open Connection to database"))
		panic(err)
	}
	err = createTable()
	if err != nil {
		log.Fatal("Error creating table:", err)
	        panic(err)
	    }

    fmt.Println("Database and table initialized successfully")
	// Test the database connection
	//err = db.Ping()
	//if err != nil {
	//	fmt.Println("Couldn't Connect to database")
	//	panic(err)
	//}

}


// ✅ Function to Create Table in Cloud SQL
func createTable() error {
	query := `CREATE TABLE IF NOT EXISTS MSDSCourseCatalog (
		course_id TEXT PRIMARY KEY,
		course_name TEXT NOT NULL,
		prerequisite TEXT
	);`
	_, err := db.Exec(query)
	return err
}

// ✅ Function to Insert Course into Cloud SQL
func insertCourse(c *MSDSCourse) error {
	_, err := db.Exec("INSERT INTO MSDSCourseCatalog (course_id, course_name, prerequisite) VALUES ($1, $2, $3)",
		c.CID, c.CNAME, c.CPREREQ)
	return err
}

// ✅ Function to Delete Course from Cloud SQL
func deleteCourse(courseID string) error {
	_, err := db.Exec("DELETE FROM MSDSCourseCatalog WHERE course_id = $1", courseID)
	return err
}

// ✅ Function to Fetch All Courses
func fetchAllCourses() (string, error) {
	rows, err := db.Query("SELECT course_id, course_name, prerequisite FROM MSDSCourseCatalog")
	if err != nil {
		return "", err
	}
	defer rows.Close()

	var result string
	for rows.Next() {
		var c MSDSCourse
		if err := rows.Scan(&c.CID, &c.CNAME, &c.CPREREQ); err != nil {
			continue
		}
		result += fmt.Sprintf("Course ID: %s | Name: %s | Prerequisite: %s\n", c.CID, c.CNAME, c.CPREREQ)
	}
	return result, nil
}

// ✅ Function to Get Total Course Count
func getTotalCourses() (int, error) {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM MSDSCourseCatalog").Scan(&count)
	return count, err
}

// ✅ Function to Search for a Course
func searchCourse(courseID string) (*MSDSCourse, error) {
	var course MSDSCourse
	err := db.QueryRow("SELECT course_id, course_name, prerequisite FROM MSDSCourseCatalog WHERE course_id = $1", courseID).
		Scan(&course.CID, &course.CNAME, &course.CPREREQ)
	if err != nil {
		return nil, err
	}
	return &course, nil
}

// ✅ Main function - Start HTTP server
func main() {
	log.Print("Starting Cloud Run service...")

	mux := http.NewServeMux()
	mux.Handle("/list", http.HandlerFunc(listHandler))
	mux.Handle("/insert/", http.HandlerFunc(insertHandler))
	mux.Handle("/search/", http.HandlerFunc(searchHandler))
	mux.Handle("/delete/", http.HandlerFunc(deleteHandler))
	mux.Handle("/status", http.HandlerFunc(statusHandler))
	mux.Handle("/", http.HandlerFunc(defaultHandler)) // Default landing page

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Listening on port %s", port)

	s := &http.Server{
		Addr:         ":" + port,
		Handler:      mux,
		IdleTimeout:  10 * time.Second,
		ReadTimeout:  time.Second,
		WriteTimeout: time.Second,
	}

	err := s.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
