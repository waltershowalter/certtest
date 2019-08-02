package backup

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"
	"strings"
)

var counter int
var mutex = &sync.Mutex{}

var domainName string

var ONE_MINUTE = 1 * time.Second.Minutes()

type certData struct {

	sync.Mutex
	createdTime time.Time
	domain string
	expired bool

}

var certList = make(map[string]*certData)

func checkExpire() {

	fmt.Println("Launching checkExpire")

	for {
		// do some job
		time.Sleep(1 * time.Second)
		timeNow := time.Now().UTC()
		minusOneMinute := timeNow.Add(time.Duration(20) * time.Second)
		fmt.Println("Minus twenty seconds : %s", minusOneMinute)

		for thisKey, thisStruct := range certList {
			// (*certData).Lock()
			fmt.Println("Testing expiry : %s", thisKey)
			if thisStruct.createdTime.Before(minusOneMinute) {
				fmt.Println("Cert expired for: %s", thisKey)
				thisStruct.expired = false
			}
			// (*certData).Unlock()
		}

	}
}

func echoString(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w,"Req: %s %s\n", r.Host, r.URL.Path)
	fmt.Fprintf(w, "hello")
}

func incrementCounter(w http.ResponseWriter, r *http.Request) {
	mutex.Lock()
	counter++
	fmt.Fprintf(w,"Req: %s %s\n", r.Host, r.URL.Path)
	fmt.Fprintf(w, strconv.Itoa(counter))
	mutex.Unlock()
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
		fmt.Println("cert already exists!!!")
	} else {
		certList[passedDomainString] = newStruct
	}


}

func serveCertificate(w http.ResponseWriter, r *http.Request) {
	//if r.URL.Path != "/cert/" {
	//	http.Error(w, "URI not found. Possible 404", http.StatusNotFound)
	//	return
	//}

	switch r.Method {
	case "GET":
		requestURL := r.URL.EscapedPath()
		requestURI := r.URL.RequestURI()
		host := r.URL.Host
		hostname := r.URL.Hostname()
		path := r.URL.Path
		fmt.Fprintf(w,"URL, Req, Host, Hostname, Path: %s, %s, %s, %s, %s \n", requestURL, requestURI, host, hostname, path)
		subPaths := strings.Split(r.URL.Path, "/")
		fmt.Fprintf(w, "Length : %d\n", len(subPaths))
		if len(subPaths) > 3 {
			fmt.Fprintf(w, "Length of URI subpath is too long")
		}
		for _, currentString := range subPaths {
			fmt.Fprintf(w, currentString + "--\n")
		}
		addDomainCert(subPaths[len(subPaths)-1])
	default:
		fmt.Fprintf(w, "Sorry, only GET method is supported.")
	}
}


func hello(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	switch r.Method {
	case "GET":
		http.ServeFile(w, r, "form.html")
	case "POST":
		// Call ParseForm() to parse the raw query and update r.PostForm and r.Form.
		if err := r.ParseForm(); err != nil {
			fmt.Fprintf(w, "ParseForm() err: %v", err)
			return
		}
		fmt.Fprintf(w, "Post from website! r.PostFrom = %v\n", r.PostForm)
		name := r.FormValue("name")
		address := r.FormValue("address")
		fmt.Fprintf(w, "Name = %s\n", name)
		fmt.Fprintf(w, "Address = %s\n", address)
	default:
		fmt.Fprintf(w, "Sorry, only GET and POST methods are supported.")
	}
}


func firstTest() {

	http.HandleFunc("/", echoString)

	http.HandleFunc("/increment", incrementCounter)

	http.HandleFunc("/hi", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hi")
	})

	log.Fatal(http.ListenAndServe(":8081", nil))

}

func secondTest() {

	http.HandleFunc("/", hello)

	fmt.Printf("Starting server for testing HTTP REST requests...\n")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}


}


func certTest() {


	go checkExpire()

	http.HandleFunc("/cert/", serveCertificate)

	fmt.Printf("Starting server for testing HTTP GET...\n")
	if err := http.ListenAndServe(":8888", nil); err != nil {
		log.Fatal(err)
	}




}

func main() {

	//firstTest()

	//secondTest()

	certTest()

}
