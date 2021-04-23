// https://github.com/slack-go/slack
// https://guzalexander.com/2017/09/15/cowsay-slack-command.html

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/hako/durafmt"
	"github.com/slack-go/slack"
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

	apitoken := os.Getenv("API_TOKEN")
	if "" == apitoken {
		panic("API_TOKEN is not set!")
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

	responseURL := r.Form["response_url"][0]

	// field info here: https://api.slack.com/interactivity/slash-commands

	command := r.Form["command"][0]
	userid := r.Form["user_id"][0]
	username := r.Form["user_name"][0]
	channelid := r.Form["channel_id"][0]
	channelname := r.Form["channel_name"][0]
	text := r.Form["text"][0]
	textArray := strings.Split(text, " ")

	fmt.Println("command: ", command, "\nuser_id: ", userid, "\nuser_name: ", username, "\nchannel_id", channelid, "\nchannel_name:", channelname, "\ntext: ", text)
	fmt.Printf("textArray:'%s'\n", textArray)
	fmt.Printf("textArray len:'%d'\n", len(textArray))

	botReply := "Processing your request, please standby ⏳"
	jsonResp, _ := json.Marshal(struct {
		Type string `json:"response_type"`
		Text string `json:"text"`
	}{
		Type: "in_channel",
		Text: fmt.Sprintf("```%s```", botReply),
	})

	w.Header().Add("Content-Type", "application/json")
	fmt.Fprintf(w, string(jsonResp))

	go handleCommand(responseURL, command, userid, textArray)
}

func handleCommand(responseURL, command, userid string, textArray []string) {
	var botReply string
	switch command {
	case "/brrr":
		botReply = brrr(userid, textArray)
	case "/send":
		botReply = send(userid, textArray)
	case "/balance":
		botReply = balance(userid, textArray)
	case "/til-brrr":
		botReply = tilbrrr(userid, textArray)
	default:
		botReply = fmt.Sprintf("Sorry I don't understand that command %s.", command)
	}

	jsonResp, _ := json.Marshal(struct {
		Type    string `json:"response_type"`
		Replace bool   `json:"replace_original"`
		Text    string `json:"text"`
	}{
		Type:    "in_channel",
		Replace: true,
		Text:    fmt.Sprintf("```%s```", botReply),
	})

	fmt.Println("responseURL", responseURL)
	resp, err := http.Post(responseURL, "application/json", bytes.NewBuffer(jsonResp))
	if err != nil {
		fmt.Println("err", err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("err", err)
	}

	fmt.Println(string(body))

}
func getUserID(userID string) (string, string, error) {
	api := slack.New(os.Getenv("API_TOKEN"))
	// fmt.Println("userID", userID)
	// fmt.Println("usernameOriginal", usernameOriginal)
	// username := strings.Split(usernameOriginal, "|")[1]
	// username = username[:len(username)-1]
	// fmt.Println("username", username)

	user, err := api.GetUserInfo(userID)
	if err != nil {
		fmt.Printf("%s\n", err)
		return userID, "", err
	}

	username := user.Profile.DisplayNameNormalized
	return user.ID, username, nil
	// fmt.Printf("ID: %s, Fullname: %s, Email: %s\n", user.ID, user.Profile.RealName, user.Profile.Email)
}

const ShellToUse = "bash"

func Shellout(command string) (error, string, string) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd := exec.Command(ShellToUse, "-c", command)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	return err, stdout.String(), stderr.String()
}

func createNewUserAccount(user, username string) error {
	fmt.Printf("createNewUserAccount(%s, %s)\n", user, username)
	username = strings.ReplaceAll(username, " ", "_")
	err, out, errout := Shellout(fmt.Sprintf("pooltoy tx pooltoy create-user $(pooltoy keys show %s -a --keyring-backend test) false %s %s --from alice -y", user, username, user))
	fmt.Println("err", err)
	fmt.Println("out", out)
	fmt.Println("errout", errout)
	time.Sleep(5 * time.Second)
	return nil
}
func createNewUserKey(user, username string) error {
	fmt.Printf("createNewUserKey(%s, %s)\n", user, username)
	err, _, errout := Shellout(fmt.Sprintf("pooltoy keys add %s --keyring-backend test", user))

	if err == nil {
		path, err := os.Getwd()
		if err != nil {
			fmt.Println(err)
			return err
		}
		filename := fmt.Sprintf("%s/keys/%s.json", path, user)
		fmt.Printf("New user %s created and backup saved at %s\n", username, filename)
		d1 := []byte(errout)
		err = ioutil.WriteFile(filename, d1, 0644)
		if err != nil {
			fmt.Println(err)
		}
	}
	return createNewUserAccount(user, username)
}
func confirmUser(user, username string) error {
	fmt.Printf("confirmUser(%s, %s)\n", user, username)
	err, out, errout := Shellout(fmt.Sprintf("pooltoy keys show %s --keyring-backend test", user))
	fmt.Println("err", err)
	fmt.Println("out", out)
	fmt.Println("errout", errout)
	if err != nil {
		// there's an error, find out if it's just that the key doesn't exist
		if errout == "ERROR: The specified item could not be found in the keyring\n" {
			return createNewUserKey(user, username)
		} else {
			fmt.Printf("'%s' didn't match\n", errout)
			return err
		}
	} else {
		err, out, errout = Shellout(fmt.Sprintf("pooltoy q bank -o json balances $(pooltoy keys show %s -a --keyring-backend test)", user))
		fmt.Println("err", err)
		fmt.Println("out", out)
		fmt.Println("errout", errout)
		if err != nil && strings.Index(errout, "ERROR: unknown address: bank") != -1 {
			return createNewUserAccount(user, username)
		}
	}
	return nil
}

// slashes
func tilbrrr(userid string, text []string) string {
	// confirm sender user id key exists
	// if not create key
	// if not create account
	senderID, senderUsername, err := getUserID(userid)
	if err != nil {
		return fmt.Sprintf("ERROR: %s (%s)", err.Error(), senderID)
	}
	err = confirmUser(senderID, senderUsername)
	if err != nil {
		return fmt.Sprintf("ERROR: %s (%s)", err.Error(), userid)
	}

	if len(text) != 1 {
		return fmt.Sprintf("Sorry %s, I don't understand that command. Please follow the format '/til-brrr [user]'", senderUsername)
	}

	// confirm queried user id key exists
	// if not create key
	// if not create account

	extractedIDArray := strings.Split(text[0], "|")
	if len(extractedIDArray) < 2 {
		return fmt.Sprintf("%s does not follow the expected user id format (no |)", text[0])
	}
	if len(extractedIDArray[0]) < 3 {
		return fmt.Sprintf("%s does not follow the expected user id format (len < 3)", text[0])
	}
	extractedID := extractedIDArray[0][2:]

	queriedID, queriedUsername, err := getUserID(extractedID)
	if err != nil {
		return fmt.Sprintf("ERROR: %s (%s)", err.Error(), queriedID)
	}
	err = confirmUser(queriedID, queriedUsername)
	if err != nil {
		return fmt.Sprintf("ERROR: %s (%s)", err.Error(), userid)
	}

	command := fmt.Sprintf("pooltoy q faucet when-brrr -o json $(pooltoy keys show %s -a --keyring-backend test) | jq \".timeLeft\"", queriedID)
	fmt.Printf("Try command '%s\n", command)

	// create the CLI command for faucet from userid to queriedID
	err, out, errout := Shellout(command)

	fmt.Println("err", err)
	fmt.Println("out", out)
	fmt.Println("errout", errout)

	// parse various responses
	if err != nil {
		return err.Error()
	}
	fmt.Println(out)

	if out == "0" {
		return fmt.Sprintf("🖨 %s is ready to brrr right now!", queriedUsername)
	}

	_, err = strconv.Atoi(out)
	if err != nil {
		return err.Error()
		// fmt.Println(err)
	}
	// timeleft := time.Duration(int64(i)).String()

	timeleft, err := durafmt.ParseString(out + "s")
	if err != nil {
		fmt.Println(err)
	}

	return fmt.Sprintf("⏳ %s til %s can brrr again.", timeleft.String(), queriedUsername)
}
func brrr(userid string, text []string) string {

	if len(text) < 2 {
		return "Sorry, I don't understand that command. Please follow the format '/brrr [recipient] [emoji]' where emoji is part of the basic emoji list outlined here: https://unicode.org/Public/emoji/5.0/emoji-test.txt"
	}

	// confirm sender user id key exists
	// if not create key
	// if not create account
	senderID, senderUsername, err := getUserID(userid)
	if err != nil {
		return fmt.Sprintf("ERROR: %s (%s)", err.Error(), senderID)
	}
	err = confirmUser(senderID, senderUsername)
	if err != nil {
		return fmt.Sprintf("ERROR: %s (%s)", err.Error(), userid)
	}

	// confirm recipientID key exists
	// if not create key
	// if not create account

	recipientID, recipientUsername, err := getUserID(strings.Split(text[0], "|")[0][2:])
	if err != nil {
		return fmt.Sprintf("ERROR: %s (%s)", err.Error(), recipientID)
	}
	err = confirmUser(recipientID, recipientUsername)
	if err != nil {
		return fmt.Sprintf("ERROR: %s (%s)", err.Error(), userid)
	}
	emoji, containsEmoji := parseEmoji(text[1])

	// throw error if no emoji is found
	if !containsEmoji {
		emojiError := fmt.Errorf("No emoji found while parsing %s", text[1])
		return emojiError.Error()
	}

	command := fmt.Sprintf("pooltoy tx faucet mintfor $(pooltoy keys show %s -a --keyring-backend test) %s --from %s -y", recipientID, emoji, senderID)
	fmt.Printf("Try command '%s\n", command)

	// create the CLI command for faucet from userid to recipientID
	err, out, errout := Shellout(command)

	fmt.Println("err", err)
	fmt.Println("out", out)
	fmt.Println("errout", errout)

	// parse various responses
	if err != nil {
		return err.Error()
	}

	type TxResult struct {
		Height string
		Txhash string
		RawLog string
	}

	var txResult TxResult
	json.Unmarshal([]byte(out), &txResult)

	fmt.Println("txResult.Txhash", txResult.Txhash)
	// wait until the tx is processed
	time.Sleep(5 * time.Second)

	query := fmt.Sprintf("pooltoy q tx %s", txResult.Txhash)
	err, out, errout = Shellout(query)

	fmt.Println("err", err)
	fmt.Println("out", out)
	fmt.Println("errout", errout)

	var qResult map[string]interface{}
	json.Unmarshal([]byte(out), &qResult)

	fmt.Println("qResult", qResult)
	fmt.Println("qResult[\"codespace\"]", qResult["codespace"])

	// codespace is part of an error log
	if qResult["codespace"] != nil {
		return fmt.Sprintf("Sorry %s, you can only send an emoji once a day. Please try again tomorrow 📆", senderUsername)
	}

	return fmt.Sprintf("Success %s! You sent %s a %s. Check their balance like: /balance @%s", senderUsername, recipientUsername, emoji, recipientUsername)
}
func send(userid string, text []string) string {

	// confirm sender user id key exists
	// if not create key
	// if not create account
	senderID, senderUsername, err := getUserID(userid)
	if err != nil {
		return fmt.Sprintf("ERROR: %s (%s)", err.Error(), senderID)
	}
	err = confirmUser(senderID, senderUsername)
	if err != nil {
		return fmt.Sprintf("ERROR: %s (%s)", err.Error(), userid)
	}

	if len(text) < 2 {
		return fmt.Sprintf("Sorry %s, I don't understand that command. Please follow the format '/send [recipient] [emoji]' where emoji is part of the basic emoji list outlined here: https://unicode.org/Public/emoji/5.0/emoji-test.txt", senderUsername)
	}

	// confirm recipientID key exists
	// if not create key
	// if not create account

	recipientID, recipientUsername, err := getUserID(strings.Split(text[0], "|")[0][2:])
	if err != nil {
		return fmt.Sprintf("ERROR: %s (%s)", err.Error(), recipientID)
	}
	err = confirmUser(recipientID, recipientUsername)
	if err != nil {
		return fmt.Sprintf("ERROR: %s (%s)", err.Error(), userid)
	}

	emoji, containsEmoji := parseEmoji(text[1])
	// throw error if no emoji is found
	if !containsEmoji {
		emojiError := fmt.Errorf("No emoji found while parsing %s", text[1])
		return emojiError.Error()
	}

	command := fmt.Sprintf("pooltoy tx send %s $(pooltoy keys show %s -a --keyring-backend test) 1%s --from %s -y", senderID, recipientID, emoji, senderID)
	fmt.Printf("Try command '%s\n", command)

	// create the CLI command for faucet from userid to recipientID
	err, out, errout := Shellout(command)

	fmt.Println("err", err)
	fmt.Println("out", out)
	fmt.Println("errout", errout)

	// parse various responses
	if err != nil {
		return err.Error()
	}

	type TxResult struct {
		Height string
		Txhash string
		RawLog string
	}

	var txResult TxResult
	json.Unmarshal([]byte(out), &txResult)

	fmt.Println("txResult.Txhash", txResult.Txhash)
	// wait until the tx is processed
	time.Sleep(5 * time.Second)

	query := fmt.Sprintf("pooltoy q tx %s", txResult.Txhash)
	err, out, errout = Shellout(query)

	fmt.Println("err", err)
	fmt.Println("out", out)
	fmt.Println("errout", errout)

	var qResult map[string]interface{}
	json.Unmarshal([]byte(out), &qResult)

	fmt.Println("qResult", qResult)
	fmt.Println("qResult[\"codespace\"]", qResult["codespace"])

	// codespace is part of an error log
	if qResult["codespace"] != nil {
		wasInsufficient := strings.Index(qResult["raw_log"].(string), "insufficient funds") != -1
		if wasInsufficient {
			return fmt.Sprintf("Sorry %s you don't have enough %s to send any to %s. Try convincing one of your co-workers to /brrr you some 🖨", senderUsername, emoji, recipientUsername)
		}
		return fmt.Sprintf("Sorry %s, something went wrong\n%s", senderUsername, out)
	}

	return fmt.Sprintf("Success %s! You sent %s a %s. Check their balance like: /balance @%s", senderUsername, recipientUsername, emoji, recipientUsername)
}

func balance(userid string, text []string) string {

	// confirm sender user id key exists
	// if not create key
	// if not create account
	senderID, senderUsername, err := getUserID(userid)
	if err != nil {
		return fmt.Sprintf("ERROR: %s (%s)", err.Error(), senderID)
	}
	err = confirmUser(senderID, senderUsername)
	if err != nil {
		return fmt.Sprintf("ERROR: %s (%s)", err.Error(), userid)
	}

	if len(text) != 1 {
		return fmt.Sprintf("Sorry %s, I don't understand that command. Please follow the format '/balance [user]'", senderUsername)
	}

	// confirm sender user id key exists
	// if not create key
	// if not create account
	queriedID, queriedUsername, err := getUserID(strings.Split(text[0], "|")[0][2:])
	if err != nil {
		return fmt.Sprintf("ERROR: %s (%s)", err.Error(), queriedID)
	}
	err = confirmUser(queriedID, queriedUsername)
	if err != nil {
		return fmt.Sprintf("ERROR: %s (%s)", err.Error(), userid)
	}

	command := fmt.Sprintf("pooltoy q bank -o json balances $(pooltoy keys show %s -a --keyring-backend test) | jq \".balances\"", queriedID)
	fmt.Printf("Try command '%s\n", command)

	// create the CLI command for faucet from userid to queriedID
	err, out, errout := Shellout(command)

	fmt.Println("err", err)
	fmt.Println("out", out)
	fmt.Println("errout", errout)

	// parse various responses
	if err != nil {
		return err.Error()
	}

	type Coin struct {
		Denom  string
		Amount string
	}
	var coins []Coin
	json.Unmarshal([]byte(out), &coins)
	// coins := result["value"].(map[string]interface{})["coins"].([]Coin)

	fmt.Println("Coins", coins)
	balancetext := fmt.Sprintf("%s's balance:\n", queriedUsername)
	for i := 0; i < len(coins); i++ {
		balancetext += coins[i].Amount + " " + coins[i].Denom + "\n"
	}

	if len(coins) == 0 {
		balancetext = fmt.Sprintf("%s is broke 🕳", queriedUsername)
	}
	return balancetext
}

func main() {
	http.HandleFunc("/", botHandler)

	crt := os.Getenv("LETSENCRYPT_CRT")
	key := os.Getenv("LETSENCRYPT_KEY")
	if crt == "" || key == "" {
		log.Fatalln(http.ListenAndServe(":"+port, nil))
	} else {
		log.Fatalln(http.ListenAndServeTLS(":"+port, crt, key, nil))
	}

}
