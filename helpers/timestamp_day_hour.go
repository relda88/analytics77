package helpers

type TimestampOffsets struct {
	OffsetUTC int64

	TimestampDSTWinter int64 // Epoch when DST ends (falls back to winter/standard time)
	TimestampDSTSpring int64 // Epoch when DST starts (springs forward to summer time)
}

func ExtractDayAndHour(timestampUTC int64, offsets *TimestampOffsets) (day int, hour int) {
	// Start with base timestamp + standard offset
	localTimestamp := timestampUTC + offsets.OffsetUTC

	// Check if the timestamp falls within the DST active window.
	// If it does, we inject the extra 1 hour (3600 seconds) savings.
	if timestampUTC >= offsets.TimestampDSTSpring && timestampUTC < offsets.TimestampDSTWinter {
		localTimestamp = localTimestamp + 3600
	}

	totalHours := localTimestamp / 3600
	hour = int(totalHours % 24)

	// Handle Go's truncated division behavior on negative local timestamps
	if hour < 0 {
		hour = hour + 24
	}

	// 86400 seconds = 1 day.
	totalDays := localTimestamp / 86400

	// Handle negative totalDays if the offset pushes the time before 1970
	if localTimestamp < 0 && localTimestamp%86400 != 0 {
		totalDays--
	}

	// Unix epoch (Jan 1, 1970) was a Thursday.
	// To convert totalDays since 1970 into a specific day of the current month
	// using pure integer math requires a fast epoch-to-date algorithm (like civil time).

	// optimized integer algorithm for UTC date extraction:
	totalDays = totalDays + 719468 // Offset to March 1, 0000

	doe := totalDays % 146097
	yoe := (doe - doe/1460 + doe/36524 - doe/146096) / 365
	doy := doe - (365*yoe + yoe/4 - yoe/100)

	mp := (5*doy + 2) / 153
	day = int(doy - (153*mp+2)/5 + 1)

	return day, hour
}
