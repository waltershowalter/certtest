package constants

import (
	"sync"
	"time"
)

var (

	DomainName string

	TEN_MINUTES = 10 * time.Second.Minutes()
	TEN_SECONDS = 10 * time.Second.Seconds()

)

type CertData struct {

	mux sync.Mutex
	createdTime time.Time
	domain string
	expired bool

}

var CertList = make(map[string]*CertData)