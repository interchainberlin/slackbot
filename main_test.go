package main

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hako/durafmt"

	"github.com/stretchr/testify/require"
)

func TestTime(t *testing.T) {
	out := "\"6400\""
	out = strings.ReplaceAll(out, "\"", "")
	// i, err := strconv.Atoi(out)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// timeleft := time.Duration(int64(i)).String()

	timeleft, err := durafmt.ParseString(out + "s")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("timeleft", timeleft)
	require.True(t, false)

}

func TestConfirmUser(t *testing.T) {
	user := "doug"
	emoji := ":boat:"
	recipientID := "recipientID"
	senderID := "senderID"
	fmt.Printf("emoji: '%s'\n", emoji)
	fmt.Println("emojiCodeMap", emojiCodeMap[emoji])

	if emojiCodeMap[emoji] != "" {
		emoji = emojiCodeMap[emoji]
	}
	command := fmt.Sprintf("pooltoycli tx faucet mintfor $(pooltoycli keys show %s -a) %s --from %s -y", recipientID, emoji, senderID)
	fmt.Printf("Try command '%s\n", command)

	err, out, errout := Shellout(fmt.Sprintf("pooltoycli q account  $(pooltoycli keys show %s  -a)", user))
	fmt.Println("err", err)
	fmt.Println("out", out)
	fmt.Println("errout", errout)

	require.True(t, false)
}
