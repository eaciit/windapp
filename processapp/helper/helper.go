package helper

import (
	"strings"
	"time"

	tk "github.com/eaciit/toolkit"
)

func GenNext10Minutes(current time.Time) time.Time {
	strTimeStamp := current.Format("15-04")
	anTimeStamp := strings.Split(strTimeStamp, "-")

	var hourInt, minuteInt int
	hourInt = tk.ToInt(anTimeStamp[0], tk.RoundingAuto)
	minuteInt = tk.ToInt(anTimeStamp[1], tk.RoundingAuto)

	var year, month, day, hour, minute string

	hour = tk.ToString(hourInt)
	if minuteInt+10 == 60 {
		minuteInt = 0
		hourInt++

		minute = "00"
		hour = tk.ToString(hourInt)

		if hourInt == 24 {
			hourInt = 0
			hour = "00"

			current = current.AddDate(0, 0, 1)
		}
	} else {
		minuteInt += 10
		minute = tk.ToString(minuteInt)
	}

	if hourInt < 10 && hour != "00" {
		hour = "0" + hour
	}

	if minuteInt < 10 && minute != "00" {
		minute = "0" + minute
	}

	current, _ = time.Parse("2006-01-02 15:04", current.Format("2006-01-02")+" "+hour+":"+minute)

	year = tk.ToString(current.Year())
	month = tk.ToString(int(current.Month()))
	day = tk.ToString(current.Day())

	if int(current.Month()) < 10 {
		month = "0" + month
	}

	if current.Day() < 10 {
		day = "0" + day
	}

	current, _ = time.Parse("20060102 15:04", year+month+day+" "+hour+":"+minute)
	return current
}
