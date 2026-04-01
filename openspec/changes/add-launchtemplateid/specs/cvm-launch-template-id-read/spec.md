## ADDED Requirements

### Requirement: Schema includes launch_template_id field
The tencentcloud_cvm_launch_template resource SHALL include a `launch_template_id` field in its schema. The field SHALL be of type String and SHALL be Computed (not Optional), indicating it is returned from the API and should not be set by the user.

#### Scenario: Schema field definition
- **WHEN** the resource schema is defined
- **THEN** it includes a `launch_template_id` field with Type: String and Computed: true
- **AND** the field is not marked as Optional or Required

### Requirement: Create operation sets launch_template_id
When creating a CVM launch template, the Create function SHALL read the LaunchTemplateId field from the CreateLaunchTemplate API response and set it to the state's `launch_template_id` field.

#### Scenario: Successful create with launch_template_id
- **WHEN** CreateLaunchTemplate API returns a LaunchTemplateId value (e.g., "lt-a5n1305lo")
- **THEN** the Create function sets d.Set("launch_template_id", "lt-a5n1305lo")
- **AND** the state contains the correct launch_template_id value

#### Scenario: Create response without launch_template_id
- **WHEN** CreateLaunchTemplate API response does not contain LaunchTemplateId field
- **THEN** the Create function ignores the missing field
- **AND** no error is raised
- **AND** the state's launch_template_id remains unset

### Requirement: Read operation updates launch_template_id
When reading a CVM launch template, the Read function SHALL read the LaunchTemplateId field from the DescribeLaunchTemplates API response and update the state's `launch_template_id` field.

#### Scenario: Successful read with launch_template_id
- **WHEN** DescribeLaunchTemplates API returns a template with LaunchTemplateId value
- **THEN** the Read function updates the state's launch_template_id field with the returned value
- **AND** the state contains the correct launch_template_id value

#### Scenario: Read response without launch_template_id
- **WHEN** DescribeLaunchTemplates API response does not contain LaunchTemplateId field for an existing template
- **THEN** the Read function ignores the missing field
- **AND** no error is raised
- **AND** the state's launch_template_id remains unchanged

### Requirement: Backward compatibility
The addition of the `launch_template_id` field SHALL NOT break existing Terraform configurations or state. Existing users' configurations and state files SHALL remain valid without modification.

#### Scenario: Existing configuration without launch_template_id
- **WHEN** a user applies an existing configuration that does not include launch_template_id
- **THEN** the terraform apply completes successfully
- **AND** the state is updated to include the computed launch_template_id value
- **AND** no configuration changes are required from the user
