package ga2

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
)

// HandleGa2ResourceNotFoundError checks if the given error is a ResourceNotFound SDK error.
// If the error is ResourceNotFound and the resource is not new (i.e., not during initial create
// where d.IsNewResource() returns true), it logs a warning, clears the resource ID, and returns
// true to indicate the error was handled (caller should return nil).
// Otherwise, it returns false, meaning the caller should propagate the error.
func HandleGa2ResourceNotFoundError(err error, d *schema.ResourceData, resourceName, logId string) bool {
	if err == nil {
		return false
	}
	if sdkErr, ok := err.(*sdkErrors.TencentCloudSDKError); ok {
		if !d.IsNewResource() && sdkErr.Code == "ResourceNotFound" {
			log.Printf("[WARN]%s resource `%s` [%s] not found, please check if it has been deleted.\n", logId, resourceName, d.Id())
			d.SetId("")
			return true
		}
	}
	return false
}

// HandleGa2ReadNotFound is a unified handler for Read operations to check whether the
// describe API indicates the resource does not exist. It handles two cases:
//
//  1. SDK ResourceNotFound error: delegates to HandleGa2ResourceNotFoundError.
//  2. Nil / empty response (respData is nil): when the resource is not new, logs a
//     warning and clears the resource ID. When d.IsNewResource() is true, it returns
//     an error so the initial Create → Read cycle can propagate the failure normally.
//
// Returns (handled, err):
//   - (true, nil) — the not-found case was handled; caller should return nil.
//   - (false, non-nil error) — caller should propagate the error.
//   - (false, nil) — neither case applies; caller should proceed with field hydration.
func HandleGa2ReadNotFound(err error, respData interface{}, d *schema.ResourceData, resourceName, logId string) (bool, error) {
	if err != nil {
		if HandleGa2ResourceNotFoundError(err, d, resourceName, logId) {
			return true, nil
		}
		return false, fmt.Errorf("[%s] describe %s failed: %w", logId, resourceName, err)
	}

	if respData == nil {
		if d.IsNewResource() {
			return false, fmt.Errorf("[%s] %s not found during create", logId, resourceName)
		}
		log.Printf("[WARN]%s resource `%s` [%s] not found, please check if it has been deleted.\n", logId, resourceName, d.Id())
		d.SetId("")
		return true, nil
	}

	return false, nil
}
