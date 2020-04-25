package main

import (
	"net/url"
	"os"
	"github.com/ChimeraCoder/anaconda"
	log "github.com/sirupsen/logrus"
)

func getkey(name string) string{
	v := os.Getenv(name)
	if v == "" {
		panic("Missing required environment key "+ name)
	}
	return v
}

func main(){

  twtPubkey := getkey("TWITTER_API_KEY")
	twtsecretKey := getkey("TWITTER_API_SECRET_KEY")
	twtPubToken := getkey("TWITTER_ACCESS_TOKEN")
	twtsecretToken := getkey("TWITTER_ACCESS_TOKEN_SECRET")

	anaconda.SetConsumerKey(twtPubkey)
	anaconda.SetConsumerSecret(twtsecretKey)
	api := anaconda.NewTwitterApi(twtPubToken, twtsecretToken)

	getLog := &logger{log.New()}

	api.SetLogger(getLog)

	stream :=	api.PublicStreamFilter(url.Values{
			"track": []string{"#rapuluchukwu"},
		})

	defer	stream.Stop()

	for v := range stream.C{
		t, ok := v.(anaconda.Tweet)
		if !ok {
			getLog.Warnf("Received an unexpected of type %T", v)
		}
		if t.RetweetedStatus != nil{
			continue
		}
		_, err := api.Retweet(t.Id, false)
		if err != nil{
			getLog.Errorf("Could not retweet this tweet %d: %v", t.Id, err)
			continue
		}
		getLog.Infof("retweeted %d", t.Id)
	}
}

type logger struct{
	*log.Logger
}

func(log *logger) Critical(arg ...interface{}){ log.Error(arg...) }
func(log *logger) Criticalf(format string, arg ...interface{}){ log.Errorf(format, arg...) }
func(log *logger) Notice(arg ...interface{}){ log.Info(arg...) }
func(log *logger) Noticef(format string, arg ...interface{}){ log.Infof(format, arg...) }