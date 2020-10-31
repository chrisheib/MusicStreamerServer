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

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

var db, _ = sql.Open("sqlite3", "songdb.sqlite")
var songlistSet = false
var songlist []string

const htmlBaseHeader = "<html>\n<body>\n"
const htmlBaseFooter = "</body>\n</html>"
const nl = "\n"
const hnl = "</br>" + nl

// https://jaxenter.de/restful-rest-api-go-golang-68845
func main() {
	//storetest()
	initDB()
	rebuildFileList()

	r := mux.NewRouter()

	r.HandleFunc("/", netSendBase)
	r.HandleFunc("/random", netSendRandomSong)
	r.HandleFunc("/randomID", netSendRandomID)
	r.HandleFunc("/init", netRebuildFileList)
	r.HandleFunc("/song/data/{id}", netSendSongDataByID)
	r.HandleFunc("/song/file/{id}", netSendSongByID)

	http.ListenAndServe(":80", r)
}

func initDB() {
	if selI("SELECT count(*) FROM sqlite_master WHERE type='table' AND name='songs';") < 1 {
		db.Exec("CREATE TABLE songs (id INTEGER not null primary key autoincrement, path TEXT unique, name TEXT, play INTEGER);")
	}
}

// https://stackoverflow.com/questions/24116147/how-to-download-file-in-browser-from-go-serverv
// https://developer.mozilla.org/en-US/docs/Web/HTTP/Basics_of_HTTP/MIME_types/Common_types
func netSendBase(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(htmlBaseHeader +
		link("random", "/random") + hnl +
		//link("init", "/init") + hnl +
		htmlBaseFooter))
}

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

func netSendRandomID(w http.ResponseWriter, r *http.Request) {
	id := selS("select id from songs ORDER BY RANDOM() LIMIT 1")
	w.Write([]byte(id))
}

func getNextSongName() string {
	//return "E:\\Musik\\1 für beat hazart\\Apply (Felix Green remix).mp3"

	rand.Seed(time.Now().UnixNano())
	if songlistSet {
		return songlist[rand.Intn(len(songlist))]
	}
	return ""
}

// https://golangcode.com/get-a-url-parameter-from-a-request/
func netRebuildFileList(w http.ResponseWriter, r *http.Request) {
	//w.Write([]byte("Rebuilding of file list database started. This will take a while."))
	w.Write([]byte(htmlBaseHeader + link("Zurück!", "/") + hnl + div(rebuildFileList()) + htmlBaseFooter))
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
				allFilesStr += path + "</br>\n"
				songlist = append(songlist, path)
				db.Exec("INSERT OR IGNORE INTO songs (path, name) VALUES (?,?)", path, filepath.Base(path))
			}
			return nil
		})
	songlistSet = true
	return allFilesStr
}

func netSendSongDataByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := vars["id"]
	name := selS("select name from songs where id = ?", id)
	w.Write([]byte(name))
}

func netSendSongByID(w http.ResponseWriter, r *http.Request) {
	//vars := mux.Vars(r)
	//id, _ := strconv.Atoi(vars["id"])
}

func link(text string, target string) string {
	return "<a href=\"" + target + "\">" + text + "</a>"
}

func div(content string) string {
	return "<div>" + nl + content + nl + "</div>" + nl
}

func selS(statement string, args ...interface{}) string {
	row := db.QueryRow(statement, args...)
	e := row.Err()
	if e != nil {
		return e.Error()
	}
	var value string
	row.Scan(&value)
	return value
}

func selI(statement string) int {
	row := db.QueryRow(statement)
	var value int
	row.Scan(&value)
	return value
}
