## Why

The Data Lake Compute (DLC) product currently has no Terraform resource for managing meta databases (`MetaDatabase`). Users must create and delete meta databases manually in the console, which breaks the infrastructure-as-code workflow. Adding `tencentcloud_dlc_meta_database` closes this gap and lets users manage the full lifecycle (create / read / delete / import) of a DLC meta database natively through Terraform.

## What Changes

- Add a new resource `tencentcloud_dlc_meta_database` backed by the `dlc` v20210125 SDK.
- Map all `CreateMetaDatabase` request parameters to schema fields, including: `datasource_connection_name`, `database_name`, `comment`, `govern_policy` (with `rule_type`, `govern_engine`), `smart_policy` (with deeply nested `base_info`, `policy` blocks for intelligent data governance configuration).
- Implement CRD operations: `CreateMetaDatabase` (async, returns `BatchId` and `TaskIdSet`), `DescribeDatabase` (read), `DeleteMetaDatabase` (async).
- Since there is **no Update API** (`ModifyMetaDatabase` / `UpdateMetaDatabase` / `AlterMetaDatabase` do not exist in the SDK), the resource follows the CRD-only pattern: all non-ID top-level schema fields are added to `immutableArgs`; if any of them changes, the Update method returns an error so Terraform triggers a recreate.
- `CreateMetaDatabase` and `DeleteMetaDatabase` are **async interfaces** that return `BatchId` and `TaskIdSet`; after calling them, poll `DescribeDatabase` until the resource is confirmed created or deleted.
- Surface read-only computed fields from the `DescribeDatabase` response (via `DatabaseResponseInfo`): `properties`, `create_time`, `modified_time`, `location`, `user_alias`, `user_sub_uin`, `database_id`, `catalog_name`, `catalog_type`, `is_information_schema`, `batch_id`, `task_id_set`.
- Resource ID uses `datasource_connection_name` + `FILED_SP` + `database_name` as a composite ID when `datasource_connection_name` is provided, otherwise just `database_name`.
- Wire the new resource into `tencentcloud/provider.go` and `tencentcloud/provider.md` under the `dlc` namespace.
- Author resource markdown documentation `resource_tc_dlc_meta_database.md` (example HCL snippet + `terraform import` syntax with composite ID).
- Author unit test file `resource_tc_dlc_meta_database_test.go` using gomonkey mock (no terraform test suite).

## Capabilities

### New Capabilities
- `dlc-meta-database-resource`: Lifecycle management (create / read / delete / import) of a Tencent Cloud DLC meta database, including async task polling via `DescribeDatabase`, full schema parity with `CreateMetaDatabase`, and immutable-field enforcement due to absence of an Update API.

### Modified Capabilities
<!-- None: this change only introduces a new resource; it does not alter requirement-level behavior of any existing capability. -->

## Impact

- **New code**:
  - `tencentcloud/services/dlc/resource_tc_dlc_meta_database.go` (CRD + schema + build/flatten helpers, single file, mirroring `tencentcloud_igtm_strategy` style).
  - `tencentcloud/services/dlc/resource_tc_dlc_meta_database.md` (resource doc + import syntax).
  - `tencentcloud/services/dlc/resource_tc_dlc_meta_database_test.go` (unit test with gomonkey mock).
- **Modified code**:
  - `tencentcloud/provider.go`: register `tencentcloud_dlc_meta_database` in `ResourcesMap`.
  - `tencentcloud/provider.md`: add documentation entry.
- **APIs consumed**: `CreateMetaDatabase`, `DescribeDatabase`, `DeleteMetaDatabase` (already vendored in `vendor/github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dlc/v20210125/`).
- **No breaking change**: purely additive; no existing schema or state is modified.
- **No SDK upgrade required**: all required APIs are already present in the vendored SDK.
