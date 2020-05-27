package api

import (
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
	code := r.FormValue("code")
	if len(code) == 0 {
		http.Redirect(w, r, "/", 303)
		return
	}
	c := r.Context()
	tok, err := conf.Exchange(c, code)
	if err != nil || !tok.Valid() {
		log.Printf("ERROR - Failed to exchange token %v: %v", tok, err)
		http.SetCookie(w, &http.Cookie{
			Name:  Cookie,
			Path:  "/",
			Value: Error,
		})
		http.Redirect(w, r, "/", 303)
		return
	}
	log.Printf("INFO - Got token %+v", tok)
	http.SetCookie(w, &http.Cookie{
		Name:  Cookie,
		Path:  "/",
		Value: Okay,
	})
	http.Redirect(w, r, "/", 303)
}
