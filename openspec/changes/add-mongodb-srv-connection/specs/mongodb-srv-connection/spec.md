## ADDED Requirements

### Requirement: SRV Connection URL Management
The system SHALL provide a Terraform resource `tencentcloud_mongodb_instance_srv_connection` to manage MongoDB instance SRV connection URL configuration through TencentCloud API.

#### Scenario: Enable SRV connection with default domain
- **GIVEN** a MongoDB instance exists with ID `cmgo-xxxxxxxx`
- **WHEN** user creates the resource without specifying `domain`
- **THEN** the system SHALL call `EnableSRVConnectionUrl` API without custom domain parameter
- **AND** wait for the async task to complete using `DescribeAsyncRequestInfo`
- **AND** call `DescribeSRVConnectionDomain` to populate the `domain` field with system default domain
- **AND** return the generated SRV URL

#### Scenario: Enable SRV connection with custom domain
- **GIVEN** a MongoDB instance exists with ID `cmgo-xxxxxxxx`
- **WHEN** user creates the resource with `domain = "example.mongodb.com"`
- **THEN** the system SHALL call `EnableSRVConnectionUrl` API with the custom domain
- **AND** wait for the async task to complete using `DescribeAsyncRequestInfo`
- **AND** the `domain` field SHALL reflect the custom domain value
- **AND** return the generated SRV URL with the custom domain

#### Scenario: Query SRV connection status
- **GIVEN** an SRV connection is enabled for instance `cmgo-xxxxxxxx`
- **WHEN** user reads the resource or triggers a refresh
- **THEN** the system SHALL call `DescribeSRVConnectionDomain` API
- **AND** populate the `srv_url` and `domain` attributes with API response values

#### Scenario: Update to add custom domain
- **GIVEN** an SRV connection exists without custom domain (using system default)
- **WHEN** user updates the resource to add `domain = "new.mongodb.com"`
- **THEN** the system SHALL call `ModifySRVConnectionUrl` API with the new domain
- **AND** wait for the async task to complete using `DescribeAsyncRequestInfo`
- **AND** update the resource state with new `domain` and `srv_url`

#### Scenario: Update to change custom domain
- **GIVEN** an SRV connection exists with custom domain `old.mongodb.com`
- **WHEN** user updates `domain` to `new.mongodb.com`
- **THEN** the system SHALL call `ModifySRVConnectionUrl` API with the new domain
- **AND** wait for the async task to complete using `DescribeAsyncRequestInfo`
- **AND** update the resource state with new `domain` and `srv_url`

#### Scenario: Disable SRV connection
- **GIVEN** an SRV connection is enabled for instance `cmgo-xxxxxxxx`
- **WHEN** user destroys the resource
- **THEN** the system SHALL call `DisableSRVConnectionUrl` API
- **AND** remove the resource from Terraform state

#### Scenario: Import existing SRV connection
- **GIVEN** an SRV connection already exists for instance `cmgo-xxxxxxxx`
- **WHEN** user imports the resource using `terraform import tencentcloud_mongodb_instance_srv_connection.example cmgo-xxxxxxxx`
- **THEN** the system SHALL call `DescribeSRVConnectionDomain` API
- **AND** populate the resource state with current `domain` and `srv_url` values

#### Scenario: Handle async task failure
- **GIVEN** user creates or updates an SRV connection resource
- **WHEN** the async task returned by `EnableSRVConnectionUrl` or `ModifySRVConnectionUrl` fails
- **THEN** the system SHALL detect the failure through `DescribeAsyncRequestInfo`
- **AND** return an error to Terraform
- **AND** not update the resource state

#### Scenario: Handle async task timeout
- **GIVEN** user creates or updates an SRV connection resource
- **WHEN** the async task does not complete within the timeout period (3 * ReadRetryTimeout)
- **THEN** the system SHALL return a timeout error to Terraform
- **AND** not update the resource state

### Requirement: Resource Schema Definition
The resource SHALL define the following schema fields with appropriate types and constraints.

#### Scenario: Required instance_id field
- **GIVEN** the resource schema definition
- **THEN** `instance_id` SHALL be a Required field of type String
- **AND** SHALL have ForceNew = true to prevent in-place updates
- **AND** SHALL accept MongoDB instance IDs in format `cmgo-*`

#### Scenario: Optional and Computed domain field
- **GIVEN** the resource schema definition
- **THEN** `domain` SHALL have both Optional = true and Computed = true
- **AND** SHALL be of type String
- **AND** SHALL allow users to omit the field (system will use default domain and populate it after Read)
- **AND** SHALL allow users to specify a custom domain name
- **AND** SHALL be populated from API response in Read operations
- **AND** changes to this field SHALL trigger Update operation

#### Scenario: Computed srv_url field
- **GIVEN** the resource schema definition
- **THEN** `srv_url` SHALL be a Computed field of type String
- **AND** SHALL be populated from `DescribeSRVConnectionDomain` API response
- **AND** SHALL contain the complete SRV connection URL

### Requirement: Service Layer Methods
The MongoDB service layer SHALL provide methods to interact with SRV connection APIs.

#### Scenario: EnableSRVConnectionUrl service method
- **GIVEN** the MongoDB service client
- **WHEN** `EnableSRVConnectionUrl` method is called with instance_id and optional domain
- **THEN** the system SHALL create a `EnableSRVConnectionUrlRequest`
- **AND** if domain parameter is provided and not empty, include it in the request
- **AND** if domain parameter is nil or empty, omit it from the request (use system default)
- **AND** call the TencentCloud SDK `EnableSRVConnectionUrl` API
- **AND** extract the FlowId/AsyncRequestId from the response
- **AND** call `DescribeAsyncRequestInfo` to wait for task completion
- **AND** return error if task fails or times out

#### Scenario: DescribeSRVConnectionDomain service method
- **GIVEN** the MongoDB service client
- **WHEN** `DescribeSRVConnectionDomain` method is called with instance_id
- **THEN** the system SHALL create a `DescribeSRVConnectionDomainRequest`
- **AND** call the TencentCloud SDK `DescribeSRVConnectionDomain` API
- **AND** return the SRV URL and domain information from the response
- **AND** handle cases where SRV connection is not enabled (return nil or empty response)

#### Scenario: ModifySRVConnectionUrl service method
- **GIVEN** the MongoDB service client
- **WHEN** `ModifySRVConnectionUrl` method is called with instance_id and domain
- **THEN** the system SHALL create a `ModifySRVConnectionUrlRequest`
- **AND** include the domain parameter in the request
- **AND** call the TencentCloud SDK `ModifySRVConnectionUrl` API
- **AND** extract the FlowId/AsyncRequestId from the response
- **AND** call `DescribeAsyncRequestInfo` to wait for task completion
- **AND** return error if task fails or times out

#### Scenario: DisableSRVConnectionUrl service method
- **GIVEN** the MongoDB service client
- **WHEN** `DisableSRVConnectionUrl` method is called with instance_id
- **THEN** the system SHALL create a `DisableSRVConnectionUrlRequest`
- **AND** call the TencentCloud SDK `DisableSRVConnectionUrl` API
- **AND** handle synchronous completion (no async task)

### Requirement: Error Handling and Logging
The resource SHALL implement comprehensive error handling and logging following project conventions.

#### Scenario: Log operation elapsed time
- **GIVEN** any CRUD operation on the resource
- **WHEN** the operation starts
- **THEN** the system SHALL defer `tccommon.LogElapsed` to log the operation duration
- **AND** include operation type in the log message

#### Scenario: Perform inconsistent check
- **GIVEN** any CRUD operation on the resource
- **WHEN** the operation starts
- **THEN** the system SHALL defer `tccommon.InconsistentCheck` for state validation
- **AND** detect state inconsistencies between Terraform and cloud provider

#### Scenario: Handle API rate limiting
- **GIVEN** any API call to TencentCloud SDK
- **WHEN** the call is about to be made
- **THEN** the system SHALL call `ratelimit.Check(request.GetAction())`
- **AND** respect the rate limiting constraints

#### Scenario: Retry transient errors
- **GIVEN** an API call fails with a transient error
- **WHEN** the error occurs within a retry block
- **THEN** the system SHALL use `resource.Retry` with appropriate timeout
- **AND** return `resource.RetryableError` for transient failures
- **AND** return `resource.NonRetryableError` for permanent failures

### Requirement: Testing Coverage
The resource SHALL have comprehensive acceptance tests covering all major scenarios.

#### Scenario: Basic lifecycle test
- **GIVEN** a test MongoDB instance
- **WHEN** running acceptance test for basic lifecycle
- **THEN** the test SHALL create the resource without domain parameter
- **AND** verify the resource is created successfully
- **AND** verify srv_url is populated
- **AND** verify domain is automatically populated with system default value
- **AND** update the resource by adding a domain value
- **AND** verify the update is applied successfully
- **AND** verify domain field reflects the new custom domain
- **AND** destroy the resource
- **AND** verify the resource is deleted successfully

#### Scenario: Custom domain test
- **GIVEN** a test MongoDB instance
- **WHEN** running acceptance test for custom domain
- **THEN** the test SHALL create the resource with domain specified
- **AND** verify the custom domain is used in srv_url
- **AND** verify the domain field matches the specified value
- **AND** verify the resource can be read and refreshed

#### Scenario: Domain update test
- **GIVEN** an existing SRV connection resource
- **WHEN** running acceptance test for domain updates
- **THEN** the test SHALL update domain from one custom value to another
- **AND** verify the update triggers ModifySRVConnectionUrl
- **AND** verify srv_url and domain are updated correctly

#### Scenario: Import test
- **GIVEN** an existing SRV connection configuration
- **WHEN** running acceptance test for import
- **THEN** the test SHALL import the resource using instance_id
- **AND** verify domain attribute is correctly populated from API
- **AND** verify srv_url attribute is correctly populated
- **AND** verify the imported state matches the actual cloud state

### Requirement: Documentation
The resource SHALL have complete documentation for users.

#### Scenario: Basic usage example
- **GIVEN** the resource documentation
- **THEN** it SHALL include a basic example without domain parameter
- **AND** show how to reference a MongoDB instance
- **AND** show the output attributes (srv_url, domain)
- **AND** explain that domain will be automatically populated with system default

#### Scenario: Custom domain example
- **GIVEN** the resource documentation
- **THEN** it SHALL include an example with domain parameter
- **AND** explain when to use custom domains
- **AND** show the complete resource configuration

#### Scenario: Domain field documentation
- **GIVEN** the resource documentation
- **THEN** it SHALL clearly document that domain has both Optional and Computed attributes
- **AND** explain that users can omit it to use system default
- **AND** explain that users can specify a custom value
- **AND** explain that it will be populated after creation/import

#### Scenario: Import documentation
- **GIVEN** the resource documentation
- **THEN** it SHALL document the import command syntax
- **AND** provide an example of importing by instance_id
- **AND** explain that domain and srv_url will be populated during import

#### Scenario: Argument reference
- **GIVEN** the resource documentation
- **THEN** it SHALL list all input arguments with types and descriptions
- **AND** mark instance_id as Required with ForceNew
- **AND** mark domain as Optional and Computed
- **AND** explain the behavior differences

#### Scenario: Attribute reference
- **GIVEN** the resource documentation
- **THEN** it SHALL list all output attributes with types and descriptions
- **AND** explain what each attribute represents
- **AND** note that domain appears in both argument and attribute sections due to Optional+Computed
- **AND** note which attributes are only available after creation
