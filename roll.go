package hackyslack

import (
	"fmt"
	"math/rand"
	"regexp"
	"strconv"
	"time"
)

func init() {
	Register("roll", roll)
}

var dice = regexp.MustCompile(`(\d{0,3})d(%|\d{0,4})([+-]\d{1,4})?`)

func roll(args Args) D {
	rand.Seed(time.Now().UnixNano())
	m := dice.FindStringSubmatch(args.Text)
	num := 2
	sides := 6
	if len(m) > 2 {
		if m[1] == "" {
			num = 1
		} else {
			num, _ = strconv.Atoi(m[1])
		}
		if m[2] == "" {
			sides = 6
		} else if m[2] == "%" {
			sides = 100
		} else {
			sides, _ = strconv.Atoi(m[2])
		}
	}
	add := 0
	if len(m) > 3 {
		add, _ = strconv.Atoi(m[3])
	}
	if num < 1 {
		num = 1
	}
	if num > 100 {
		num = 100
	}
	if sides < 1 {
		sides = 1
	}
	if sides > 1000 {
		sides = 1000
	}
	total := 0
	rolls := []int{}
	for i := 0; i < num; i++ {
		n := rand.Intn(sides) + 1
		total += n
		rolls = append(rolls, n)
	}

	single := num * sides / 3.0
	rollText := fmt.Sprint(rolls)
	var color string
	if total > single * 2 {
		color = "good"
	} else if total > single + num - 1 {
		color = "warning"
	} else {
		color = "danger"
	}
	fields := []D{
		D{
			"title": "Dice",
			"value": fmt.Sprint(num, "d", sides),
			"short": true,
		},
		D{
			"title": "Rolls",
			"value": rollText[1 : len(rollText)-1],
			"short": true,
		},
	}
	if add != 0 {
		fields = append(fields, D{
			"title": "Raw",
			"value": strconv.Itoa(total),
			"short": true,
		}, D{
			"title": "Modifier",
			"value": strconv.Itoa(add),
			"short": true,
		})
	}
	total += add
	return D{
		"response_type": "in_channel",
		"attachments": []D{
			D{
				"fallback":  fmt.Sprint("@", args.UserName, " rolled ", total),
				"text":      fmt.Sprint("@", args.UserName, " rolled *", total, "*"),
				"color":     color,
				"fields":    fields,
				"mrkdwn_in": []string{"text"},
			},
		},
	}
}
