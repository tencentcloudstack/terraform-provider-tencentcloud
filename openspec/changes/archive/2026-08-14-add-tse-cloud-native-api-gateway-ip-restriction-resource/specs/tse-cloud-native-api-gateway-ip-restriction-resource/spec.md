## ADDED Requirements

### Requirement: Resource MUST be registered as `tencentcloud_tse_cloud_native_api_gateway_ip_restriction`

The provider SHALL register a new general-type resource named `tencentcloud_tse_cloud_native_api_gateway_ip_restriction` whose Create/Update callbacks invoke the `CreateOrModifyCloudNativeAPIGatewayIPRestriction` API, Read invokes `DescribeCloudNativeAPIGatewayIPRestriction`, and Delete invokes `DeleteCloudNativeAPIGatewayIPRestriction` of `tencentcloud-sdk-go/tencentcloud/tse/v20201207`.

#### Scenario: Resource registered in provider map

- **WHEN** the provider is loaded
- **THEN** `provider.go` exposes the resource via key `"tencentcloud_tse_cloud_native_api_gateway_ip_restriction"` mapped to `tse.ResourceTencentCloudTseCloudNativeAPIGatewayIPRestriction()`, alongside existing `tencentcloud_tse_*` resource entries.

#### Scenario: Resource appears in gendoc index

- **WHEN** `tencentcloud/provider.md` is scanned by `make doc`
- **THEN** the TSE Resource section MUST include `tencentcloud_tse_cloud_native_api_gateway_ip_restriction` so that `website/docs/r/tse_cloud_native_api_gateway_ip_restriction.html.markdown` is generated.

### Requirement: Schema MUST mirror the IP restriction API inputs

The resource schema SHALL declare exactly these top-level argument keys, with semantics matching the SDK request fields of `CreateOrModifyCloudNativeAPIGatewayIPRestrictionRequest`:

| HCL key | SDK field | Type | Required | ForceNew | Computed |
|---|---|---|---|---|---|
| `gateway_id` | `GatewayId` | TypeString | Yes | **Yes** | No |
| `source_type` | `SourceType` | TypeString | Yes | **Yes** | No |
| `source_id` | `SourceId` | TypeString | Yes | **Yes** | No |
| `enabled` | `Enabled` | TypeBool | No | No | Yes |
| `restriction_type` | `RestrictionType` | TypeString | No | No | Yes |
| `address_list` | `AddressList` | TypeSet of TypeString | No | No | Yes |

The schema MUST NOT declare any additional fields. `enabled`, `restriction_type`, and `address_list` are `Optional: true, Computed: true` (configurable by the user and refreshed from the Describe API), while `gateway_id`, `source_type`, and `source_id` are `Required`. `address_list` is a `TypeSet` so that element ordering is not significant.

#### Scenario: Required identity fields enforce on plan

- **WHEN** the user writes a config that omits any of `gateway_id`, `source_type`, or `source_id`
- **THEN** `terraform plan` SHALL fail validation pointing at the missing required attribute.

#### Scenario: Configurable fields are optional and computed

- **GIVEN** the user omits `enabled`, `restriction_type`, or `address_list` from the config
- **WHEN** `terraform plan` and `terraform apply` run
- **THEN** the omitted fields are populated from the Describe API response and are not treated as missing required attributes.

#### Scenario: address_list is order-insensitive

- **GIVEN** state holds `address_list` with a set of IP/CIDR entries
- **WHEN** the user reorders those entries in HCL, or the Describe API returns them in a different order
- **THEN** `terraform plan` reports no diff for `address_list` because it is a `TypeSet`.

#### Scenario: Changing identity triplet forces replacement

- **GIVEN** state has `gateway_id = "gw-1"`, `source_type = "route"`, `source_id = "r-1"`
- **WHEN** the user changes any one of the triplet in HCL
- **THEN** Terraform's plan reports a destroy + create cycle because each of `gateway_id`, `source_type`, `source_id` is ForceNew.

### Requirement: Resource ID MUST be the composite of gateway_id, source_type and source_id

After Create, the Terraform resource ID SHALL be `strings.Join([]string{gateway_id, source_type, source_id}, tccommon.FILED_SP)` (i.e. `gateway_id#source_type#source_id`). Read/Update/Delete SHALL split `d.Id()` by `tccommon.FILED_SP` and fail with `"id is broken"` when the segment count is not 3.

#### Scenario: ID is composed from the identity triplet

- **GIVEN** HCL declares `gateway_id = "gw-1"`, `source_type = "route"`, `source_id = "r-1"`
- **WHEN** Create completes
- **THEN** `d.Id()` returns `"gw-1#route#r-1"`.

### Requirement: Create MUST delegate to Update

Because `CreateOrModifyCloudNativeAPIGatewayIPRestriction` is an upsert (overwrite-semantic) API, the Create callback SHALL:
1. Read `gateway_id`, `source_type`, `source_id` and set `d.SetId(strings.Join([]string{gateway_id, source_type, source_id}, tccommon.FILED_SP))`.
2. Return `resourceTencentCloudTseCloudNativeAPIGatewayIPRestrictionUpdate(d, meta)`.

The Update callback contains the actual SDK call.

#### Scenario: Create produces a single upsert API call

- **GIVEN** HCL with the full field set
- **WHEN** the user runs `terraform apply`
- **THEN** exactly one `CreateOrModifyCloudNativeAPIGatewayIPRestriction` request is issued with all six fields populated; `d.Id()` returns the composite triplet.

### Requirement: Update MUST call CreateOrModifyCloudNativeAPIGatewayIPRestriction with all current values

The Update callback SHALL build a `CreateOrModifyCloudNativeAPIGatewayIPRestrictionRequest` populated from the current `d.Get(...)` values for `gateway_id`, `source_type`, `source_id`, `enabled`, `restriction_type`, and `address_list`, and call `UseTseClient().CreateOrModifyCloudNativeAPIGatewayIPRestrictionWithContext(ctx, request)` inside `resource.Retry(tccommon.WriteRetryTimeout, ...)`. Any change to `enabled`, `restriction_type`, or `address_list` therefore triggers a full overwrite. Update MUST end by calling Read.

#### Scenario: Editing address_list triggers exactly one upsert call

- **GIVEN** state has `address_list = ["1.1.1.1"]`
- **WHEN** the user changes it to `["2.2.2.2"]`
- **THEN** Update issues exactly one `CreateOrModifyCloudNativeAPIGatewayIPRestriction` request carrying the new `AddressList`.

#### Scenario: Transient SDK error is retried

- **GIVEN** the first invocation of `CreateOrModifyCloudNativeAPIGatewayIPRestriction` returns a retriable error
- **WHEN** the retry callback re-runs
- **THEN** the second attempt's response is observed; the original error surfaces only after the retry budget is exhausted.

### Requirement: Read MUST refresh state from DescribeCloudNativeAPIGatewayIPRestriction

The Read callback SHALL split `d.Id()` into `gateway_id`, `source_type`, `source_id`, call `UseTseClient().DescribeCloudNativeAPIGatewayIPRestrictionWithContext(ctx, request)` inside `resource.Retry(tccommon.ReadRetryTimeout, ...)`, and:
- If `result == nil` or `result.Response == nil` or `result.Response.Result == nil`, the implementation SHALL log `[CRUD] tse_cloud_native_api_gateway_ip_restriction id=<d.Id()>` then `d.SetId("")` and return nil.
- Otherwise, set `gateway_id`, `source_type`, `source_id` from the split ID, and set `enabled`, `restriction_type`, `address_list` from `result.Response.Result` when the respective fields are non-nil.

#### Scenario: Read populates state from API response

- **GIVEN** the API returns `Result{Enabled: true, RestrictionType: "whiteList", AddressList: ["10.0.0.0/8"]}`
- **WHEN** Read runs
- **THEN** state holds `enabled = true`, `restriction_type = "whiteList"`, `address_list = ["10.0.0.0/8"]`.

#### Scenario: Read clears id when resource is gone

- **GIVEN** the API returns a nil `Result`
- **WHEN** Read runs
- **THEN** the log line `[CRUD] tse_cloud_native_api_gateway_ip_restriction id=<id>` is emitted before `d.SetId("")`.

### Requirement: Delete MUST call DeleteCloudNativeAPIGatewayIPRestriction

The Delete callback SHALL split `d.Id()` into `gateway_id`, `source_type`, `source_id`, build a `DeleteCloudNativeAPIGatewayIPRestrictionRequest` with those three fields, and call `UseTseClient().DeleteCloudNativeAPIGatewayIPRestrictionWithContext(ctx, request)` inside `resource.Retry(tccommon.WriteRetryTimeout, ...)`.

#### Scenario: Delete removes the IP restriction

- **GIVEN** state has a bound IP restriction
- **WHEN** the user runs `terraform destroy`
- **THEN** exactly one `DeleteCloudNativeAPIGatewayIPRestriction` request is issued with the identity triplet; the resource is removed from state.

### Requirement: Every API call MUST be wrapped in resource.Retry

Every SDK invocation (`CreateOrModify...`, `Describe...`, `Delete...`) SHALL be wrapped in `resource.Retry(...)` with errors forwarded via `tccommon.RetryError(e)`. Write operations use `tccommon.WriteRetryTimeout`; Read uses `tccommon.ReadRetryTimeout`. Setting `d.SetId` and field `d.Set` MUST happen outside the retry block.

#### Scenario: Retry wrapper present

- **WHEN** a code reviewer inspects each CRUD callback
- **THEN** the SDK call is inside a `resource.Retry(...)` block, not a bare invocation.

### Requirement: Response reads MUST be nil-safe

The retry callback SHALL guard `result == nil || result.Response == nil` and return `resource.NonRetryableError` when either is nil, before further use of `result`.

#### Scenario: Nil response is detected

- **GIVEN** the SDK returns `(result == nil, err == nil)`
- **WHEN** the retry callback runs
- **THEN** the callback returns `resource.NonRetryableError(...)` and does NOT panic.

### Requirement: Resource MUST support import

The resource SHALL declare an `Importer` with `schema.ImportStatePassthrough`. Import reads the composite ID `gateway_id#source_type#source_id` and delegates to Read.

#### Scenario: Import populates state from composite id

- **WHEN** the user runs `terraform import tencentcloud_tse_cloud_native_api_gateway_ip_restriction.example gw-1#route#r-1`
- **THEN** Read is invoked with `d.Id() = "gw-1#route#r-1"` and state is populated from the Describe API.

### Requirement: Documentation and unit test MUST follow project conventions

- The HCL example markdown SHALL live at `tencentcloud/services/tse/resource_tc_tse_cloud_native_api_gateway_ip_restriction.md`, containing a one-line summary mentioning TSE, an `Example Usage` block with all fields, and an `Import` section stating the composite ID format. It MUST NOT include `Argument Reference` or `Attribute Reference` sections (auto-generated).
- A unit test SHALL live at `tencentcloud/services/tse/resource_tc_tse_cloud_native_api_gateway_ip_restriction_test.go`, package `tse_test`, using gomonkey to mock the TSE client and covering Create, Read, Update, Delete callbacks.
- Running `make doc` SHALL regenerate `website/docs/r/tse_cloud_native_api_gateway_ip_restriction.html.markdown`.

#### Scenario: Generated website doc lists the resource

- **WHEN** `make doc` runs
- **THEN** `website/docs/r/tse_cloud_native_api_gateway_ip_restriction.html.markdown` exists and lists every schema attribute defined in the Schema requirement.

#### Scenario: Unit test structure

- **WHEN** the test file is opened
- **THEN** the package is `tse_test`, tests mock `UseTseClient` and the three `...WithContext` methods, and exercise Create/Read/Update/Delete.
