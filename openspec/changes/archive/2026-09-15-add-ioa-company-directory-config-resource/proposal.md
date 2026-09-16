## Why

The TencentCloud IOA (Identity and Access Management / Open APIs) service exposes a `CreateCompanyDirectoryConfig` API for creating enterprise directory configurations (e.g. WeCom, Lark, DingTalk, MicrosoftEntraID identity sources), but the Terraform provider currently has no `tencentcloud_ioa_*` resource to manage them. Operators therefore cannot manage enterprise directory lifecycle (create / read / update / delete) through infrastructure-as-code. This change introduces the first IOA resource so that enterprise directory configuration can be provisioned and tracked declaratively.

## What Changes

- Add a new GENERAL resource `tencentcloud_ioa_company_directory_config` under the `tencentcloud/services/ioa/` package that manages the full lifecycle of an enterprise directory configuration.
- Implement the four CRUD operations backed by the ioa v20220601 SDK:
  - **Create** → `CreateCompanyDirectoryConfig`, returning the new directory `Id` (and related auth-source fields) which becomes the resource ID.
  - **Read** → `DescribeCompanyDirectoryConfig`, hydrating state from `DirectoryConfigData`.
  - **Update** → `ModifyCompanyDirectoryConfig`, re-sending the full configuration (Create-style parameters plus `Id`).
  - **Delete** → `DeleteAccountGroup` (the SDK API that deletes an account group / directory), using the directory `Id` as `AccountGroupId`.
- Add a new IOA service client accessor `UseIoaV20220601Client` in `tencentcloud/connectivity/client.go` (and the corresponding import + connection field), since no IOA client exists yet.
- Add a service layer file `service_tencentcloud_ioa.go` with a `DescribeCompanyDirectoryConfigById` helper following the existing retry/error pattern.
- Register the new resource in `tencentcloud/provider.go` and add the entry in `tencentcloud/provider.md`.
- Add a resource example markdown `resource_tc_ioa_company_directory_config.md` (later consumed by `make doc`).
- Add unit tests using gomonkey mocks (no Terraform test suite) in `resource_tc_ioa_company_directory_config_test.go`.

## Capabilities

### New Capabilities
- `ioa-company-directory-config-resource`: Manages a TencentCloud IOA enterprise directory configuration (`tencentcloud_ioa_company_directory_config`) through Create/Read/Update/Delete operations, including the ioa v20220601 SDK client accessor and service-layer helper.

### Modified Capabilities
<!-- No existing spec-level behavior is being modified; this is a net-new resource. -->
(none)

## Impact

- **New code**:
  - `tencentcloud/services/ioa/resource_tc_ioa_company_directory_config.go` — resource definition + CRUD.
  - `tencentcloud/services/ioa/service_tencentcloud_ioa.go` — IOA service layer (Describe helper).
  - `tencentcloud/services/ioa/resource_tc_ioa_company_directory_config_test.go` — gomonkey-mocked unit tests.
  - `tencentcloud/services/ioa/resource_tc_ioa_company_directory_config.md` — usage example.
- **Modified code**:
  - `tencentcloud/connectivity/client.go` — add ioa import, `ioaV20220601Conn` field and `UseIoaV20220601Client()` accessor.
  - `tencentcloud/provider.go` — register `tencentcloud_ioa_company_directory_config`.
  - `tencentcloud/provider.md` — add the resource entry (consumed by `make doc`).
- **Dependencies**: The ioa v20220601 SDK is already vendored under `vendor/github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ioa/v20220601/`; `go.mod` / `go.sum` already reference it. No new external dependencies.
- **Backward compatibility**: Fully additive. No existing resource schema is modified and no state is affected.
