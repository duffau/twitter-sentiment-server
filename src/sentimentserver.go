package main

import (
	"net/http"
	"log"
	"time"
	"math"
	"errors"
	"github.com/gorilla/websocket"
	"github.com/dghubble/go-twitter/twitter"
	_ "github.com/mattn/go-sqlite3"
	"bitbucket.com/chrduffau/wordvec"
	"bitbucket.com/chrduffau/twitter-sentiment-server/utils"
	"github.com/ekzhu/go-fasttext"
	"github.com/cdipaolo/goml/base"
	"github.com/cdipaolo/goml/linear"
)

func True(r *http.Request) bool {
	return true
}

var upgrader = websocket.Upgrader{
	CheckOrigin: True,
}

var clients =  make(map[*websocket.Conn]bool)


type SentimentMessage struct {
	Timestamp int64
	Sentiment float64
	SentimentFilter float64
	Text string
}

type ClientMessage struct {
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


func handleWebSocketConnections(w http.ResponseWriter, r *http.Request) {
		wsConn, err := upgrader.Upgrade(w, r, nil)	
		logError(err,"WebsocketHandler")
		clients[wsConn] = true
		log.Printf("New client connected. Clients: %v", len(clients))
}


func publishSentimentMessages(ch <-chan twitter.Tweet, ft *fasttext.FastText, sentimentmodel *linear.Logistic) {

	var (
		err error
		msg SentimentMessage
		sentimentFilter float64
	)

	for tweet := range ch {

		msg = calcSentiment(tweet, ft, sentimentmodel, &sentimentFilter)

		for clientConn := range clients {
			err = clientConn.WriteJSON(msg)
			if err != nil {
				logError(err, "WebsocketHandler")
				clientConn.Close()
				delete(clients, clientConn)
				log.Printf("Client connection closed. Clients: %v", len(clients))
			}

		}
	}

}


func calcSentiment(tweet  twitter.Tweet, ft *fasttext.FastText, sentimentmodel *linear.Logistic, sentimentFilter *float64) SentimentMessage {
		var (
			sentimentOut []float64
			sentiment float64
			tweetVec  []float64
			tokenizedTweetText string
			err error
		)

		const ExpMAweight = 0.8

		tokenizedTweetText = utils.Tokenize(tweet.Text)
		tweetVec, _ = wordvec.SentenceVec(ft, tokenizedTweetText)
		//logError(err, "calcSentiment")
		if len(tweetVec) > 0 {
			sentimentOut, err = sentimentmodel.Predict(tweetVec)
			logError(err, "calcSentiment")
		} else {
			sentimentOut = []float64{0.5}
			err = errors.New("Tweet Vector has length zero. Sentiment set to 0.5.")
			logError(err, "calcSentiment")
		}
		
		if math.IsNaN(sentimentOut[0]) {
			sentiment = 0.5
		} else {
			sentiment = sentimentOut[0]
			logError(errors.New("Sentiment evaluated to NaN, set to 0.5."), "calcSentiment")
		}
		sentiment = rescaleSentiment(sentiment)
		expMovAvgFilter(&sentiment, sentimentFilter, ExpMAweight)

		return SentimentMessage{time.Now().Unix(), sentiment, *sentimentFilter, tweet.Text}
}


func rescaleSentiment(sentiment float64) float64 {
	return sentiment*2.0 - 1.0
}

func main() {

	//twitterChannel := mockTwitterChannel()
	twitterChannel := make(chan twitter.Tweet)
	twitterDemux := TwitterDemux(twitterChannel)
	twitterStream, err := TwitterStream("#bitcoin", "en")
	logError(err, "main")
	if err != nil {
		return
	}

	go twitterDemux.HandleChan(twitterStream.Messages)

	//listen(twitterChannel)

	ft := fasttext.NewFastText("./database/glove25.db")
	sentimentModel := linear.NewLogistic(base.BatchGA, .0001, 1, 100, nil, nil, 25)
	err = sentimentModel.RestoreFromFile("logistic_sentiment.model")
	logError(err, "main")

	go publishSentimentMessages(twitterChannel, ft, sentimentModel)

	http.Handle("/", http.FileServer(http.Dir("public/")))
	http.HandleFunc("/websocket", handleWebSocketConnections)
	http.ListenAndServe(":8080", nil)
}