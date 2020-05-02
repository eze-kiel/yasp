package handlers

import (
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/eze-kiel/yasp/utils"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

const maxUploadSize = 5 * 1024 * 1024 // 5 Mo

// Transaction contains informations about the upload
type Transaction struct {
	Success bool
	ID      string
}

// HandleFunc handle funcs
func HandleFunc() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/", homePage)
	r.HandleFunc("/upload", uploadPage).Methods("GET")
	r.HandleFunc("/upload", uploadData).Methods("POST")
	r.HandleFunc("/download", downloadPage).Methods("GET")
	r.HandleFunc("/download", downloadData).Methods("POST")

	return r
}

func homePage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("views/home.html", "views/templates/navbar.html", "views/templates/footer.html")
	if err != nil {
		log.Fatal(err)
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		log.Fatalf("Can not execute templates for home page : %v", err)
	}

}

func uploadPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("views/upload.html", "views/templates/navbar.html", "views/templates/footer.html")
	if err != nil {
		log.Fatal(err)
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		log.Fatalf("Can not execute templates for upload page : %v", err)
	}
}

func downloadPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("views/download.html", "views/templates/navbar.html", "views/templates/footer.html")
	if err != nil {
		log.Fatal(err)
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		log.Fatalf("Can not execute templates for donwload page : %v", err)
	}
}

func uploadData(w http.ResponseWriter, r *http.Request) {
	// Try to parse data from post form
	if err := r.ParseMultipartForm(maxUploadSize); err != nil {
		fmt.Printf("Could not parse multipart form: %v\n", err)
		return
	}

	// Get the file from form
	file, fileHeader, err := r.FormFile("filename")
	if err != nil {
		return
	}
	defer file.Close()

	fileSize := fileHeader.Size

	// Check if the file's size is accepted
	if fileSize > maxUploadSize {
		return
	}

	// implement io.Reader
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		return
	}

	detectedFileType := http.DetectContentType(fileBytes)

	// Create a new name based on a random token
	// should use UUID in production
	fileName := utils.RandToken(12)

	fileEndings, err := mime.ExtensionsByType(detectedFileType)
	if err != nil {
		return
	}

	path := "./uploads/" + fileName[0:2] + "/" + fileName[2:4] + "/"
	err = os.MkdirAll(path, 0700)
	if err != nil && !os.IsExist(err) {
		return
	}

	newPath := filepath.Join(path, fileName+fileEndings[0])
	fmt.Printf("FileType: %s, File: %s\n", detectedFileType, newPath)

	// write file
	newFile, err := os.Create(newPath)
	if err != nil {
		return
	}
	defer newFile.Close()

	if _, err := newFile.Write(fileBytes); err != nil || newFile.Close() != nil {
		return
	}

	uploadState := Transaction{
		Success: true,
		ID:      fileName,
	}

	// Parse templates to display file's id
	tmpl, err := template.ParseFiles("views/upload.html", "views/templates/navbar.html", "views/templates/footer.html")
	if err != nil {
		log.Fatal(err)
	}

	err = tmpl.Execute(w, uploadState)
	if err != nil {
		log.Fatalf("Can not execute templates for donwload page : %v", err)
	}

}

func downloadData(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("fileID")

	if len(id) < 12 {
		log.Printf("Error in id length : %d instead of 12\n", len(id))
		return
	}

	path := "./uploads/" + id[0:2] + "/" + id[2:4] + "/"

	// check if the file exists
	fileName := utils.FindRealFilename(path, id)

	if fileName == "" {
		log.Printf("Filename is %s\n", fileName)
		http.Error(w, "File not found.", 404)
		return
	}

	openfile, err := os.Open(fileName)
	if err != nil {
		logrus.Errorf("error opening file %s: %v", fileName, err)
		http.Error(w, "File not found.", 404)
		return
	}
	defer openfile.Close()

	// Read the 512 first bytes of the file's headers
	FileHeader := make([]byte, 512)
	openfile.Read(FileHeader)

	FileContentType := http.DetectContentType(FileHeader)

	// Get informations about to file for the headers
	FileStat, _ := openfile.Stat()
	FileSize := strconv.FormatInt(FileStat.Size(), 10)

	parts := strings.Split(fileName, "/")
	// Send the headers
	w.Header().Set("Content-Disposition", "attachment; filename="+parts[len(parts)-1])
	w.Header().Set("Content-Type", FileContentType)
	w.Header().Set("Content-Length", FileSize)

	// Send the file
	openfile.Seek(0, 0)
	io.Copy(w, openfile)
}
