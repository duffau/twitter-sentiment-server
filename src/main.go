package main

import (
	_ "github.com/mattn/go-sqlite3"
	"github.com/ekzhu/go-fasttext"
	"time"
	"net/http"
	// "io/ioutil"
	"github.com/gorilla/websocket"
	"log"
	"bitbucket.com/chrduffau/wordvec"
	"github.com/dghubble/go-twitter/twitter"
)

var upgrader = websocket.Upgrader{}


type Message struct {
	Timestamp int64
	Sentiment float64
	SentimentFilter float64
	Text string
}



func logError(err error, fname string) {
	if err != nil {
		log.Printf("%s failed with message: %v \n",fname, err)
	}	
}


func expMovAvgFilter(currentValue *float64, FilterValue *float64, weight float64) {
	*FilterValue = *FilterValue * weight + *currentValue * (1 - weight) 
}

func makeWebsocketHandler(ch <-chan twitter.Tweet, ft *fasttext.FastText,vGood []float64, vBad []float64) http.HandlerFunc {

	var msg Message
	var sentiment float64
	var sentimentFilter float64
	var tweetVec  []float64
	const ExpMAweight = 0.8

	return func (w http.ResponseWriter, r *http.Request) {
		
		conn, err := upgrader.Upgrade(w, r, nil)	
		logError(err,"makeWebsocketHandler")

		for tweet := range ch {

			tweetVec, err = wordvec.SentenceVec(ft, tweet.Text)
			logError(err, "makeWebsocketHandler")

			sentiment, err = wordvec.Sentiment(vGood, vBad, tweetVec, wordvec.CosineSimilarity)
			logError(err, "makeWebsocketHandler")
			expMovAvgFilter(&sentiment, &sentimentFilter, ExpMAweight)

			msg = Message{time.Now().Unix(), sentiment, sentimentFilter, tweet.Text}
			err = conn.WriteJSON(msg)
			if err != nil {
				logError(err, "makeWebsocketHandler")
				return
			}
		}
	}

}

func main() {

	//twitterChannel := mockTwitterChannel()
	twitterChannel := make(chan twitter.Tweet)
	twitterDemux, twitterSream := TwitterDemux("SecretSettings.yml", twitterChannel)
	go twitterDemux.HandleChan(twitterSream.Messages)
	
	//listen(twitterChannel)

	ft := fasttext.NewFastText("./database/fasttext.db")
	vecNormalizedGood, _ := wordvec.Normalize(wordvec.WordVec(ft, "good"))
	vecNormalizedBad, _ := wordvec.Normalize(wordvec.WordVec(ft, "bad"))

	http.Handle("/", http.FileServer(http.Dir("public/")))
	http.HandleFunc("/websocket", makeWebsocketHandler(twitterChannel, ft, vecNormalizedGood, vecNormalizedBad))
	http.ListenAndServe(":8080", nil)
}