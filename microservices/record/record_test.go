// record_test.go
package record

import (
	//change here

	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-sql-driver/mysql"
)

func TestDeleteRecordHandler(t *testing.T) {
	// Create a new mock database connection
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	SetDB(db)

	// Set up expectations for the Prepare call
	mock.ExpectPrepare("DELETE FROM Record WHERE RecordID = ?")

	mock.ExpectExec("DELETE FROM Record WHERE RecordID = ?").WithArgs(sqlmock.AnyArg()).WillReturnResult(sqlmock.NewResult(0, 1))

	// recordID follows existing record for deletion with recordID=4 in record_db for testing deletion
	recordID := "3"

	req, err := http.NewRequest("DELETE", fmt.Sprintf("/api/v1/records/delete?recordID=%s", recordID), nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	DeleteRecordHandler(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check the response body
	expected := "Record deleted successfully\n"
	if rr.Body.String() != expected {
		t.Errorf("Handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

func TestDeleteRecordHandler_NoID(t *testing.T) {
	// accID follows the existing acc with pending status in record_db for testing approval
	recordID := ""

	// Create a new mock database connection
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	SetDB(db)

	// Set up expectations for the Prepare call
	mock.ExpectPrepare("DELETE FROM Record WHERE RecordID = ?")

	mock.ExpectExec("DELETE FROM Record WHERE RecordID = ?").WithArgs(sqlmock.AnyArg()).WillReturnResult(sqlmock.NewResult(0, 1))

	req, err := http.NewRequest("DELETE", fmt.Sprintf("/api/v1/records/delete?recordID=%s", recordID), nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	DeleteRecordHandler(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}
}

func TestDeleteRecordHandler_Prepare(t *testing.T) {
	// Create a new mock database connection
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	SetDB(db)

	// Create a mock MySQL error
	mockError := &mysql.MySQLError{
		Number:  1062,                                      // MySQL error number (example)
		Message: "Duplicate entry 'xyz' for key 'PRIMARY'", // MySQL error message (example)
	}

	// Mock Prepare to return the mock MySQL error
	mock.ExpectPrepare("DELETE FROM Record").WillReturnError(mockError)

	// recordID follows existing record for deletion with recordID=4 in record_db for testing deletion
	recordID := "3"

	req, err := http.NewRequest("DELETE", fmt.Sprintf("/api/v1/records/delete?recordID=%s", recordID), nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	// Call the handler directly
	DeleteRecordHandler(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusInternalServerError)
	}
}

func TestDeleteRecordHandler_Exec(t *testing.T) {
	// Create a new mock database connection
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	SetDB(db)

	// Create a mock MySQL error
	mockError := &mysql.MySQLError{
		Number:  1062,                                      // MySQL error number (example)
		Message: "Duplicate entry 'xyz' for key 'PRIMARY'", // MySQL error message (example)
	}

	// Set up expectations for your query
	mock.ExpectPrepare("DELETE FROM Record").ExpectExec().
		WillReturnError(mockError)

	// recordID follows existing record for deletion with recordID=4 in record_db for testing deletion
	recordID := "3"

	req, err := http.NewRequest("DELETE", fmt.Sprintf("/api/v1/records/delete?recordID=%s", recordID), nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	// Call the handler directly
	DeleteRecordHandler(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusInternalServerError)
	}
}

func TestListAllRecordsHandler(t *testing.T) {
	// Create a new mock database connection
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	// Replace the actual database connection with the mock
	SetDB(db)

	// Test case for successful query execution
	t.Run("Success", func(t *testing.T) {
		// Set up expected database query and result
		rows := sqlmock.NewRows([]string{"RecordID", "Name", "RoleOfContact", "NoOfStudents", "AcadYr", "CapstoneTitle", "CompanyName", "CompanyContact", "ProjDesc"}).
			AddRow(1, "Test Name1", "Student", 3, "2022/2023", "Title1", "Company1", "Contact Name1", "Description").
			AddRow(2, "Test Name2", "Staff", 4, "2023/2024", "Title2", "Company2", "Contact Name2", "Description")

		mock.ExpectQuery("SELECT RecordID, Name, RoleOfContact, NoOfStudents, AcadYr, CapstoneTitle, CompanyName, CompanyContact, ProjDesc FROM Record").
			WillReturnRows(rows)

		req, err := http.NewRequest("GET", "/api/v1/records", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()

		// Call the handler
		ListAllRecordsHandler(rr, req)

		// Check the status code
		if status := rr.Code; status != http.StatusOK {
			t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}

		// Verify that the expectations were met
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	// Test case for database query error
	t.Run("DatabaseError", func(t *testing.T) {
		// Set up mock to return an error
		mock.ExpectQuery("SELECT RecordID, Name, RoleOfContact, NoOfStudents, AcadYr, CapstoneTitle, CompanyName, CompanyContact, ProjDesc FROM Record").
			WillReturnError(errors.New("database error"))

		req, err := http.NewRequest("GET", "/api/v1/records", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()

		// Call the handler
		ListAllRecordsHandler(rr, req)

		// Check the status code
		if status := rr.Code; status != http.StatusInternalServerError {
			t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusInternalServerError)
		}

		// Verify that the expectations were met
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})
}

func TestCreateRecordHandler(t *testing.T) {
	// Create a new mock database connection
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	// Replace the actual database connection with the mock
	SetDB(db)

	// Set up expected database query and result
	mock.ExpectPrepare(regexp.QuoteMeta("INSERT INTO Record (Name, RoleOfContact, NoOfStudents, AcadYr, CapstoneTitle, CompanyName, CompanyContact, ProjDesc) VALUES (?, ?, ?, ?, ?, ?, ?, ?)")).
		ExpectExec().
		WithArgs("Test Create Reecord", "Student", 3, "2022/2023", "Title", "Company", "Contact Name", "Description").
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Create JSON request body
	requestBody := `{"Name": "Test Create Reecord", "RoleOfContact": "Student", "NoOfStudents": 3, "AcadYr": "2022/2023", "CapstoneTitle": "Title", "CompanyName": "Company", "CompanyContact": "Contact Name", "ProjDesc": "Description"}`

	req, err := http.NewRequest("POST", "/api/v1/records", strings.NewReader(requestBody))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	// Call the handler
	CreateRecordHandler(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusCreated)
	}

	// Verify that the expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestCreateRecordHandler_Error(t *testing.T) {
	// Create a new mock database connection
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	// Replace the actual database connection with the mock
	SetDB(db)

	// Simulate a database error
	mock.ExpectPrepare(regexp.QuoteMeta("INSERT INTO Record (Name, RoleOfContact, NoOfStudents, AcadYr, CapstoneTitle, CompanyName, CompanyContact, ProjDesc) VALUES (?, ?, ?, ?, ?, ?, ?, ?)")).
		ExpectExec().
		WithArgs("Create Error", "Staff", 3, "2022/2023", "Title", "Company", "Contact Name", "Description").
		WillReturnError(fmt.Errorf("database error"))

	// Create JSON request body
	requestBody := `{"Name": "Create Error", "RoleOfContact": "Staff", "NoOfStudents": 3, "AcadYr": "2022/2023", "CapstoneTitle": "Title", "CompanyName": "Company", "CompanyContact": "Contact Name", "ProjDesc": "Description"}`

	req, err := http.NewRequest("POST", "/api/v1/records", strings.NewReader(requestBody))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	// Call the handler
	CreateRecordHandler(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusInternalServerError)
	}

	// Check the response body
	expectedBody := "Internal server error\n"
	if rr.Body.String() != expectedBody {
		t.Errorf("Handler returned unexpected body: got %v want %v", rr.Body.String(), expectedBody)
	}

	// Verify that the expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestCreateRecordHandler_InvalidRequestPayload(t *testing.T) {
	req, err := http.NewRequest("POST", "/api/v1/records", strings.NewReader("invalid json"))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	// Call the handler
	CreateRecordHandler(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}

	// Check the response body
	expectedBody := "Invalid request payload\n"
	if rr.Body.String() != expectedBody {
		t.Errorf("Handler returned unexpected body: got %v want %v", rr.Body.String(), expectedBody)
	}
}

func TestCreateRecordHandler_ErrorPreparingStatement(t *testing.T) {
	// Create a new mock database connection
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	// Replace the actual database connection with the mock
	SetDB(db)

	// Simulate an error when preparing the SQL statement
	mock.ExpectPrepare(regexp.QuoteMeta("INSERT INTO Record (Name, RoleOfContact, NoOfStudents, AcadYr, CapstoneTitle, CompanyName, CompanyContact, ProjDesc) VALUES (?, ?, ?, ?, ?, ?, ?, ?)")).
		WillReturnError(fmt.Errorf("failed to prepare statement"))

	// Create JSON request body
	requestBody := `{"Name": "Create Error", "RoleOfContact": "Staff", "NoOfStudents": 3, "AcadYr": "2022/2023", "CapstoneTitle": "Title", "CompanyName": "Company", "CompanyContact": "Contact Name", "ProjDesc": "Description"}`

	req, err := http.NewRequest("POST", "/api/v1/records", strings.NewReader(requestBody))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	// Call the handler
	CreateRecordHandler(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusInternalServerError)
	}

	// Check the response body
	expectedBody := "Internal server error\n"
	if rr.Body.String() != expectedBody {
		t.Errorf("Handler returned unexpected body: got %v want %v", rr.Body.String(), expectedBody)
	}

	// Verify that the expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
