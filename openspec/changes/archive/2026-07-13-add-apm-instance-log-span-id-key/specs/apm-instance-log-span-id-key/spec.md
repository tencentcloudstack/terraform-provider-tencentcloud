## ADDED Requirements

### Requirement: APM instance supports log_span_id_key parameter
The `tencentcloud_apm_instance` resource SHALL support the `log_span_id_key` parameter, which maps to the `LogSpanIdKey` field of the TencentCloud APM `ModifyApmInstance` API. It specifies the CLS index key for `spanId` and is valid when the CLS index type is key-value index (`log_index_type = 1`).

**Rationale**: The API supports configuring the spanId index key, but the provider currently only exposes the sibling `log_trace_id_key`. Exposing `log_span_id_key` completes the key-value CLS index configuration for APM business systems.

#### Scenario: Create instance with log_span_id_key set
- **WHEN** the user creates a `tencentcloud_apm_instance` with `log_span_id_key` configured
- **THEN** the provider SHALL create the instance via `CreateApmInstance` and apply the `log_span_id_key` value through a subsequent `ModifyApmInstance` call (because `CreateApmInstance` does not accept `LogSpanIdKey`)
- **AND** the Terraform state SHALL reflect the configured `log_span_id_key` value after the Read operation

#### Scenario: Update log_span_id_key on an existing instance
- **WHEN** the user updates `log_span_id_key` on an existing `tencentcloud_apm_instance`
- **THEN** the provider SHALL send the new value in the `ModifyApmInstance` request (`LogSpanIdKey` field)
- **AND** the update SHALL be performed in-place without recreating the resource

#### Scenario: Read log_span_id_key from the API
- **WHEN** the resource is read or refreshed
- **THEN** the provider SHALL read `LogSpanIdKey` from the `ApmInstanceDetail` response of `DescribeApmInstances`
- **AND** the provider SHALL set `log_span_id_key` in Terraform state only when the API-returned `LogSpanIdKey` is not nil

#### Scenario: Omit log_span_id_key
- **WHEN** the user does not specify `log_span_id_key` in the configuration
- **THEN** the provider SHALL NOT send `LogSpanIdKey` in the `ModifyApmInstance` request
- **AND** the API default behavior SHALL apply

#### Scenario: Backward compatibility
- **WHEN** an existing Terraform configuration without `log_span_id_key` is applied
- **THEN** the configuration SHALL continue to work without modification
- **AND** no state migration SHALL be required
