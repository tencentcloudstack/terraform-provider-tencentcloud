## ADDED Requirements

### Requirement: Resource registration
The provider SHALL expose a resource type named `tencentcloud_bdrc_instance_copy_pair` that manages a single Tencent Cloud BDRC CVM instance copy pair per resource block. The resource MUST be registered in `tencentcloud/provider.go` under the `bdrc` namespace. The `bdrc` service package and connectivity client binding (`UseBdrcV20260330Client`) SHALL be created in this change.

#### Scenario: Resource type is discoverable
- **WHEN** an operator runs `terraform plan` against a configuration that references `resource "tencentcloud_bdrc_instance_copy_pair" "<name>"`
- **THEN** Terraform resolves the type without an "unknown resource" error and shows the planned create.

#### Scenario: Provider compiles
- **WHEN** the codebase is built with `go build ./tencentcloud/...`
- **THEN** the build succeeds with no compilation errors related to the new resource, service, or connectivity code.

#### Scenario: Connectivity client wired
- **WHEN** a CRUD function calls `meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseBdrcV20260330Client()`
- **THEN** a non-nil `*bdrcv20260330.Client` is returned, with the client lazily initialized and cached on the `TencentCloudClient` struct.

### Requirement: Schema mirrors CreateInstanceCopyPair request
The resource schema SHALL expose every input parameter accepted by the `CreateInstanceCopyPair` API:

Top-level fields:
- `protect_group_id` (string, Required, ForceNew): the protect group ID the copy pair belongs to.
- `create_target_instance_parameters` (list, Required, ForceNew, MaxItems=1, MinItems=1): target CVM creation parameters. Each element is a nested resource mirroring `CreateInstanceModel`.
- `instance_copy_pair_name` (string, Optional, updatable): copy pair name; if omitted the API defaults to "未命名".
- `client_token` (string, Optional, ForceNew): idempotency token, max 64 ASCII chars.
- `recovery_point_objective` (int, Optional, ForceNew): user-desired RPO in minutes (currently only 15 is supported).
- `delete_target_resource` (bool, Optional): whether to delete the disaster-recovery site disk on destroy; defaults to API behavior (true).

Nested fields inside `create_target_instance_parameters`:
- `source_instance_id` (string, Required, ForceNew): source CVM ID.
- `instance_charge_type` (string, Required, ForceNew): instance billing mode.
- `placement` (list, Required, ForceNew, MaxItems=1): placement with nested `zone` (string, Required), `project_id` (int, Optional), `host_id` (string, Optional), `host_ids` (list of string, Optional), `project_name` (string, Optional).
- `image_id` (string, Required, ForceNew): image ID.
- `system_disk` (list, Required, ForceNew, MaxItems=1): system disk with nested `disk_type` (string, Optional), `disk_size` (int, Optional), `delete_with_instance` (bool, Optional).
- `instance_charge_prepaid` (list, Optional, ForceNew, MaxItems=1): prepaid settings with nested `period` (int, Required), `renew_flag` (string, Optional).
- `instance_type` (string, Optional, ForceNew): instance type.
- `data_disks` (list, Optional, ForceNew): data disk list, each with `disk_type` (string, Optional), `disk_size` (int, Optional), `delete_with_instance` (bool, Optional).
- `virtual_private_cloud` (list, Optional, ForceNew, MaxItems=1): VPC config with nested `vpc_id` (string, Required), `subnet_id` (string, Required), `subnet_name` (string, Optional), `as_vpc_gateway` (bool, Optional), `private_ip_addresses` (list of string, Optional), `vpc_name` (string, Optional), `ipv6address_count` (int, Optional).
- `internet_accessible` (list, Optional, ForceNew, MaxItems=1): public bandwidth with nested `internet_charge_type` (string, Optional), `internet_max_bandwidth_out` (int, Optional), `public_ip_assigned` (bool, Optional), `internet_service_provider` (string, Optional).
- `instance_name` (string, Optional, ForceNew): instance display name.
- `login_settings` (list, Optional, ForceNew, MaxItems=1): login settings with nested `password` (string, Optional, Sensitive), `key_ids` (list of string, Optional), `keep_image_login` (string, Optional).
- `enhanced_service` (list, Optional, ForceNew, MaxItems=1): enhanced service with nested `security_service` (list, MaxItems=1, with `enabled` bool), `monitor_service` (list, MaxItems=1, with `enabled` bool), `automation_service` (list, MaxItems=1, with `enabled` bool), `basic_service` (list, MaxItems=1, with `enabled` bool).
- `spot_price` (string, Optional, ForceNew): spot instance max bid.
- `host_name` (string, Optional, ForceNew): instance hostname.
- `user_data` (string, Optional, ForceNew): user data for the instance.
- `disaster_recover_group_ids` (list of string, Optional, ForceNew): placement group IDs.
- `stopped_mode` (string, Optional, ForceNew): shutdown billing mode.
- `copy_pair_id` (string, Optional, ForceNew): copy pair ID for drill scenarios.
- `recovery_time` (string, Optional, ForceNew): recovery time point for drill scenarios.

The resource SHALL additionally expose the following read-only computed attributes hydrated from the `DescribeCopyPairs` response (`CopyPair` struct):
- `copy_pair_ids` (list of string, Computed): created copy pair IDs from `CreateInstanceCopyPair` response.
- `copy_pair_id` (string, Computed): the copy pair ID (also stored as `d.Id()`).
- `copy_pair_name` (string, Computed): copy pair name from the describe response.
- `copy_pair_state` (string, Computed): copy pair state (INIT/RUNNING/FULL_COPYING/etc.).
- `copy_pair_type` (string, Computed): copy pair type (INSTANCE).
- `site_pair_id` (string, Computed): site pair ID.
- `site_pair_name` (string, Computed): site pair name.
- `protect_group_name` (string, Computed): protect group name.
- `source_region` (string, Computed): source region.
- `source_zone` (string, Computed): source zone.
- `source_vpc` (string, Computed): source VPC.
- `target_region` (string, Computed): target region.
- `target_zone` (string, Computed): target zone.
- `target_vpc` (string, Computed): target VPC.
- `source_resource_id` (string, Computed): source resource ID.
- `target_resource_id` (string, Computed): target resource ID.
- `instance_id` (string, Computed): instance ID.
- `instance_copy_pair_id` (string, Computed): CVM copy pair ID.
- `percent` (int, Computed): replication progress percent.
- `latest_protection_time` (string, Computed): latest protection time.
- `data_direction` (string, Computed): data direction (POSITIVE/REVERSE).
- `create_from` (string, Computed): creation source (LOCAL/PEER).
- `disaster_recovery_type` (string, Computed): disaster recovery type.
- `peer_cloud_name` (string, Computed): peer cloud name.
- `rollbacking` (int, Computed): whether rollback is in progress.
- `rollback_percent` (int, Computed): rollback progress.
- `create_time` (string, Computed): creation time.
- `account_uin` (string, Computed): account Uin.
- `sub_account_uin` (string, Computed): sub-account Uin.
- `drill_group_id` (string, Computed): drill group ID.
- `protection_time_set` (list of string, Computed): protection time points.
- `disk_copy_pair_set` (list, Computed): disk copy pair list for CVM, each with `copy_pair_id`, `copy_pair_name`, `source_resource_id`, `target_resource_id`, `create_time`.
- `deferred_create` (bool, Computed): whether deferred creation mode.
- `target_cvm_created` (bool, Computed): whether target CVM is actually created.
- `cvm_create_params` (string, Computed): CVM creation params JSON string.

#### Scenario: Required SDK fields are present
- **WHEN** a developer inspects the resource schema
- **THEN** every field declared in `CreateInstanceCopyPairRequestParams` (ProtectGroupId, CreateTargetInstanceParameters, InstanceCopyPairName, ClientToken, RecoveryPointObjective) appears in the schema with semantically equivalent typing, and the `CreateTargetInstanceParameters` nested resource mirrors all fields of `CreateInstanceModel`.

#### Scenario: No undocumented schema fields
- **WHEN** a developer inspects the resource schema
- **THEN** there are no fields beyond those listed above and those returned by `DescribeCopyPairs`; no derived flags or synthetic toggles are introduced.

### Requirement: Resource ID
The resource ID SHALL be the bare `CopyPairId` returned by `CreateInstanceCopyPair` (the first element of the `CopyPairIds` response array). The resource SHALL support `terraform import` using just that ID.

#### Scenario: Create sets the ID
- **WHEN** `CreateInstanceCopyPair` returns a non-nil, non-empty `CopyPairIds` slice
- **THEN** the resource takes the first element, calls `d.SetId(<copyPairId>)`, and proceeds to poll the read interface until the copy pair state is no longer `INIT`.

#### Scenario: Import by ID
- **WHEN** an operator runs `terraform import tencentcloud_bdrc_instance_copy_pair.x cvmcopypair-xxxxxxxx`
- **THEN** the resource state is hydrated from `DescribeCopyPairs` using the imported ID, with no manual composite-ID parsing required.

### Requirement: Async create with read polling
On Create, the resource SHALL invoke `CreateInstanceCopyPair` (wrapped in `resource.Retry(tccommon.WriteRetryTimeout, ...)`), then poll `DescribeCopyPairs` until `CopyPairState` transitions out of `INIT` (e.g., to `RUNNING`, `FULL_COPYING`, or `NORMAL`) or the read retry budget is exhausted.

#### Scenario: Successful async create
- **WHEN** `CreateInstanceCopyPair` succeeds and the subsequent `DescribeCopyPairs` poll shows `CopyPairState` is not `INIT`
- **THEN** the resource sets the ID, invokes Read, and returns no error.

#### Scenario: Empty CopyPairIds
- **WHEN** `CreateInstanceCopyPair` returns a nil or empty `CopyPairIds` slice
- **THEN** the resource returns a `NonRetryableError` with a descriptive message rather than proceeding to set an empty ID.

#### Scenario: Nil response defense
- **WHEN** `CreateInstanceCopyPair` returns a nil `Response`
- **THEN** the retry wrapper returns `resource.NonRetryableError` with a descriptive message rather than dereferencing the nil pointer.

### Requirement: Read with retry and pagination
On Read, the resource SHALL call `BdrcService.DescribeBdrcInstanceCopyPairById(ctx, copyPairId)`, which:
- Wraps the SDK call `DescribeCopyPairsWithContext` in `resource.Retry(tccommon.ReadRetryTimeout, ...)`.
- Sets `CopyPairType = "INSTANCE"` and `CopyPairIds = [copyPairId]`.
- Returns `(nil, nil)` when the instance is absent (after the retry budget is exhausted with empty results, logging `[CRUD] bdrc_instance_copy_pair id=<id>` before returning).

When the helper returns `(nil, nil)`, the resource SHALL log a `[WARN]` line, call `d.SetId("")`, and return no error.

#### Scenario: Resource present
- **WHEN** the helper finds a matching `CopyPair`
- **THEN** the resource populates all schema fields (input + computed) from the response, checking each response field for nil before calling `d.Set(...)`.

#### Scenario: Resource removed externally
- **WHEN** the helper returns `(nil, nil)` (no matching copy pair)
- **THEN** the resource logs `[CRUD] bdrc_instance_copy_pair id=<id>`, calls `d.SetId("")`, and returns no error.

#### Scenario: Empty response in retry block returns NonRetryableError
- **WHEN** `DescribeCopyPairs` returns an empty `CopyPairSet` inside the service-layer retry block
- **THEN** the retry block returns `resource.NonRetryableError` (not `d.SetId("")` directly) to allow the outer retry to continue, and only after the budget is exhausted does the service return `(nil, nil)`.

### Requirement: Update path limited to copy pair name
The Update function SHALL call `ModifyCopyPairAttribute` (wrapped in `resource.Retry(tccommon.WriteRetryTimeout, ...)`) **only when** `instance_copy_pair_name` has changed. The request MUST include `CopyPairId` (from `d.Id()`), `CopyPairType = "INSTANCE"`, and `CopyPairName`. All other schema fields are ForceNew; if any ForceNew field is detected as changed, the Update function SHALL return an error.

#### Scenario: Name change
- **WHEN** only `instance_copy_pair_name` changes
- **THEN** the resource calls `ModifyCopyPairAttribute` with the new name and returns to Read.

#### Scenario: ForceNew field changed
- **WHEN** a ForceNew field (e.g., `protect_group_id`, `create_target_instance_parameters`) is detected as changed
- **THEN** the Update function returns an error indicating the field is immutable and requires recreation.

### Requirement: Delete with CopyPairType=INSTANCE
On Delete, the resource SHALL call `DeleteCopyPairs` (wrapped in `resource.Retry(tccommon.WriteRetryTimeout, ...)`) with `CopyPairIds = [d.Id()]`, `CopyPairType = "INSTANCE"`, and `DeleteTargetResource` from the schema (if set). The resource SHALL defend against nil `Response`.

#### Scenario: Successful delete
- **WHEN** `DeleteCopyPairs` succeeds
- **THEN** the resource returns no error and Terraform marks the resource as destroyed.

#### Scenario: Nil response defense
- **WHEN** `DeleteCopyPairs` returns a nil `Response`
- **THEN** the retry wrapper returns `resource.NonRetryableError` with a descriptive message.

### Requirement: Retry coverage
Every SDK call (`CreateInstanceCopyPairWithContext`, `DescribeCopyPairsWithContext`, `ModifyCopyPairAttributeWithContext`, `DeleteCopyPairsWithContext`) SHALL be invoked from inside a `resource.Retry` block. The retry budget is `tccommon.WriteRetryTimeout` for write operations and `tccommon.ReadRetryTimeout` for read operations. On SDK error, the wrapper SHALL use `tccommon.RetryError(e)` to classify retryability. No nested retry inside the retry block.

#### Scenario: Transient SDK error
- **WHEN** any SDK call returns a transient TencentCloud SDK error
- **THEN** the call is retried via `tccommon.RetryError(e)` until it succeeds or the retry budget is exhausted.

### Requirement: Logging conventions
The resource SHALL emit:
- `defer tccommon.LogElapsed("resource.tencentcloud_bdrc_instance_copy_pair.<op>")()` at the top of every CRUD function.
- `defer tccommon.InconsistentCheck(d, meta)()` at the top of every CRUD function.
- A `[DEBUG]` line per SDK invocation containing the request action, request body, and response body.
- A `[CRITAL]%s ... failed, reason:%+v` line on every retry-block failure, using the resource name `bdrc_instance_copy_pair` (not "该资源"/"当前资源").
- A `[CRUD] bdrc_instance_copy_pair id=<id>` line before `d.SetId("")` when clearing the ID due to empty read response.

#### Scenario: Standard log lines emitted
- **WHEN** any CRUD operation runs
- **THEN** the operation's elapsed time is logged via `tccommon.LogElapsed` and inconsistency is checked via `tccommon.InconsistentCheck`, and all log messages use the literal resource name `bdrc_instance_copy_pair`.

### Requirement: Documentation and tests
The change SHALL include:
- A markdown document at `tencentcloud/services/bdrc/resource_tc_bdrc_instance_copy_pair.md` containing a one-line description mentioning "BDRC", an HCL example using `jsonencode()` for any JSON string fields, and an `Import` section explaining `terraform import tencentcloud_bdrc_instance_copy_pair.example cvmcopypair-xxxxxxxx`. The document SHALL NOT include `Argument Reference` or `Attribute Reference` sections (auto-generated by `make doc`).
- A unit test file at `tencentcloud/services/bdrc/resource_tc_bdrc_instance_copy_pair_test.go` using gomonkey mock (not terraform test suite) to mock cloud API calls and test business logic of Create/Read/Update/Delete.

#### Scenario: Documentation present
- **WHEN** the change is merged
- **THEN** the markdown documentation file exists and contains both an HCL example and an `import` example using the bare `CopyPairId`.

#### Scenario: Test file present
- **WHEN** the change is merged
- **THEN** the `_test.go` file declares unit tests using gomonkey mocks for `UseBdrcV20260330Client` methods, covering the CRUD logic branches.

### Requirement: SDK constraint
The implementation SHALL NOT modify any file under `vendor/github.com/tencentcloud/tencentcloud-sdk-go/`. If a required API is missing from the vendored SDK, the implementer MUST halt and request an SDK upgrade rather than authoring or patching SDK source.

#### Scenario: Vendored SDK is sufficient
- **WHEN** the implementation begins
- **THEN** `CreateInstanceCopyPair`, `DescribeCopyPairs`, `ModifyCopyPairAttribute`, and `DeleteCopyPairs` are confirmed present under `vendor/github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/bdrc/v20260330/` before any code is written.
