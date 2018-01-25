package main

import (
	"log"
	//"os"
	//"os/signal"
	//"syscall"
	//"fmt"
	"gopkg.in/yaml.v2"
    "io/ioutil"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
)

func check(e error) {
    if e != nil {
        panic(e)
    }
}

type config struct {
	ConsumerKey string `yaml:"TWITTER_PUBLIC_CONSUMER_KEY"`
	ConsumerSecret string `yaml:"TWITTER_SECRET_CONSUMER_KEY"`
	AccessToken string `yaml:"TWTTER_PUBLIC_ACCESS_TOKEN"`
	AccessSecret string `yaml:"TWTTER_SECRET_ACCESS_TOKEN"`
} 

func (c *config) getConf(configFilePath string) *config {

    yamlFile, err := ioutil.ReadFile(configFilePath)
    if err != nil {
        log.Printf("yamlFile.Get err   #%v ", err)
    }
    err = yaml.Unmarshal(yamlFile, c)
    if err != nil {
        log.Fatalf("Unmarshal: %v", err)
    }

    return c
}


func TwitterDemux(configFilePath string, ch chan twitter.Tweet) (twitter.SwitchDemux, *twitter.Stream) {
	return makeTwitterDemux(configFilePath, ch)
}

func makeTwitterDemux(configFilePath string, ch chan twitter.Tweet) (twitter.SwitchDemux, *twitter.Stream) {
    var c config 
    c.getConf(configFilePath)

    config := oauth1.NewConfig(c.ConsumerKey, c.ConsumerSecret)
	token := oauth1.NewToken(c.AccessToken, c.AccessSecret)
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