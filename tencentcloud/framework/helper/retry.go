package helper

import (
	"context"
	"errors"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// RetryableErrorFn is the callback used by RetryFramework to classify an
// error as "retryable". Returning true means the error is transient and
// should be retried within the eventual-consistency window; false means
// the call should fail immediately.
//
// Default (when nil is passed): every non-nil error is treated as
// non-retryable. Framework callers should supply their own predicate (for
// instance, white-listing specific TencentCloud SDK error codes).
type RetryableErrorFn func(err error) bool

// RetryFramework is the framework-flavoured equivalent of helper.Retry in
// the SDKv2 world.
//
// Differences from SDKv2 helper.Retry:
//   - Callers no longer need to wrap *resource.RetryError manually; the
//     business logic only returns (T, error) and the retryable predicate
//     decides whether the error is retryable.
//   - It accepts ctx, so timeout and cancellation are governed by both ctx
//     and the timeout argument.
//
// Typical usage:
//
//	out, err := helper.RetryFramework(ctx, 30*time.Second, isResourceTransient,
//	    func() (*sdk.DescribeFooResponse, error) {
//	        return client.UseFooClient().DescribeFoo(req)
//	    })
//
// When err is nil it returns directly; when err is classified as retryable
// it enters the retry loop, otherwise it terminates immediately.
func RetryFramework[T any](
	ctx context.Context,
	timeout time.Duration,
	retryable RetryableErrorFn,
	fn func() (T, error),
) (T, error) {
	if retryable == nil {
		retryable = func(error) bool { return false }
	}

	var (
		zero T
		out  T
	)

	retryErr := resource.RetryContext(ctx, timeout, func() *resource.RetryError {
		got, err := fn()
		if err == nil {
			out = got
			return nil
		}
		if retryable(err) {
			return resource.RetryableError(err)
		}
		return resource.NonRetryableError(err)
	})
	if retryErr != nil {
		return zero, retryErr
	}
	return out, nil
}

// IsTimeoutError reports whether helper.RetryContext exited because of an
// overall timeout, so that callers can translate the timeout into a more
// user-friendly diagnostic.
func IsTimeoutError(err error) bool {
	if err == nil {
		return false
	}
	var te *resource.TimeoutError
	return errors.As(err, &te)
}
