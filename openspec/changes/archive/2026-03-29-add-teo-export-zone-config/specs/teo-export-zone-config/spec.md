## ADDED Requirements

### Requirement: Data source can query zone configuration by zone ID
The system SHALL allow users to query the complete configuration of a TEO zone by providing the zone ID as an input parameter.

#### Scenario: Successful query with valid zone ID
- **WHEN** user provides a valid zone_id to the tencentcloud_teo_export_zone_config data source
- **THEN** system returns the complete zone configuration information
- **AND** the configuration includes zone_id, zone_name, status, area, cname_status, and other basic information
- **AND** the configuration includes domain configurations, security policies, cache rules if available

#### Scenario: Query with non-existent zone ID
- **WHEN** user provides a zone_id that does not exist
- **THEN** system returns an appropriate error message indicating the zone was not found
- **AND** the error message follows the standard error format of the Terraform Provider

### Requirement: Data source supports all zone configuration attributes
The system SHALL return all available configuration attributes for the TEO zone, including but not limited to basic information, domain settings, security policies, and cache configurations.

#### Scenario: Query returns comprehensive configuration
- **WHEN** user queries a zone with complete configuration
- **THEN** system returns all configuration attributes supported by the TEO API
- **AND** the returned data is structured and accessible in Terraform configuration
- **AND** complex nested configurations (e.g., security policies) are properly structured as maps or lists

#### Scenario: Query returns partial configuration for zones with minimal settings
- **WHEN** user queries a zone with minimal configuration
- **THEN** system returns only the available configuration attributes
- **AND** optional fields that are not configured are omitted or set to null
- **AND** the data source handles missing attributes gracefully

### Requirement: Data source handles API errors and retries
The system SHALL implement proper error handling and retry mechanisms for API calls to TEO services.

#### Scenario: API call succeeds on retry after transient failure
- **WHEN** the initial TEO API call fails due to a transient error (e.g., rate limiting, temporary network issue)
- **THEN** system retries the API call with exponential backoff
- **AND** after successful retry, system returns the zone configuration data
- **AND** retry attempts are logged for debugging purposes

#### Scenario: API call fails with permanent error
- **WHEN** the TEO API call fails with a permanent error (e.g., invalid credentials, unauthorized access)
- **THEN** system does not retry the API call
- **AND** system returns a clear error message to the user
- **AND** the error message includes the root cause of the failure

### Requirement: Data source validates input parameters
The system SHALL validate input parameters before making API calls to TEO services.

#### Scenario: Validation fails with missing zone_id
- **WHEN** user omits the required zone_id parameter
- **THEN** system returns a validation error before making any API calls
- **AND** the error message clearly indicates that zone_id is required

#### Scenario: Validation fails with invalid zone_id format
- **WHEN** user provides a zone_id in an invalid format
- **THEN** system returns a validation error
- **AND** the error message indicates the expected format of zone_id

### Requirement: Data source provides consistent and idempotent results
The system SHALL return consistent and idempotent results for multiple queries of the same zone within a short time period.

#### Scenario: Multiple queries return consistent results
- **WHEN** user queries the same zone_id multiple times
- **AND** the zone configuration has not changed between queries
- **THEN** system returns the same configuration data for each query
- **AND** the results are consistent regardless of the number of queries

#### Scenario: Queries reflect eventual consistency
- **WHEN** user queries a zone shortly after configuration changes
- **THEN** system may return the previous configuration due to eventual consistency
- **AND** if the configuration is not yet available, system retries the query
- **AND** after the retry timeout, system returns the most recent available configuration
