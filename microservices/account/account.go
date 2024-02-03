package account

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type Account struct {
	AccID     int    `json:"accId"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	AccType   string `json:"accType"`
	AccStatus string `json:"accStatus"`
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
	router.HandleFunc("/api/v1/accounts", CreateAccHandler).Methods("POST")
	router.HandleFunc("/api/v1/accounts", GetAccHandler).Methods("GET")
	router.HandleFunc("/api/v1/accounts/all", ListAllAccsHandler).Methods("GET")
	router.HandleFunc("/api/v1/accounts/approve", ApproveAccHandler).Methods("POST")
	router.HandleFunc("/api/v1/accounts", AdminCreateAccHandler).Methods("POST")
	router.HandleFunc("/api/v1/accounts/delete", DeleteAccHandler).Methods("DELETE")
	router.HandleFunc("/api/v1/accounts/get", GetSpecificAccHandler).Methods("GET")
	router.HandleFunc("/api/v1/accounts/{accID}", UpdateAccHandler).Methods("PUT")

	fmt.Println("Listening at port 5001")
	http.ListenAndServe(":5001",
		handlers.CORS(
			handlers.AllowedOrigins([]string{"*"}),
			handlers.AllowedMethods([]string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}),
			handlers.AllowedHeaders([]string{"Origin", "X-Api-Key", "X-Requested-With", "Content-Type", "Accept", "Authorization"}),
			handlers.AllowCredentials(),
		)(router))
}

func CreateAccHandler(w http.ResponseWriter, r *http.Request) {
	var newAcc Account
	err := json.NewDecoder(r.Body).Decode(&newAcc)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Insert the new account into the database
	stmt, err := db.Prepare("INSERT INTO Account (Username, Password, AccType, AccStatus) VALUES (?, ?, ?, ?)")
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(newAcc.Username, newAcc.Password, newAcc.AccType, newAcc.AccStatus)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintln(w, "Account created successfully")
}

func GetAccHandler(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	password := r.URL.Query().Get("password")

	if username == "" || password == "" {
		http.Error(w, "Username and Password parameters are required", http.StatusBadRequest)
		return
	}

	var acc Account
	err := db.QueryRow("SELECT AccID, Username, Password, AccType, AccStatus FROM Account WHERE Username = ? AND Password = ?", username, password).Scan(&acc.AccID, &acc.Username, &acc.Password, &acc.AccType, &acc.AccStatus)
	if err == sql.ErrNoRows {
		http.Error(w, "Anvalid Username or Password", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "Bnternal server error", http.StatusInternalServerError)
		return
	}

	// Respond with user information
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(acc)
}

func ListAllAccsHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT AccID, Username, AccType, AccStatus FROM Account")
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var accs []Account
	for rows.Next() {
		var acc Account
		err := rows.Scan(&acc.AccID, &acc.Username, &acc.AccType, &acc.AccStatus)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		accs = append(accs, acc)
	}

	// Respond with the list of users
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(accs)
}

func ApproveAccHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the account ID from the request parameters
	accID := r.URL.Query().Get("accID")
	if accID == "" {
		http.Error(w, "Account ID parameter is required", http.StatusBadRequest)
		return
	}

	// Update the account status in the database
	stmt, err := db.Prepare("UPDATE Account SET AccStatus = 'Created' WHERE AccID = ?")
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(accID)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Account approved successfully")
}

func AdminCreateAccHandler(w http.ResponseWriter, r *http.Request) {
	var newAcc Account
	err := json.NewDecoder(r.Body).Decode(&newAcc)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Insert the new account into the database
	stmt, err := db.Prepare("INSERT INTO Account (Username, Password, AccType, AccStatus) VALUES (?, ?, ?, ?)")
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(newAcc.Username, newAcc.Password, newAcc.AccType, newAcc.AccStatus)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintln(w, "Account created successfully")
}

func DeleteAccHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the account ID from the request parameters
	accID := r.URL.Query().Get("accID")
	if accID == "" {
		http.Error(w, "Account ID parameter is required", http.StatusBadRequest)
		return
	}

	// Delete the account from the database
	stmt, err := db.Prepare("DELETE FROM Account WHERE AccID = ?")
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(accID)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Account deleted successfully")
}

func GetSpecificAccHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the account ID from the request parameters
	accID := r.URL.Query().Get("accID")

	// get the account from the database
	var acc Account
	db.QueryRow("Select * FROM Account WHERE AccID = ?", accID).Scan(&acc.AccID, &acc.Username, &acc.Password, &acc.AccType, &acc.AccStatus)

	w.Header().Set("Content-Type", "application/json")
	fmt.Println(acc)
	json.NewEncoder(w).Encode(acc)
}

func UpdateAccHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the user ID from the request URL
	vars := mux.Vars(r)
	accID, err := strconv.Atoi(vars["accID"])
	if err != nil {
		http.Error(w, "Invalid Account ID", http.StatusBadRequest)
		return
	}

	var updatedAcc Account
	err = json.NewDecoder(r.Body).Decode(&updatedAcc)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Update the user's information in the database
	stmt, err := db.Prepare("UPDATE Account SET Username=?, AccType=? WHERE AccID=?")
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(updatedAcc.Username, updatedAcc.AccType, accID)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
	fmt.Fprintln(w, "Account updated successfully!")
}
