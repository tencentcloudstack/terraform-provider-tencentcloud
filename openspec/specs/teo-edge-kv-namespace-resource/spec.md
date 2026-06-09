## ADDED Requirements

### Requirement: Resource schema definition

The resource `tencentcloud_teo_edge_k_v_namespace` SHALL define the following schema fields:
- `zone_id` (Required, ForceNew, String): 站点 ID
- `namespace` (Required, ForceNew, String): 命名空间名称，1-50 个字符，允许 a-z、A-Z、0-9、-
- `remark` (Optional, String): 命名空间描述，最大 256 个字符

The resource SHALL use `zone_id#namespace` as the composite resource ID (using `tccommon.FILED_SP` as separator).

#### Scenario: Schema fields are correctly defined
- **WHEN** user defines a `tencentcloud_teo_edge_k_v_namespace` resource block with `zone_id`, `namespace`, and `remark`
- **THEN** Terraform SHALL accept the configuration and plan the resource creation

#### Scenario: Import with composite ID
- **WHEN** user runs `terraform import tencentcloud_teo_edge_k_v_namespace.example zone-xxx#my-namespace`
- **THEN** the resource SHALL be imported successfully by parsing `zone_id` and `namespace` from the composite ID

### Requirement: Create operation

The resource Create function SHALL call `CreateEdgeKVNamespace` API with `ZoneId`, `Namespace`, and `Remark` parameters. After successful creation, the resource ID SHALL be set to `zone_id#namespace`. The Create function SHALL verify that the API call succeeds before setting the ID.

#### Scenario: Successful creation
- **WHEN** user applies a configuration with valid `zone_id`, `namespace`, and `remark`
- **THEN** the system SHALL call `CreateEdgeKVNamespace` API and set the resource ID to `zone_id#namespace`

#### Scenario: Creation with retry on failure
- **WHEN** the `CreateEdgeKVNamespace` API call fails with a retryable error
- **THEN** the system SHALL retry the request using `tccommon.ReadRetryTimeout` and `resource.RetryContext`

### Requirement: Read operation

The resource Read function SHALL call `DescribeEdgeKVNamespaces` API with `ZoneId` and a `namespace` filter to retrieve the namespace details. It SHALL set `zone_id`, `namespace`, and `remark` from the response. If the namespace is not found, it SHALL remove the resource from state.

#### Scenario: Successful read
- **WHEN** the resource exists in the cloud
- **THEN** the system SHALL call `DescribeEdgeKVNamespaces` with namespace filter, find the matching namespace, and set `zone_id`, `namespace`, `remark` in state

#### Scenario: Resource not found
- **WHEN** the `DescribeEdgeKVNamespaces` API returns empty results for the namespace
- **THEN** the system SHALL remove the resource from Terraform state (call `d.SetId("")`)

### Requirement: Update operation

The resource Update function SHALL call `ModifyEdgeKVNamespace` API when `remark` changes. Since `zone_id` and `namespace` are ForceNew, they cannot be updated in-place.

#### Scenario: Update remark
- **WHEN** user changes the `remark` field value
- **THEN** the system SHALL call `ModifyEdgeKVNamespace` API with the new `remark` value

#### Scenario: Update with retry on failure
- **WHEN** the `ModifyEdgeKVNamespace` API call fails with a retryable error
- **THEN** the system SHALL retry the request using `tccommon.ReadRetryTimeout` and `resource.RetryContext`

### Requirement: Delete operation

The resource Delete function SHALL call `DeleteEdgeKVNamespace` API with `ZoneId` and `Namespace` parameters to remove the namespace.

#### Scenario: Successful deletion
- **WHEN** user destroys the resource
- **THEN** the system SHALL call `DeleteEdgeKVNamespace` API with the correct `ZoneId` and `Namespace`

#### Scenario: Delete with retry on failure
- **WHEN** the `DeleteEdgeKVNamespace` API call fails with a retryable error
- **THEN** the system SHALL retry the request using `tccommon.ReadRetryTimeout` and `resource.RetryContext`

### Requirement: Provider registration

The resource SHALL be registered in `tencentcloud/provider.go` under the TEO service section, and referenced in `tencentcloud/provider.md`.

#### Scenario: Resource is available in provider
- **WHEN** user configures the tencentcloud provider
- **THEN** the resource `tencentcloud_teo_edge_k_v_namespace` SHALL be available for use

### Requirement: Unit tests with gomonkey mock

The resource SHALL have unit tests in `resource_tc_teo_edge_k_v_namespace_test.go` that use gomonkey to mock cloud API calls. Tests SHALL cover Create, Read, Update, and Delete operations and SHALL pass with `go test -gcflags=all=-l`.

#### Scenario: Unit tests pass
- **WHEN** running `go test -gcflags=all=-l` on the test file
- **THEN** all unit tests SHALL pass without requiring real cloud credentials

### Requirement: Resource documentation

The resource SHALL have a documentation file at `tencentcloud/services/teo/resource_tc_teo_edge_k_v_namespace.md` with:
- A one-line description mentioning TEO product
- Example Usage section
- Import section showing the composite ID format `zone_id#namespace`

#### Scenario: Documentation is complete
- **WHEN** the documentation file is generated
- **THEN** it SHALL contain description, Example Usage, and Import sections
