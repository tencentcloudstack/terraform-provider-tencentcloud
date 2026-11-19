## ADDED Requirements

### Requirement: Resource registration
The provider SHALL expose a resource type named `tencentcloud_ioa_company_directory_config` that manages a single TencentCloud IOA enterprise directory configuration per resource block. The resource MUST be registered in `tencentcloud/provider.go` under the `ioa` namespace and documented in `tencentcloud/provider.md`.

#### Scenario: Resource type is discoverable
- **WHEN** an operator runs `terraform plan` against a configuration that references `resource "tencentcloud_ioa_company_directory_config" "<name>"`
- **THEN** Terraform resolves the type without an "unknown resource" error and shows the planned create.

#### Scenario: Provider compiles
- **WHEN** the codebase is built with `go build ./tencentcloud/...`
- **THEN** the build succeeds with no compilation errors related to the new resource.

### Requirement: IOA SDK client accessor
The connectivity layer in `tencentcloud/connectivity/client.go` SHALL provide a `UseIoaV20220601Client()` method that lazily creates and returns a cached `*ioav20220601.Client` (package `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ioa/v20220601`), mirroring the existing `UseIgtmV20231024Client` pattern.

#### Scenario: Client is cached
- **WHEN** `UseIoaV20220601Client()` is called multiple times on the same `TencentCloudClient`
- **THEN** the same cached `*ioav20220601.Client` instance is returned without re-initializing.

### Requirement: Schema mirrors CreateCompanyDirectoryConfig
The resource schema SHALL expose every input parameter accepted by `CreateCompanyDirectoryConfig` as a top-level attribute:
- `type` (string, required) — enterprise directory type (e.g. `WeCom`, `Lark`, `DingTalk`, `MicrosoftEntraID`).
- `name` (string, required) — enterprise directory name.
- `config` (string, required, sensitive) — SM2-encrypted-then-Hex configuration data.
- `sync_enable` (bool, required) — whether scheduled sync is enabled.
- `sync_policy` (string, required) — sync strategy (`4hours` / `daily` / `weekly`).
- `sync_policy_params` (string, required) — JSON string of sync strategy parameters.
- `create_auth_config` (bool, required) — whether to synchronously create an auth source.
- `display_on_login_page` (bool, required) — whether to display on the login page.
- `description` (string, optional) — directory description.
- `scene` (string, optional) — usage scene (API created / quick start / normal config); Create-only, not present in Describe output.
- `name_i18n` (list of objects, optional) — multi-language names, each element has `lang` (string) and `value` (string).

The resource SHALL NOT introduce a synthetic wrapping list that nests these fields under a single parent attribute.

#### Scenario: Required SDK fields are present in schema
- **WHEN** a developer inspects the resource schema
- **THEN** every field declared in `ioav20220601.CreateCompanyDirectoryConfigRequest` (Type, Name, Config, SyncEnable, SyncPolicy, SyncPolicyParams, CreateAuthConfig, DisplayOnLoginPage, Description, Scene, NameI18n) appears in the schema with semantically equivalent typing.

#### Scenario: Config field is sensitive
- **WHEN** a developer inspects the `config` attribute
- **THEN** the attribute is marked `Sensitive: true` because it carries encrypted configuration data.

#### Scenario: name_i18n is a list of objects
- **WHEN** a developer inspects the `name_i18n` attribute
- **THEN** it is a `schema.TypeList` whose element schema contains `lang` and `value` string fields, matching `ioav20220601.I18nString`.

### Requirement: Computed and read-only attributes
The resource SHALL expose the following read-only attributes:
- `id` (string, computed) — the enterprise directory id; also assigned as the resource ID.
- `source_id` (string, computed) — hydrated from `DescribeCompanyDirectoryConfig.Response.Data.SourceId`.
- `identify_source_id` (string, computed) — from Create/Modify result `DirectoryConfigResultData.IdentifySourceId`.
- `auth_source_id` (string, computed) — from `DirectoryConfigResultData.AuthSourceId`.
- `auth_config_id` (int, computed) — from `DirectoryConfigResultData.AuthConfigId`.
- `auth_policy_id` (int, computed) — from `DirectoryConfigResultData.AuthPolicyId`.
- `auth_support_platforms` (list of string, computed) — from `DirectoryConfigResultData.AuthSupportPlatforms`.
- `auth_methods` (list of string, computed) — from `DirectoryConfigResultData.AuthMethods`.

#### Scenario: Read-only attributes are not settable
- **WHEN** a user writes a configuration that assigns values to `source_id`, `identify_source_id`, `auth_source_id`, `auth_config_id`, `auth_policy_id`, `auth_support_platforms`, or `auth_methods`
- **THEN** Terraform reports these as computed/read-only and rejects explicit values.

### Requirement: Create operation
The Create function SHALL call `CreateCompanyDirectoryConfig` with all configured input parameters, wrapped in `resource.Retry(tccommon.WriteRetryTimeout, …)` with errors passed through `tccommon.RetryError`. After a successful response, the function MUST validate that `Response` and `Response.Data` and `Response.Data.Id` are non-nil/non-empty; if any is empty it SHALL return a `NonRetryableError` and MUST NOT call `d.SetId`. On success the directory id (from `Response.Data.Id`) SHALL be assigned as the resource ID, after which the computed result fields SHALL be set into state and Read SHALL be invoked to refresh state.

#### Scenario: Successful create sets the resource id
- **WHEN** `CreateCompanyDirectoryConfig` returns a non-nil `Response.Data.Id`
- **THEN** the resource calls `d.SetId(<id-as-string>)` using the returned directory id.

#### Scenario: Empty create response does not corrupt state
- **WHEN** `CreateCompanyDirectoryConfig` returns `Response == nil` or `Response.Data == nil` or `Response.Data.Id == nil`
- **THEN** the Create function returns a `NonRetryableError` describing the empty result and never calls `d.SetId`.

#### Scenario: Create API errors are retried
- **WHEN** `CreateCompanyDirectoryConfig` returns a retryable SDK error
- **THEN** the error is wrapped with `tccommon.RetryError` so the retry loop continues until `WriteRetryTimeout` is exhausted.

### Requirement: Read operation
The Read function SHALL call the service-layer helper `DescribeCompanyDirectoryConfigById` (which wraps `DescribeCompanyDirectoryConfig` with `resource.Retry(tccommon.ReadRetryTimeout, …)`). Each attribute returned in `Response.Data` (`DirectoryConfigData`) SHALL be set into state only when the corresponding field is non-nil. If the API returns empty data (`Data == nil`), the function SHALL first log `log.Printf("[CRUD] ioa_company_directory_config id=%s", d.Id())` and then call `d.SetId("")`.

#### Scenario: Existing resource is hydrated
- **WHEN** Read is called for an existing directory whose `DescribeCompanyDirectoryConfig.Response.Data` is populated
- **THEN** every non-nil field in `DirectoryConfigData` is written to the corresponding state attribute.

#### Scenario: Missing resource clears state with a log
- **WHEN** `DescribeCompanyDirectoryConfig.Response.Data` is nil
- **THEN** the function logs the current id via `log.Printf("[CRUD] ioa_company_directory_config id=%s", d.Id())` before calling `d.SetId("")`, so the id is preserved in logs.

#### Scenario: nil fields are skipped
- **WHEN** a field in `DirectoryConfigData` is nil
- **THEN** the corresponding `d.Set(...)` call is skipped for that field.

### Requirement: Update operation
The Update function SHALL detect changes to the mutable top-level fields and call `ModifyCompanyDirectoryConfig` with the full configuration (all Required + Optional API params plus `Id` parsed from the resource id), wrapped in `resource.Retry(tccommon.WriteRetryTimeout, …)` with `tccommon.RetryError`. After a successful Modify, the computed result fields from `DirectoryConfigResultData` SHALL be set when present, and Read SHALL be invoked to refresh state.

#### Scenario: Changed field triggers Modify
- **WHEN** any of `type`, `name`, `config`, `sync_enable`, `sync_policy`, `sync_policy_params`, `create_auth_config`, `display_on_login_page`, `description`, `name_i18n` changes
- **THEN** `ModifyCompanyDirectoryConfig` is called once with the complete updated configuration and the directory `Id`.

#### Scenario: No tracked change skips Modify
- **WHEN** none of the mutable fields change
- **THEN** `ModifyCompanyDirectoryConfig` is not called and Read is still invoked to refresh state.

### Requirement: Delete operation
The Delete function SHALL call `DeleteAccountGroup` with `AccountGroupId` set to the directory id parsed to `uint64`, wrapped in `resource.Retry(tccommon.WriteRetryTimeout, …)` with `tccommon.RetryError`. `DomainInstanceId` SHALL be left unset (root domain default). On success the function returns without error.

#### Scenario: Successful delete
- **WHEN** `DeleteAccountGroup` succeeds for an existing directory
- **THEN** the Delete function returns nil and the resource is removed from Terraform state.

#### Scenario: Delete API errors are retried
- **WHEN** `DeleteAccountGroup` returns a retryable SDK error
- **THEN** the error is wrapped with `tccommon.RetryError` so the retry loop continues.

### Requirement: Import support
The resource SHALL support `terraform import` using the bare enterprise directory id, via `schema.ImportStatePassthrough`.

#### Scenario: Import by id
- **WHEN** an operator runs `terraform import tencentcloud_ioa_company_directory_config.x <directory-id>`
- **THEN** the resource state is hydrated from `DescribeCompanyDirectoryConfig` using the imported id, with no manual composite-id parsing required.

### Requirement: Documentation and tests
The change SHALL add a usage example file `tencentcloud/services/ioa/resource_tc_ioa_company_directory_config.md` (consumed by `make doc`) and a unit-test file `tencentcloud/services/ioa/resource_tc_ioa_company_directory_config_test.go` that uses gomonkey mocks for the ioa SDK client (no Terraform acceptance test suite), covering Create/Read/Update/Delete business logic.

#### Scenario: Example markdown exists
- **WHEN** `make doc` is run
- **THEN** a docs page for `tencentcloud_ioa_company_directory_config` is generated from the `.md` example file.

#### Scenario: Unit tests mock the SDK
- **WHEN** the generated `*_test.go` is compiled
- **THEN** it uses `gomonkey` to mock the ioa v20220601 client methods and asserts CRUD business logic without calling the real cloud API.
