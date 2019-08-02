**Certificate Domain Service**

Tested on: MAC OS X, High Sierra

Go version used :  go1.10.1 darwin/amd64

To run type: go run certificate_server.go in main under the root folder

This certtest program returns a mock cert based off the the URL you give in the GET request

To use open a browser and type : http://localhost:8888/cert/www.abc.com

You can swap out the URL at the end of the path with your own test URL.

Certificates will time out after 10 minutes and be marked as expired until a new request comes in. At that time the already created certificate will be marked as non expired and a new time stamp will be added to allow for another 10 minute valid certificate period.

NOTES based off assignment:

Implement a code base that:

1) properly handles certificates for different domains - DONE
2) allows certificates to live for 10 minutes before expiring - DONE
sleeps for 10 seconds when a new certificate is generated - NOT DONE
generates its own certificate and keeps its certificate up-to-date when it expires - DONE
can generate multiple certificates at the same time - DONE
