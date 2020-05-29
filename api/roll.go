package api

import (
	"encoding/json"
	"fmt"
	"github.com/arkie/dicebot/roll"
	"go.uber.org/zap"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type D map[string]interface{}

func helpGen(markdown bool) string {
	backtick := ""
	if markdown {
		backtick = "`"
	}
	return backtick + `/roll [XdY]` + backtick + ` - Roll X Y-sided dice

	Note that the following standard notation dice expression are also supported:

	- ` + backtick + `df` + backtick + ` for fudge dice
	- ` + backtick + `d%` + backtick + ` for percentile dice
	- A ` + backtick + `!` + backtick + ` to indicate that the dice "explode", rolling an additional die of the same size when the max value is rolled.
	- A ` + backtick + `<5` + backtick + ` or ` + backtick + `>5` + backtick + ` (for any number between 2 and the die size - 1) to indicate that you want to count up how many dice rolled below or above the value.
	- A ` + backtick + `k5` + backtick + ` or ` + backtick + `k-5` + backtick + ` (for any number between -999 and 999) to indicate that only the best K (for positive values) or worst K (for negative values) rolls should be counted in the result.
	`
}

var (
	help         = helpGen(true)
	helpFallback = helpGen(false)
)

func Roll(w http.ResponseWriter, r *http.Request) {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Printf("can't initialize zap logger: %v", err)
	}
	defer logger.Sync()
	rand.Seed(time.Now().UnixNano())
	text := r.FormValue("text")
	mini := strings.HasPrefix(text, "mini")
	silent := strings.HasPrefix(text, "silent")
	result := roll.Parse(text)
	for _, roll := range result {
		roll.Roll()
	}
	logger.Info("rolling", zap.Any("form", r.Form))
	if text == "help" {
		writeJson(w, r, D{
			"response_type": "ephemeral",
			"attachments": []D{
				{
					"fallback":  helpFallback,
					"text":      help,
					"mrkdwn_in": []string{"text"},
				},
			},
		})
		return
	}
	data := formatRoll(r.FormValue("user_id"), mini, silent, result)
	writeJson(w, r, data)
}

func writeJson(w http.ResponseWriter, r *http.Request, data D) {
	bytes, err := json.Marshal(data)
	if err != nil {
		log.Printf("ERROR - Failed to mashal %v: %v", data, err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(bytes)
}

func formatRoll(id string, mini bool, silent bool, results []*roll.Dice) D {
	var (
		color     string
		fields    []D
		final     int
		text      string
		fallback  string
		rollCount int
		finalSum  int
		forCount  int
	)
	for i, result := range results {
		op := result.Operator
		switch result.Operator {
		case roll.Add:
			final += result.Total
		case roll.Subtract:
			final -= result.Total
		case roll.Multiply:
			final *= result.Total
			op = "×"
		case roll.Divide:
			// Just ignoring divide by zero.
			if result.Total != 0 {
				final /= result.Total
			}
		case roll.Max:
			if result.Total > final {
				final = result.Total
			}
		case roll.Min:
			if result.Total < final {
				final = result.Total
			}
		}
		if i != 0 {
			text += fmt.Sprint(" ", op, " ")
			fallback += fmt.Sprint(" ", op, " ")
		}
		text += fmt.Sprint("*", result.Total, "*")
		fallback += fmt.Sprint(result.Total)
		rollCount++
		if result.For != "" {
			forCount++
			if rollCount > 1 {
				text += fmt.Sprint(" = *", final, "*")
				fallback += fmt.Sprint(" = ", final)
			}
			finalSum += final
			final = 0
			rollCount = 0
			text += fmt.Sprint(" for *", result.For, "*")
			fallback += fmt.Sprint(" for ", result.For)
		}
		if i == len(results)-1 && i > 0 && (forCount != 1 || rollCount > 0) {
			text += fmt.Sprint(" = *", finalSum+final, "*")
			fallback += fmt.Sprint(" = ", finalSum+final)
		}
		if result.Sides <= 1 {
			continue
		}
		rollText := fmt.Sprint(result.Rolls)
		if result.Fudge {
			if result.Total > 0 {
				color = "good"
			} else if result.Total == 0 {
				color = "warning"
			} else {
				color = "danger"
			}
		} else {
			single := result.Number * result.Sides / 3.0
			if result.Total > single*2 {
				color = "good"
			} else if result.Total > single+result.Number-1 {
				color = "warning"
			} else {
				color = "danger"
			}
		}
		dice := fmt.Sprint(result.Number, "d", result.Sides)
		if result.Fudge {
			dice = fmt.Sprint(result.Number, "df")
		} else if result.Explode {
			dice += "!"
		}
		fields = append(fields, D{
			"title": "Dice",
			"value": dice,
			"short": true,
		}, D{
			"title": "Rolls",
			"value": rollText[1 : len(rollText)-1],
			"short": true,
		})
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
	}
	if mini {
		fields = []D{}
	}
	response := "in_channel"
	if silent {
		response = "ephemeral"
	}
	return D{
		"response_type": response,
		"attachments": []D{
			{
				"fallback": fmt.Sprint("<@", id, "> rolled ", fallback),
				"text":     fmt.Sprint("<@", id, "> rolled ", text),
				// TODO: Color just uses the last color chosen.
				"color":     color,
				"fields":    fields,
				"mrkdwn_in": []string{"text"},
			},
		},
	}
}
