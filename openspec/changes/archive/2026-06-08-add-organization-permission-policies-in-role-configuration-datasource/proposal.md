## Why

Users need to query the list of permission policies attached to a role configuration in TencentCloud Organization's identity center (CIC). Currently there is no Terraform data source to retrieve this information, requiring users to use the console or CLI directly.

## What Changes

- Add a new Terraform data source `tencentcloud_organization_permission_policies_in_role_configuration` that calls the `ListPermissionPoliciesInRoleConfiguration` API to retrieve permission policies associated with a given role configuration.
- Register the new data source in the provider.
- Add corresponding documentation.

## Capabilities

### New Capabilities

- `permission-policies-in-role-configuration-datasource`: A read-only data source that queries permission policies (system and custom) attached to a specific role configuration in an organization's identity center zone.

### Modified Capabilities

(none)

## Impact

- New files in `tencentcloud/services/tco/`:
  - `data_source_tc_organization_permission_policies_in_role_configuration.go`
  - `data_source_tc_organization_permission_policies_in_role_configuration_test.go`
  - `data_source_tc_organization_permission_policies_in_role_configuration.md`
- Modified files:
  - `tencentcloud/provider.go` (register data source)
  - `tencentcloud/provider.md` (add documentation entry)
- Dependencies: uses existing `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/organization/v20210331` SDK package (already vendored).
