package roll

import (
	"math/rand"
	"regexp"
	"sort"
	"strconv"
)

var regex = regexp.MustCompile(`(?i)(?P<op>[×*/+-])?\s*((?P<num>\d{0,3})d(?P<sides>%|\d{1,4})(?P<max>[<>]\d{1,4})?(?P<keep>k\d{1,3})?|(?P<const>\d{1,5}))`)

const (
	Add      = "+"
	Subtract = "-"
	Multiply = "*"
	Divide   = "/"
)

type Dice struct {
	Operator string
	Number   int
	Sides    int
	Minimum  int
	Maximum  int
	Keep     int
	Rolls    []int
	Removed  []int
	Total    int
}

func Parse(text string) []*Dice {
	var rolls []*Dice
	for _, m := range regex.FindAllStringSubmatch(text, 5) {
		dice := &Dice{
			Operator: Add,
			Number:   2,
			Sides:    6,
		}
		for i, name := range regex.SubexpNames() {
			switch name {
			case "op":
				if m[i] == "×" {
					dice.Operator = "*"
				} else if m[i] != "" {
					dice.Operator = m[i]
				}
			case "const":
				if m[i] != "" {
					num, _ := strconv.Atoi(m[i])
					dice.Number = num
					dice.Sides = 1
					break
				}
			case "num":
				num, _ := strconv.Atoi(m[i])
				if num < 1 {
					num = 1
				}
				if num > 100 {
					num = 100
				}
				dice.Number = num
			case "sides":
				if m[i] == "" {
					dice.Sides = 6
				} else if m[i] == "%" {
					dice.Sides = 100
				} else {
					dice.Sides, _ = strconv.Atoi(m[i])
					if dice.Sides < 1 {
						dice.Sides = 1
					}
					if dice.Sides > 1000 {
						dice.Sides = 1000
					}
				}
			case "keep":
				if m[i] == "" {
					break
				}
				dice.Keep, _ = strconv.Atoi(m[i][1:])
				if dice.Keep > dice.Number {
					dice.Keep = dice.Number
				}
			case "max":
				if m[i] == "" {
					break
				}
				if m[i][0] == '>' {
					dice.Minimum, _ = strconv.Atoi(m[i][1:])
					if dice.Minimum >= dice.Sides {
						dice.Minimum = dice.Sides - 1
					}
				} else {
					dice.Maximum, _ = strconv.Atoi(m[i][1:])
					if dice.Maximum < 2 {
						dice.Maximum = 2
					}
				}
			}
		}
		rolls = append(rolls, dice)
	}
	if len(rolls) == 0 {
		rolls = append(rolls, &Dice{
			Operator: Add,
			Number:   2,
			Sides:    6,
		})
	}
	return rolls
}

func (r *Dice) Roll() {
	r.Total = 0
	if r.Sides == 1 {
		r.Total = r.Number
		return
	}
	for i := 0; i < r.Number; i++ {
		n := rand.Intn(r.Sides) + 1
		r.Total += n
		r.Rolls = append(r.Rolls, n)
	}
	if r.Keep != 0 {
		sort.Ints(r.Rolls)
		r.Removed = r.Rolls[:len(r.Rolls)-r.Keep]
		r.Rolls = r.Rolls[len(r.Rolls)-r.Keep:]
		r.Total = 0
		for _, n := range r.Rolls {
			r.Total += n
		}
	}
}
