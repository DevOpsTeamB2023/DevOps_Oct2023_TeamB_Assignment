// account_test.go
package account

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestCreateAccHandler(t *testing.T) {
	// Create a new mock database connection
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	// Replace the actual database connection with the mock
	SetDB(db)

	// Set up expected database query and result
	mock.ExpectPrepare("INSERT INTO Account").
		ExpectExec().
		WithArgs("testacc", "testpwd", "User", "Pending").
		WillReturnResult(sqlmock.NewResult(1, 1))

	newAcc := Account{
		Username:  "testacc",
		Password:  "testpwd",
		AccType:   "User",
		AccStatus: "Pending",
	}
	payload, err := json.Marshal(newAcc)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("POST", "/api/v1/accounts", bytes.NewBuffer(payload))
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Call the handler directly
	CreateAccHandler(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusCreated)
	}

	// Check the response body
	expected := "Account created successfully\n"
	if rr.Body.String() != expected {
		t.Errorf("Handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}

	// Verify that the expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	defer mock.ExpectationsWereMet() // Ensure expectations are checked even if the test fails early
}

func TestGetAccHandler(t *testing.T) {
	username := "testacc"
	password := "testpwd"

	// Create a new mock database connection
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	// Replace the actual database connection with the mock
	SetDB(db)

	// Set up expectations for the query and scan to return sql.ErrNoRows
	mock.ExpectQuery(regexp.QuoteMeta("SELECT AccID, Username, Password, AccType, AccStatus FROM Account WHERE Username = ? AND Password = ?")).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"AccID", "Username", "Password", "AccType", "AccStatus"}).
			AddRow(1, "testacc", "testpwd", "user", "active")) // Simulating a successful row

	req, err := http.NewRequest("GET", fmt.Sprintf("/api/v1/accounts?username=%s&password=%s", username, password), nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	GetAccHandler(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Unmarshal the response body into an Account struct
	var acc Account
	err = json.NewDecoder(rr.Body).Decode(&acc)
	if err != nil {
		t.Fatal(err)
	}

	// Check the retrieved account information
	expectedUsername := "testacc"
	if acc.Username != expectedUsername {
		t.Errorf("Handler returned unexpected username: got %v want %v", acc.Username, expectedUsername)
	}

	// Verify that the expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestApproveAccHandler(t *testing.T) {
	// accID follows the existing acc with pending status in record_db for testing approval
	accID := "2004"

	req, err := http.NewRequest("POST", fmt.Sprintf("/api/v1/accounts/approve?accID=%s", accID), nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	ApproveAccHandler(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check the response body
	expected := "Account approved successfully\n"
	if rr.Body.String() != expected {
		t.Errorf("Handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

func TestAdminCreateAccHandler(t *testing.T) {
	//dB()

	newAcc := Account{
		Username:  "admincreatedacc",
		Password:  "admincreatedpwd",
		AccType:   "User",
		AccStatus: "Created",
	}
	payload, err := json.Marshal(newAcc)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("POST", "/api/v1/accounts", bytes.NewBuffer(payload))
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Call the handler directly
	AdminCreateAccHandler(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusCreated)
	}

	// Check the response body
	expected := "Account created successfully\n"
	if rr.Body.String() != expected {
		t.Errorf("Handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

func TestDeleteAccHandler(t *testing.T) {
	// accID follows existing account for deletion with AccID=2003 in record_db for testing deletion
	accID := "2003"

	req, err := http.NewRequest("DELETE", fmt.Sprintf("/api/v1/accounts/delete?accID=%s", accID), nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	DeleteAccHandler(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check the response body
	expected := "Account deleted successfully\n"
	if rr.Body.String() != expected {
		t.Errorf("Handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

// not working - request not passing to account.go
func TestUpdateAccHandler(t *testing.T) {
	// accID follows existing account for update with AccID=2005 in record_db for testing update
	accID := "2005"

	// Create a request with a JSON payload for updating the account
	updatedAcc := Account{
		Username: "testupdatepass",
		AccType:  "Admin",
	}
	payload, err := json.Marshal(updatedAcc)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("/api/v1/accounts/%s", accID), bytes.NewBuffer(payload))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	UpdateAccHandler(rr, req)

	if status := rr.Code; status != http.StatusAccepted {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusAccepted)
	}

	// Check the response body
	expected := "Account updated successfully!\n"
	if rr.Body.String() != expected {
		t.Errorf("Handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}
