package main

import (
	"fmt"
	"os"
	_ "github.com/mattn/go-sqlite3"
	"github.com/ekzhu/go-fasttext"
	"bitbucket.com/chrduffau/wordvec"
	"bitbucket.com/chrduffau/twitter-sentiment-server/utils"
	"encoding/csv"
	"github.com/cdipaolo/goml/base"
	"github.com/cdipaolo/goml/linear"
)

func bool2int(b bool) int {
   if b {
      return 1
   }
   return 0
} 

func readSemEvalTXTtoFloatArrays(path string, ft *fasttext.FastText) ([][]float64, []float64) {
	csvFile, _ := os.Open(path)
	var sentencevec []float64
	var X [][]float64
	var Y []float64
	var tokenizedText string
	
	var label2sent = map[string]float64{
		"negative": 0.0,
		"positive": 1.0,	
	}

	defer csvFile.Close()
	linecount := 0

	reader := csv.NewReader(csvFile)
	reader.Comma = '\t' // Use tab-delimited instead of comma.
	reader.FieldsPerRecord = -1
	csvData, _ := reader.ReadAll()

	for _, line := range csvData {
		linecount += 1
		if line[1] != "neutral" {
			tokenizedText = utils.Tokenize(line[2])
			sentencevec, _ = wordvec.SentenceVec(ft, tokenizedText)
			if len(sentencevec) == 25 {
				X = append(X, sentencevec)
				Y = append(Y, label2sent[line[1]])
			}
	 	}
	 	if linecount % 500 == 0 {
	 		fmt.Printf("Processed %v lines of %v.\n", linecount, path)
	 	}
	}
	fmt.Printf("X has %v elements.\n", len(X))
	fmt.Printf("Y has %v elements.\n", len(Y))
	fmt.Printf("First element of X has %v elements.\n", len(X[0]))
	return X, Y
} 


func main()  {

	// create the channel of data and errors
	ft := fasttext.NewFastText("../database/glove25.db")

	trainX, trainY := readSemEvalTXTtoFloatArrays("./SemEval_training.txt", ft)
	testX, testY := readSemEvalTXTtoFloatArrays("./SemEval_test.txt", ft)


	model := linear.NewLogistic(base.BatchGA, 1e-4, 5, 800, trainX, trainY)
	model.Learn()
	fmt.Println(model)

	correctPreds := 0
	totalPreds := 0
	
	for i, testy := range(testY) {
		predy, _ := model.Predict(testX[i])
		correctPreds += bool2int(((predy[0] > 0.5) && (testy > 0.5)) || ((predy[0] <= 0.5) && (testy < 0.5))) 
		totalPreds += 1
	}

	fmt.Printf("Correct predictions: %v out of %v. Equivalent to %.2f pct.\n", correctPreds, totalPreds, float64(correctPreds)/float64(totalPreds)*100.0)
	
	model.PersistToFile("../logistic_sentiment.model")

}