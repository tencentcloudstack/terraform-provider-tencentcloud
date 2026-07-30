# cls-console-resource Specification

## ADDED Requirements

### Requirement: Resource MUST be registered as `tencentcloud_cls_console`

The provider SHALL register a new CRUD-type resource named `tencentcloud_cls_console` whose Create/Read/Update/Delete callbacks invoke the CLS DataSight console APIs (`CreateConsole`, `DescribeConsoles`, `ModifyConsole`, `DeleteConsole`) of `tencentcloud-sdk-go/tencentcloud/cls/v20201016`.

#### Scenario: Resource registered in provider map

- **WHEN** the provider is loaded
- **THEN** `provider.go` exposes the resource via key `"tencentcloud_cls_console"` mapped to `cls.ResourceTencentCloudClsConsole()`, alongside existing `tencentcloud_cls_*` resource entries.

#### Scenario: Importer is configured

- **WHEN** a user runs `terraform import tencentcloud_cls_console.example <ConsoleId>`
- **THEN** the importer SHALL pass the supplied `ConsoleId` straight into the Read callback (using `schema.ImportStatePassthrough`), and Read SHALL hydrate state from `DescribeConsoles` filtered by that ConsoleId.

### Requirement: Schema MUST mirror the CreateConsole API input 1:1

The resource schema SHALL declare exactly these top-level argument keys, with semantics matching the SDK request fields of `CreateConsoleRequestParams`:

| HCL key | SDK field | Type | Required | ForceNew | Computed |
|---|---|---|---|---|---|
| `access_mode` | `AccessMode` | TypeList of TypeString | Yes | No | No |
| `login_mode` | `LoginMode` | TypeInt | Yes | No | No |
| `domain_prefix` | `DomainPrefix` | TypeString | Yes | No | No |
| `accounts` | `Accounts` | TypeList of nested object | No | No | No |
| `anonymous_login` | `AnonymousLogin` | TypeList of nested object, MaxItems:1 | No | No | No |
| `intranet_type` | `IntranetType` | TypeInt | No | No | No |
| `intranet_region` | `IntranetRegion` | TypeString | No | No | No |
| `vpc_id` | `VpcId` | TypeString | No | No | No |
| `subnet_id` | `SubnetId` | TypeString | No | No | No |
| `auth_roles` | `AuthRoles` | TypeList of nested object | No | No | No |
| `tags` | `Tags` | TypeList of nested object | No | **Yes** | No |
| `hide_params` | `HideParams` | TypeList of TypeString | No | No | No |
| `access_control_rules` | `AccessControlRules` | TypeList of nested object | No | No | No |
| `remarks` | `Remarks` | TypeString | No | No | No |
| `menus` | `Menus` | TypeList of TypeString | No | No | No |
| `console_id` | `ConsoleId` | TypeString | — | — | **Yes** |
| `domain` | `Domain` | TypeString | — | — | **Yes** |
| `intranet_domain` | `IntranetDomain` | TypeString | — | — | **Yes** |

Inner schema for nested objects:

- `accounts` element: `user_name` (Optional, TypeString), `password` (Optional, TypeString, Sensitive), `secret_id` (Optional, TypeString, Sensitive), `secret_key` (Optional, TypeString, Sensitive), `email` (Optional, TypeString).
- `anonymous_login` element: `secret_id` (Optional, TypeString, Sensitive), `secret_key` (Optional, TypeString, Sensitive).
- `auth_roles` element: `role_name` (Optional, TypeString), `secret_id` (Optional, TypeString, Sensitive), `secret_key` (Optional, TypeString, Sensitive).
- `tags` element: `key` (Required, TypeString), `value` (Required, TypeString).
- `access_control_rules` element: `access_mode` (Optional, TypeString) — the only field on the SDK `AccessControlRule` struct.

Every nested element schema MUST set its inner string credentials (`password`, `secret_id`, `secret_key`) with `Sensitive: true` to keep them out of plan output.

#### Scenario: New required field appears in plan

- **WHEN** the user writes a config containing only `nat_gateway_id`-equivalent identity fields and omits `access_mode` / `login_mode` / `domain_prefix`
- **THEN** `terraform plan` SHALL fail validation pointing at the missing required attributes.

#### Scenario: Anonymous login uses a single nested object

- **WHEN** the user writes two `anonymous_login { ... }` blocks in HCL
- **THEN** `terraform plan` SHALL fail with a `MaxItems: 1` validation error.

### Requirement: Resource ID MUST be the API-returned ConsoleId

After Create, the resource's Terraform ID SHALL be set to `*response.Response.ConsoleId` (no compound separator). Read SHALL use `d.Id()` directly as the lookup key.

#### Scenario: ID is set from CreateConsole response

- **GIVEN** a successful `CreateConsole` returns `ConsoleId="ds-abc123"`
- **WHEN** Create completes
- **THEN** `d.Id()` returns `"ds-abc123"` and `console_id` state attribute equals `"ds-abc123"`.

#### Scenario: Create response missing ConsoleId is fatal

- **WHEN** `CreateConsole` returns a non-nil response with `ConsoleId == nil`
- **THEN** the Create callback SHALL return a non-retryable error mentioning `ConsoleId is nil`.

### Requirement: Read MUST locate the resource via DescribeConsoles + ConsoleId filter

The service-layer helper `(*ClsService).DescribeClsConsoleById(ctx, consoleId) (*cls.Console, error)` SHALL build a `DescribeConsolesRequest` with:

- `Limit = 100` (the API documented maximum)
- `Filters = [{Key: "ConsoleId", Values: [consoleId]}]`

…and invoke `UseClsClient().DescribeConsoles` inside `resource.Retry(ReadRetryTimeout, ...)`.

If `Response.Consoles` is empty, the helper SHALL return `(nil, nil)`. The Read callback SHALL detect that case and call `d.SetId("")` to mark the resource gone.

#### Scenario: Resource exists

- **GIVEN** the API holds a Console with `ConsoleId="ds-abc"`
- **WHEN** Read runs
- **THEN** all schema attributes are populated from the API response with nil-safe pointer dereferencing.

#### Scenario: Resource removed out-of-band

- **GIVEN** the user runs `terraform refresh` after the Console was deleted in the web console
- **WHEN** Read calls `DescribeConsoles` and `len(Response.Consoles) == 0`
- **THEN** Read SHALL `d.SetId("")` and return `nil` (no error), so the next plan proposes a re-create.

### Requirement: Update MUST call ModifyConsole when any mutable field changes

The Update callback SHALL detect changes against the field set:

```
access_mode, login_mode, domain_prefix, accounts, anonymous_login,
intranet_type, intranet_region, vpc_id, subnet_id, auth_roles,
hide_params, access_control_rules, remarks, menus
```

`tags` is **not** in this set (it is `ForceNew`). When any field in the set has changed, Update SHALL build a single `ModifyConsoleRequest` populated from current `d.Get(...)` values for ALL mutable fields (full overwrite, matching API semantics) and call `UseClsClient().ModifyConsole` inside `resource.Retry(WriteRetryTimeout, ...)`. Each call MUST include `ConsoleId = d.Id()`. Update MUST end by re-invoking Read.

#### Scenario: Editing only `remarks`

- **GIVEN** state has `remarks = "old"`
- **WHEN** the user changes `remarks` to `"new"`
- **THEN** Update issues exactly one `ModifyConsole` request with `ConsoleId = d.Id()`, `Remarks = "new"`, and the unchanged values of all other mutable fields.

#### Scenario: Editing `tags` triggers replacement

- **GIVEN** state has `tags = [{key="k", value="v"}]`
- **WHEN** the user changes `tags` to `[{key="k2", value="v2"}]`
- **THEN** Terraform's plan reports a destroy + create cycle (resource is replaced because `tags` is ForceNew). Update is NOT invoked.

### Requirement: Delete MUST call DeleteConsole

The Delete callback SHALL build a `DeleteConsoleRequest` with `ConsoleId = d.Id()` and call `UseClsClient().DeleteConsole` inside `resource.Retry(WriteRetryTimeout, ...)`. No additional fields are sent.

#### Scenario: Idempotent delete on a missing resource

- **GIVEN** the resource is already gone in the cloud (e.g. external deletion between Read and Delete)
- **WHEN** Delete is called
- **THEN** the SDK error is propagated to Terraform, which surfaces it to the user; the provider does NOT swallow the error to avoid masking real problems.

### Requirement: Every API call MUST be wrapped in resource.Retry

Every invocation of `UseClsClient().CreateConsole`, `DescribeConsoles`, `ModifyConsole`, `DeleteConsole` SHALL be wrapped in `resource.Retry(...)` with the appropriate timeout (`ReadRetryTimeout` for read paths, `WriteRetryTimeout` for write paths) and SHALL forward errors via `tccommon.RetryError(e)` so transient errors are retried per the standard provider retry policy.

#### Scenario: Transient SDK error is retried

- **GIVEN** the first invocation of `CreateConsole` returns a retriable error (e.g. internal error)
- **WHEN** the retry callback re-runs
- **THEN** the second attempt's response is observed; the original error is surfaced only after the retry budget is exhausted.

### Requirement: All response field reads MUST be nil-safe

Every `_ = d.Set("<key>", respData.<Field>)` access SHALL be guarded by `if respData.<Field> != nil` before dereference, including nested struct fields (e.g. `respData.AnonymousLogin.SecretId`). Map and slice flatten loops SHALL also nil-check each item before reading.

#### Scenario: Optional string field absent in response

- **GIVEN** the API returns `Remarks == nil`
- **WHEN** Read runs
- **THEN** `d.Set("remarks", ...)` is NOT called for that field; state retains whatever the previous value was (or stays unset for fresh imports).

### Requirement: Documentation and acceptance test MUST follow project conventions

- The HCL example markdown SHALL live at `tencentcloud/services/cls/resource_tc_cls_console.md` (mirroring `resource_tc_config_compliance_pack.md`) and contain at least one full HCL example using account-password login and one Import example.
- An acceptance test SHALL live at `tencentcloud/services/cls/resource_tc_cls_console_test.go` (mirroring `resource_tc_config_compliance_pack_test.go`) covering: basic Create, Update of `remarks`, and `ImportState` round-trip.
- Running `make doc` SHALL regenerate `website/docs/r/cls_console.html.markdown`.

#### Scenario: Generated website doc lists the resource

- **WHEN** `make doc` runs
- **THEN** `website/docs/r/cls_console.html.markdown` exists and lists every schema attribute defined in the Schema requirement above.

#### Scenario: Acceptance test name and structure

- **WHEN** the test file is opened
- **THEN** the package is `cls_test`, the test function is `TestAccTencentCloudClsConsoleResource_basic`, and it includes `tcacctest.AccPreCheck`, two HCL configs (initial and updated), and an `ImportState` step with `ImportStateVerify: true`.
