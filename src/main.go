package main

import (
	"net/http"
	"log"
	"io/ioutil"
	"strings"
	//"github.com/dghubble/go-twitter/twitter"
	//"github.com/dghubble/oauth1"
)


func myHandlerFunc(w http.ResponseWriter, req *http.Request) {
	path := string(req.URL.Path[1:])
	log.Println(path)
	data, err := ioutil.ReadFile(path)
	addContentType(w, path)
	if err == nil {
		w.Write(data)
	} else {
		w.WriteHeader(404)
		w.Write([]byte("404 - " + http.StatusText(404)))
	}
}


func addContentType(w http.ResponseWriter, path string){
	var contentType string

	if strings.HasSuffix(path, ".html") {
        contentType = "text/html"
	} else if strings.HasSuffix(path, ".js") {
		contentType = "application/javascript"
	} else {
		contentType = "text/plain"			
	}
	w.Header().Add("Content Type", contentType)
}

func main() {
	http.HandleFunc("/", myHandlerFunc)
	http.ListenAndServe(":8080", nil)
}
