## Why

The CVM `RunInstances` and `ModifyInstancesAttribute` APIs support an `EnableJumboFrame` parameter to enable/disable jumbo frames on an instance, but the Terraform resource `tencentcloud_instance` does not expose it. Users who need jumbo frame support (for supported instance types) cannot manage this setting through Terraform and must fall back to the console or API directly.

## What Changes

- Add `enable_jumbo_frame` (Optional, Computed, `TypeBool`) parameter to the `tencentcloud_instance` resource to allow users to enable or disable jumbo frames during instance creation and updates.
- Wire `enable_jumbo_frame` through the Create flow: pass it to `RunInstancesRequest.EnableJumboFrame` when set.
- Wire `enable_jumbo_frame` through the Update flow: when the field changes, call a new `ModifyInstancesAttribute`-based service method (`ModifyEnableJumboFrame`) that sets `EnableJumboFrame` on the instance.
- Wire `enable_jumbo_frame` through the Read flow: read `EnableJumboFrame` from the `DescribeInstances` response (`Instance` struct) and set it into state.

## Capabilities

### New Capabilities
- `cvm-instance-enable-jumbo-frame`: Expose the `enable_jumbo_frame` parameter on the `tencentcloud_instance` resource, covering Create (`RunInstances`), Update (`ModifyInstancesAttribute`) and Read (`DescribeInstances`) support for the jumbo frame setting.

### Modified Capabilities
<!-- No existing specs require modification -->

## Impact

- **Affected files:**
  - `tencentcloud/services/cvm/resource_tc_instance.go` — add `enable_jumbo_frame` schema field, set it in Create (`RunInstancesRequest`), read it in Read (`Instance.EnableJumboFrame`), and update it in Update via a new service call.
  - `tencentcloud/services/cvm/service_tencentcloud_cvm.go` — add a `ModifyEnableJumboFrame` service method that builds a `ModifyInstancesAttributeRequest` with `EnableJumboFrame` set and calls `ModifyInstancesAttribute`.
  - `tencentcloud/services/cvm/resource_tc_instance.md` — update the resource documentation/example to include `enable_jumbo_frame`.
- **SDK dependency:** The vendored `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312` already includes `EnableJumboFrame` on `RunInstancesRequest`, `ModifyInstancesAttributeRequest`, and the `Instance` struct returned by `DescribeInstances`. No SDK upgrade is required.
- **Backward compatibility:** fully backward compatible — the new parameter is Optional and Computed, so existing configurations and state continue to work without changes.
- **API constraints:** Only instance types that support jumbo frames accept `EnableJumboFrame=true`; the value is also returned by `DescribeInstances`, so Read can refresh it. The parameter is mutable via `ModifyInstancesAttribute` (not ForceNew).