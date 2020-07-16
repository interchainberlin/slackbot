package main

import (
	"fmt"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConfirmUser(t *testing.T) {
	user := "asdf"

	app := "pooltoycli"
	arg0 := "keys"
	arg1 := "show"
	arg2 := user
	// arg3 := ""

	cmd := exec.Command(app, arg0, arg1, arg2)
	fmt.Println("cmd", cmd)

	stdout, err := cmd.Output()
	fmt.Println("stdout", stdout)
	fmt.Println("err", err)
	require.True(t, false)
}
