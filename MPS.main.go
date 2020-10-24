package main

import (
	"database/sql"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

var db, _ = sql.Open("sqlite3", "songdb.sqlite")
var songlistSet = false
var songlist []string

// https://jaxenter.de/restful-rest-api-go-golang-68845
func main() {
	//storetest()
	rebuildFileList()

	http.HandleFunc("/", netSendRandomSong)
	http.HandleFunc("/init", netRebuildFileList)
	http.ListenAndServe(":80", nil)
}

// https://stackoverflow.com/questions/24116147/how-to-download-file-in-browser-from-go-serverv
// https://developer.mozilla.org/en-US/docs/Web/HTTP/Basics_of_HTTP/MIME_types/Common_types

func netSendRandomSong(w http.ResponseWriter, r *http.Request) {
	// Get filename from DB
	var filePath string = getNextSongName()

	var fileContent, err = os.Open(filePath)
	if err != nil {
		log.Print(err)
	}
	var stat, _ = fileContent.Stat()
	var size = stat.Size()
	// Für den Download als datei
	w.Header().Set("Content-Disposition", "attachment; filename="+filepath.Base(filePath))
	//w.Header().Set("Content-Disposition", "inline")
	w.Header().Set("Content-Type", "audio/mpeg")
	w.Header().Set("Content-Length", strconv.FormatInt(size, 10))
	io.Copy(w, fileContent)
}

func getNextSongName() string {
	//return "E:\\Musik\\1 für beat hazart\\Apply (Felix Green remix).mp3"

	rand.Seed(time.Now().UnixNano())
	if songlistSet {
		return songlist[rand.Intn(len(songlist))]
	}
	return ""
}

func netRebuildFileList(w http.ResponseWriter, r *http.Request) {
	//w.Write([]byte("Rebuilding of file list database started. This will take a while."))
	w.Write([]byte(rebuildFileList()))
}

func rebuildFileList() string {
	//var allFilesArr []string
	var allFilesStr string
	songlistSet = false
	songlist = []string{}

	filepath.Walk("E:\\Musik", //".",
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if filepath.Ext(path) == ".mp3" {
				allFilesStr += path + "\n"
				songlist = append(songlist, path)
			}
			return nil
		})
	songlistSet = true
	return allFilesStr
}
