## ADDED Requirements

### Requirement: Create DNS record
The system SHALL allow users to create a TEO DNS record with the specified parameters.

#### Scenario: Successfully create an AAAA DNS record
- **WHEN** user creates a DNS record with zone_id, name, type="AAAA", content, and optional parameters
- **THEN** system creates the DNS record and returns a record_id
- **AND** the record is accessible via the record_id

#### Scenario: Successfully create a DNS record with location and weight
- **WHEN** user creates an A/AAAA/CNAME type DNS record with location and weight parameters
- **THEN** system creates the DNS record with the specified location and weight
- **AND** the record is accessible via the record_id

#### Scenario: Successfully create an MX DNS record with priority
- **WHEN** user creates an MX type DNS record with priority parameter
- **THEN** system creates the DNS record with the specified priority
- **AND** the record is accessible via the record_id

### Requirement: Read DNS record
The system SHALL allow users to read an existing TEO DNS record by its ID.

#### Scenario: Successfully read a DNS record
- **WHEN** user reads a DNS record by zone_id and record_id
- **THEN** system returns the DNS record with all its attributes including zone_id, record_id, name, type, content, location, ttl, weight, priority, status, created_on

### Requirement: Update DNS record
The system SHALL allow users to update an existing TEO DNS record.

#### Scenario: Successfully update DNS record content
- **WHEN** user updates a DNS record's content
- **THEN** system updates the DNS record with the new content
- **AND** subsequent reads return the updated content

#### Scenario: Successfully update DNS record TTL
- **WHEN** user updates a DNS record's TTL
- **THEN** system updates the DNS record with the new TTL
- **AND** subsequent reads return the updated TTL

#### Scenario: Successfully update DNS record weight
- **WHEN** user updates an A/AAAA/CNAME type DNS record's weight
- **THEN** system updates the DNS record with the new weight
- **AND** subsequent reads return the updated weight

#### Scenario: Successfully update MX record priority
- **WHEN** user updates an MX type DNS record's priority
- **THEN** system updates the DNS record with the new priority
- **AND** subsequent reads return the updated priority

### Requirement: Delete DNS record
The system SHALL allow users to delete an existing TEO DNS record.

#### Scenario: Successfully delete a DNS record
- **WHEN** user deletes a DNS record by zone_id and record_id
- **THEN** system deletes the DNS record
- **AND** subsequent read operations return "not found" error

### Requirement: DNS record attributes validation
The system SHALL validate DNS record attributes according to TEO API requirements.

#### Scenario: Validate TTL range
- **WHEN** user creates or updates a DNS record with TTL
- **THEN** system accepts TTL values between 60 and 86400
- **AND** rejects TTL values outside this range

#### Scenario: Validate weight range
- **WHEN** user creates or updates an A/AAAA/CNAME type DNS record with weight
- **THEN** system accepts weight values between -1 and 100
- **AND** rejects weight values outside this range

#### Scenario: Validate MX priority range
- **WHEN** user creates or updates an MX type DNS record with priority
- **THEN** system accepts priority values between 0 and 50
- **AND** rejects priority values outside this range

#### Scenario: Validate location only for specific record types
- **WHEN** user creates or updates a DNS record with location
- **THEN** system accepts location only for A, AAAA, CNAME record types
- **AND** ignores or rejects location for other record types

### Requirement: DNS record type support
The system SHALL support all TEO DNS record types.

#### Scenario: Support A record type
- **WHEN** user creates an A type DNS record with IPv4 address content
- **THEN** system creates the DNS record successfully

#### Scenario: Support AAAA record type
- **WHEN** user creates an AAAA type DNS record with IPv6 address content
- **THEN** system creates the DNS record successfully

#### Scenario: Support CNAME record type
- **WHEN** user creates a CNAME type DNS record with domain name content
- **THEN** system creates the DNS record successfully

#### Scenario: Support TXT record type
- **WHEN** user creates a TXT type DNS record with text content
- **THEN** system creates the DNS record successfully

#### Scenario: Support MX record type
- **WHEN** user creates an MX type DNS record with mail server content and priority
- **THEN** system creates the DNS record successfully

#### Scenario: Support NS record type
- **WHEN** user creates an NS type DNS record with nameserver content
- **THEN** system creates the DNS record successfully for subdomains
- **AND** rejects NS records for root domain

#### Scenario: Support CAA record type
- **WHEN** user creates a CAA type DNS record with CA content
- **THEN** system creates the DNS record successfully

#### Scenario: Support SRV record type
- **WHEN** user creates an SRV type DNS record with service content
- **THEN** system creates the DNS record successfully

### Requirement: DNS record ID format
The system SHALL use a composite ID format for DNS records.

#### Scenario: Resource ID format
- **WHEN** a DNS record is created
- **THEN** system assigns a composite ID in format "zoneId#recordId"
- **AND** the ID can be parsed to extract zone_id and record_id
