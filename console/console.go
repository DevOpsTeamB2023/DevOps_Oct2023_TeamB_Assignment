package main

import (
	"DevOps_Oct2023_TeamB_Assignment/microservices/account" //change here
	"DevOps_Oct2023_TeamB_Assignment/microservices/record"  //change here
)

func main() {
	go account.InitHTTPServer()
	go record.InitHTTPServer()

	select {}
}

// to check if backend works (developing using console first before changing to html)
// package main

// import (
// 	"DevOps_Oct2023_TeamB_Assignment/microservices/account" //change here
// 	"DevOps_Oct2023_TeamB_Assignment/microservices/record"  //change here
// 	"bufio"
// 	"bytes"
// 	"encoding/json"
// 	"fmt"
// 	"net/http"
// 	"os"
// )

// type Account struct {
// 	AccID     int    `json:"accId"`
// 	Username  string `json:"username"`
// 	Password  string `json:"password"`
// 	AccType   string `json:"accType"`
// 	AccStatus string `json:"accStatus"`
// }

// type Record struct {
// 	RecordID       int    `json:"recordId"`
// 	Name           string `json:"name"`
// 	RoleOfContact  string `json:"roleOfContact"`
// 	NoOfStudents   int    `json:"noOfStudents"`
// 	AcadYr         string `json:"acadYr"`
// 	CapstoneTitle  string `json:"capstoneTitle"`
// 	CompanyName    string `json:"companyName"`
// 	CompanyContact string `json:"companyContact"`
// 	ProjDesc       string `json:"projDesc"`
// }

// func main() {
// 	go account.InitHTTPServer()
// 	go record.InitHTTPServer()

// outer:
// 	for {
// 		fmt.Println("===============================================")
// 		fmt.Println("Welcome to the Capstone Records System!")
// 		fmt.Println("1. Create User Account")
// 		fmt.Println("2. Login")
// 		fmt.Println("0. Exit")
// 		fmt.Print("Enter an option: ")

// 		var choice int
// 		fmt.Scanf("%d", &choice)

// 		switch choice {
// 		case 1:
// 			// user account creation
// 			fmt.Println("----Create User Account----")
// 			createAcc()
// 		case 2:
// 			// user login
// 			fmt.Println("----Login----")
// 			acc, err := login()
// 			if err != nil {
// 				fmt.Println("Login failed:", err)
// 				return
// 			}
// 			//after login display user main menu
// 			if acc.AccStatus == "Created" {
// 				if acc.AccType == "User" {
// 					userMainMenu()
// 				} else {
// 					adminMainMenu()
// 				}
// 			} else {
// 				fmt.Println("Your account has not been approved yet. Please try again later.")
// 			}
// 		case 0:
// 			break outer
// 		default:
// 			fmt.Println("Invalid option")
// 		}
// 		fmt.Scanln()
// 	}
// }

// // creates user account
// func createAcc() {
// 	var acc Account
// 	reader := bufio.NewReader(os.Stdin)
// 	reader.ReadString('\n')
// 	fmt.Print("Enter Username: ")
// 	fmt.Scanf("%v", &acc.Username)
// 	reader.ReadString('\n')
// 	fmt.Print("Enter Password: ")
// 	fmt.Scanf("%v", &acc.Password)

// 	acc.AccType = "User"
// 	acc.AccStatus = "Pending"

// 	postBody, _ := json.Marshal(acc)

// 	client := &http.Client{}
// 	if req, err := http.NewRequest(http.MethodPost, "http://localhost:5001/api/v1/accounts", bytes.NewBuffer(postBody)); err == nil {
// 		if res, err := client.Do(req); err == nil {
// 			if res.StatusCode == 201 {
// 				fmt.Println("Account request sent. Please wait for admin approval.")
// 			} else {
// 				fmt.Println("Error creating user account")
// 			}
// 		} else {
// 			fmt.Println(2, err)
// 		}
// 	} else {
// 		fmt.Println(3, err)
// 	}

// }

// func login() (*Account, error) {
// 	var (
// 		username string
// 		password string
// 	)
// 	reader := bufio.NewReader(os.Stdin)

// 	reader.ReadString('\n')
// 	fmt.Print("Enter Username: ")
// 	fmt.Scanf("%v", &username)

// 	reader.ReadString('\n')
// 	fmt.Print("Enter Password: ")
// 	fmt.Scanf("%v", &password)

// 	// Perform login check
// 	client := &http.Client{}
// 	if req, err := http.NewRequest(http.MethodGet, "http://localhost:5001/api/v1/accounts?username="+username+"&password="+password, nil); err == nil {
// 		if res, err := client.Do(req); err == nil {
// 			defer res.Body.Close()

// 			if res.StatusCode == http.StatusOK {
// 				var acc Account
// 				err := json.NewDecoder(res.Body).Decode(&acc)
// 				if err == nil {
// 					fmt.Printf("Welcome back, %s!\n", acc.Username)
// 					return &acc, nil
// 				} else {
// 					return nil, fmt.Errorf("Error decoding response: %v", err)
// 				}
// 			} else {
// 				return nil, fmt.Errorf("Inavlid Username or Password")
// 			}
// 		} else {
// 			return nil, fmt.Errorf("Error making request: %v", err)
// 		}
// 	} else {
// 		return nil, fmt.Errorf("Error creating request: %v", err)
// 	}
// }

// // user features - main menu selection
// func userMainMenu() {
// 	for {
// 		var choice int
// 		reader := bufio.NewReader(os.Stdin)
// 		fmt.Println("===============================================")
// 		fmt.Println("--------------User Main Menu-------------")
// 		// list all capstone entries
// 		listAllRecords()
// 		// feature selection
// 		fmt.Println("\n1. Create a Capstone Entry")
// 		fmt.Println("2. Search") //search based on acad year and/or keywords -> displays project title and name of person in charge
// 		fmt.Println("3. Edit Capstone Entry")
// 		fmt.Println("4. Delete Capstone Entry")
// 		fmt.Println("0. Exit")
// 		reader.ReadString('\n')
// 		fmt.Print("Enter an option: ")

// 		fmt.Scanf("%d", &choice)

// 		switch choice {
// 		case 1:
// 			// Create a capstone entry
// 			fmt.Println("----Create Capstone Entry----")
// 			createRecord()
// 		case 2:
// 			// Search
// 			fmt.Println("----Search----")
// 			queryRecord()
// 		case 3:
// 			// Edit capstone entry
// 			fmt.Println("----Edit Capstone Entry----")
// 			editRecord()
// 		case 4:
// 			// Delete capstone entry
// 			fmt.Println("----Delete Capstone Entry----")
// 			deleteRecord()
// 		case 0:
// 			os.Exit(0)
// 		default:
// 			fmt.Println("Invalid option")
// 		}
// 	}
// }

// // main page after logging in as admin
// func adminMainMenu() {
// 	for {
// 		var choice int
// 		reader := bufio.NewReader(os.Stdin)
// 		fmt.Println("===============================================")
// 		fmt.Println("--------------Admin Main Menu--------------")
// 		fmt.Println("1. List all User Accounts")
// 		fmt.Println("2. List all Capstone Entries")
// 		fmt.Println("0. Exit")
// 		reader.ReadString('\n')
// 		fmt.Print("Enter an option: ")

// 		fmt.Scanf("%d", &choice)

// 		switch choice {
// 		case 1:
// 			// List all user accounts
// 			fmt.Println("----All User Accounts----")
// 			listAllAccs()
// 			manageAccsMenu()
// 		case 2:
// 			// List all capstone entries
// 			fmt.Println("----All Capstone Entries----")
// 			listAllRecords()
// 			manageRecordsMenu()
// 		case 0:
// 			os.Exit(0)
// 		default:
// 			fmt.Println("Invalid option")
// 		}
// 	}
// }

// // admin account functions
// func listAllAccs() error {
// 	// Perform list all users request
// 	client := &http.Client{}
// 	if req, err := http.NewRequest(http.MethodGet, "http://localhost:5001/api/v1/accounts/all", nil); err == nil {
// 		if res, err := client.Do(req); err == nil {
// 			defer res.Body.Close()

// 			if res.StatusCode == http.StatusOK {
// 				var accs []Account
// 				err := json.NewDecoder(res.Body).Decode(&accs)
// 				if err == nil {
// 					fmt.Println("List of all users:")
// 					for _, acc := range accs {
// 						fmt.Printf("Account ID: %d \nUsername: %s \nAccount Type: %s \nAccount Status: %s\n\n", acc.AccID, acc.Username, acc.AccType, acc.AccStatus)
// 					}
// 					return nil
// 				} else {
// 					return fmt.Errorf("Error decoding response: %v", err)
// 				}
// 			} else {
// 				return fmt.Errorf("Error fetching user list")
// 			}
// 		} else {
// 			return fmt.Errorf("Error making request: %v", err)
// 		}
// 	} else {
// 		return fmt.Errorf("Error creating request: %v", err)
// 	}
// }

// func manageAccsMenu() {
// 	var choice int
// 	reader := bufio.NewReader(os.Stdin)
// 	fmt.Println("\nManage Accounts:")
// 	fmt.Println("1. Approve Pending User Account")
// 	fmt.Println("2. Modify User Account")
// 	fmt.Println("3. Delete User Account")
// 	fmt.Println("4. Create User Account")
// 	fmt.Println("0. Go Back")
// 	reader.ReadString('\n')
// 	fmt.Print("Enter an option: ")
// 	fmt.Scanf("%d", &choice)

// 	switch choice {
// 	case 1:
// 		//get accID and update accStatus based on selected accID
// 		fmt.Println("----Approve Account----")
// 		approveAcc()
// 	case 2:
// 		//get accID and allow modifications to acc details based on selected accID
// 		fmt.Println("----Modify Account----")
// 		editAcc()
// 	case 3:
// 		//get accID and delete selected account
// 		fmt.Println("----Delete Account----")
// 		deleteAcc()
// 	case 4:
// 		//set status as created since done by admin
// 		fmt.Println("----Create Account----")
// 		adminCreateAcc()
// 	case 0:
// 		// Go back
// 	default:
// 		fmt.Println("Invalid option")
// 	}
// }

// func approveAcc() {
// 	var accID int

// 	reader := bufio.NewReader(os.Stdin)

// 	reader.ReadString('\n')
// 	fmt.Print("Enter Account ID to approve: ")
// 	fmt.Scanf("%d", &accID)

// 	client := &http.Client{}
// 	if req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("http://localhost:5001/api/v1/accounts/approve?accID=%d", accID), nil); err == nil {
// 		if res, err := client.Do(req); err == nil {
// 			defer res.Body.Close()

// 			if res.StatusCode == http.StatusOK {
// 				fmt.Println("Account approved successfully")
// 			} else {
// 				fmt.Println("Error approving account")
// 			}
// 		} else {
// 			fmt.Println("Error making request:", err)
// 		}
// 	} else {
// 		fmt.Println("Error creating request:", err)
// 	}
// }

// // admin creates user acc
// func adminCreateAcc() {
// 	var acc Account
// 	reader := bufio.NewReader(os.Stdin)
// 	reader.ReadString('\n')
// 	fmt.Print("Enter Username: ")
// 	fmt.Scanf("%v", &acc.Username)
// 	reader.ReadString('\n')
// 	fmt.Print("Enter Password: ")
// 	fmt.Scanf("%v", &acc.Password)

// 	acc.AccType = "User"
// 	acc.AccStatus = "Created"

// 	postBody, _ := json.Marshal(acc)

// 	client := &http.Client{}
// 	if req, err := http.NewRequest(http.MethodPost, "http://localhost:5001/api/v1/accounts", bytes.NewBuffer(postBody)); err == nil {
// 		if res, err := client.Do(req); err == nil {
// 			if res.StatusCode == 201 {
// 				fmt.Println("User account created successfully.")
// 			} else {
// 				fmt.Println("Error creating user account")
// 			}
// 		} else {
// 			fmt.Println(2, err)
// 		}
// 	} else {
// 		fmt.Println(3, err)
// 	}

// }

// func deleteAcc() {
// 	var accID int
// 	reader := bufio.NewReader(os.Stdin)
// 	reader.ReadString('\n')
// 	fmt.Print("Enter Account ID to delete: ")
// 	fmt.Scanf("%d", &accID)

// 	client := &http.Client{}
// 	if req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("http://localhost:5001/api/v1/accounts/delete?accID=%d", accID), nil); err == nil {
// 		if res, err := client.Do(req); err == nil {
// 			defer res.Body.Close()

// 			if res.StatusCode == http.StatusOK {
// 				fmt.Println("Account deleted successfully")
// 			} else {
// 				fmt.Println("Error deleting account")
// 			}
// 		} else {
// 			fmt.Println("Error making request:", err)
// 		}
// 	} else {
// 		fmt.Println("Error creating request:", err)
// 	}
// }

// func editAcc() {
// 	var accID int
// 	reader := bufio.NewReader(os.Stdin)
// 	reader.ReadString('\n')
// 	fmt.Print("Enter Account ID to edit: ")
// 	fmt.Scanf("%d", &accID)

// 	// Request updated information from the user
// 	var updatedAcc Account
// 	reader.ReadString('\n')
// 	fmt.Print("Enter updated Username: ")
// 	fmt.Scanf("%v", &updatedAcc.Username)
// 	reader.ReadString('\n')
// 	fmt.Print("Enter updated AccType (User, Admin): ")
// 	fmt.Scanf("%v", &updatedAcc.AccType)

// 	// Perform the update by making a PUT request to the API
// 	postBody, _ := json.Marshal(updatedAcc)
// 	client := &http.Client{}
// 	if req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("http://localhost:5001/api/v1/accounts/%d", accID), bytes.NewBuffer(postBody)); err == nil {
// 		if res, err := client.Do(req); err == nil {
// 			if res.StatusCode == http.StatusAccepted {
// 				fmt.Println("Profile updated successfully!")
// 			} else {
// 				fmt.Println("Error updating profile")
// 			}
// 		} else {
// 			fmt.Println("Error making request", err)
// 		}
// 	} else {
// 		fmt.Println("Error creating request", err)
// 	}
// }

// //admin capstone entry/records functions

// // list all capstone entries
// func listAllRecords() error {
// 	// Perform list all capstone entries
// 	client := &http.Client{}
// 	if req, err := http.NewRequest(http.MethodGet, "http://localhost:5002/api/v1/records/all", nil); err == nil {
// 		if res, err := client.Do(req); err == nil {
// 			defer res.Body.Close()
// 			if res.StatusCode == http.StatusOK {
// 				var records []Record
// 				err := json.NewDecoder(res.Body).Decode(&records)
// 				if err == nil {
// 					fmt.Println("List of all Capstone Entries:")
// 					for _, record := range records {
// 						fmt.Printf("\nRecord ID: %d \nName: %s \nRole of Contact: %s \nNo of Students: %d \nAcadamic Year: %s\nCapstone Title: %s \nCompany Name: %s \nCompany Contact: %s \nProject Desc: %s\n\n", record.RecordID, record.Name, record.RoleOfContact, record.NoOfStudents, record.AcadYr, record.CapstoneTitle, record.CompanyName, record.CompanyContact, record.ProjDesc)
// 					}
// 					return nil
// 				} else {
// 					return fmt.Errorf("Error decoding response: %v", err)
// 				}
// 			} else {
// 				return fmt.Errorf("Error fetching user list")
// 			}
// 		} else {
// 			return fmt.Errorf("Error making request: %v", err)
// 		}
// 	} else {
// 		return fmt.Errorf("Error creating request: %v", err)
// 	}
// }

// // manage record menu
// func manageRecordsMenu() {
// 	var choice int
// 	reader := bufio.NewReader(os.Stdin)
// 	fmt.Println("\nManage Records:")
// 	fmt.Println("1. Create Record")
// 	fmt.Println("2. Modify Record")
// 	fmt.Println("3. Delete Record")
// 	fmt.Println("4. Query Record by Academic Year or Keyword")
// 	fmt.Println("0. Go Back")
// 	reader.ReadString('\n')
// 	fmt.Print("Enter an option: ")
// 	fmt.Scanf("%d", &choice)

// 	switch choice {
// 	case 1:
// 		//post new record
// 		fmt.Println("----Create Record----")
// 		createRecord()
// 	case 2:
// 		//get recordID and allow modifications to acc details based on selected recordID
// 		fmt.Println("----Modify Record----")
// 		editRecord()
// 	case 3:
// 		//get recordID and delete selected account
// 		fmt.Println("----Delete Record----")
// 		deleteRecord()
// 	case 4:
// 		//search for record by academic year or keyword/capstone title
// 		fmt.Println("----Query Record----")
// 		queryRecord()
// 	case 0:
// 		// Go back
// 	default:
// 		fmt.Println("Invalid option")
// 	}
// }

// // start of all capstone related features (both admin and users share common record functions)
// // create new capstone entry
// func createRecord() {
// 	var record Record
// 	reader := bufio.NewReader(os.Stdin)
// 	reader.ReadString('\n')
// 	fmt.Print("Enter Name: ")
// 	fmt.Scanf("%v", &record.Name)
// 	reader.ReadString('\n')
// 	fmt.Print("Enter Role of Contact (Staff or Student): ")
// 	fmt.Scanf("%v", &record.RoleOfContact)
// 	reader.ReadString('\n')
// 	fmt.Print("Enter Number of Students: ")
// 	fmt.Scanf("%d", &record.NoOfStudents)
// 	reader.ReadString('\n')
// 	fmt.Print("Enter Academic Year: ")
// 	fmt.Scanf("%v", &record.AcadYr)
// 	reader.ReadString('\n')
// 	fmt.Print("Enter Capstone Title: ")
// 	fmt.Scanf("%v", &record.CapstoneTitle)
// 	reader.ReadString('\n')
// 	fmt.Print("Enter Name of Company: ")
// 	fmt.Scanf("%v", &record.CompanyName)
// 	reader.ReadString('\n')
// 	fmt.Print("Enter Company Point of Contact: ")
// 	fmt.Scanf("%v", &record.CompanyContact)
// 	reader.ReadString('\n')
// 	fmt.Print("Enter Brief Description of the Project: ")
// 	fmt.Scanf("%v", &record.ProjDesc)

// 	postBody, _ := json.Marshal(record)

// 	client := &http.Client{}
// 	if req, err := http.NewRequest(http.MethodPost, "http://localhost:5002/api/v1/records", bytes.NewBuffer(postBody)); err == nil {
// 		if res, err := client.Do(req); err == nil {
// 			if res.StatusCode == 201 {
// 				fmt.Println("Capstone Entry created successfully.")
// 			} else {
// 				fmt.Println("Error creating new entry")
// 			}
// 		} else {
// 			fmt.Println(2, err)
// 		}
// 	} else {
// 		fmt.Println(3, err)
// 	}

// }

// // delete capstone record
// func deleteRecord() {
// 	var recordID int
// 	reader := bufio.NewReader(os.Stdin)
// 	reader.ReadString('\n')
// 	fmt.Print("Enter Record ID to delete: ")
// 	fmt.Scanf("%d", &recordID)

// 	client := &http.Client{}
// 	if req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("http://localhost:5002/api/v1/records/delete?recordID=%d", recordID), nil); err == nil {
// 		if res, err := client.Do(req); err == nil {
// 			defer res.Body.Close()

// 			if res.StatusCode == http.StatusOK {
// 				fmt.Println("Record deleted successfully")
// 			} else {
// 				fmt.Println("Error deleting capstone entry")
// 			}
// 		} else {
// 			fmt.Println("Error making request:", err)
// 		}
// 	} else {
// 		fmt.Println("Error creating request:", err)
// 	}
// }

// // edit capstone entry
// func editRecord() {
// 	var recordID int
// 	reader := bufio.NewReader(os.Stdin)
// 	reader.ReadString('\n')
// 	fmt.Print("Enter Record ID to edit: ")
// 	fmt.Scanf("%d", &recordID)

// 	// Request updated information from the user
// 	var updatedrecord Record
// 	reader.ReadString('\n')
// 	fmt.Print("Enter Updated Name: ")
// 	fmt.Scanf("%v", &updatedrecord.Name)
// 	reader.ReadString('\n')
// 	fmt.Print("Enter Updated Role of Contact (Staff or Student): ")
// 	fmt.Scanf("%v", &updatedrecord.RoleOfContact)
// 	reader.ReadString('\n')
// 	fmt.Print("Enter Updated Number of Students: ")
// 	fmt.Scanf("%d", &updatedrecord.NoOfStudents)
// 	reader.ReadString('\n')
// 	fmt.Print("Enter Updated Academic Year: ")
// 	fmt.Scanf("%v", &updatedrecord.AcadYr)
// 	reader.ReadString('\n')
// 	fmt.Print("Enter Updated Capstone Title: ")
// 	fmt.Scanf("%v", &updatedrecord.CapstoneTitle)
// 	reader.ReadString('\n')
// 	fmt.Print("Enter Updated Name of Company: ")
// 	fmt.Scanf("%v", &updatedrecord.CompanyName)
// 	reader.ReadString('\n')
// 	fmt.Print("Enter Updated Company Point of Contact: ")
// 	fmt.Scanf("%v", &updatedrecord.CompanyContact)
// 	reader.ReadString('\n')
// 	fmt.Print("Enter Updated Brief Description of the Project: ")
// 	fmt.Scanf("%v", &updatedrecord.ProjDesc)

// 	// Perform the update by making a PUT request to the API
// 	postBody, _ := json.Marshal(updatedrecord)
// 	client := &http.Client{}
// 	if req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("http://localhost:5002/api/v1/records/%d", recordID), bytes.NewBuffer(postBody)); err == nil {
// 		if res, err := client.Do(req); err == nil {
// 			if res.StatusCode == http.StatusAccepted {
// 				fmt.Println("Capstone entry updated successfully!")
// 			} else {
// 				fmt.Println("Error updating capstone entry")
// 			}
// 		} else {
// 			fmt.Println("Error making request", err)
// 		}
// 	} else {
// 		fmt.Println("Error creating request", err)
// 	}
// }

// // search for capstone entry based on academic year or keyword/capstone title
// func queryRecord() {
// 	var query string
// 	reader := bufio.NewReader(os.Stdin)
// 	reader.ReadString('\n')
// 	fmt.Print("Enter Query (Year or Capstone Title): ")
// 	fmt.Scanf("%s", &query)

// 	client := &http.Client{}
// 	if req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://localhost:5002/api/v1/records/search?query=%s", query), nil); err == nil {
// 		if res, err := client.Do(req); err == nil {
// 			defer res.Body.Close()

// 			if res.StatusCode == http.StatusOK {
// 				var searchResults []Record
// 				err := json.NewDecoder(res.Body).Decode(&searchResults)
// 				if err == nil {
// 					fmt.Println("Search Results:")
// 					for _, record := range searchResults {
// 						fmt.Printf("Record ID: %d \nName: %s \nRole of Contact: %s \nNo of Students: %d \nAcadamic Year: %s\nCapstone Title: %s \nCompany Name: %s \nCompany Contact: %s \nProject Desc: %s\n\n", record.RecordID, record.Name, record.RoleOfContact, record.NoOfStudents, record.AcadYr, record.CapstoneTitle, record.CompanyName, record.CompanyContact, record.ProjDesc)
// 					}
// 				} else {
// 					fmt.Println("Error decoding response:", err)
// 				}
// 			} else {
// 				fmt.Println("No Records found")
// 			}
// 		} else {
// 			fmt.Println("Error making request:", err)
// 		}
// 	} else {
// 		fmt.Println("Error creating request:", err)
// 	}
// }
