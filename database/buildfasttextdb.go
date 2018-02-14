package main

import (
	_ "github.com/mattn/go-sqlite3"
	"github.com/ekzhu/go-fasttext"
	"os"
)

func main() {
	ft := fasttext.NewFastText("glove25.db")
	vec_file, _ := os.Open("../data/glove.twitter.27B.25d.txt")
	ft.BuildDB(vec_file)	
	ft.Close()
}


