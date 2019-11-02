package main

import (
	"fmt"
	"github.com/arkie/dicebot/roll"
	"github.com/arkie/dicebot/slack"
	"testing"
)

const (
	userID = "1234"
)

var formatTests = []string{
	"",
	"d6",
	"10d10-5",
	"123456789d123456789+123456789",
	"1d1 2d2 3d3",
	"10d10>5",
	"10d10<5",
	"10d10k5",
	"2d6 / 2d6",
	"mini 2d6 / 2d6",
}

func TestRollFormat(t *testing.T) {
	for _, test := range formatTests {
		rolls := roll.Parse(test)
		verify := fmt.Sprint("<@", userID, "> rolled ")
		sum := 0
		for i, result := range rolls {
			result.Roll()
			// Use the last roll as the final.
			total := result.Total
			switch result.Operator {
			case "+":
				sum += result.Total
			case "-":
				sum -= result.Total
			case "*":
				sum *= result.Total
			case "/":
				sum /= result.Total
			}
			if i == 0 {
				verify += fmt.Sprint("*", total, "*")
			} else {
				verify += fmt.Sprint(" ", result.Operator, " *", total, "*")
			}
		}
		if len(rolls) > 1 {
			verify += fmt.Sprint(" = *", sum, "*")
		}
		resp := formatRoll(userID, false, false, rolls)
		if resp["response_type"] != "in_channel" {
			t.Error("Incorrect response type", resp["response_type"])
		}
		attach := resp["attachments"].([]slack.D)[0]
		if attach["color"] == "" {
			t.Error("Missing color")
		}
		if attach["text"] != verify {
			t.Error("Incorrect response text", attach["text"], "instead of", verify)
		}
		resp = formatRoll(userID, false, true, rolls)
		if resp["response_type"] != "ephemeral" {
			t.Error("Incorrect response type", resp["response_type"])
		}
	}
}

// Test for panics.
func TestCommand(t *testing.T) {
	for _, test := range formatTests {
		resp := command(slack.Args{
			Text:     test,
			UserName: "CommandTest",
		})
		if resp["response_type"] != "in_channel" {
			t.Error("Incorrect response type", resp["response_type"])
		}
	}
}

func TestFor(t *testing.T) {
	rolls := roll.Parse("d20 for initiative, 2d6 + 5 for attack")
	resp := formatRoll(userID, false, false, rolls)
	attach := resp["attachments"].([]slack.D)[0]
	verify := "<@1234> rolled 0 for initiative + 0 + 0 = 0 for attack = 0"
	if attach["fallback"] != verify {
		t.Error("Incorrect response text", attach["fallback"], "instead of", verify)
	}
}
