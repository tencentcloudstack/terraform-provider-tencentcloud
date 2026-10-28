## Why

The `tencentcloud_waf_cc_session` resource manages WAF session definitions but does not expose the `Key` parameter, which is required by the TencentCloud WAF `UpsertSession` API for precise matching scenarios. The API already supports this field (`request.Key`, documented as "精准匹配时配置的key" / "key configured for precise matching"), and both the request struct (`UpsertSessionRequest`) and response struct (`SessionItem`) contain the `Key` field. Exposing it lets users configure precise session key matching directly from Terraform.

## What Changes

- Add a new optional string parameter `key` to the `tencentcloud_waf_cc_session` resource schema.
- Populate `request.Key` from the schema in the `Create` (`resourceTencentCloudWafCcSessionCreate`) and `Update` (`resourceTencentCloudWafCcSessionUpdate`) operations when calling the `UpsertSession` API.
- Read `Key` from the `DescribeSession` API response (`SessionItem.Key`) and set it into Terraform state via `d.Set("key", ...)` in the `Read` operation (`resourceTencentCloudWafCcSessionRead`).
- Update the resource documentation (`resource_tc_waf_cc_session.md`) example to show usage of the new `key` parameter.
- Add unit test coverage in `resource_tc_waf_cc_session_test.go` for the new `key` parameter.

## Capabilities

### New Capabilities
- `waf-cc-session-key`: Support for configuring the precise-match session key (`Key`) parameter on the `tencentcloud_waf_cc_session` resource.

### Modified Capabilities
<!-- None. This is a pure addition of an optional parameter. -->

## Impact

- **Affected code**:
  - `tencentcloud/services/waf/resource_tc_waf_cc_session.go` — schema, create, read, update
  - `tencentcloud/services/waf/resource_tc_waf_cc_session_test.go` — test cases
  - `tencentcloud/services/waf/resource_tc_waf_cc_session.md` — documentation example
- **APIs**: TencentCloud WAF `v20180125` — `UpsertSession` (create/update) and `DescribeSession` (read)
- **Dependencies**: None new; the SDK field already exists in `vendor/github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/waf/v20180125/models.go`
- **Backward compatibility**: Fully backward compatible — new optional parameter, no changes to existing schema fields or state.
