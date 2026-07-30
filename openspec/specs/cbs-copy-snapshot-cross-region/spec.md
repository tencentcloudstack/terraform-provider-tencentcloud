## ADDED Requirements

### Requirement: Resource schema definition for cbs_copy_snapshot_cross_region
The resource `tencentcloud_cbs_copy_snapshot_cross_region` SHALL define a schema with the following fields:
- `snapshot_id` (TypeString, Required, ForceNew): Source snapshot ID for cross-region copy
- `destination_regions` (TypeList with TypeString elements, Required, ForceNew): Target region names for cross-region copy
- `snapshot_name` (TypeString, Optional, ForceNew): Name of the copied snapshot, defaults to "Copied <source-snapshot-id> from <region-name>"
- `delete_bind_images` (TypeBool, Optional, ForceNew): Whether to force-delete images associated with the snapshots when deleting, defaults to false
- `snapshot_copy_result_set` (TypeList, Computed): Cross-region copy results, containing snapshot_id, code, message, destination_region for each target region
- `id` (Computed): Composite resource ID in the format `snapshot_id#copied_snapshot_id` using FILED_SP separator

#### Scenario: Required fields are present
- **WHEN** a user creates a `tencentcloud_cbs_copy_snapshot_cross_region` resource
- **THEN** the resource SHALL require `snapshot_id` and `destination_regions` fields, and all Required/Optional fields SHALL have ForceNew set to true

#### Scenario: Computed fields populated from API response
- **WHEN** a user reads a `tencentcloud_cbs_copy_snapshot_cross_region` resource
- **THEN** `snapshot_copy_result_set` SHALL be populated from the CopySnapshotCrossRegions API response, containing each destination region's copy result

### Requirement: Create operation for cbs_copy_snapshot_cross_region
The Create operation SHALL call the `CopySnapshotCrossRegions` API with snapshot_id, destination_regions, and optional snapshot_name parameters. Since this is an async API, the Create function SHALL poll `DescribeSnapshots` until all copied snapshots reach NORMAL state.

#### Scenario: Successful cross-region copy creation
- **WHEN** a user creates a `tencentcloud_cbs_copy_snapshot_cross_region` resource with snapshot_id and destination_regions
- **THEN** the system SHALL call CopySnapshotCrossRegions API, check the response for SnapshotCopyResultSet, verify that Response and SnapshotCopyResultSet are not nil, set the composite ID using `snapshot_id#copied_snapshot_id` (using the first copied snapshot ID from the result set), and poll DescribeSnapshots until each copied snapshot's state is NORMAL

#### Scenario: Async polling for copy completion
- **WHEN** CopySnapshotCrossRegions returns successfully with SnapshotCopyResult entries
- **THEN** the system SHALL for each SnapshotCopyResult poll DescribeSnapshots by the copied snapshot_id and check SnapshotState until it becomes NORMAL, using resource.Retry with tccommon.ReadRetryTimeout

#### Scenario: Copy creation returns empty response
- **WHEN** CopySnapshotCrossRegions returns nil Response or nil SnapshotCopyResultSet or empty SnapshotCopyResultSet
- **THEN** the system SHALL return NonRetryableError and log the logId and current d.Id() for troubleshooting

#### Scenario: SnapshotCopyResult contains error code
- **WHEN** a SnapshotCopyResult entry has Code that is not "Success"
- **THEN** the system SHALL return an error with the Code and Message from the SnapshotCopyResult

### Requirement: Read operation for cbs_copy_snapshot_cross_region
The Read operation SHALL parse the composite ID to extract snapshot_id and copied_snapshot_id, then call `DescribeSnapshots` to query the copied snapshot details and update the state.

#### Scenario: Successful read of copied snapshot
- **WHEN** a user reads a `tencentcloud_cbs_copy_snapshot_cross_region` resource
- **THEN** the system SHALL parse d.Id() into snapshot_id and copied_snapshot_id using FILED_SP, call DescribeSnapshots with the copied_snapshot_id, and populate all schema fields from the response

#### Scenario: Copied snapshot not found during read
- **WHEN** DescribeSnapshots returns empty result for the copied snapshot
- **THEN** the system SHALL first log `log.Printf("[CRUD] cbs_copy_snapshot_cross_region id=%s", d.Id())`, then set d.SetId("") to signal the resource no longer exists

#### Scenario: DescribeSnapshots API failure
- **WHEN** DescribeSnapshots API call fails
- **THEN** the system SHALL wrap the error with tccommon.RetryError and return it for retry

### Requirement: Delete operation for cbs_copy_snapshot_cross_region
The Delete operation SHALL call `DeleteSnapshots` API to delete all copied snapshots identified by the snapshot_copy_result_set stored in the state.

#### Scenario: Successful deletion of copied snapshots
- **WHEN** a user deletes a `tencentcloud_cbs_copy_snapshot_cross_region` resource
- **THEN** the system SHALL read all copied snapshot IDs from the state's snapshot_copy_result_set, call DeleteSnapshots with those snapshot IDs and the delete_bind_images flag, and the resource SHALL be removed from state

#### Scenario: Delete with delete_bind_images enabled
- **WHEN** a user has set delete_bind_images to true
- **THEN** the DeleteSnapshots request SHALL include DeleteBindImages set to true

#### Scenario: Delete with delete_bind_images disabled
- **WHEN** a user has not set delete_bind_images or set it to false
- **THEN** the DeleteSnapshots request SHALL include DeleteBindImages set to false

### Requirement: No Update operation for cbs_copy_snapshot_cross_region
The resource SHALL NOT define an Update function. It SHALL only support Create, Read, and Delete operations (CRD pattern). All other top-level fields besides Id() SHALL be added to the immutableArgs array. If a user attempts to update an immutable field, the resource SHALL return an error.

#### Scenario: No update function defined
- **WHEN** the resource schema is defined
- **THEN** the resource SHALL only have Create, Read, and Delete functions, and no Update function SHALL be registered

#### Scenario: Immutable field change attempted
- **WHEN** a user changes any top-level field (other than Id()) in the terraform configuration
- **THEN** Terraform SHALL detect the ForceNew flag and trigger a destroy-and-recreate cycle

### Requirement: Import support for cbs_copy_snapshot_cross_region
The resource SHALL support import via `terraform import` using the composite ID format `snapshot_id#copied_snapshot_id`.

#### Scenario: Import using composite ID
- **WHEN** a user runs `terraform import tencentcloud_cbs_copy_snapshot_cross_region.example snap-xxx#snap-yyy`
- **THEN** the system SHALL parse the composite ID into snapshot_id and copied_snapshot_id, call Read to populate the state

### Requirement: Resource registration in provider.go
The resource SHALL be registered in `provider.go` ResourcesMap with key `"tencentcloud_cbs_copy_snapshot_cross_region"` and value `cbs.ResourceTencentCloudCbsCopySnapshotCrossRegionAttachment()`. The provider.go file comment SHALL be updated with the new resource name under the CBS Resource section. The `provider.md` SHALL be updated with the resource name under the Cloud Block Storage(CBS) Resource section.

#### Scenario: Provider registration
- **WHEN** the provider is initialized
- **THEN** `tencentcloud_cbs_copy_snapshot_cross_region` SHALL be available as a registered resource in the provider

### Requirement: Documentation file for cbs_copy_snapshot_cross_region
A documentation file `resource_tc_cbs_copy_snapshot_cross_region_attachment.md` SHALL be created in `tencentcloud/services/cbs/` directory, containing a one-line description, Example Usage section, and Import section (since this is RESOURCE_KIND_ATTACHMENT type which supports import).

#### Scenario: Documentation file exists
- **WHEN** the documentation generation is complete
- **THEN** a markdown file SHALL exist at `tencentcloud/services/cbs/resource_tc_cbs_copy_snapshot_cross_region_attachment.md` with description "Provides a CBS copy snapshot cross region resource." and example HCL usage showing how to create the resource, and an Import section explaining the composite ID format

### Requirement: Unit tests for cbs_copy_snapshot_cross_region
Unit test files SHALL be created using mock (gomonkey) approach for the cloud API, testing business logic only. Tests SHALL be runnable with `go test -gcflags=all=-l` and SHALL NOT use terraform test suite.

#### Scenario: Create function unit test
- **WHEN** unit tests are executed
- **THEN** there SHALL be a test case that mocks CopySnapshotCrossRegions API and verifies the Create function correctly calls the API, handles the async response, and sets the composite ID

#### Scenario: Read function unit test
- **WHEN** unit tests are executed
- **THEN** there SHALL be a test case that mocks DescribeSnapshots API and verifies the Read function correctly parses the composite ID, calls the API, and populates the state

#### Scenario: Delete function unit test
- **WHEN** unit tests are executed
- **THEN** there SHALL be a test case that mocks DeleteSnapshots API and verifies the Delete function correctly extracts copied snapshot IDs from state and calls the delete API

### Requirement: Retry and error handling patterns
All API calls SHALL use retry logic with `resource.Retry` and `tccommon.RetryTimeout`. All CRUD functions SHALL use `defer tccommon.LogElapsed()` and Read/Delete SHALL use `defer tccommon.InconsistentCheck()`. Error handling SHALL follow the established patterns with `tccommon.RetryError()`.

#### Scenario: Retry on API failure
- **WHEN** an API call fails with a retryable error
- **THEN** the system SHALL retry the call using resource.Retry with appropriate timeout (WriteRetryTimeout for Create/Delete, ReadRetryTimeout for Read)

#### Scenario: Log elapsed time
- **WHEN** any CRUD function is executed
- **THEN** the function SHALL defer `tccommon.LogElapsed("resource.tencentcloud_cbs_copy_snapshot_cross_region.<operation>")()` at the beginning

#### Scenario: Inconsistent check in Read
- **WHEN** the Read function is executed
- **THEN** the function SHALL defer `tccommon.InconsistentCheck(d, meta)()` at the beginning
