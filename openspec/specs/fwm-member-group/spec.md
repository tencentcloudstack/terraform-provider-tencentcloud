# fwm-member-group Specification

## Purpose
TBD - created by archiving change add-fwm-member-group. Update Purpose after archive.
## Requirements
### Requirement: fwm member group config resource schema
The system SHALL provide a Terraform resource `tencentcloud_fwm_member_group` of config kind whose schema exposes `member_id` (required), `app_id` (required, ForceNew), and `group_id` (optional, single-value string, empty = ungrouped). The resource description SHALL be in Chinese and field descriptions in English. The resource MUST NOT expose drift-prone derived fields (nickname/uin/role/cfw_managed/etc.) as comparable arguments.

#### Scenario: Schema fields and types
- **WHEN** the resource schema is defined
- **THEN** `member_id` is `schema.TypeString` Required, `app_id` is `schema.TypeString` Required with `ForceNew: true`, and `group_id` is `schema.TypeString` Optional (no ForceNew, clearable to empty string)

#### Scenario: group_id is single value not list
- **WHEN** the cloud API returns `AccountGroup` as a single object `*AccountGroupInfo`
- **THEN** the resource SHALL model `group_id` as a single string, NOT a list, because the DB unique index `uk_admin_appid_appid(admin_appid,appid)` enforces at most one group per member

### Requirement: config-kind CRUD with empty create and delete
The system SHALL implement Create as an empty operation that calls `d.SetId(member_id)` (using the composite `app_id#member_id` id) and then delegates to Update. The system SHALL implement Delete as a no-op returning nil, because member-group membership cannot be destroyed (it can only be cleared by setting `group_id` to empty). Only Read and Update SHALL contain real logic.

#### Scenario: Create sets id and delegates to update
- **WHEN** Terraform creates the resource
- **THEN** the handler computes the composite id `app_id#member_id`, calls `d.SetId(...)`, and invokes the Update handler to apply the initial `group_id`

#### Scenario: Delete is a no-op
- **WHEN** Terraform destroys the resource
- **THEN** the handler returns nil without calling any cloud API, leaving the membership relation untouched (clearing must be done via `group_id = ""`)

### Requirement: read locates member via AppId filter
The Read handler SHALL locate the target member by calling `DescribeOrganMembers` with `Filters=[{Name:"AppId",Values:[appId],OperatorType:1}]` (1=EQUAL), because the filter mapping does NOT support `MemberId`. After receiving `Members`, the handler SHALL match the entry whose `MemberId` equals the target member id.

#### Scenario: Successful read of a grouped member
- **WHEN** the member exists and belongs to a group (`AccountGroup != nil`)
- **THEN** the handler SHALL set `group_id` to `AccountGroup.GroupId` and MUST NOT clear the resource id

#### Scenario: Read of an ungrouped member
- **WHEN** the member exists but `AccountGroup` is null (ungrouped)
- **THEN** the handler SHALL set `group_id` to empty string and MUST NOT treat the resource as deleted (no `d.SetId("")`)

#### Scenario: Member not found
- **WHEN** no entry in `Members` matches the target `MemberId`
- **THEN** the handler SHALL print a log with the current id and then call `d.SetId("")`

### Requirement: update via ModifyMemberGroup with MemberId
The Update handler SHALL call `ModifyMemberGroup` with `UserMemberId` set to the member's `MemberId` (NOT the AppId) and `GroupId` set to the target group id or empty string (to move out). The handler SHALL skip the write API call when the target `group_id` equals the current value. The operation is synchronous (single transaction, no task id), so the handler SHALL NOT poll or waitReady afterward.

#### Scenario: Move member into a group
- **WHEN** `group_id` is set to a non-empty group id and differs from current
- **THEN** the handler SHALL call `ModifyMemberGroup` with `UserMemberId=<member_id>`, `GroupId=<group_id>`

#### Scenario: Move member out of a group
- **WHEN** `group_id` is set to empty string and member was grouped
- **THEN** the handler SHALL call `ModifyMemberGroup` with `UserMemberId=<member_id>`, `GroupId=""`, and the server removes the relation (soft delete + member_count-1)

#### Scenario: Skip update when unchanged
- **WHEN** the target `group_id` equals the current value
- **THEN** the handler SHALL NOT call `ModifyMemberGroup`

### Requirement: retry and error handling
The Read handler SHALL wrap the `DescribeOrganMembers` call in `resource.Retry(tccommon.ReadRetryTimeout, ...)` and convert errors via `tccommon.RetryError`. The Update handler SHALL wrap the `ModifyMemberGroup` call in `resource.Retry(tccommon.WriteRetryTimeout, ...)` and convert errors via `tccommon.RetryError`. Cloud API errors (ResourceNotFound / InvalidParameter) SHALL be propagated as-is without silently mutating state.

#### Scenario: Read retry on transient error
- **WHEN** `DescribeOrganMembers` returns a transient error
- **THEN** the retry wrapper SHALL retry within `ReadRetryTimeout` and surface the final error

### Requirement: composite id and import
The resource id SHALL be the composite string `app_id` + `#` (`tccommon.FILED_SP`) + `member_id`. The resource SHALL support import via `schema.ImportStatePassthrough`; users import using the two-segment id `app_id#member_id` because Read needs app_id to locate the member. The Read, Update handlers SHALL parse `app_id` and `member_id` from `d.Id()`.

#### Scenario: Import two-segment id
- **WHEN** a user runs `terraform import tencentcloud_fwm_member_group.foo "app_id#member_id"`
- **THEN** the resource SHALL populate state from the parsed app_id and member_id by running Read

### Requirement: provider registration and client onboarding
The system SHALL register `tencentcloud_fwm_member_group` in `tencentcloud/provider.go` and `tencentcloud/provider.md`. The system SHALL add `UseFwmClient()` and a `fwmConn *fwm.Client` field to `tencentcloud/connectivity/client.go` mirroring `UseCfwClient`, importing `fwm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/fwm/v20250611"`.

#### Scenario: Resource registered in provider
- **WHEN** the provider initializes
- **THEN** the resources map SHALL contain `"tencentcloud_fwm_member_group"` mapped to `fwm.ResourceTencentCloudFwmMemberGroup()`

### Requirement: documentation and test
The system SHALL create `resource_tc_fwm_member_group.md` with a one-line description (mentioning FWM), Example Usage (with `depends_on` pointing to `tencentcloud_fwm_cfw_managed` and `group_id` referencing a group id), and an Import section explaining the two-segment id. The system SHALL create `resource_tc_fwm_member_group_test.go` using gomonkey mock (not TF suite) for this new resource. The `website/docs/` doc is generated by `make doc`.

#### Scenario: HCL example reflects managed-member prerequisite
- **WHEN** the example sets a non-empty `group_id`
- **THEN** the example SHALL include `depends_on = [tencentcloud_fwm_cfw_managed.<name>]` and the member group resource description SHALL document the reverse coupling with cfw_managed

