## ADDED Requirements

### Requirement: Preview Log Statistics Support
The resource SHALL support specifying test data for preview processing results through the `preview_log_statistics` field.

#### Scenario: Configure preview log data for processing preview
- **WHEN** user sets `task_type` to 1 or 2 (preview mode)
- **AND** user provides `preview_log_statistics` with log content, line numbers, and destination topic IDs
- **THEN** the data transform task is created with the preview log statistics
- **AND** the preview data is used for processing preview

### Requirement: Backup Data Policy Support
The resource SHALL support configuring backup data handling policy through the `backup_give_up_data` field when `func_type` is 2.

#### Scenario: Configure data discard policy for dynamic creation
- **WHEN** user sets `func_type` to 2 (dynamic creation)
- **AND** user sets `backup_give_up_data` to true
- **THEN** logs are discarded when dynamically created logset/topic count exceeds product limits
- **AND** no backup logset/topic is created

#### Scenario: Configure data backup policy for dynamic creation
- **WHEN** user sets `func_type` to 2 (dynamic creation)
- **AND** user sets `backup_give_up_data` to false or omits it (default false)
- **THEN** backup logset and topic are created when limits are exceeded
- **AND** logs are written to the backup topic

### Requirement: Service Log Delivery Support
The resource SHALL support enabling or disabling service log delivery through the `has_services_log` field.

#### Scenario: Enable service log delivery
- **WHEN** user sets `has_services_log` to 2
- **THEN** the data transform task delivers service logs
- **AND** service operation logs are available for monitoring

#### Scenario: Disable service log delivery
- **WHEN** user sets `has_services_log` to 1
- **THEN** the data transform task does not deliver service logs

### Requirement: Data Transform Type Support
The resource SHALL support specifying the data transform type through the `data_transform_type` field.

#### Scenario: Configure standard data transform task
- **WHEN** user sets `data_transform_type` to 0 or omits it
- **THEN** a standard data transform task is created
- **AND** logs are processed after being written to the source topic

#### Scenario: Configure pre-processing data transform task
- **WHEN** user sets `data_transform_type` to 1
- **THEN** a pre-processing data transform task is created
- **AND** collected logs are processed before being written to the source topic

### Requirement: Failure Log Retention Support
The resource SHALL support configuring failure log retention policy through the `keep_failure_log` and `failure_log_key` fields.

#### Scenario: Enable failure log retention
- **WHEN** user sets `keep_failure_log` to 2
- **AND** user provides `failure_log_key` with the field name for failure logs
- **THEN** failed processing logs are retained in the specified field
- **AND** failure logs are available for troubleshooting

#### Scenario: Disable failure log retention
- **WHEN** user sets `keep_failure_log` to 1 or omits it (default 1)
- **THEN** failed processing logs are not retained

### Requirement: Time Range Processing Support
The resource SHALL support specifying the processing time range through the `process_from_timestamp` and `process_to_timestamp` fields.

#### Scenario: Configure historical data processing with time range
- **WHEN** user sets `process_from_timestamp` to a past timestamp
- **AND** user sets `process_to_timestamp` to an end timestamp
- **THEN** the data transform task processes logs only within the specified time range
- **AND** processing stops after reaching the end timestamp

#### Scenario: Configure continuous processing from specific time
- **WHEN** user sets `process_from_timestamp` to a start timestamp
- **AND** user omits `process_to_timestamp`
- **THEN** the data transform task processes logs from the start timestamp continuously
- **AND** processing continues for new incoming logs

### Requirement: SQL Data Source Integration Support
The resource SHALL support associating external database data sources through the `data_transform_sql_data_sources` field.

#### Scenario: Configure MySQL data source for data enrichment
- **WHEN** user provides `data_transform_sql_data_sources` with MySQL connection details
- **AND** user specifies data_source type, region, instance_id, user, password, and alias_name
- **THEN** the data transform task can reference the external database in ETL content
- **AND** logs can be enriched with data from the MySQL database using the specified alias

#### Scenario: Configure multiple data sources
- **WHEN** user provides multiple entries in `data_transform_sql_data_sources`
- **AND** each entry has a unique alias_name
- **THEN** all data sources are registered with the data transform task
- **AND** ETL content can reference multiple databases by their aliases

### Requirement: Environment Variables Support
The resource SHALL support setting environment variables through the `env_infos` field.

#### Scenario: Configure environment variables for ETL processing
- **WHEN** user provides `env_infos` with key-value pairs
- **THEN** the environment variables are available in the ETL content processing context
- **AND** ETL functions can reference these variables during log processing

#### Scenario: Use environment variables for configuration management
- **WHEN** user sets sensitive or environment-specific values in `env_infos`
- **THEN** ETL content can use these variables without hardcoding values
- **AND** configuration can be easily updated by changing environment variables

### Requirement: Read and Update Operations for New Fields
The resource SHALL support reading and updating all newly added fields through standard Terraform operations.

#### Scenario: Read resource state with new fields
- **WHEN** user imports or refreshes an existing data transform task
- **THEN** all new fields are read from the API response
- **AND** the Terraform state is updated with current values

#### Scenario: Update resource with new field changes
- **WHEN** user modifies any of the new fields in Terraform configuration
- **THEN** the resource update operation includes the changed fields
- **AND** the data transform task is updated accordingly
- **AND** immutable fields (if any) prevent resource recreation when changed
