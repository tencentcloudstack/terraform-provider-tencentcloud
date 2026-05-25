# ga2-accelerate-area Specification

## Purpose
TBD - created by archiving change add-ga2-accelerate-area-resource. Update Purpose after archive.
## Requirements
### Requirement: Resource schema definition
The `tencentcloud_ga2_accelerate_area` resource SHALL define the following schema:
- `global_accelerator_id` (String, Required, ForceNew): The global accelerator instance ID.
- `accelerator_areas` (List, Required): List of accelerate area configurations for Create/Update.
  - `accelerate_region` (String, Required): The accelerate region.
  - `bandwidth` (Integer, Required): The bandwidth in Mbps.
  - `isp_type` (String, Optional, Default "BGP"): ISP type, supports "BGP", "三网", "精品".
  - `ip_version` (String, Optional, Default "IPv4"): IP version, only supports "IPv4".
- `accelerate_area_set` (List, Computed): List of accelerate area details returned by Read.
  - `accelerate_region` (String): The accelerate region.
  - `bandwidth` (Integer): The bandwidth.
  - `isp_type` (String): ISP type.
  - `ip_version` (String): IP version.
  - `accelerator_area_id` (String): The accelerate area ID.
  - `ip_address` (List of String): IP addresses.
  - `ip_address_info_set` (List): IP address info.
    - `ip_address` (String): IP address.
    - `isp_type` (String): IP type.
- `task_id` (String, Computed): The last async task ID.
- Timeouts block: Create (default 10m), Update (default 10m), Delete (default 10m).

#### Scenario: Schema validation on apply
- **WHEN** user provides a valid `global_accelerator_id` and at least one `accelerator_areas` entry with `accelerate_region` and `bandwidth`
- **THEN** Terraform plan SHALL succeed without validation errors

#### Scenario: ForceNew on global_accelerator_id change
- **WHEN** user changes `global_accelerator_id` in configuration
- **THEN** Terraform SHALL plan a destroy-and-recreate operation

### Requirement: Create accelerate areas
The resource Create function SHALL call `CreateAccelerateAreas` API with the configured `global_accelerator_id` and `accelerator_areas`, then poll `DescribeTaskResult` until the task completes, and finally call Read to populate state.

#### Scenario: Successful creation
- **WHEN** user applies a new `tencentcloud_ga2_accelerate_area` resource
- **THEN** the system SHALL call CreateAccelerateAreas API, poll DescribeTaskResult until task status indicates success, set the resource ID to `global_accelerator_id`, and read back the state

#### Scenario: Create API returns empty response
- **WHEN** CreateAccelerateAreas returns nil Response or empty TaskId
- **THEN** the system SHALL return a non-retryable error

### Requirement: Read accelerate areas
The resource Read function SHALL call `DescribeAccelerateAreas` API with pagination to fetch all accelerate areas for the given `global_accelerator_id` and populate `accelerate_area_set`.

#### Scenario: Successful read with pagination
- **WHEN** the resource Read function is called
- **THEN** the system SHALL paginate through all results of DescribeAccelerateAreas and set `accelerate_area_set` with the complete list

#### Scenario: Resource not found
- **WHEN** DescribeAccelerateAreas returns an empty AccelerateAreaSet
- **THEN** the system SHALL remove the resource from state (d.SetId(""))

### Requirement: Update accelerate areas
The resource Update function SHALL call `ModifyAccelerateAreas` API with the updated `accelerator_areas` configuration, then poll `DescribeTaskResult` until the task completes.

#### Scenario: Successful update
- **WHEN** user modifies `accelerator_areas` (e.g., changes bandwidth)
- **THEN** the system SHALL call ModifyAccelerateAreas API, poll DescribeTaskResult until task completes, and read back the updated state

### Requirement: Delete accelerate areas
The resource Delete function SHALL first call `DescribeAccelerateAreas` to get all `AcceleratorAreaId` values, then call `DeleteAccelerateAreas` with those IDs, and poll `DescribeTaskResult` until the task completes.

#### Scenario: Successful deletion
- **WHEN** user destroys the resource
- **THEN** the system SHALL read current area IDs, call DeleteAccelerateAreas with all IDs, and poll DescribeTaskResult until task completes

#### Scenario: No areas to delete
- **WHEN** DescribeAccelerateAreas returns empty list during delete
- **THEN** the system SHALL return successfully (resource already gone)

### Requirement: Import support
The resource SHALL support Terraform import using `global_accelerator_id` as the import ID.

#### Scenario: Import existing accelerate areas
- **WHEN** user runs `terraform import tencentcloud_ga2_accelerate_area.example <global_accelerator_id>`
- **THEN** the system SHALL call Read to populate the state with current accelerate area information

### Requirement: Async task polling
All async operations (Create, Modify, Delete) SHALL poll `DescribeTaskResult` with the returned `TaskId` until the task status indicates completion, respecting the configured Terraform timeout.

#### Scenario: Task completes successfully
- **WHEN** DescribeTaskResult returns Status indicating success
- **THEN** the polling loop SHALL exit and the operation SHALL proceed

#### Scenario: Task polling timeout
- **WHEN** the Terraform timeout is reached before task completion
- **THEN** the system SHALL return a timeout error

### Requirement: Provider registration
The resource SHALL be registered in `tencentcloud/provider.go` under the GA2 service section and documented in `tencentcloud/provider.md`.

#### Scenario: Resource available after registration
- **WHEN** the provider is initialized
- **THEN** `tencentcloud_ga2_accelerate_area` SHALL be available as a resource type

