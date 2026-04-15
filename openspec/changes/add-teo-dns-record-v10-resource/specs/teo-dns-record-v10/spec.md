## ADDED Requirements

### Requirement: Resource creation
The system SHALL allow users to create a TEO DNS record by providing zone_id, name, type, and content parameters.

#### Scenario: Successful DNS record creation
- **WHEN** user provides valid zone_id, name, type, and content parameters
- **THEN** system creates the DNS record and returns a record_id
- **AND** system waits for the DNS record to become active
- **AND** Terraform state stores the resource with id "zone_id#record_id"

#### Scenario: DNS record creation with optional parameters
- **WHEN** user provides zone_id, name, type, content, and optional parameters (location, ttl, weight, priority)
- **THEN** system creates the DNS record with all specified parameters
- **AND** optional parameters are set according to user input
- **AND** default values are used for unspecified optional parameters (TTL: 300, Weight: -1, Priority: 0)

#### Scenario: DNS record creation failure due to invalid zone_id
- **WHEN** user provides an invalid zone_id
- **THEN** system returns an error indicating the zone does not exist
- **AND** Terraform resource creation fails

#### Scenario: DNS record creation with invalid record type
- **WHEN** user provides an invalid DNS record type
- **THEN** system returns an error indicating invalid type
- **AND** Terraform resource creation fails

### Requirement: Resource read
The system SHALL allow users to read the details of an existing TEO DNS record by its resource id.

#### Scenario: Successful DNS record read
- **WHEN** user queries an existing DNS record by its resource id (zone_id#record_id)
- **THEN** system returns all DNS record attributes including computed fields (status, created_on, modified_on)
- **AND** Terraform state is updated with the current resource state

#### Scenario: DNS record not found
- **WHEN** user queries a non-existent DNS record
- **THEN** system returns an error indicating the resource was not found
- **AND** Terraform removes the resource from state

#### Scenario: DNS record read with computed fields
- **WHEN** system reads an existing DNS record
- **THEN** system includes all computed fields (status, created_on, modified_on) in the resource state
- **AND** computed fields are read-only and cannot be modified by users

### Requirement: Resource update
The system SHALL allow users to update mutable attributes of an existing TEO DNS record.

#### Scenario: Successful DNS record update with content change
- **WHEN** user updates the content of an existing DNS record
- **THEN** system modifies the DNS record content
- **AND** system waits for the change to take effect
- **AND** Terraform state is updated with the new content

#### Scenario: Successful DNS record update with multiple parameters
- **WHEN** user updates multiple mutable parameters (content, ttl, weight, location)
- **THEN** system modifies all specified parameters
- **AND** system waits for the changes to take effect
- **AND** Terraform state is updated with the new values

#### Scenario: DNS record update with invalid parameter for record type
- **WHEN** user attempts to set a parameter that is not valid for the record type (e.g., location for TXT record)
- **THEN** cloud API validates and rejects the update
- **AND** system returns an error from the cloud API
- **AND** Terraform update operation fails

#### Scenario: DNS record update with immutable parameters
- **WHEN** user attempts to modify immutable parameters (name, type, zone_id)
- **THEN** system returns an error indicating these parameters cannot be changed
- **AND** user must recreate the resource with new values

#### Scenario: DNS record update failure due to record not found
- **WHEN** user attempts to update a non-existent DNS record
- **THEN** system returns an error indicating the resource was not found
- **AND** Terraform marks the resource as tainted

### Requirement: Resource deletion
The system SHALL allow users to delete an existing TEO DNS record.

#### Scenario: Successful DNS record deletion
- **WHEN** user deletes an existing DNS record
- **THEN** system removes the DNS record from the cloud service
- **AND** Terraform removes the resource from state
- **AND** resource is no longer managed by Terraform

#### Scenario: DNS record deletion failure due to record not found
- **WHEN** user attempts to delete a non-existent DNS record
- **THEN** system returns an error indicating the resource was not found
- **AND** Terraform removes the resource from state (idempotent deletion)

#### Scenario: DNS record deletion with dependencies
- **WHEN** user attempts to delete a DNS record that is referenced by other resources
- **THEN** system returns an error from the cloud API
- **AND** Terraform deletion operation fails
- **AND** user must resolve dependencies before deletion

### Requirement: Resource import
The system SHALL allow users to import an existing TEO DNS record into Terraform state management.

#### Scenario: Successful DNS record import
- **WHEN** user imports a DNS record using the command with id "zone_id#record_id"
- **THEN** system queries the DNS record from the cloud service
- **AND** Terraform creates a resource state with the current configuration
- **AND** resource is now managed by Terraform

#### Scenario: DNS record import failure due to invalid id format
- **WHEN** user attempts to import with an invalid id format
- **THEN** system returns an error indicating the id format is incorrect
- **AND** Terraform import operation fails

#### Scenario: DNS record import failure due to record not found
- **WHEN** user attempts to import a non-existent DNS record
- **THEN** system returns an error indicating the resource was not found
- **AND** Terraform import operation fails

### Requirement: DNS record type support
The system SHALL support all TEO DNS record types: A, AAAA, MX, CNAME, TXT, NS, CAA, SRV.

#### Scenario: A record creation
- **WHEN** user creates an A record with a valid IPv4 address
- **THEN** system creates the A record successfully
- **AND** record type is set to "A"
- **AND** content is validated as a valid IPv4 address by the cloud API

#### Scenario: AAAA record creation
- **WHEN** user creates an AAAA record with a valid IPv6 address
- **THEN** system creates the AAAA record successfully
- **AND** record type is set to "AAAA"
- **AND** content is validated as a valid IPv6 address by the cloud API

#### Scenario: CNAME record creation
- **WHEN** user creates a CNAME record with a domain name
- **THEN** system creates the CNAME record successfully
- **AND** record type is set to "CNAME"
- **AND** content is set to the provided domain name

#### Scenario: MX record creation with priority
- **WHEN** user creates an MX record with a mail server and priority value
- **THEN** system creates the MX record successfully
- **AND** record type is set to "MX"
- **AND** priority is set to the provided value (0-50)
- **AND** default priority is 0 if not specified

#### Scenario: TXT record creation
- **WHEN** user creates a TXT record with text content
- **THEN** system creates the TXT record successfully
- **AND** record type is set to "TXT"
- **AND** content is set to the provided text

#### Scenario: NS record creation
- **WHEN** user creates an NS record with a name server
- **THEN** system creates the NS record successfully
- **AND** record type is set to "NS"
- **AND** content is set to the provided name server
- **AND** root domain cannot have NS records

#### Scenario: CAA record creation
- **WHEN** user creates a CAA record with CAA parameters
- **THEN** system creates the CAA record successfully
- **AND** record type is set to "CAA"
- **AND** content is set to the provided CAA record value

#### Scenario: SRV record creation
- **WHEN** user creates an SRV record with SRV parameters
- **THEN** system creates the SRV record successfully
- **AND** record type is set to "SRV"
- **AND** content is set to the provided SRV record value

### Requirement: DNS record parameters
The system SHALL support all TEO DNS record parameters according to the cloud API specification.

#### Scenario: TTL parameter validation
- **WHEN** user provides a TTL value within the valid range (60-86400 seconds)
- **THEN** system creates the DNS record with the specified TTL
- **AND** default TTL is 300 if not specified

#### Scenario: TTL parameter out of range
- **WHEN** user provides a TTL value outside the valid range (< 60 or > 86400)
- **THEN** cloud API rejects the request
- **AND** system returns an error indicating invalid TTL value

#### Scenario: Weight parameter for A/AAAA/CNAME records
- **WHEN** user provides a weight value (-1 to 100) for A, AAAA, or CNAME record
- **THEN** system creates the DNS record with the specified weight
- **AND** -1 means no weight assigned
- **AND** 0 means the record does not resolve
- **AND** default weight is -1 if not specified

#### Scenario: Weight parameter for other record types
- **WHEN** user attempts to set weight for non-A/AAAA/CNAME record types
- **THEN** cloud API validates and rejects the parameter
- **AND** system returns an error from the cloud API

#### Scenario: Location parameter for A/AAAA/CNAME records
- **WHEN** user provides a location value for A, AAAA, or CNAME record
- **THEN** system creates the DNS record with the specified location
- **AND** default location is "Default" if not specified
- **AND** location parameter is only available for standard and enterprise editions

#### Scenario: Location parameter for other record types
- **WHEN** user attempts to set location for non-A/AAAA/CNAME record types
- **THEN** cloud API validates and rejects the parameter
- **AND** system returns an error from the cloud API

#### Scenario: Priority parameter for MX records
- **WHEN** user provides a priority value (0-50) for MX record
- **THEN** system creates the MX record with the specified priority
- **AND** lower values indicate higher priority
- **AND** default priority is 0 if not specified

#### Scenario: Priority parameter for non-MX record types
- **WHEN** user attempts to set priority for non-MX record types
- **THEN** cloud API validates and ignores the parameter
- **AND** parameter has no effect on non-MX records

### Requirement: Resource timeouts
The system SHALL support configurable timeouts for resource operations.

#### Scenario: Default timeouts
- **WHEN** user creates the DNS record without specifying timeouts
- **THEN** system uses default timeout values
- **AND** create timeout is 20 minutes
- **AND** read timeout is 3 minutes
- **AND** update timeout is 20 minutes
- **AND** delete timeout is 20 minutes

#### Scenario: Custom timeouts
- **WHEN** user specifies custom timeout values
- **THEN** system uses the specified timeout values for each operation
- **AND** operations fail if they exceed the specified timeout

#### Scenario: Timeout during record creation
- **WHEN** DNS record creation takes longer than the specified timeout
- **THEN** system returns a timeout error
- **AND** Terraform marks the resource as tainted
- **AND** user can retry the operation

### Requirement: Resource state management
The system SHALL maintain accurate resource state throughout its lifecycle.

#### Scenario: State synchronization after creation
- **WHEN** DNS record creation completes successfully
- **THEN** system performs a read operation to fetch the final state
- **AND** Terraform state is synchronized with the actual resource state
- **AND** all computed fields are populated in the state

#### Scenario: State synchronization after update
- **WHEN** DNS record update completes successfully
- **THEN** system performs a read operation to fetch the final state
- **AND** Terraform state is synchronized with the actual resource state
- **AND** all modified fields are updated in the state

#### Scenario: State drift detection
- **WHEN** actual resource state differs from Terraform state
- **THEN** system detects the drift during the next Terraform operation
- **AND** system proposes a plan to reconcile the state

### Requirement: Error handling
The system SHALL provide clear error messages for all failure scenarios.

#### Scenario: API error handling
- **WHEN** cloud API returns an error
- **THEN** system propagates the error to Terraform
- **AND** error message includes relevant details from the API response
- **AND** user can understand the cause of the failure

#### Scenario: Network error handling
- **WHEN** network error occurs during API call
- **THEN** system retries the operation up to the retry limit
- **AND** if retries are exhausted, system returns an error
- **AND** error message indicates a network error occurred

#### Scenario: Validation error handling
- **WHEN** user input fails validation
- **THEN** system returns a clear error message
- **AND** error message indicates which parameter is invalid and why

### Requirement: Resource documentation
The system SHALL provide comprehensive documentation for the DNS record resource.

#### Scenario: Documentation completeness
- **WHEN** user views the resource documentation
- **THEN** documentation includes all resource parameters
- **AND** documentation includes parameter types and constraints
- **AND** documentation includes usage examples
- **AND** documentation includes information about DNS record types

#### Scenario: Documentation accuracy
- **WHEN** resource implementation changes
- **THEN** documentation is updated to reflect the changes
- **AND** documentation remains in sync with the actual implementation
