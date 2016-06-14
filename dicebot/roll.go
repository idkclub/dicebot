package dicebot

import (
	"fmt"
	"github.com/arkie/hackyslack2"
	"github.com/arkie/hackyslack2/dicebot/roll"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

func init() {
	hackyslack.Register("roll", command)
}

func formatRoll(name string, results []*roll.Dice) hackyslack.D {
	var (
		color  string
		fields []hackyslack.D
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
		fields = append(fields, hackyslack.D{
			"title": "Dice",
			"value": fmt.Sprint(result.Number, "d", result.Sides),
			"short": true,
		}, hackyslack.D{
			"title": "Rolls",
			"value": rollText[1 : len(rollText)-1],
			"short": true,
		})
		if result.Modifier != 0 {
			fields = append(fields, hackyslack.D{
				"title": "Raw",
				"value": strconv.Itoa(result.Total),
				"short": true,
			}, hackyslack.D{
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
			fields = append(fields, hackyslack.D{
				"title": "Minimum",
				"value": strconv.Itoa(result.Minimum),
				"short": true,
			}, hackyslack.D{
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
			fields = append(fields, hackyslack.D{
				"title": "Maximum",
				"value": strconv.Itoa(result.Maximum),
				"short": true,
			}, hackyslack.D{
				"title": "Under",
				"value": strconv.Itoa(count),
				"short": true,
			})
		}
		if result.Keep != 0 {
			removed := fmt.Sprint(result.Removed)
			fields = append(fields, hackyslack.D{
				"title": "Keep",
				"value": strconv.Itoa(result.Keep),
				"short": true,
			}, hackyslack.D{
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
	return hackyslack.D{
		"response_type": "in_channel",
		"attachments": []hackyslack.D{
			{
				"fallback": fmt.Sprint("@", name, " rolled ", fallback),
				"text":     fmt.Sprint("@", name, " rolled ", text),
				// TODO: Color just uses the last color chosen.
				"color":     color,
				"fields":    fields,
				"mrkdwn_in": []string{"text"},
			},
		},
	}
}

func command(args hackyslack.Args) hackyslack.D {
	rand.Seed(time.Now().UnixNano())
	result := roll.Parse(args.Text)
	for _, roll := range result {
		roll.Roll()
	}
	return formatRoll(args.UserName, result)
}
