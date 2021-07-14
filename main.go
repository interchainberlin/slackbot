// https://github.com/slack-go/slack
// https://guzalexander.com/2017/09/15/cowsay-slack-command.html

package main

import (
	"bytes"
	"encoding/json"
	"errors"
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

//
//type TxResult struct {
//	Height string
//	Txhash string
//	RawLog string
//}
type TxResult struct {
	Code       int
	Codespace  string
	Data       string
	Gas_used   string
	Gas_wanted string
	Height     string
	Info       string
	Logs       string
	Raw_log    string
	Timestamp  string
	Tx         string
	Txhash     string
}

type EventAttribute struct {
	Key   []byte
	Value []byte
	Index bool
}

type Event struct {
	Type       string
	Attributes []EventAttribute
}

type ResponseCheckTx struct {
	Code      uint32
	Data      []byte
	Log       string
	Info      string
	GasWanted int64
	GasUsed   int64
	Events    []Event
	Codespace string
}

type UserGetter func(string) (string, string, error)

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

	err := r.ParseForm()
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

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

	go handleCommand(responseURL, command, userid, textArray, getUserID)
}

func handleCommand(responseURL, command, userid string, textArray []string, getUserID UserGetter) {
	var botReply string
	switch command {
	case "/brrr":
		botReply = brrr(userid, textArray, getUserID)
	case "/send":
		botReply = send(userid, textArray, getUserID)
	case "/balance":
		botReply = balance(userid, textArray, getUserID)
	case "/til-brrr":
		botReply = tilbrrr(userid, textArray, getUserID)
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

var getUserID = func(userID string) (string, string, error) {
	fmt.Println("THIS IS THE REAL GETUSERID AND SHOULD NOT RUN IN TESTS")
	api := slack.New(os.Getenv("API_TOKEN"))

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

<<<<<<< HEAD
// Aadd a slack user to the pooltoy blockchain using the admin user alice and giving the new user non-admin permissions
// wait 5 seconds so that the transactions is completed
// return an error if the transaction failed
=======
// working directory is not the directory keeps the key, the~/.pooltoy/keyring_test is? however, the key info is stored in file with suffix .info rathter than .json
>>>>>>> 199f573 (test edge cases and change confirmUser function)
func createNewUserAccount(user, username string) error {
	fmt.Printf("createNewUserAccount(%s, %s)\n", user, username)
	username = strings.ReplaceAll(username, " ", "_")
	err, out, errout := Shellout(fmt.Sprintf("pooltoy tx pooltoy create-user $(pooltoy keys show %s -a --keyring-backend test) false %s %s --from alice -y --keyring-backend test --chain-id pooltoy-5", user, username, user))
	fmt.Println("err", err)
	fmt.Println("out", out)
	fmt.Println("errout", errout)
	time.Sleep(5 * time.Second)
	// TODO: check the tx hash to see whether the transaction was successful or not.
	// return the error if the transaction failed
	return nil
}

<<<<<<< HEAD
// Add a local key representing this slack user.
// Return an error if this fails
// Do not return an error, if the user key already exists
=======
>>>>>>> 199f573 (test edge cases and change confirmUser function)
func createNewUserKey(user, username string) error {
	fmt.Printf("createNewUserKey(%s, %s)\n", user, username)
	err, out, errout := Shellout(fmt.Sprintf("pooltoy keys add %s --keyring-backend test", user))
	fmt.Println("err", err)
	fmt.Println("out", out)
	fmt.Println("errout", errout)
	// TODO: figure out if errout contains actual errors.
	return err
	// 	this is currently broken vvvvvvvv
	// 	path, err := os.Getwd()
	// 	if err != nil {
	// 		fmt.Println(err)
	// 		return err
	// 	}
	// 	filename := fmt.Sprintf("%s/keys/%s.json", path, user)
	// 	fmt.Printf("New user %s created and backup saved at %s\n", username, filename)
	// 	d1 := []byte(errout)
	// 	err = ioutil.WriteFile(filename, d1, 0644)
	// 	if err != nil {
	// 		fmt.Println(err)
	// 	}
}

// Checks if a slack user key and a slack user account exists. If not create that key and/or account.
// return an error when the creation of the key or the account fails.
func confirmUser(user, username string) error {
	fmt.Printf("confirmUser(%s, %s)\n", user, username)
	err, out, errout := Shellout(fmt.Sprintf("pooltoy keys show %s --keyring-backend test", user))
	if err!= nil{
		fmt.Println("confirmUser err", err)
		fmt.Println("confirmUser out", out)
		fmt.Println("confirmUser errout", errout)
	}

	if err != nil {
<<<<<<< HEAD
		// there's an error, find out if it's just that the key doesn't exist
		if user != "" {
			err = createNewUserKey(user, username)
			if err != nil {
				return err
			}
		} else {
			fmt.Printf("'%s' didn't match\n", errout)
			return err
		}
	}
	err, out, errout = Shellout(fmt.Sprintf("pooltoy q auth account $(pooltoy keys show %s -a --keyring-backend test) -o json", user))
	fmt.Println("err", err)
	fmt.Println("out", out)
	fmt.Println("errout", errout)
	if err != nil && strings.Index(errout, "Error: rpc error: code = NotFound") != -1 {
		return createNewUserAccount(user, username)
=======

			return errors.New(fmt.Sprintf("%s is not found", user))

	} else {
		err, out, errout = Shellout(fmt.Sprintf("pooltoy q bank -o json balances $(pooltoy keys show %s -a --keyring-backend test) | jq \".balances\"", user))
		fmt.Println(fmt.Sprintf("pooltoy q bank -o json balances $(pooltoy keys show %s -a --keyring-backend test) | jq \".balances\"", user))
		if err!= nil{
			fmt.Println("err", err)
			fmt.Println("out", out)
			fmt.Println("confirmUser errout", errout)
		}
		if err != nil && strings.Index(errout, "ERROR: unknown address: bank") != -1 {
			return createNewUserAccount(user, username)
		}
>>>>>>> 199f573 (test edge cases and change confirmUser function)
	}
	return nil
}

func checkTimeLeft(queriedID string) string {
	command := fmt.Sprintf("pooltoy q faucet when-brrr -o json $(pooltoy keys show %s -a --keyring-backend test)", queriedID)
	fmt.Printf("Try command '%s\n", command)

	// create the CLI command for faucet from userid to queriedID
	err, out, errout := Shellout(command)
	if err != nil {
		fmt.Println("err", err)
		fmt.Println("out", out)
		fmt.Println("errout", errout)
	}


	// parse various responses
	if err != nil {
		return err.Error()
	}

	type TimeLeft struct {
		TimeLeft string
	}
	time := TimeLeft{}
	json.Unmarshal([]byte(out), &time)
	return strings.ReplaceAll(time.TimeLeft, "\"", "")
}

// from send result
func parseTxResult(out string) (map[string]string, error) {
	fields := []string{
		"code",
		"codespace",
		"data",
		"gas_used",
		"gas_wanted",
		"height",
		"info",
		"logs",
		"raw_log",
		"timestamp",
		"tx",
		"txhash",
	}
	tx := make(map[string]string, len(fields))
	s := strings.Fields("\n" + out)

	for i := 0; i < len(fields); i++ {
		if fields[i] != s[i*2][:len(s[i*2])-1] {
			return nil, errors.New(fmt.Sprintf("field %s is not found, %s is found", fields[i], s[i*2]))
		}

		tx[fields[i]] = s[i*2+1]
	}

	return tx, nil
}

// temporal solution for the cosmos sdk issue #9663, output format issue
func txErr(out string) bool {
	txMap, err := parseTxResult(out)
	if err != nil {
		fmt.Println("failed to parse tx out")
	}

	time.Sleep(5 * time.Second)

	query := fmt.Sprintf("pooltoy q tx %s", txMap["txhash"])
	err, out1, errout := Shellout(query)
	if err!= nil{
		fmt.Println("err", err)
		fmt.Println("out", out1)
		fmt.Println("txErr errout", errout)
	}

	c, _ := strconv.Atoi(txMap["code"])
	fmt.Println("[\"code\"]", c)

	if c != 0 {
		return true
	}
	return false
}

// slashes
func tilbrrr(userid string, text []string, getUserID UserGetter) string {
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

	if len(text) != 1 || text[0] == "" {
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

	timeLeft := checkTimeLeft(queriedID)

	if timeLeft == "0" {
		return fmt.Sprintf("🖨 %s is ready to brrr right now!", queriedUsername)
	}

	_, err = strconv.Atoi(timeLeft)
	if err != nil {
		return err.Error()
	}

	t, err := durafmt.ParseString(timeLeft + "s")
	if err != nil {
		fmt.Println(err)
	}

	return fmt.Sprintf("⏳ %s til %s can brrr again.", t.String(), queriedUsername)
}
func brrr(userid string, text []string, getUserID UserGetter) string {

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

	if recipientID == senderID {
		return fmt.Sprintf("Nice try %s, but you gotta work for your emojis!", senderUsername)
	}

	emoji, containsEmoji := parseEmoji(text[1])

	// throw error if no emoji is found
	if !containsEmoji {
		emojiError := fmt.Errorf("No emoji found while parsing %s", text[1])
		return emojiError.Error()
	}

	timeLeft := checkTimeLeft(senderID)
	if timeLeft != "0" {
		return fmt.Sprintf("Sorry %s, you can only send an emoji once a day. Please try again tomorrow 📆", senderUsername)
	}

	command := fmt.Sprintf("pooltoy tx faucet mintfor $(pooltoy keys show %s -a --keyring-backend test) %s --from %s -y --keyring-backend test --chain-id pooltoy-5", recipientID, emoji, senderID)
	fmt.Printf("Try command '%s\n", command)

	// create the CLI command for faucet from userid to recipientID
	err, out, errout := Shellout(command)
	if err!= nil{
		fmt.Println("err", err)
		fmt.Println("out", out)
		fmt.Println("errout", errout)
	}

	// parse various responses
	if err != nil {
		return err.Error()
	}

	if txErr(out) {
		return fmt.Sprintf("There has been an error: %s", err)
	}

	return fmt.Sprintf("Success %s! You sent %s a %s. Check their balance like: /balance @%s", senderUsername, recipientUsername, emoji, recipientUsername)
}
func send(userid string, text []string, getUserID UserGetter) string {

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

	command := fmt.Sprintf("pooltoy tx bank send %s $(pooltoy keys show %s -a --keyring-backend test) 1%s --from %s -y --keyring-backend test --chain-id pooltoy-5", senderID, recipientID, emoji, senderID)
	fmt.Printf("Try command '%s\n", command)

	// create the CLI command for faucet from userid to recipientID
	err, out, errout := Shellout(command)
	if err!= nil {
		fmt.Println("err", err)
		fmt.Println("send out", out)
		fmt.Println("errout", errout)
	}
	// parse various responses
	if err != nil {
		return err.Error()
	}

	// TODO: add logging about insufficient funds
	// if wasInsufficient {
	// 	return fmt.Sprintf("Sorry %s you don't have enough %s to send any to %s. Try convincing one of your co-workers to /brrr you some 🖨", senderUsername, emoji, recipientUsername)
	// }

	if txErr(out) {
		return fmt.Sprintf("Sorry %s, something went wrong\n%s", senderUsername, out)
	}

	return fmt.Sprintf("Success %s! You sent %s a %s. Check their balance like: /balance @%s", senderUsername, recipientUsername, emoji, recipientUsername)
}

func balance(userid string, text []string, getUserID UserGetter) string {

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

	if len(text) != 1 || text[0] == "" {
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
	if err!= nil{
		fmt.Println("err", err)
		fmt.Println("out", out)
		fmt.Println("errout", errout)
	}

	// parse various responses
	if err != nil {
		return err.Error()
	}

	type Coin struct {
		Denom  string
		Amount string
	}
	var coins []Coin
	err = json.Unmarshal([]byte(out), &coins)
	if err != nil {
		fmt.Println("unmarshal error", err)
	}

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


