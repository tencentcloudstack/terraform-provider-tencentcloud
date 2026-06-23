## ADDED Requirements

### Requirement: Create Edge KV Namespace
The system SHALL allow users to create a TEO Edge KV namespace by specifying `zone_id`, `namespace`, and optionally `remark`. The resource SHALL call the `CreateEdgeKVNamespace` API and use `zone_id#namespace` (separated by `tccommon.FILED_SP`) as the composite resource ID. After creation, the system SHALL verify the API response is not nil and return a `NonRetryableError` if the response is empty.

#### Scenario: Successful creation with all parameters
- **WHEN** user provides `zone_id`, `namespace`, and `remark` in the Terraform configuration
- **THEN** the system calls `CreateEdgeKVNamespace` API with the provided parameters, sets the resource ID to `zone_id#namespace`, and reads back the resource state

#### Scenario: Successful creation without optional remark
- **WHEN** user provides only `zone_id` and `namespace` without `remark`
- **THEN** the system calls `CreateEdgeKVNamespace` API with `zone_id` and `namespace`, sets the resource ID, and reads back the resource state

#### Scenario: API returns empty response on create
- **WHEN** the `CreateEdgeKVNamespace` API returns a nil or empty response
- **THEN** the system logs the current `logId` and `d.Id()`, and returns a `NonRetryableError`

### Requirement: Read Edge KV Namespace
The system SHALL read the Edge KV namespace state by calling `DescribeEdgeKVNamespaces` API with `ZoneId` and a `Filters` parameter filtering by `namespace` name. The system SHALL set `Limit` to 1000 (the API maximum). The system SHALL set computed fields (`capacity`, `capacity_used`, `created_on`, `modified_on`) from the API response only when the corresponding response fields are not nil.

#### Scenario: Successful read
- **WHEN** the resource exists and `DescribeEdgeKVNamespaces` returns a matching namespace
- **THEN** the system sets `zone_id`, `namespace`, `remark`, `capacity`, `capacity_used`, `created_on`, and `modified_on` from the response (checking each field for nil before setting)

#### Scenario: Resource not found
- **WHEN** `DescribeEdgeKVNamespaces` returns an empty `KVNamespaces` list or no matching namespace
- **THEN** the system logs `[CRUD] teo_edge_kv_namespace id=<current_id>` before calling `d.SetId("")` to remove the resource from state

### Requirement: Update Edge KV Namespace
The system SHALL allow updating the `remark` field of an existing Edge KV namespace by calling the `ModifyEdgeKVNamespace` API. The `zone_id` and `namespace` fields SHALL be ForceNew, meaning changes to these fields trigger resource recreation.

#### Scenario: Update remark
- **WHEN** user changes the `remark` field in the Terraform configuration
- **THEN** the system calls `ModifyEdgeKVNamespace` API with `zone_id`, `namespace`, and the new `remark` value, then reads back the updated state

#### Scenario: Change zone_id or namespace triggers recreation
- **WHEN** user changes `zone_id` or `namespace` in the Terraform configuration
- **THEN** Terraform destroys the existing resource and creates a new one (ForceNew behavior)

### Requirement: Delete Edge KV Namespace
The system SHALL delete an Edge KV namespace by calling the `DeleteEdgeKVNamespace` API with `zone_id` and `namespace` extracted from the composite resource ID.

#### Scenario: Successful deletion
- **WHEN** user runs `terraform destroy` or removes the resource from configuration
- **THEN** the system parses `zone_id` and `namespace` from the resource ID, calls `DeleteEdgeKVNamespace` API, and removes the resource from state

### Requirement: Import Edge KV Namespace
The system SHALL support importing an existing Edge KV namespace using the composite ID format `zone_id#namespace`.

#### Scenario: Successful import
- **WHEN** user runs `terraform import tencentcloud_teo_edge_kv_namespace.example zone_id#namespace`
- **THEN** the system parses the composite ID, calls `DescribeEdgeKVNamespaces` to read the resource state, and populates all fields

### Requirement: Provider Registration
The system SHALL register `tencentcloud_teo_edge_kv_namespace` in `tencentcloud/provider.go` resource map and add the corresponding entry in `tencentcloud/provider.md`.

#### Scenario: Resource available after registration
- **WHEN** the provider is initialized
- **THEN** `tencentcloud_teo_edge_kv_namespace` is available as a valid resource type

### Requirement: Unit Tests with Gomonkey Mock
The system SHALL provide unit tests in `resource_tc_teo_edge_kv_namespace_test.go` using gomonkey to mock cloud API calls. Tests SHALL cover Create, Read, Update, and Delete operations and SHALL pass with `go test -gcflags=all=-l`.

#### Scenario: Test CRUD operations
- **WHEN** unit tests are executed with `go test -gcflags=all=-l`
- **THEN** all tests pass, verifying the resource's Create, Read, Update, and Delete logic through mocked API responses
