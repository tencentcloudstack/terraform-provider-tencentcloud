## ADDED Requirements

### Requirement: Create DNS Record
The system SHALL create a new DNS record for a specified TEO zone.

#### Scenario: Successful DNS record creation
- **WHEN** user provides valid zone_id, record_name, record_type, and record_content
- **THEN** system creates a DNS record with the specified parameters
- **AND** system returns a unique record_id
- **AND** system stores the DNS record in Terraform state

#### Scenario: DNS record creation with optional fields
- **WHEN** user provides zone_id, record_name, record_type, record_content, and optional ttl, weight, priority, or location
- **THEN** system creates a DNS record with all specified parameters
- **AND** system applies the optional parameters according to their constraints

#### Scenario: DNS record creation with invalid zone_id
- **WHEN** user provides an invalid zone_id
- **THEN** system returns an error indicating the zone does not exist
- **AND** no DNS record is created

#### Scenario: DNS record creation with missing required fields
- **WHEN** user omits required fields (zone_id, record_name, record_type, or record_content)
- **THEN** system returns an error indicating the missing required fields
- **AND** no DNS record is created

### Requirement: Read DNS Record
The system SHALL retrieve an existing DNS record by its unique identifier.

#### Scenario: Successful DNS record read
- **WHEN** user provides a valid DNS record_id
- **THEN** system retrieves the DNS record with all its attributes
- **AND** system returns record details including zone_id, name, type, content, ttl, weight, priority, location, status, created_on, and modified_on

#### Scenario: DNS record read with non-existent record_id
- **WHEN** user provides a record_id that does not exist
- **THEN** system returns ResourceNotFound error
- **AND** Terraform state is updated to reflect the resource does not exist

#### Scenario: DNS record read after update
- **WHEN** user reads a DNS record after it has been updated
- **THEN** system returns the most recent state of the DNS record
- **AND** all modified fields reflect their updated values

### Requirement: Update DNS Record
The system SHALL update an existing DNS record with new values.

#### Scenario: Successful DNS record update
- **WHEN** user provides valid record_id and one or more updatable fields (content, ttl, weight, priority, or location)
- **THEN** system updates the DNS record with the new values
- **AND** system preserves the record_id, zone_id, name, type, and other immutable fields

#### Scenario: DNS record update of name field
- **WHEN** user attempts to update the record_name field
- **THEN** system returns an error indicating that name cannot be modified
- **AND** no changes are applied to the DNS record

#### Scenario: DNS record update of type field
- **WHEN** user attempts to update the record_type field
- **THEN** system returns an error indicating that type cannot be modified
- **AND** no changes are applied to the DNS record

#### Scenario: DNS record update with non-existent record_id
- **WHEN** user provides a record_id that does not exist
- **THEN** system returns ResourceNotFound error
- **AND** no updates are applied

### Requirement: Delete DNS Record
The system SHALL delete an existing DNS record by its unique identifier.

#### Scenario: Successful DNS record deletion
- **WHEN** user provides a valid DNS record_id
- **THEN** system deletes the DNS record from the cloud
- **AND** system removes the DNS record from Terraform state

#### Scenario: DNS record deletion with non-existent record_id
- **WHEN** user provides a record_id that does not exist
- **THEN** system returns ResourceNotFound error
- **AND** Terraform state is updated to reflect the resource does not exist

#### Scenario: DNS record deletion during read-after-write consistency
- **WHEN** system deletes a DNS record and immediately reads it
- **THEN** system implements retry logic to handle eventual consistency
- **AND** system confirms deletion before removing from state

### Requirement: DNS Record Types
The system SHALL support multiple DNS record types with appropriate constraints.

#### Scenario: A record creation
- **WHEN** user creates a DNS record with type "A"
- **THEN** system validates that the content is a valid IPv4 address
- **AND** system creates the A record

#### Scenario: AAAA record creation
- **WHEN** user creates a DNS record with type "AAAA"
- **THEN** system validates that the content is a valid IPv6 address
- **AND** system creates the AAAA record

#### Scenario: CNAME record creation
- **WHEN** user creates a DNS record with type "CNAME"
- **THEN** system validates that the content is a valid domain name
- **AND** system creates the CNAME record

#### Scenario: MX record creation with priority
- **WHEN** user creates a DNS record with type "MX"
- **THEN** system validates that the priority is in range 0-50
- **AND** system creates the MX record with the specified priority

#### Scenario: TXT record creation
- **WHEN** user creates a DNS record with type "TXT"
- **THEN** system accepts any text content for the record
- **AND** system creates the TXT record

#### Scenario: NS record creation
- **WHEN** user creates a DNS record with type "NS"
- **THEN** system validates that the content is a valid domain name
- **AND** system creates the NS record

#### Scenario: CAA record creation
- **WHEN** user creates a DNS record with type "CAA"
- **THEN** system validates that the content follows CAA record format
- **AND** system creates the CAA record

#### Scenario: SRV record creation
- **WHEN** user creates a DNS record with type "SRV"
- **THEN** system validates that the content follows SRV record format
- **AND** system creates the SRV record

### Requirement: DNS Record Weight Configuration
The system SHALL support DNS record weight configuration for load balancing.

#### Scenario: Weight configuration for A record
- **WHEN** user creates or updates an A record with weight value between -1 and 100
- **THEN** system applies the weight configuration for load balancing
- **AND** weight value 0 means the record will not resolve

#### Scenario: Weight configuration with invalid value
- **WHEN** user provides weight value outside the range -1 to 100
- **THEN** system returns an error indicating the weight value is invalid

#### Scenario: Weight configuration for unsupported record type
- **WHEN** user attempts to set weight for a record type that does not support weight (non-A, non-AAAA, non-CNAME)
- **THEN** system returns an error indicating weight is not supported for this record type

### Requirement: DNS Record TTL Configuration
The system SHALL support DNS record TTL (Time To Live) configuration.

#### Scenario: TTL configuration with valid range
- **WHEN** user provides TTL value between 60 and 86400 seconds
- **THEN** system applies the TTL configuration
- **AND** DNS resolvers cache the record for the specified duration

#### Scenario: TTL configuration with invalid range
- **WHEN** user provides TTL value outside the range 60 to 86400 seconds
- **THEN** system returns an error indicating the TTL value is invalid

#### Scenario: TTL configuration default value
- **WHEN** user does not provide TTL value
- **THEN** system uses the default value of 300 seconds

### Requirement: DNS Record Location Configuration
The system SHALL support DNS record location (resolution line) configuration.

#### Scenario: Location configuration for supported record types
- **WHEN** user creates or updates an A, AAAA, or CNAME record with location value
- **THEN** system applies the location configuration for geographical routing
- **AND** DNS queries from specified regions resolve to this record

#### Scenario: Location configuration for unsupported record types
- **WHEN** user attempts to set location for a record type that does not support location (non-A, non-AAAA, non-CNAME)
- **THEN** system returns an error indicating location is not supported for this record type

#### Scenario: Location configuration default value
- **WHEN** user does not provide location value
- **THEN** system uses the default value "Default" (all regions)

### Requirement: DNS Record Idempotency
The system SHALL ensure idempotent operations for DNS record management.

#### Scenario: Create with same parameters
- **WHEN** user creates a DNS record with the same parameters as an existing record
- **THEN** system detects the duplicate
- **AND** system returns the existing record_id
- **AND** no duplicate record is created

#### Scenario: Update with same values
- **WHEN** user updates a DNS record with values identical to current state
- **THEN** system detects no changes
- **AND** system returns success without making API calls

#### Scenario: Delete non-existent resource
- **WHEN** user deletes a DNS record that does not exist
- **THEN** system returns success without error
- **AND** Terraform state reflects the resource as deleted

### Requirement: DNS Record Error Handling
The system SHALL provide clear error messages for DNS record operations.

#### Scenario: API rate limit exceeded
- **WHEN** system exceeds the API rate limit
- **THEN** system returns a clear error message indicating rate limit exceeded
- **AND** system implements retry logic with exponential backoff

#### Scenario: Network error during operation
- **WHEN** system encounters a network error during DNS record operation
- **THEN** system returns a clear error message indicating network failure
- **AND** system implements retry logic to handle transient errors

#### Scenario: Validation error for invalid input
- **WHEN** user provides invalid input (e.g., invalid IP address for A record)
- **THEN** system returns a clear error message indicating validation failure
- **AND** system specifies which field is invalid and why

### Requirement: DNS Record Timeout Configuration
The system SHALL support timeout configuration for DNS record operations.

#### Scenario: Create operation timeout
- **WHEN** create operation exceeds the configured timeout
- **THEN** system cancels the operation
- **AND** system returns a timeout error

#### Scenario: Update operation timeout
- **WHEN** update operation exceeds the configured timeout
- **THEN** system cancels the operation
- **AND** system returns a timeout error

#### Scenario: Delete operation timeout
- **WHEN** delete operation exceeds the configured timeout
- **THEN** system cancels the operation
- **AND** system returns a timeout error

### Requirement: DNS Record Composite ID
The system SHALL use a composite ID format for DNS record identification.

#### Scenario: Composite ID format
- **WHEN** system creates or reads a DNS record
- **THEN** system uses composite ID in format "zoneId#recordId"
- **AND** system can parse the composite ID to extract zoneId and recordId

#### Scenario: Composite ID parsing
- **WHEN** system needs to extract zoneId and recordId from composite ID
- **THEN** system splits the ID by "#" delimiter
- **AND** first part represents zoneId
- **AND** second part represents recordId

#### Scenario: Invalid composite ID format
- **WHEN** system encounters a composite ID with invalid format (missing # delimiter)
- **THEN** system returns an error indicating invalid ID format
- **AND** system provides guidance on correct ID format
