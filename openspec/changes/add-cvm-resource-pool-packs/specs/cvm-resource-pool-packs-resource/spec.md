## ADDED Requirements

### Requirement: Resource creation via PurchaseResourcePoolPacks API

The resource SHALL support creating CVM resource pool packs by calling the `PurchaseResourcePoolPacks` API with user-provided configuration.

#### Scenario: Successful resource pool pack creation
- **WHEN** user defines `tencentcloud_cvm_resource_pool_packs` resource with valid parameters
- **THEN** provider calls `PurchaseResourcePoolPacks` API and stores the returned pack ID in Terraform state

#### Scenario: Creation with invalid parameters
- **WHEN** user provides invalid parameters (e.g., negative quantity)
- **THEN** provider returns clear error message from API validation

#### Scenario: API call failure during creation
- **WHEN** `PurchaseResourcePoolPacks` API call fails due to network or service error
- **THEN** provider retries with exponential backoff and returns error if all retries exhausted

### Requirement: Resource reading via DescribeResourcePoolPacks API

The resource SHALL support reading resource pool pack details by calling the `DescribeResourcePoolPacks` API with the pack ID from Terraform state.

#### Scenario: Successful resource state refresh
- **WHEN** Terraform performs a refresh operation
- **THEN** provider calls `DescribeResourcePoolPacks` API with pack ID and updates local state with current values

#### Scenario: Resource not found during read
- **WHEN** pack ID no longer exists in cloud (manually deleted or expired)
- **THEN** provider removes resource from Terraform state without error

#### Scenario: API call failure during read
- **WHEN** `DescribeResourcePoolPacks` API call fails due to network or service error
- **THEN** provider retries with exponential backoff and returns error if all retries exhausted

### Requirement: Resource deletion via TerminateResourcePoolPacks API

The resource SHALL support deleting resource pool packs by calling the `TerminateResourcePoolPacks` API with the pack ID.

#### Scenario: Successful resource deletion
- **WHEN** user runs `terraform destroy` or removes the resource
- **THEN** provider calls `TerminateResourcePoolPacks` API and removes resource from state upon success

#### Scenario: Resource already deleted
- **WHEN** pack ID no longer exists during deletion
- **THEN** provider treats as successful deletion and removes from state

#### Scenario: Deletion blocked by active resources
- **WHEN** pack has active instances preventing termination
- **THEN** provider returns clear error message explaining the blocking condition

### Requirement: All fields marked as ForceNew

The resource SHALL mark all configuration fields as `ForceNew: true` in the schema, requiring resource recreation for any changes.

#### Scenario: User modifies any field
- **WHEN** user changes any resource configuration parameter
- **THEN** Terraform plans a destroy and recreate operation

#### Scenario: No update operation attempted
- **WHEN** user attempts to update resource
- **THEN** provider never calls any Update API (none exists)

### Requirement: Service layer abstraction

The resource SHALL use service layer methods in `service_tencentcloud_cvm.go` for all API interactions, not direct SDK calls.

#### Scenario: Create operation uses service layer
- **WHEN** resource Create function is called
- **THEN** it invokes `CreateCvmResourcePoolPacks()` service method

#### Scenario: Read operation uses service layer
- **WHEN** resource Read function is called
- **THEN** it invokes `DescribeCvmResourcePoolPackById()` service method

#### Scenario: Delete operation uses service layer
- **WHEN** resource Delete function is called
- **THEN** it invokes `DeleteCvmResourcePoolPacks()` service method

### Requirement: Retry logic for eventual consistency

The resource SHALL implement retry logic using `resource.Retry` with `tccommon.ReadRetryTimeout` for query operations.

#### Scenario: Query retries on transient errors
- **WHEN** `DescribeResourcePoolPacks` returns transient error (rate limit, temporary unavailable)
- **THEN** provider retries with exponential backoff up to configured timeout

#### Scenario: Query succeeds within retry window
- **WHEN** subsequent retry succeeds after initial failure
- **THEN** provider continues normally without exposing transient error to user

### Requirement: Standard error handling patterns

The resource SHALL use `defer tccommon.LogElapsed()` and `defer tccommon.InconsistentCheck()` for error handling and logging.

#### Scenario: Operations log elapsed time
- **WHEN** any CRUD operation executes
- **THEN** elapsed time is logged via `tccommon.LogElapsed()`

#### Scenario: Inconsistent state detection
- **WHEN** resource state becomes inconsistent
- **THEN** `tccommon.InconsistentCheck()` detects and logs the inconsistency

### Requirement: Resource registration in provider

The resource SHALL be registered in `tencentcloud/provider.go` ResourcesMap and declared in `tencentcloud/provider.md`.

#### Scenario: Provider includes new resource
- **WHEN** provider initializes
- **THEN** `tencentcloud_cvm_resource_pool_packs` is available in ResourcesMap

#### Scenario: Documentation lists new resource
- **WHEN** documentation is generated
- **THEN** resource appears in provider.md resource list

### Requirement: Comprehensive test coverage

The resource SHALL have acceptance tests covering create, read, delete, and import operations.

#### Scenario: Basic resource lifecycle test
- **WHEN** acceptance test runs
- **THEN** test successfully creates, reads, and destroys resource

#### Scenario: Import test
- **WHEN** resource is imported via `terraform import`
- **THEN** resource state is populated correctly from API

#### Scenario: ForceNew behavior test
- **WHEN** test modifies resource field
- **THEN** test verifies destroy and recreate plan is generated

### Requirement: Resource documentation

The resource SHALL have a `.md` documentation file following provider conventions with usage examples.

#### Scenario: Documentation includes basic example
- **WHEN** user views resource documentation
- **THEN** at least one complete usage example is provided

#### Scenario: All fields documented
- **WHEN** user views argument reference
- **THEN** all schema fields are documented with type and description

#### Scenario: ForceNew behavior documented
- **WHEN** user views documentation
- **THEN** ForceNew behavior is clearly explained
