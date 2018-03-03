package dicebot

import (
	"github.com/arkie/dicebot/slack"
	"html/template"
	"net/http"
	"os"
)

var (
	clientId     = os.Getenv("SLACK_ID")
	clientSecret = os.Getenv("SLACK_SECRET")
	templates    = template.Must(template.ParseGlob("template/*.html"))
)

func init() {
	slack.Configure(clientId, clientSecret)

	http.HandleFunc("/command", slack.Route)
	http.HandleFunc("/oauth", slack.Oauth)

	http.HandleFunc("/", index)
	http.HandleFunc("/contact", contact)
	http.HandleFunc("/privacy", privacy)
}

type page struct {
	Client string
	Status string
}

func index(w http.ResponseWriter, r *http.Request) {
	c, _ := r.Cookie(slack.Cookie)
	http.SetCookie(w, &http.Cookie{
		Name:   slack.Cookie,
		MaxAge: -1,
	})
	if c != nil {
		s := "Installed Dicebot."
		if c.Value != slack.Okay {
			s = "Error Installing"
		}
		templates.ExecuteTemplate(w, "index.html", page{
			Client: clientId,
			Status: s,
		})
	} else {
		templates.ExecuteTemplate(w, "index.html", page{Client: clientId})
	}
}

func contact(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "contact.html", nil)
}

func privacy(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "privacy.html", nil)
}
