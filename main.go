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
		panic("COWSAY_TOKEN is not set!")
	}

	if "" != os.Getenv("PORT") {
		port = os.Getenv("PORT")
	}
}

func cowHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	if token != r.FormValue("token") {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	r.ParseForm()
	text := r.Form["text"]
	fmt.Println(text[0])

	botReply := "moo"

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
	http.HandleFunc("/", cowHandler)
	log.Fatalln(http.ListenAndServe(":"+port, nil))
}
