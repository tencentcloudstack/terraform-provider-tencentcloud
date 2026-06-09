## ADDED Requirements

### Requirement: Create Just-In-Time Transcode Template
The system SHALL allow users to create a TEO just-in-time transcode template with specified configuration parameters.

#### Scenario: Create template with minimal configuration
- **WHEN** user provides zone_id, template_name, and required stream configurations
- **THEN** system creates a new just-in-time transcode template
- **AND** system returns template_id in the response
- **AND** system generates a composite resource ID as zone_id#template_id
- **AND** system waits for the template to become available via polling

#### Scenario: Create template with full configuration
- **WHEN** user provides all optional parameters including comment, video_stream_switch, audio_stream_switch, video_template, and audio_template
- **THEN** system creates the template with all specified configurations
- **AND** video_stream_switch is set to "on" or "off"
- **AND** audio_stream_switch is set to "on" or "off"
- **AND** video_template is required when video_stream_switch is "on"
- **AND** audio_template is required when audio_stream_switch is "on"

#### Scenario: Create template with video stream disabled
- **WHEN** user sets video_stream_switch to "off"
- **THEN** system creates the template without video configuration
- **AND** system does not require video_template parameter

#### Scenario: Create template with audio stream disabled
- **WHEN** user sets audio_stream_switch to "off"
- **THEN** system creates the template without audio configuration
- **AND** system does not require audio_template parameter

#### Scenario: Handle creation timeout
- **WHEN** template creation does not complete within the configured timeout period
- **THEN** system returns a timeout error
- **AND** system includes information about the timeout configuration

### Requirement: Read Just-In-Time Transcode Template
The system SHALL allow users to retrieve the configuration of an existing just-in-time transcode template.

#### Scenario: Read template by composite ID
- **WHEN** user queries a template with a valid composite ID (zone_id#template_id)
- **THEN** system returns the complete template configuration
- **AND** response includes all parameters: zone_id, template_id, template_name, comment, video_stream_switch, audio_stream_switch, video_template, and audio_template
- **AND** system correctly parses the composite ID to extract zone_id and template_id

#### Scenario: Read non-existent template
- **WHEN** user queries a template with an invalid composite ID
- **THEN** system returns an error indicating the template does not exist
- **AND** system does not attempt to create or modify any resources

#### Scenario: Read template with retry on inconsistency
- **WHEN** there is a temporary inconsistency in the TEO service
- **THEN** system retries the read operation up to the retry limit
- **AND** system returns the template configuration when consistency is restored
- **AND** system returns an error if retry limit is exceeded

#### Scenario: Read template with eventual consistency
- **WHEN** template was just created or updated
- **THEN** system waits and polls until the template becomes readable
- **AND** system returns the latest configuration when available

### Requirement: Delete Just-In-Time Transcode Template
The system SHALL allow users to delete an existing just-in-time transcode template.

#### Scenario: Delete template by composite ID
- **WHEN** user requests deletion with a valid composite ID (zone_id#template_id)
- **THEN** system calls the delete API with zone_id and template_id
- **AND** system waits for the deletion to complete
- **AND** system removes the resource from Terraform state
- **AND** system returns success when deletion is confirmed

#### Scenario: Delete non-existent template
- **WHEN** user attempts to delete a template that does not exist
- **THEN** system returns an error indicating the template was not found
- **AND** system does not modify any existing resources

#### Scenario: Handle deletion timeout
- **WHEN** template deletion does not complete within the configured timeout period
- **THEN** system returns a timeout error
- **AND** system includes information about the timeout configuration
- **AND** system does not remove the resource from Terraform state

#### Scenario: Delete template with retry on failure
- **WHEN** the delete operation fails due to temporary service issues
- **THEN** system retries the delete operation up to the retry limit
- **AND** system returns success when deletion eventually succeeds
- **AND** system returns an error if retry limit is exceeded

### Requirement: Update Just-In-Time Transcode Template
The system SHALL force new resource creation when any parameter changes are requested (since Update API is not available).

#### Scenario: Update any parameter forces recreation
- **WHEN** user modifies any parameter (template_name, comment, video_stream_switch, audio_stream_switch, video_template, or audio_template)
- **THEN** system marks the resource for recreation
- **AND** system deletes the existing template
- **AND** system creates a new template with updated parameters
- **AND** system generates a new template_id
- **AND** system updates the Terraform state with the new composite ID

#### Scenario: Update zone_id forces recreation
- **WHEN** user changes the zone_id parameter
- **THEN** system marks the resource for recreation
- **AND** system deletes the template from the original zone
- **AND** system creates a new template in the new zone
- **AND** system updates the Terraform state with the new zone_id#template_id

### Requirement: Resource Schema Validation
The system SHALL validate all input parameters according to the TEO API constraints.

#### Scenario: Validate template_name length
- **WHEN** user provides a template_name exceeding 64 characters
- **THEN** system returns a validation error
- **AND** system includes information about the 64 character limit

#### Scenario: Validate comment length
- **WHEN** user provides a comment exceeding 256 characters
- **THEN** system returns a validation error
- **AND** system includes information about the 256 character limit

#### Scenario: Validate video_stream_switch values
- **WHEN** user provides a value other than "on" or "off" for video_stream_switch
- **THEN** system returns a validation error
- **AND** system lists the valid values: "on" and "off"

#### Scenario: Validate audio_stream_switch values
- **WHEN** user provides a value other than "on" or "off" for audio_stream_switch
- **THEN** system returns a validation error
- **AND** system lists the valid values: "on" and "off"

#### Scenario: Validate required video_template when video enabled
- **WHEN** user sets video_stream_switch to "on" but does not provide video_template
- **THEN** system returns a validation error
- **AND** system indicates that video_template is required when video_stream_switch is "on"

#### Scenario: Validate required audio_template when audio enabled
- **WHEN** user sets audio_stream_switch to "on" but does not provide audio_template
- **THEN** system returns a validation error
- **AND** system indicates that audio_template is required when audio_stream_switch is "on"

### Requirement: Timeout Configuration
The system SHALL support configurable timeouts for asynchronous operations.

#### Scenario: Configure create timeout
- **WHEN** user specifies a custom create timeout
- **THEN** system uses the custom timeout for template creation
- **AND** system waits up to the specified duration before returning a timeout error

#### Scenario: Configure delete timeout
- **WHEN** user specifies a custom delete timeout
- **THEN** system uses the custom timeout for template deletion
- **AND** system waits up to the specified duration before returning a timeout error

#### Scenario: Use default timeouts
- **WHEN** user does not specify custom timeouts
- **THEN** system uses the default timeout values (e.g., 10 minutes)
- **AND** system applies the same defaults for create, delete, and read operations

### Requirement: Video Template Configuration
The system SHALL support video template parameters for configuring video stream transcoding.

#### Scenario: Configure video codec
- **WHEN** user specifies video_codec in video_template
- **THEN** system passes the codec to the TEO API
- **AND** system stores the codec in Terraform state

#### Scenario: Configure frame rate
- **WHEN** user specifies fps in video_template
- **THEN** system passes the frame rate to the TEO API
- **AND** system stores the frame rate in Terraform state

#### Scenario: Configure video bitrate
- **WHEN** user specifies bitrate in video_template
- **THEN** system passes the bitrate to the TEO API
- **AND** system stores the bitrate in Terraform state

#### Scenario: Configure video resolution
- **WHEN** user specifies resolution parameters (width, height) in video_template
- **THEN** system passes the resolution to the TEO API
- **AND** system stores the resolution in Terraform state

### Requirement: Audio Template Configuration
The system SHALL support audio template parameters for configuring audio stream transcoding.

#### Scenario: Configure audio codec
- **WHEN** user specifies audio_codec in audio_template
- **THEN** system passes the codec to the TEO API
- **AND** system stores the codec in Terraform state

#### Scenario: Configure audio bitrate
- **WHEN** user specifies bitrate in audio_template
- **THEN** system passes the bitrate to the TEO API
- **AND** system stores the bitrate in Terraform state

#### Scenario: Configure audio sample rate
- **WHEN** user specifies sample_rate in audio_template
- **THEN** system passes the sample rate to the TEO API
- **AND** system stores the sample rate in Terraform state

#### Scenario: Configure audio channels
- **WHEN** user specifies channels in audio_template
- **THEN** system passes the channels to the TEO API
- **AND** system stores the channels in Terraform state

### Requirement: Error Handling
The system SHALL provide clear error messages for all failure scenarios.

#### Scenario: Handle API errors
- **WHEN** the TEO API returns an error
- **THEN** system propagates the error to the user
- **AND** system includes relevant error details from the API response
- **AND** system does not crash or leave the resource in an inconsistent state

#### Scenario: Handle network errors
- **WHEN** a network error occurs during API call
- **THEN** system retries the operation according to the retry policy
- **AND** system returns an error if retries are exhausted
- **AND** system includes information about the retry attempts

#### Scenario: Handle authentication errors
- **WHEN** authentication fails (invalid credentials)
- **THEN** system returns a clear authentication error
- **AND** system guides user to check credentials
- **AND** system does not retry authentication failures

#### Scenario: Handle permission errors
- **WHEN** user lacks permission to perform the operation
- **THEN** system returns a clear permission error
- **AND** system indicates the required permission
- **AND** system does not retry permission errors
