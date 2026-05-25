# ga2-accelerate-area-crud Specification

## Purpose
TBD - created by archiving change add-ga2-accelerate-area. Update Purpose after archive.
## Requirements
### Requirement: Create accelerate areas
The system SHALL allow users to create accelerate areas for a GA2 global accelerator instance by specifying the global_accelerator_id and a list of accelerator_areas configurations.

The resource SHALL call the CreateAccelerateAreas API with GlobalAcceleratorId and AcceleratorAreas parameters. After the async API returns a TaskId, the resource SHALL poll DescribeAccelerateAreas until the accelerate areas are created and visible.

#### Scenario: Successful creation of accelerate areas
- **WHEN** user applies a terraform config with tencentcloud_ga2_accelerate_area resource specifying global_accelerator_id and accelerator_areas
- **THEN** the resource calls CreateAccelerateAreas API, polls DescribeAccelerateAreas until areas appear, sets the resource ID to global_accelerator_id, and stores the returned accelerate_area_set in state

#### Scenario: Creation API returns empty response
- **WHEN** CreateAccelerateAreas API returns nil response or nil TaskId
- **THEN** the resource SHALL return a NonRetryableError

### Requirement: Read accelerate areas
The system SHALL allow reading the current state of accelerate areas for a given global accelerator instance by calling DescribeAccelerateAreas API with pagination support.

#### Scenario: Successful read with pagination
- **WHEN** the resource read function is called
- **THEN** the resource SHALL call DescribeAccelerateAreas with the GlobalAcceleratorId, handle pagination (using maximum Limit value), and set accelerate_area_set in state

#### Scenario: Resource not found on read
- **WHEN** DescribeAccelerateAreas returns empty AccelerateAreaSet
- **THEN** the resource SHALL remove the resource from state (d.SetId(""))

### Requirement: Update accelerate areas
The system SHALL allow users to modify accelerate areas configuration by calling ModifyAccelerateAreas API. The global_accelerator_id is ForceNew and cannot be changed.

#### Scenario: Successful modification of accelerate areas
- **WHEN** user modifies accelerator_areas in terraform config (e.g., changes bandwidth or adds/removes areas)
- **THEN** the resource calls ModifyAccelerateAreas API with updated AcceleratorAreas, polls DescribeAccelerateAreas until changes are reflected

### Requirement: Delete accelerate areas
The system SHALL allow users to delete all accelerate areas for a global accelerator instance. The delete operation SHALL first query current AcceleratorAreaIds via DescribeAccelerateAreas, then call DeleteAccelerateAreas with those IDs.

#### Scenario: Successful deletion of accelerate areas
- **WHEN** user destroys the tencentcloud_ga2_accelerate_area resource
- **THEN** the resource first calls DescribeAccelerateAreas to get all AcceleratorAreaIds, then calls DeleteAccelerateAreas with those IDs, and polls until areas are removed

#### Scenario: No areas to delete
- **WHEN** DescribeAccelerateAreas returns empty AccelerateAreaSet during delete
- **THEN** the resource SHALL return successfully (nothing to delete)

### Requirement: Import support
The system SHALL support importing existing accelerate areas by global_accelerator_id.

#### Scenario: Import by global_accelerator_id
- **WHEN** user runs terraform import with a global_accelerator_id
- **THEN** the resource SHALL read the current accelerate areas state and populate terraform state

### Requirement: Retry and error handling
The system SHALL implement retry logic with tccommon.ReadRetryTimeout for all API calls, wrapping errors with tccommon.RetryError().

#### Scenario: Transient API failure during create
- **WHEN** CreateAccelerateAreas API returns a transient error
- **THEN** the resource SHALL retry the operation within the configured timeout

#### Scenario: Instance state not allowed
- **WHEN** API returns UnsupportedOperation.InstanceStateNotAllowedOperate
- **THEN** the resource SHALL retry the operation (instance may be processing another async operation)

