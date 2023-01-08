package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/go-resty/resty/v2"
)

var URL = flag.String("url", "http://localhost:80", "URL to the server")
var APIKEY = flag.String("apikey", "", "API Key for the server")

func main() {
	flag.Parse()
	filesToSubmit := readFileNamesStdIn()
	for _, filename := range filesToSubmit {
		fileBytes, err := readFile(filename)
		if err != nil {
			fmt.Printf("Error reading file %s: %s", filename, err)
			continue
		}
		fmt.Printf("[+] Submitting file %s\n", filename)
		postFileToServer(filename, fileBytes)
	}
}

func readFileNamesStdIn() []string {
	fileBytes, err := io.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}
	fileNames := bytes.Split(fileBytes, []byte{'\n'})
	// Remove empty file names
	var filesToSubmit []string
	for _, fileName := range fileNames {
		if len(fileName) > 0 {
			filesToSubmit = append(filesToSubmit, string(fileName))
		}
	}
	return filesToSubmit
}

func readFile(file string) ([]byte, error) {
	filePointer, err := os.Open(file)
	if err != nil {
		return []byte{}, err
	}
	defer filePointer.Close()
	fileBytes, err := io.ReadAll(filePointer)
	if err != nil {
		return []byte{}, err
	}
	return fileBytes, nil
}

func postFileToServer(filename_full_path string, fileBytes []byte) {
	filename_split := strings.Split(filename_full_path, "/")
	filename := string(filename_split[len(filename_split)-1])
	if strings.HasSuffix(*URL, "/") {
		*URL = strings.TrimSuffix(*URL, "/")
	}
	full_url := fmt.Sprintf("%s/api/file", *URL)
	client := newClient()
	resp, err := client.R().SetFileReader("file", filename, bytes.NewReader(fileBytes)).Post(full_url)
	if err != nil {
		fmt.Printf("Error submitting file %s: %s", filename, err)
		return
	}
	if resp.StatusCode() != 200 {
		fmt.Printf("Error submitting file %s: %s", filename, resp.String())
		return
	}
	fmt.Printf("[+] File %s submitted successfully\n", filename)

}

func newClient() *resty.Client {
	client := resty.New()
	client.SetHeader("Content-Type", "application/json")
	client.SetHeader("Authorization", fmt.Sprintf("Bearer %s", *APIKEY))
	return client
}
