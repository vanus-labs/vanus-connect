# Source Twitter Proposal

## Description

The Twitter Source is used to pull data from a Twitter user's account, transform it into CloudEvents and send the events to the target. 

This page describes the design of the Twitter Source in detail.

## Programming Language

-[x] Golang
-[ ] Java

## Prerequisites

- A Twitter account
- Registered developer account at [Twitter Developer Portal](https://developer.twitter.com/)

## Connector Details

### Configuration

- install module [go-twitter API library](github.com/dghubble/go-twitter/twitter)
- install module [OAuth1](github.com/dghubble/oauth1)
- add an application at [Twitter Developer Portal](https://developer.twitter.com/)
- from developer portal get consumer Keys and Tokens and add them to configuration file

### Connector Behavior

* The connector receives events from Twitter: when a user’s tweet is liked, retweeted, quote-tweeted, or commented or when someone Dms the user.
* The connector puts the data into a CloudEvent to send it to the target.
* Each event type (like, dm, comment, etc ...) is specified in the event type specification of the CloudEvent.

Draft code:

```
import (
	"fmt"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
)

func main() {
	// Replace the values below with your own Twitter API credentials
	consumerKey := "YOUR_CONSUMER_KEY"
	consumerSecret := "YOUR_CONSUMER_SECRET"
	accessToken := "YOUR_ACCESS_TOKEN"
	accessSecret := "YOUR_ACCESS_SECRET"

	// Set up the OAuth1 config with your credentials
	config := oauth1.NewConfig(consumerKey, consumerSecret)
	token := oauth1.NewToken(accessToken, accessSecret)

	// Use the config and token to create an http.Client
	httpClient := config.Client(oauth1.NoContext, token)

	// Use the http.Client to create a twitter client
	client := twitter.NewClient(httpClient)

	// Set up the filter to listen for events
	filterParams := &twitter.StreamFilterParams{
		Track:         []string{"@username"}, // replace with your own keyword(s) to track
		StallWarnings: twitter.Bool(true),
	}

	stream, err := client.Streams.Filter(filterParams)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Print each event to the console
	go func() {
		for {
			event := <-stream.Messages
			switch tweet := event.(type) {
			case *twitter.Tweet:
				fmt.Printf("Received a tweet: ID: %v, Text: %v\n", tweet.ID, tweet.Text)
			case *twitter.DirectMessage:
				fmt.Printf("Received a DM: ID: %v, Text: %v\n", tweet.ID, tweet.Text)
			case *twitter.Retweet:
				fmt.Printf("Received a retweet: ID: %v, Retweeted tweet ID: %v\n", tweet.ID, tweet.IDStr)
			case *twitter.TweetDeleteEvent:
				fmt.Printf("Received a tweet delete event: ID: %v, Deleted status ID: %v\n", tweet.ID, tweet.Delete.Status.ID)
			case *twitter.DirectMessageDeleteEvent:
				fmt.Printf("Received a DM delete event: ID: %v, Deleted DM ID: %v\n", tweet.ID, tweet.Delete.DirectMessage.ID)
			case *twitter.EventTweet:
				fmt.Printf("Received a tweet event: ID: %v, Event type: %v\n", tweet.ID, tweet.Event)
			default:
				fmt.Printf("Received an unhandled event of type: %T\n", tweet)
			}
		}
	}()

	// Wait for a signal to stop the stream
	ch := make(chan struct{})
	go func() {
		fmt.Println("Press ENTER to stop the stream")
		fmt.Scanln()
		ch <- struct{}{}
	}()
	<-ch

	// Stop the stream
	fmt.Println("Stopping the stream...")
	stream.Stop()
	fmt.Println("Stream stopped")
}
```

### Used Libraries/APIs

- [go-twitter API library](github.com/dghubble/go-twitter/twitter)
- [OAuth1](github.com/dghubble/oauth1)