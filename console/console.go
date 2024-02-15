//developing using console first before changing to html

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
