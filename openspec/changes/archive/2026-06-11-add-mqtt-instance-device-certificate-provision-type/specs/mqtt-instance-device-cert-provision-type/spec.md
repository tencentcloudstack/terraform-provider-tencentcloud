## ADDED Requirements

### Requirement: User can configure device_certificate_provision_type on mqtt instance

The `tencentcloud_mqtt_instance` resource SHALL expose `device_certificate_provision_type` as an optional, configurable parameter. The field accepts string values `JITP` (automatic registration) or `API` (manual registration). When specified, the value SHALL be sent to the `ModifyInstance` API to update the instance configuration. When not specified, the server-side default SHALL be used and read back into state.

#### Scenario: User sets device_certificate_provision_type to JITP
- **WHEN** user specifies `device_certificate_provision_type = "JITP"` in the resource configuration
- **THEN** the Update function SHALL set `request.DeviceCertificateProvisionType` to `"JITP"` and call the `ModifyInstance` API

#### Scenario: User sets device_certificate_provision_type to API
- **WHEN** user specifies `device_certificate_provision_type = "API"` in the resource configuration
- **THEN** the Update function SHALL set `request.DeviceCertificateProvisionType` to `"API"` and call the `ModifyInstance` API

#### Scenario: User does not specify device_certificate_provision_type
- **WHEN** user does not include `device_certificate_provision_type` in the resource configuration
- **THEN** the resource SHALL read the server-assigned default value from the `DescribeInstance` API response and store it in state

#### Scenario: User modifies device_certificate_provision_type from JITP to API
- **WHEN** user changes `device_certificate_provision_type` from `"JITP"` to `"API"` in an existing resource
- **THEN** the Update function SHALL detect the change, set `request.DeviceCertificateProvisionType` to `"API"`, and call the `ModifyInstance` API without recreating the instance

### Requirement: Schema field is backward compatible

The schema change from `Computed: true` to `Optional: true, Computed: true` SHALL NOT cause plan diffs or state corruption for existing users who do not specify the field.

#### Scenario: Existing resource without explicit device_certificate_provision_type
- **WHEN** an existing `tencentcloud_mqtt_instance` resource has `device_certificate_provision_type` in state (from prior Computed read) and the user has not added the field to their configuration
- **THEN** Terraform plan SHALL show no changes for this field
