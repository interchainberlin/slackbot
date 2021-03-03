package main

import (
	"testing"
)

func TestEmojiParserCorrectSingleEmoji(t *testing.T) {
	_, contains := parseEmoji(":female-artist:")
	if !contains {
		t.Error("Does not contain emoji")
	}
}

func TestEmojiParserCorrectDoubleEmojiWithModifier(t *testing.T) {
	_, contains := parseEmoji(":female-artist::skin-tone-3:")
	if !contains {
		t.Error("Does not contain emoji")
	}
}

func TestEmojiParserIncorrectDoubleEmojiWithModifier(t *testing.T) {
	_, contains := parseEmoji(":dinosour-artist::skin-tone-3:")
	if contains {
		t.Error("Should not contain emoji code")
	}
}

type testSyntax struct {
	arg1     string
	expected bool
}

var syntaxTests = []testSyntax{
	{":female-artist:", true},
	{"female-artist:", false},
	{"::female-artist:", false},
	{":female-artist:skin-tone-2:", false},
	{":female-artist::skin-tone-2:", true},
}

func TestEmojiParserInputSyntax(t *testing.T) {
	for _, test := range syntaxTests {
		if _, contains := parseEmoji(test.arg1); contains != test.expected {
			t.Error("Incorect syntax")
		}
	}
}
