package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
)

var tracks []string

func readMusic() []string {
	path := "./music"
	files, err := os.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}
	var fileNames []string
	for _, file := range files {
		if !file.IsDir() {
			fileNames = append(fileNames, file.Name())
		}
	}
	return fileNames
}

func musicHandler(w http.ResponseWriter, r *http.Request) {

	indexStr := r.URL.Query().Get("index")
	if indexStr == "" {
		http.Error(w, "Index parameter is missing", http.StatusBadRequest)
		return
	}

	index, err := strconv.Atoi(indexStr)
	if err != nil {
		http.Error(w, "Invalid index parameter", http.StatusBadRequest)
		return
	}

	if index < 0 || index >= len(tracks) {
		http.Error(w, "Index out of range", http.StatusBadRequest)
		return
	}

	filename := tracks[index]
	file, err := os.Open("music/" + filename)
	if err != nil {
		http.Error(w, "Could not open music file", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		http.Error(w, "Could not get file info", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "audio/mpeg")
	w.Header().Set("Content-Length", fmt.Sprint(fileInfo.Size()))

	http.ServeContent(w, r, filename, fileInfo.ModTime(), file)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {

	tmpl, err := template.ParseFiles("index.html")
	if err != nil {
		http.Error(w, "Could not read HTML template", http.StatusInternalServerError)
		return
	}

	trackNames := make([]string, len(tracks))
	for i, track := range tracks {
		trackNames[i] = track
	}

	data := struct {
		TrackCount int
		Tracks     []string
	}{
		TrackCount: len(tracks),
		Tracks:     trackNames,
	}

	w.Header().Set("Content-Type", "text/html")
	tmpl.Execute(w, data)
}

func main() {
	tracks = readMusic()

	http.HandleFunc("/music", musicHandler)
	http.HandleFunc("/", indexHandler)

	fmt.Println("Server is running on http://127.0.0.1:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Printf("Ошибка старта сервера: %s\n", err)
	}
}
