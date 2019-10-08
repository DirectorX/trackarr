package strings

import (
	"strconv"
	"strings"
	"time"
)

func s(x int) string {
	if x == 1 {
		return ""
	}

	return "s"
}

func timeDiff(a, b time.Time) (year, month, day, hour, min, sec int) {
	if a.Location() != b.Location() {
		b = b.In(a.Location())
	}

	if a.After(b) {
		a, b = b, a
	}

	y1, M1, d1 := a.Date()
	y2, M2, d2 := b.Date()

	h1, m1, s1 := a.Clock()
	h2, m2, s2 := b.Clock()

	year = y2 - y1
	month = int(M2 - M1)
	day = d2 - d1
	hour = h2 - h1
	min = m2 - m1
	sec = s2 - s1

	// Normalize negative values
	if sec < 0 {
		sec += 60
		min--
	}

	if min < 0 {
		min += 60
		hour--
	}

	if hour < 0 {
		hour += 24
		day--
	}

	if day < 0 {
		// days in month:
		t := time.Date(y1, M1, 32, 0, 0, 0, 0, time.UTC)
		day += 32 - t.Day()
		month--
	}

	if month < 0 {
		month += 12
		year--
	}

	return year, month, day, hour, min, sec
}

func TimeDiffDurationString(now time.Time, then time.Time, full bool) string {
	/* credits:
	https://www.socketloop.com/tutorials/golang-human-readable-time-elapsed-format-such-as-5-days-ago
	https://stackoverflow.com/a/36531443
	*/
	var parts []string

	var text string

	// get diff
	year, month, day, hour, minute, second := timeDiff(then, now)

	// process diff
	if year > 0 {
		parts = append(parts, strconv.Itoa(year)+" year"+s(year))
	}

	if month > 0 {
		parts = append(parts, strconv.Itoa(month)+" month"+s(month))
	}

	if day > 0 {
		parts = append(parts, strconv.Itoa(day)+" day"+s(day))
	}

	if hour > 0 {
		parts = append(parts, strconv.Itoa(hour)+" hour"+s(hour))
	}

	if minute > 0 {
		parts = append(parts, strconv.Itoa(minute)+" minute"+s(minute))
	}

	if second > 0 {
		parts = append(parts, strconv.Itoa(second)+" second"+s(second))
	}

	if len(parts) == 0 {
		return "just now"
	}

	if full {
		return strings.Join(parts, ", ") + text
	}
	return parts[0] + text
}
