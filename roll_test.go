package hackyslack

import (
	"fmt"
	"testing"
)

var (
	parseTests = []RollTest{
		RollTest{"", []DiceRoll{DiceRoll{Number: 2, Sides: 6}}},
		RollTest{"blah", []DiceRoll{DiceRoll{Number: 2, Sides: 6}}},
		RollTest{"2", []DiceRoll{DiceRoll{Number: 2, Sides: 6}}},
		RollTest{"d", []DiceRoll{DiceRoll{Number: 1, Sides: 6}}},
		RollTest{"d%", []DiceRoll{DiceRoll{Number: 1, Sides: 100}}},
		RollTest{"0d1", []DiceRoll{DiceRoll{Number: 1, Sides: 1}}},
		RollTest{"1d20-1", []DiceRoll{DiceRoll{Number: 1, Sides: 20, Modifier: -1}}},
		RollTest{"2d20+12345", []DiceRoll{DiceRoll{Number: 2, Sides: 20, Modifier: 1234}}},
		RollTest{"1234567890d123456790-1234567890", []DiceRoll{DiceRoll{Number: 100, Sides: 1000}}},
		RollTest{"2d2+1 1d6", []DiceRoll{
			DiceRoll{Number: 2, Sides: 2, Modifier: 1},
			DiceRoll{Number: 1, Sides: 6},
		}},
		RollTest{"1d20, 2d6-10", []DiceRoll{
			DiceRoll{Number: 1, Sides: 20},
			DiceRoll{Number: 2, Sides: 6, Modifier: -10},
		}},
		RollTest{"2d6>5", []DiceRoll{DiceRoll{Number: 2, Sides: 6, Minimum: 5}}},
		RollTest{"2d6<2", []DiceRoll{DiceRoll{Number: 2, Sides: 6, Maximum: 2}}},
		RollTest{"6d6k5", []DiceRoll{DiceRoll{Number: 6, Sides: 6, Keep: 5}}},
		RollTest{"2d6k5", []DiceRoll{DiceRoll{Number: 2, Sides: 6, Keep: 2}}},
	}
)

type RollTest struct {
	Text  string
	Rolls []DiceRoll
}

func compareRolls(a DiceRoll, b DiceRoll) bool {
	return a.Number == b.Number && a.Sides == b.Sides &&
		a.Modifier == b.Modifier && a.Keep == b.Keep &&
		a.Minimum == b.Minimum && a.Maximum == b.Maximum
}

func TestRoll(t *testing.T) {
	username := "TestUser"
	for _, test := range parseTests {
		rolls := parseRoll(test.Text)
		verify := fmt.Sprint("@", username, " rolled ")
		sum := 0
		for i, result := range rolls {
			if !compareRolls(*result, test.Rolls[i]) {
				t.Error("Failed", test, "got", *result)
			}
			for i := 0; i < 10; i++ {
				result.Roll()
				if result.Total < result.Number {
					t.Error("Rolled too low", *result)
				}
				if result.Total > result.Number*result.Sides {
					t.Error("Rolled too high", *result)
				}
			}
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
		attach := resp["attachments"].([]D)[0]
		if attach["color"] == "" {
			t.Error("Missing color")
		}
		if attach["text"] != verify {
			t.Log(rolls)
			t.Error("Incorrect response text", attach["text"], "instead of", verify)
		}
	}
}
