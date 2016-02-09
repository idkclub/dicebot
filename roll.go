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

var dice = regexp.MustCompile(`(\d{1,3})d(\d{1,4})([+-]\d{1,4})?`)

func roll(args Args) D {
	rand.Seed(time.Now().UnixNano())
	m := dice.FindStringSubmatch(args.Text)
	num := 2
	sides := 6
	if len(m) > 2 {
		num, _ = strconv.Atoi(m[1])
		sides, _ = strconv.Atoi(m[2])
	}
	add := 0
	if len(m) > 3 {
		add, _ = strconv.Atoi(m[3])
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

	max := num * sides
	rollText := fmt.Sprint(rolls)
	var color string
	if total > max/2 {
		color = "good"
	} else {
		color = "danger"
	}
	total += add
	return D{
		"response_type": "in_channel",
		"attachments": []D{
			D{
				"fallback": fmt.Sprint("@", args.UserName, " rolled ", total),
				"text":     fmt.Sprint("@", args.UserName, " rolled *", total, "*"),
				"color":    color,
				"fields": []D{
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
				},
				"mrkdwn_in": []string{"text"},
			},
		},
	}
}
