## Context

The TencentCloud IOA service (`github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ioa/v20220601`) exposes APIs to manage enterprise directory configurations (identity sources such as WeCom, Lark, DingTalk, MicrosoftEntraID). The Terraform provider currently has no `tencentcloud_ioa_*` resource and no registered IOA SDK client, so enterprise directories cannot be managed declaratively.

The vendored SDK provides four relevant APIs:
- `CreateCompanyDirectoryConfig(request)` → returns `Response.Data` of type `DirectoryConfigResultData` whose `Id` (`*int64`) is the new directory identifier. Optional fields include `Description`, `Scene`, and `NameI18n []*I18nString`.
- `DescribeCompanyDirectoryConfig(request)` with `Id` (`*int64`) → returns `Response.Data` of type `DirectoryConfigData`, the full configuration used to hydrate state.
- `ModifyCompanyDirectoryConfig(request)` → accepts the full configuration (Create-style params plus `Id`) and returns `Response.Data` (`DirectoryConfigResultData`).
- `DeleteAccountGroup(request)` with `DomainInstanceId` (`*string`) and `AccountGroupId` (`*uint64`) → deletes an account group / directory. There is no dedicated `DeleteCompanyDirectoryConfig`; this API is the documented delete path for a directory.

The provider follows a RESOURCE_KIND_GENERAL pattern (reference: `tencentcloud_igtm_strategy`): a single `resource_tc_<service>_<name>.go` implements Create/Read/Update/Delete directly against the SDK client, plus a thin service layer for the Read query with retry, and the SDK client accessor lives in `tencentcloud/connectivity/client.go`.

## Goals / Non-Goals

**Goals:**
- Provide a GENERAL resource `tencentcloud_ioa_company_directory_config` with full CRUD lifecycle backed by the four ioa v20220601 APIs above.
- Register a new IOA SDK client accessor `UseIoaV20220601Client()` so the resource (and future ioa resources) can call the SDK.
- Add a service-layer `DescribeCompanyDirectoryConfigById` helper using the standard `resource.Retry(tccommon.ReadRetryTimeout, …)` + `tccommon.RetryError` pattern.
- Keep the resource schema flat (no synthetic wrapping list), mapping schema attributes 1:1 to the API fields so each field can be `set`/`read` independently.
- Add gomonkey-mocked unit tests (no Terraform acceptance suite) per the codegen requirements.
- Provide `.md` example content for `make doc`.

**Non-Goals:**
- Do not implement any IOA data source (RESOURCE_KIND_DATASOURCE) in this change.
- Do not add `Timeouts` block / async poll: none of the four APIs are documented as asynchronous; they return synchronously, so no Read-poll loop is needed after Create/Update/Delete.
- Do not introduce `_extension.go` files.
- Do not modify any existing resource schema or state.

## Decisions

### D1: Resource ID strategy — bare directory `Id`
The resource ID SHALL be the integer directory `Id` returned by `CreateCompanyDirectoryConfig.Response.Data.Id`, stored as a string. The directory is uniquely identified by a single `Id` across Describe/Modify/Delete, so there is no need for a composite (`FIELD_SP`-joined) ID.

**Rationale:** `DescribeCompanyDirectoryConfig` and `ModifyCompanyDirectoryConfig` both key off the single `Id` field, and `DeleteAccountGroup` keys off `AccountGroupId` which is the same directory id. A single-field ID keeps import trivial (`terraform import tencentcloud_ioa_company_directory_config.x <id>`).

**Alternatives considered:** composite `id` — rejected because there is no second identifying component to encode.

### D2: Delete via `DeleteAccountGroup`
The Delete operation SHALL call `DeleteAccountGroup` with `AccountGroupId` set to the directory id (parsed to `uint64`) and `DomainInstanceId` left unset (the API treats absence as the root domain "1", covering the common case). The resource supports `terraform import`, so import uses the bare directory id.

**Rationale:** The SDK has no `DeleteCompanyDirectoryConfig`; `DeleteAccountGroup` is the documented API for deleting an account group / directory. `DomainInstanceId` is optional ("可直接传入根域'1'"), and omitting it is the least-surprise default for the resource lifecycle.

**Alternatives considered:** Require users to pass `domain_instance_id` into delete — rejected because it is not a natural attribute of a directory config resource and would force an extra schema field that does not belong to Create/Read.

### D3: Schema shape — flat top-level fields, `name_i18n` as a list-of-objects
The schema SHALL expose every Create/Describe input as a top-level attribute:
- Required: `type` (string), `name` (string), `config` (string, sensitive), `sync_enable` (bool), `sync_policy` (string), `sync_policy_params` (string), `create_auth_config` (bool), `display_on_login_page` (bool).
- Optional: `description` (string), `scene` (string — Create-only, not in Describe output), `name_i18n` (TypeList of objects with `lang` + `value`, optional).
- Computed/read-only (hydrated from Describe `DirectoryConfigData`): `source_id` (string). Computed fields from the Create/Modify result `DirectoryConfigResultData` (`identify_source_id`, `auth_source_id`, `auth_config_id`, `auth_policy_id`, `auth_support_platforms`, `auth_methods`) are also exposed as computed because the API returns them on write responses but NOT on Describe; they are set only after Create/Update when present.

**Rationale:** Mapping each API field to its own schema field is the established convention and satisfies the "no synthetic wrapping list" rule. `config` holds SM2-encrypted+Hex configuration and should be marked `Sensitive: true`. `name_i18n` mirrors the SDK's `[]*I18nString` (Lang, Value).

**Alternatives considered:** Flatten `lang`/`value` to the top level — rejected because `NameI18n` is an array and flattening loses the per-entry pairing.

### D4: Update sends the full configuration
Update SHALL detect changes across the mutable top-level fields and call `ModifyCompanyDirectoryConfig` with the full configuration (all Required + Optional params the API accepts, plus `Id`). The API explicitly requires the complete config to be resent.

**Rationale:** Matches `ModifyCompanyDirectoryConfig` semantics and the existing GENERAL-resource update pattern (see `tencentcloud_igtm_strategy`).

### D5: New IOA client accessor
Add `ioaV20220601Conn *ioav20220601.Client` field, the `ioav20220601` import, and a `UseIoaV20220601Client()` accessor in `tencentcloud/connectivity/client.go`, mirroring `UseIgtmV20231024Client`.

**Rationale:** No IOA client exists today; every resource accesses the SDK through such an accessor.

### D6: Read-empty handling
In Read, if `DescribeCompanyDirectoryConfig` returns empty data (`Data == nil`), the code SHALL first `log.Printf("[CRUD] ioa_company_directory_config id=%s", d.Id())` and then `d.SetId("")`, preserving the id in logs as required by the codegen rules.

### D7: Create result validation
After Create, the code SHALL validate `response.Response != nil` and `response.Response.Data != nil` and `Data.Id != nil` (and non-zero/non-empty). On any "empty" form it SHALL return `NonRetryableError` so an empty id is never written to state.

## Risks / Trade-offs

- **[DeleteAccountGroup semantics]** `DeleteAccountGroup` deletes "an account group OR a directory"; passing the directory id as `AccountGroupId` relies on the directory id being accepted in that role. → Mitigation: This is the documented delete path for directories; if a future SDK adds a dedicated delete API, the resource can be migrated transparently since the id contract is unchanged.
- **[Scene is Create-only]** `Scene` appears only on Create and is absent from `DescribeCompanyDirectoryConfig` output, so it cannot be round-tripped into state. → Mitigation: Mark `scene` as `Optional` without `Computed`; do not attempt to `d.Set` it in Read. It will not trigger spurious diffs because Terraform keeps the configured value.
- **[Write-only computed fields]** `DirectoryConfigResultData` fields (`identify_source_id`, `auth_source_id`, `auth_config_id`, `auth_policy_id`, `auth_support_platforms`, `auth_methods`) are returned by Create/Modify but not by Describe. → Mitigation: Expose them as `Computed` and populate them after Create/Update; Read (backed by Describe) will leave them as last-known. This matches how other GENERAL resources handle write-response-only data.
- **[DomainInstanceId omitted on delete]** Omitting `DomainInstanceId` assumes the root domain. → Mitigation: Acceptable default for the common case; documented in the resource description.
