package dicebot

import (
	"fmt"
	"github.com/arkie/hackyslack2"
	"github.com/arkie/hackyslack2/dicebot/roll"
	"testing"
)

var FormatTests = []string{
	"",
	"d6",
	"10d10-5",
	"123456789d123456789+123456789",
	"1d1 2d2 3d3",
	"10d10>5",
	"10d10<5",
	"10d10k5",
}

func TestRollFormat(t *testing.T) {
	username := "TestUser"
	for _, test := range FormatTests {
		rolls := roll.Parse(test)
		verify := fmt.Sprint("@", username, " rolled ")
		sum := 0
		for i, result := range rolls {
			result.Roll()
			// Use the last roll as the final.
			total := result.Total + result.Modifier
			sum += total
			if i == 0 {
				verify += fmt.Sprint("*", total, "*")
			} else {
				verify += fmt.Sprint(" + *", total, "*")
			}
		}
		if len(rolls) > 1 {
			verify += fmt.Sprint(" = *", sum, "*")
		}
		resp := formatRoll(username, rolls)
		if resp["response_type"] != "in_channel" {
			t.Error("Incorrect response type", resp["response_type"])
		}
		attach := resp["attachments"].([]hackyslack.D)[0]
		if attach["color"] == "" {
			t.Error("Missing color")
		}
		if attach["text"] != verify {
			t.Log(rolls)
			t.Error("Incorrect response text", attach["text"], "instead of", verify)
		}
	}
}
