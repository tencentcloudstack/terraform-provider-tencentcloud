# Change: Fix GWLB Target Group Port Field Drift

## Why
The `port` field in `tencentcloud_gwlb_target_group` resource causes drift because the API returns a default value even when the user does not provide it. Currently, the `port` field lacks the `Computed: true` attribute, causing Terraform to detect configuration drift on every plan/apply operation.

## What Changes
- Add `Computed: true` attribute to the `port` field in the resource schema
- This allows Terraform to accept API-provided default values without treating them as drift

## Impact
- **Affected specs**: gwlb-target-group
- **Affected code**: `tencentcloud/services/gwlb/resource_tc_gwlb_target_group.go`
- **User impact**: Eliminates false-positive drift detection for users who don't explicitly set the `port` field
- **Breaking change**: No - this is a bug fix that improves user experience
