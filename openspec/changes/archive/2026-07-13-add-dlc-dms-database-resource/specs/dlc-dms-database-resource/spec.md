# dlc-dms-database-resource Specification

## ADDED Requirements

### Requirement: Resource MUST be registered as `tencentcloud_dlc_dms_database`

The provider SHALL register a new RESOURCE_KIND_GENERAL resource named `tencentcloud_dlc_dms_database` whose Create/Read/Update/Delete callbacks invoke the DMS Database APIs of `tencentcloud-sdk-go/tencentcloud/dlc/v20210125`.

#### Scenario: Resource registered in provider map

- **WHEN** the provider is loaded
- **THEN** `provider.go` exposes the resource via key `"tencentcloud_dlc_dms_database"` mapped to `dlc.ResourceTencentCloudDlcDmsDatabase()`, alongside existing `tencentcloud_dlc_*` resource entries.

#### Scenario: Resource appears in gendoc index

- **WHEN** `tencentcloud/provider.md` is scanned by `make doc`
- **THEN** the DLC Resource section MUST include `tencentcloud_dlc_dms_database` so that `website/docs/r/dlc_dms_database.html.markdown` is generated.

### Requirement: Schema MUST mirror the DMS Database API input and output

The resource schema SHALL declare these top-level argument keys, with semantics matching the SDK request/response fields:

| HCL key | SDK field | Type | Required | ForceNew | Computed |
|---|---|---|---|---|---|
| `name` | `Name` | TypeString | Yes | Yes | No |
| `schema_name` | `SchemaName` | TypeString | Yes | Yes | No |
| `datasource_connection_name` | `DatasourceConnectionName` | TypeString | Yes | Yes | No |
| `location` | `Location` | TypeString | Optional | No | No |
| `delete_data` | `DeleteData` | TypeBool | Optional | No | No |
| `cascade` | `Cascade` | TypeBool | Optional | No | No |
| `asset` | `Asset` | TypeList, MaxItems=1 | Optional | No | No |

The `asset` block SHALL be a single-element list whose nested resource declares: `id`(TypeInt, Computed), `name`(TypeString), `guid`(TypeString, Computed), `catalog`(TypeString), `description`(TypeString), `owner`(TypeString), `owner_account`(TypeString), `perm_values`(TypeList of KVPair), `params`(TypeList of KVPair), `biz_params`(TypeList of KVPair), `data_version`(TypeInt, Computed), `create_time`(TypeString, Computed), `modified_time`(TypeString, Computed), `datasource_id`(TypeInt, Computed). Each KVPair list element SHALL declare `key`(TypeString, Required) and `value`(TypeString, Optional).

The schema MUST NOT expose `Pattern` (a fuzzy-match parameter used only internally by Read for precise lookups).

#### Scenario: Required fields enforce on plan

- **WHEN** the user writes a config that omits `name`, `schema_name`, or `datasource_connection_name`
- **THEN** `terraform plan` SHALL fail validation pointing at the missing required attribute.

#### Scenario: Changing a ForceNew field forces replacement

- **GIVEN** state has `name = "db1"`
- **WHEN** the user changes `name` to `"db2"` in HCL
- **THEN** Terraform's plan reports a destroy + create cycle (the resource is replaced because `name` is ForceNew).

#### Scenario: Asset nested block accepts a single configuration

- **WHEN** the user declares an `asset` block with `name`, `description`, `owner`, and one `params` entry `{ key = "k", value = "v" }`
- **THEN** `terraform plan` accepts the configuration without error.

### Requirement: Resource ID MUST be a compound of name, schema_name and datasource_connection_name

After Create, the Terraform resource ID SHALL be `name + "#" + schema_name + "#" + datasource_connection_name` using `tccommon.FILED_SP` ("#") as the separator. Read, Update, and Delete SHALL split `d.Id()` by `tccommon.FILED_SP` to recover the three parts for API requests.

#### Scenario: ID is composed from the three identity fields

- **GIVEN** HCL declares `name = "mydb"`, `schema_name = "myschema"`, `datasource_connection_name = "conn1"`
- **WHEN** Create completes
- **THEN** `d.Id()` returns `"mydb#myschema#conn1"`.

#### Scenario: Import requires the compound id

- **GIVEN** an existing DMS database `mydb` under schema `myschema` with datasource connection `conn1`
- **WHEN** the user runs `terraform import tencentcloud_dlc_dms_database.foo "mydb#myschema#conn1"`
- **THEN** Read populates state from `name="mydb"`, `schema_name="myschema"`, `datasource_connection_name="conn1"`.

### Requirement: Create MUST call CreateDMSDatabase and set the compound ID

The Create callback SHALL build a `CreateDMSDatabaseRequest` populated from HCL `name`, `schema_name`, `datasource_connection_name`, `location` (if set), and `asset` (if set), and call `UseDlcClient().CreateDMSDatabaseWithContext(ctx, request)` inside `resource.Retry(tccommon.WriteRetryTimeout, ...)`. After success it MUST verify `response != nil && response.Response != nil` (returning `NonRetryableError` otherwise), print the logId and name, set `d.SetId(<compound id>)`, then call Read.

#### Scenario: Create issues a single CreateDMSDatabase call

- **GIVEN** HCL with `name`, `schema_name`, `datasource_connection_name`, `location`, and an `asset` block
- **WHEN** the user runs `terraform apply`
- **THEN** exactly one `CreateDMSDatabase` request is issued with the supplied fields; `d.Id()` returns the compound id.

#### Scenario: Nil response is detected

- **GIVEN** the SDK returns `(result == nil, err == nil)`
- **WHEN** the retry callback runs
- **THEN** the callback returns `resource.NonRetryableError(...)` and does NOT set the resource id.

### Requirement: Read MUST call DescribeDMSDatabase with precise identity

The Read callback SHALL build a `DescribeDMSDatabaseRequest` from the split id parts (`name`, `schema_name`, `datasource_connection_name`) WITHOUT setting `Pattern`, and call `UseDlcClient().DescribeDMSDatabaseWithContext(ctx, request)` inside `resource.Retry(tccommon.ReadRetryTimeout, ...)`. Inside the retry block, if `response == nil || response.Response == nil`, it MUST return `NonRetryableError` (it MUST NOT clear `d.Id()` there). After retry, if the database is not found, it SHALL first `log.Printf("[CRUD] dlc_dms_database id=%s", d.Id())` and then `d.SetId("")`. Each `d.Set(...)` MUST be nil-safe against the response field.

#### Scenario: Transient read error is retried without clearing id

- **GIVEN** the first `DescribeDMSDatabase` call returns a retriable error
- **WHEN** the retry callback re-runs
- **THEN** the second attempt's response is observed; `d.Id()` is NOT cleared while retries remain.

#### Scenario: Missing database clears id after logging

- **GIVEN** `DescribeDMSDatabase` returns an empty `Response` (Name is nil)
- **WHEN** Read completes outside retry
- **THEN** a `[CRUD] dlc_dms_database id=...` log line is emitted and `d.SetId("")` is called.

### Requirement: Update MUST call AlterDMSDatabase for mutable fields

The Update callback SHALL treat `name`, `schema_name`, `datasource_connection_name` as immutable (list them in `immutableArgs` and return an error if changed). For changes to `location` or `asset`, it SHALL build an `AlterDMSDatabaseRequest` with `CurrentName` set to the current `name`, plus `SchemaName`, `DatasourceConnectionName`, `Location` (if set), and `Asset` (if set), then call `UseDlcClient().AlterDMSDatabaseWithContext(ctx, request)` inside `resource.Retry(tccommon.WriteRetryTimeout, ...)`. Update MUST end by calling Read.

#### Scenario: Editing location triggers AlterDMSDatabase

- **GIVEN** state has `location = "cosn://old/path"`
- **WHEN** the user changes `location` to `"cosn://new/path"`
- **THEN** Update issues exactly one `AlterDMSDatabase` request with `CurrentName = <state name>`, `Location = "cosn://new/path"`.

#### Scenario: Changing immutable id field is rejected

- **GIVEN** a config change attempts to alter `datasource_connection_name` (ForceNew already forces replacement)
- **WHEN** the immutableArgs check runs
- **THEN** the Update returns an error for the immutable argument (defense-in-depth alongside ForceNew).

### Requirement: Delete MUST call DropDMSDatabase

The Delete callback SHALL build a `DropDMSDatabaseRequest` with `Name` and `DatasourceConnectionName` recovered from the split id, plus `DeleteData` and `Cascade` from the HCL config (defaulting to false when unset), and call `UseDlcClient().DropDMSDatabaseWithContext(ctx, request)` inside `resource.Retry(tccommon.WriteRetryTimeout, ...)`.

#### Scenario: Delete propagates user options

- **GIVEN** HCL declares `delete_data = true` and `cascade = true`
- **WHEN** the user runs `terraform destroy`
- **THEN** exactly one `DropDMSDatabase` request is issued with `Name=<id name>`, `DatasourceConnectionName=<id conn>`, `DeleteData=true`, `Cascade=true`.

#### Scenario: Delete defaults are false

- **GIVEN** HCL omits `delete_data` and `cascade`
- **WHEN** the user runs `terraform destroy`
- **THEN** the `DropDMSDatabase` request is issued with `DeleteData=false` and `Cascade=false`.

### Requirement: Every API call MUST be wrapped in resource.Retry

All four SDK invocations (`CreateDMSDatabaseWithContext`, `DescribeDMSDatabaseWithContext`, `AlterDMSDatabaseWithContext`, `DropDMSDatabaseWithContext`) SHALL be wrapped in `resource.Retry(tccommon.<Read|Write>RetryTimeout, ...)` and forward errors via `tccommon.RetryError(e)`.

#### Scenario: Retry wrapper present

- **WHEN** a code reviewer inspects each CRUD callback
- **THEN** the SDK call is inside a `resource.Retry(...)` block, not a bare invocation.

### Requirement: Response field reads MUST be nil-safe

For every `d.Set(...)` in Read, the implementation SHALL guard the corresponding response field against nil before setting. Inside retry blocks for Read, a nil response MUST return `NonRetryableError` rather than clearing the id.

#### Scenario: Nil response is detected in Read retry

- **GIVEN** `DescribeDMSDatabase` returns `(result == nil, err == nil)`
- **WHEN** the retry callback runs
- **THEN** it returns `resource.NonRetryableError(...)` and does NOT panic or clear the id.

### Requirement: Documentation and tests MUST follow project conventions

- The HCL example markdown SHALL live at `tencentcloud/services/dlc/resource_tc_dlc_dms_database.md`. It SHALL contain: a one-line description mentioning DLC, an Example Usage section (using `jsonencode()` for any json-string field values if applicable), and an Import section explaining the compound id (`name#schema_name#datasource_connection_name`). It MUST NOT include `Argument Reference` or `Attribute Reference` sections (auto-generated).
- The unit test SHALL live at `tencentcloud/services/dlc/resource_tc_dlc_dms_database_test.go` using gomonkey to mock the DLC cloud API (NOT the terraform test suite), testing business logic only, and be runnable via `go test -gcflags=all=-l`.
- Running `make doc` SHALL regenerate `website/docs/r/dlc_dms_database.html.markdown`.

#### Scenario: Generated website doc lists the resource

- **WHEN** `make doc` runs
- **THEN** `website/docs/r/dlc_dms_database.html.markdown` exists and lists every schema attribute defined in the Schema requirement.

#### Scenario: Unit test uses gomonkey mocks

- **WHEN** the test file is opened
- **THEN** the package is `dlc_test`, the test functions mock `CreateDMSDatabase`/`DescribeDMSDatabase`/`AlterDMSDatabase`/`DropDMSDatabase` via gomonkey, and no `TF_ACC` acceptance test is required.
