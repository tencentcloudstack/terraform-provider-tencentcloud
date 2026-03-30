## ADDED Requirements

### Requirement: TaskId based asynchronous operation waiting
The tencentcloud_teo_l7_acc_rule resource SHALL support asynchronous operation waiting when updating via ImportZoneConfig API. The resource SHALL extract the TaskId from the API response and wait for the task to complete before returning.

#### Scenario: Successful task completion
- **WHEN** ImportZoneConfig API is called during resource update
- **AND** API returns a TaskId in the response
- **THEN** the provider SHALL poll the task status using the TaskId
- **AND** SHALL wait until the task status indicates success
- **AND** SHALL return successfully after task completion

#### Scenario: Task failure
- **WHEN** ImportZoneConfig API is called during resource update
- **AND** API returns a TaskId in the response
- **AND** task polling detects a failure status
- **THEN** the provider SHALL return an error with task failure details
- **AND** SHALL not mark the resource update as successful

#### Scenario: Task timeout
- **WHEN** ImportZoneConfig API is called during resource update
- **AND** API returns a TaskId in the response
- **AND** task does not complete within the configured timeout period
- **THEN** the provider SHALL return a timeout error
- **AND** SHALL include the TaskId and timeout duration in the error message

#### Scenario: Task status polling retry on API error
- **WHEN** the provider is polling task status
- **AND** the task status query API returns a transient error (e.g., rate limit, network error)
- **THEN** the provider SHALL retry the query according to the retry policy
- **AND** SHALL continue waiting if retries succeed
- **AND** SHALL return an error only if all retries fail

### Requirement: Task status polling configuration
The task status polling mechanism SHALL be configurable with reasonable default values for polling interval and timeout duration.

#### Scenario: Default polling interval
- **WHEN** polling task status
- **THEN** the provider SHALL use a default polling interval of 5 seconds
- **AND** SHALL make API requests at this interval until task completes or times out

#### Scenario: Default timeout duration
- **WHEN** waiting for task completion
- **THEN** the provider SHALL use a default timeout of 10 minutes
- **AND** SHALL fail with timeout error if task exceeds this duration

#### Scenario: Respect resource timeout configuration
- **WHEN** the resource schema defines a Timeout configuration
- **THEN** the provider SHALL use the configured timeout value instead of the default
- **AND** SHALL apply the timeout to the task waiting operation

### Requirement: Error handling and logging
The provider SHALL implement comprehensive error handling and logging for the task waiting mechanism to facilitate debugging and troubleshooting.

#### Scenario: Log task waiting start
- **WHEN** starting to wait for a task
- **THEN** the provider SHALL log the TaskId and start time
- **AND** SHALL log at debug level

#### Scenario: Log task polling progress
- **WHEN** polling task status
- **THEN** the provider SHALL log the current task status at appropriate intervals
- **AND** SHALL log at debug level to avoid excessive logging

#### Scenario: Log task completion
- **WHEN** task completes successfully or fails
- **THEN** the provider SHALL log the final task status and total waiting time
- **AND** SHALL log at info level for success and error level for failure

#### Scenario: Error message includes TaskId
- **WHEN** an error occurs during task waiting
- **THEN** the error message SHALL include the TaskId
- **AND** SHALL include the last known task status if available
