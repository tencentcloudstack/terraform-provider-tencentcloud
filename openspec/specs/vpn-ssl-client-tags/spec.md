## ADDED Requirements

### Requirement: Tags field in schema
The `tencentcloud_vpn_ssl_client` resource SHALL include a `tags` field in its schema to support tag management.

#### Scenario: Schema includes tags field
- **WHEN** user defines the resource schema
- **THEN** the schema MUST include a `tags` field of type `TypeMap` with `TypeString` elements
- **THEN** the `tags` field MUST be Optional to maintain backward compatibility
- **THEN** the `tags` field MUST NOT have `ForceNew: true` to allow in-place updates
- **THEN** the `tags` field MUST have a description explaining its purpose

### Requirement: Set tags on creation
The resource SHALL support setting tags when creating a VPN SSL Client through the API.

#### Scenario: Create with tags
- **WHEN** user specifies tags in the resource configuration
- **THEN** the Create function MUST convert the tags map to API Tag format
- **THEN** the Create function MUST include tags in the `CreateVpnGatewaySslClient` API request
- **THEN** the tags MUST be set on the resource during creation

#### Scenario: Create without tags
- **WHEN** user does not specify tags in the resource configuration
- **THEN** the Create function MUST succeed without passing tags to the API
- **THEN** the resource MUST be created successfully

### Requirement: Read tags from resource
The resource SHALL synchronize tags from the cloud resource to Terraform state during Read operations.

#### Scenario: Read tags using Tag Service
- **WHEN** the Read function executes
- **THEN** the function MUST use Tag Service to query tags for the resource
- **THEN** the function MUST use service type "vpc", resource type "vpnx", and the SSL client ID
- **THEN** the function MUST set the tags in state using `d.Set("tags", ...)`

#### Scenario: Handle missing tags
- **WHEN** the resource has no tags
- **THEN** the Read function MUST handle empty tag list gracefully
- **THEN** the state MUST reflect an empty or nil tags value

### Requirement: Update tags in-place
The resource SHALL support updating tags without recreating the resource.

#### Scenario: Update function exists
- **WHEN** resource definition is loaded
- **THEN** the resource MUST define an Update callback function
- **THEN** the Update function MUST be registered in the resource schema

#### Scenario: Update tags using Tag Service
- **WHEN** user modifies tags in the configuration and applies changes
- **THEN** the Update function MUST detect the change using `d.HasChange("tags")`
- **THEN** the Update function MUST use `d.GetChange("tags")` to get old and new tags
- **THEN** the Update function MUST use `svctag.DiffTags` to calculate tags to add/modify and delete
- **THEN** the Update function MUST use `tagService.ModifyTags()` to apply changes
- **THEN** the Update function MUST use resource name format: `BuildTagResourceName("vpc", "vpnx", region, sslClientId)`

#### Scenario: Update tags - add new tags
- **WHEN** user adds new tags to existing configuration
- **THEN** the new tags MUST be added to the resource
- **THEN** existing tags MUST be preserved

#### Scenario: Update tags - remove tags
- **WHEN** user removes tags from configuration
- **THEN** the removed tags MUST be deleted from the resource
- **THEN** remaining tags MUST be preserved

#### Scenario: Update tags - modify tag values
- **WHEN** user changes tag values in configuration
- **THEN** the tag values MUST be updated on the resource
- **THEN** other tags MUST be preserved

#### Scenario: No tag changes
- **WHEN** user applies configuration without tag changes
- **THEN** the Update function MUST skip tag modification
- **THEN** no Tag Service API calls MUST be made

### Requirement: Documentation and examples
The resource SHALL provide complete documentation and examples for the tags feature.

#### Scenario: Documentation includes tags
- **WHEN** user reads the resource documentation
- **THEN** the documentation MUST include tags field in the argument reference
- **THEN** the documentation MUST include at least one example using tags
- **THEN** the example MUST show valid tag key-value pairs
- **THEN** the documentation MUST indicate that tags can be updated in-place

#### Scenario: Test coverage for tags
- **WHEN** running acceptance tests
- **THEN** tests MUST include test cases that verify tags are created correctly
- **THEN** tests MUST verify tags are read back into state correctly
- **THEN** tests MUST include tags update scenarios (add, remove, modify)
- **THEN** tests MUST include both scenarios: with tags and without tags
