## Why

The `tencentcloud_dcx` resource exposes a `bandwidth` parameter, but it is marked `ForceNew: true`, so any change to the bandwidth destroys and recreates the dedicated tunnel. The DC `ModifyDirectConnectTunnelAttribute` API now accepts a `Bandwidth` parameter (`*int64`, "专用通道带宽值，单位为M"), so bandwidth can be updated in place without recreating the resource. Exposing this update path avoids unnecessary disruption to dedicated tunnels when users only need to adjust the bandwidth.

## What Changes

- Remove `ForceNew: true` from the `bandwidth` schema field of the `tencentcloud_dcx` resource so that bandwidth becomes an updatable parameter (still `Optional`, `Computed`, `TypeInt`).
- In the `resourceTencentCloudDcxInstanceUpdate` function, handle `bandwidth` changes by calling `ModifyDirectConnectTunnelAttribute` with the `Bandwidth` request parameter (alongside the existing `name` update path).
- The Read function already reads `Bandwidth` from the Describe API response (`DirectConnectTunnel.Bandwidth`), so no Read changes are needed.
- Update the resource documentation (`tencentcloud/services/dc/resource_tc_dcx.md`) to reflect that `bandwidth` is now updatable.

## Capabilities

### New Capabilities
- `dcx-bandwidth-update`: Enable in-place update of the `bandwidth` parameter on the `tencentcloud_dcx` resource via the `ModifyDirectConnectTunnelAttribute` API, removing the previous `ForceNew` behavior.

### Modified Capabilities
<!-- No existing specs require modification -->

## Impact

- **Affected files:**
  - `tencentcloud/services/dc/resource_tc_dcx.go` — remove `ForceNew: true` from the `bandwidth` schema field; add `bandwidth` change handling in the Update function that calls `ModifyDirectConnectTunnelAttribute` with `Bandwidth`
  - `tencentcloud/services/dc/resource_tc_dcx.md` — update documentation to reflect that `bandwidth` is updatable
- **SDK dependency:** No SDK update needed — `ModifyDirectConnectTunnelAttributeRequest` in `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dc/v20180410` already includes the `Bandwidth *int64` field.
- **API constraints:** `Bandwidth` is accepted by `ModifyDirectConnectTunnelAttribute` (Update). The Create API `CreateDirectConnectTunnel` also accepts `Bandwidth`, so the Create path continues to work unchanged. The Describe API (`DescribeDirectConnectTunnels`) returns `Bandwidth` in the `DirectConnectTunnel` struct, so Read continues to refresh the value.
- **Backward compatibility:** Fully backward compatible. Existing configurations that set `bandwidth` continue to work; the only behavioral change is that updating `bandwidth` no longer forces resource recreation but instead updates in place.
