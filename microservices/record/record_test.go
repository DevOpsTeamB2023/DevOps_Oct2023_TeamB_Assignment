// record_test.go
package record

import (
	//change here
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
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
