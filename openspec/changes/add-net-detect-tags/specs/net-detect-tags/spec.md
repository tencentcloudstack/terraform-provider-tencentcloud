## ADDED Requirements

### Requirement: User can create net detect resource with tags
The system SHALL allow users to specify tags when creating a network detect resource.

#### Scenario: Successful creation with tags
- **WHEN** user creates a network detect resource with tags parameter containing key-value pairs
- **THEN** the tags SHALL be passed to the CreateNetDetect API
- **AND** the network detect resource SHALL be created with the specified tags

### Requirement: User can read tags from net detect resource
The system SHALL allow users to read the tags associated with a network detect resource.

#### Scenario: Successful read of existing tags
- **WHEN** user reads a network detect resource that has tags
- **THEN** the system SHALL retrieve the tags from the API response
- **AND** the tags SHALL be set in the Terraform state

#### Scenario: Read resource without tags
- **WHEN** user reads a network detect resource that has no tags
- **THEN** the tags field in the Terraform state SHALL be empty or null
- **AND** no error SHALL be raised

### Requirement: User can update tags on net detect resource
The system SHALL allow users to update the tags of an existing network detect resource.

#### Scenario: Successful update adding new tags
- **WHEN** user updates a network detect resource by adding new tags
- **THEN** the system SHALL call the ModifyTags API with the new tags
- **AND** the tags SHALL be updated in the cloud
- **AND** the tags in the Terraform state SHALL reflect the changes

#### Scenario: Successful update removing tags
- **WHEN** user updates a network detect resource by removing tags (setting tags to empty map)
- **THEN** the system SHALL call the ModifyTags API to remove all tags
- **AND** all tags SHALL be removed from the cloud resource
- **AND** the tags in the Terraform state SHALL be empty

#### Scenario: Successful update modifying existing tags
- **WHEN** user updates a network detect resource by modifying existing tag values
- **THEN** the system SHALL call the ModifyTags API with the modified tags
- **AND** the tag values SHALL be updated in the cloud
- **AND** the tags in the Terraform state SHALL reflect the new values

#### Scenario: No API call when tags unchanged
- **WHEN** user updates a network detect resource but tags remain unchanged
- **THEN** the system SHALL NOT call the ModifyTags API
- **AND** no tag update SHALL occur

### Requirement: Tags parameter is optional
The system SHALL allow users to create or update a network detect resource without specifying tags.

#### Scenario: Creation without tags
- **WHEN** user creates a network detect resource without the tags parameter
- **THEN** the resource SHALL be created successfully
- **AND** the resource SHALL have no tags
- **AND** no error SHALL be raised

#### Scenario: Update without tags parameter
- **WHEN** user updates a network detect resource without modifying the tags parameter
- **THEN** the existing tags SHALL remain unchanged
- **AND** no tag update API SHALL be called

### Requirement: Tags are computed fields
The system SHALL populate the tags field from the API response to ensure consistency between Terraform state and cloud resources.

#### Scenario: Tags populated after create
- **WHEN** user creates a network detect resource with tags
- **THEN** the system SHALL read the created resource
- **AND** the tags SHALL be populated from the API response
- **AND** the tags in Terraform state SHALL match the API response

#### Scenario: Tags populated after read
- **WHEN** user reads an existing network detect resource
- **THEN** the system SHALL retrieve the tags from the API
- **AND** the tags in Terraform state SHALL match the API response

### Requirement: Tags support multiple key-value pairs
The system SHALL support multiple tags on a single network detect resource.

#### Scenario: Multiple tags on resource
- **WHEN** user creates a network detect resource with multiple tags (e.g., "Environment": "test", "Owner": "devops", "Project": "network")
- **THEN** all tags SHALL be created successfully
- **AND** all tags SHALL be readable and updatable

### Requirement: Tags are backward compatible
The system SHALL maintain backward compatibility with existing configurations that do not use tags.

#### Scenario: Existing configuration without tags
- **WHEN** user applies an existing Terraform configuration that does not include tags
- **THEN** the resource SHALL continue to work without changes
- **AND** no tags SHALL be added to the resource
- **AND** no drift SHALL be detected

### Requirement: Tag key and value constraints
The system SHALL validate tag keys and values according to TencentCloud tag constraints.

#### Scenario: Valid tag keys and values
- **WHEN** user specifies tags with valid keys and values (following TencentCloud tag naming rules)
- **THEN** the tags SHALL be accepted and created successfully

#### Scenario: Invalid tag keys or values
- **WHEN** user specifies tags with invalid keys or values
- **THEN** the API SHALL return an error
- **AND** the error SHALL be propagated to the user
- **AND** the resource SHALL not be created or updated
