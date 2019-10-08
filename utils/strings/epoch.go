package strings

import "time"

/* Public */

func UnixEpochToUtcTimestamp(epoch int64) string {
	if epoch == 0 {
		return "Unknown"
	}
	secs, millis := divmod(epoch, 1000)

	return time.Unix(secs, millis).UTC().Format(time.RFC822)
}

/* Private */

func divmod(numerator, denominator int64) (quotient, remainder int64) {
	/* Credits: https://stackoverflow.com/a/43945812 */
	quotient = numerator / denominator // integer division, decimals are truncated
	remainder = numerator % denominator
	return
}
