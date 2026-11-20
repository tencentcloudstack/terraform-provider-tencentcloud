package ratelimit

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestProCheckContextCanceled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if err := ProCheckContext(ctx, t.Name(), "GetVectorBucket"); !errors.Is(err, context.Canceled) {
		t.Fatalf("cancellation ignored: %v", err)
	}
}

func TestProCheckContextDeadline(t *testing.T) {
	namespace, action := t.Name(), "GetVectorBucket"
	locker.Lock()
	limitConfig[namespace+"."+action] = 1
	locker.Unlock()
	if err := ProCheckContext(context.Background(), namespace, action); err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancel()
	start := time.Now()
	if err := ProCheckContext(ctx, namespace, action); !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("saturated limiter ignored deadline: %v", err)
	}
	if time.Since(start) > time.Second {
		t.Fatal("limiter did not exit promptly")
	}
}

func TestCheckContextCallerNamespace(t *testing.T) {
	action := t.Name()
	if err := CheckContext(context.Background(), action); err != nil {
		t.Fatal(err)
	}
	locker.Lock()
	defer locker.Unlock()
	if limitContainer["limit_test."+action] == nil || limitContainer["limit."+action] != nil {
		t.Fatal("context wrapper changed the caller namespace")
	}
}
