package main

import (
	"log"
	//"gopkg.in/yaml.v2"
    //"io/ioutil"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
)

func check(e error) {
    if e != nil {
        panic(e)
    }
}




func TwitterDemux(configFilePath string, ch chan twitter.Tweet) (twitter.SwitchDemux, *twitter.Stream) {
	return makeTwitterDemux(configFilePath, ch)
}

func makeTwitterDemux(configFilePath string, ch chan twitter.Tweet) (twitter.SwitchDemux, *twitter.Stream) {
    //var c config 
    //c.getConf(configFilePath)

    config := oauth1.NewConfig(secretConfig.ConsumerKey, secretConfig.ConsumerSecret)
	token := oauth1.NewToken(secretConfig.AccessToken, secretConfig.AccessSecret)
	httpClient := config.Client(oauth1.NoContext, token)
	client := twitter.NewClient(httpClient)
	
	params := &twitter.StreamFilterParams{
		Language: []string{"en"},
    	Track: []string{"#bitcoin"},
    	StallWarnings: twitter.Bool(true),
	}
	
	stream, err := client.Streams.Filter(params)
	
	demux := twitter.NewSwitchDemux()

	demux.Tweet = func(tweet *twitter.Tweet) {
		//log.Println(tweet)
		ch <- *tweet
	}
	
	if err != nil {
		log.Printf("Error client.Streams: %v\n", err)
	} 
    return demux, stream
}

func listen(ch <-chan twitter.Tweet) {
	for tweet := range ch {
		log.Println(tweet)
	} 

}