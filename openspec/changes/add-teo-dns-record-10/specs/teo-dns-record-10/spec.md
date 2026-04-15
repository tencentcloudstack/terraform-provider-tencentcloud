## ADDED Requirements

### Requirement: Resource identifier

The resource SHALL be identified by a composite ID in the format `zone_id#record_id`, where:
- `zone_id`: The zone ID from the cloud API
- `record_id`: The DNS record ID returned by the CreateDnsRecord API

This composite ID uniquely identifies a DNS record and is used throughout CRUD operations.

#### Scenario: Parse composite ID
- **WHEN** the Terraform resource receives an ID in the format `zone_id#record_id`
- **THEN** the provider SHALL parse the ID into zone_id and record_id components
- **AND** use the zone_id to identify the zone
- **AND** use the record_id to identify the specific DNS record

#### Scenario: Build composite ID after creation
- **WHEN** the CreateDnsRecord API returns a successful response
- **THEN** the provider SHALL extract the zone_id from the request
- **AND** extract the record_id from the response
- **AND** build a composite ID in the format `zone_id#record_id`
- **AND** set this as the Terraform resource ID

### Requirement: Create operation

The provider SHALL support creating a DNS record through the CreateDnsRecord API.

Required parameters:
- `zone_id`: Zone ID (Required)
- `name`: DNS record name (Required)
- `type`: DNS record type - one of A, AAAA, MX, CNAME, TXT, NS, CAA, SRV (Required)
- `content`: DNS record content (Required)

Optional parameters:
- `location`: DNS record resolution line, applicable for A, AAAA, CNAME types only (Optional, default: "Default")
- `ttl`: Cache time in seconds, range 60~86400 (Optional, default: 300)
- `weight`: DNS record weight, range -1~100, applicable for A, AAAA, CNAME types only (Optional, default: -1)
- `priority`: MX record priority, range 0~50, applicable for MX type only (Optional, default: 0)

#### Scenario: Successful creation
- **WHEN** the user provides all required parameters (zone_id, name, type, content)
- **AND** the provider calls the CreateDnsRecord API
- **AND** the API returns successfully with a record_id
- **THEN** the provider SHALL build a composite ID
- **AND** store the resource in the Terraform state
- **AND** return the resource to the user

#### Scenario: Creation with optional parameters
- **WHEN** the user provides optional parameters (location, ttl, weight, priority)
- **AND** the provider calls the CreateDnsRecord API with all parameters
- **THEN** the provider SHALL create the DNS record with the specified configuration
- **AND** store the resource in the Terraform state

#### Scenario: Creation failure - invalid parameters
- **WHEN** the user provides invalid parameters (e.g., type not in allowed values, ttl out of range)
- **AND** the provider calls the CreateDnsRecord API
- **AND** the API returns an error
- **THEN** the provider SHALL return a descriptive error message to the user
- **AND** not store any resource in the Terraform state

#### Scenario: Creation failure - missing required parameters
- **WHEN** the user does not provide one or more required parameters
- **THEN** the provider SHALL return a validation error
- **AND** not call the CreateDnsRecord API
- **AND** not store any resource in the Terraform state

#### Scenario: Async operation handling
- **WHEN** the CreateDnsRecord API is asynchronous
- **AND** the API returns successfully
- **THEN** the provider SHALL call the DescribeDnsRecords API to poll for the record
- **AND** wait until the record is available
- **OR** timeout according to the user-configured timeout value

### Requirement: Read operation

The provider SHALL support reading a DNS record through the DescribeDnsRecords API.

The read operation SHALL:
1. Parse the composite ID to extract zone_id and record_id
2. Call DescribeDnsRecords with filters to find the specific record
3. Map all returned fields to the Terraform schema

Returned attributes:
- `zone_id`: Zone ID (Computed)
- `record_id`: DNS record ID (Computed)
- `name`: DNS record name (Computed)
- `type`: DNS record type (Computed)
- `content`: DNS record content (Computed)
- `location`: DNS record resolution line (Computed)
- `ttl`: Cache time in seconds (Computed)
- `weight`: DNS record weight (Computed)
- `priority`: MX record priority (Computed)
- `status`: DNS record resolution status - "enable" or "disable" (Computed)
- `created_on`: Creation time (Computed)
- `modified_on`: Modification time (Computed)

#### Scenario: Successful read
- **WHEN** the provider parses the composite ID to extract zone_id and record_id
- **AND** calls DescribeDnsRecords with zone_id and filter by record_id
- **AND** the API returns the DNS record
- **THEN** the provider SHALL populate all computed attributes
- **AND** return the resource to the user

#### Scenario: Read after creation
- **WHEN** a resource has been created
- **AND** the provider performs a read operation
- **THEN** the provider SHALL call DescribeDnsRecords with the zone_id and record_id
- **AND** verify that the record exists with the expected attributes
- **AND** update the Terraform state with the actual cloud resource state

#### Scenario: Read after update
- **WHEN** a resource has been updated
- **AND** the provider performs a read operation
- **THEN** the provider SHALL call DescribeDnsRecords with the zone_id and record_id
- **AND** verify that the record has been updated with the new attributes
- **AND** update the Terraform state with the actual cloud resource state

#### Scenario: Record not found
- **WHEN** the provider calls DescribeDnsRecords with the zone_id and record_id
- **AND** the API returns an empty result set
- **THEN** the provider SHALL remove the resource from the Terraform state
- **AND** mark the resource as deleted

#### Scenario: Read error
- **WHEN** the provider calls DescribeDnsRecords
- **AND** the API returns an error
- **THEN** the provider SHALL return a descriptive error message to the user
- **AND** not modify the Terraform state

### Requirement: Update operation

The provider SHALL support updating a DNS record through the ModifyDnsRecords API.

Updatable parameters:
- `name`: DNS record name
- `type`: DNS record type
- `content`: DNS record content
- `location`: DNS record resolution line
- `ttl`: Cache time in seconds
- `weight`: DNS record weight
- `priority`: MX record priority

Non-updatable parameters (computed only):
- `zone_id`: Zone ID
- `record_id`: DNS record ID
- `status`: DNS record resolution status
- `created_on`: Creation time
- `modified_on`: Modification time

#### Scenario: Successful update
- **WHEN** the user modifies one or more updatable parameters
- **AND** the provider calls ModifyDnsRecords with the zone_id and the modified DnsRecord
- **AND** the API returns successfully
- **THEN** the provider SHALL update the Terraform state
- **AND** call DescribeDnsRecords to verify the update
- **AND** return the updated resource to the user

#### Scenario: Update with new type
- **WHEN** the user changes the DNS record type (e.g., from A to CNAME)
- **AND** the provider calls ModifyDnsRecords with the new type
- **AND** the API returns successfully
- **THEN** the provider SHALL update the Terraform state
- **AND** verify that the type has been changed

#### Scenario: Update with invalid parameters
- **WHEN** the user provides invalid parameters for update
- **AND** the provider calls ModifyDnsRecords
- **AND** the API returns an error
- **THEN** the provider SHALL return a descriptive error message to the user
- **AND** not modify the Terraform state

#### Scenario: Update non-updatable parameter
- **WHEN** the user attempts to update a non-updatable parameter (e.g., zone_id, record_id)
- **THEN** the provider SHALL ignore the change
- **AND** not call the ModifyDnsRecords API
- **AND** keep the Terraform state unchanged

#### Scenario: Async operation handling
- **WHEN** the ModifyDnsRecords API is asynchronous
- **AND** the API returns successfully
- **THEN** the provider SHALL call the DescribeDnsRecords API to poll for the record
- **AND** wait until the record is updated
- **OR** timeout according to the user-configured timeout value

### Requirement: Delete operation

The provider SHALL support deleting a DNS record through the DeleteDnsRecords API.

#### Scenario: Successful deletion
- **WHEN** the user deletes the resource
- **AND** the provider parses the composite ID to extract zone_id and record_id
- **AND** calls DeleteDnsRecords with zone_id and record_id
- **AND** the API returns successfully
- **THEN** the provider SHALL remove the resource from the Terraform state
- **AND** return success to the user

#### Scenario: Delete non-existent record
- **WHEN** the user deletes a resource
- **AND** the provider calls DeleteDnsRecords
- **AND** the API returns an error indicating the record does not exist
- **THEN** the provider SHALL remove the resource from the Terraform state
- **AND** return success to the user (idempotent delete)

#### Scenario: Delete error
- **WHEN** the provider calls DeleteDnsRecords
- **AND** the API returns an error
- **THEN** the provider SHALL return a descriptive error message to the user
- **AND** not remove the resource from the Terraform state

#### Scenario: Async operation handling
- **WHEN** the DeleteDnsRecords API is asynchronous
- **AND** the API returns successfully
- **THEN** the provider SHALL call the DescribeDnsRecords API to poll for the record
- **AND** wait until the record is deleted (not found)
- **OR** timeout according to the user-configured timeout value

### Requirement: Timeout configuration

The provider SHALL support timeout configuration for async operations.

The timeout configuration SHALL include:
- `create`: Timeout for create operation (default: 10 minutes)
- `update`: Timeout for update operation (default: 10 minutes)
- `delete`: Timeout for delete operation (default: 10 minutes)
- `read`: Timeout for read operation (default: 5 minutes)

#### Scenario: Custom timeout
- **WHEN** the user configures a custom timeout value
- **THEN** the provider SHALL use the custom timeout value for the corresponding operation
- **AND** return a timeout error if the operation exceeds the configured timeout

#### Scenario: Default timeout
- **WHEN** the user does not configure a custom timeout value
- **THEN** the provider SHALL use the default timeout value for the corresponding operation

#### Scenario: Timeout during async operation
- **WHEN** an async operation exceeds the configured timeout
- **THEN** the provider SHALL return a timeout error to the user
- **AND** include details about the operation that timed out
- **AND** suggest adjusting the timeout configuration

### Requirement: Error handling and retry

The provider SHALL implement proper error handling and retry logic for API calls.

#### Scenario: Retryable error
- **WHEN** the API returns a retryable error (e.g., rate limit, temporary network issue)
- **THEN** the provider SHALL retry the operation with exponential backoff
- **AND** not exceed the maximum retry limit
- **AND** return a descriptive error if all retries fail

#### Scenario: Non-retryable error
- **WHEN** the API returns a non-retryable error (e.g., invalid parameters, permission denied)
- **THEN** the provider SHALL not retry the operation
- **AND** return a descriptive error to the user immediately

#### Scenario: Final consistency check
- **WHEN** the provider completes an operation
- **THEN** the provider SHALL perform a final consistency check
- **AND** log the operation duration
- **AND** detect and report any inconsistent state

### Requirement: Schema validation

The provider SHALL validate schema parameters before making API calls.

#### Scenario: Type validation
- **WHEN** the user provides an invalid DNS record type
- **THEN** the provider SHALL return a validation error
- **AND** list the allowed values: A, AAAA, MX, CNAME, TXT, NS, CAA, SRV

#### Scenario: TTL range validation
- **WHEN** the user provides a TTL value outside the range 60~86400
- **THEN** the provider SHALL return a validation error
- **AND** indicate the valid range

#### Scenario: Weight range validation
- **WHEN** the user provides a weight value outside the range -1~100
- **THEN** the provider SHALL return a validation error
- **AND** indicate the valid range

#### Scenario: Priority range validation
- **WHEN** the user provides a priority value outside the range 0~50
- **THEN** the provider SHALL return a validation error
- **AND** indicate the valid range

#### Scenario: Conditional parameter validation
- **WHEN** the user sets type to MX and provides a weight parameter
- **THEN** the provider SHALL return a validation error
- **AND** indicate that weight is not applicable for MX type

#### Scenario: Conditional parameter validation
- **WHEN** the user sets type to MX and does not provide a priority parameter
- **THEN** the provider SHALL use the default value of 0
- **AND** not return a validation error

### Requirement: Import operation

The provider SHALL support importing existing DNS records.

#### Scenario: Successful import
- **WHEN** the user imports a DNS record with the composite ID `zone_id#record_id`
- **AND** the provider calls DescribeDnsRecords with the zone_id and record_id
- **AND** the API returns the DNS record
- **THEN** the provider SHALL populate all attributes
- **AND** store the resource in the Terraform state
- **AND** return the imported resource to the user

#### Scenario: Import non-existent record
- **WHEN** the user attempts to import a DNS record
- **AND** the provider calls DescribeDnsRecords
- **AND** the API returns an empty result set
- **THEN** the provider SHALL return an error
- **AND** indicate that the record does not exist
- **AND** not store any resource in the Terraform state

#### Scenario: Import with invalid ID format
- **WHEN** the user attempts to import a DNS record with an invalid ID format
- **THEN** the provider SHALL return a validation error
- **AND** indicate the expected format: `zone_id#record_id`
- **AND** not store any resource in the Terraform state
