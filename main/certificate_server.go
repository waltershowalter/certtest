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

func checkSameCertRequestInSeconds(passedDomainString string) {

}

func newDomainCert(passedDomainString string) *certData{
	return &certData{
		createdTime: time.Now().UTC(),
		domain: passedDomainString,
		expired: false,
	}
}

func checkURL(passedDomainString string) {


}

func addDomainCert(passedDomainString string, w http.ResponseWriter) {

	newStruct := newDomainCert(passedDomainString)

	_, exists := certList[passedDomainString]
	if exists {
		fmt.Println("Certificate already exists!!!")
		fmt.Fprintf(w, "Certificate already exists!!!\n")
		fmt.Println("\tChecking if it's expired")
		fmt.Fprintf(w, "Checking if it's expired\n")
		fmt.Println("\tStatus: " + strconv.FormatBool(certList[passedDomainString].expired))
		if certList[passedDomainString].expired == true {
			fmt.Println("\tRecertifying certificate....")
			fmt.Fprintf(w, "Recertifying certificate....\n")
			certList[passedDomainString].createdTime = time.Now().UTC()
			certList[passedDomainString].expired = false
			fmt.Fprintf(w, "%s has been recertified", passedDomainString)
		} else {
			fmt.Fprintf(w, "%s is still valid. No need to be recertified", passedDomainString)
		}
	} else {
		certList[passedDomainString] = newStruct
		fmt.Fprintf(w, "%s certificate has been created", passedDomainString)
	}


}

func serveCertificate(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case "GET":
		requestURL := r.URL.EscapedPath()
		path := r.URL.Path
		fmt.Println("URL,  Path: %s, %s \n", requestURL, path)
		subPaths := strings.Split(r.URL.Path, "/")
		if len(subPaths) > 3 {
			fmt.Fprintf(w, "Length of URI subpath is too long")
		}
		fmt.Println("Browken up path:")
		for _, currentString := range subPaths {
			fmt.Println("Subpath : %s", currentString)
		}
		addDomainCert(subPaths[len(subPaths)-1], w)
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
