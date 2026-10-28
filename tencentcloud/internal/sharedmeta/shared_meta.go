package sharedmeta

import (
	"sync/atomic"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
)

// sharedMeta carries the configured *connectivity.TencentCloudClient between
// the SDKv2 provider and the framework provider within the same process.
//
// Design notes:
//   - atomic.Pointer guarantees concurrency safety (mux may invoke
//     ConfigureProvider on multiple stacks concurrently).
//   - Only the SDKv2 ConfigureContextFunc writes; the framework Configure is
//     read-only.
//   - The exported API never exposes pointer-assignment side effects; the
//     caller always observes the same *Client instance.
var sharedMeta atomic.Pointer[connectivity.TencentCloudClient]

// SetSharedMeta is invoked by the SDKv2 provider after credentials are
// resolved and the client is constructed, exposing the same client to the
// framework provider for reuse.
//
// Multiple invocations are legal: last write wins (e.g. acceptance tests that
// reset state between cases).
func SetSharedMeta(c *connectivity.TencentCloudClient) {
	sharedMeta.Store(c)
}

// GetSharedMeta is invoked by the framework provider in Configure.
// A nil return means the SDKv2 provider has not finished configuring; the
// caller should emit a diagnostic instead of panicking.
func GetSharedMeta() *connectivity.TencentCloudClient {
	return sharedMeta.Load()
}

// ResetSharedMetaForTest is intended only for unit tests to reset state
// between cases. Do not call it from production code.
func ResetSharedMetaForTest() {
	sharedMeta.Store(nil)
}
