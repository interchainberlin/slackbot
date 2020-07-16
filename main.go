// https://github.com/slack-go/slack
// https://guzalexander.com/2017/09/15/cowsay-slack-command.html

package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"encoding/json"
)

var (
	port  string = "80"
	token string
)

func init() {
	token = os.Getenv("VERIFICATION_TOKEN")
	if "" == token {
		panic("VERIFICATION_TOKEN is not set!")
	}

	if "" != os.Getenv("PORT") {
		port = os.Getenv("PORT")
	}
}

func botHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	if token != r.FormValue("token") {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	r.ParseForm()

	// field info here: https://api.slack.com/interactivity/slash-commands

	command := r.Form["command"][0]
	userid := r.Form["user_id"][0]
	username := r.Form["user_name"][0]
	channelid := r.Form["channel_id"][0]
	channelname := r.Form["channel_name"][0]
	text := r.Form["text"][0]

	fmt.Println("command: ", command, "\nuser_id: ", userid, "\nuser_name: ", username, "\nchannel_id", channelid, "\nchannel_name:", channelname, "\ntext: ", text)

	botReply := "hello!"

	jsonResp, _ := json.Marshal(struct {
		Type string `json:"response_type"`
		Text string `json:"text"`
	}{
		Type: "in_channel",
		Text: fmt.Sprintf("```%s```", botReply),
	})

	w.Header().Add("Content-Type", "application/json")
	fmt.Fprintf(w, string(jsonResp))
}

func main() {
	http.HandleFunc("/", botHandler)
	log.Fatalln(http.ListenAndServe(":"+port, nil))
}
