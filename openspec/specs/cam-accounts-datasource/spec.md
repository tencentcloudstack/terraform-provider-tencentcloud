# cam-accounts-datasource Specification

## Purpose
Provides a Terraform data source `tencentcloud_cam_accounts` that queries the full list of CAM accounts (Owner, SubUser, CICUser, WechatCorpUser, AgentIdentity, Collaborator, MessageReceiver) via the CAM `ListAccounts` API. Pagination is handled automatically in the service layer (up to 100 records per call, looping on the response's `IsTruncated` field and `Marker` until all pages are fetched), and no paging-related parameters are exposed in the schema.
## Requirements
### Requirement: CAM accounts data source SHALL query accounts via ListAccounts
The `tencentcloud_cam_accounts` data source SHALL call the CAM `ListAccounts` API and expose the returned account list to Terraform.

#### Scenario: Query all accounts with no filters
- **WHEN** a user declares `data "tencentcloud_cam_accounts" "all"` with no optional arguments
- **THEN** the provider calls `ListAccounts` without a `UserType` filter (each page requested with `MaxItems = 100`) and flattens the collected `response.Response.Users` into the `users` list

#### Scenario: Query accounts filtered by user type
- **WHEN** a user sets `user_type = "SubUser"`
- **THEN** the provider sets `request.UserType = "SubUser"` so the API returns only sub-users

### Requirement: Pagination SHALL be handled automatically in the service layer and SHALL NOT be exposed in the schema
The service method `DescribeCamAccountsByFilter` SHALL fetch all pages: it SHALL request at most 100 records per call (the `ListAccounts` `MaxItems` upper bound) and, based on the response's `IsTruncated` field, continue calling the API with the returned `Marker` until `IsTruncated` is false. The data source schema SHALL NOT expose any paging-related parameters (`max_items`, `marker`, or `is_truncated`).

#### Scenario: Truncated response triggers a follow-up request
- **WHEN** a `ListAccounts` call returns `IsTruncated = true` with `Marker = "abc"`
- **THEN** the service issues another `ListAccounts` request with `request.Marker = "abc"`, appends the new `Users` to the result, and repeats until `IsTruncated = false`

#### Scenario: Non-truncated response ends the loop
- **WHEN** a `ListAccounts` call returns `IsTruncated = false`
- **THEN** the service stops looping and returns the accumulated account list

#### Scenario: Schema exposes no paging parameters
- **WHEN** the data source schema is inspected
- **THEN** it contains only `user_type`, `users`, and `result_output_file` (plus the `users` element fields), and does not contain `max_items`, `marker`, or `is_truncated`

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

### Requirement: The read handler SHALL retry on transient errors
The read handler SHALL wrap the `DescribeCamAccountsByFilter` call in `resource.Retry(tccommon.ReadRetryTimeout, ...)`. An empty `Users` list SHALL NOT be treated as an error: the handler proceeds and produces an empty `users` list, consistent with the existing CAM data sources.

#### Scenario: Transient API error triggers retry
- **WHEN** the `ListAccounts` call fails with a transient error
- **THEN** the retry block retries the call until it succeeds or `tccommon.ReadRetryTimeout` is exhausted, and the provider logs `[DATASOURCE] read empty, skip SetId` before returning the error

#### Scenario: Empty account list completes normally
- **WHEN** the API returns no users
- **THEN** the read handler completes without error and the data source reports an empty `users` list

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
