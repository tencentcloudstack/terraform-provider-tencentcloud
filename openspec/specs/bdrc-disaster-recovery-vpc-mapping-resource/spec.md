# bdrc-disaster-recovery-vpc-mapping-resource Specification

## Purpose
TBD - created by archiving change add-bdrc-disaster-recovery-vpc-mapping-resource. Update Purpose after archive.
## Requirements
### Requirement: Resource MUST be registered as `tencentcloud_bdrc_disaster_recovery_vpc_mapping`

The provider SHALL register a new general-type resource named `tencentcloud_bdrc_disaster_recovery_vpc_mapping` whose Create / Read / Delete callbacks invoke the BDRC cloud APIs of `tencentcloud-sdk-go/tencentcloud/bdrc/v20260330`. The Update callback SHALL NOT invoke any cloud API (the API set is CRD-only); Update MUST fall back to Read after an immutable-args guard.

The connectivity layer SHALL expose a `UseBdrcV20260330Client()` method on `*TencentCloudClient` returning a `*bdrcv20260330.Client` (mirroring `UseIgtmV20231024Client`), since bdrc has no existing client in `tencentcloud/connectivity/client.go`.

#### Scenario: Resource registered in provider map

- **WHEN** the provider is loaded
- **THEN** `provider.go` exposes the resource via key `"tencentcloud_bdrc_disaster_recovery_vpc_mapping"` mapped to `bdrc.ResourceTencentCloudBdrcDisasterRecoveryVpcMapping()`, and imports the `bdrc` service package.

#### Scenario: BDRC client accessor exists

- **WHEN** any CRUD callback calls `meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseBdrcV20260330Client()`
- **THEN** a non-nil `*bdrcv20260330.Client` is returned, constructed from the provider credential and region.

#### Scenario: Resource appears in gendoc index

- **WHEN** `tencentcloud/provider.md` is scanned by `make doc`
- **THEN** a BDRC section includes `tencentcloud_bdrc_disaster_recovery_vpc_mapping` so that `website/docs/r/bdrc_disaster_recovery_vpc_mapping.html.markdown` is generated.

### Requirement: Schema MUST mirror the BDRC VPC mapping APIs (inputs flattened, outputs flattened to top level)

The resource schema SHALL declare exactly these top-level argument keys:

| HCL key | API source | Type | Required | ForceNew | Computed |
|---|---|---|---|---|---|
| `site_pair_id` | Create/Describe input `SitePairId`; Describe output `VpcMappingSet[].SitePairId` | TypeString | Yes | **Yes** | - |
| `source_vpc_id` | Create input `SourceVpcId` | TypeString | Yes | **Yes** | - |
| `source_subnet_id` | Create input `SourceSubnetId` | TypeString | Yes | **Yes** | - |
| `target_vpc_id` | Create input `TargetVpcId` | TypeString | Yes | **Yes** | - |
| `target_subnet_id` | Create input `TargetSubnetId` | TypeString | Yes | **Yes** | - |
| `id` | Describe output `VpcMappingSet[].Id` (uint64) | TypeInt | - | - | Yes |
| `source_vpc` | Describe output `VpcMappingSet[].SourceVpc` | TypeString | - | - | Yes |
| `source_subnet` | Describe output `VpcMappingSet[].SourceSubnet` | TypeString | - | - | Yes |
| `target_vpc` | Describe output `VpcMappingSet[].TargetVpc` | TypeString | - | - | Yes |
| `target_subnet` | Describe output `VpcMappingSet[].TargetSubnet` | TypeString | - | - | Yes |
| `status` | Describe output `VpcMappingSet[].Status` | TypeString | - | - | Yes |
| `life_state` | Describe output `VpcMappingSet[].LifeState` | TypeString | - | - | Yes |

The schema MUST NOT introduce a `vpc_mapping_set` nested layer; the Describe response list SHALL be flattened so each element's fields are top-level schema keys.

#### Scenario: Required input fields enforce on plan

- **WHEN** the user writes a config that omits any of `site_pair_id`, `source_vpc_id`, `source_subnet_id`, `target_vpc_id`, `target_subnet_id`
- **THEN** `terraform plan` SHALL fail validation pointing at the missing required attribute.

#### Scenario: Changing any business argument forces replacement

- **GIVEN** state has `source_vpc_id = "vpc-old"`
- **WHEN** the user changes `source_vpc_id` to `"vpc-new"` in HCL
- **THEN** Terraform's plan reports a destroy + create cycle (the argument is ForceNew).

#### Scenario: No `vpc_mapping_set` nesting in schema

- **WHEN** the schema is inspected
- **THEN** there is no top-level key named `vpc_mapping_set`; each Describe output field (`id`, `source_vpc`, `source_subnet`, `target_vpc`, `target_subnet`, `status`, `life_state`) is a standalone top-level schema key.

### Requirement: Resource ID MUST be the compound `sitePairId#vpcMappingId`

After Create, the Terraform resource ID SHALL be `strings.Join([]string{sitePairId, vpcMappingIdStr}, tccommon.FILED_SP)` where `vpcMappingIdStr` is the decimal string form of the `Id` (uint64) obtained from the post-create Describe query. Read and Delete SHALL split `d.Id()` by `tccommon.FILED_SP` to recover `sitePairId` and `vpcMappingId`.

#### Scenario: ID is a two-part compound string

- **GIVEN** HCL declares `site_pair_id = "sp-001"` and the post-create query returns `Id = 88`
- **WHEN** Create completes
- **THEN** `d.Id()` returns `"sp-001#88"` (using the project FILED_SP separator).

#### Scenario: Broken compound ID is rejected

- **WHEN** Read or Delete receives `d.Id()` that does not contain exactly one `tccommon.FILED_SP` separator
- **THEN** the callback returns `fmt.Errorf("id is broken,%s", d.Id())` and does NOT call any SDK API.

### Requirement: Create MUST issue CreateDisasterRecoveryVpcMapping then re-query DescribeVpcMappings to obtain the mapping ID

The Create callback SHALL:

1. Build a `CreateDisasterRecoveryVpcMappingRequest` populated from `site_pair_id`, `source_vpc_id`, `source_subnet_id`, `target_vpc_id`, `target_subnet_id` (all via `helper.String(...)`).
2. Call `UseBdrcV20260330Client().CreateDisasterRecoveryVpcMappingWithContext(ctx, request)` inside `resource.Retry(tccommon.WriteRetryTimeout, ...)`. Errors MUST be wrapped via `tccommon.RetryError(e)`. On success, guard `result == nil || result.Response == nil` and return `resource.NonRetryableError(...)` if either is nil.
3. After the create retry succeeds (outside the retry block), call `DescribeVpcMappingsWithContext` inside `resource.Retry(tccommon.ReadRetryTimeout, ...)` with `SitePairId` set and `Filters` containing `source-vpc-id` and `source-subnet-id`; set `Limit` to 100 (the API-documented maximum). Inside the Describe retry: if `result == nil || result.Response == nil` return `tccommon.RetryError(e)`-style retryable behavior on SDK error; if `len(result.Response.VpcMappingSet) == 0` return a retryable error to wait for eventual consistency; if the list is non-empty but no entry matches all four fields (SourceVpc/SourceSubnet/TargetVpc/TargetSubnet) or the matched entry's `Id` is nil, return `resource.NonRetryableError(...)`.
4. After the Describe retry succeeds (outside the retry block), set `d.SetId(...)` to the compound `sitePairId#vpcMappingIdStr`.
5. End by calling Read.

#### Scenario: Create issues exactly one Create API call and one Describe query

- **GIVEN** HCL with all five input fields
- **WHEN** the user runs `terraform apply`
- **THEN** exactly one `CreateDisasterRecoveryVpcMapping` request is issued with all five fields populated; followed by one or more `DescribeVpcMappings` queries until the new mapping is found; `d.Id()` returns the compound `sitePairId#vpcMappingId`.

#### Scenario: Create response missing Response is rejected

- **GIVEN** the SDK returns `(result != nil, result.Response == nil)`
- **WHEN** the create retry callback runs
- **THEN** the callback returns `resource.NonRetryableError(...)` with a clear message and does NOT proceed to the Describe query.

#### Scenario: Post-create Describe empty list is retried

- **GIVEN** Create succeeded but the first Describe returns an empty `VpcMappingSet` (eventual consistency)
- **WHEN** the Describe retry callback re-runs
- **THEN** a retryable error is returned so the retry budget continues; the resource ID is NOT set to empty.

### Requirement: Read MUST query DescribeVpcMappings by SitePairId and match by mapping Id

The Read callback SHALL split the compound ID, set `SitePairId` on the Describe request, set `Limit` to 100, call `DescribeVpcMappingsWithContext` inside `resource.Retry(tccommon.ReadRetryTimeout, ...)`, and match the entry whose `Id` equals the `vpcMappingId` portion. Read MUST guard `result == nil || result.Response == nil` inside the retry.

If no matching entry is found, Read SHALL first `log.Printf("[CRUD] bdrc_disaster_recovery_vpc_mapping id=%s", d.Id())` (preserving the id in logs BEFORE clearing), then `d.SetId("")`. Read MUST NOT call `d.SetId("")` before logging the id.

When a matching entry is found, Read SHALL set each output field (`id`, `site_pair_id`, `source_vpc`, `source_subnet`, `target_vpc`, `target_subnet`, `status`, `life_state`) ONLY when the corresponding `VpcMapping` field is non-nil.

#### Scenario: Read populates computed fields from the matched mapping

- **GIVEN** state id `"sp-001#88"`
- **WHEN** `terraform refresh` runs and DescribeVpcMappings returns an entry with `Id=88, Status="ok"`
- **THEN** `d.Set("status", "ok")` is called; `d.Id()` remains `"sp-001#88"`.

#### Scenario: Read clears id only after logging

- **GIVEN** state id `"sp-001#88"`
- **WHEN** DescribeVpcMappings returns no entry with `Id=88`
- **THEN** `log.Printf("[CRUD] bdrc_disaster_recovery_vpc_mapping id=%s", "sp-001#88")` is emitted BEFORE `d.SetId("")`.

#### Scenario: Read skips set for nil response fields

- **GIVEN** the matched entry has `Status == nil`
- **WHEN** Read sets fields
- **THEN** `d.Set("status", ...)` is NOT called for that field.

### Requirement: Update MUST NOT call any cloud API and MUST reject changes to immutable args

The Update callback SHALL declare `immutableArgs := []string{"site_pair_id", "source_vpc_id", "source_subnet_id", "target_vpc_id", "target_subnet_id"}`. If any of these has changed (`d.HasChange`), Update MUST return an error explaining the API is CRD-only and the resource must be recreated. If no immutable arg changed, Update MUST call Read and return its result. Update MUST NOT invoke any BDRC SDK method.

#### Scenario: Changing a business arg in Update path is rejected

- **GIVEN** a change reaches the Update callback with `d.HasChange("target_vpc_id") == true`
- **WHEN** Update runs
- **THEN** Update returns an error mentioning the BDRC API is CRD-only and the resource must be recreated; no SDK call is made.

### Requirement: Delete MUST call DeleteDisasterRecoveryVpcMapping with a single-element VpcMappingIds list

The Delete callback SHALL split the compound ID, convert the `vpcMappingId` portion to `*uint64` via `helper.StrToUint64Point`, build `request.VpcMappingIds = []*uint64{vpcMappingIdPoint}` (exactly one element), and call `DeleteDisasterRecoveryVpcMappingWithContext` inside `resource.Retry(tccommon.WriteRetryTimeout, ...)`. Errors MUST be wrapped via `tccommon.RetryError(e)`. On success, guard `result == nil || result.Response == nil` and return `resource.NonRetryableError(...)`.

#### Scenario: Delete issues exactly one Delete call with a single id

- **GIVEN** state id `"sp-001#88"`
- **WHEN** `terraform destroy` runs
- **THEN** exactly one `DeleteDisasterRecoveryVpcMapping` request is issued with `VpcMappingIds = [88]`; no other mapping ids are included.

#### Scenario: Delete response missing Response is rejected

- **GIVEN** the SDK returns `(result != nil, result.Response == nil)`
- **WHEN** the delete retry callback runs
- **THEN** the callback returns `resource.NonRetryableError(...)` with a clear message.

### Requirement: Every cloud API call MUST be wrapped in resource.Retry

Each of the three SDK invocations (`CreateDisasterRecoveryVpcMappingWithContext`, `DescribeVpcMappingsWithContext` in both Create-post-query and Read, `DeleteDisasterRecoveryVpcMappingWithContext`) SHALL be wrapped in a `resource.Retry(...)` block. Create and Delete use `tccommon.WriteRetryTimeout`; Describe (Read and post-create query) use `tccommon.ReadRetryTimeout`. Inside the retry, errors MUST be forwarded via `tccommon.RetryError(e)`. Setting `d.SetId(...)`, logging success, and other success-path operations MUST occur outside the retry block (after retry error handling), NOT inside it.

#### Scenario: Retry wrapper present on every SDK call

- **WHEN** a code reviewer inspects the Create, Read, and Delete callbacks
- **THEN** every `*WithContext` SDK call is inside a `resource.Retry(...)` block; success-path `d.SetId(...)` and `log.Printf` calls are outside the retry block.

#### Scenario: Transient SDK error is retried

- **GIVEN** the first invocation of `DescribeVpcMappings` returns a retriable error
- **WHEN** the retry callback re-runs
- **THEN** the second attempt's response is observed; the original error surfaces only after the retry budget is exhausted.

### Requirement: Response field reads MUST be nil-safe

For every SDK call, the implementation SHALL guard `result == nil || result.Response == nil` and return a `resource.NonRetryableError(...)` (for write/create/delete) or propagate via `tccommon.RetryError` (for Describe on transient errors) when either is nil, before any further use of `result`. Additionally, when reading `VpcMappingSet` entries, each field MUST be checked for nil before `d.Set(...)`.

#### Scenario: Nil response is detected without panic

- **GIVEN** the SDK returns `(result == nil, err == nil)` (defensive case)
- **WHEN** the retry callback runs
- **THEN** the callback returns `resource.NonRetryableError(...)` with a clear message and does NOT panic.

### Requirement: Documentation and unit test MUST follow project conventions

- The resource markdown SHALL live at `tencentcloud/services/bdrc/resource_tc_bdrc_disaster_recovery_vpc_mapping.md` and contain: a one-line summary mentioning the BDRC product; an `Example Usage` HCL block with all five input fields; an `Import` section explaining the compound id format `sitePairId#vpcMappingId`. It MUST NOT contain hand-written `Argument Reference` or `Attribute Reference` sections (those are auto-generated).
- A unit test SHALL live at `tencentcloud/services/bdrc/resource_tc_bdrc_disaster_recovery_vpc_mapping_test.go` and use gomonkey to mock the BDRC SDK methods (`CreateDisasterRecoveryVpcMappingWithContext`, `DescribeVpcMappingsWithContext`, `DeleteDisasterRecoveryVpcMappingWithContext`), testing only business logic; it MUST NOT use the terraform acceptance test suite.
- Running `make doc` (in the finalize phase) SHALL regenerate `website/docs/r/bdrc_disaster_recovery_vpc_mapping.html.markdown`.

#### Scenario: Resource markdown structure

- **WHEN** the markdown file is opened
- **THEN** it contains a one-line summary mentioning BDRC, an `Example Usage` HCL block, and an `Import` section with the compound id note; it does NOT contain `Argument Reference` or `Attribute Reference`.

#### Scenario: Unit test uses gomonkey mocks

- **WHEN** the test file is opened
- **THEN** it uses `gomonkey` to mock the three BDRC SDK methods and does NOT reference `tfprotov5`/`terraform-plugin-sdk` test helpers like `testAccProvider`.

