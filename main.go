package main

import (
	"database/sql"
	"encoding/csv"
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

	// Test the database connection
	//err = db.Ping()
	//if err != nil {
	//	fmt.Println("Couldn't Connect to database")
	//	panic(err)
	//}

}

// Read courses from CSV (creates if missing)
func readCSVFile(filepath string) error {
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		file, err := os.Create(filepath)
		if err != nil {
			return err
		}
		defer file.Close()
		fmt.Println("Created new CSV file:", filepath)
		return nil
	}

	f, err := os.Open(filepath)
	if err != nil {
		return err
	}
	defer f.Close()

	lines, err := csv.NewReader(f).ReadAll()
	if err != nil {
		return err
	}

	for _, line := range lines {
		temp := MSDSCourse{
			CID:     line[0],
			CNAME:   line[1],
			CPREREQ: line[2],
		}
		data = append(data, temp)
	}
	return nil
}

// Save courses to CSV
func saveCSVFile(filepath string) error {
	csvfile, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer csvfile.Close()

	csvwriter := csv.NewWriter(csvfile)
	for _, row := range data {
		temp := []string{row.CID, row.CNAME, row.CPREREQ}
		_ = csvwriter.Write(temp)
	}
	csvwriter.Flush()
	return nil
}

// Create an index for quick lookup
func createIndex() error {
	index = make(map[string]int)
	for i, k := range data {
		index[k.CID] = i
	}
	return nil
}

// Insert a new course
func insertCourse(c *MSDSCourse) error {
	if _, ok := index[c.CID]; ok {
		return fmt.Errorf("%s already exists", c.CID)
	}

	data = append(data, *c)
	_ = createIndex()

	err := saveCSVFile(CSVFILE)
	if err != nil {
		return err
	}
	return nil
}

// Delete a course
func deleteCourse(courseID string) error {
	i, ok := index[courseID]
	if !ok {
		return fmt.Errorf("%s not found!", courseID)
	}

	data = append(data[:i], data[i+1:]...)
	delete(index, courseID)

	err := saveCSVFile(CSVFILE)
	if err != nil {
		return err
	}
	return nil
}

// Search for a course
func searchCourse(courseID string) *MSDSCourse {
	i, ok := index[courseID]
	if !ok {
		return nil
	}
	return &data[i]
}

// List all courses
func listCourses() string {
	var all string
	for _, k := range data {
		all = all + k.CID + " | " + k.CNAME + " | Prerequisite: " + k.CPREREQ + "\n"
	}
	return all
}

// Main function - Start HTTP server
func main() {

	log.Print("Starting logging ...")
	err := readCSVFile(CSVFILE)
	if err != nil {
		fmt.Println(err)
		return
	}

	http.HandleFunc("/", handler)

	// Determine port for HTTP service.
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("defaulting to port %s", port)
	}

	// Start HTTP server.
	log.Printf("listening on port %s", port)
	log.Print("Navigate to Cloud Run services and find the URL of your service")
	log.Print("Use the browser and navigate to your service URL to to check your service has started")

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}

	err = createIndex()
	if err != nil {
		fmt.Println("Cannot create index.")
		return
	}

	mux := http.NewServeMux()
	s := &http.Server{
		Addr:         ":1234",
		Handler:      mux,
		IdleTimeout:  10 * time.Second,
		ReadTimeout:  time.Second,
		WriteTimeout: time.Second,
	}

	// Register handlers (ensure these exist in handlers.go)
	mux.Handle("/list", http.HandlerFunc(listHandler))
	mux.Handle("/insert/", http.HandlerFunc(insertHandler))
	mux.Handle("/insert", http.HandlerFunc(insertHandler))
	mux.Handle("/search", http.HandlerFunc(searchHandler))
	mux.Handle("/search/", http.HandlerFunc(searchHandler))
	mux.Handle("/delete/", http.HandlerFunc(deleteHandler))
	mux.Handle("/status", http.HandlerFunc(statusHandler))
	mux.Handle("/", http.HandlerFunc(defaultHandler))

	fmt.Println("Ready to serve at :1234")
	err = s.ListenAndServe()
	if err != nil {
		fmt.Println(err)
		return
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	name := os.Getenv("PROJECT_ID")
	if name == "" {
		name = "MSDS"
	}

	fmt.Fprintf(w, "started MSDS categlog %s!\n", name)
}
