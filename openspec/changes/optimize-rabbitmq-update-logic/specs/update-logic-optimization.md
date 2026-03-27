# Spec: Optimized Update Logic for RabbitMQ Instance

This specification defines the optimized update logic requirements for the `tencentcloud_tdmq_rabbitmq_vip_instance` Terraform resource.

## Requirements

### Requirement: Change Detection Optimization
The Update function SHALL use `d.HasChange()` to detect actual configuration changes and only trigger API calls when necessary.

#### Scenario: User updates cluster name
- **GIVEN** a user changes the `cluster_name` field in the configuration
- **WHEN** the Update function detects the change via `d.HasChange("cluster_name")`
- **THEN** only the ModifyRabbitMQVipInstance API is called with cluster_name parameter
- **AND** other unchanged fields are not sent in the request
- **AND** no unnecessary API calls are made

#### Scenario: User applies same configuration
- **GIVEN** a user applies Terraform configuration without any changes
- **WHEN** the Update function executes
- **THEN** all `d.HasChange()` checks return false
- **AND** no ModifyRabbitMQVipInstance API calls are made
- **AND** only the Read function is called to refresh state
- **AND** the operation completes quickly without network latency

#### Scenario: User updates multiple properties simultaneously
- **GIVEN** a user changes both `cluster_name` and `remark` fields
- **WHEN** the Update function detects both changes
- **THEN** both properties are updated in sequence
- **AND** API calls are made for each changed property
- **AND** unchanged properties are not sent in requests

### Requirement: Independent Update Functions
The Update function SHALL delegate property-specific updates to independent functions.

#### Scenario: Update cluster name via dedicated function
- **GIVEN** the `updateClusterName()` function exists
- **WHEN** called with the resource data and metadata
- **THEN** it constructs a ModifyRabbitMQVipInstance request with only ClusterId and ClusterName
- **AND** it calls the API with retry mechanism
- **AND** it waits for task completion if a task ID is returned
- **AND** it updates the local state with the new value
- **AND** it returns appropriate error messages

#### Scenario: Update remark via dedicated function
- **GIVEN** the `updateRemark()` function exists
- **WHEN** called with the resource data and metadata
- **THEN** it constructs a ModifyRabbitMQVipInstance request with only ClusterId and Remark
- **AND** it follows the same retry and error handling pattern
- **AND** it updates the local state accordingly

#### Scenario: Update spec via dedicated function
- **GIVEN** the `updateSpec()` function exists
- **WHEN** called with the resource data and metadata
- **THEN** it constructs a ModifyRabbitMQVipInstanceSpec request with ClusterId and spec parameters
- **AND** it handles node_count and spec_name changes separately
- **AND** it follows the same retry and task completion pattern

### Requirement: Timeout Configuration
The Update operation SHALL support timeout configuration to prevent long-running operations from blocking.

#### Scenario: Default timeout configuration
- **GIVEN** the resource schema includes Timeout configuration
- **WHEN** examining the Timeout block
- **THEN** the Update timeout is set to a reasonable default (e.g., 20 minutes)
- **AND** the timeout is less than Create/Delete timeouts if appropriate
- **AND** the timeout is documented in the resource schema

#### Scenario: User configures custom timeout
- **GIVEN** a user defines a custom timeout in the Terraform configuration
- **WHEN** the Update function executes
- **THEN** it creates a context with the user-specified timeout
- **AND** the context is used for all API calls within the Update operation
- **AND** the operation is cancelled if it exceeds the timeout
- **AND** an appropriate timeout error is returned to the user

#### Scenario: Timeout cancellation
- **GIVEN** an Update operation is in progress
- **WHEN** the operation exceeds the configured timeout
- **THEN** the context cancellation signal is sent
- **AND** all in-flight API calls are cancelled
- **AND** an error message indicating timeout is returned
- **AND** the resource state is not partially updated

### Requirement: Task Completion Waiting
The Update function SHALL wait for asynchronous tasks to complete before proceeding.

#### Scenario: API returns task ID
- **GIVEN** the Modify API call succeeds and returns a task ID
- **WHEN** the response contains a TaskId field
- **THEN** the Update function calls `waitForTaskCompletion()` with the task ID
- **AND** it polls the DescribeTaskDetail API periodically
- **AND** it waits until the task status is "success"
- **AND** it fails if the task status is "failed"
- **AND** it retries if the task status is "running"

#### Scenario: Task completion timeout
- **GIVEN** a task is in progress
- **WHEN** the task exceeds the retry timeout
- **THEN** the wait function returns an error
- **AND** the error indicates task completion timeout
- **AND** the Update function propagates the error to the user

#### Scenario: Task completion success
- **GIVEN** a task completes successfully
- **WHEN** the DescribeTaskDetail API returns status "success"
- **THEN** the wait function returns nil
- **AND** the Update function proceeds to the next step
- **AND** no error is logged

### Requirement: Error Handling Enhancement
The Update function SHALL provide detailed error messages and handle failures gracefully.

#### Scenario: API call fails
- **GIVEN** an API call during Update operation fails
- **WHEN** the API returns an error
- **THEN** the error is logged with context (log ID, instance ID, property name)
- **AND** the error message includes which property update failed
- **AND** the error is propagated to Terraform
- **AND** the resource state remains unchanged

#### Scenario: Retry mechanism on transient errors
- **GIVEN** an API call fails with a transient error
- **WHEN** the error is retryable (e.g., rate limit, network issue)
- **THEN** the retry mechanism is triggered
- **AND** the API call is retried with exponential backoff
- **AND** the retry continues until success or maximum retries reached
- **AND** the final error is returned if all retries fail

#### Scenario: Task execution fails
- **GIVEN** a task is submitted but fails during execution
- **WHEN** the DescribeTaskDetail API returns status "failed"
- **THEN** the wait function returns an error with the task error message
- **AND** the Update function propagates the error to the user
- **AND** the error includes details about why the task failed

### Requirement: Performance Logging
The Update function SHALL log performance metrics to help identify bottlenecks.

#### Scenario: Log update duration
- **GIVEN** an Update operation starts
- **WHEN** the operation completes successfully or fails
- **THEN** the total duration is logged in milliseconds or seconds
- **AND** the log message includes the instance ID
- **AND** the log message indicates success or failure

#### Scenario: Log individual property updates
- **GIVEN** multiple properties are updated in one operation
- **WHEN** each property update function is called
- **THEN** the start of each update is logged
- **AND** the completion of each update is logged
- **AND** the duration of each update is logged
- **AND** the property name is included in the log

#### Scenario: Log API call details
- **GIVEN** an API call is made during Update
- **WHEN** the API call is initiated
- **THEN** the API name and parameters are logged (excluding sensitive data)
- **AND** the API call duration is logged
- **AND** the success or failure is logged

### Requirement: State Management
The Update function SHALL efficiently manage Terraform state updates.

#### Scenario: Local state update for simple properties
- **GIVEN** a simple property like cluster_name is updated
- **WHEN** the API call succeeds
- **THEN** the local state is updated using `d.Set("cluster_name", newValue)`
- **AND** the value matches what was sent to the API
- **AND** the state update occurs without an additional Read call

#### Scenario: Final Read for complex properties
- **GIVEN** complex properties or multiple properties are updated
- **WHEN** all update operations complete successfully
- **THEN** a final Read call is made to refresh the entire state
- **AND** the Read call ensures state consistency with the cloud resource
- **AND** the Read call uses the same context with timeout

#### Scenario: State update failure handling
- **GIVEN** an API call succeeds but state update fails
- **WHEN** `d.Set()` returns an error
- **THEN** the error is logged and returned to the user
- **AND** the state inconsistency is clearly indicated
- **AND** the user may need to run terraform refresh to resolve

### Requirement: Backward Compatibility
The optimized Update logic SHALL maintain backward compatibility with existing configurations and state.

#### Scenario: Existing resource state migration
- **GIVEN** a resource was created with the previous Update implementation
- **WHEN** the provider is upgraded to the new implementation
- **THEN** the existing state is compatible with the new Update function
- **AND** terraform plan shows no changes for unchanged resources
- **AND** all Update operations work correctly with existing state

#### Scenario: Existing configuration compatibility
- **GIVEN** a user has existing Terraform configuration
- **WHEN** the user upgrades the provider version
- **THEN** the existing configuration works without modification
- **AND** all Update scenarios work as expected
- **AND** no breaking changes are introduced

#### Scenario: API version compatibility
- **GIVEN** the new Update function uses the same TencentCloud APIs
- **WHEN** the Update function is called
- **THEN** it uses the same API endpoints and versions
- **AND** API request and response formats are unchanged
- **AND** no new API permissions are required

### Requirement: No Changes Detection
The Update function SHALL handle the case where no actual changes are detected.

#### Scenario: No changes detected
- **GIVEN** a user runs terraform apply with no configuration changes
- **WHEN** the Update function executes and all `d.HasChange()` checks return false
- **THEN** the function logs a warning message indicating no changes
- **AND** no Modify API calls are made
- **AND** the function returns immediately
- **AND** only a minimal state refresh may occur if needed

#### Scenario: Immutable fields unchanged
- **GIVEN** a user modifies a field that is actually immutable
- **WHEN** the Update function executes
- **THEN** the immutable field change is detected but handled appropriately
- **AND** an error is returned indicating the field cannot be changed
- **AND** no API calls are made for the immutable field

### Requirement: Code Quality and Formatting
All code changes SHALL follow project coding standards and formatting requirements.

#### Scenario: Go code formatting
- **GIVEN** code changes are made to the Update function
- **WHEN** the development is complete
- **THEN** `go fmt` is executed on all modified files
- **AND** all code follows Go standard formatting
- **AND** no formatting warnings or errors exist

#### Scenario: Linting compliance
- **GIVEN** the modified code is ready for review
- **WHEN** `make lint` is executed
- **THEN** no linting errors are reported
- **AND** all code quality checks pass
- **AND** the code follows project linting rules

#### Scenario: Code documentation
- **GIVEN** new functions are added (e.g., `updateClusterName`)
- **WHEN** the functions are implemented
- **THEN** each function has appropriate comments
- **AND** parameters and return values are documented
- **AND** complex logic is explained in comments
