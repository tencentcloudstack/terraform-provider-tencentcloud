## ADDED Requirements

### Requirement: Resource creation
The system SHALL allow users to create a teo DNS record by providing zone_id, name, type, and content parameters.

#### Scenario: Create A record successfully
- **WHEN** user provides zone_id, name="www.example.com", type="A", content="1.2.3.4"
- **THEN** system creates the DNS record and returns record_id

#### Scenario: Create CNAME record successfully
- **WHEN** user provides zone_id, name="www.example.com", type="CNAME", content="alias.example.com"
- **THEN** system creates the DNS record and returns record_id

#### Scenario: Create MX record with priority
- **WHEN** user provides zone_id, name="example.com", type="MX", content="mail.example.com", priority=10
- **THEN** system creates the DNS record with priority and returns record_id

#### Scenario: Create record with TTL
- **WHEN** user provides zone_id, name, type, content, ttl=600
- **THEN** system creates the DNS record with specified TTL

#### Scenario: Create record with location and weight
- **WHEN** user provides zone_id, name="www.example.com", type="A", content="1.2.3.4", location="Asia", weight=10
- **THEN** system creates the DNS record with location and weight

#### Scenario: Create record fails with invalid zone_id
- **WHEN** user provides non-existent zone_id
- **THEN** system returns an error indicating zone not found

#### Scenario: Create record fails with missing required fields
- **WHEN** user omits zone_id, name, type, or content
- **THEN** system returns validation error

### Requirement: Resource read
The system SHALL allow users to read a teo DNS record by its resource ID (zone_id#record_id).

#### Scenario: Read existing record successfully
- **WHEN** user provides valid resource_id "zone123#record456"
- **THEN** system returns the DNS record with all fields including status, created_on, modified_on

#### Scenario: Read record fails with invalid resource_id
- **WHEN** user provides invalid resource_id format
- **THEN** system returns error indicating invalid ID format

#### Scenario: Read record fails with non-existent record
- **WHEN** user provides resource_id for non-existent record
- **THEN** system returns error indicating record not found

#### Scenario: Read record returns computed fields
- **WHEN** user reads an existing record
- **THEN** system returns status, created_on, modified_on as computed fields

### Requirement: Resource update
The system SHALL allow users to update a teo DNS record by modifying its configurable fields (name, type, content, location, ttl, weight, priority).

#### Scenario: Update record content successfully
- **WHEN** user updates content from "1.2.3.4" to "5.6.7.8"
- **THEN** system updates the DNS record content

#### Scenario: Update record TTL successfully
- **WHEN** user updates ttl from 300 to 600
- **THEN** system updates the DNS record TTL

#### Scenario: Update record weight successfully
- **WHEN** user updates weight from 10 to 20
- **THEN** system updates the DNS record weight

#### Scenario: Update record priority successfully
- **WHEN** user updates priority from 10 to 5
- **THEN** system updates the DNS record priority

#### Scenario: Update record with multiple fields
- **WHEN** user updates content, ttl, and weight simultaneously
- **THEN** system updates all specified fields

#### Scenario: Update record type successfully
- **WHEN** user updates type from "A" to "CNAME" and content accordingly
- **THEN** system updates the DNS record type and content

#### Scenario: Update record fails with non-existent record
- **WHEN** user attempts to update non-existent record
- **THEN** system returns error indicating record not found

#### Scenario: Update record ignores read-only fields
- **WHEN** user tries to update status, created_on, or modified_on
- **THEN** system ignores these fields and does not update them

### Requirement: Resource deletion
The system SHALL allow users to delete a teo DNS record by its resource ID.

#### Scenario: Delete existing record successfully
- **WHEN** user deletes a record with resource_id "zone123#record456"
- **THEN** system deletes the DNS record

#### Scenario: Delete record fails with non-existent record
- **WHEN** user attempts to delete non-existent record
- **THEN** system returns error indicating record not found

#### Scenario: Delete record fails with invalid resource_id
- **WHEN** user provides invalid resource_id format
- **THEN** system returns error indicating invalid ID format

### Requirement: Record type validation
The system SHALL validate that location and weight fields are only used with A, AAAA, or CNAME record types, and priority field is only used with MX record type.

#### Scenario: Validate A record with location and weight
- **WHEN** user creates A record with location and weight
- **THEN** system accepts the configuration

#### Scenario: Validate AAAA record with location and weight
- **WHEN** user creates AAAA record with location and weight
- **THEN** system accepts the configuration

#### Scenario: Validate CNAME record with location and weight
- **WHEN** user creates CNAME record with location and weight
- **THEN** system accepts the configuration

#### Scenario: Validate MX record with priority
- **WHEN** user creates MX record with priority
- **THEN** system accepts the configuration

#### Scenario: Reject location and weight for non-A/AAAA/CNAME records
- **WHEN** user creates TXT record with location or weight
- **THEN** system returns validation error

#### Scenario: Reject priority for non-MX records
- **WHEN** user creates A record with priority
- **THEN** system returns validation error

### Requirement: TTL value validation
The system SHALL validate that TTL value is within the allowed range (60-86400 seconds).

#### Scenario: Accept valid TTL value
- **WHEN** user provides ttl=600
- **THEN** system accepts the TTL value

#### Scenario: Reject TTL value below minimum
- **WHEN** user provides ttl=30
- **THEN** system returns validation error

#### Scenario: Reject TTL value above maximum
- **WHEN** user provides ttl=100000
- **THEN** system returns validation error

### Requirement: Weight value validation
The system SHALL validate that weight value is within the allowed range (-1 to 100), where -1 means no weight and 0 means no resolution.

#### Scenario: Accept valid weight value
- **WHEN** user provides weight=50
- **THEN** system accepts the weight value

#### Scenario: Accept weight=-1 for no weight
- **WHEN** user provides weight=-1
- **THEN** system accepts the weight value indicating no weight set

#### Scenario: Accept weight=0 for no resolution
- **WHEN** user provides weight=0
- **THEN** system accepts the weight value indicating no resolution

#### Scenario: Reject weight value below minimum
- **WHEN** user provides weight=-10
- **THEN** system returns validation error

#### Scenario: Reject weight value above maximum
- **WHEN** user provides weight=150
- **THEN** system returns validation error

### Requirement: Priority value validation
The system SHALL validate that priority value is within the allowed range (0-50) for MX records.

#### Scenario: Accept valid priority value
- **WHEN** user provides priority=10
- **THEN** system accepts the priority value

#### Scenario: Reject priority value below minimum
- **WHEN** user provides priority=-1
- **THEN** system returns validation error

#### Scenario: Reject priority value above maximum
- **WHEN** user provides priority=100
- **THEN** system returns validation error

### Requirement: Async operation handling
The system SHALL handle DNS record creation, update, and deletion as async operations by polling the read operation until the changes are reflected.

#### Scenario: Wait for record creation to complete
- **WHEN** system creates a DNS record
- **THEN** system polls read operation until the record appears in the list

#### Scenario: Wait for record update to complete
- **WHEN** system updates a DNS record
- **THEN** system polls read operation until the updated values are reflected

#### Scenario: Wait for record deletion to complete
- **WHEN** system deletes a DNS record
- **THEN** system polls read operation until the record is no longer found

#### Scenario: Timeout after maximum retries
- **WHEN** system exceeds maximum retry count waiting for operation
- **THEN** system returns timeout error

### Requirement: Resource ID format
The system SHALL use a composite resource ID format "zone_id#record_id" to uniquely identify teo DNS records.

#### Scenario: Parse resource ID correctly
- **WHEN** system parses resource_id "zone123#record456"
- **THEN** system extracts zone_id="zone123" and record_id="record456"

#### Scenario: Generate resource ID on creation
- **WHEN** system creates a DNS record with zone_id="zone123"
- **THEN** system generates resource_id="zone123#<returned_record_id>"

### Requirement: Timeout configuration
The system SHALL support configurable timeouts for create, update, and delete operations through Terraform's Timeout block.

#### Scenario: Use default timeout
- **WHEN** user does not specify timeout
- **THEN** system uses default timeout values

#### Scenario: Use custom create timeout
- **WHEN** user specifies custom create timeout
- **THEN** system uses the specified timeout for create operation

#### Scenario: Use custom update timeout
- **WHEN** user specifies custom update timeout
- **THEN** system uses the specified timeout for update operation

#### Scenario: Use custom delete timeout
- **WHEN** user specifies custom delete timeout
- **THEN** system uses the specified timeout for delete operation

### Requirement: Error handling
The system SHALL provide clear error messages for all failure scenarios including invalid parameters, API errors, and timeout errors.

#### Scenario: Return clear error for missing zone_id
- **WHEN** user omits zone_id
- **THEN** system returns error message indicating zone_id is required

#### Scenario: Return clear error for API failure
- **WHEN** cloud API returns error
- **THEN** system propagates the API error message to user

#### Scenario: Return clear error for timeout
- **WHEN** operation times out
- **THEN** system returns error message indicating timeout and current state
