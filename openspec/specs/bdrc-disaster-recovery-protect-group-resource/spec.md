# bdrc-disaster-recovery-protect-group-resource Specification

## Requirements

### Requirement: Resource MUST be registered as `tencentcloud_bdrc_disaster_recovery_protect_group`

The provider SHALL register a new general-type resource named `tencentcloud_bdrc_disaster_recovery_protect_group` whose Create / Read / Update / Delete callbacks invoke the BDRC `v20260330` APIs (`CreateDisasterRecoveryProtectGroup`, `DescribeDisasterRecoveryProtectGroups`, `ModifyProtectGroupAttribute`, `DeleteDisasterRecoveryProtectGroups`).

#### Scenario: Resource registered in provider map

- **WHEN** the provider is loaded
- **THEN** `provider.go` exposes the resource via key `"tencentcloud_bdrc_disaster_recovery_protect_group"` mapped to `bdrc.ResourceTencentCloudBdrcDisasterRecoveryProtectGroup()`.

#### Scenario: Resource appears in gendoc index

- **WHEN** `tencentcloud/provider.md` is scanned by `make doc`
- **THEN** the Business Disaster Recovery Center(BDRC) Resource section MUST include `tencentcloud_bdrc_disaster_recovery_protect_group` so that `website/docs/r/bdrc_disaster_recovery_protect_group.html.markdown` is generated.

### Requirement: BDRC client method MUST be added to connectivity layer

The provider SHALL add a `UseBdrcV20260330Client()` method to `TencentCloudClient` in `tencentcloud/connectivity/client.go` that lazily constructs and caches a `*bdrcv20260330.Client`, mirroring the `UseIgtmV20231024Client()` pattern.

#### Scenario: Client method exists and is reusable

- **WHEN** any resource callback calls `meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseBdrcV20260330Client()`
- **THEN** a non-nil `*bdrcv20260330.Client` is returned, and repeated calls return the same cached instance.

### Requirement: Schema input fields MUST mirror CreateDisasterRecoveryProtectGroup API input

The resource schema SHALL declare these top-level argument keys, with semantics matching the SDK request fields of `CreateDisasterRecoveryProtectGroupRequestParams`:

| HCL key | SDK field | Type | Required | ForceNew |
|---|---|---|---|---|
| `site_pair_id` | `SitePairId` | TypeString | Yes | **Yes** |
| `protect_group_type` | `ProtectGroupType` | TypeString | Yes | **Yes** |
| `recovery_point_objective` | `RecoveryPointObjective` | TypeInt | Yes | **Yes** |
| `protect_group_name` | `ProtectGroupName` | TypeString | No | No |
| `data_direction` | `DataDirection` | TypeString | No | **Yes** |

`site_pair_id`, `protect_group_type`, `recovery_point_objective`, and `data_direction` are `ForceNew` because `ModifyProtectGroupAttribute` only supports renaming. `protect_group_name` is updatable.

#### Scenario: Required fields enforce on plan

- **WHEN** the user writes a config that omits `site_pair_id`, `protect_group_type`, or `recovery_point_objective`
- **THEN** `terraform plan` SHALL fail validation pointing at the missing required attribute.

#### Scenario: Changing a ForceNew field forces replacement

- **GIVEN** state has `site_pair_id = "sp-123"`
- **WHEN** the user changes `site_pair_id` to `"sp-456"` in HCL
- **THEN** Terraform's plan reports a destroy + create cycle.

### Requirement: Schema computed fields MUST mirror ProtectGroup response (flattened, no list-wrapper layer)

The resource schema SHALL declare these top-level Computed keys sourced from the `ProtectGroup` struct returned by `DescribeDisasterRecoveryProtectGroups`. The list wrapper `ProtectGroupSet` MUST NOT be exposed as a schema key; instead each element field is flattened to the top level (the query returns a single element by ID).

| HCL key | SDK field | Type | Computed |
|---|---|---|---|
| `app_id` | `AppId` | TypeInt | Yes |
| `site_pair_name` | `SitePairName` | TypeString | Yes |
| `source_region` | `SourceRegion` | TypeString | Yes |
| `source_zone` | `SourceZone` | TypeString | Yes |
| `source_vpc` | `SourceVpc` | TypeString | Yes |
| `target_region` | `TargetRegion` | TypeString | Yes |
| `target_zone` | `TargetZone` | TypeString | Yes |
| `target_vpc` | `TargetVpc` | TypeString | Yes |
| `copy_type` | `CopyType` | TypeString | Yes |
| `disaster_recovery_type` | `DisasterRecoveryType` | TypeString | Yes |
| `peer_cloud_name` | `PeerCloudName` | TypeString | Yes |
| `create_from` | `CreateFrom` | TypeString | Yes |
| `life_state` | `LifeState` | TypeString | Yes |
| `account_uin` | `AccountUin` | TypeString | Yes |
| `sub_account_uin` | `SubAccountUin` | TypeString | Yes |
| `create_time` | `CreateTime` | TypeString | Yes |
| `modify_time` | `ModifyTime` | TypeString | Yes |
| `bind_protected_resource_count` | `BindProtectedResourceCount` | TypeInt | Yes |
| `error_recovery_point_objective_count` | `ErrorRecoveryPointObjectiveCount` | TypeInt | Yes |
| `protected_resource_status_set` | `ProtectedResourceStatusSet` | TypeList | Yes |

`protected_resource_status_set` is a `TypeList` whose `Elem` is a `schema.Resource` with Computed keys `status` (TypeString) and `count` (TypeInt), mirroring `ProtectedResourceStatus`.

#### Scenario: No list-wrapper schema key

- **WHEN** a code reviewer inspects the schema map
- **THEN** there is NO key named `protect_group_set` or `protect_group_list`; every response field is a top-level scalar or the flattened `protected_resource_status_set` list.

### Requirement: Resource ID MUST be ProtectGroupId (single field, importable)

After Create, the Terraform resource ID SHALL be the value of `response.Response.ProtectGroupId`. No compound separator is used. The resource SHALL declare `Importer: &schema.ResourceImporter{State: schema.ImportStatePassthrough}` so users can `terraform import tencentcloud_bdrc_disaster_recovery_protect_group.xxx <protect_group_id>`.

#### Scenario: ID is set from the API response

- **GIVEN** Create returns `ProtectGroupId = "pg-abc"`
- **WHEN** Create completes
- **THEN** `d.Id()` returns `"pg-abc"`.

#### Scenario: Import by ProtectGroupId

- **WHEN** the user runs `terraform import tencentcloud_bdrc_disaster_recovery_protect_group.foo pg-abc`
- **THEN** the resource is imported with `d.Id() == "pg-abc"` and a subsequent Read populates all fields from the Describe API.

### Requirement: Create MUST call CreateDisasterRecoveryProtectGroup with retry and nil-safe checks

The Create callback SHALL build a `CreateDisasterRecoveryProtectGroupRequest` populated from `d.Get(...)` for all five input fields, call `UseBdrcV20260330Client().CreateDisasterRecoveryProtectGroupWithContext(ctx, request)` inside `resource.Retry(tccommon.WriteRetryTimeout, ...)`, and guard `result == nil || result.Response == nil` with `resource.NonRetryableError`. After the retry, it MUST verify `response.Response.ProtectGroupId` is non-nil (logging logId and d.Id for troubleshooting) before calling `d.SetId(...)`. Create MUST end by calling Read.

#### Scenario: Create issues the API call and sets the ID

- **GIVEN** HCL with all required fields
- **WHEN** the user runs `terraform apply`
- **THEN** exactly one `CreateDisasterRecoveryProtectGroup` request is issued with all input fields populated, and `d.Id()` returns the returned `ProtectGroupId`.

#### Scenario: Nil response is detected

- **GIVEN** the SDK returns `(result == nil, err == nil)`
- **WHEN** the retry callback runs
- **THEN** the callback returns `resource.NonRetryableError(...)` and does NOT panic.

#### Scenario: Nil ProtectGroupId is rejected

- **GIVEN** Create returns a response where `ProtectGroupId` is nil
- **WHEN** Create processes the response
- **THEN** Create returns an error mentioning `ProtectGroupId is nil` and does NOT call `d.SetId`.

#### Scenario: Transient SDK error is retried

- **GIVEN** the first invocation returns a retriable error
- **WHEN** the retry callback re-runs
- **THEN** the second attempt's response is observed; the original error surfaces only after the retry budget is exhausted.

### Requirement: Service layer MUST provide DescribeDisasterRecoveryProtectGroupById with pagination

A `BdrcService` struct SHALL be created in `tencentcloud/services/bdrc/service_tencentcloud_bdrc.go` with method `DescribeDisasterRecoveryProtectGroupById(ctx, protectGroupId, protectGroupType string) (*bdrcv20260330.ProtectGroup, error)`. The method SHALL call `DescribeDisasterRecoveryProtectGroups` with `ProtectGroupIds = [&protectGroupId]`, set `ProtectGroupType` only when non-empty, use `Limit = 100` (API max), and paginate via `Offset` until the target ID is found or pages are exhausted. The SDK call SHALL be wrapped in `resource.Retry(tccommon.ReadRetryTimeout, ...)` with `result == nil || result.Response == nil` guarded by `NonRetryableError`.

#### Scenario: Found by ID returns the element

- **GIVEN** a protect group with id `pg-1` exists
- **WHEN** `DescribeDisasterRecoveryProtectGroupById(ctx, "pg-1", "DISK")` is called
- **THEN** a non-nil `*ProtectGroup` with `ProtectGroupId == "pg-1"` is returned.

#### Scenario: Not found returns nil without error

- **GIVEN** no protect group with id `pg-missing` exists
- **WHEN** `DescribeDisasterRecoveryProtectGroupById(ctx, "pg-missing", "DISK")` is called
- **THEN** `(nil, nil)` is returned after exhausting pages.

#### Scenario: Empty type omits the filter

- **GIVEN** `protectGroupType == ""` (e.g. right after import)
- **WHEN** `DescribeDisasterRecoveryProtectGroupById(ctx, "pg-1", "")` is called
- **THEN** the request does NOT set `ProtectGroupType` (nil pointer), relying on `ProtectGroupIds` filtering.

### Requirement: Read MUST use the service layer and flatten the single element

The Read callback SHALL call `BdrcService.DescribeDisasterRecoveryProtectGroupById(ctx, d.Id(), protectGroupType)`. If `respData == nil`, it MUST first log `log.Printf("[CRUD] bdrc disaster_recovery_protect_group id=%s", d.Id())` to preserve context, then `d.SetId("")`, then return nil. Otherwise it SHALL set every input and computed field (nil-checked before each `d.Set`), including flattening `ProtectedResourceStatusSet` into `protected_resource_status_set`.

#### Scenario: Resource not found clears state with logging

- **GIVEN** the Describe returns nil for an existing state id `pg-x`
- **WHEN** Read runs
- **THEN** a `[CRUD]` log line containing `id=pg-x` is emitted, `d.SetId("")` is called, and Read returns nil.

#### Scenario: Fields are populated from response

- **GIVEN** Describe returns a `ProtectGroup` with `ProtectGroupName = "foo"`, `LifeState = "NORMAL"`, and one `ProtectedResourceStatus{Status:"AVAILABLE", Count:2}`
- **WHEN** Read runs
- **THEN** state has `protect_group_name = "foo"`, `life_state = "NORMAL"`, and `protected_resource_status_set = [{status="AVAILABLE", count=2}]`.

### Requirement: Update MUST reject immutable fields and only rename via ModifyProtectGroupAttribute

The Update callback SHALL define `immutableArgs = ["site_pair_id", "protect_group_type", "recovery_point_objective", "data_direction"]`. For each, if `d.HasChange(v)`, Update SHALL return an error stating the field is immutable and must be recreated. If `d.HasChange("protect_group_name")`, Update SHALL call `ModifyProtectGroupAttributeWithContext` with `ProtectGroupId = d.Id()` and `ProtectGroupName` from `d.GetOk`, wrapped in `resource.Retry(tccommon.WriteRetryTimeout, ...)` with nil-safe guards. Update MUST end by calling Read.

#### Scenario: Changing an immutable field errors

- **GIVEN** state has `protect_group_type = "DISK"`
- **WHEN** the user changes `protect_group_type` to `"INSTANCE"` in HCL (without ForceNew catching it)
- **THEN** Update returns an error mentioning `protect_group_type` is immutable and cannot be updated.

#### Scenario: Renaming triggers ModifyProtectGroupAttribute

- **GIVEN** state has `protect_group_name = "old"`
- **WHEN** the user changes it to `"new"`
- **THEN** Update issues exactly one `ModifyProtectGroupAttribute` request with `ProtectGroupId = <id>` and `ProtectGroupName = "new"`.

#### Scenario: No change skips the API call

- **GIVEN** no schema field changed
- **WHEN** Update runs
- **THEN** no `ModifyProtectGroupAttribute` request is issued and Read is called.

### Requirement: Delete MUST call DeleteDisasterRecoveryProtectGroups with retry

The Delete callback SHALL build a `DeleteDisasterRecoveryProtectGroupsRequest` with `ProtectGroups = [&d.Id()]`, call `DeleteDisasterRecoveryProtectGroupsWithContext` inside `resource.Retry(tccommon.WriteRetryTimeout, ...)`, guard `result == nil || result.Response == nil`, and return nil on success.

#### Scenario: Delete issues the API call

- **GIVEN** state has id `pg-del`
- **WHEN** the user runs `terraform destroy`
- **THEN** exactly one `DeleteDisasterRecoveryProtectGroups` request is issued with `ProtectGroups = ["pg-del"]`.

#### Scenario: Nil response is detected

- **GIVEN** the SDK returns `(result == nil, err == nil)`
- **WHEN** the retry callback runs
- **THEN** the callback returns `resource.NonRetryableError(...)` and does NOT panic.

### Requirement: Every API call MUST be wrapped in resource.Retry

All four SDK invocations (Create / Describe / Modify / Delete) SHALL be wrapped in `resource.Retry(...)` with the appropriate timeout (`tccommon.WriteRetryTimeout` for Create/Modify/Delete, `tccommon.ReadRetryTimeout` for Describe) and forward errors via `tccommon.RetryError(e)`.

#### Scenario: Retry wrapper present on all call sites

- **WHEN** a code reviewer inspects each CRUD callback and the service layer method
- **THEN** every SDK call is inside a `resource.Retry(...)` block, not a bare invocation.

### Requirement: Documentation and unit test MUST follow project conventions

- The HCL example markdown SHALL live at `tencentcloud/services/bdrc/resource_tc_bdrc_disaster_recovery_protect_group.md` and contain a one-line summary (mentioning BDRC), an `Example Usage` HCL block, and an `Import` section explaining import uses `ProtectGroupId`.
- A unit test SHALL live at `tencentcloud/services/bdrc/resource_tc_bdrc_disaster_recovery_protect_group_test.go` (package `bdrc_test`) using gomonkey to mock the BDRC client methods (NOT the terraform acceptance test suite), covering Create → Read → Update (rename) → Delete business logic.
- Running `make doc` SHALL regenerate `website/docs/r/bdrc_disaster_recovery_protect_group.html.markdown`.

#### Scenario: Generated website doc lists the resource

- **WHEN** `make doc` runs
- **THEN** `website/docs/r/bdrc_disaster_recovery_protect_group.html.markdown` exists and lists every schema attribute defined above.

#### Scenario: Unit test uses gomonkey mocks

- **WHEN** the test file is opened
- **THEN** the package is `bdrc_test`, gomonkey is used to patch `UseBdrcV20260330Client` and/or the client methods, and NO `TF_ACC` acceptance test step is present.
