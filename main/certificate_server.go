// Author : Andrew Fernandez
// Assignment : For Fanatics
// Date : July 28th, 2019

package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"
	"strings"
	"testing"

)

var (

	domainName string

	TEN_MINUTES = 10 * time.Second.Minutes()
	TEN_SECONDS = 10 * time.Second.Seconds()

)

type certData struct {

	mux sync.Mutex
	createdTime time.Time
	domain string
	expired bool

}

var certList = make(map[string]*certData)

func checkExpire() {

	fmt.Println("Launching checkExpire")

	for {
		time.Sleep(1 * time.Second)
		timeNow := time.Now().UTC()
		minusTenMinutes := timeNow.Add(time.Duration(-30) * time.Second)

		for thisKey, thisStruct := range certList {
			thisStruct.mux.Lock()
			fmt.Println("Testing expiry : %s", thisKey)
			if thisStruct.createdTime.Before(minusTenMinutes) {
				fmt.Println("\tCert expired for: %s", thisKey)
				thisStruct.expired = true
			}
			thisStruct.mux.Unlock()
		}

	}
}

func newDomainCert(passedDomainString string) *certData{
	return &certData{
		createdTime: time.Now().UTC(),
		domain: passedDomainString,
		expired: false,
	}
}

func addDomainCert(passedDomainString string) {

	newStruct := newDomainCert(passedDomainString)

	_, exists := certList[passedDomainString]
	if exists {
		fmt.Println("Certificate already exists!!!")
		fmt.Println("\tChecking if it's expired")
		fmt.Println("\tStatus: " + strconv.FormatBool(certList[passedDomainString].expired))
		if certList[passedDomainString].expired == true {
			fmt.Println("\tRecertifying certificate....")
			certList[passedDomainString].createdTime = time.Now().UTC()
			certList[passedDomainString].expired = false
		}
	} else {
		certList[passedDomainString] = newStruct
	}


}

func serveCertificate(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case "GET":
		requestURL := r.URL.EscapedPath()
		path := r.URL.Path
		fmt.Fprintf(w,"URL,  Path: %s, %s \n", requestURL, path)
		subPaths := strings.Split(r.URL.Path, "/")
		// fmt.Fprintf(w, "Length : %d\n", len(subPaths))
		if len(subPaths) > 3 {
			fmt.Fprintf(w, "Length of URI subpath is too long")
		}
		for _, currentString := range subPaths {
			fmt.Fprintf(w, currentString + "\n")
		}
		addDomainCert(subPaths[len(subPaths)-1])
	default:
		fmt.Fprintf(w, "Sorry, only GET method is supported.")
	}
}


func certTest() {

	go checkExpire()

	http.HandleFunc("/cert/", serveCertificate)

	fmt.Printf("Starting server for testing HTTP REST GET Cert calls...\n")
	if err := http.ListenAndServe(":8888", nil); err != nil {
		log.Fatal(err)
	}

}

func TestCertificateServerGetRequest(t *testing.T) {

	// TODO : called with go test
	// TODO : Test 100 concurrent requests

	certTest()

	resp, err := http.Get("http://localhost:8888/cert/www.cnn.com")
	if err != nil {
		fmt.Println("Error using GET to retrieve a certificate")
	}
	if resp.StatusCode != 200 {
		fmt.Println("Non-200 code returned")

	}

}

func main() {

	certTest()

}
