# dts-sync-job-tags-set-type Specification

## Purpose
TBD - created by archiving change fix-dts-sync-job-tags-list-to-set. Update Purpose after archive.
## Requirements
### Requirement: DTS sync job tags field uses Set type

The `tags` field of the `tencentcloud_dts_sync_job` resource SHALL use `schema.TypeSet` with a nested `schema.Resource` element containing `tag_key` (string) and `tag_value` (string) sub-fields. The Set SHALL use the default `schema.HashResource` hash function so that two tag collections with identical key/value pairs — regardless of element order — are treated as equal. This eliminates spurious diffs caused by the DTS cloud API `DescribeSyncJobs` returning tags in an unordered manner.

#### Scenario: Tags are written to CreateSyncJob request from Set

- **WHEN** a user configures the `tags` block on `tencentcloud_dts_sync_job` and the provider executes the Create handler
- **THEN** the provider SHALL read the `tags` field via `v.(*schema.Set).List()` and iterate the resulting elements to build `[]*dts.TagItem` for the `CreateSyncJob` request `Tags` parameter

#### Scenario: Tags are read back into Set state

- **WHEN** the provider reads the resource via `DescribeSyncJobs` and the returned `SyncJob.Tags` is non-nil
- **THEN** the provider SHALL build a `[]interface{}` of `map[string]interface{}` (each containing `tag_key` and `tag_value`) and call `d.Set("tags", tagsList)` to populate the Set state

#### Scenario: Unordered tag response does not produce a diff

- **WHEN** the user's configuration contains tags `[{tag_key="a", tag_value="1"}, {tag_key="b", tag_value="2"}]` and the cloud API returns the same tags in reverse order `[{tag_key="b", tag_value="2"}, {tag_key="a", tag_value="1"}]`
- **THEN** `terraform plan` SHALL report no diff on the `tags` field, because the Set type is order-insensitive

#### Scenario: Tags field is omitted from configuration

- **WHEN** a user does not set the `tags` field in the `tencentcloud_dts_sync_job` configuration
- **THEN** the provider SHALL NOT set the `Tags` field on the `CreateSyncJob` request, preserving existing behavior

#### Scenario: Tags returned as nil from cloud API

- **WHEN** the provider reads the resource and the returned `SyncJob.Tags` is nil
- **THEN** the provider SHALL NOT call `d.Set("tags", ...)`, leaving existing state untouched

#### Scenario: HCL configuration syntax compatibility

- **WHEN** a user writes tags using the existing block-list HCL syntax (repeated `tags { tag_key = ... tag_value = ... }` blocks)
- **THEN** Terraform SHALL accept the configuration without error, because TypeList and TypeSet share the same HCL block syntax

