package main

import (
	"errors"
	"fmt"
	"testing"
	"time"

	//"time"
)

//func setup() {
//	os.Setenv("VERIFICATION_TOKEN", "foobar")
//	startChain()
//}

//func TestMain(m *testing.M) {
//	go startChain()
//	time.Sleep(20 * time.Second)
//	code := m.Run()
//   stopChain()
//	os.Exit(code)
//}

//func TestTime(t *testing.T) {
//	out := "\"6400\""
//	out = strings.ReplaceAll(out, "\"", "")
//	//i, err := strconv.Atoi(out)
//	//if err != nil {
//	//	fmt.Println(err)
//	//}
//	//timeleft := time.Duration(int64(i)).String()
//
//	timeleft, err := durafmt.ParseString(out + "s")
//	if err != nil {
//		fmt.Println(err)
//	}
//	fmt.Println("timeleft", timeleft)
//	require.True(t, false)
//}
//
//func TestConfirmUser(t *testing.T) {
//	user := "doug"
//	emoji := ":boat:"
//	recipientID := "recipientID"
//	senderID := "senderID"
//	fmt.Printf("emoji: '%s'\n", emoji)
//	fmt.Println("emojiCodeMap", emojiCodeMap[emoji])
//
//	if emojiCodeMap[emoji] != "" {
//		emoji = emojiCodeMap[emoji]
//	}
//	command := fmt.Sprintf("pooltoy tx faucet mintfor $(pooltoy keys show %s -a) %s --from %s -y", recipientID, emoji, senderID)
//	fmt.Printf("Try command '%s\n", command)
//
//	err, out, errout := Shellout(fmt.Sprintf("pooltoy q account  $(pooltoy keys show %s  -a)", user))
//	fmt.Println("err", err)
//	fmt.Println("out", out)
//	fmt.Println("errout", errout)
//
//	require.True(t, false)
//}

func startChain() {
	stopChain()
	err, out, errout := Shellout(`./init.sh`)
	if err != nil {
		fmt.Println("err", err)
		fmt.Println("out", out)
		fmt.Println("errout", errout)
	}
	Shellout(`pooltoy start`)
}

func stopChain() {
	Shellout(`killall -9 pooltoy`)
	fmt.Println("pooltoy stoped")
}

var GetUserID = func(userID string) (string, string, error) {
	fmt.Println("SWITCH:", userID)
	switch userID {
	case UserID1:
		return UserID1, UserName1, nil
	case UserID2:
		return UserID2, UserName2, nil
	case UserID3:
		return UserID3, UserName3, nil
	default:
		return "", "", errors.New(fmt.Sprintf("user id (%s) not found", userID))
	}
}

func TestConfirmUser(t *testing.T){
	go startChain()
	time.Sleep(20 * time.Second)

	err := confirmUser(UserID1, UserName1)
	if err != nil {
		t.Fatal(err)
	}
}

func TestBrrr(t *testing.T) {
	go startChain()
	time.Sleep(20 * time.Second)

	var testCommand1 = []string{User3, ":wave:"}
	brResp := brrr(UserID1, testCommand1, GetUserID)
	fmt.Println("brResp", brResp)
	var testCommand2 = []string{User3}
	balanceResp := balance(UserID1, testCommand2, GetUserID)
	if balanceResp != "Onur's balance:\n1 🆗\n2 🍍\n1 🎙️\n1 🎬\n1 👋\n1 🥅\n1 🥞\n"{
		t.Fatal(balanceResp)
	}
	//not sure stopchain is needed
	stopChain()
}

func TestBalance(t *testing.T) {
	go startChain()
	time.Sleep(20 * time.Second)

	var testCommand = []string{User3}
	balanceResp := balance(UserID1, testCommand, GetUserID)
	if balanceResp != "Onur's balance:\n1 🆗\n2 🍍\n1 🎙️\n1 🎬\n1 🥅\n1 🥞\n" {
		t.Fatal(balanceResp)
	}

	stopChain()
}

func TestSend(t *testing.T) {
	go startChain()
	time.Sleep(20 * time.Second)
	var testCommand1 = []string{User3, ":avocado:"}
	sendResp := send(UserID1, testCommand1, GetUserID)
	fmt.Println("sendResp:", sendResp)
	var testCommand2 = []string{User3}
	balances := balance(UserID1, testCommand2, GetUserID)
	if balances != "Onur's balance:\n1 🆗\n2 🍍\n1 🎙️\n1 🎬\n1 🥅\n1 🥑\n1 🥞\n"{
		t.Fatal(balances)
	}

	stopChain()
}

func TestTilBrrr(t *testing.T){
	go startChain()
	time.Sleep(20 * time.Second)

	til :=tilbrrr(UserID1,[]string{User3}, GetUserID)
	if til != "🖨 Onur is ready to brrr right now!"{
		t.Fatal(til)
	}

	brrr(UserID1, []string{User3, ":wave:"}, GetUserID)
	til =tilbrrr(UserID1,[]string{User1}, GetUserID)
	fmt.Println(til)
	if til != "⏳ 1 day til billy can brrr again."{
		t.Fatal(til)
	}

	stopChain()
}

