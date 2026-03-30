## ADDED Requirements

### Requirement: Data source provides Offset parameter for pagination
The data source SHALL support an `offset` parameter to skip a specified number of results in the paginated list of realtime log delivery tasks.

#### Scenario: Offset skips first N results
- **WHEN** user specifies `offset = 10`
- **THEN** data source returns results starting from the 11th task
- **AND** skips the first 10 tasks in the result set

#### Scenario: Offset with zero value
- **WHEN** user specifies `offset = 0` or omits the parameter
- **THEN** data source returns results from the beginning
- **AND** behaves as if no offset was specified

### Requirement: Data source provides Limit parameter for pagination size
The data source SHALL support a `limit` parameter to specify the maximum number of results to return in each page.

#### Scenario: Limit restricts result count
- **WHEN** user specifies `limit = 20`
- **THEN** data source returns at most 20 tasks
- **AND** API request includes the limit value

#### Scenario: Limit with large value
- **WHEN** user specifies `limit = 100`
- **THEN** data source returns at most 100 tasks
- **AND** respects API's maximum allowed limit

### Requirement: Data source supports filtering by zone_id
The data source SHALL allow users to filter realtime log delivery tasks by `zone_id` using the filters parameter.

#### Scenario: Filter by single zone_id
- **WHEN** user specifies filter `name = "zone_id"` with `values = ["zone-abc123"]`
- **THEN** data source returns only tasks belonging to zone "zone-abc123"
- **AND** API request includes the zone_id filter

#### Scenario: Filter by multiple zone_ids
- **WHEN** user specifies filter `name = "zone_id"` with `values = ["zone-abc123", "zone-def456"]`
- **THEN** data source returns tasks from both specified zones
- **AND** API request includes multiple zone_id values

### Requirement: Data source supports filtering by task_id
The data source SHALL allow users to filter realtime log delivery tasks by `task_id` using the filters parameter.

#### Scenario: Filter by single task_id
- **WHEN** user specifies filter `name = "task-id"` with `values = ["task-123"]`
- **THEN** data source returns only the task with ID "task-123"
- **AND** API request includes the task-id filter

### Requirement: Data source supports filtering by task_name
The data source SHALL allow users to filter realtime log delivery tasks by `task_name` using the filters parameter with fuzzy matching.

#### Scenario: Filter by exact task_name
- **WHEN** user specifies filter `name = "task-name"` with `values = ["my-task"]`
- **THEN** data source returns tasks matching the exact name "my-task"
- **AND** API request includes the task-name filter

#### Scenario: Filter by partial task_name with fuzzy matching
- **WHEN** user specifies filter `name = "task-name"` with `values = ["my-task"]` and `fuzzy = true`
- **THEN** data source returns tasks containing "my-task" in their name
- **AND** API request includes fuzzy matching parameter

### Requirement: Data source supports filtering by task_type
The data source SHALL allow users to filter realtime log delivery tasks by `task_type` using the filters parameter.

#### Scenario: Filter by task_type
- **WHEN** user specifies filter `name = "task-type"` with `values = ["cls"]`
- **THEN** data source returns only tasks with type "cls"
- **AND** API request includes the task-type filter

#### Scenario: Filter by multiple task_types
- **WHEN** user specifies filter `name = "task-type"` with `values = ["cls", "s3"]`
- **THEN** data source returns tasks with either "cls" or "s3" type
- **AND** API request includes multiple task-type values

### Requirement: Data source returns task list with all required fields
The data source SHALL return a list of realtime log delivery tasks, where each task contains all essential fields including task_id, zone_id, task_name, task_type, status, and other relevant metadata.

#### Scenario: Return complete task information
- **WHEN** user queries the data source with valid zone_id
- **THEN** data source returns a list of tasks
- **AND** each task includes task_id, zone_id, task_name, task_type, delivery_status, and other fields from the API response
- **AND** task structure matches the API's RealtimeLogDeliveryTask model

### Requirement: Data source handles empty results gracefully
The data source SHALL return an empty list when no realtime log delivery tasks match the specified filters or criteria.

#### Scenario: No matching tasks found
- **WHEN** user queries with filters that match no tasks
- **THEN** data source returns an empty task list
- **AND** does not produce an error

#### Scenario: Offset exceeds total count
- **WHEN** user specifies `offset = 1000` when only 100 tasks exist
- **THEN** data source returns an empty task list
- **AND** does not produce an error

### Requirement: Data source handles API errors appropriately
The data source SHALL propagate API errors to Terraform with meaningful error messages.

#### Scenario: API authentication failure
- **WHEN** TencentCloud API returns authentication error
- **THEN** data source returns error with clear message about authentication failure
- **AND** Terraform displays the error to the user

#### Scenario: API rate limit exceeded
- **WHEN** TencentCloud API returns rate limit error
- **THEN** data source returns error indicating rate limit was exceeded
- **AND** suggests retry after backoff period

### Requirement: Data source supports ordering results
The data source SHALL support an optional `order` parameter to specify the field used for sorting results.

#### Scenario: Order by task_id
- **WHEN** user specifies `order = "task-id"`
- **THEN** API request includes sorting by task_id
- **AND** results are returned ordered by task_id

### Requirement: Data source supports sort direction
The data source SHALL support an optional `direction` parameter to specify the sort direction (asc or desc).

#### Scenario: Ascending order
- **WHEN** user specifies `direction = "asc"`
- **THEN** API request includes ascending sort direction
- **AND** results are sorted in ascending order

#### Scenario: Descending order
- **WHEN** user specifies `direction = "desc"`
- **THEN** API request includes descending sort direction
- **AND** results are sorted in descending order
