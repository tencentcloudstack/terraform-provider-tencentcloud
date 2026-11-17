## Why

The DBDC `CreateDBCustomCluster` API supports a `DeletionProtection` parameter that controls whether cluster deletion protection is enabled, but the Terraform resource `tencentcloud_dbdc_db_custom_cluster` does not expose this parameter. Users cannot manage deletion protection declaratively through Terraform, forcing them to use the console or API directly. The `ModifyDBCustomClusterAttributes` API also supports updating `DeletionProtection`, enabling full lifecycle management.

## What Changes

- Add `deletion_protection` (Optional, TypeBool) parameter to `tencentcloud_dbdc_db_custom_cluster` resource schema to control whether cluster deletion protection is enabled. Valid values: `true` (enabled, API default), `false` (disabled).
- Pass `DeletionProtection` to the `CreateDBCustomCluster` API request when the user specifies the parameter.
- Read `DeletionProtection` from the `DescribeDBCustomClusterDetail` API response to support state refresh and import.
- Update the `Update` function to detect changes to `deletion_protection` and call `ModifyDBCustomClusterAttributes` API to update the deletion protection setting.

## Capabilities

### New Capabilities
- `dbdc-db-custom-cluster-deletion-protection`: Enable the `deletion_protection` parameter on the `tencentcloud_dbdc_db_custom_cluster` resource to allow users to specify and manage cluster deletion protection through the full lifecycle (create, read, update).

### Modified Capabilities
<!-- No existing specs require modification -->

## Impact

- **Affected files:**
  - `tencentcloud/services/dbdc/resource_tc_dbdc_db_custom_cluster.go` — add `deletion_protection` schema field, wire through Create flow, add Read support, add Update support via `ModifyDBCustomClusterAttributes`
  - `tencentcloud/services/dbdc/resource_tc_dbdc_db_custom_cluster_test.go` — add unit test cases for the new parameter
  - `tencentcloud/services/dbdc/resource_tc_dbdc_db_custom_cluster.md` — update documentation with `deletion_protection` usage example
- **SDK dependency:** No SDK upgrade required. The vendored SDK `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dbdc/v20201029` already includes `DeletionProtection` in `CreateDBCustomClusterRequest`, `DescribeDBCustomClusterDetailResponseParams`, and `ModifyDBCustomClusterAttributesRequest`.
- **Backward compatibility:** fully backward compatible — the new parameter is Optional and defaults to not being set, so existing configurations continue to work unchanged.
- **API constraints:** `DeletionProtection` is available in `CreateDBCustomCluster` (create), `DescribeDBCustomClusterDetail` (read), and `ModifyDBCustomClusterAttributes` (update), enabling full CRUD support for this parameter.
