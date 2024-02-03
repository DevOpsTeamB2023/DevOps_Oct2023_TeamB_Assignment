//developing using console first before changing to html

package main

import (
	"CICD_TEST/microservices/account" //change here
	"CICD_TEST/microservices/record"  //change here
)

func main() {
	account.InitHTTPServer()
	record.InitHTTPServer()
}
