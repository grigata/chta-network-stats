package parser

import (
	"encoding/hex"
	"strings"
	"unicode"
)

func DecodeCoinbaseText(coinbaseHex string) string {
	data, err := hex.DecodeString(coinbaseHex)
	if err != nil {
		return ""
	}

	var builder strings.Builder

	for _, r := range string(data) {
		if unicode.IsPrint(r) {
			builder.WriteRune(r)
		} else {
			builder.WriteByte(' ')
		}
	}

	return strings.Join(strings.Fields(builder.String()), " ")
}

func DetectPool(text string) string {
	lower := strings.ToLower(text)

	switch {
	case strings.Contains(lower, "tinywinypool"):
		return "TinyWinyPool"

	case strings.Contains(lower, "chta-pool"):
		return "TinyWinyPool"

	case strings.Contains(lower, "ckpool-lhr-chta"):
		return "HeliosPool"

	case strings.Contains(lower, "mined by satoshi"):
		return "Satoshi Pool"

	case strings.Contains(lower, "bowserlab"):
		return "Bowserlab"

	case strings.Contains(lower, "rt-pool.cc"):
		return "RT-Pool"

	default:
		return "Unknown"
	}
}
