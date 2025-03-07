package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

const PORT = ":1234"

func defaultHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Serving:", r.URL.Path, "from", r.Host)
	w.WriteHeader(http.StatusOK)
	Body := "Welcome to the MSDS Course Catalog!\n"
	fmt.Fprintf(w, "%s", Body)
}

// ✅ DELETE Handler: Remove a course by ID
func deleteHandler(w http.ResponseWriter, r *http.Request) {
	// Extract course ID from URL
	paramStr := strings.Split(r.URL.Path, "/")
	fmt.Println("Path:", paramStr)
	if len(paramStr) < 3 {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintln(w, "Not found: "+r.URL.Path)
		return
	}

	log.Println("Serving:", r.URL.Path, "from", r.Host)

	courseID := paramStr[2]
	err := deleteCourse(courseID) // Calls deleteCourse() from courses.go
	if err != nil {
		fmt.Println(err)
		Body := err.Error() + "\n"
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "%s", Body)
		return
	}

	Body := "Course " + courseID + " deleted!\n"
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%s", Body)
}

// ✅ GET Handler: List all courses
func listHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Serving:", r.URL.Path, "from", r.Host)
	w.WriteHeader(http.StatusOK)
	Body := listCourses() // Calls listCourses() from courses.go
	fmt.Fprintf(w, "%s", Body)
}

// ✅ STATUS Handler: Display total number of courses
func statusHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Serving:", r.URL.Path, "from", r.Host)
	w.WriteHeader(http.StatusOK)
	Body := fmt.Sprintf("Total courses: %d\n", len(data)) // Uses data from courses.go
	fmt.Fprintf(w, "%s", Body)
}

// ✅ POST Handler: Insert a new course
func insertHandler(w http.ResponseWriter, r *http.Request) {
	// Extract course details from URL path
	paramStr := strings.Split(r.URL.Path, "/")
	fmt.Println("Path:", paramStr)

	if len(paramStr) < 5 {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintln(w, "Not enough arguments: "+r.URL.Path)
		return
	}

	cid := paramStr[2]
	cname := paramStr[3]
	cprereq := paramStr[4]

	// Create a new course struct
	newCourse := &MSDSCourse{CID: cid, CNAME: cname, CPREREQ: cprereq}
	err := insertCourse(newCourse) // Calls insertCourse() from courses.go

	if err != nil {
		w.WriteHeader(http.StatusNotModified)
		Body := "Failed to add course\n"
		fmt.Fprintf(w, "%s", Body)
	} else {
		log.Println("Serving:", r.URL.Path, "from", r.Host)
		Body := "New course added successfully\n"
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "%s", Body)
	}

	log.Println("Serving:", r.URL.Path, "from", r.Host)
}

// ✅ GET Handler: Search for a course by ID
func searchHandler(w http.ResponseWriter, r *http.Request) {
	// Extract course ID from URL
	paramStr := strings.Split(r.URL.Path, "/")
	fmt.Println("Path:", paramStr)

	if len(paramStr) < 3 {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintln(w, "Not found: "+r.URL.Path)
		return
	}

	var Body string
	courseID := paramStr[2]
	c := searchCourse(courseID) // Calls searchCourse() from courses.go

	if c == nil {
		w.WriteHeader(http.StatusNotFound)
		Body = "Could not find course: " + courseID + "\n"
	} else {
		w.WriteHeader(http.StatusOK)
		Body = "Course ID: " + c.CID + "\nCourse Name: " + c.CNAME + "\nPrerequisite: " + c.CPREREQ + "\n"
	}

	fmt.Println("Serving:", r.URL.Path, "from", r.Host)
	fmt.Fprintf(w, "%s", Body)
}
