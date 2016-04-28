package dicebot

import (
	"github.com/arkie/hackyslack2"
	"html/template"
	"net/http"
)

var (
	templates = template.Must(template.ParseGlob("template/*.html"))
)

func init() {
	hackyslack.Configure(clientId, clientSecret)

	http.HandleFunc("/command", hackyslack.Route)
	http.HandleFunc("/oauth", hackyslack.Oauth)

	http.HandleFunc("/", index)
	http.HandleFunc("/contact", contact)
	http.HandleFunc("/privacy", privacy)
}

func index(w http.ResponseWriter, r *http.Request) {
	c, _ := r.Cookie(hackyslack.Cookie)
	http.SetCookie(w, &http.Cookie{
		Name:   hackyslack.Cookie,
		MaxAge: -1,
	})
	if c != nil {
		s := "Installed Dicebot."
		if c.Value != hackyslack.Okay {
			s = "Error Installing"
		}
		templates.ExecuteTemplate(w, "index.html", s)
	} else {
		templates.ExecuteTemplate(w, "index.html", "")
	}
}

func contact(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "contact.html", true)
}

func privacy(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "privacy.html", true)
}
