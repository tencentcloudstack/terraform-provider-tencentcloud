## ADDED Requirements

### Requirement: Resource schema definition
The `tencentcloud_cls_splunk_deliver` resource SHALL expose all parameters supported by the CLS Splunk Deliver APIs in its Terraform schema, including `topic_id`, `name`, `net_info`, `metadata_info`, `has_service_log`, `index_ack`, `source`, `source_type`, `index`, `channel`, `dsl_filter`, `external_role`, and `enable`.

#### Scenario: Schema contains all required fields
- **WHEN** a user inspects the `tencentcloud_cls_splunk_deliver` resource schema
- **THEN** the schema SHALL include `topic_id` (Required), `name` (Required), `net_info` (Required, TypeList, MaxItems:1), and `metadata_info` (Required, TypeList, MaxItems:1)

#### Scenario: Schema contains all optional fields
- **WHEN** a user inspects the `tencentcloud_cls_splunk_deliver` resource schema
- **THEN** the schema SHALL include `has_service_log` (Optional), `index_ack` (Optional), `source` (Optional), `source_type` (Optional), `index` (Optional), `channel` (Optional), `dsl_filter` (Optional), `external_role` (Optional, TypeList, MaxItems:1), and `enable` (Optional, Computed)

#### Scenario: Computed attributes
- **WHEN** a user inspects the `tencentcloud_cls_splunk_deliver` resource schema
- **THEN** the schema SHALL include `task_id` (Computed) as the unique identifier returned by the Create API

### Requirement: Create operation
The resource SHALL create a Splunk deliver task via the `CreateSplunkDeliver` API, setting the resource ID to `task_id#topic_id` on success.

#### Scenario: Successful creation
- **WHEN** a user applies a Terraform configuration with valid `topic_id`, `name`, `net_info`, and `metadata_info`
- **THEN** the provider SHALL call `CreateSplunkDeliver` API and set the resource ID to the format `{task_id}#{topic_id}`

#### Scenario: Create API returns empty task_id
- **WHEN** the `CreateSplunkDeliver` API returns a response with nil or empty `TaskId`
- **THEN** the provider SHALL return a `NonRetryableError` with an appropriate error message

### Requirement: Read operation
The resource SHALL read the Splunk deliver task state via the `DescribeSplunkDelivers` API using `topic_id` and `task_id` filter.

#### Scenario: Successful read
- **WHEN** the provider reads an existing Splunk deliver task
- **THEN** the provider SHALL call `DescribeSplunkDelivers` with `TopicId` and `Filters` containing `taskId` filter, and set all schema attributes from the response

#### Scenario: Task not found in read
- **WHEN** the `DescribeSplunkDelivers` API returns empty `Infos` array
- **THEN** the provider SHALL log the current resource ID and call `d.SetId("")` to remove the resource from state

#### Scenario: Read response fields are nil
- **WHEN** the `DescribeSplunkDelivers` API returns an `Infos` entry with some nil fields
- **THEN** the provider SHALL skip setting those nil fields in the Terraform state

### Requirement: Update operation
The resource SHALL update a Splunk deliver task via the `ModifySplunkDeliver` API.

#### Scenario: Successful update
- **WHEN** a user modifies the `name` or other mutable field of an existing Splunk deliver task
- **THEN** the provider SHALL call `ModifySplunkDeliver` API with the updated parameters and then call `Read` to refresh state

#### Scenario: Update with immutable fields
- **WHEN** a user attempts to modify `topic_id` (which is immutable)
- **THEN** the `topic_id` field SHALL be marked as `ForceNew`, triggering resource recreation instead of update

### Requirement: Delete operation
The resource SHALL delete a Splunk deliver task via the `DeleteSplunkDeliver` API.

#### Scenario: Successful deletion
- **WHEN** a user destroys a Terraform-managed Splunk deliver task
- **THEN** the provider SHALL call `DeleteSplunkDeliver` API with the `TaskId` and `TopicId` parsed from the resource ID

### Requirement: Import operation
The resource SHALL support importing existing Splunk deliver tasks using the `task_id#topic_id` format as the import ID.

#### Scenario: Successful import
- **WHEN** a user imports a resource with `terraform import tencentcloud_cls_splunk_deliver.example task-xxx#topic-yyy`
- **THEN** the provider SHALL parse the import ID, set the resource ID, and call `Read` to populate the state

### Requirement: Retry logic for API calls
All CRUD operations SHALL use `tccommon.ReadRetryTimeout` as the timeout and wrap errors with `tccommon.RetryError()` for transient failures.

#### Scenario: Transient API failure during create
- **WHEN** the `CreateSplunkDeliver` API returns a transient error
- **THEN** the provider SHALL retry the request within the retry timeout period

#### Scenario: Persistent API failure
- **WHEN** the API consistently returns errors beyond the retry timeout
- **THEN** the provider SHALL return the final error to the user

### Requirement: Provider registration
The resource SHALL be registered in the TencentCloud provider so it is available for use.

#### Scenario: Resource registered in provider
- **WHEN** the provider is initialized
- **THEN** `tencentcloud_cls_splunk_deliver` SHALL be a valid resource type that maps to the resource implementation