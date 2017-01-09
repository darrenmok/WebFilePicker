package main

// first try of golang
import (
	"net/http"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"io"
)

// html template location of all the html files
var templates = template.Must(template.ParseGlob("template/*.html"))

// simple struct for storing slices of filenames
type FileList struct {
	FileNames []string
}

// to render the view page for uploading form
func viewHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "view", nil)
}

// render page to list all files
func browserHandler(w http.ResponseWriter, r *http.Request) {
	files, err := ioutil.ReadDir("uploaded")
	if err != nil {
		renderTemplate(w, err.Error(), nil)
		return
	}

	fileNames := make([]string, len(files))

	log.Println("Length of files = ", len(files))

	for i := 0; i < len(files); i++ {
		log.Println("Loop = ", i)
		fileNames[i] = files[i].Name()
		log.Println("value = ", fileNames[i])
	}

	//for _, file := range files {
	//	fileNames = append(fileNames, file.Name())
	//	log.Println("Filenames " + file.Name())
	//}
	data := &FileList{FileNames:fileNames}

	templates.ExecuteTemplate(w, "browser.html", data)
}

// handler to handle file upload
func uploadHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Enter upload handler")
	fileUpload, fileHeader, err := r.FormFile("file")

	buff := make([]byte, 512)
	_, err = fileUpload.Read(buff)
	if err != nil {
		log.Println("buffer error")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer fileUpload.Close()

	contentType := http.DetectContentType(buff)
	log.Println("Content Type = " + contentType)

	if err != nil {
		log.Println("FormFile error")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Println("Uploading file " + fileHeader.Filename)

	fileTemp, err := os.Create("uploaded/" + fileHeader.Filename)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer fileTemp.Close()

	_, err = io.Copy(fileTemp, fileUpload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	browserHandler(w, r)
}

// abstract function to execute the template to render the html
func renderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	err := templates.ExecuteTemplate(w, tmpl + ".html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// the app entry point
func main() {
	http.HandleFunc("/view/", viewHandler)
	http.HandleFunc("/browser/", browserHandler)
	http.HandleFunc("/upload/", uploadHandler)

	http.ListenAndServe(":8080", nil)
}
