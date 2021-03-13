package main

import (
	"go.uber.org/zap"
	"golang.org/x/oauth2"
	"log"
	"net/http"
	"os"
)

var (
	clientID     = os.Getenv("SLACK_ID")
	clientSecret = os.Getenv("SLACK_SECRET")
	conf         = oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Scopes:       []string{},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://slack.com/oauth/v2/authorize",
			TokenURL: "https://slack.com/api/oauth.v2.access",
		},
	}
)

const (
	Cookie = "c"
	Okay   = "Okay"
	Error  = "Error"
)

func Oauth(w http.ResponseWriter, r *http.Request) {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Printf("can't initialize zap logger: %v", err)
	}
	code := r.FormValue("code")
	if len(code) == 0 {
		http.Redirect(w, r, "/", 303)
		return
	}
	c := r.Context()
	tok, err := conf.Exchange(c, code)
	if err != nil || !tok.Valid() {
		logger.Error("oauth error", zap.Any("token", tok), zap.Error(err))
		http.SetCookie(w, &http.Cookie{
			Name:  Cookie,
			Path:  "/",
			Value: Error,
		})
		http.Redirect(w, r, "/", 303)
		return
	}
	logger.Info("oauth token", zap.Any("token", tok))
	http.SetCookie(w, &http.Cookie{
		Name:  Cookie,
		Path:  "/",
		Value: Okay,
	})
	http.Redirect(w, r, "/", 303)
}
