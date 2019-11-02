package slack

import (
	"cloud.google.com/go/datastore"
	"context"
	"encoding/json"
	"golang.org/x/oauth2"
	"log"
	"net/http"
	"os"
	"time"
)

const (
	Cookie = "c"
	Okay   = "Okay"
	Error  = "Error"
)

var (
	commands = map[string]Command{}
	conf     = oauth2.Config{
		ClientID:     "",
		ClientSecret: "",
		Scopes:       []string{},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://slack.com/oauth/authorize",
			TokenURL: "https://slack.com/api/oauth.access",
		},
	}
	client *datastore.Client
)

func init() {
	ctx := context.Background()
	var err error
	client, err = datastore.NewClient(ctx, os.Getenv("PROJECT_ID"))
	if err != nil {
		log.Printf("ERROR - Failed to connect to datastore: %v", err)
	}
}

type D map[string]interface{}
type Args struct {
	TeamId      string
	TeamDomain  string
	ChannelID   string
	ChannelName string
	UserID      string
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

func Configure(clientId string, clientSecret string) {
	conf.ClientID = clientId
	conf.ClientSecret = clientSecret
}

func Register(name string, cmd Command) {
	commands["/"+name] = cmd
}

func writeJson(w http.ResponseWriter, r *http.Request, data D) {
	bytes, err := json.Marshal(data)
	if err != nil {
		log.Printf("ERROR - Failed to mashal %v: %v", data, err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(bytes)
}

func Route(w http.ResponseWriter, r *http.Request) {
	args := Args{
		TeamId:      r.FormValue("team_id"),
		TeamDomain:  r.FormValue("team_domain"),
		ChannelID:   r.FormValue("channel_id"),
		ChannelName: r.FormValue("channel_name"),
		UserID:      r.FormValue("user_id"),
		UserName:    r.FormValue("user_name"),
		Command:     r.FormValue("command"),
		Text:        r.FormValue("text"),
		ResponseUrl: r.FormValue("response_url"),
	}
	log.Printf("INFO - Got command %v", args)
	cmd, ok := commands[args.Command]
	if ok {
		writeJson(w, r, cmd(args))
	} else {
		w.Write([]byte("Command not found."))
	}
}

func Oauth(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query()["code"]
	if len(code) == 0 {
		http.Redirect(w, r, "/", 303)
		return
	}
	c := r.Context()
	tok, err := conf.Exchange(c, code[0])
	if err != nil || !tok.Valid() {
		log.Printf("ERROR - Failed to exchange token %v: %v", tok, err)
		http.SetCookie(w, &http.Cookie{
			Name:  Cookie,
			Value: Error,
		})
		http.Redirect(w, r, "/", 303)
		return
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
	key := datastore.NameKey("token", team.TeamId, nil)
	if client == nil {
		log.Printf("ERROR - Datastore unusable, got team %+v", team)
	} else {
		_, err := client.Put(c, key, &team)
		if err != nil {
			log.Printf("ERROR - Failed to save team %+v", team)
		}
	}
	http.SetCookie(w, &http.Cookie{
		Name:  Cookie,
		Value: Okay,
	})
	http.Redirect(w, r, "/", 303)
}
