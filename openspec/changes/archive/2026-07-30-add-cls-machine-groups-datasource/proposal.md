## Why

The Terraform Provider for Tencent Cloud currently has no data source to query the list of CLS (Cloud Log Service) machine groups. Users manage machine groups through the `tencentcloud_cls_machine_group` resource, but cannot look up existing machine groups (for example to reference a group's ID, name, tags, or OS type) from within Terraform. Adding a `tencentcloud_cls_machine_groups` data source lets users query and reference existing machine groups for automation, validation, and cross-resource referencing.

## What Changes

- Add a new data source `tencentcloud_cls_machine_groups` that wraps the CLS `DescribeMachineGroups` API (package `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cls/v20201016`).
- The data source exposes a `filters` input (matching the API's `Filters` field) and a computed `machine_groups` list output, plus `total_count` and `result_output_file`.
- Register the data source in `tencentcloud/provider.go` and document it in `tencentcloud/provider.md`.
- Add a website docs page (`website/docs/d/cls_machine_groups.html.markdown`) generated via `make doc`.
- Add unit tests using gomonkey mocks (no terraform acceptance suite) per project rules for new data sources.

## Capabilities

### New Capabilities
- `cls-machine-groups-datasource`: A read-only data source to query CLS machine groups via the `DescribeMachineGroups` API, with filter support and paginated retrieval.

### Modified Capabilities
<!-- No existing spec-level behavior is changing. The new data source is additive. -->
- _None_

## Impact

- **New file**: `tencentcloud/services/cls/data_source_tc_cls_machine_groups.go` — data source schema + Read function.
- **New file**: `tencentcloud/services/cls/data_source_tc_cls_machine_groups_test.go` — unit tests with gomonkey mocks.
- **New file**: `tencentcloud/services/cls/data_source_tc_cls_machine_groups.md` — docs source for `make doc`.
- **Modified file**: `tencentcloud/provider.go` — register `tencentcloud_cls_machine_groups` data source.
- **Modified file**: `tencentcloud/provider.md` — add data source entry.
- **APIs**: Uses `DescribeMachineGroups` from the CLS `2020-10-16` SDK (already vendored).
- **Backward compatibility**: Fully additive; no existing resources or schemas are modified.
