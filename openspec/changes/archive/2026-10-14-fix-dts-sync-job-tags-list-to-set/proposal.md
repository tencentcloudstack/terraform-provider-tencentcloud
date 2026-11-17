## Why

The `tags` field of `tencentcloud_dts_sync_job` is currently defined as `schema.TypeList`, but the DTS cloud API `DescribeSyncJobs` returns tags in an unordered manner. This causes Terraform to detect a perpetual diff on every `terraform plan`/`refresh` — the tag order in state never matches the order returned by the API, so users see non-empty diffs even when the tag set is unchanged. Changing the field to `schema.TypeSet` eliminates order-sensitive diffs and matches the unordered semantics of the cloud API.

## What Changes

- **BREAKING** (state format only): Change the `tags` field in `tencentcloud_dts_sync_job` from `schema.TypeList` to `schema.TypeSet`. The nested `schema.Resource` element (with `tag_key` and `tag_value` fields) remains unchanged. Using `TypeSet` with the default `schema.HashResource` hashing, Terraform treats two tag sets with the same key/value pairs — regardless of order — as equal.
- Update the Create handler to read the `tags` field via `v.(*schema.Set).List()` instead of `v.([]interface{})`.
- Update the Read handler so that `d.Set("tags", tagsList)` correctly populates the Set (the existing `[]interface{}` of maps already works with TypeSet).
- Update the test file `resource_tc_dts_sync_job_test.go` to use `resource.TestCheckTypeSetElemNestedAttrs` for tag verification and update the HCL config to use the set syntax.
- Update the `.md` example doc to keep it consistent.

## Capabilities

### New Capabilities
- `dts-sync-job-tags-set-type`: Defines that the `tags` field of `tencentcloud_dts_sync_job` SHALL use `schema.TypeSet` so that tag ordering differences between configuration and cloud API responses do not produce spurious diffs.

### Modified Capabilities
<!-- (no existing capability spec to modify) -->

## Impact

- **Code**: `tencentcloud/services/dts/resource_tc_dts_sync_job.go` — schema definition (line 71–91), Create handler (line 168–180), Read handler (line 283–297).
- **Tests**: `tencentcloud/services/dts/resource_tc_dts_sync_job_test.go` — update tag assertions and HCL config.
- **Docs**: `tencentcloud/services/dts/resource_tc_dts_sync_job.md` — example syntax (tags block unchanged for users but doc regenerated via `make doc`).
- **API**: No cloud API change. `CreateSyncJob` request and `DescribeSyncJobs` response schemas are unchanged — only the Terraform-side collection type changes from List to Set.
- **State migration**: On the first `terraform plan` after provider upgrade, users will see a one-time state refresh converting the `tags` representation from List to Set. No actual cloud resource modification occurs. The HCL configuration syntax (block-list style) is identical for List and Set, so user `.tf` files require no changes.