package helper

import (
	"errors"
	"strings"
	"testing"

	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
)

func TestWrapSDKError_TencentCloudSDKError(t *testing.T) {
	err := sdkErrors.NewTencentCloudSDKError("FailedOperation", "boom", "req-123")

	d := WrapSDKError("Failed to do thing", err)
	if d.Severity().String() != "Error" {
		t.Fatalf("expected Error severity, got %s", d.Severity())
	}
	if d.Summary() != "Failed to do thing" {
		t.Fatalf("unexpected summary: %s", d.Summary())
	}

	detail := d.Detail()
	for _, want := range []string{"FailedOperation", "boom", "req-123"} {
		if !strings.Contains(detail, want) {
			t.Errorf("expected detail to contain %q, got: %s", want, detail)
		}
	}
}

func TestWrapSDKError_PlainError(t *testing.T) {
	d := WrapSDKError("ctx", errors.New("oops"))
	if d.Detail() != "oops" {
		t.Fatalf("expected detail oops, got %s", d.Detail())
	}
}

func TestWrapSDKError_NilError(t *testing.T) {
	d := WrapSDKError("ctx", nil)
	if d.Detail() == "" {
		t.Fatalf("expected non-empty detail even for nil error")
	}
}

func TestWrapSDKError_EmptySummaryFallback(t *testing.T) {
	d := WrapSDKError("", errors.New("oops"))
	if d.Summary() == "" {
		t.Fatalf("summary should be replaced when empty")
	}
}

func TestIsSDKErrorCode(t *testing.T) {
	err := sdkErrors.NewTencentCloudSDKError("InternalError", "x", "rid")

	if !IsSDKErrorCode(err, "InternalError") {
		t.Fatalf("expected match")
	}
	if !IsSDKErrorCode(err, "Other", "InternalError") {
		t.Fatalf("expected match in list")
	}
	if IsSDKErrorCode(err, "FailedOperation") {
		t.Fatalf("unexpected match")
	}
	if IsSDKErrorCode(errors.New("plain"), "InternalError") {
		t.Fatalf("plain error should not match")
	}
	if IsSDKErrorCode(nil, "X") {
		t.Fatalf("nil should not match")
	}
}
