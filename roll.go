package hackyslack

import (
	"fmt"
	"math/rand"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

var diceRegEx = regexp.MustCompile(`(?P<num>\d{0,3})d(?P<sides>%|\d{0,4})(?P<max>[<>]\d{1,4})?(?P<keep>k\d{1,4})?(?P<mod>[+-]\d{1,4})?`)

type DiceRoll struct {
	Number   int
	Sides    int
	Modifier int
	Minimum  int
	Maximum  int
	Keep     int
	Rolls    []int
	Removed  []int
	Total    int
}

func init() {
	Register("roll", roll)
}

func parseRoll(text string) []*DiceRoll {
	var rolls []*DiceRoll
	for _, m := range diceRegEx.FindAllStringSubmatch(text, 5) {
		dice := &DiceRoll{
			Number: 2,
			Sides:  6,
		}
		for i, name := range diceRegEx.SubexpNames() {
			switch name {
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
			case "mod":
				dice.Modifier, _ = strconv.Atoi(m[i])
			}
		}
		rolls = append(rolls, dice)
	}
	if len(rolls) == 0 {
		rolls = append(rolls, &DiceRoll{
			Number: 2,
			Sides:  6,
		})
	}
	return rolls
}

func (r *DiceRoll) Roll() {
	r.Total = 0
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

func formatRoll(name string, results []*DiceRoll) D {
	var (
		color  string
		fields []D
		totals []string
		final  int
	)
	for _, result := range results {
		single := result.Number * result.Sides / 3.0
		rollText := fmt.Sprint(result.Rolls)
		if result.Total > single*2 {
			color = "good"
		} else if result.Total > single+result.Number-1 {
			color = "warning"
		} else {
			color = "danger"
		}
		fields = append(fields, D{
			"title": "Dice",
			"value": fmt.Sprint(result.Number, "d", result.Sides),
			"short": true,
		}, D{
			"title": "Rolls",
			"value": rollText[1 : len(rollText)-1],
			"short": true,
		})
		if result.Modifier != 0 {
			fields = append(fields, D{
				"title": "Raw",
				"value": strconv.Itoa(result.Total),
				"short": true,
			}, D{
				"title": "Modifier",
				"value": strconv.Itoa(result.Modifier),
				"short": true,
			})
		}
		if result.Minimum != 0 {
			count := 0
			for _, r := range result.Rolls {
				if r > result.Minimum {
					count++
				}
			}
			fields = append(fields, D{
				"title": "Minimum",
				"value": strconv.Itoa(result.Minimum),
				"short": true,
			}, D{
				"title": "Over",
				"value": strconv.Itoa(count),
				"short": true,
			})
		}
		if result.Maximum != 0 {
			count := 0
			for _, r := range result.Rolls {
				if r < result.Maximum {
					count++
				}
			}
			fields = append(fields, D{
				"title": "Maximum",
				"value": strconv.Itoa(result.Maximum),
				"short": true,
			}, D{
				"title": "Under",
				"value": strconv.Itoa(count),
				"short": true,
			})
		}
		if result.Keep != 0 {
			removed := fmt.Sprint(result.Removed)
			fields = append(fields, D{
				"title": "Keep",
				"value": strconv.Itoa(result.Keep),
				"short": true,
			}, D{
				"title": "Removed",
				"value": removed[1 : len(removed)-1],
				"short": true,
			})
		}
		total := result.Total + result.Modifier
		totals = append(totals, strconv.Itoa(total))
		final += total
	}
	var text, fallback string
	if len(totals) > 1 {
		text = fmt.Sprint("*", strings.Join(totals, "* + *"), "* = *", final, "*")
		fallback = fmt.Sprint(strings.Join(totals, "+"), " = ", final)
	} else {
		text = fmt.Sprint("*", strconv.Itoa(final), "*")
		fallback = strconv.Itoa(final)
	}
	return D{
		"response_type": "in_channel",
		"attachments": []D{
			D{
				"fallback":  fmt.Sprint("@", name, " rolled ", fallback),
				"text":      fmt.Sprint("@", name, " rolled ", text),
				// TODO: Color just uses the last color chosen.
				"color":     color,
				"fields":    fields,
				"mrkdwn_in": []string{"text"},
			},
		},
	}
}

func roll(args Args) D {
	rand.Seed(time.Now().UnixNano())
	result := parseRoll(args.Text)
	for _, roll := range result {
		roll.Roll()
	}
	return formatRoll(args.UserName, result)
}
