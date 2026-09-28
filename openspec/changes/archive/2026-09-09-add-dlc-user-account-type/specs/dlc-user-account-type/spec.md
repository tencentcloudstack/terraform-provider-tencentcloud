## ADDED Requirements

### Requirement: DLC user resource supports account_type parameter

The `tencentcloud_dlc_user` resource SHALL support an optional `account_type` parameter that specifies the DLC user account type. The parameter SHALL be passed to the `CreateUser`, `DescribeUsers`, `ModifyUser`, and `DeleteUser` cloud API calls. The value SHALL be read back from the `DescribeUsers` API response and stored in Terraform state.

**Rationale**: The DLC cloud API supports the `AccountType` parameter across all user CRUD operations. Exposing it in Terraform allows users to manage user accounts and role accounts distinctly, and ensures Terraform state accurately reflects the cloud resource configuration.

#### Scenario: Create a DLC user with account_type

- **WHEN** a user creates a `tencentcloud_dlc_user` resource with `account_type` set to a valid value (e.g., `UserAccount` or `RoleAccount`)
- **THEN** the `CreateUser` API SHALL be called with the `AccountType` parameter set to the specified value
- **AND** the resource SHALL be created successfully
- **AND** the `account_type` value SHALL be read back from the `DescribeUsers` response and stored in Terraform state

#### Scenario: Create a DLC user without account_type

- **WHEN** a user creates a `tencentcloud_dlc_user` resource without specifying `account_type`
- **THEN** the `CreateUser` API SHALL be called without the `AccountType` parameter (or with it empty)
- **AND** the resource SHALL be created successfully using the API default account type
- **AND** no errors SHALL be raised due to the missing parameter

#### Scenario: Read account_type from DescribeUsers response

- **WHEN** the `DescribeUsers` API returns a `UserInfo` with `AccountType` set to a non-nil value
- **THEN** the resource read operation SHALL set `account_type` in Terraform state to the returned value
- **WHEN** the `DescribeUsers` API returns a `UserInfo` with `AccountType` set to nil
- **THEN** the resource read operation SHALL NOT call `d.Set("account_type", ...)` (skip setting the field)

#### Scenario: Update account_type in-place

- **WHEN** a user updates the `account_type` parameter on an existing `tencentcloud_dlc_user` resource
- **THEN** the `ModifyUser` API SHALL be called with the `AccountType` parameter set to the new value
- **AND** the resource SHALL be updated in-place without recreation
- **AND** the `account_type` parameter SHALL NOT be in the `immutableArgs` list

#### Scenario: Delete a DLC user with account_type

- **WHEN** a user deletes a `tencentcloud_dlc_user` resource that has `account_type` set
- **THEN** the `DeleteUser` API SHALL be called with the `AccountType` parameter passed from the resource schema
- **AND** the resource SHALL be deleted successfully

#### Scenario: Delete a DLC user without account_type

- **WHEN** a user deletes a `tencentcloud_dlc_user` resource that does not have `account_type` set
- **THEN** the `DeleteUser` API SHALL be called without the `AccountType` parameter (or with it empty)
- **AND** the resource SHALL be deleted successfully

#### Scenario: Import a DLC user with account_type

- **WHEN** a `tencentcloud_dlc_user` resource is imported
- **THEN** the `account_type` value SHALL be read from the `DescribeUsers` API response and populated in Terraform state
- **AND** the import SHALL succeed regardless of whether `account_type` was previously configured
