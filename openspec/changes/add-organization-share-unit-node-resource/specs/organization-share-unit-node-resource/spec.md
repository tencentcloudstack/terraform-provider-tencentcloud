## ADDED Requirements

### Requirement: Resource SHALL support creating share unit nodes
The `tencentcloud_organization_org_share_unit_node` resource SHALL support creating a new share unit node by calling the AddShareUnitNode API.

#### Scenario: Successfully create share unit node
- **WHEN** user provides valid unit_id and node_id
- **THEN** resource MUST call AddShareUnitNode API with provided parameters
- **THEN** resource MUST set resource ID to composite format `{unit_id}#{node_id}`
- **THEN** resource MUST read back the created node to populate state

#### Scenario: Handle share unit node over limit error
- **WHEN** creating share unit node exceeds the limit
- **THEN** resource MUST return error with message indicating LimitExceeded.ShareUnitNodeOverLimit
- **THEN** resource MUST NOT set resource ID

#### Scenario: Handle organization node not exist error
- **WHEN** specified node_id does not exist in organization
- **THEN** resource MUST return error with message indicating ResourceNotFound.OrganizationNodeNotExist
- **THEN** resource MUST NOT set resource ID

### Requirement: Resource SHALL support reading share unit nodes
The resource SHALL support reading share unit node details by calling the DescribeShareUnitNodes API.

#### Scenario: Successfully read existing share unit node
- **WHEN** resource ID exists and node is present in share unit
- **THEN** resource MUST parse composite ID to extract unit_id and node_id
- **THEN** resource MUST call DescribeShareUnitNodes API with unit_id and node_id
- **THEN** resource MUST populate state with unit_id and node_id

#### Scenario: Handle missing share unit node gracefully
- **WHEN** share unit node does not exist in API response
- **THEN** resource MUST set resource ID to empty string
- **THEN** resource MUST NOT return an error
- **THEN** Terraform MUST mark resource as destroyed

#### Scenario: Handle invalid composite ID format
- **WHEN** resource ID does not match format `{unit_id}#{node_id}`
- **THEN** resource MUST return error with message "id is broken,{actual_id}"

### Requirement: Resource SHALL support deleting share unit nodes
The resource SHALL support deleting share unit nodes by calling the DeleteShareUnitNode API.

#### Scenario: Successfully delete share unit node
- **WHEN** user destroys the resource
- **THEN** resource MUST parse composite ID to extract unit_id and node_id
- **THEN** resource MUST call DeleteShareUnitNode API with unit_id and node_id
- **THEN** resource MUST handle API call success

#### Scenario: Handle already deleted share unit node
- **WHEN** share unit node does not exist during deletion
- **THEN** resource MUST handle FailedOperation.ShareNodeNotExist error gracefully
- **THEN** resource MUST NOT return an error to user

### Requirement: Resource SHALL support importing share unit nodes
The resource SHALL support importing existing share unit nodes using composite ID format.

#### Scenario: Successfully import share unit node
- **WHEN** user runs `terraform import` with format `{unit_id}#{node_id}`
- **THEN** resource MUST set resource ID to provided value
- **THEN** resource MUST call Read function to populate state
- **THEN** resource MUST verify node exists in share unit

#### Scenario: Import with invalid ID format
- **WHEN** user provides ID not matching `{unit_id}#{node_id}` format
- **THEN** resource MUST return error during Read phase
- **THEN** import operation MUST fail with clear error message

### Requirement: Resource schema SHALL enforce ForceNew on all attributes
All resource attributes SHALL be marked as ForceNew since the API does not support update operations.

#### Scenario: Changing unit_id forces recreation
- **WHEN** user modifies unit_id in Terraform configuration
- **THEN** Terraform MUST plan to destroy existing resource
- **THEN** Terraform MUST plan to create new resource with new unit_id

#### Scenario: Changing node_id forces recreation
- **WHEN** user modifies node_id in Terraform configuration
- **THEN** Terraform MUST plan to destroy existing resource  
- **THEN** Terraform MUST plan to create new resource with new node_id

### Requirement: Resource SHALL handle API rate limiting
The resource SHALL implement retry logic to handle API rate limiting (20 requests/second).

#### Scenario: Retry on transient failures
- **WHEN** API call fails with transient error
- **THEN** resource MUST retry using helper.Retry() with tccommon.WriteRetryTimeout
- **THEN** resource MUST use ratelimit.Check() before each API call

#### Scenario: Fail after retry timeout
- **WHEN** API call continues to fail after retry timeout
- **THEN** resource MUST return error to user
- **THEN** error message MUST include original API error details

### Requirement: Resource SHALL provide comprehensive error logging
The resource SHALL log all API interactions for debugging purposes.

#### Scenario: Log successful API calls
- **WHEN** any API call succeeds
- **THEN** resource MUST log request body using request.ToJsonString()
- **THEN** resource MUST log response body using response.ToJsonString()
- **THEN** logs MUST include logId for correlation

#### Scenario: Log failed API calls
- **WHEN** any API call fails
- **THEN** resource MUST log error with CRITAL level
- **THEN** log MUST include request body and error reason
- **THEN** log MUST include logId for correlation

### Requirement: Resource SHALL use consistent error handling patterns
The resource SHALL use defer statements for logging and consistency checks.

#### Scenario: Always log elapsed time
- **WHEN** any CRUD operation executes
- **THEN** resource MUST defer tccommon.LogElapsed() at function start
- **THEN** elapsed time MUST be logged regardless of success or failure

#### Scenario: Always check consistency
- **WHEN** any CRUD operation executes  
- **THEN** resource MUST defer tccommon.InconsistentCheck() at function start
- **THEN** consistency check MUST run regardless of success or failure

### Requirement: Resource SHALL have complete documentation
The resource SHALL provide comprehensive documentation for users.

#### Scenario: Documentation includes usage example
- **WHEN** user views resource documentation
- **THEN** documentation MUST include complete Terraform configuration example
- **THEN** example MUST show required attributes (unit_id, node_id)
- **THEN** example MUST demonstrate import command format

#### Scenario: Documentation describes all arguments
- **WHEN** user views argument reference
- **THEN** documentation MUST list unit_id with description "共享单元ID"
- **THEN** documentation MUST list node_id with description "组织部门ID"
- **THEN** documentation MUST mark all arguments as Required and ForceNew

#### Scenario: Documentation includes generated website version
- **WHEN** make doc is executed
- **THEN** system MUST generate website/docs/r/organization_org_share_unit_node.html.markdown
- **THEN** generated file MUST match source .md file structure
