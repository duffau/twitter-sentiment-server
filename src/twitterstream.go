package main

import (
	"log"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
)

func check(e error) {
    if e != nil {
        panic(e)
    }
}


func TwitterStream(trackQuery string, language string) (*twitter.Stream, error) {
	
	config := oauth1.NewConfig(secretConfig.ConsumerKey, secretConfig.ConsumerSecret)
	token := oauth1.NewToken(secretConfig.AccessToken, secretConfig.AccessSecret)
	httpClient := config.Client(oauth1.NoContext, token)
	client := twitter.NewClient(httpClient)
	
	params := &twitter.StreamFilterParams{
		Language: []string{language},
    	Track: []string{trackQuery},
    	StallWarnings: twitter.Bool(true),
	}
	
	return client.Streams.Filter(params)
}

func TwitterDemux(ch chan twitter.Tweet) twitter.SwitchDemux {
	
	demux := twitter.NewSwitchDemux()

	demux.Tweet = func(tweet *twitter.Tweet) {
		//log.Println(tweet)
		ch <- *tweet
	}
	
    return demux
}

func listen(ch <-chan twitter.Tweet) {
	for tweet := range ch {
		log.Println(tweet)
	} 

}