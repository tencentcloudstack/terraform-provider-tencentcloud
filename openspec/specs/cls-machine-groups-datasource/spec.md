# cls-machine-groups-datasource Specification

## Purpose
TBD - created by archiving change add-cls-machine-groups-datasource. Update Purpose after archive.
## Requirements
### Requirement: Data source to query CLS machine groups
The system SHALL provide a Terraform data source named `tencentcloud_cls_machine_groups` that queries CLS machine groups via the `DescribeMachineGroups` API (CLS SDK package `cls/v20201016`) and exposes the returned machine groups as a computed list.

#### Scenario: Query all machine groups with no filters
- **WHEN** the user declares `data "tencentcloud_cls_machine_groups" "all" {}` with no filters
- **THEN** the provider calls `DescribeMachineGroups` with no `Filters` and internal pagination (limit=100) until all pages are fetched, and sets `machine_groups` to the full list of returned groups and `total_count` to the API-reported total.

#### Scenario: Query machine groups with filters
- **WHEN** the user provides one or more `filters` blocks (each with `name` and `values`)
- **THEN** the provider converts each filter to a `cls.Filter` (Name + Values), sends them in the `DescribeMachineGroups` request, and returns only the matching machine groups.

#### Scenario: Pagination across multiple pages
- **WHEN** the account has more than 100 machine groups matching the query
- **THEN** the provider iterates `offset` by `limit` (100) per page, appending each page's `MachineGroups` to the result, and stops when a page returns fewer than `limit` items.

### Requirement: Machine group output fields
The data source SHALL expose, for each element of the `machine_groups` list, the fields returned by `MachineGroupInfo`: `group_id`, `group_name`, `machine_group_type` (with `type` and `values`), `create_time`, `tags` (list of key/value), `auto_update`, `update_start_time`, `update_end_time`, `service_logging`, `delay_cleanup_time`, `meta_tags` (list of key/value), and `os_type`. Fields that are nil in the API response SHALL be omitted from the mapped output rather than causing a panic.

#### Scenario: Field mapping from API response
- **WHEN** `DescribeMachineGroups` returns a `MachineGroupInfo` with all fields populated
- **THEN** the provider maps each non-nil field into the corresponding `machine_groups` element attribute, including nested `machine_group_type`, `tags`, and `meta_tags`.

#### Scenario: Nil-safe field handling
- **WHEN** a `MachineGroupInfo` field (e.g. `MachineGroupType`, `Tags`, `MetaTags`) is nil
- **THEN** the provider skips setting that nested attribute instead of dereferencing a nil pointer.

### Requirement: Empty-response handling for data source
The data source Read function SHALL NOT clear the Terraform state ID when the CLS API returns an empty response due to a transient error. When the API returns `nil` response, `nil` `Response`, or an empty `MachineGroups` list inside the retry block, the provider SHALL return a `NonRetryableError` so the outer retry continues. On the outer retry-exhausted path, the provider SHALL log `[DATASOURCE] read empty, skip SetId` and propagate the error without calling `d.SetId("")`.

#### Scenario: Transient empty API response
- **WHEN** `DescribeMachineGroups` returns a nil response or nil `Response` during a retry attempt
- **THEN** the provider returns `NonRetryableError` within the retry block so the retry loop continues, and does not clear the data source ID.

#### Scenario: Genuine empty result (no matching groups)
- **WHEN** the API returns a valid non-nil response with an empty `MachineGroups` list and zero `TotalCount`
- **THEN** the provider sets `machine_groups` to an empty list and sets a deterministic ID from the empty ids hash, without error.

### Requirement: Data source registration and documentation
The provider SHALL register the `tencentcloud_cls_machine_groups` data source in `tencentcloud/provider.go` and list it in `tencentcloud/provider.md`. A docs source file `data_source_tc_cls_machine_groups.md` SHALL be created (one-line description mentioning CLS, Example Usage, and Import section only where applicable) so `make doc` generates the website documentation.

#### Scenario: Registered in provider
- **WHEN** a user runs `terraform providers` / plan referencing `data.tencentcloud_cls_machine_groups`
- **THEN** the data source is recognized because it is registered in `provider.go`.

#### Scenario: Documentation generated
- **WHEN** `make doc` is run during the finalize phase
- **THEN** the website docs page for `tencentcloud_cls_machine_groups` is generated from the `data_source_tc_cls_machine_groups.md` source.

### Requirement: Result output file
The data source SHALL support an optional `result_output_file` parameter. When set to a non-empty path, the provider SHALL write the mapped `machine_groups` list as JSON to that file.

#### Scenario: Export results to file
- **WHEN** the user sets `result_output_file = "groups.json"`
- **THEN** the provider writes the mapped machine groups list to `groups.json` in JSON format.

### Requirement: Unit tests with mocked API
The data source SHALL include unit tests in `data_source_tc_cls_machine_groups_test.go` that use gomonkey to mock the CLS `DescribeMachineGroups` API (no terraform acceptance test suite), verifying field mapping, nil-safe handling, and multi-page pagination.

#### Scenario: Unit test covers field mapping and pagination
- **WHEN** the unit tests run
- **THEN** they mock `DescribeMachineGroups` to return canned `MachineGroupInfo` across two pages and assert the Read function correctly maps fields and aggregates pages.

