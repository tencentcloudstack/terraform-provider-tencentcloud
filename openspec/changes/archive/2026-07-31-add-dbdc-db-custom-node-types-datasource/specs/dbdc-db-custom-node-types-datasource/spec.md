## ADDED Requirements

### Requirement: Datasource schema definition
The system SHALL define a Terraform datasource schema for `tencentcloud_dbdc_db_custom_node_types` with the following input parameters:
- `filters`: Optional, TypeList of schema.Resource with sub-fields `name` (Required, TypeString) and `values` (Required, TypeList of TypeString) - filter by region, zone, node-family, node-type
- `result_output_file`: Optional, TypeString - used to save results to file

The system SHALL define computed output parameters:
- `node_type_set`: Computed, TypeList of schema.Resource containing node type info attributes
- Each node type element SHALL contain: `zone` (TypeString), `node_type` (TypeString), `node_family` (TypeString), `cpu` (TypeInt), `memory` (TypeInt), `status` (TypeString), `system_disk_types` (TypeList of TypeString), `data_disk_types` (TypeList of TypeString)

#### Scenario: Datasource with filters input
- **WHEN** a user provides `filters` with `name` = "zone" and `values` = ["ap-guangzhou-6"]
- **THEN** the system SHALL call DescribeDBCustomNodeTypes API with the corresponding Filters parameter and return matching node types

#### Scenario: Datasource with node-family filter
- **WHEN** a user provides `filters` with `name` = "node-family" and `values` = ["DB.SA5"]
- **THEN** the system SHALL call DescribeDBCustomNodeTypes API with the corresponding Filters parameter and return matching node types

#### Scenario: Datasource with no filter inputs
- **WHEN** a user provides no filter inputs (filters empty)
- **THEN** the system SHALL call DescribeDBCustomNodeTypes API without filters and return all node types the account has access to

### Requirement: Read operation with retry
The system SHALL implement a Read function that calls `DescribeDBCustomNodeTypes` API with `tccommon.ReadRetryTimeout` retry logic. Since the `DescribeDBCustomNodeTypes` API has no pagination parameters (no Offset/Limit), the system SHALL NOT implement pagination logic and SHALL NOT expose limit/offset parameters to users in the schema.

#### Scenario: Successful read
- **WHEN** the API returns node type data
- **THEN** the system SHALL populate `node_type_set` with all returned NodeTypeSet elements

#### Scenario: API retry on transient failure
- **WHEN** the DescribeDBCustomNodeTypes API call fails with a transient error
- **THEN** the system SHALL retry the call using `tccommon.ReadRetryTimeout` and `tccommon.RetryError()` for error wrapping

#### Scenario: Empty API response
- **WHEN** the DescribeDBCustomNodeTypes API returns nil response, nil Response, or empty NodeTypeSet
- **THEN** the system SHALL return `NonRetryableError` instead of clearing `d.SetId("")`, and SHALL log `log.Printf("[DATASOURCE] read empty, skip SetId")`

### Requirement: Nil field handling in response
The system SHALL check each response field for nil before adding it to the output map. Pointer fields (`Zone`, `NodeType`, `NodeFamily`, `CPU`, `Memory`, `Status`) SHALL be skipped when nil. Slice fields (`SystemDiskTypes`, `DataDiskTypes`) SHALL be set only when non-nil.

#### Scenario: Node type with nil SystemDiskTypes
- **WHEN** a DBCustomNodeTypeInfo has nil `SystemDiskTypes` field
- **THEN** the system SHALL skip setting `system_disk_types` for that node type element

#### Scenario: Node type with nil DataDiskTypes
- **WHEN** a DBCustomNodeTypeInfo has nil `DataDiskTypes` field
- **THEN** the system SHALL skip setting `data_disk_types` for that node type element

### Requirement: Datasource ID generation
The system SHALL use `helper.BuildToken()` as the datasource ID after successful read, following the standard pattern for list-type datasources.

#### Scenario: Successful datasource read
- **WHEN** the Read function completes successfully with data
- **THEN** the system SHALL set `d.SetId(helper.BuildToken())`

### Requirement: Provider registration
The system SHALL register the new datasource `tencentcloud_dbdc_db_custom_node_types` in `tencentcloud/provider.go` DataSourcesMap with the datasource mapping entry.

#### Scenario: Provider registration
- **WHEN** the provider is initialized
- **THEN** `tencentcloud_dbdc_db_custom_node_types` SHALL be available as a datasource in Terraform configurations

### Requirement: Documentation
The system SHALL provide documentation in `data_source_tc_dbdc_db_custom_node_types.md` with:
- One-sentence description using format "Use this data source to query ..." mentioning the dbdc product
- Example Usage section showing filter and output usage
- No Import section (datasource type does not have import)
- No Argument Reference or Attribute Reference sections (auto-generated)

#### Scenario: Documentation file creation
- **WHEN** the datasource is added
- **THEN** the .md file SHALL contain description, example usage, and follow the standard datasource documentation format

### Requirement: Unit tests with mock
The system SHALL provide unit tests in `data_source_tc_dbdc_db_custom_node_types_test.go` using gomonkey mock approach (not Terraform test suite). The tests SHALL cover the Read function logic including request construction (filters) and response parsing (NodeTypeSet fields). Tests SHALL be runnable with `go test -gcflags=all=-l`.

#### Scenario: Mock-based unit test for Read
- **WHEN** unit tests are executed
- **THEN** the gomonkey mock SHALL replace the API client call and verify correct schema population from mock response data
