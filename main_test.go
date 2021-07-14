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
	fmt.Println("after init")
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
	case UserID4:
		return UserID4, UserName4, nil
	case UserID5:
		return UserID5, UserName5, nil
	case UserID6:
		return UserID6, UserName6, nil
	default:
		return "", "", errors.New(fmt.Sprintf("user id (%s) not found", userID))
	}
}

// -----------------------------------------------------------------------------
// test ConfirmUser
// -----------------------------------------------------------------------------
func TestConfirmUser_ExistentUser(t *testing.T) {
	go startChain()
	time.Sleep(20 * time.Second)

	err := confirmUser(UserID1, UserName1)
	if err != nil {
		t.Fatal(err)
	}
}

func TestConfirmUser_NonExistentUser(t *testing.T) {
	go startChain()
	time.Sleep(20 * time.Second)

	err := confirmUser(UserID4, UserName4)
	if errors.Is(err, errors.New(fmt.Sprintf("%s is not found", UserID4))){
		t.Fatal(err)
	}
}

// -----------------------------------------------------------------------------
// test createNewUser
// -----------------------------------------------------------------------------
// todo check why failed
//func TestCreateNewUserAccount(t *testing.T){
//	go startChain()
//	time.Sleep(20 * time.Second)
//
//	var userID ="U01JT5U5WL0"
//	var userName = "newuser"
//	err := confirmUser(userID, userName)
//	//if !errors.Is(err, errors.New(fmt.Sprintf("%s is not found", userID))){
//	//	t.Fatal(err)
//	//}
//	if err == nil{
//		t.Fatal(err)
//	}
//	err = createNewUserAccount(userID, userName)
//	if err!= nil{
//		t.Fatal(err)
//	}
//	err = confirmUser(userID, userName)
//	if errors.Is(err, errors.New(fmt.Sprintf("%s is not found", userID))){
//		fmt.Println("=======created user is not found==")
//		t.Fatal(err)
//	}
//}
//
//func TestCreateNewUserKey(t *testing.T){
//	go startChain()
//	time.Sleep(20 * time.Second)
//
//	var userID ="U01JT5U5WL0"
//	var userName = "newuser"
//	err := confirmUser(userID, userName)
//	//if !errors.Is(err, errors.New(fmt.Sprintf("%s is not found", userID))){
//	//	t.Fatal(err)
//	//}
//	if err == nil{
//		t.Fatal(err)
//	}
//	err = createNewUserKey(userID, userName)
//	if err!= nil{
//		t.Fatal(err)
//	}
//	err = confirmUser(userID, userName)
//	if errors.Is(err, errors.New(fmt.Sprintf("%s is not found", userID))){
//		fmt.Println("=======created user is not found==")
//		t.Fatal(err)
//	}
//	balance1 := balance(userID,[]string{"@<U01JT5U5WL0|newuser"} , GetUserID)
//	fmt.Println("balance of new user", balance1)
//}



// -----------------------------------------------------------------------------
// test Brrr
// -----------------------------------------------------------------------------
func TestBrrr(t *testing.T) {
	go startChain()
	time.Sleep(20 * time.Second)

	var testCommand1 = []string{User3, ":wave:"}
	brResp := brrr(UserID1, testCommand1, GetUserID)
	fmt.Println("brResp", brResp)
	var testCommand2 = []string{User3}
	balanceResp := balance(UserID1, testCommand2, GetUserID)
	if balanceResp != fmt.Sprintf("%s's balance:\n1 🆗\n2 🍍\n1 🎙️\n1 🎬\n1 👋\n1 🥅\n1 🥞\n", UserName3) {
		t.Fatal(balanceResp)
	}

	stopChain()
}

func TestBrrr_TwiceWithin1Day(t *testing.T) {
	go startChain()
	time.Sleep(20 * time.Second)

	var testCommand1 = []string{User3, ":wave:"}
	var testCommand2 = []string{User3}

	// first Brrr and check balance
	brResp := brrr(UserID1, testCommand1, GetUserID)
	fmt.Println("first brrr response", brResp)
	balanceResp := balance(UserID1, testCommand2, GetUserID)
	if balanceResp != fmt.Sprintf("%s's balance:\n1 🆗\n2 🍍\n1 🎙️\n1 🎬\n1 👋\n1 🥅\n1 🥞\n", UserName3) {
		t.Fatal(balanceResp)
	}

	// second Brrr
	brResp = brrr(UserID1, testCommand1, GetUserID)
	fmt.Println("second brrr response", brResp)
	if brResp != fmt.Sprintf("Sorry %s, you can only send an emoji once a day. Please try again tomorrow 📆", UserName1) {
		t.Fatal(brResp)
	}
	stopChain()
}

func TestTilBrrr(t *testing.T) {
	go startChain()
	time.Sleep(20 * time.Second)

	til := tilbrrr(UserID1, []string{User3}, GetUserID)
	if til != fmt.Sprintf("🖨 %s is ready to brrr right now!", UserName3) {
		t.Fatal(til)
	}

	brrr(UserID1, []string{User3, ":wave:"}, GetUserID)
	til = tilbrrr(UserID1, []string{User1}, GetUserID)
	fmt.Println(til)
	if til != fmt.Sprintf("⏳ 1 day til %s can brrr again.", UserName1) {
		t.Fatal(til)
	}

	stopChain()
}

// -----------------------------------------------------------------------------
// test Balance
// -----------------------------------------------------------------------------
func TestBalance(t *testing.T) {
	go startChain()
	time.Sleep(20 * time.Second)

	var testCommand = []string{User3}
	balanceResp := balance(UserID1, testCommand, GetUserID)
	if balanceResp != fmt.Sprintf("%s's balance:\n1 🆗\n2 🍍\n1 🎙️\n1 🎬\n1 🥅\n1 🥞\n", UserName3) {
		t.Fatal(balanceResp)
	}

	stopChain()
}

//  existent user checks nonexistent user's balance
func TestBalance_UserNotExist(t *testing.T) {
	go startChain()
	time.Sleep(20 * time.Second)

	var testCommand = []string{User4}
	balanceResp := balance(UserID1, testCommand, GetUserID)
	if balanceResp != fmt.Sprintf("ERROR: %s is not found (%s)", UserID4, UserID1) {
		t.Fatal(balanceResp)
	}

	stopChain()
}

//non-existent user checks nonexistent user's balance
func TestBalance_UserNotExist1(t *testing.T) {
	go startChain()
	time.Sleep(15 * time.Second)

	var testCommand = []string{User6}
	balanceResp := balance(UserID4, testCommand, GetUserID)
	if balanceResp != fmt.Sprintf("ERROR: %s is not found (%s)", UserID4, UserID4) {
		t.Fatal(balanceResp)
	}

	stopChain()
}

// -----------------------------------------------------------------------------
// test send
// -----------------------------------------------------------------------------
func TestSend(t *testing.T) {
	go startChain()
	time.Sleep(20 * time.Second)
	var testCommand1 = []string{User3, ":avocado:"}
	sendResp := send(UserID1, testCommand1, GetUserID)
	fmt.Println("sendResp:", sendResp)
	var testCommand2 = []string{User3}
	balances := balance(UserID1, testCommand2, GetUserID)
	if balances != fmt.Sprintf("%s's balance:\n1 🆗\n2 🍍\n1 🎙️\n1 🎬\n1 🥅\n1 🥑\n1 🥞\n", UserName3) {
		t.Fatal(balances)
	}

	stopChain()
}

func TestSend_OutOfBalance(t *testing.T) {
	go startChain()
	time.Sleep(20 * time.Second)

	var testCommand = []string{User1, "🍍"}
	sendResp := send(UserID3, testCommand, GetUserID)
	fmt.Println("sendResp:", sendResp)
	if sendResp != "No emoji found while parsing 🍍"{
		t.Fatal(sendResp)
	}

	stopChain()
}

//// can we send 2 ????
//func TestSend_MoreEmoji(t *testing.T) {
//	go startChain()
//	time.Sleep(25 * time.Second)
//
//	var testCommand = []string{User1, "2:birthday_cake:"}
//	sendResp := send(UserID5, testCommand, GetUserID)
//	fmt.Println("sendResp:", sendResp)
//	if sendResp != "No emoji found while parsing 2🎂"{
//		t.Fatal(sendResp)
//	}
//
//	stopChain()
//}

func TestSend_FromNonExistentUser(t *testing.T) {
	go startChain()
	time.Sleep(20 * time.Second)

	var testCommand = []string{User1, ":avocado:"}
	sendResp := send(UserID4, testCommand, GetUserID)
	fmt.Println("sendResp:", sendResp)
	if sendResp != fmt.Sprintf("ERROR: %s is not found (%s)",UserID4, UserID4) {
		t.Fatal(sendResp)
	}

	stopChain()
}

func TestSendTo_NonExistentUser(t *testing.T) {
	go startChain()
	time.Sleep(20 * time.Second)

	var testCommand = []string{User4, ":avocado:"}
	sendResp := send(UserID1, testCommand, GetUserID)
	fmt.Println("sendResp:", sendResp)
	if sendResp != fmt.Sprintf("ERROR: %s is not found (%s)",UserID4, UserID1) {
		t.Fatal(sendResp)
	}

	stopChain()
}


