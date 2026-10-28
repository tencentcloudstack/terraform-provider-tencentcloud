package helper

import (
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"

	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
)

// WrapSDKError translates an arbitrary error into a structured framework
// diag.Diagnostic.
//
// When err is *sdkErrors.TencentCloudSDKError:
//   - Summary uses the supplied summary (the business-side context, e.g.
//     "Failed to create foo").
//   - Detail includes Code, Message and RequestId so that troubleshooting
//     and ticket filing are easier.
//
// When err is any other type:
//   - Falls back to a generic error rendering; Detail contains only
//     err.Error().
//
// summary must not be empty; an empty summary is replaced with
// "TencentCloud API error" so the framework runtime does not reject the
// diagnostic as malformed.
func WrapSDKError(summary string, err error) diag.Diagnostic {
	if summary == "" {
		summary = "TencentCloud API error"
	}

	var sdkErr *sdkErrors.TencentCloudSDKError
	if errors.As(err, &sdkErr) {
		detail := fmt.Sprintf(
			"Code: %s\nMessage: %s\nRequestId: %s",
			sdkErr.GetCode(),
			sdkErr.GetMessage(),
			sdkErr.GetRequestId(),
		)
		return diag.NewErrorDiagnostic(summary, detail)
	}

	if err == nil {
		return diag.NewErrorDiagnostic(summary, "(no error detail provided)")
	}
	return diag.NewErrorDiagnostic(summary, err.Error())
}

// IsSDKErrorCode is the standard retryable predicate template: it returns
// true when err is a SDK error and its Code is contained in codes. A
// resource can use it inside RetryFramework like this:
//
//	helper.RetryFramework(ctx, timeout,
//	    func(err error) bool {
//	        return helper.IsSDKErrorCode(err,
//	            "InternalError", "ResourceInUse", "FailedOperation")
//	    },
//	    func() (T, error) { ... })
func IsSDKErrorCode(err error, codes ...string) bool {
	var sdkErr *sdkErrors.TencentCloudSDKError
	if !errors.As(err, &sdkErr) {
		return false
	}
	for _, c := range codes {
		if sdkErr.GetCode() == c {
			return true
		}
	}
	return false
}
