## ADDED Requirements

### Requirement: WAF CC Session supports the precise-match key parameter

The `tencentcloud_waf_cc_session` resource SHALL support an optional `key` string parameter that configures the precise-match session key (`request.Key`) for the TencentCloud WAF `UpsertSession` API. The parameter SHALL be sent on both create and update operations and SHALL be read from the `DescribeSession` API response (`SessionItem.Key`).

**Rationale**: The TencentCloud WAF `UpsertSession` API already supports the `Key` field for precise-matching scenarios, but the Terraform resource does not expose it. Users need to configure precise-match session keys directly from Terraform.

#### Scenario: Create a WAF CC session with the key parameter

- **WHEN** a user creates a `tencentcloud_waf_cc_session` resource with `key` set to a non-empty string
- **THEN** the provider SHALL include the value in `request.Key` when calling the `UpsertSession` API
- **AND** the created session SHALL have the `key` value reflected in Terraform state after read

**Example Configuration**:
```hcl
resource "tencentcloud_waf_cc_session" "example" {
  domain           = "www.demo.com"
  source           = "get"
  category         = "match"
  key_or_start_mat = "key_a=123"
  end_mat          = "&"
  start_offset     = "-1"
  end_offset       = "-1"
  edition          = "sparta-waf"
  session_name     = "terraformDemo"
  key              = "sessionId"
}
```

**Acceptance Criteria**:
- The `key` parameter is accepted as a string value
- The value is sent to the TencentCloud API in `request.Key`
- Terraform state reflects the configured value after creation

#### Scenario: Read the key parameter from an existing session

- **WHEN** the provider reads an existing WAF CC session via the `DescribeSession` API and the response `SessionItem.Key` is non-nil
- **THEN** the provider SHALL set `key` in Terraform state via `d.Set("key", ...)`

**Acceptance Criteria**:
- The read operation correctly retrieves the value from `SessionItem.Key`
- Nil values from the API response are handled gracefully (the provider SHALL NOT call `d.Set` when the API returns nil)

#### Scenario: Update the key parameter in place

- **WHEN** a user updates the `key` parameter on an existing `tencentcloud_waf_cc_session` resource
- **THEN** the provider SHALL send the updated value in `request.Key` when calling the `UpsertSession` API
- **AND** the resource SHALL be updated in-place without recreation (no ForceNew)

**Acceptance Criteria**:
- In-place update is performed (no ForceNew)
- The `UpsertSession` API is called with the updated `Key`
- Terraform plan shows only the changed attribute

#### Scenario: Omit the key parameter

- **WHEN** a user does not specify the `key` parameter in the configuration
- **THEN** the provider SHALL NOT set `request.Key` (the API default applies)
- **AND** no error SHALL be raised

**Acceptance Criteria**:
- The parameter is optional
- Omitting the parameter does not cause errors
- The API default behavior is preserved

#### Scenario: Import an existing WAF CC session with key configured

- **WHEN** a WAF CC session exists in TencentCloud with `Key` configured and the resource is imported or refreshed
- **THEN** the `key` value SHALL be correctly read from the API response and stored in Terraform state

**Acceptance Criteria**:
- Read operation retrieves the value from `DescribeSession` API
- Value is set in Terraform state via `d.Set`
- Import functionality works correctly
- Nil values are handled gracefully

---

## MODIFIED Requirements

None. This change adds a new parameter without modifying existing requirements.

---

## REMOVED Requirements

None. This is a pure addition with no deprecations or removals.
