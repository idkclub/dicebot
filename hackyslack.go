package hackyslack

import (
	"encoding/json"
	"golang.org/x/oauth2"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
	"html/template"
	"net/http"
	"time"
)

var (
	templates = template.Must(template.ParseGlob("*.html"))
	commands  = map[string]Command{}
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
	http.HandleFunc("/command", command)
	http.HandleFunc("/contact", contact)
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

// Manually inline oauth2.Token for datastore.
type TeamToken struct {
	AccessToken  string    `json:"access_token"`
	TokenType    string    `json:"token_type,omitempty"`
	RefreshToken string    `json:"refresh_token,omitempty"`
	Expiry       time.Time `json:"expiry,omitempty"`
	TeamId       string    `json:"team_id,omitempty"`
	TeamName     string    `json:"team_name,omitempty"`
	Scope        string    `json:"scope,omitempty"`
	Created      time.Time `json:"created,omitempty"`
}

func Register(name string, cmd Command) {
	commands["/"+name] = cmd
	http.HandleFunc("/"+name, func(w http.ResponseWriter, r *http.Request) {
		args := Args{
			TeamId:      r.FormValue("team_id"),
			TeamDomain:  r.FormValue("team_domain"),
			ChannelId:   r.FormValue("channel_id"),
			ChannelName: r.FormValue("channel_name"),
			UserId:      r.FormValue("user_id"),
			UserName:    r.FormValue("user_name"),
			Command:     r.FormValue("command"),
			Text:        r.FormValue("text"),
			ResponseUrl: r.FormValue("response_url"),
		}
		c := appengine.NewContext(r)
		log.Infof(c, "Got command %v", args)
		writeJson(w, r, cmd(args))
	})
}

func writeJson(w http.ResponseWriter, r *http.Request, data D) {
	bytes, err := json.Marshal(data)
	if err != nil {
		c := appengine.NewContext(r)
		log.Errorf(c, "Failed to mashal %v: %v", data, err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(bytes)
}

func command(w http.ResponseWriter, r *http.Request) {
	args := Args{
		TeamId:      r.FormValue("team_id"),
		TeamDomain:  r.FormValue("team_domain"),
		ChannelId:   r.FormValue("channel_id"),
		ChannelName: r.FormValue("channel_name"),
		UserId:      r.FormValue("user_id"),
		UserName:    r.FormValue("user_name"),
		Command:     r.FormValue("command"),
		Text:        r.FormValue("text"),
		ResponseUrl: r.FormValue("response_url"),
	}
	c := appengine.NewContext(r)
	log.Infof(c, "Got command %v", args)
	cmd, ok := commands[args.Command]
	if ok {
		writeJson(w, r, cmd(args))
	} else {
		w.Write([]byte("Command not found."))
	}
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
	c := appengine.NewContext(r)
	tok, err := conf.Exchange(c, code[0])
	if err != nil {
		log.Errorf(c, "Failed to exchange token %v: %v", tok, err)
	}
	team := TeamToken{
		AccessToken:  tok.AccessToken,
		TokenType:    tok.TokenType,
		RefreshToken: tok.RefreshToken,
		Expiry:       tok.Expiry,
		TeamId:       tok.Extra("team_id").(string),
		TeamName:     tok.Extra("team_name").(string),
		Scope:        tok.Extra("scope").(string),
		Created:      time.Now(),
	}
	key := datastore.NewKey(c, "token", team.TeamId, 0, nil)
	datastore.Put(c, key, &team)
	http.SetCookie(w, &http.Cookie{
		Name:  "s",
		Value: "1",
	})
	http.Redirect(w, r, "/", 303)
}

func contact(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "contact.html", true)
}

func privacy(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "privacy.html", true)
}
