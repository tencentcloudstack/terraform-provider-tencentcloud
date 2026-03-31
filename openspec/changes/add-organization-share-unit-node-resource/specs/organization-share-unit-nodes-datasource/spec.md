## ADDED Requirements

### Requirement: Data source SHALL support querying share unit nodes
The `tencentcloud_organization_org_share_unit_nodes` data source SHALL support querying share unit nodes by calling the DescribeShareUnitNodes API.

#### Scenario: Successfully query share unit nodes with unit_id
- **WHEN** user provides valid unit_id
- **THEN** data source MUST call DescribeShareUnitNodes API with unit_id
- **THEN** data source MUST populate items list with all returned nodes
- **THEN** data source MUST populate total count

#### Scenario: Query with pagination parameters
- **WHEN** user provides offset and limit parameters
- **THEN** data source MUST pass offset and limit to API call
- **THEN** data source MUST return nodes matching pagination range
- **THEN** data source MUST respect limit range of 1-50

#### Scenario: Query with search_key filter
- **WHEN** user provides search_key parameter
- **THEN** data source MUST pass search_key to API call
- **THEN** data source MUST return only nodes matching search criteria
- **THEN** search MUST support node ID matching

#### Scenario: Query returns empty result
- **WHEN** no nodes exist for given criteria
- **THEN** data source MUST set items to empty list
- **THEN** data source MUST set total to 0
- **THEN** data source MUST NOT return an error

### Requirement: Data source SHALL implement pagination automatically
The data source SHALL handle API pagination when result count exceeds limit.

#### Scenario: Fetch all nodes across multiple pages
- **WHEN** share unit has more than 50 nodes
- **THEN** data source MUST make multiple API calls with increasing offset
- **THEN** data source MUST aggregate results from all pages
- **THEN** data source MUST continue until fewer than limit nodes returned

#### Scenario: Stop pagination when no more results
- **WHEN** API returns fewer nodes than limit
- **THEN** data source MUST stop making additional API calls
- **THEN** data source MUST return all accumulated nodes

### Requirement: Data source schema SHALL match API response structure
The data source output SHALL accurately reflect DescribeShareUnitNodes API response.

#### Scenario: Items contain share_node_id field
- **WHEN** data source populates items
- **THEN** each item MUST include share_node_id field of type Integer
- **THEN** share_node_id MUST match ShareNodeId from API response

#### Scenario: Items contain create_time field
- **WHEN** data source populates items
- **THEN** each item MUST include create_time field of type String
- **THEN** create_time MUST match CreateTime from API response
- **THEN** create_time format MUST be "YYYY-MM-DD HH:mm:ss"

### Requirement: Data source SHALL validate input parameters
The data source SHALL validate user-provided parameters before making API calls.

#### Scenario: Require unit_id parameter
- **WHEN** user does not provide unit_id
- **THEN** Terraform MUST return validation error
- **THEN** error MUST indicate unit_id is required

#### Scenario: Validate limit range
- **WHEN** user provides limit outside 1-50 range
- **THEN** data source MUST use default value 10
- **THEN** data source MUST NOT return validation error

#### Scenario: Validate offset is non-negative
- **WHEN** user provides negative offset
- **THEN** data source MUST use default value 0
- **THEN** data source MUST NOT return validation error

### Requirement: Data source SHALL handle API errors gracefully
The data source SHALL handle various API error conditions appropriately.

#### Scenario: Handle share unit not exist error
- **WHEN** specified unit_id does not exist
- **THEN** data source MUST return error with FailedOperation.ShareUnitNotExist
- **THEN** error message MUST indicate share unit not found

#### Scenario: Handle internal API errors
- **WHEN** API returns internal error
- **THEN** data source MUST return error to user
- **THEN** error message MUST include original API error details

#### Scenario: Retry on transient failures
- **WHEN** API call fails with transient error
- **THEN** data source MUST use ratelimit.Check() before retry
- **THEN** data source MUST retry with exponential backoff

### Requirement: Data source SHALL support result output to file
The data source SHALL support writing results to a file for external processing.

#### Scenario: Write results when result_output_file specified
- **WHEN** user provides result_output_file parameter
- **THEN** data source MUST serialize items to JSON format
- **THEN** data source MUST write JSON to specified file path
- **THEN** file MUST contain properly formatted JSON array

#### Scenario: Skip file output when parameter not specified
- **WHEN** user does not provide result_output_file parameter
- **THEN** data source MUST NOT create any output file
- **THEN** data source MUST still populate Terraform state normally

### Requirement: Data source SHALL provide comprehensive logging
The data source SHALL log all API interactions for debugging purposes.

#### Scenario: Log successful API calls
- **WHEN** API call succeeds
- **THEN** data source MUST log request body using request.ToJsonString()
- **THEN** data source MUST log response body using response.ToJsonString()
- **THEN** logs MUST include DEBUG level prefix
- **THEN** logs MUST include logId for correlation

#### Scenario: Log failed API calls
- **WHEN** API call fails
- **THEN** data source MUST log error with CRITAL level
- **THEN** log MUST include request body and error reason
- **THEN** log MUST include logId for correlation

### Requirement: Data source SHALL have complete documentation
The data source SHALL provide comprehensive documentation for users.

#### Scenario: Documentation includes usage example
- **WHEN** user views data source documentation
- **THEN** documentation MUST include complete Terraform configuration example
- **THEN** example MUST show required attributes (unit_id)
- **THEN** example MUST demonstrate optional parameters (offset, limit, search_key)
- **THEN** example MUST show how to access output items

#### Scenario: Documentation describes all arguments
- **WHEN** user views argument reference
- **THEN** documentation MUST list unit_id as Required with description "共享单元ID"
- **THEN** documentation MUST list offset as Optional with default 0
- **THEN** documentation MUST list limit as Optional with default 10
- **THEN** documentation MUST list search_key as Optional
- **THEN** documentation MUST list result_output_file as Optional

#### Scenario: Documentation describes all attributes
- **WHEN** user views attribute reference
- **THEN** documentation MUST describe items as list of share unit nodes
- **THEN** documentation MUST describe items.share_node_id as node ID
- **THEN** documentation MUST describe items.create_time as creation timestamp
- **THEN** documentation MUST describe total as total count

#### Scenario: Documentation includes generated website version
- **WHEN** make doc is executed
- **THEN** system MUST generate website/docs/d/organization_org_share_unit_nodes.html.markdown
- **THEN** generated file MUST match source .md file structure
