package progress

import (
	"errors"
	"testing"
)

func TestNewValidatesFraction(t *testing.T) {
	for _, f := range []float64{-0.01, 1.01} {
		if _, err := New(f, ""); !errors.Is(err, ErrInvalidFraction) {
			t.Errorf("fraction %v: expected ErrInvalidFraction, got %v", f, err)
		}
	}
	for _, f := range []float64{0, 0.5, 1} {
		if _, err := New(f, ""); err != nil {
			t.Errorf("fraction %v: unexpected error %v", f, err)
		}
	}
}

func TestAtTime(t *testing.T) {
	p, err := AtTime(754.5, 1509)
	if err != nil {
		t.Fatal(err)
	}
	if p.Locator != "t=754.5" || p.Fraction != 0.5 {
		t.Fatalf("got %+v", p)
	}
	secs, ok := p.Seconds()
	if !ok || secs != 754.5 {
		t.Fatalf("Seconds() = %v, %v", secs, ok)
	}
	if _, ok := (Progress{Locator: "page=3"}).Seconds(); ok {
		t.Fatal("page locator should not parse as seconds")
	}
}

func TestAtTimeUnknownDurationAndClamp(t *testing.T) {
	p, err := AtTime(30, 0)
	if err != nil || p.Fraction != 0 || p.Locator != "t=30" {
		t.Fatalf("unknown duration: %+v, %v", p, err)
	}
	p, err = AtTime(120, 100)
	if err != nil || p.Fraction != 1 {
		t.Fatalf("past end should clamp to 1: %+v, %v", p, err)
	}
	if _, err := AtTime(-1, 100); !errors.Is(err, ErrInvalidPosition) {
		t.Fatalf("negative position: %v", err)
	}
}

func TestAtPage(t *testing.T) {
	p, err := AtPage(12, 48)
	if err != nil || p.Locator != "page=12" || p.Fraction != 0.25 {
		t.Fatalf("got %+v, %v", p, err)
	}
	page, ok := p.Page()
	if !ok || page != 12 {
		t.Fatalf("Page() = %v, %v", page, ok)
	}
	if _, err := AtPage(0, 10); !errors.Is(err, ErrInvalidPosition) {
		t.Fatalf("page 0: %v", err)
	}
}

func TestFinished(t *testing.T) {
	if !(Progress{Fraction: 0.95}).Finished() || (Progress{Fraction: 0.5}).Finished() {
		t.Fatal("Finished threshold wrong")
	}
}
