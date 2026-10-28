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

### Requirement: User can configure message_rate on mqtt instance

The `tencentcloud_mqtt_instance` resource SHALL expose `message_rate` as an optional, configurable integer parameter representing the single client message send/receive rate limit (messages/second). When specified, the value SHALL be sent to the `ModifyInstance` API to update the instance configuration. When not specified, the server-side default SHALL be used and read back into state.

#### Scenario: User sets message_rate
- **WHEN** user specifies `message_rate = 100` in the resource configuration
- **THEN** the Update function SHALL set `request.MessageRate` to `100` and call the `ModifyInstance` API

#### Scenario: User does not specify message_rate
- **WHEN** user does not include `message_rate` in the resource configuration
- **THEN** the resource SHALL read the server-assigned default value from the `DescribeInstance` API response and store it in state

#### Scenario: User modifies message_rate
- **WHEN** user changes `message_rate` from `100` to `200` in an existing resource
- **THEN** the Update function SHALL detect the change, set `request.MessageRate` to `200`, and call the `ModifyInstance` API without recreating the instance

### Requirement: User can configure use_default_server_cert on mqtt instance

The `tencentcloud_mqtt_instance` resource SHALL expose `use_default_server_cert` as an optional, configurable boolean parameter indicating whether to use the default server certificate. When specified, the value SHALL be sent to the `ModifyInstance` API to update the instance configuration. When not specified, the server-side default SHALL be used and read back into state.

#### Scenario: User sets use_default_server_cert to true
- **WHEN** user specifies `use_default_server_cert = true` in the resource configuration
- **THEN** the Update function SHALL set `request.UseDefaultServerCert` to `true` and call the `ModifyInstance` API

#### Scenario: User sets use_default_server_cert to false
- **WHEN** user specifies `use_default_server_cert = false` in the resource configuration
- **THEN** the Update function SHALL set `request.UseDefaultServerCert` to `false` and call the `ModifyInstance` API

#### Scenario: User does not specify use_default_server_cert
- **WHEN** user does not include `use_default_server_cert` in the resource configuration
- **THEN** the resource SHALL read the server-assigned default value from the `DescribeInstance` API response and store it in state

### Requirement: Schema fields are backward compatible

The schema changes SHALL NOT cause plan diffs or state corruption for existing users who do not specify the new fields.

#### Scenario: Existing resource without explicit device_certificate_provision_type
- **WHEN** an existing `tencentcloud_mqtt_instance` resource has `device_certificate_provision_type` in state (from prior Computed read) and the user has not added the field to their configuration
- **THEN** Terraform plan SHALL show no changes for this field

#### Scenario: Existing resource without explicit message_rate or use_default_server_cert
- **WHEN** an existing `tencentcloud_mqtt_instance` resource does not have `message_rate` or `use_default_server_cert` in the user configuration
- **THEN** Terraform plan SHALL show no changes for these fields (values are read from server)
