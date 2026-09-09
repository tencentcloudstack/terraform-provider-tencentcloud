## Why

The `tencentcloud_dlc_user` resource currently does not expose the `AccountType` parameter, which is supported by the DLC (Data Lake Compute) cloud API in all CRUD operations (`CreateUser`, `DescribeUsers`, `ModifyUser`, `DeleteUser`). This parameter allows users to specify the account type (e.g., `UserAccount` for user accounts, `RoleAccount` for role accounts) when managing DLC users. Without this parameter, Terraform users cannot configure or track the account type, limiting their ability to fully manage DLC user resources.

## What Changes

- Add a new optional `account_type` parameter (string type) to the `tencentcloud_dlc_user` resource schema
- Include `AccountType` in the `CreateUser` API request during resource creation
- Read `AccountType` from the `DescribeUsers` API response (`UserInfo.AccountType`) and set it in Terraform state during resource read
- Include `AccountType` in the `ModifyUser` API request during resource update (support in-place updates)
- Include `AccountType` in the `DeleteUser` API request during resource deletion (passed through the service layer)
- Update the service layer methods (`DescribeDlcUserById`, `DeleteDlcUserById`) to accept and pass the `AccountType` parameter
- Update the resource documentation (`resource_tc_dlc_user.md`) with the new parameter
- Add unit test coverage for the new parameter

## Capabilities

### New Capabilities

- `dlc-user-account-type`: Adds support for the `account_type` parameter in the `tencentcloud_dlc_user` resource, enabling users to specify and track the DLC user account type across create, read, update, and delete operations.

### Modified Capabilities

None. There are no existing specs for `tencentcloud_dlc_user`, so this is a pure addition.

## Impact

### Affected Code
- `tencentcloud/services/dlc/resource_tc_dlc_user.go` - Add `account_type` to schema, CRUD operations
- `tencentcloud/services/dlc/service_tencentcloud_dlc.go` - Update `DescribeDlcUserById` and `DeleteDlcUserById` to support `AccountType`
- `tencentcloud/services/dlc/resource_tc_dlc_user_test.go` - Add test coverage for the new parameter
- `tencentcloud/services/dlc/resource_tc_dlc_user.md` - Update documentation

### API References
- **DLC API** (`dlc/v20210125`):
  - `CreateUser`: Accepts `AccountType` as an optional input parameter
  - `DescribeUsers`: Returns `AccountType` in `UserInfo` within `UserSet`
  - `ModifyUser`: Accepts `AccountType` as an optional input parameter
  - `DeleteUser`: Accepts `AccountType` as an optional input parameter

### Backward Compatibility
- **Fully Compatible**: The new `account_type` parameter is optional. Existing configurations without this parameter continue to work unchanged. No state migration is required.
