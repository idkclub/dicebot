package hackyslack

import (
	"encoding/json"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"html/template"
	"log"
	"net/http"
)

var (
	templates = template.Must(template.ParseGlob("*.html"))
	conf      = oauth2.Config{
		ClientID:     clientId,
		ClientSecret: clientSecret,
		Scopes:       []string{},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://slack.com/oauth/authorize",
			TokenURL: "https://slack.com/api/oauth.access",
		},
	}
)

func init() {
	http.HandleFunc("/", index)
	http.HandleFunc("/oauth", oauth)
	http.HandleFunc("/privacy", privacy)
}

type D map[string]interface{}
type Args struct {
	TeamId      string
	TeamDomain  string
	ChannelId   string
	ChannelName string
	UserId      string
	UserName    string
	Command     string
	Text        string
	ResponseUrl string
}
type Command func(Args) D

func Register(name string, cmd Command) {
	http.HandleFunc("/"+name, func(w http.ResponseWriter, r *http.Request) {
		writeJson(w, cmd(Args{
			TeamId:      r.FormValue("team_id"),
			TeamDomain:  r.FormValue("team_domain"),
			ChannelId:   r.FormValue("channel_id"),
			ChannelName: r.FormValue("channel_name"),
			UserId:      r.FormValue("user_id"),
			UserName:    r.FormValue("user_name"),
			Command:     r.FormValue("command"),
			Text:        r.FormValue("text"),
			ResponseUrl: r.FormValue("response_url"),
		}))
	})
}

func writeJson(w http.ResponseWriter, data D) {
	bytes, err := json.Marshal(data)
	if err != nil {
		log.Fatal(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(bytes)
}

func index(w http.ResponseWriter, r *http.Request) {
	c, _ := r.Cookie("s")
	http.SetCookie(w, &http.Cookie{
		Name:   "s",
		MaxAge: -1,
	})
	templates.ExecuteTemplate(w, "index.html", c != nil)
}

func oauth(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query()["code"]
	if len(code) == 0 {
		http.Redirect(w, r, "/", 303)
		return
	}
	var ctx context.Context = appengine.NewContext(r)
	tok, err := conf.Exchange(ctx, code[0])
	if err != nil {
		log.Fatal(err)
	}
	log.Print(w, "%s", tok)
	key := datastore.NewKey(ctx, "token", tok.Extra("team_id").(string), 0, nil)
	datastore.Put(ctx, key, tok)
	http.SetCookie(w, &http.Cookie{
		Name:  "s",
		Value: "1",
	})
	http.Redirect(w, r, "/", 303)
}

func privacy(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "privacy.html", true)
}
