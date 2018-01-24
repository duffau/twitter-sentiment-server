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


type AppData struct {
	ch <-chan Tweet
	ft *fasttext.FastText
	vgood []float64
	vbad []float64
}

/*func indexHandler(w http.ResponseWriter, r *http.Request) {
	html, _ := ioutil.ReadFile("./public/index.html")
	w.Write(html)
}
*/

func logError(err error) {
	if err != nil {
		log.Println(err)
	}	
}


func (ad *AppData) websocketHandler(w http.ResponseWriter, r *http.Request) {
	
	var msg Message
	var sentiment float64
	var tweetVec  []float64
	var tweet Tweet

	conn, err := upgrader.Upgrade(w, r, nil)	
	logError(err)

	for {
		tweet = <-ad.ch
		log.Println(tweet)
		tweetVec = wordvec.SentenceVec(ad.ft, tweet.Text)
		sentiment = wordvec.Sentiment(ad.vgood, ad.vbad, tweetVec, wordvec.CosineSimilarity)
		msg = Message{time.Now().Unix(), sentiment, tweet.Text}
		err = conn.WriteJSON(msg)
		if err != nil {
			logError(err)
			return
		}
	}
}


func main() {

	twitterChannel := mockTwitterChannel() 
	ft := fasttext.NewFastText("./database/fasttext.db")
	vecNormalizedGood, _ := wordvec.Normalize(wordvec.WordVec(ft, "good"))
	vecNormalizedBad, _ := wordvec.Normalize(wordvec.WordVec(ft, "bad"))

	ad := AppData{twitterChannel, ft, vecNormalizedGood, vecNormalizedBad}

	http.Handle("/", http.FileServer(http.Dir("public/")))
	http.HandleFunc("/websocket", ad.websocketHandler)
	http.ListenAndServe(":8080", nil)
}