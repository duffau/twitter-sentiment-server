package main

import (
	_ "github.com/mattn/go-sqlite3"
	"github.com/ekzhu/go-fasttext"
	"fmt"
	"bitbucket.com/chrduffau/wordvec"
)

func main() {
	ft := fasttext.NewFastText("../database/fasttext.db")

	vecNormalizedGood, _ := wordvec.Normalize(wordvec.WordVec(ft, "good"))
	vecNormalizedBad, _ := wordvec.Normalize(wordvec.WordVec(ft, "bad"))

	tweet1 := "South Korea Urges 23 Countries, EU, and IMF to Collaborate on Curbing Crypto Trading #cryptocurrencynews… https://t.co/uRsu8z1sQG"
	tweet2 := "Prof. Jeremy Siegel offers his views on the Dow, the Fed, GDP growth, inflation, wage growth, a #Bitcoin bubble, and why future GOP economic legislation faces a tough road in 2018."

	vecTweet1 := wordvec.SentenceVec(ft, tweet1)
	vecTweet2 := wordvec.SentenceVec(ft, tweet2)
	fmt.Printf("Tweet1 similarity with good: %.3f\n", wordvec.CosineSimilarity(vecTweet1, vecNormalizedGood))
	fmt.Printf("Tweet1 similarity with bad: %.3f\n", wordvec.CosineSimilarity(vecTweet1, vecNormalizedBad))
	fmt.Printf("Tweet2 similarity with good: %.3f\n", wordvec.CosineSimilarity(vecTweet2, vecNormalizedGood))
	fmt.Printf("Tweet2 similarity with bad: %.3f\n", wordvec.CosineSimilarity(vecTweet2, vecNormalizedBad))

	fmt.Printf("sentiment Tweet1  = %.3f, sentiment Tweet2 = %.3f\n",
		wordvec.Sentiment(vecNormalizedGood, vecNormalizedBad, vecTweet1, wordvec.CosineSimilarity),
		wordvec.Sentiment(vecNormalizedGood, vecNormalizedBad, vecTweet2, wordvec.CosineSimilarity))
}



func printWordVec(ft *fasttext.FastText, s string) {
	vec, err := ft.GetEmb(s)
	if err == nil {
		norm := wordvec.Norm(vec)
		fmt.Printf("word = %s, vector elements = %d, vector norm = %f\nvector =\n%v\n", s, len(vec), norm, vec)
	} else {
		fmt.Printf("Error: %s. Word: \"%s\"\n", err, s)
	}

}
