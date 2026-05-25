## ADDED Requirements

### Requirement: Resource creates accelerate areas

The `tencentcloud_ga2_accelerate_area` resource SHALL call the `CreateAccelerateAreas` API with the specified `GlobalAcceleratorId` and `AcceleratorAreas` list when created. After the async call returns a `TaskId`, the resource SHALL poll `DescribeAccelerateAreas` until the created areas appear in the response.

#### Scenario: Successful creation of accelerate areas
- **WHEN** user applies a `tencentcloud_ga2_accelerate_area` resource with valid `global_accelerator_id` and `accelerator_areas` configuration
- **THEN** the resource calls `CreateAccelerateAreas`, waits for the areas to appear in `DescribeAccelerateAreas`, and stores `global_accelerator_id` as the resource ID

#### Scenario: Creation API returns empty response
- **WHEN** `CreateAccelerateAreas` returns a nil response or nil TaskId
- **THEN** the resource SHALL return a non-retryable error

### Requirement: Resource reads accelerate areas

The resource SHALL call `DescribeAccelerateAreas` with the stored `GlobalAcceleratorId` to read the current state. It SHALL handle pagination (using Offset/Limit) to retrieve all areas. The result SHALL be stored in the `accelerate_area_set` computed attribute.

#### Scenario: Successful read with pagination
- **WHEN** the resource performs a Read operation and there are multiple pages of accelerate areas
- **THEN** the resource iterates through all pages and populates `accelerate_area_set` with all areas

#### Scenario: Read returns empty result
- **WHEN** `DescribeAccelerateAreas` returns an empty `AccelerateAreaSet`
- **THEN** the resource SHALL remove itself from state (resource has been deleted externally)

### Requirement: Resource updates accelerate areas

The resource SHALL call `ModifyAccelerateAreas` with the updated `AcceleratorAreas` list when the `accelerator_areas` attribute changes. After the async call, it SHALL poll `DescribeAccelerateAreas` to confirm the modification took effect.

#### Scenario: Successful modification of accelerate areas
- **WHEN** user changes `accelerator_areas` configuration and applies
- **THEN** the resource calls `ModifyAccelerateAreas` with the new configuration and waits for the change to be reflected in `DescribeAccelerateAreas`

### Requirement: Resource deletes accelerate areas

The resource SHALL first call `DescribeAccelerateAreas` to collect all `AcceleratorAreaId` values, then call `DeleteAccelerateAreas` with those IDs. After the async call, it SHALL poll `DescribeAccelerateAreas` until the areas are no longer present.

#### Scenario: Successful deletion of all accelerate areas
- **WHEN** user destroys the `tencentcloud_ga2_accelerate_area` resource
- **THEN** the resource reads current area IDs, calls `DeleteAccelerateAreas`, and waits until `DescribeAccelerateAreas` returns empty for that accelerator

#### Scenario: No areas exist during delete
- **WHEN** the resource attempts to delete but `DescribeAccelerateAreas` returns no areas
- **THEN** the delete operation SHALL succeed without calling `DeleteAccelerateAreas`

### Requirement: Resource supports import

The resource SHALL support import using `global_accelerator_id` as the import ID.

#### Scenario: Import existing accelerate areas
- **WHEN** user runs `terraform import tencentcloud_ga2_accelerate_area.example <global_accelerator_id>`
- **THEN** the resource reads the current state from `DescribeAccelerateAreas` and populates all computed attributes

### Requirement: Resource handles async operations with retry

All write operations (Create/Modify/Delete) SHALL use `resource.Retry` with the configured timeout to poll `DescribeAccelerateAreas` until the expected state is observed. The retry SHALL use `tccommon.ReadRetryTimeout` for API calls and respect the schema-defined Timeouts.

#### Scenario: Async operation completes within timeout
- **WHEN** an async write operation is performed and the state converges within the configured timeout
- **THEN** the operation succeeds and the resource state is updated

#### Scenario: Async operation exceeds timeout
- **WHEN** an async write operation is performed but the state does not converge within the configured timeout
- **THEN** the resource SHALL return a timeout error

### Requirement: Resource registers in provider

The resource SHALL be registered in `tencentcloud/provider.go` under the GA2 service section and documented in `tencentcloud/provider.md`.

#### Scenario: Provider includes ga2 accelerate area resource
- **WHEN** the provider is initialized
- **THEN** `tencentcloud_ga2_accelerate_area` is available as a managed resource
