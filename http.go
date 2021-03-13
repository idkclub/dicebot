package main

import (
	"html/template"
	"net/http"
)

var (
	templates = template.Must(template.ParseGlob("resources/templates/*.html"))
)

type page struct {
	Client string
	Status string
}

func index(w http.ResponseWriter, r *http.Request) {
	c, _ := r.Cookie(Cookie)
	http.SetCookie(w, &http.Cookie{
		Name:   Cookie,
		MaxAge: -1,
	})
	if c != nil {
		s := "Installed Dicebot."
		if c.Value != Okay {
			s = "Error Installing"
		}
		templates.ExecuteTemplate(w, "index.html", page{
			Client: clientID,
			Status: s,
		})
	} else {
		templates.ExecuteTemplate(w, "index.html", page{Client: clientID})
	}
}

func privacy(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "privacy.html", nil)
}

func main() {
	fs := http.FileServer(http.Dir("./resources/assets"))
	http.Handle("/assets/", http.StripPrefix("/assets/", fs))
	http.HandleFunc("/oauth", Oauth)
	http.HandleFunc("/roll", Roll)
	http.HandleFunc("/", index)
	http.HandleFunc("/privacy", privacy)
	http.ListenAndServe(":8080", nil)
}
