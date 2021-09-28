package helpers

import (
	"fmt"
	"time"
)

var (
	indonesiaDay = map[int]string{
		0: "Minggu",
		1: "Senin",
		2: "Selasa",
		3: "Rabu",
		4: "Kamis",
		5: "Jum'at",
		6: "Sabtu",
	}

	indonesiaMonth = map[int]string{
		1:  "Januari",
		2:  "Februari",
		3:  "Maret",
		4:  "April",
		5:  "Mei",
		6:  "Juni",
		7:  "Juli",
		8:  "Agustus",
		9:  "September",
		10: "Oktober",
		11: "Nopember",
		12: "Desember",
	}
)

func IndonesiaDateTimeFull(date time.Time, shortMonth bool) string {
	var month string = indonesiaMonth[int(date.Month())]
	if shortMonth {
		month = indonesiaMonth[int(date.Month())][:3]
	}
	hour, minute := convertHourMinute(date.Hour(), date.Minute())

	return fmt.Sprintf(
		"%s, %d %s %d %s:%s WIB",
		indonesiaDay[int(date.Weekday())],
		date.Day(), month,
		date.Year(), hour, minute,
	)
}

func IndonesiaDateTime(date time.Time, shortMonth bool) string {
	var month string = indonesiaMonth[int(date.Month())]
	if shortMonth {
		month = indonesiaMonth[int(date.Month())][:3]
	}
	hour, minute := convertHourMinute(date.Hour(), date.Minute())

	return fmt.Sprintf(
		"%d %s %d %s:%s WIB",
		date.Day(), month,
		date.Year(), hour, minute,
	)
}

func DateTime(date time.Time) string {
	var month string = fmt.Sprintf("%d", int(date.Month()))
	if int(date.Month()) < 10 {
		month = fmt.Sprintf("0%d", int(date.Month()))
	}

	hour, minute := convertHourMinute(date.Hour(), date.Minute())

	return fmt.Sprintf(
		"%d/%s/%d %s:%s WIB",
		date.Day(), month,
		date.Year(), hour, minute,
	)
}

func convertHourMinute(hour, minute int) (string, string) {
	var h, m string

	h = fmt.Sprintf("%d", hour)
	if hour < 10 {
		h = fmt.Sprintf("0%d", hour)
	}

	m = fmt.Sprintf("%d", minute)
	if minute < 10 {
		m = fmt.Sprintf("0%d", minute)
	}

	return h, m
}
