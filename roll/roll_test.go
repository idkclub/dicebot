package roll

import (
	"testing"
)

var ParseTests = []RollTest{
	RollTest{"", []Dice{Dice{Number: 2, Sides: 6}}},
	RollTest{"blah", []Dice{Dice{Number: 2, Sides: 6}}},
	RollTest{"2", []Dice{Dice{Number: 2, Sides: 6}}},
	RollTest{"0d0", []Dice{Dice{Number: 1, Sides: 1}}},
	RollTest{"d", []Dice{Dice{Number: 1, Sides: 6}}},
	RollTest{"d%", []Dice{Dice{Number: 1, Sides: 100}}},
	RollTest{"0d1", []Dice{Dice{Number: 1, Sides: 1}}},
	RollTest{"1d20-1", []Dice{Dice{Number: 1, Sides: 20, Modifier: -1}}},
	RollTest{"2d20+12345", []Dice{Dice{Number: 2, Sides: 20, Modifier: 1234}}},
	RollTest{"1234567890d123456790-1234567890", []Dice{Dice{Number: 100, Sides: 1000}}},
	RollTest{"2d2+1 1d6", []Dice{
		Dice{Number: 2, Sides: 2, Modifier: 1},
		Dice{Number: 1, Sides: 6},
	}},
	RollTest{"1d20, 2d6-10", []Dice{
		Dice{Number: 1, Sides: 20},
		Dice{Number: 2, Sides: 6, Modifier: -10},
	}},
	RollTest{"1d1+1 2d2-2 3d3+3", []Dice{
		Dice{Number: 1, Sides: 1, Modifier: 1},
		Dice{Number: 2, Sides: 2, Modifier: -2},
		Dice{Number: 3, Sides: 3, Modifier: 3},
	}},
	RollTest{"2d6>5", []Dice{Dice{Number: 2, Sides: 6, Minimum: 5}}},
	RollTest{"2d6<2", []Dice{Dice{Number: 2, Sides: 6, Maximum: 2}}},
	RollTest{"2d6>6", []Dice{Dice{Number: 2, Sides: 6, Minimum: 5}}},
	RollTest{"2d6<1", []Dice{Dice{Number: 2, Sides: 6, Maximum: 2}}},
	RollTest{"6d6k5", []Dice{Dice{Number: 6, Sides: 6, Keep: 5}}},
	RollTest{"2d6k5", []Dice{Dice{Number: 2, Sides: 6, Keep: 2}}},
}

type RollTest struct {
	Text  string
	Rolls []Dice
}

func compareRolls(a Dice, b Dice) bool {
	return a.Number == b.Number && a.Sides == b.Sides &&
		a.Modifier == b.Modifier && a.Keep == b.Keep &&
		a.Minimum == b.Minimum && a.Maximum == b.Maximum
}

func TestRoll(t *testing.T) {
	for _, test := range ParseTests {
		rolls := Parse(test.Text)
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
		}
	}
}
