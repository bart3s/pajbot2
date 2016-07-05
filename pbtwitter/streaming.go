package pbtwitter

import (
	"strings"

	"github.com/dghubble/go-twitter/twitter"
)

// Follow follows given users timeline stream
func (bot *Bot) Follow(user string) {
	if user == "ALL" {
		bot.Following = bot.Client.followedUsers
		return
	}
	u := strings.ToLower(user)
	for _, usr := range bot.Following {
		if usr == u {
			return
		}
	}
	bot.Following = append(bot.Following, u)
	go bot.Client.Follow(u)
}

func (c *Client) streamToBots(tweet *twitter.Tweet) {
	log.Debug(tweet.Text)
	log.Debug(tweet.Retweeted)
	log.Debugf("%+v", tweet)
	if tweet.RetweetedStatus != nil || tweet.QuotedStatus != nil {
		log.Debug("RETWEETED OR QUOTED TWEET, NOT STREAMING")
		return
	}
	for _, bot := range c.Bots {
		for _, followedUser := range bot.Following {
			if strings.ToLower(tweet.User.ScreenName) == followedUser {
				bot.Stream <- tweet
			}
		}
	}
}

// stream starts the stream
func (c *Client) stream() {
	demux := twitter.NewSwitchDemux()
	demux.Tweet = func(tweet *twitter.Tweet) {
		go c.streamToBots(tweet)
	}
	params := &twitter.StreamUserParams{
		With:          "followings",
		StallWarnings: twitter.Bool(true),
	}
	stream, err := c.StreamClient.Streams.User(params)
	if err != nil {
		log.Fatal(err)
	}
	demux.HandleChan(stream.Messages)
}
