## ADDED Requirements

### Requirement: Resource schema definition

The `tencentcloud_ga2_accelerate_area` resource SHALL define the following schema:

- `global_accelerator_id` (Required, ForceNew, String): The global accelerator instance ID.
- `accelerator_areas` (Required, List of Object): The accelerate area configurations. Each object contains:
  - `accelerate_region` (Required, String): The accelerate region.
  - `bandwidth` (Required, Int): The bandwidth in Mbps.
  - `isp_type` (Optional, String): ISP type, supports "BGP", "三网", "精品". Defaults to "BGP".
  - `ip_version` (Optional, String): IP version, only supports "IPv4". Defaults to "IPv4".
  - `accelerator_area_id` (Optional, Computed, String): The accelerator area ID, assigned by the backend.
  - `ip_address` (Computed, List of String): The assigned IP addresses.
  - `ip_address_info_set` (Computed, List of Object): Detailed IP address information.
    - `ip_address` (Computed, String): The IP address.
    - `isp_type` (Computed, String): The ISP type of the IP.
- `accelerate_area_set` (Computed, List of Object): The full accelerate area set returned by Read, with same structure as `accelerator_areas`.
- `task_id` (Computed, String): The async task ID from the last write operation.

The resource SHALL declare Timeouts for Create, Update, and Delete (default 10 minutes each).

#### Scenario: Schema validates required fields
- **WHEN** a user provides a configuration without `global_accelerator_id` or `accelerator_areas`
- **THEN** Terraform SHALL report a validation error indicating the missing required field

#### Scenario: ForceNew on global_accelerator_id change
- **WHEN** a user changes the `global_accelerator_id` value in configuration
- **THEN** Terraform SHALL plan to destroy and recreate the resource

### Requirement: Create accelerate areas

The resource Create function SHALL call the `CreateAccelerateAreas` API with `GlobalAcceleratorId` and `AcceleratorAreas` from the user configuration. After a successful API call, it SHALL set the resource ID to `global_accelerator_id`. Since the API is asynchronous, the Create function SHALL poll `DescribeAccelerateAreas` until the created areas appear in the response, using `resource.Retry` with the configured Create timeout.

#### Scenario: Successful creation
- **WHEN** a user applies a valid configuration with `global_accelerator_id` and `accelerator_areas`
- **THEN** the resource SHALL call `CreateAccelerateAreas`, poll until areas are visible via `DescribeAccelerateAreas`, and store the `global_accelerator_id` as the resource ID

#### Scenario: Create API returns empty response
- **WHEN** the `CreateAccelerateAreas` API returns a nil response or nil TaskId
- **THEN** the resource SHALL return a non-retryable error

### Requirement: Read accelerate areas

The resource Read function SHALL call `DescribeAccelerateAreas` with the `GlobalAcceleratorId` from the resource ID. It SHALL handle pagination (using Offset/Limit) to retrieve all areas. The response `AccelerateAreaSet` SHALL be flattened and set into the state as `accelerate_area_set`.

#### Scenario: Successful read
- **WHEN** Terraform refreshes the resource state
- **THEN** the resource SHALL call `DescribeAccelerateAreas` and populate `accelerate_area_set` with all returned areas

#### Scenario: Resource not found
- **WHEN** `DescribeAccelerateAreas` returns an empty `AccelerateAreaSet`
- **THEN** the resource SHALL remove the resource from state (set ID to empty)

### Requirement: Update accelerate areas

The resource Update function SHALL call `ModifyAccelerateAreas` with the updated `AcceleratorAreas` list (including `AcceleratorAreaId` for existing areas). Since the API is asynchronous, the Update function SHALL poll `DescribeAccelerateAreas` until the modifications are reflected, using `resource.Retry` with the configured Update timeout.

#### Scenario: Successful update
- **WHEN** a user modifies `accelerator_areas` (e.g., changes bandwidth) and applies
- **THEN** the resource SHALL call `ModifyAccelerateAreas` with the full updated area list and poll until changes are visible

### Requirement: Delete accelerate areas

The resource Delete function SHALL first call `DescribeAccelerateAreas` to collect all `AcceleratorAreaId` values, then call `DeleteAccelerateAreas` with those IDs. Since the API is asynchronous, the Delete function SHALL poll `DescribeAccelerateAreas` until the areas are no longer present, using `resource.Retry` with the configured Delete timeout.

#### Scenario: Successful deletion
- **WHEN** a user destroys the resource
- **THEN** the resource SHALL collect all area IDs via `DescribeAccelerateAreas`, call `DeleteAccelerateAreas`, and poll until no areas remain

#### Scenario: No areas to delete
- **WHEN** `DescribeAccelerateAreas` returns an empty set during Delete
- **THEN** the resource SHALL succeed without calling `DeleteAccelerateAreas`

### Requirement: Import support

The resource SHALL support import using the `global_accelerator_id` as the import ID.

#### Scenario: Import existing resource
- **WHEN** a user runs `terraform import tencentcloud_ga2_accelerate_area.example <global_accelerator_id>`
- **THEN** Terraform SHALL call Read and populate the state with the current accelerate areas

### Requirement: Retry and error handling

All API calls SHALL be wrapped with `resource.Retry` using `tccommon.ReadRetryTimeout` for read operations and the configured Timeouts for write operations. Errors SHALL be wrapped with `tccommon.RetryError()`.

#### Scenario: Transient API error during read
- **WHEN** `DescribeAccelerateAreas` returns a transient error
- **THEN** the resource SHALL retry the call within the timeout period

### Requirement: Provider registration

The resource SHALL be registered in `tencentcloud/provider.go` under the GA2 service section and documented in `tencentcloud/provider.md`.

#### Scenario: Resource available after registration
- **WHEN** the provider is initialized
- **THEN** `tencentcloud_ga2_accelerate_area` SHALL be available as a resource type
