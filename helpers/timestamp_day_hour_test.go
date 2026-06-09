package helpers

import "testing"

func TestExtractDayAndHour(t *testing.T) {
	tests := []struct {
		name string

		offsets      TimestampOffsets
		timestampUTC int64
		expectedDay  int
		expectedHour int
	}{
		{
			name:         "1. Pure UTC - No Offset, Mid-day",
			timestampUTC: 1780932600, // 2026-06-08 15:30:00 UTC
			offsets: TimestampOffsets{
				OffsetUTC:          0,
				TimestampDSTSpring: 0,
				TimestampDSTWinter: 0,
			},
			expectedDay:  8,
			expectedHour: 15,
		},
		{
			name:         "2. New York Standard Time (-5h) - No DST active",
			timestampUTC: 1767225600, // 2026-01-01 00:00:00 UTC
			offsets: TimestampOffsets{
				OffsetUTC:          -18000,     // -5 hours
				TimestampDSTSpring: 1773481200, // March DST (future)
				TimestampDSTWinter: 1761901200, // Nov DST (past)
			},
			expectedDay:  31, // Shipped back to Dec 31, 2025
			expectedHour: 19, // 19:00 PM NY Time
		},
		{
			name:         "3. London Summer Time (+1h DST Active)",
			timestampUTC: 1780932600, // 2026-06-08 15:30:00 UTC
			offsets: TimestampOffsets{
				OffsetUTC:          0,          // London standard is 0
				TimestampDSTSpring: 1774755600, // March 29, 2026 01:00:00 UTC
				TimestampDSTWinter: 1792890000, // October 25, 2026 01:00:00 UTC
			},
			expectedDay:  8,
			expectedHour: 16, // 15:30 + 1 hour DST = 16:30
		},
		{
			name:         "4. Exactly on DST Spring Boundary Start",
			timestampUTC: 1774746000, // 2026-03-29 01:00:00 UTC
			offsets: TimestampOffsets{
				OffsetUTC:          0,
				TimestampDSTSpring: 1774746000,
				TimestampDSTWinter: 1792890000,
			},
			expectedDay:  29,
			expectedHour: 2,
		},
		{
			name:         "5. Exactly on DST Winter Boundary End (Back to Standard)",
			timestampUTC: 1792890000, // 2026-10-25 01:00:00 UTC
			offsets: TimestampOffsets{
				OffsetUTC:          0,
				TimestampDSTSpring: 1774746000,
				TimestampDSTWinter: 1792890000,
			},
			expectedDay:  25,
			expectedHour: 1,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			day, hour := ExtractDayAndHour(tc.timestampUTC, &tc.offsets)

			if day != tc.expectedDay || hour != tc.expectedHour {
				t.Errorf(
					"ExtractDayAndHour() failed for '%s'\nGot:  Day %d, Hour %d\nWant: Day %d, Hour %d",
					tc.name,
					day,
					hour,
					tc.expectedDay,
					tc.expectedHour,
				)
			}
		})
	}
}
