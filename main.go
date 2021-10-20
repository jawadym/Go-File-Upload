package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func setupRoutes() {
	// http handler
	http.HandleFunc("/upload", uploadFile)
	http.HandleFunc("/temp-files/", temporaryServe)

	// static files
	static_files := http.FileServer(http.Dir("./static"))
	http.Handle("/", static_files)
	http.ListenAndServe(":8080", nil)
}

func temporaryServe(w http.ResponseWriter, r *http.Request) {

	fileBytes, err := ioutil.ReadFile(fmt.Sprintf(".%s", r.URL))
	if err != nil {
		panic(err)
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(fileBytes)
	return
}

func uploadFile(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Uploading File")

	// parse input of type multipart/form-data
	// max 10MB
	r.ParseMultipartForm(10 << 20)

	// get file
	file, handler, err := r.FormFile("myFile")
	if err != nil {
		fmt.Println("Error getting file from form data")
		return
	}
	defer file.Close()
	fmt.Printf("Uploaded file: %+v\n", handler.Filename)
	fmt.Printf("Uploaded size: %+v\n", handler.Size)
	fmt.Printf("MIME TYPE: %+v\n", handler.Header)

	// write temporary file
	tempfile, err := ioutil.TempFile("temp-files", "uploaded-*.png")
	if err != nil {
		fmt.Println("Error saving file")
		return
	}
	defer tempfile.Close()

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println("Error reading uploaded file")
	}

	tempfile.Write(fileBytes)

	// return status
	fmt.Fprintf(w, "\nSuccessfully uploaded! Available at %s", tempfile.Name())
}

func main() {
	fmt.Println("Go file upload")
	setupRoutes()
}
