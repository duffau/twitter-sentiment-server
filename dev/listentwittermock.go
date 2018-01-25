package notmain


import (
	"fmt"
	"time"
	"math/rand"
)

type Tweet struct {
	Text string
}

func MagicEightBallTweet() Tweet {
	answers := []string{
		"It is certain",
		"It is decidedly so",
		"Without a doubt",
		"Yes definitely",
		"You may rely on it",
		"As I see it yes",
		"Most likely",
		"Outlook good",
		"Yes",
		"Signs point to yes",
		"Reply hazy try again",
		"Ask again later",
		"Better not tell you now",
		"Cannot predict now",
		"Concentrate and ask again",
		"Don't count on it",
		"My reply is no",
		"My sources say no",
		"Outlook not so good",
		"Very doubtful",
	}
	return Tweet{answers[rand.Intn(20)]}
}


func sendMagicEightBallTweetsToChannel(limit int, ch chan Tweet) {
	var randomDuration time.Duration
	for i := 0; i < limit; i++ {
		ch <- MagicEightBallTweet()
		randomDuration = time.Duration(500 + rand.Intn(1000)) * time.Millisecond 
		//fmt.Printf("randomDuration = %v\n", randomDuration)
		time.Sleep(randomDuration)
	}
	close(ch)
}

func mockTwitterChannel() <-chan Tweet {
	ch := make(chan Tweet)
	go sendMagicEightBallTweetsToChannel(2000, ch)
	return ch	
}

func listen(ch <-chan Tweet) {
	for tweet := range ch {
		fmt.Println(tweet)
	} 

}
