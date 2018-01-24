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
)

var upgrader = websocket.Upgrader{}


type Message struct {
	Timestamp int64
	Sentiment float64
	Text string
}



func logError(err error) {
	if err != nil {
		log.Println(err)
	}	
}




func makeWebsocketHandler(ch <-chan Tweet, ft *fasttext.FastText,vGood []float64, vBad []float64) http.HandlerFunc {

	var msg Message
	var sentiment float64
	var tweetVec  []float64

	return func (w http.ResponseWriter, r *http.Request) {
		
		conn, err := upgrader.Upgrade(w, r, nil)	
		logError(err)

		for tweet := range ch {
			//log.Println(tweet)
			tweetVec = wordvec.SentenceVec(ft, tweet.Text)
			sentiment = wordvec.Sentiment(vGood, vBad, tweetVec, wordvec.CosineSimilarity)
			msg = Message{time.Now().Unix(), sentiment, tweet.Text}
			err = conn.WriteJSON(msg)
			if err != nil {
				logError(err)
				return
			}
		}
	}

}

func main() {

	twitterChannel := mockTwitterChannel() 
	ft := fasttext.NewFastText("./database/fasttext.db")
	vecNormalizedGood, _ := wordvec.Normalize(wordvec.WordVec(ft, "good"))
	vecNormalizedBad, _ := wordvec.Normalize(wordvec.WordVec(ft, "bad"))

	http.Handle("/", http.FileServer(http.Dir("public/")))
	http.HandleFunc("/websocket", makeWebsocketHandler(twitterChannel, ft, vecNormalizedGood, vecNormalizedBad))
	http.ListenAndServe(":8080", nil)
}