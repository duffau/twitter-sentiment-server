 package main

import (
	"fmt"
	"encoding/csv"
	"os"
	_ "github.com/mattn/go-sqlite3"
	"github.com/ekzhu/go-fasttext"
	"bitbucket.com/chrduffau/wordvec"

)


type Tweetline struct {
	SentimentLabel string
	SentimentID int
	Text string
	TextVec []float64
}




func main() {
	
	ft := fasttext.NewFastText("../data/fasttext.db")

	var sent2id = map[string]int{
		"negative": 0,
		"positive": 1,	
	}

	csvFile, err := os.Open("./tabdata.txt")
	defer csvFile.Close()
	reader := csv.NewReader(csvFile)

	reader.Comma = '\t' // Use tab-delimited instead of comma <---- here!

	reader.FieldsPerRecord = -1

	csvData, err := reader.ReadAll()
	if err != nil {
	     fmt.Println(err)
	     os.Exit(1)
	}

	var oneRecord Tweetline

	for _, line := range csvData {
		if line[1] != "neutral" {
			oneRecord.SentimentLabel = line[1]
			oneRecord.SentimentID = sent2id[line[1]]
			oneRecord.Text = line[2]
			oneRecord.TextVec, _ = wordvec.SentenceVec(ft, line[2])
			fmt.Println(oneRecord)

	 	}
	}
}