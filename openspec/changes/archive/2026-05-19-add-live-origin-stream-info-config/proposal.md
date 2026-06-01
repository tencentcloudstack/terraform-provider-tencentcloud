## Why

The Tencent Cloud Live (CSS) service provides a `ModifyOriginStreamInfo` API to configure origin stream settings for a playback domain, and `DescribeOriginStreamInfo` to query those settings. These APIs are not yet covered by the Terraform provider, leaving users unable to manage live origin stream configuration as infrastructure-as-code.

## What Changes

- Add new config-type resource `tencentcloud_live_origin_stream_info_config` under the `live` service.
- Resource schema fields mirror the `ModifyOriginStreamInfo` request parameters exactly (required and optional).
- Async API handling: after calling `ModifyOriginStreamInfo`, poll `DescribeOriginStreamInfo` until `Status` is `1` (success) or `3` (closed successfully).
- Resource ID is `DomainName` (unique key for the configuration).
- Add service-layer method `DescribeLiveOriginStreamInfo` to `service_tencentcloud_live.go`.
- Register resource in `provider.go`.
- Generate `.md` documentation and `_test.go` acceptance test file.

## Capabilities

### New Capabilities
- `live-origin-stream-info-config`: Full CRUD lifecycle (config pattern: Create→Update, Read, Update, Delete→no-op) for live origin stream configuration, including async wait on modify operations.

### Modified Capabilities

## Impact

- New files: `tencentcloud/services/live/resource_tc_live_origin_stream_info_config.go`, `.md`, `_test.go`
- Modified files: `tencentcloud/services/live/service_tencentcloud_live.go`, `tencentcloud/provider.go`
- SDK dependency: `live/v20180801` — `ModifyOriginStreamInfo` and `DescribeOriginStreamInfo` both present in vendored SDK.
