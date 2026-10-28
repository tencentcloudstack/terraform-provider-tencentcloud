# Spec: EMR Cluster v2 Resource

**Capability**: `emr-cluster-v2-resource`
**Related Change**: `add-emr-cluster-v2-resource`
**Status**: Draft

---

## ADDED Requirements

### Requirement: Resource Registration

The provider SHALL register a new Terraform resource named `tencentcloud_emr_cluster_v2` that manages Tencent Cloud EMR clusters created through the `CreateCluster` (EMR `2019-01-03`) API.

The resource MUST be exposed in `tencentcloud/provider.go` ResourcesMap as `"tencentcloud_emr_cluster_v2": emr.ResourceTencentCloudEmrClusterV2()` and its implementation MUST live at `tencentcloud/services/emr/resource_tc_emr_cluster_v2.go`.

#### Scenario: Resource is discoverable

- **WHEN** a user runs `terraform providers schema -json` against a provider binary built from this branch
- **THEN** the output contains a resource entry for `tencentcloud_emr_cluster_v2`
- **AND** the resource file `resource_tc_emr_cluster_v2.go` defines the exported constructor `ResourceTencentCloudEmrClusterV2`

#### Scenario: Companion artifacts exist

- **WHEN** the change is merged
- **THEN** `tencentcloud/services/emr/resource_tc_emr_cluster_v2.md` exists and contains an `Example Usage` block with an `hcl` code fence
- **AND** `tencentcloud/services/emr/resource_tc_emr_cluster_v2_test.go` exists and defines at least one `TestAccTencentCloudEmrClusterV2_*` acceptance test

---

### Requirement: Schema Must Mirror CreateCluster Request

The resource schema MUST expose one Terraform attribute per top-level field of `CreateClusterRequestParams` from `tencentcloud-sdk-go/tencentcloud/emr/v20190103/models.go`, with no field merging, renaming, or omission.

The top-level attributes are, in snake_case: `product_version`, `enable_support_ha_flag`, `instance_name`, `instance_charge_type`, `login_settings`, `scene_software_config`, `instance_charge_prepaid`, `security_group_ids`, `script_bootstrap_action_config`, `client_token`, `need_master_wan`, `enable_remote_login_flag`, `enable_kerberos_flag`, `custom_conf`, `tags`, `disaster_recover_group_ids`, `enable_cbs_encrypt_flag`, `meta_db_info`, `depend_service`, `zone_resource_configuration`, `cos_bucket`, `node_marks`, `load_balancer_id`, `default_meta_version`, `need_cdb_audit`, `sg_ip`, `partition_number`, `web_ui_version`.

Nested SDK structs MUST be represented as `TypeList` blocks preserving the SDK hierarchy (`LoginSettings`, `SceneSoftwareConfig`, `InstanceChargePrepaid`, `ScriptBootstrapActionConfig`, `Tag`, `CustomMetaDBInfo`, `DependService`, `ZoneResourceConfiguration → VirtualPrivateCloud/Placement/AllNodeResourceSpec → NodeResourceSpec → DiskSpecInfo/Tag`, `NodeMark`).

No schema field in this change SHALL have `ForceNew: true`.

#### Scenario: Every CreateCluster parameter has a schema attribute

- **WHEN** the schema map of `ResourceTencentCloudEmrClusterV2` is enumerated
- **THEN** for every exported field F in `CreateClusterRequestParams`, the schema map contains a key equal to the snake_case form of F
- **AND** the corresponding Go type in the Create request builder uses the SDK pointer type of F (`*string`, `*bool`, `*int64`, `[]*T`)

#### Scenario: No field carries ForceNew

- **WHEN** the schema map is enumerated
- **THEN** every `*schema.Schema` in the map has `ForceNew == false` (default)

#### Scenario: Sensitive fields are marked Sensitive

- **WHEN** the schema is inspected
- **THEN** `login_settings.password` has `Sensitive: true`
- **AND** `meta_db_info.meta_data_pass` has `Sensitive: true`

#### Scenario: Required-vs-Optional matches the API contract

- **WHEN** a user writes a minimal config containing only `product_version`, `enable_support_ha_flag`, `instance_name`, `instance_charge_type`, `login_settings`, and `scene_software_config`
- **THEN** `terraform validate` passes
- **AND** omitting any one of those six top-level attributes causes `terraform validate` to fail with a "Missing required argument" error

---

### Requirement: Create Operation

The Create handler MUST call `CreateCluster` (EMR v20190103), set the Terraform resource ID to the returned `InstanceId`, and then poll `DescribeInstances` until `Clusters[0].Status == 2` before returning.

The Create handler MUST:
- Wrap the SDK call in `resource.Retry(tccommon.WriteRetryTimeout, ...)`.
- Use the Terraform-configured `schema.TimeoutCreate` (default **60 minutes**) as the polling budget.
- Return a non-retryable error if the cluster enters a terminated/failed status (`Status` in the EMR terminated set).
- Emit standard log lines with `tccommon.LogElapsed("resource.tencentcloud_emr_cluster_v2.create")` and `tccommon.InconsistentCheck`.

#### Scenario: Successful cluster creation

- **WHEN** a valid config is applied
- **THEN** the provider calls `CreateCluster` exactly once (modulo retry on transient errors)
- **AND** the resource ID is set to the `InstanceId` returned by the API
- **AND** the provider polls `DescribeInstances` until `Clusters[0].Status == 2`
- **AND** Read is invoked before Create returns

#### Scenario: Async provisioning failure

- **WHEN** `CreateCluster` returns `InstanceId` but the cluster subsequently enters a failed status (e.g. `Status == 6`)
- **THEN** Create returns a non-retryable error describing the terminal status
- **AND** the resource stays in state with the `InstanceId` so the user can run `terraform destroy`

#### Scenario: Create timeout honoured

- **WHEN** a user sets `timeouts { create = "30m" }`
- **THEN** the Create polling budget equals 30 minutes (minus a safety margin)
- **AND** if the cluster has not reached `Status == 2` by then, Create returns a wrapping error from `resource.Retry`

---

### Requirement: Read Operation

The Read handler MUST populate state from `DescribeInstances` (cluster-level fields) and, best-effort, from `DescribeClusterNodes` (node-level fields under `zone_resource_configuration`).

If `DescribeInstances` returns an empty `Clusters` list for the configured `InstanceId`, the Read handler MUST call `d.SetId("")` and return `nil` so Terraform treats the resource as externally deleted.

Fields that are **not** returned by `DescribeInstances`/`DescribeClusterNodes` (secrets, `client_token`, `custom_conf`) MUST NOT be overwritten from API data — they remain whatever the user's config declared.

#### Scenario: Normal read round-trip

- **WHEN** a cluster exists and the Read handler is invoked
- **THEN** `DescribeInstances` is called with `InstanceIds=[d.Id()]`
- **AND** state is populated for all cluster-level fields reported in `ClusterInstancesInfo`
- **AND** `DescribeClusterNodes` is called at least once with `NodeFlag="all"` to hydrate node specs

#### Scenario: Externally deleted cluster

- **WHEN** `DescribeInstances` returns `Clusters: []` (or the cluster's terminated flag is set)
- **THEN** the Read handler calls `d.SetId("")`
- **AND** returns `nil`

#### Scenario: Sensitive fields preserved

- **WHEN** the Read handler runs against an existing cluster
- **THEN** `login_settings.password` in state is NOT replaced by API data
- **AND** `meta_db_info.meta_data_pass` in state is NOT replaced by API data

---

### Requirement: Update Operation (No-op in This Change)

The Update handler MUST be registered (required by Terraform Plugin SDK because no field has `ForceNew`) but MUST NOT call any modify API. It SHALL simply re-invoke the Read handler.

The handler MUST include a code comment identifying the follow-up change that will add real modify support.

#### Scenario: terraform apply after config change

- **WHEN** a user changes a non-ForceNew field and runs `terraform apply`
- **THEN** the Update handler is invoked
- **AND** no remote modify API is called
- **AND** the handler returns the result of Read

---

### Requirement: Delete Operation

The Delete handler MUST call `TerminateInstance` with the cluster's `InstanceId` and then poll `DescribeInstances` until the cluster is no longer listed (or has reached a terminated status) before returning.

The Delete handler MUST use `schema.TimeoutDelete` (default **30 minutes**).

#### Scenario: Successful deletion

- **WHEN** `terraform destroy` is issued
- **THEN** `TerminateInstance` is called with the resource's `InstanceId`
- **AND** the handler polls `DescribeInstances` until `Clusters` is empty or the cluster status is terminated
- **AND** returns `nil`

#### Scenario: Delete timeout honoured

- **WHEN** a user sets `timeouts { delete = "10m" }`
- **THEN** the polling budget equals 10 minutes
- **AND** if the cluster has not disappeared by then, Delete returns a wrapping error

---

### Requirement: Resource ID Is the EMR InstanceId

The Terraform resource ID MUST be the EMR `InstanceId` returned by `CreateCluster` (e.g. `emr-f2da1cd`), with no prefix, suffix, or separator.

Import MUST be supported via `schema.ImportStatePassthrough` so that `terraform import tencentcloud_emr_cluster_v2.example emr-xxxxxxxx` works.

#### Scenario: ID format

- **WHEN** Create succeeds
- **THEN** `d.Id()` returns a string matching the regex `^emr-[a-z0-9]+$`
- **AND** this string equals `*response.Response.InstanceId`

#### Scenario: Import works

- **WHEN** a user runs `terraform import tencentcloud_emr_cluster_v2.example emr-abc123`
- **THEN** the provider stores the ID directly
- **AND** the next `terraform plan` reads back the cluster via `DescribeInstances`

---

### Requirement: Async Operations Must Declare Timeouts

The resource schema MUST declare `Timeouts` for `Create`, `Read`, and `Delete` operations because all three are asynchronous or paginated.

Default values MUST be:
- `Create`: 60 minutes
- `Read`: 20 minutes
- `Delete`: 30 minutes

The CRUD handlers MUST consume these via `d.Timeout(schema.TimeoutCreate/Read/Delete)` and pass the derived context to `resource.Retry`.

#### Scenario: Timeout block is honoured

- **WHEN** the schema is inspected
- **THEN** the resource's `Timeouts` field is non-nil with the three defaults above set

---

### Requirement: Service Layer Helpers

`tencentcloud/services/emr/service_tencentcloud_emr.go` MUST be extended with helper methods used by this resource only. Existing methods MUST NOT be modified.

Required new helpers:
- `DescribeEmrClusterV2ById(ctx, instanceId) (*emr.ClusterInstancesInfo, error)` — wraps `DescribeInstances` and returns the first matching cluster, or `nil` if not found.
- `DescribeEmrClusterV2Nodes(ctx, instanceId, nodeFlag) ([]*emr.NodeHardwareInfo, error)` — paginates `DescribeClusterNodes` with `Limit=100` until exhausted.

Both helpers MUST wrap SDK calls in `resource.Retry(tccommon.ReadRetryTimeout, ...)` and emit standard debug log lines.

#### Scenario: Helpers exist and are package-local

- **WHEN** the emr package is compiled
- **THEN** `(*EMRService).DescribeEmrClusterV2ById` and `(*EMRService).DescribeEmrClusterV2Nodes` are defined
- **AND** no other resource file is changed to consume them

---

### Requirement: Code Style Matches tencentcloud_igtm_strategy

The Create/Read/Update/Delete handlers MUST follow the exact style patterns of `resource_tc_igtm_strategy.go`:

- `defer tccommon.LogElapsed("resource.tencentcloud_emr_cluster_v2.<op>")()` at the top.
- `defer tccommon.InconsistentCheck(d, meta)()` at the top.
- `logId := tccommon.GetLogId(tccommon.ContextNil)` then `ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)`.
- Request builder pattern: one `if v, ok := d.GetOk("<field>"); ok {}` block per top-level field, with nested loops for list/nested fields.
- SDK call wrapped in `resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError { ... })`.
- Nil-safe helper usage: `helper.String`, `helper.IntInt64`, `helper.BoolToPtr`, `helper.Strings`, `helper.IntUint64` as appropriate.
- Importer: `Importer: &schema.ResourceImporter{State: schema.ImportStatePassthrough}`.

#### Scenario: Lint and format pass

- **WHEN** `make fmt` and `make lint` run against the new files
- **THEN** both commands exit 0

#### Scenario: Pattern conformance

- **WHEN** a reviewer diffs `resource_tc_emr_cluster_v2.go` against `resource_tc_igtm_strategy.go`
- **THEN** both files use the same top-level defer pattern
- **AND** both files use the same `resource.Retry` + `tccommon.RetryError` wrapping

---

### Requirement: Documentation

`tencentcloud/services/emr/resource_tc_emr_cluster_v2.md` MUST follow the layout of `resource_tc_config_compliance_pack.md`:

1. A one-line "Provides a resource to …" lead-in.
2. An `Example Usage` header followed by a triple-backtick `hcl` block with a realistic end-to-end example covering at minimum: `product_version`, `instance_charge_type=POSTPAID_BY_HOUR`, `login_settings`, `scene_software_config` (with Software list), and one `zone_resource_configuration` entry (VPC + Placement + MasterResourceSpec + CoreResourceSpec).
3. An `Import` section with a shell-fenced `terraform import tencentcloud_emr_cluster_v2.example emr-xxxxxxxx` command.
4. Optional `~> **NOTE:**` block explaining that Update is not yet supported.

The website markdown at `website/docs/r/emr_cluster_v2.html.markdown` MUST be generated by `make doc` and MUST NOT be hand-edited.

#### Scenario: Example MD is present and well-formed

- **WHEN** `resource_tc_emr_cluster_v2.md` is parsed
- **THEN** it contains exactly one `Example Usage` H2-style section
- **AND** it contains at least one triple-backtick `hcl` fenced code block
- **AND** the last section is `Import` with a fenced shell block starting with `terraform import tencentcloud_emr_cluster_v2.`

---

### Requirement: Acceptance Test Coverage

`resource_tc_emr_cluster_v2_test.go` MUST define an acceptance test function `TestAccTencentCloudEmrClusterV2_basic` (following the project convention `TestAccTencentCloud<Service>_<scenario>`).

The test MUST:
- Use `tcacctest.AccPreCheck(t)`.
- Reference `tcacctest.AccProviders`.
- Include a `resource.TestStep` that applies a minimal valid config and asserts `resource.TestCheckResourceAttrSet("tencentcloud_emr_cluster_v2.example", "id")`.
- Include an Import step with `ImportStateVerify: true` and `ImportStateVerifyIgnore: ["login_settings", "client_token", "meta_db_info"]` (fields known not to round-trip).

#### Scenario: Test compiles

- **WHEN** `go test -run TestAccTencentCloudEmrClusterV2_basic -count=1 ./tencentcloud/services/emr/...` is invoked without `TF_ACC=1`
- **THEN** the test is skipped (per `tcacctest` convention) but the compilation succeeds with exit code 0

#### Scenario: Import step skips non-readable fields

- **WHEN** the test's Import step runs
- **THEN** `ImportStateVerifyIgnore` includes at least `login_settings` and `client_token`

---

### Requirement: Backward Compatibility

This change MUST NOT modify or remove any existing resource, data source, schema field, or service-layer method. In particular, `tencentcloud_emr_cluster` MUST remain byte-identical in its schema contract.

#### Scenario: Existing resource untouched

- **WHEN** the diff of this change is inspected
- **THEN** `tencentcloud/services/emr/resource_tc_emr_cluster.go` is unchanged
- **AND** no existing entry in `provider.go` ResourcesMap is removed or renamed
