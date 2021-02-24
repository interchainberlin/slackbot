package main

// Parse an emoji with optional modifier from a string
func parseEmoji(inputEmoji string) string {
	emoji := strings.TrimSpace(inputEmoji)
	// if slack emoji format
	if strings.Index(emoji, ":") == 0 {
		emoji = strings.ReplaceAll(inputEmoji, "️", "") // this removes Variation Selector-16 (https://emojipedia.org/variation-selector-16/)
	}

	//return emoji if it exists in default codemap
	if emojiCodeMap[emoji] != "" {
		emoji = emojiCodeMap[emoji]
		return emoji
	}

	//check for modifier if emoji not found in default codemap
	if emojiCodeMap[emoji] == "" {
		arr := strings.Split(emoji, "::")
		if len(arr) > 1 {
			//modifier will always be the final code
			modifier := ":" + arr[len(arr)-1]
			parsedEmoji := arr[0] + ":"

			_, hasEmoji := emojiCodeMap[parsedEmoji]
			_, hasModifier := emojiModifiers[modifier]

			if hasModifier && hasEmoji {
				return emojiCodeMap[parsedEmoji] + emojiModifiers[modifier]
			}
		}
	}

	return ""
}
