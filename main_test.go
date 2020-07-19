package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConfirmUser(t *testing.T) {
	user := "doug"
	emoji := ":boat:️"
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
