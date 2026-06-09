## ADDED Requirements

### Requirement: Resource manages live origin stream configuration
The provider SHALL expose a resource `tencentcloud_live_origin_stream_info_config` that manages the live origin stream settings for a playback domain via `ModifyOriginStreamInfo` and `DescribeOriginStreamInfo` APIs.

#### Scenario: Create applies configuration
- **WHEN** a user defines the resource with `domain_name` and required fields
- **THEN** the provider SHALL call `ModifyOriginStreamInfo`, poll until `Status` is `1` or `3`, and set the resource ID to `domain_name`

#### Scenario: Read reflects current state
- **WHEN** the provider reads the resource
- **THEN** it SHALL call `DescribeOriginStreamInfo` and map all returned fields to state

#### Scenario: Update modifies configuration
- **WHEN** any mutable field changes
- **THEN** the provider SHALL call `ModifyOriginStreamInfo` with all current field values and poll for terminal status

#### Scenario: Delete is a no-op
- **WHEN** the resource is destroyed
- **THEN** the provider SHALL return nil without calling any API (configuration cannot be deleted)

### Requirement: Schema matches ModifyOriginStreamInfo request parameters
The resource schema SHALL contain fields corresponding to every parameter in `ModifyOriginStreamInfo`:
- Required: `domain_name`, `origin_stream_play_type`, `cdn_stream_play_type`, `origin_stream_type`, `origin_address`, `origin_address_type`
- Optional: `customer_name`, `origin_host`, `origin_timeout`, `origin_retry_times`, `pass_through_http_header`, `pass_through_response`, `pass_through_param`, `indexer_cache`, `fragment_cache`, `hls_play_fragment_count`, `hls_play_fragment_duration`, `time_jitter`, `using_https`, `cache_follow_origin`, `cache_status_code`, `url_replace_rules`, `options_request`, `follow_redirect`, `indexer_keep_param`, `fragment_keep_param`, `media_package_type`, `media_package_channel_types`, `indexer_header`, `fragment_header`, `customization_rules`, `cache_format_rule`
- Computed: `status`

#### Scenario: Required fields missing
- **WHEN** a user omits `domain_name` or other required fields
- **THEN** the provider SHALL return a validation error before making any API call

### Requirement: Async polling after modify
The provider SHALL poll `DescribeOriginStreamInfo` after each `ModifyOriginStreamInfo` call.

#### Scenario: Configuration succeeds
- **WHEN** `DescribeOriginStreamInfo` returns `Status == 1`
- **THEN** polling SHALL stop and the operation SHALL be considered successful

#### Scenario: Configuration closed successfully
- **WHEN** `DescribeOriginStreamInfo` returns `Status == 3`
- **THEN** polling SHALL stop and the operation SHALL be considered successful

#### Scenario: Configuration in progress
- **WHEN** `DescribeOriginStreamInfo` returns `Status == 0` or `Status == 2`
- **THEN** polling SHALL continue retrying within WriteRetryTimeout

### Requirement: customization_rules nested block
The `customization_rules` field SHALL be a `TypeList` where each element maps to `OriginStreamCustomizationRule` struct fields.

#### Scenario: Multiple rules configured
- **WHEN** a user specifies multiple `customization_rules` blocks
- **THEN** all rules SHALL be passed as a slice to the API request
