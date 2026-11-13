## 1. Connectivity & Provider Registration

- [x] 1.1 Add the ioa v20220601 SDK import and a `ioaV20220601Conn *ioav20220601.Client` field to `tencentcloud/connectivity/client.go`
- [x] 1.2 Add `UseIoaV20220601Client()` accessor method to `tencentcloud/connectivity/client.go` (mirror `UseIgtmV20231024Client`, cached + `LogRoundTripper`)
- [x] 1.3 Register `tencentcloud_ioa_company_directory_config` resource in `tencentcloud/provider.go` (map key + `ResourceTencentCloudIoaCompanyDirectoryConfig()` factory)
- [x] 1.4 Add the `tencentcloud_ioa_company_directory_config` entry to `tencentcloud/provider.md`

## 2. Service Layer

- [x] 2.1 Create `tencentcloud/services/ioa/service_tencentcloud_ioa.go` with `IoaService` struct and `NewIoaService` constructor
- [x] 2.2 Implement `DescribeCompanyDirectoryConfigById(ctx, id)` helper wrapping `DescribeCompanyDirectoryConfig` in `resource.Retry(tccommon.ReadRetryTimeout, …)` with `tccommon.RetryError`, returning `*ioav20220601.DirectoryConfigData`

## 3. Resource Schema & CRUD

- [x] 3.1 Create `tencentcloud/services/ioa/resource_tc_ioa_company_directory_config.go` with `ResourceTencentCloudIoaCompanyDirectoryConfig()` defining the flat schema:
  - Required: `type`, `name`, `config` (Sensitive), `sync_enable`, `sync_policy`, `sync_policy_params`, `create_auth_config`, `display_on_login_page`
  - Optional: `description`, `scene`, `name_i18n` (TypeList of {lang, value})
  - Computed: `id`, `source_id`, `identify_source_id`, `auth_source_id`, `auth_config_id`, `auth_policy_id`, `auth_support_platforms`, `auth_methods`
  - Importer: `schema.ImportStatePassthrough`
- [x] 3.2 Implement `resourceTencentCloudIoaCompanyDirectoryConfigCreate`:
  - Build `CreateCompanyDirectoryConfigRequest` from schema (including `name_i18n` → `[]*I18nString`)
  - Call SDK inside `resource.Retry(tccommon.WriteRetryTimeout, …)` with `tccommon.RetryError`
  - Validate `Response`/`Response.Data`/`Response.Data.Id` non-nil/non-empty; on empty return `NonRetryableError` (log logId + d.Id first)
  - `d.SetId(id)` then set computed result fields from `DirectoryConfigResultData`; call Read
- [x] 3.3 Implement `resourceTencentCloudIoaCompanyDirectoryConfigRead`:
  - Parse `d.Id()`, call `IoaService.DescribeCompanyDirectoryConfigById`
  - On empty `Data`: `log.Printf("[CRUD] ioa_company_directory_config id=%s", d.Id())` then `d.SetId("")`
  - Set each non-nil `DirectoryConfigData` field into state (skip nil fields)
- [x] 3.4 Implement `resourceTencentCloudIoaCompanyDirectoryConfigUpdate`:
  - Detect changes across mutable fields (`type`, `name`, `config`, `sync_enable`, `sync_policy`, `sync_policy_params`, `create_auth_config`, `display_on_login_page`, `description`, `name_i18n`)
  - On change, build `ModifyCompanyDirectoryConfigRequest` with full config + `Id`, call SDK inside `resource.Retry(tccommon.WriteRetryTimeout, …)` with `tccommon.RetryError`
  - Set computed result fields from `DirectoryConfigResultData` when present; call Read
- [x] 3.5 Implement `resourceTencentCloudIoaCompanyDirectoryConfigDelete`:
  - Parse `d.Id()` to uint64, call `DeleteAccountGroup` with `AccountGroupId` (omit `DomainInstanceId`) inside `resource.Retry(tccommon.WriteRetryTimeout, …)` with `tccommon.RetryError`

## 4. Documentation Example

- [x] 4.1 Create `tencentcloud/services/ioa/resource_tc_ioa_company_directory_config.md` with one-line description (mention IOA), Example Usage (use `jsonencode()` for any JSON-string fields), and Import section (import by bare directory id)

## 5. Unit Tests (gomonkey mocks)

- [x] 5.1 Create `tencentcloud/services/ioa/resource_tc_ioa_company_directory_config_test.go` using gomonkey to mock the ioa v20220601 client methods (Create/Describe/Modify/Delete), covering Create/Read/Update/Delete business logic — no Terraform acceptance suite

## 6. Verification (separate from code edits)

- [ ] 6.1 Run `gofmt` on all new/modified Go files (finalization phase only, via tfpacer-finalize skill)
- [ ] 6.2 Run `make doc` to generate `website/docs/` documentation (finalization phase only, via tfpacer-finalize skill)
- [x] 6.3 Confirm the change proposal passes `openspec validate add-ioa-company-directory-config-resource`
