## ADDED Requirements

### Requirement: Create DNS record
The system SHALL allow users to create a TEO DNS record with specified parameters.

#### Scenario: Successful creation of A record
- **WHEN** user creates a DNS record with Type="A", Name="www", Content="1.2.3.4", ZoneId="zone-123"
- **THEN** system calls CreateDnsRecord API with provided parameters
- **AND** system returns record ID
- **AND** system stores record state with composite ID "zone-123#recordId"

#### Scenario: Successful creation of MX record with priority
- **WHEN** user creates a DNS record with Type="MX", Name="@", Content="mail.example.com", Priority=10, ZoneId="zone-123"
- **THEN** system calls CreateDnsRecord API with provided parameters including Priority
- **AND** system returns record ID
- **AND** system stores record state

#### Scenario: Create record with optional parameters
- **WHEN** user creates a DNS record with Location="China", TTL=600, Weight=50
- **THEN** system calls CreateDnsRecord API with all specified parameters
- **AND** system stores all parameters in state

#### Scenario: Idempotent creation
- **WHEN** user creates a DNS record with same ZoneId, Name, Type, and Content that already exists
- **THEN** system checks existing records using DescribeDnsRecords
- **AND** system returns existing record ID without creating duplicate
- **AND** system maintains idempotency

### Requirement: Read DNS record
The system SHALL allow users to read a TEO DNS record by its ID.

#### Scenario: Successful read
- **WHEN** user reads a DNS record with ID "zone-123#recordId"
- **THEN** system extracts ZoneId="zone-123" and RecordId="recordId"
- **AND** system calls DescribeDnsRecords API with ZoneId and filter by RecordId
- **AND** system returns record details including all attributes
- **AND** system updates Terraform state with current record data

#### Scenario: Read non-existent record
- **WHEN** user reads a DNS record that does not exist
- **THEN** system calls DescribeDnsRecords API
- **AND** system returns empty result
- **AND** system removes record from Terraform state

#### Scenario: Read with all record types
- **WHEN** user reads DNS records of types A, AAAA, MX, CNAME, TXT, NS, CAA, SRV
- **THEN** system retrieves records for all specified types
- **AND** system correctly parses Type attribute for each record

### Requirement: Update DNS record
The system SHALL allow users to update a TEO DNS record.

#### Scenario: Successful update of record content
- **WHEN** user updates a DNS record with ID "zone-123#recordId" changing Content from "1.2.3.4" to "5.6.7.8"
- **THEN** system extracts ZoneId and RecordId from ID
- **AND** system calls ModifyDnsRecords API with RecordId and new Content
- **AND** system waits for update to complete by polling DescribeDnsRecords
- **AND** system updates Terraform state with new Content value

#### Scenario: Update multiple parameters
- **WHEN** user updates a DNS record changing Name, Content, Location, TTL, and Weight
- **THEN** system calls ModifyDnsRecords API with all updated parameters
- **AND** system updates all specified parameters in Terraform state

#### Scenario: Update MX record priority
- **WHEN** user updates an MX record changing Priority from 10 to 20
- **THEN** system calls ModifyDnsRecords API with new Priority value
- **AND** system updates Priority in Terraform state

#### Scenario: No-op update
- **WHEN** user updates a DNS record with same values as current state
- **THEN** system detects no changes
- **AND** system does not call ModifyDnsRecords API
- **AND** system maintains existing state

### Requirement: Delete DNS record
The system SHALL allow users to delete a TEO DNS record.

#### Scenario: Successful deletion
- **WHEN** user deletes a DNS record with ID "zone-123#recordId"
- **THEN** system extracts ZoneId="zone-123" and RecordId="recordId"
- **AND** system calls DeleteDnsRecords API with ZoneId and RecordIds
- **AND** system removes record from Terraform state

#### Scenario: Delete non-existent record
- **WHEN** user deletes a DNS record that does not exist
- **THEN** system calls DeleteDnsRecords API
- **AND** system handles API response appropriately
- **AND** system removes record from Terraform state

### Requirement: Handle record types and their constraints
The system SHALL support all DNS record types with their specific constraints.

#### Scenario: A record with IPv4 address
- **WHEN** user creates an A record with IPv4 address like "1.2.3.4"
- **THEN** system validates and accepts the IPv4 address format

#### Scenario: AAAA record with IPv6 address
- **WHEN** user creates an AAAA record with IPv6 address like "2001:0db8:85a3:0000:0000:8a2e:0370:7334"
- **THEN** system validates and accepts the IPv6 address format

#### Scenario: CNAME record with domain
- **WHEN** user creates a CNAME record pointing to another domain
- **THEN** system accepts domain name as Content value

#### Scenario: MX record with priority constraint
- **WHEN** user creates an MX record with Priority outside range 0-50
- **THEN** system validates Priority is within acceptable range
- **AND** system returns validation error if out of range

#### Scenario: Weight only applicable to certain types
- **WHEN** user creates a DNS record with Type="TXT" and Weight=50
- **THEN** system accepts Weight parameter (validation handled by API)
- **AND** system stores Weight value in state

### Requirement: Handle TTL parameter
The system SHALL support TTL parameter with appropriate constraints.

#### Scenario: Create record with TTL
- **WHEN** user creates a DNS record with TTL=300
- **THEN** system calls CreateDnsRecord API with TTL value
- **AND** system stores TTL in Terraform state

#### Scenario: Update TTL
- **WHEN** user updates a DNS record changing TTL from 300 to 600
- **THEN** system calls ModifyDnsRecords API with new TTL value
- **AND** system updates TTL in Terraform state

#### Scenario: Default TTL when not specified
- **WHEN** user creates a DNS record without specifying TTL
- **THEN** system omits TTL from API call
- **AND** system uses API default value (typically 300)
- **AND** system stores actual TTL value returned by API

### Requirement: Handle Location parameter for routing
The system SHALL support Location parameter for DNS record routing.

#### Scenario: Create record with Default location
- **WHEN** user creates a DNS record without specifying Location
- **THEN** system omits Location from API call
- **AND** system uses API default value (typically "Default")

#### Scenario: Create record with specific location
- **WHEN** user creates a DNS record with Location="China"
- **THEN** system calls CreateDnsRecord API with Location value
- **AND** system stores Location in Terraform state

#### Scenario: Update location
- **WHEN** user updates a DNS record changing Location from "Default" to "China"
- **THEN** system calls ModifyDnsRecords API with new Location value
- **AND** system updates Location in Terraform state

### Requirement: Handle Weight parameter for load balancing
The system SHALL support Weight parameter for DNS record load balancing.

#### Scenario: Create record with weight
- **WHEN** user creates a DNS record with Weight=50
- **THEN** system calls CreateDnsRecord API with Weight value
- **AND** system stores Weight in Terraform state

#### Scenario: Create record without weight
- **WHEN** user creates a DNS record without specifying Weight
- **THEN** system omits Weight from API call
- **AND** system uses API default value (typically -1, meaning no weight)
- **AND** system stores actual Weight value returned by API

#### Scenario: Weight value of 0 disables resolution
- **WHEN** user creates a DNS record with Weight=0
- **THEN** system calls CreateDnsRecord API with Weight=0
- **AND** system stores Weight=0 in Terraform state
- **AND** record will not resolve (per API behavior)

### Requirement: Handle Status attribute
The system SHALL correctly handle the Status attribute as a read-only computed field.

#### Scenario: Read record with enabled status
- **WHEN** system reads a DNS record that is active
- **THEN** system returns Status="enable"
- **AND** system stores Status as computed field in state

#### Scenario: Read record with disabled status
- **WHEN** system reads a DNS record that is inactive
- **THEN** system returns Status="disable"
- **AND** system stores Status as computed field in state

#### Scenario: Status cannot be modified
- **WHEN** user attempts to update a DNS record's Status
- **THEN** system ignores Status in Update request
- **AND** system does not include Status in ModifyDnsRecords API call

### Requirement: Handle CreatedOn timestamp
The system SHALL correctly handle the CreatedOn attribute as a read-only computed field.

#### Scenario: Read record creation timestamp
- **WHEN** system reads a DNS record
- **THEN** system returns CreatedOn timestamp from API
- **AND** system stores CreatedOn as computed field in state

#### Scenario: CreatedOn cannot be modified
- **WHEN** user attempts to update a DNS record's CreatedOn
- **THEN** system ignores CreatedOn in Update request
- **AND** system does not include CreatedOn in ModifyDnsRecords API call

### Requirement: Handle asynchronous operations
The system SHALL properly handle asynchronous API operations with polling.

#### Scenario: Wait for record creation to complete
- **WHEN** system creates a DNS record
- **THEN** system calls CreateDnsRecord API
- **AND** system polls DescribeDnsRecords API until record is available
- **AND** system returns success only after record is confirmed created

#### Scenario: Wait for record update to complete
- **WHEN** system updates a DNS record
- **THEN** system calls ModifyDnsRecords API
- **AND** system polls DescribeDnsRecords API until changes are reflected
- **AND** system returns success only after update is confirmed

#### Scenario: Handle polling timeout
- **WHEN** system polls for record status but timeout is reached
- **THEN** system returns timeout error
- **AND** system includes timeout in error message
- **AND** user can configure timeout via Timeouts block

### Requirement: Handle Timeouts configuration
The system SHALL support user-configurable Timeouts for async operations.

#### Scenario: Use default timeouts
- **WHEN** user does not specify Timeouts in resource configuration
- **THEN** system uses default timeout values for create, update, delete, read operations

#### Scenario: Use custom create timeout
- **WHEN** user specifies Timeouts.create=10m
- **THEN** system uses 10 minute timeout for create operations
- **AND** system polls DescribeDnsRecords API for up to 10 minutes

#### Scenario: Use custom update timeout
- **WHEN** user specifies Timeouts.update=5m
- **THEN** system uses 5 minute timeout for update operations
- **AND** system polls DescribeDnsRecords API for up to 5 minutes

### Requirement: Handle error conditions
The system SHALL properly handle and report error conditions.

#### Scenario: Invalid ZoneId
- **WHEN** user creates a DNS record with invalid ZoneId
- **THEN** system receives error from API
- **AND** system returns clear error message indicating invalid ZoneId

#### Scenario: Invalid record type
- **WHEN** user creates a DNS record with unsupported Type
- **THEN** system receives error from API
- **AND** system returns clear error message indicating invalid Type

#### Scenario: Network timeout
- **WHEN** API call times out due to network issues
- **THEN** system returns timeout error
- **AND** system includes retry information if applicable

#### Scenario: Authentication failure
- **WHEN** API call fails due to authentication issues
- **THEN** system returns authentication error
- **AND** system includes credential configuration guidance
