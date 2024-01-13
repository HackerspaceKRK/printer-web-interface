package main

import (
	"embed"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

//go:embed src/*
var content embed.FS

func main() {
	// Set up HTTP server\
	subfs, err := fs.Sub(content, "src")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	http.HandleFunc("/", http.FileServer(http.FS(subfs)).ServeHTTP)
	http.HandleFunc("/print", printHandler)

	// Start the server
	port := 80
	fmt.Printf("Server is running on http://localhost:%d\n", port)
	err = http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		fmt.Println("Error:", err)
	}
}

type Data struct {
	Png string `json:"png"`
}

func printHandler(w http.ResponseWriter, r *http.Request) {

	// parse json body
	data := Data{}
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, "Unable to parse json body", http.StatusBadRequest)
		return
	}

	base64Str := strings.TrimPrefix(data.Png, "data:image/png;base64,")

	dataDecoded, err := base64.StdEncoding.DecodeString(base64Str)
	if err != nil {
		http.Error(w, "Unable to decode base64 string", http.StatusBadRequest)
		log.Println(err)
		return
	}

	// Save the file to a temporary directory
	tempDir := os.TempDir()
	tempFile, err := os.CreateTemp(tempDir, "uploaded-image-*.png")
	if err != nil {
		http.Error(w, "Unable to create temporary file", http.StatusInternalServerError)
		return
	}

	_, err = tempFile.Write(dataDecoded)
	if err != nil {
		http.Error(w, "Unable to save file", http.StatusInternalServerError)
		return
	}

	tempFile.Close()

	log.Println("File saved to", tempFile.Name())

	// Execute the lp command
	printCommand := exec.Command("lp", "-d", "CITIZEN_CL-S700Z", tempFile.Name())
	err = printCommand.Run()
	if err != nil {
		log.Println(printCommand.Stderr)
		log.Println(printCommand.Stdout)
		log.Println(err)
		http.Error(w, "Unable to execute print command", http.StatusInternalServerError)
		return
	}

	// Delete the temporary file
	err = os.Remove(tempFile.Name())
	if err != nil {
		log.Println(err)
	}

	// Respond with success message
	w.Write([]byte("File printed successfully"))
}
