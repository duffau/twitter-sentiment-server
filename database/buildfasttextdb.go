package main

import (
	_ "github.com/mattn/go-sqlite3"
	"github.com/ekzhu/go-fasttext"
	"os"
)

func main() {
	ft := fasttext.NewFastText("../database/fasttext.db")
	vec_file, _ := os.Open("../wiki.en.vec")
	ft.BuildDB(vec_file)	
	ft.Close()
}


