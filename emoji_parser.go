package main

import "strings"

// Parse an emoji with optional modifier from a string
func parseEmoji(inputEmoji string) (string, bool) {
	emoji := strings.TrimSpace(inputEmoji)
	// if slack emoji format
	if strings.Index(emoji, ":") == 0 {
		emoji = strings.ReplaceAll(inputEmoji, "️", "") // this removes Variation Selector-16 (https://emojipedia.org/variation-selector-16/)
	}

	//return emoji if it exists in default codemap
	if emojiCodeMap[emoji] != "" {
		emoji = emojiCodeMap[emoji]
		return emoji, true
	}

	//check for modifier if emoji not found in default codemap
	arr := strings.Split(emoji, "::")
	if len(arr) > 1 {
		//modifier will always be the final code
		modifier := ":" + arr[len(arr)-1]
		parsedEmoji := strings.ReplaceAll(emoji, modifier, "")

		_, hasEmoji := emojiCodeMap[parsedEmoji]
		_, hasModifier := emojiModifierCodeMap[modifier]

		if hasModifier && hasEmoji {
			return emojiCodeMap[parsedEmoji] + emojiModifierCodeMap[modifier], true
		}
	}

	return "", false
}
