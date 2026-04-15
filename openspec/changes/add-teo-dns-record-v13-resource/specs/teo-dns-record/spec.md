## ADDED Requirements

### Requirement: Create DNS record
The system SHALL allow users to create a TEO DNS record with specified configuration including zone ID, record name, type, content, and optional parameters such as TTL, location, weight, and priority.

#### Scenario: Create A record with basic configuration
- **WHEN** user provides zone_id, name="www", type="A", content="1.2.3.4"
- **THEN** system creates an A DNS record and returns record_id

#### Scenario: Create CNAME record with TTL and location
- **WHEN** user provides zone_id, name="api", type="CNAME", content="example.com", ttl=600, location="Mainland"
- **THEN** system creates a CNAME DNS record with specified TTL and location

#### Scenario: Create MX record with priority
- **WHEN** user provides zone_id, name="@", type="MX", content="mail.example.com", priority=10
- **THEN** system creates an MX DNS record with specified priority

#### Scenario: Create AAAA record with weight
- **WHEN** user provides zone_id, name="ipv6", type="AAAA", content="2001:0db8:85a3:0000:0000:8a2e:0370:7334", weight=10
- **THEN** system creates an AAAA DNS record with specified weight

#### Scenario: Create TXT record
- **WHEN** user provides zone_id, name="_dmarc", type="TXT", content="v=DMARC1; p=none"
- **THEN** system creates a TXT DNS record

#### Scenario: Create SRV record
- **WHEN** user provides zone_id, name="_sip._tcp", type="SRV", content="10 60 5060 sipserver.example.com"
- **THEN** system creates an SRV DNS record

#### Scenario: Create CAA record
- **WHEN** user provides zone_id, name="@", type="CAA", content="0 issue \"letsencrypt.org\""
- **THEN** system creates a CAA DNS record

#### Scenario: Create NS record for subdomain
- **WHEN** user provides zone_id, name="sub", type="NS", content="ns.example.com"
- **THEN** system creates an NS DNS record

### Requirement: Read DNS record by ID
The system SHALL allow users to retrieve DNS record details including all configuration parameters and metadata such as status, creation time, and modification time.

#### Scenario: Read existing DNS record
- **WHEN** user queries by record_id
- **THEN** system returns DNS record details including zone_id, record_id, name, type, content, location, ttl, weight, priority, status, created_on, modified_on

#### Scenario: Query DNS records by zone with filters
- **WHEN** user queries by zone_id with filters for record name
- **THEN** system returns list of DNS records matching the criteria

#### Scenario: Query DNS records with sorting
- **WHEN** user queries by zone_id with sort_by="name" and sort_order="asc"
- **THEN** system returns DNS records sorted by name in ascending order

### Requirement: Update DNS record configuration
The system SHALL allow users to modify DNS record parameters including content, TTL, location, weight, and priority.

#### Scenario: Update DNS record content
- **WHEN** user provides record_id and new content="5.6.7.8"
- **THEN** system updates the DNS record content

#### Scenario: Update DNS record TTL
- **WHEN** user provides record_id and new ttl=300
- **THEN** system updates the DNS record TTL

#### Scenario: Update DNS record location
- **WHEN** user provides record_id and new location="Overseas"
- **THEN** system updates the DNS record location

#### Scenario: Update DNS record weight
- **WHEN** user provides record_id and new weight=20
- **THEN** system updates the DNS record weight

#### Scenario: Update DNS record priority
- **WHEN** user provides record_id and new priority=5
- **THEN** system updates the DNS record priority

#### Scenario: Update multiple DNS record parameters
- **WHEN** user provides record_id and updates content, ttl, and weight simultaneously
- **THEN** system updates all specified DNS record parameters

### Requirement: Delete DNS record
The system SHALL allow users to delete DNS records by record ID.

#### Scenario: Delete single DNS record
- **WHEN** user provides record_id
- **THEN** system deletes the DNS record

#### Scenario: Delete multiple DNS records
- **WHEN** user provides multiple record_ids
- **THEN** system deletes all specified DNS records

### Requirement: DNS record validation
The system SHALL validate DNS record parameters according to record type specifications and cloud API constraints.

#### Scenario: Validate record type constraints
- **WHEN** user attempts to create a record with invalid type
- **THEN** system rejects the request with appropriate error message

#### Scenario: Validate TTL range
- **WHEN** user provides TTL value outside valid range (60-86400)
- **THEN** system rejects the request with appropriate error message

#### Scenario: Validate weight range
- **WHEN** user provides weight value outside valid range (-1 to 100)
- **THEN** system rejects the request with appropriate error message

#### Scenario: Validate priority range
- **WHEN** user provides priority value outside valid range (0-50)
- **THEN** system rejects the request with appropriate error message

#### Scenario: Validate priority only for MX records
- **WHEN** user provides priority for non-MX record type
- **THEN** system ignores the priority parameter

### Requirement: DNS record idempotency
The system SHALL ensure that creating, updating, or deleting DNS records are idempotent operations.

#### Scenario: Recreate DNS record with same parameters
- **WHEN** user attempts to create a DNS record with identical parameters that already exists
- **THEN** system returns existing record without duplication

#### Scenario: Update DNS record with same values
- **WHEN** user updates a DNS record with unchanged parameter values
- **THEN** system performs no-op update without errors

#### Scenario: Delete non-existent DNS record
- **WHEN** user attempts to delete a DNS record that does not exist
- **THEN** system returns success without errors

### Requirement: Async operation handling
The system SHALL handle asynchronous DNS record operations by polling the DescribeDnsRecords API until the operation completes.

#### Scenario: Wait for DNS record creation to complete
- **WHEN** user creates a DNS record
- **THEN** system polls until record status is "enable" before returning

#### Scenario: Wait for DNS record update to complete
- **WHEN** user updates a DNS record
- **THEN** system polls until record is updated with new values before returning

#### Scenario: Wait for DNS record deletion to complete
- **WHEN** user deletes a DNS record
- **THEN** system polls until record no longer exists before returning success

### Requirement: Terraform state management
The system SHALL maintain Terraform state consistency with actual DNS record state, including all resource attributes.

#### Scenario: Import existing DNS record into Terraform state
- **WHEN** user imports an existing DNS record by record_id
- **THEN** system populates Terraform state with all current DNS record attributes

#### Scenario: Refresh Terraform state with latest DNS record
- **WHEN** user runs terraform refresh
- **THEN** system updates Terraform state with latest DNS record data

#### Scenario: Handle DNS record state drift
- **WHEN** actual DNS record differs from Terraform state
- **THEN** system detects and reports state drift on next plan

### Requirement: DNS record status management
The system SHALL support managing DNS record status (enable/disable) through the resource lifecycle.

#### Scenario: Disable DNS record
- **WHEN** user disables a DNS record
- **THEN** system sets record status to "disable"

#### Scenario: Enable DNS record
- **WHEN** user enables a disabled DNS record
- **THEN** system sets record status to "enable"

#### Scenario: Read DNS record includes current status
- **WHEN** user reads a DNS record
- **THEN** system returns current status (enable/disable) in response

### Requirement: Error handling and retries
The system SHALL handle API errors gracefully and implement retry logic for transient failures.

#### Scenario: Retry on temporary API failure
- **WHEN** cloud API returns temporary error during operation
- **THEN** system retries the operation with exponential backoff

#### Scenario: Handle authentication errors
- **WHEN** cloud API returns authentication error
- **THEN** system returns clear error message to user without retrying

#### Scenario: Handle rate limiting
- **WHEN** cloud API returns rate limit error
- **THEN** system waits and retries operation after appropriate delay

### Requirement: Timeouts configuration
The system SHALL allow users to configure timeout values for create, update, and delete operations.

#### Scenario: Configure custom create timeout
- **WHEN** user specifies create timeout of 10 minutes
- **THEN** system uses specified timeout for create operation

#### Scenario: Configure custom update timeout
- **WHEN** user specifies update timeout of 5 minutes
- **THEN** system uses specified timeout for update operation

#### Scenario: Configure custom delete timeout
- **WHEN** user specifies delete timeout of 5 minutes
- **THEN** system uses specified timeout for delete operation

#### Scenario: Use default timeout when not specified
- **WHEN** user does not specify timeout values
- **THEN** system uses default timeout values for all operations
