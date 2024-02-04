package record

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type Record struct {
	RecordID       int    `json:"recordId"`
	Name           string `json:"name"`
	RoleOfContact  string `json:"roleOfContact"`
	NoOfStudents   int    `json:"noOfStudents"`
	AcadYr         string `json:"acadYr"`
	CapstoneTitle  string `json:"capstoneTitle"`
	CompanyName    string `json:"companyName"`
	CompanyContact string `json:"companyContact"`
	ProjDesc       string `json:"projDesc"`
}

var (
	db  *sql.DB
	err error
)

func SetDB(database *sql.DB) {
	db = database
}

func DB() {
	db, err = sql.Open("mysql", "record_system:dopasgpwd@tcp(127.0.0.1:3306)/record_db")
	if err != nil {
		log.Fatal(err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to the database")
}

func InitHTTPServer() {
	DB()

	router := mux.NewRouter()
	router.Use(corsMiddleware)

	router.HandleFunc("/api/v1/records/all", ListAllRecordsHandler).Methods("GET")
	router.HandleFunc("/api/v1/records", CreateRecordHandler).Methods("POST")
	router.HandleFunc("/api/v1/records/delete", DeleteRecordHandler).Methods("DELETE")
	router.HandleFunc("/api/v1/records/{recordID}", UpdateRecordHandler).Methods("PUT")
	router.HandleFunc("/api/v1/records/search", QueryRecordHandler).Methods("GET")

	fmt.Println("Listening at port 5002")
	go func() {
		log.Fatal(http.ListenAndServe(":5002", router))
	}()
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, PUT, PATCH, GET, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Api-Key, X-Requested-With, Content-Type, Accept, Authorization")
		next.ServeHTTP(w, r)
	})
}

// gets and lists all capstone records
func ListAllRecordsHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT RecordID, Name, RoleOfContact, NoOfStudents, AcadYr, CapstoneTitle, CompanyName, CompanyContact, ProjDesc FROM Record")
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var records []Record
	for rows.Next() {
		var record Record
		err := rows.Scan(&record.RecordID, &record.Name, &record.RoleOfContact, &record.NoOfStudents, &record.AcadYr, &record.CapstoneTitle, &record.CompanyName, &record.CompanyContact, &record.ProjDesc)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		records = append(records, record)
	}

	// Respond with the list of records
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(records)
}

// create a capstone record
func CreateRecordHandler(w http.ResponseWriter, r *http.Request) {
	var newRecord Record
	err := json.NewDecoder(r.Body).Decode(&newRecord)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Insert the new record into the database
	stmt, err := db.Prepare("INSERT INTO Record (Name, RoleOfContact, NoOfStudents, AcadYr, CapstoneTitle, CompanyName, CompanyContact, ProjDesc) VALUES (?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(newRecord.Name, newRecord.RoleOfContact, newRecord.NoOfStudents, newRecord.AcadYr, newRecord.CapstoneTitle, newRecord.CompanyName, newRecord.CompanyContact, newRecord.ProjDesc)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintln(w, "Record created successfully")
}

func DeleteRecordHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the record ID from the request parameters
	recordID := r.URL.Query().Get("recordID")
	if recordID == "" {
		http.Error(w, "Record ID parameter is required", http.StatusBadRequest)
		return
	}

	// Delete the record from the database
	stmt, err := db.Prepare("DELETE FROM Record WHERE RecordID = ?")
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(recordID)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Record deleted successfully")
}

func UpdateRecordHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the record ID from the request URL
	vars := mux.Vars(r)
	recordID, err := strconv.Atoi(vars["recordID"])
	if err != nil {
		http.Error(w, "Invalid record ID", http.StatusBadRequest)
		return
	}

	var updatedRecord Record
	err = json.NewDecoder(r.Body).Decode(&updatedRecord)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Update the record's information in the database
	stmt, err := db.Prepare("UPDATE Record SET Name=?, RoleOfContact=?, NoOfStudents=?, AcadYr=?, CapstoneTitle=?, CompanyName=?, CompanyContact=?, ProjDesc=? WHERE RecordID=?")
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(updatedRecord.Name, updatedRecord.RoleOfContact, updatedRecord.NoOfStudents, updatedRecord.AcadYr, updatedRecord.CapstoneTitle, updatedRecord.CompanyName, updatedRecord.CompanyContact, updatedRecord.ProjDesc, recordID)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
	fmt.Fprintln(w, "Record updated successfully!")
}

func QueryRecordByAcadYrHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the acadYr from the query parameters
	acadYr := r.URL.Query().Get("acadYr")

	// Query the database to search for trips based on the acadYr
	rows, err := db.Query("SELECT RecordID, Name, RoleOfContact, NoOfStudents, AcadYr, CapstoneTitle, CompanyName, CompanyContact, ProjDesc FROM Record WHERE AcadYr LIKE ?", "%"+acadYr+"%")
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Create a slice to store the search results
	var searchResults []Record

	// Iterate through the rows and populate the search results slice
	for rows.Next() {
		var record Record
		if err := rows.Scan(&record.RecordID, &record.Name, &record.RoleOfContact, &record.NoOfStudents, &record.AcadYr, &record.CapstoneTitle, &record.CompanyName, &record.CompanyContact, &record.ProjDesc); err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		searchResults = append(searchResults, record)
	}

	// Check for errors during row iteration
	if err := rows.Err(); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Encode the search results as JSON and send the response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(searchResults)
}

func QueryRecordHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the capstoneTitle from the query parameters
	query := r.URL.Query().Get("query")

	// Query the database to search for trips based on the acadYr
	rows, err := db.Query("SELECT RecordID, Name, RoleOfContact, NoOfStudents, AcadYr, CapstoneTitle, CompanyName, CompanyContact, ProjDesc FROM Record WHERE AcadYr LIKE ? OR CapstoneTitle LIKE ?", "%"+query+"%", "%"+query+"%")
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Create a slice to store the search results
	var searchResults []Record

	// Iterate through the rows and populate the search results slice
	for rows.Next() {
		var record Record
		if err := rows.Scan(&record.RecordID, &record.Name, &record.RoleOfContact, &record.NoOfStudents, &record.AcadYr, &record.CapstoneTitle, &record.CompanyName, &record.CompanyContact, &record.ProjDesc); err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		searchResults = append(searchResults, record)
	}

	// Check for errors during row iteration
	if err := rows.Err(); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Encode the search results as JSON and send the response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(searchResults)
}
