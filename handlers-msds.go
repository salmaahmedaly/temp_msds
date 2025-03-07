package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"database/sql"
)

const PORT = ":1234"


// Default handler - Landing page
func defaultHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Serving:", r.URL.Path, "from", r.Host)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Welcome to the MSDS Course Catalog API!\n")
}

// ✅ DELETE Handler: Remove a course by ID
func deleteHandler(w http.ResponseWriter, r *http.Request) {
	paramStr := strings.Split(r.URL.Path, "/")
	if len(paramStr) < 3 {
		http.Error(w, "Not enough arguments. Usage: /delete/CID", http.StatusBadRequest)
		return
	}

	courseID := paramStr[2]
	err := deleteCourse(courseID) // Calls deleteCourse() from main.go

	if err != nil {
		http.Error(w, "Failed to delete course: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Course %s deleted!\n", courseID)
}

// ✅ GET Handler: List all courses from Cloud SQL
func listHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Serving:", r.URL.Path, "from", r.Host)

	courses, err := fetchAllCourses() // Calls fetchAllCourses() from main.go
	if err != nil {
		http.Error(w, "Failed to fetch courses: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%s", courses)
}

// ✅ STATUS Handler: Display total number of courses
func statusHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Serving:", r.URL.Path, "from", r.Host)

	count, err := getTotalCourses() // Calls getTotalCourses() from main.go
	if err != nil {
		http.Error(w, "Failed to get course count: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Total courses: %d\n", count)
}

// ✅ POST Handler: Insert a new course into Cloud SQL
func insertHandler(w http.ResponseWriter, r *http.Request) {
	paramStr := strings.Split(r.URL.Path, "/")

	if len(paramStr) < 5 {
		http.Error(w, "Not enough arguments. Usage: /insert/CID/CNAME/CPREREQ", http.StatusBadRequest)
		return
	}

	cid, cname, cprereq := paramStr[2], paramStr[3], paramStr[4]

	newCourse := &MSDSCourse{CID: cid, CNAME: cname, CPREREQ: cprereq}
	err := insertCourse(newCourse) // Calls insertCourse() from main.go

	if err != nil {
		http.Error(w, "Failed to insert course: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Course added successfully: %s\n", cid)
}

// ✅ GET Handler: Search for a course by ID
func searchHandler(w http.ResponseWriter, r *http.Request) {
	paramStr := strings.Split(r.URL.Path, "/")

	if len(paramStr) < 3 {
		http.Error(w, "Not enough arguments. Usage: /search/CID", http.StatusBadRequest)
		return
	}

	courseID := paramStr[2]
	course, err := searchCourse(courseID) // Calls searchCourse() from main.go

	if err == sql.ErrNoRows {
		http.Error(w, "Course not found: "+courseID, http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Course ID: %s\nCourse Name: %s\nPrerequisite: %s\n", course.CID, course.CNAME, course.CPREREQ)
}
