## ADDED Requirements

### Requirement: Datasource schema definition
The system SHALL define a Terraform datasource schema for `tencentcloud_dbdc_db_custom_nodes` with the following input parameters:
- `node_ids`: Optional, TypeList of string - filter by one or more node IDs (max 100 per request)
- `filters`: Optional, TypeList of schema.Resource with sub-fields `name` (Required, string) and `values` (Required, TypeSet of string) - filter by cluster-id, node-name, status, zone
- `tags`: Optional, TypeList of schema.Resource with sub-fields `key` (Required, string) and `value` (Required, string) - filter by tag key-value pairs
- `result_output_file`: Optional, TypeString - used to save results to file

The system SHALL define computed output parameters:
- `node_set`: Computed, TypeList of schema.Resource containing flattened node attributes
- Each node element SHALL contain: `node_id`, `node_name`, `ssh_endpoint`, `lan_ip`, `cluster_id`, `zone`, `node_type`, `cpu`, `memory`, `system_disk` (TypeList of schema.Resource with `disk_type`, `disk_size`), `data_disks` (TypeList of schema.Resource with `disk_type`, `disk_size`, `disk_name`), `os_name`, `image_id`, `vpc_id`, `subnet_id`, `status`, `charge_type`, `expire_time`, `created_time`, `isolated_time`, `tags` (TypeList of schema.Resource with `key`, `value`), `auto_renew`, `switch_id`, `rack_id`, `host_ip`

#### Scenario: Datasource with filters input
- **WHEN** a user provides `filters` with `name` = "cluster-id" and `values` = ["cluster-123"]
- **THEN** the system SHALL call DescribeDBCustomNodes API with the corresponding Filter parameter and return matching nodes

#### Scenario: Datasource with node_ids input
- **WHEN** a user provides `node_ids` = ["node-1", "node-2"]
- **THEN** the system SHALL call DescribeDBCustomNodes API with NodeIds parameter containing those IDs

#### Scenario: Datasource with tags input
- **WHEN** a user provides `tags` with `key` = "env" and `value` = "prod"
- **THEN** the system SHALL call DescribeDBCustomNodes API with Tags parameter containing that key-value pair

#### Scenario: Datasource with no filter inputs
- **WHEN** a user provides no filter inputs (node_ids, filters, tags all empty)
- **THEN** the system SHALL call DescribeDBCustomNodes API without filters and return all nodes the account has access to

### Requirement: Read operation with retry and pagination
The system SHALL implement a Read function that calls `DescribeDBCustomNodes` API with `tccommon.ReadRetryTimeout` retry logic. The system SHALL implement internal pagination by setting `Limit` to 100 (API max) and incrementing `Offset` until all results are collected. The system SHALL NOT expose `limit` or `offset` parameters to users in the schema.

#### Scenario: Successful read with pagination
- **WHEN** the API returns more than 100 nodes total
- **THEN** the system SHALL make multiple API calls with incremented Offset values and concatenate all NodeSet results into a single `node_set` output

#### Scenario: API retry on transient failure
- **WHEN** the DescribeDBCustomNodes API call fails with a transient error
- **THEN** the system SHALL retry the call using `tccommon.ReadRetryTimeout` and `tccommon.RetryError()` for error wrapping

#### Scenario: Empty API response
- **WHEN** the DescribeDBCustomNodes API returns nil response or empty NodeSet
- **THEN** the system SHALL return `NonRetryableError` instead of clearing `d.SetId("")`, and SHALL log `log.Printf("[DATASOURCE] read empty, skip SetId")`

### Requirement: Nil field handling in response
The system SHALL check each response field for nil before calling `d.Set()` or adding it to the output map. Fields that may be nil according to API documentation (`SystemDisk`, `DataDisks`, `Tags` in DBCustomNode) SHALL be skipped when nil rather than set to empty values.

#### Scenario: Node with nil SystemDisk
- **WHEN** a DBCustomNode has nil `SystemDisk` field
- **THEN** the system SHALL skip setting `system_disk` for that node element

#### Scenario: Node with nil DataDisks
- **WHEN** a DBCustomNode has nil `DataDisks` field
- **THEN** the system SHALL skip setting `data_disks` for that node element

### Requirement: Datasource ID generation
The system SHALL use `helper.BuildToken()` as the datasource ID after successful read, following the standard pattern for list-type datasources.

#### Scenario: Successful datasource read
- **WHEN** the Read function completes successfully with data
- **THEN** the system SHALL set `d.SetId(helper.BuildToken())`

### Requirement: Provider registration
The system SHALL register the new datasource `tencentcloud_dbdc_db_custom_nodes` in `tencentcloud/provider.go` with the datasource mapping entry.

#### Scenario: Provider registration
- **WHEN** the provider is initialized
- **THEN** `tencentcloud_dbdc_db_custom_nodes` SHALL be available as a datasource in Terraform configurations

### Requirement: Documentation
The system SHALL provide documentation in `data_source_tc_dbdc_db_custom_nodes.md` with:
- One-sentence description using format "Use this data source to query ..." mentioning the dbdc product
- Example Usage section showing filter and output usage
- No Import section (datasource type does not have import)
- No Argument Reference or Attribute Reference sections (auto-generated)

#### Scenario: Documentation file creation
- **WHEN** the datasource is added
- **THEN** the .md file SHALL contain description, example usage, and follow the standard datasource documentation format

### Requirement: Unit tests with mock
The system SHALL provide unit tests in `data_source_tc_dbdc_db_custom_nodes_test.go` using gomonkey mock approach (not Terraform test suite). The tests SHALL cover the Read function logic including request construction and response parsing. Tests SHALL be runnable with `go test -gcflags=all=-l`.

#### Scenario: Mock-based unit test for Read
- **WHEN** unit tests are executed
- **THEN** the gomonkey mock SHALL replace the API client call and verify correct schema population from mock response data
