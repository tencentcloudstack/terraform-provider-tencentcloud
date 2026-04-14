## ADDED Requirements

### Requirement: Resource creation
The system SHALL allow users to create a TEO alias domain resource through Terraform. The CreateAliasDomain API SHALL be called with zone_id, alias_name, and target_name parameters. The system SHALL wait for the resource to be available via DescribeAliasDomains API before returning success.

#### Scenario: Successful resource creation
- **WHEN** user applies a Terraform configuration with a new tc_teo_alias_domain resource
- **THEN** system calls CreateAliasDomain API with zone_id, alias_name, target_name
- **AND** system polls DescribeAliasDomains API until the alias domain appears in the response
- **AND** system stores the resource in Terraform state with id as zone_id#alias_name

#### Scenario: Resource creation with paused state
- **WHEN** user creates a tc_teo_alias_domain resource with paused=false
- **THEN** system calls CreateAliasDomain API to create the resource
- **AND** system calls ModifyAliasDomainStatus API to set paused=false after creation
- **AND** system polls DescribeAliasDomains API until paused status reflects the desired state

#### Scenario: Resource creation with zone_id only
- **WHEN** user specifies only zone_id without alias_name or target_name
- **THEN** system returns validation error
- **AND** Terraform reports required parameter missing

### Requirement: Resource reading
The system SHALL allow users to read the state of a TEO alias domain resource. The DescribeAliasDomains API SHALL be called with zone_id to retrieve alias domain information. The system SHALL map API response fields to Terraform schema fields including zone_id, alias_name, target_name, and paused.

#### Scenario: Successful resource reading
- **WHEN** user runs terraform refresh or plan
- **THEN** system calls DescribeAliasDomains API with zone_id from resource id
- **AND** system matches the alias_name from resource id
- **AND** system updates Terraform state with current zone_id, alias_name, target_name, and paused values

#### Scenario: Resource reading for non-existent resource
- **WHEN** system reads a resource that no longer exists in cloud
- **THEN** system removes the resource from Terraform state
- **AND** system reports the resource has been deleted externally

### Requirement: Resource update
The system SHALL allow users to update a TEO alias domain resource. The system SHALL detect which parameters have changed and call the appropriate API:
- If target_name changes: call ModifyAliasDomain API
- If paused changes: call ModifyAliasDomainStatus API
- If both change: call both APIs sequentially
The system SHALL wait for each operation to complete before proceeding.

#### Scenario: Update target_name
- **WHEN** user changes target_name in the Terraform configuration
- **THEN** system calls ModifyAliasDomain API with zone_id, alias_name, and new target_name
- **AND** system polls DescribeAliasDomains API until target_name reflects the new value

#### Scenario: Update paused state from false to true
- **WHEN** user changes paused from false to true
- **THEN** system calls ModifyAliasDomainStatus API with zone_id, alias_names=[alias_name], paused=true
- **AND** system polls DescribeAliasDomains API until paused status is true

#### Scenario: Update paused state from true to false
- **WHEN** user changes paused from true to false
- **THEN** system calls ModifyAliasDomainStatus API with zone_id, alias_names=[alias_name], paused=false
- **AND** system polls DescribeAliasDomains API until paused status is false

#### Scenario: Update both target_name and paused
- **WHEN** user changes both target_name and paused simultaneously
- **THEN** system calls ModifyAliasDomain API first to update target_name
- **AND** system waits for target_name change to complete
- **AND** system then calls ModifyAliasDomainStatus API to update paused
- **AND** system waits for paused status change to complete

#### Scenario: No changes detected
- **WHEN** user applies configuration with no parameter changes
- **THEN** system does not call any modification APIs
- **AND** system reports no changes to apply

### Requirement: Resource deletion
The system SHALL allow users to delete a TEO alias domain resource. The DeleteAliasDomain API SHALL be called with zone_id and alias_name. The system SHALL wait for the resource to be removed from DescribeAliasDomains response before returning success.

#### Scenario: Successful resource deletion
- **WHEN** user runs terraform destroy or removes resource from configuration
- **THEN** system calls DeleteAliasDomain API with zone_id and alias_name from resource id
- **AND** system polls DescribeAliasDomains API until the alias domain no longer appears
- **AND** system removes the resource from Terraform state

#### Scenario: Delete already deleted resource
- **WHEN** system attempts to delete a resource that no longer exists
- **THEN** system calls DescribeAliasDomains API and confirms resource is not found
- **AND** system removes the resource from Terraform state without error
- **AND** deletion operation is idempotent

### Requirement: Resource identifier
The system SHALL use zone_id#alias_name as the composite resource identifier. The # separator SHALL be used to join zone_id and alias_name into a unique id that can be parsed back into its components.

#### Scenario: Parse resource id during read
- **WHEN** system reads a resource with id zone-123#example.com
- **THEN** system extracts zone_id as zone-123
- **AND** system extracts alias_name as example.com
- **AND** system uses these values for API calls

#### Scenario: Generate resource id during create
- **WHEN** system creates a resource with zone_id=zone-123 and alias_name=example.com
- **THEN** system generates id as zone-123#example.com
- **AND** system stores this id in Terraform state

### Requirement: Asynchronous operation handling
The system SHALL handle all asynchronous operations (CreateAliasDomain, ModifyAliasDomain, ModifyAliasDomainStatus, DeleteAliasDomain) by polling the DescribeAliasDomains API until the expected state is reached. The system SHALL use the retry mechanism with configurable timeout.

#### Scenario: Poll until resource appears after creation
- **WHEN** system creates a new alias domain
- **THEN** system polls DescribeAliasDomains API at regular intervals
- **AND** system continues until the alias domain appears in the response
- **OR** system times out and returns error

#### Scenario: Poll until target_name changes after update
- **WHEN** system updates target_name
- **THEN** system polls DescribeAliasDomains API at regular intervals
- **AND** system continues until target_name reflects the new value
- **OR** system times out and returns error

#### Scenario: Poll until paused status changes after status update
- **WHEN** system updates paused status
- **THEN** system polls DescribeAliasDomains API at regular intervals
- **AND** system continues until paused reflects the new status
- **OR** system times out and returns error

#### Scenario: Poll until resource disappears after deletion
- **WHEN** system deletes an alias domain
- **THEN** system polls DescribeAliasDomains API at regular intervals
- **AND** system continues until the alias domain is not in the response
- **OR** system times out and returns error

### Requirement: Timeouts configuration
The system SHALL support configurable timeouts for create, update, and delete operations. The schema SHALL declare a Timeouts block with default values for each operation. Users SHALL be able to customize timeout durations through Terraform configuration.

#### Scenario: Use default timeout for creation
- **WHEN** user creates a resource without specifying timeout
- **THEN** system uses default timeout duration for create operation
- **AND** system polls DescribeAliasDomains API until success or default timeout expires

#### Scenario: Use custom timeout for update
- **WHEN** user specifies custom timeout in timeouts block
- **THEN** system uses the specified timeout duration for update operation
- **AND** system polls DescribeAliasDomains API until success or custom timeout expires

### Requirement: Schema parameters
The system SHALL define the following schema parameters for the tc_teo_alias_domain resource:
- zone_id: (String, Required) The zone ID of the site
- alias_name: (String, Required) The alias domain name
- target_name: (String, Required) The target domain name
- paused: (Bool, Optional, Computed) The paused status of the alias domain (false=enabled, true=paused)

#### Scenario: Validate required parameters
- **WHEN** user creates a resource without zone_id
- **THEN** system returns validation error
- **AND** Terraform reports zone_id is required

- **WHEN** user creates a resource without alias_name
- **THEN** system returns validation error
- **AND** Terraform reports alias_name is required

- **WHEN** user creates a resource without target_name
- **THEN** system returns validation error
- **AND** Terraform reports target_name is required

#### Scenario: Use optional paused parameter
- **WHEN** user creates a resource without paused
- **THEN** system creates resource with default paused status from API
- **AND** system reads paused value from DescribeAliasDomains response
- **AND** paused is marked as Computed in state

#### Scenario: Update paused parameter
- **WHEN** user sets paused=false explicitly
- **THEN** system calls ModifyAliasDomainStatus API to enable the alias domain
- **AND** system confirms paused=false in state

- **WHEN** user sets paused=true explicitly
- **THEN** system calls ModifyAliasDomainStatus API to pause the alias domain
- **AND** system confirms paused=true in state

### Requirement: Error handling
The system SHALL handle errors from all API calls appropriately. The system SHALL retry transient errors, return meaningful error messages for permanent failures, and log all API errors for debugging.

#### Scenario: Handle API authentication failure
- **WHEN** API call fails due to invalid credentials
- **THEN** system returns authentication error
- **AND** Terraform displays user-friendly error message

#### Scenario: Handle resource not found error during read
- **WHEN** DescribeAliasDomains API returns ResourceNotFound for existing resource
- **THEN** system removes resource from Terraform state
- **AND** system does not return error

#### Scenario: Handle resource not found error during delete
- **WHEN** DeleteAliasDomain API returns ResourceNotFound
- **THEN** system confirms deletion success
- **AND** system proceeds to remove resource from Terraform state

#### Scenario: Handle invalid parameter error during create
- **WHEN** CreateAliasDomain API returns InvalidParameter
- **THEN** system returns error with parameter details
- **AND** Terraform displays which parameter is invalid

#### Scenario: Handle quota exceeded error
- **WHEN** API call returns QuotaExceeded
- **THEN** system returns quota limit error
- **AND** Terraform suggests contacting support to increase quota
