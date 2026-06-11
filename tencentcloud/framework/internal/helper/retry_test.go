package helper

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestRetryFramework_Success(t *testing.T) {
	ctx := context.Background()
	calls := 0

	got, err := RetryFramework(ctx, time.Second, nil, func() (string, error) {
		calls++
		return "ok", nil
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "ok" {
		t.Fatalf("expected ok, got %q", got)
	}
	if calls != 1 {
		t.Fatalf("expected 1 call, got %d", calls)
	}
}

func TestRetryFramework_NonRetryableExitsImmediately(t *testing.T) {
	ctx := context.Background()
	calls := 0
	target := errors.New("permanent failure")

	_, err := RetryFramework(ctx, time.Second, nil, func() (int, error) {
		calls++
		return 0, target
	})
	if err == nil || !errors.Is(err, target) {
		t.Fatalf("expected to surface %v, got %v", target, err)
	}
	if calls != 1 {
		t.Fatalf("expected to fail fast on first call, got %d calls", calls)
	}
}

func TestRetryFramework_RetryableEventuallySucceeds(t *testing.T) {
	ctx := context.Background()
	transient := errors.New("transient")
	calls := 0

	got, err := RetryFramework(ctx, 5*time.Second,
		func(err error) bool { return errors.Is(err, transient) },
		func() (string, error) {
			calls++
			if calls < 3 {
				return "", transient
			}
			return "done", nil
		},
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "done" {
		t.Fatalf("expected done, got %q", got)
	}
	if calls < 3 {
		t.Fatalf("expected at least 3 calls, got %d", calls)
	}
}

func TestRetryFramework_RetryableUntilTimeout(t *testing.T) {
	ctx := context.Background()
	transient := errors.New("transient")

	_, err := RetryFramework(ctx, 200*time.Millisecond,
		func(err error) bool { return errors.Is(err, transient) },
		func() (struct{}, error) {
			return struct{}{}, transient
		},
	)
	if err == nil {
		t.Fatalf("expected error after timeout")
	}
	if !IsTimeoutError(err) && !errors.Is(err, transient) {
		t.Fatalf("expected timeout or transient error, got %v", err)
	}
}

func TestIsTimeoutError(t *testing.T) {
	if IsTimeoutError(nil) {
		t.Fatalf("nil should not be timeout")
	}
	if IsTimeoutError(errors.New("plain")) {
		t.Fatalf("plain error should not be timeout")
	}
}
