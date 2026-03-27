# Tasks: Optimize RabbitMQ Instance Update Logic

## Implementation Tasks

### 1. Analyze Current Implementation

#### Task: Analyze current Update function
- **Description**: Review the existing `resourceTencentCloudTdmqRabbitmqVipInstanceUpdate` function implementation to understand current behavior, identify bottlenecks, and document issues.
- **Acceptance Criteria**:
  - Document the current Update function structure
  - Identify all API calls made during Update
  - List all properties that can be updated
  - Document current error handling mechanisms
  - Identify performance bottlenecks (e.g., unnecessary API calls, redundant reads)

---

### 2. Implement Change Detection Logic

#### Task: Refactor Update function to use `d.HasChange()`
- **Description**: Modify the Update function to use `d.HasChange()` for each updatable property instead of always calling the API.
- **Acceptance Criteria**:
  - Update function checks `d.HasChange()` for `cluster_name`
  - Update function checks `d.HasChange()` for `remark`
  - Update function checks `d.HasChange()` for `node_count`
  - Update function checks `d.HasChange()` for `spec_name`
  - Update function checks `d.HasChange()` for `auto_renew_flag`
  - API calls are only made when changes are detected
  - No-op scenario (no changes) skips all Modify API calls

---

### 3. Create Independent Update Functions

#### Task: Create `updateClusterName()` function
- **Description**: Extract cluster name update logic into a separate function with proper error handling and retry mechanism.
- **Acceptance Criteria**:
  - Function signature: `updateClusterName(ctx context.Context, d *schema.ResourceData, meta interface{}) error`
  - Constructs ModifyRabbitMQVipInstance request with ClusterId and ClusterName
  - Implements retry mechanism using `resource.RetryContext()`
  - Handles TaskId in response and waits for task completion
  - Updates local state using `d.Set()`
  - Returns descriptive error messages on failure
  - Includes logging for operation start, completion, and duration

#### Task: Create `updateRemark()` function
- **Description**: Extract remark update logic into a separate function with proper error handling and retry mechanism.
- **Acceptance Criteria**:
  - Function signature: `updateRemark(ctx context.Context, d *schema.ResourceData, meta interface{}) error`
  - Constructs ModifyRabbitMQVipInstance request with ClusterId and Remark
  - Implements retry mechanism using `resource.RetryContext()`
  - Handles TaskId in response and waits for task completion
  - Updates local state using `d.Set()`
  - Returns descriptive error messages on failure
  - Includes logging for operation start, completion, and duration

#### Task: Create `updateSpec()` function
- **Description**: Extract spec update logic into a separate function using the ModifyRabbitMQVipInstanceSpec API.
- **Acceptance Criteria**:
  - Function signature: `updateSpec(ctx context.Context, d *schema.ResourceData, meta interface{}) error`
  - Constructs ModifyRabbitMQVipInstanceSpec request with ClusterId
  - Handles node_count and spec_name changes separately
  - Implements retry mechanism using `resource.RetryContext()`
  - Handles TaskId in response and waits for task completion
  - Returns descriptive error messages on failure
  - Includes logging for operation start, completion, and duration

#### Task: Create `updateAutoRenewFlag()` function
- **Description**: Extract auto renew flag update logic into a separate function using the billing API.
- **Acceptance Criteria**:
  - Function signature: `updateAutoRenewFlag(ctx context.Context, d *schema.ResourceData, meta interface{}) error`
  - Constructs ModifyAutoRenewFlag request with ResourceId
  - Implements retry mechanism using `resource.RetryContext()`
  - Updates local state if needed
  - Returns descriptive error messages on failure
  - Includes logging for operation start, completion, and duration

---

### 4. Implement Task Completion Waiting

#### Task: Create `waitForTaskCompletion()` function
- **Description**: Implement a helper function to poll and wait for asynchronous tasks to complete.
- **Acceptance Criteria**:
  - Function signature: `waitForTaskCompletion(ctx context.Context, meta interface{}, taskId string) error`
  - Constructs DescribeTaskDetail request with taskId
  - Implements polling loop with retry mechanism
  - Handles task statuses: "success", "running", "failed"
  - Returns nil on success
  - Returns error on task failure with error message
  - Retries while status is "running"
  - Logs task status checks at appropriate intervals
  - Uses appropriate timeout for task completion

---

### 5. Add Timeout Configuration

#### Task: Add Timeout block to resource schema
- **Description**: Update the resource schema to include Timeout configuration for Update operation.
- **Acceptance Criteria**:
  - Timeout block is added to resource definition
  - Default Update timeout is set (e.g., 20 minutes)
  - Timeout is documented in schema description
  - Create/Delete timeouts remain unchanged

#### Task: Update Update function to use timeout context
- **Description**: Modify the Update function to create a context with timeout and pass it to all sub-functions.
- **Acceptance Criteria**:
  - Update function creates context with `context.WithTimeout(ctx, d.Timeout(schema.TimeoutUpdate))`
  - Context is deferred to prevent resource leaks
  - All sub-functions (updateClusterName, updateRemark, etc.) receive the context
  - All API calls use the context
  - Timeout cancellation is handled gracefully

---

### 6. Enhance Error Handling

#### Task: Improve error messages and logging
- **Description**: Enhance error handling to provide more descriptive error messages and comprehensive logging.
- **Acceptance Criteria**:
  - Each update function logs operation start with property name and instance ID
  - Each update function logs operation completion with duration
  - Error messages include which property update failed
  - Error messages include instance ID and log ID for troubleshooting
  - API errors are logged with full context
  - Success cases are logged appropriately
  - No-op case (no changes) is logged with a warning

#### Task: Implement retry mechanism with exponential backoff
- **Description**: Ensure all API calls use the project's retry mechanism with exponential backoff.
- **Acceptance Criteria**:
  - All update functions use `resource.RetryContext()`
  - Retry timeout is configured appropriately (e.g., `tccommon.WriteRetryTimeout`)
  - Retryable errors trigger retries
  - Non-retryable errors are returned immediately
  - Retry logic follows project patterns

---

### 7. Optimize State Management

#### Task: Implement local state updates for simple properties
- **Description**: For simple property updates (cluster_name, remark), update local state directly instead of calling Read.
- **Acceptance Criteria**:
  - `updateClusterName()` uses `d.Set("cluster_name", newValue)` on success
  - `updateRemark()` uses `d.Set("remark", newValue)` on success
  - State updates are atomic (no partial updates)
  - State update failures return errors to caller

#### Task: Implement final Read for complex updates
- **Description**: Ensure that after all updates complete, a final Read is called to refresh state from the cloud.
- **Acceptance Criteria**:
  - Final Read is called only when at least one update was performed
  - Final Read is skipped when no changes were detected
  - Final Read uses the same context with timeout
  - Final Read ensures state consistency with cloud resource

---

### 8. Update Main Update Function

#### Task: Refactor main Update function to use new architecture
- **Description**: Rewrite the main Update function to orchestrate the independent update functions with proper change detection.
- **Acceptance Criteria**:
  - Update function creates timeout context
  - Update function checks changes for each updatable property
  - Update function calls appropriate update functions for changed properties
  - Update function tracks if any changes were made
  - Update function skips API calls when no changes detected
  - Update function returns early when no changes detected
  - Update function calls final Read when changes were made
  - Update function returns appropriate diagnostics

---

## Verification Tasks

### 9. Unit Tests

#### Task: Write unit tests for `updateClusterName()`
- **Description**: Create comprehensive unit tests for the cluster name update function.
- **Acceptance Criteria**:
  - Test successful update scenario
  - Test API error scenario
  - Test retry mechanism
  - Test task completion waiting
  - Test state update
  - Mock API client appropriately
  - All tests pass

#### Task: Write unit tests for `updateRemark()`
- **Description**: Create comprehensive unit tests for the remark update function.
- **Acceptance Criteria**:
  - Test successful update scenario
  - Test API error scenario
  - Test retry mechanism
  - Test task completion waiting
  - Test state update
  - Mock API client appropriately
  - All tests pass

#### Task: Write unit tests for `updateSpec()`
- **Description**: Create comprehensive unit tests for the spec update function.
- **Acceptance Criteria**:
  - Test successful update scenario
  - Test API error scenario
  - Test retry mechanism
  - Test task completion waiting
  - Mock API client appropriately
  - All tests pass

#### Task: Write unit tests for `waitForTaskCompletion()`
- **Description**: Create comprehensive unit tests for the task completion waiting function.
- **Acceptance Criteria**:
  - Test successful task completion
  - Test task failure scenario
  - Test timeout scenario
  - Test retryable status (running)
  - Mock DescribeTaskDetail API appropriately
  - All tests pass

#### Task: Write unit tests for main Update function
- **Description**: Create comprehensive unit tests for the main Update function orchestrator.
- **Acceptance Criteria**:
  - Test single property update (cluster_name)
  - Test single property update (remark)
  - Test single property update (spec)
  - Test multiple property updates simultaneously
  - Test no changes scenario
  - Test error handling for failed updates
  - Test timeout handling
  - All tests pass

---

### 10. Integration Tests

#### Task: Write integration test for Update cluster_name
- **Description**: Create an acceptance test that updates the cluster_name property.
- **Acceptance Criteria**:
  - Test uses `TF_ACC=1` environment variable
  - Test creates an instance with initial cluster_name
  - Test updates cluster_name to a new value
  - Test verifies the update was successful via Read
  - Test resource is cleaned up after test
  - Test passes in CI environment

#### Task: Write integration test for Update remark
- **Description**: Create an acceptance test that updates the remark property.
- **Acceptance Criteria**:
  - Test uses `TF_ACC=1` environment variable
  - Test creates an instance with initial remark
  - Test updates remark to a new value
  - Test verifies the update was successful via Read
  - Test resource is cleaned up after test
  - Test passes in CI environment

#### Task: Write integration test for Update spec
- **Description**: Create an acceptance test that updates the spec (node_count and/or spec_name).
- **Acceptance Criteria**:
  - Test uses `TF_ACC=1` environment variable
  - Test creates an instance with initial spec
  - Test updates node_count or spec_name
  - Test verifies the update was successful via Read
  - Test resource is cleaned up after test
  - Test passes in CI environment

#### Task: Write integration test for Update multiple properties
- **Description**: Create an acceptance test that updates multiple properties simultaneously.
- **Acceptance Criteria**:
  - Test uses `TF_ACC=1` environment variable
  - Test creates an instance with initial values
  - Test updates cluster_name, remark, and spec simultaneously
  - Test verifies all updates were successful via Read
  - Test resource is cleaned up after test
  - Test passes in CI environment

#### Task: Write integration test for Update no changes
- **Description**: Create an acceptance test that applies configuration without any changes.
- **Acceptance Criteria**:
  - Test uses `TF_ACC=1` environment variable
  - Test creates an instance with specific configuration
  - Test applies the same configuration (no changes)
  - Test verifies no Modify API calls are made
  - Test verifies operation completes quickly
  - Test resource is cleaned up after test
  - Test passes in CI environment

#### Task: Write integration test for Update timeout
- **Description**: Create an acceptance test that tests timeout configuration.
- **Acceptance Criteria**:
  - Test uses `TF_ACC=1` environment variable
  - Test creates an instance
  - Test applies a custom timeout in configuration
  - Test verifies timeout is respected (may require mocking slow API)
  - Test resource is cleaned up after test
  - Test passes in CI environment

---

### 11. Code Quality Tasks

#### Task: Format code with go fmt
- **Description**: Ensure all modified code files are formatted using `go fmt`.
- **Acceptance Criteria**:
  - Run `go fmt` on all modified Go files
  - No formatting warnings or errors
  - Code follows Go standard formatting

#### Task: Run linting checks
- **Description**: Ensure all modified code passes project linting rules.
- **Acceptance Criteria**:
  - Run `make lint` or equivalent
  - No linting errors are reported
  - Code follows project linting configuration

#### Task: Run tfproviderlint checks
- **Description**: Ensure all modified code passes Terraform-specific linting rules.
- **Acceptance Criteria**:
  - Run tfproviderlint on modified files
  - No Terraform-specific linting errors
  - Code follows Terraform provider best practices

---

### 12. Documentation Tasks

#### Task: Update resource documentation
- **Description**: Update the resource documentation to reflect the timeout configuration and any user-visible changes.
- **Acceptance Criteria**:
  - Update the resource's markdown documentation file
  - Document the new Timeout configuration
  - Update examples if needed
  - Document any behavior changes (e.g., improved performance)
  - Documentation is generated using `make doc`

#### Task: Update code comments
- **Description**: Add or update comments in the code to explain the new architecture.
- **Acceptance Criteria**:
  - Each new function has a comment explaining its purpose
  - Complex logic has inline comments
  - Parameters and return values are documented
  - Comments follow project documentation standards

---

### 13. Validation Tasks

#### Task: Run existing test suite
- **Description**: Run the existing test suite to ensure no regressions.
- **Acceptance Criteria**:
  - Run `make test` for unit tests
  - Run `make testacc` for acceptance tests (if environment available)
  - All existing tests pass
  - No new test failures introduced

#### Task: Verify backward compatibility
- **Description**: Ensure the changes are backward compatible with existing state and configurations.
- **Acceptance Criteria**:
  - Create a test scenario with existing state format
  - Run `terraform plan` and `terraform apply` on existing configuration
  - Verify no unexpected changes or errors
  - Verify existing resources continue to work correctly

#### Task: Performance testing
- **Description**: Verify that the optimizations improve performance as expected.
- **Acceptance Criteria**:
  - Measure Update operation time before and after changes
  - Verify that no-change scenario is faster (reduced API calls)
  - Verify that single-change scenario performance is comparable
  - Document performance improvements

---

### 14. Final Preparation Tasks

#### Task: Create changelog entry
- **Description**: Create a changelog entry describing the optimization.
- **Acceptance Criteria**:
  - Create a file in `.changelog/` directory
  - File name follows pattern: `PR_NUMBER.txt`
  - Content includes summary of optimization and benefits
  - Content is clear and concise

#### Task: Prepare for code review
- **Description**: Prepare all materials for code review.
- **Acceptance Criteria**:
  - All code changes are ready for review
  - All tests pass
  - Documentation is updated
  - Changelog entry is created
  - Pull request description explains the changes clearly

---

## Task Order

The tasks should be completed in the following order:

1. **Phase 1 - Analysis**: Task 1
2. **Phase 2 - Change Detection**: Task 2
3. **Phase 3 - Independent Functions**: Tasks 3.1-3.4
4. **Phase 4 - Task Waiting**: Task 4
5. **Phase 5 - Timeout**: Tasks 5.1-5.2
6. **Phase 6 - Error Handling**: Tasks 6.1-6.2
7. **Phase 7 - State Management**: Tasks 7.1-7.2
8. **Phase 8 - Main Function**: Task 8
9. **Phase 9 - Unit Tests**: Tasks 9.1-9.5
10. **Phase 10 - Integration Tests**: Tasks 10.1-10.6
11. **Phase 11 - Code Quality**: Tasks 11.1-11.3
12. **Phase 12 - Documentation**: Tasks 12.1-12.2
13. **Phase 13 - Validation**: Tasks 13.1-13.3
14. **Phase 14 - Final**: Tasks 14.1-14.2

Each phase should be completed and tested before proceeding to the next phase to ensure incremental validation and early bug detection.
