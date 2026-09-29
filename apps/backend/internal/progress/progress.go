// Package progress is the single model for "how far through a work is someone",
// shared by history, watch parties and device hand-off.
//
// Fraction is medium-independent (0 = not started, 1 = done) so the Hub can
// sort, filter and scrobble without knowing the medium. Locator is the exact
// resume point, owned by whichever engine plays the work. The Hub treats it as
// opaque; the helpers here cover the two formats the Hub itself produces:
// "t=<seconds>" for time-based media and "page=<n>" for paged media.
package progress

import (
	"errors"
	"strconv"
	"strings"
)

// FinishedThreshold is the fraction at which a work counts as finished; credits
// and back matter mean people rarely reach exactly 1.
const FinishedThreshold = 0.9

var (
	ErrInvalidFraction = errors.New("progress fraction must be between 0 and 1")
	ErrInvalidPosition = errors.New("progress position is out of range")
)

type Progress struct {
	Fraction float64 `json:"fraction" db:"fraction"`
	Locator  string  `json:"locator" db:"locator"`
}

func New(fraction float64, locator string) (Progress, error) {
	if fraction < 0 || fraction > 1 {
		return Progress{}, ErrInvalidFraction
	}
	return Progress{Fraction: fraction, Locator: locator}, nil
}

// AtTime builds progress for time-based media. A duration <= 0 means unknown,
// which leaves Fraction at 0 but still records the resume point.
func AtTime(seconds, duration float64) (Progress, error) {
	if seconds < 0 {
		return Progress{}, ErrInvalidPosition
	}
	return Progress{
		Fraction: fraction(seconds, duration),
		Locator:  "t=" + strconv.FormatFloat(seconds, 'f', -1, 64),
	}, nil
}

// AtPage builds progress for paged media; pages are 1-based.
func AtPage(page, total int) (Progress, error) {
	if page < 1 {
		return Progress{}, ErrInvalidPosition
	}
	return Progress{
		Fraction: fraction(float64(page), float64(total)),
		Locator:  "page=" + strconv.Itoa(page),
	}, nil
}

// Seconds returns the position of a "t=" locator.
func (p Progress) Seconds() (float64, bool) {
	raw, ok := strings.CutPrefix(p.Locator, "t=")
	if !ok {
		return 0, false
	}
	secs, err := strconv.ParseFloat(raw, 64)
	return secs, err == nil
}

// Page returns the page of a "page=" locator.
func (p Progress) Page() (int, bool) {
	raw, ok := strings.CutPrefix(p.Locator, "page=")
	if !ok {
		return 0, false
	}
	page, err := strconv.Atoi(raw)
	return page, err == nil
}

func (p Progress) Finished() bool {
	return p.Fraction >= FinishedThreshold
}

func fraction(pos, total float64) float64 {
	if total <= 0 {
		return 0
	}
	return min(pos/total, 1)
}
