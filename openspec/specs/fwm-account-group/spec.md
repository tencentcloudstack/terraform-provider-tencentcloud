### Requirement: FWM account group resource lifecycle

The system SHALL provide a Terraform resource `tencentcloud_fwm_account_group` that manages the full lifecycle (create, read, update, delete, import) of an FWM account group via the cloud APIs CreateAccountGroup, DescribeAccountGroups, ModifyAccountGroup, and DeleteAccountGroup (SDK package `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/fwm/v20250611`).

#### Scenario: Create account group synchronously
- **WHEN** a user applies a configuration defining `tencentcloud_fwm_account_group` with `group_name` (and optional `description`)
- **THEN** the provider calls CreateAccountGroup, and on success uses the returned `GroupId` (e.g. `acg-xxxxxxxx`) as the resource id, then calls Read to refresh state
- **AND** the create operation MUST be synchronous (no polling/waitReady), and MUST fail with NonRetryableError when the response or GroupId is empty

#### Scenario: Read uses list with GroupId filter
- **WHEN** the provider reads an account group by its id
- **THEN** it calls DescribeAccountGroups with `Filters=[{Name:"GroupId",Values:[<id>],OperatorType:1}]` and takes the first element of `AccountGroups`
- **AND** when `AccountGroups` is empty the provider MUST treat the group as deleted (`d.SetId("")`), without erroring or retrying

#### Scenario: Update group_name in place
- **WHEN** the user changes `group_name`
- **THEN** the provider calls ModifyAccountGroup with `GroupId` and the new `GroupName`, then calls Read
- **AND** the change MUST be applied in place without recreating the resource

#### Scenario: description change forces recreation
- **WHEN** the user changes `description`
- **THEN** the provider MUST recreate the resource (description is ForceNew), because ModifyAccountGroup does not accept Description

#### Scenario: Delete propagates business rejections
- **WHEN** the user destroys the resource
- **THEN** the provider calls DeleteAccountGroup with `GroupId`
- **AND** InvalidParameter errors (member_count>0 or referenced by an active policy) MUST be propagated to the user as-is, prompting cleanup, NOT treated as success
- **AND** ResourceNotFound (group already absent) is acceptable for destroy-friendly behavior

### Requirement: FWM account group schema

The resource schema SHALL define: `group_name` (Required, mutable), `description` (Optional, ForceNew), and computed read-only fields `group_id`, `member_count`, `can_delete`, `create_time`. Members MUST NOT be inlined (managed by a separate resource).

#### Scenario: Required and ForceNew attributes
- **WHEN** the schema is defined
- **THEN** `group_name` MUST be Required and mutable, `description` MUST be Optional with ForceNew=true

#### Scenario: Computed read-only attributes
- **WHEN** the schema is defined
- **THEN** `group_id`, `member_count`, `can_delete`, `create_time` MUST be Computed and read-only
- **AND** `can_delete` MUST NOT be writable or used as the authoritative deletion predicate

#### Scenario: Import by GroupId
- **WHEN** a user runs `terraform import tencentcloud_fwm_account_group.example acg-1a2b3c4d`
- **THEN** the provider imports using the GroupId directly (ImportStatePassthrough), and the id is NOT a composite id

### Requirement: FWM connectivity client

The provider SHALL register an FWM SDK client accessor `UseFwmV20250611Client()` in `tencentcloud/connectivity/client.go`, backed by the vendored `fwm/v20250611` package, and register the resource in `tencentcloud/provider.go`.

#### Scenario: FWM client accessor available
- **WHEN** a CRUD function obtains the client via `meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseFwmV20250611Client()`
- **THEN** it MUST return a cached `*fwmv20250611.Client` initialized with the provider credential and region

#### Scenario: Resource registered in provider
- **WHEN** the provider resource map is built
- **THEN** `"tencentcloud_fwm_account_group"` MUST map to `fwm.ResourceTencentCloudFwmAccountGroup()`

### Requirement: FWM account group unit tests with mocks

The new resource MUST ship a `resource_tc_fwm_account_group_test.go` that uses gomonkey to mock the FWM client methods (NOT the Terraform acceptance test suite), covering create/read/update/delete business logic.

#### Scenario: Mocked CRUD tests
- **WHEN** the unit tests run
- **THEN** they MUST mock `CreateAccountGroupWithContext`, `DescribeAccountGroupsWithContext`, `ModifyAccountGroupWithContext`, and `DeleteAccountGroupWithContext` on the FWM client
- **AND** assert the resource schema/data behavior without requiring real cloud credentials
