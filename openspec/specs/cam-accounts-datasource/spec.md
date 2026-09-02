# cam-accounts-datasource Specification

## Purpose
TBD - created by archiving change add-cam-accounts-datasource. Update Purpose after archive.
## Requirements
### Requirement: CAM accounts data source SHALL query accounts via ListAccounts
The `tencentcloud_cam_accounts` data source SHALL call the CAM `ListAccounts` API and expose the returned account list to Terraform.

#### Scenario: Query all accounts with no filters
- **WHEN** a user declares `data "tencentcloud_cam_accounts" "all"` with no optional arguments
- **THEN** the provider calls `ListAccounts` with no `MaxItems`, `Marker`, or `UserType` set and flattens `response.Response.Users` into the `users` list

#### Scenario: Query accounts filtered by user type
- **WHEN** a user sets `user_type = "SubUser"`
- **THEN** the provider sets `request.UserType = "SubUser"` so the API returns only sub-users

#### Scenario: Query accounts with paging
- **WHEN** a user sets `max_items = 50` and `marker = "<previous-marker>"`
- **THEN** the provider sets `request.MaxItems = 50` and `request.Marker = "<previous-marker>"`, and the data source outputs the updated `marker` and `is_truncated` values from the response

### Requirement: The users list SHALL flatten each account's fields without an extra nesting layer
The `users` schema SHALL be a `TypeList` whose element is a `schema.Resource` containing `uin`, `name`, `uid`, `remark`, `console_login`, `phone_num`, `country_code`, `email`, `create_time`, and `user_type` fields directly — there SHALL NOT be an intermediate wrapper block.

#### Scenario: Flattened account fields
- **WHEN** the API returns a user with `Uin=123`, `Name="alice"`, `Uid=456`
- **THEN** the `users` list element contains `uin = 123`, `name = "alice"`, `uid = 456` as top-level attributes of the element

### Requirement: Each account field SHALL be set only when the API response value is non-nil
Before calling `d.Set(...)` / building the map for each `users` element, the provider SHALL check whether the corresponding `ListAllUser` field is nil and skip setting it when nil.

#### Scenario: Nil remark is skipped
- **WHEN** the API returns a user whose `Remark` field is nil
- **THEN** the `remark` attribute is not set in that element's map

### Requirement: Paging outputs marker and is_truncated SHALL be exposed
The data source SHALL expose `marker` (string, Optional + Computed) and `is_truncated` (bool, Computed) reflecting `response.Response.Marker` and `response.Response.IsTruncated`.

#### Scenario: Truncated result reports next marker
- **WHEN** the API returns `IsTruncated = true` and `Marker = "abc"`
- **THEN** the data source outputs `is_truncated = true` and `marker = "abc"`

#### Scenario: Non-truncated result
- **WHEN** the API returns `IsTruncated = false` and no `Marker`
- **THEN** the data source outputs `is_truncated = false` and `marker` remains empty

### Requirement: The read handler SHALL retry on transient errors and not clear the id on empty response
The read handler SHALL wrap the `ListAccounts` call in `resource.Retry(tccommon.ReadRetryTimeout, ...)`. Inside the retry block, if the response or `Users` list is empty, the handler SHALL return `NonRetryableError` instead of clearing the Terraform id, so transient API blips do not wipe local state.

#### Scenario: Transient empty response triggers retry
- **WHEN** the `ListAccounts` call returns an empty `Users` list due to a transient API issue
- **THEN** the retry block returns `NonRetryableError`, causing the outer retry to keep retrying until the timeout is exhausted, and the provider logs `[DATASOURCE] read empty, skip SetId` rather than clearing the id

### Requirement: The data source SHALL be registered in the provider
The provider SHALL register `tencentcloud_cam_accounts` mapped to `cam.DataSourceTencentCloudCamAccounts()` in `provider.go`.

#### Scenario: Provider registration
- **WHEN** the provider is initialized
- **THEN** `tencentcloud_cam_accounts` is available as a data source

### Requirement: Optional input parameters SHALL NOT enforce SDK enum ranges as validate funcs
The `user_type` input parameter SHALL be a plain `schema.TypeString` with no `ValidateFunc`. The provider SHALL pass the user-supplied value directly to the SDK without validating it against the documented enum range.

#### Scenario: User supplies a user_type value
- **WHEN** a user sets `user_type = "Collaborator"`
- **THEN** the provider forwards `"Collaborator"` to `request.UserType` without schema-level validation

