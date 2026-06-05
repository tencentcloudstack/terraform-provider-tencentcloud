## ADDED Requirements

### Requirement: Data source queries Edge KV key names

The `tencentcloud_teo_edge_k_v_list` data source SHALL call the `EdgeKVList` API to retrieve a list of KV key names from a specified TEO zone namespace.

#### Scenario: Query keys with required parameters only

- **WHEN** user specifies `zone_id` and `namespace` in the data source configuration
- **THEN** the data source SHALL call `EdgeKVList` with `ZoneId`, `Namespace`, and `Limit=1000`, and return the list of key names in the `keys` attribute and the next cursor position in the `cursor` attribute

#### Scenario: Query keys with prefix filter

- **WHEN** user specifies `zone_id`, `namespace`, and `prefix` in the data source configuration
- **THEN** the data source SHALL call `EdgeKVList` with the `Prefix` parameter set, returning only keys that start with the specified prefix

#### Scenario: Query keys with cursor for pagination

- **WHEN** user specifies `zone_id`, `namespace`, and `cursor` in the data source configuration
- **THEN** the data source SHALL call `EdgeKVList` with the `Cursor` parameter set to continue traversal from the specified position

### Requirement: Data source schema definition

The data source schema SHALL define the following attributes:

- `zone_id` (Required, String): The TEO zone ID.
- `namespace` (Required, String): The namespace name.
- `prefix` (Optional, String): Key name prefix filter.
- `cursor` (Optional/Computed, String): Cursor position for pagination (input for query start position, output for next page position).
- `keys` (Computed, List of String): The list of key names returned.
- `result_output_file` (Optional, String): Used to save results to a file.

#### Scenario: Schema validates required fields

- **WHEN** user omits `zone_id` or `namespace`
- **THEN** Terraform SHALL report a validation error indicating the required field is missing

#### Scenario: Schema allows optional fields to be omitted

- **WHEN** user omits `prefix` and `cursor`
- **THEN** the data source SHALL query without prefix filtering and start from the beginning

### Requirement: Retry on transient API failures

The data source SHALL retry API calls using `tccommon.ReadRetryTimeout` and wrap errors with `tccommon.RetryError`.

#### Scenario: Transient API error triggers retry

- **WHEN** the `EdgeKVList` API returns a transient error
- **THEN** the data source SHALL retry the request within the configured timeout period

#### Scenario: Non-retryable error returns immediately

- **WHEN** the `EdgeKVList` API returns a non-retryable error
- **THEN** the data source SHALL return the error to the user without further retries

### Requirement: Provider registration

The data source SHALL be registered in `tencentcloud/provider.go` and documented in `tencentcloud/provider.md`.

#### Scenario: Data source is available in provider

- **WHEN** user references `data.tencentcloud_teo_edge_k_v_list` in their Terraform configuration
- **THEN** the provider SHALL recognize and execute the data source

### Requirement: Limit parameter hardcoded to maximum

The `Limit` parameter SHALL be set to 1000 (the API maximum) internally and SHALL NOT be exposed to users.

#### Scenario: Internal limit is set to maximum

- **WHEN** the data source calls the `EdgeKVList` API
- **THEN** the request SHALL include `Limit=1000` regardless of user configuration
