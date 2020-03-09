package main

import (
	"io"
	"log"
	"net/http"
	"os"
)

// https://jaxenter.de/restful-rest-api-go-golang-68845
func main() {
	http.HandleFunc("/", productsHandler)
	http.ListenAndServe(":8080", nil)
}

// https://stackoverflow.com/questions/24116147/how-to-download-file-in-browser-from-go-serverv
// https://developer.mozilla.org/en-US/docs/Web/HTTP/Basics_of_HTTP/MIME_types/Common_types
func productsHandler(w http.ResponseWriter, r *http.Request) {
	//w.Write([]byte("Schuhe, Hose, Hemd"))
	var filePath string = "E:\\Musik\\1 für beat hazart\\Apply (Felix Green remix).mp3"
	var fileContent, err = os.Open(filePath)
	if err != nil {
		log.Print(err)
	}
	w.Header().Set("Content-Disposition", "attachment; filename=music.mp3")
	w.Header().Set("Content-Type", "audio/mpeg")
	io.Copy(w, fileContent)
}
