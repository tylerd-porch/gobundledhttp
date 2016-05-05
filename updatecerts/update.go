// Updates certificates by downloading from https://curl.haxx.se/ca/cacert.pem
// Requires a machine with exisitng certificate pool in order to run
package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
)

var header []byte = []byte(fmt.Sprintf("package updatecerts\n\nvar pemCerts = []byte(`"))
var footer []byte = []byte(fmt.Sprintf("`)\n"))

func openCacerts() []byte {
	certFileHandle, err := ioutil.ReadFile("cacert.pem")
	if err != nil {
		log.Fatalf("Attempted to read expected file: cacert.pem \n")
	}
	return certFileHandle
}

func fetchCacerts() bool {
	filename := "cacert.pem"
	urlString := "https://curl.haxx.se/ca/cacert.pem"
	log.Println("Downloading cacert.pem file")
	// Using the basic http client because we can't load certs ourselves,
	// wget in busybox can't contact https sites either.
	downloadresp, _ := http.Get(urlString)
	caFileHandle, _ := os.Create(filename)
	defer caFileHandle.Close()
	_, err := io.Copy(caFileHandle, downloadresp.Body)
	if err != nil {
		log.Printf("HTTP error: %v \n", err)
		return false
	}
	return true
}

func main() {
	// Want to retrieve https://curl.haxx.se/ca/cacert.pem
	// Process it to be a full go-style file
	reg := regexp.MustCompile(`-----BEGIN CERTIFICATE-----[\n|\S]+-----END CERTIFICATE-----`)

	// if cacerts file does not exist, download it
	var certFile []byte
	if _, err := os.Stat("cacert.pem"); os.IsNotExist(err) {
		if ok := fetchCacerts(); ok {
			certFile = openCacerts()
		} else {
			log.Fatalf("Unable to download or open file cacert.pem \n")
		}
	} else {
		certFile = openCacerts()
	}
	matches := reg.FindAll(certFile, -1) // Find all matches

	certsgo := make([]byte, 0)
	//Byte slice to construct the go file must be zero cap so there are no leading null bytes in resulting file

	certsgo = append(certsgo, header...)
	for _, b := range matches {
		x := append([]byte{'\n'}, b...)
		certsgo = append(certsgo, x...)
	}
	certsgo = append(certsgo, footer...)

	// Write out the resulting go file
	fh, _ := os.Create("certificates.go")
	defer fh.Close()
	_, err := fh.Write(certsgo)
	if err != nil {
		log.Fatalf("Error writing certificates.go: %v \n", err)
	}
	log.Println("Finished writing file.")
	log.Println("MOVE certificates.go into the /certificates/ folder")
}
