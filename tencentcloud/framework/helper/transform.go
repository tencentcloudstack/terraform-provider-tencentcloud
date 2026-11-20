package helper

import "time"

// The functions below mirror the frequently-used value-to-pointer and
// time-format helpers from tencentcloud/internal/helper so that framework
// resources can stay on this package alone. The SDKv2-specific parts of
// internal/helper (schema.ResourceData plumbing, flatmap Interfaces*
// conversions, ...) are deliberately NOT mirrored: on the framework side
// those jobs are done by fwHelper.XValueOrNull / XPointerFromValue.

func Bool(i bool) *bool { return &i }

func String(i string) *string { return &i }

func Int(i int) *int { return &i }

func Int64(i int64) *int64 { return &i }

func Uint64(i uint64) *uint64 { return &i }

func Float64(i float64) *float64 { return &i }

// Strings converts a []string into []*string; an empty input yields nil so
// the request field stays "omitted".
func Strings(strs []string) []*string {
	if len(strs) == 0 {
		return nil
	}

	sp := make([]*string, 0, len(strs))
	for _, s := range strs {
		sp = append(sp, String(s))
	}

	return sp
}

// FormatUnixTime renders a Unix-seconds timestamp as the RFC3339 UTC string
// TencentCloud APIs commonly use ("2006-01-02T15:04:05Z").
func FormatUnixTime(n uint64) string {
	return time.Unix(int64(n), 0).UTC().Format("2006-01-02T15:04:05Z")
}

// ParseTime parses the RFC3339-style UTC timestamp emitted by TencentCloud
// APIs ("2006-01-02T15:04:05Z").
func ParseTime(s string) (time.Time, error) {
	return time.ParseInLocation("2006-01-02T15:04:05Z", s, time.UTC)
}
