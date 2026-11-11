# dcx-bandwidth-update Specification

## Purpose
TBD - created by archiving change add-dcx-bandwidth-update-param. Update Purpose after archive.
## Requirements
### Requirement: Bandwidth is updatable in place
The `tencentcloud_dcx` resource SHALL support updating the `bandwidth` parameter in place (without forcing resource recreation) by removing `ForceNew: true` from the `bandwidth` schema field. The field SHALL remain `TypeInt`, `Optional`, and `Computed`.

#### Scenario: Update bandwidth without recreating the resource
- **WHEN** a user changes the `bandwidth` value in an existing `tencentcloud_dcx` resource configuration
- **THEN** the provider SHALL NOT destroy and recreate the resource
- **AND** the provider SHALL call `ModifyDirectConnectTunnelAttribute` with the new `Bandwidth` value

#### Scenario: Bandwidth schema field attributes
- **WHEN** the `bandwidth` schema field is defined in `ResourceTencentCloudDcxInstance()`
- **THEN** the field SHALL be `TypeInt`, `Optional`, `Computed`
- **AND** the field SHALL NOT have `ForceNew: true`

### Requirement: Bandwidth change handled in Update function
The `resourceTencentCloudDcxInstanceUpdate` function SHALL detect changes to `bandwidth` via `d.HasChange("bandwidth")` and call the `ModifyDirectConnectTunnelAttribute` API with the `Bandwidth` request parameter set from the schema value and `DirectConnectTunnelId` set from `d.Id()`.

#### Scenario: Bandwidth change triggers ModifyDirectConnectTunnelAttribute
- **WHEN** `d.HasChange("bandwidth")` returns true in the Update function
- **THEN** the provider SHALL build a `ModifyDirectConnectTunnelAttributeRequest`
- **AND** SHALL set `DirectConnectTunnelId` from `d.Id()`
- **AND** SHALL set `Bandwidth` from the current schema value (`d.Get("bandwidth")`)
- **AND** SHALL call `ModifyDirectConnectTunnelAttribute` wrapped in `resource.Retry(tccommon.WriteRetryTimeout, ...)` with `tccommon.RetryError(e)` error wrapping

#### Scenario: Bandwidth unchanged does not trigger API call
- **WHEN** `d.HasChange("bandwidth")` returns false in the Update function
- **THEN** the provider SHALL NOT call `ModifyDirectConnectTunnelAttribute` for bandwidth

#### Scenario: Bandwidth update failure is retried
- **WHEN** the `ModifyDirectConnectTunnelAttribute` API call fails with a retryable error
- **THEN** the provider SHALL retry within `tccommon.WriteRetryTimeout`
- **AND** SHALL wrap the error using `tccommon.RetryError(e)`

### Requirement: Read refreshes bandwidth from Describe API
The `resourceTencentCloudDcxInstanceRead` function SHALL continue to read `Bandwidth` from the `DescribeDirectConnectTunnels` API response (`DirectConnectTunnel.Bandwidth`) and set it in Terraform state when the value is not nil.

#### Scenario: Read populates bandwidth after update
- **WHEN** the provider reads an existing `tencentcloud_dcx` resource after a bandwidth update
- **THEN** `bandwidth` SHALL be refreshed from the `DirectConnectTunnel.Bandwidth` field of the Describe API response
