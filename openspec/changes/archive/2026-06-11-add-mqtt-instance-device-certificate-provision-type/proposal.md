## Why

The `tencentcloud_mqtt_instance` resource currently exposes `device_certificate_provision_type` as a read-only (`Computed`) attribute. However, the `ModifyInstance` API supports setting this parameter, allowing users to configure the client certificate registration method (JITP for automatic registration or API for manual registration). Additionally, the `ModifyInstance` API supports `MessageRate` (single client message send/receive rate limit) and `UseDefaultServerCert` (whether to use the default server certificate), which are also returned by `DescribeInstance` but not yet configurable via Terraform. Making these fields writable enables Terraform users to manage these instance settings as part of their infrastructure-as-code workflow.

## What Changes

- Change the `device_certificate_provision_type` schema field from `Computed: true` (read-only) to `Optional: true, Computed: true` (user-configurable with server default).
- Add `message_rate` schema field as `Optional: true, Computed: true` (TypeInt) to allow configuring single client message rate limit.
- Add `use_default_server_cert` schema field as `Optional: true, Computed: true` (TypeBool) to allow configuring whether to use the default server certificate.
- Add logic in the Update function to set `request.DeviceCertificateProvisionType`, `request.MessageRate`, and `request.UseDefaultServerCert` when the fields are provided, enabling modification via the `ModifyInstance` API.
- Add logic in the Create function (post-create modify) to set these fields via `ModifyInstance` since `CreateInstance` API does not support them.
- Read `MessageRate` and `UseDefaultServerCert` from `DescribeInstance` response in the Read function.
- Update unit tests to cover the new writable behavior.
- Update the resource documentation (.md) to reflect the fields are now configurable.

## Capabilities

### New Capabilities
- `mqtt-instance-device-cert-provision-type`: Allow users to configure the `device_certificate_provision_type` parameter on `tencentcloud_mqtt_instance` resource via the `ModifyInstance` API.
- `mqtt-instance-message-rate`: Allow users to configure the `message_rate` parameter on `tencentcloud_mqtt_instance` resource via the `ModifyInstance` API.
- `mqtt-instance-use-default-server-cert`: Allow users to configure the `use_default_server_cert` parameter on `tencentcloud_mqtt_instance` resource via the `ModifyInstance` API.

### Modified Capabilities

## Impact

- **Code**: `tencentcloud/services/mqtt/resource_tc_mqtt_instance.go` — schema additions and Create/Read/Update function modifications.
- **Tests**: `tencentcloud/services/mqtt/resource_tc_mqtt_instance_test.go` — add unit test coverage for the new writable parameters.
- **Documentation**: `tencentcloud/services/mqtt/resource_tc_mqtt_instance.md` — update to show the fields as configurable.
- **APIs**: Uses existing `ModifyInstance` and `DescribeInstance` APIs (`github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/mqtt/v20240516`). No new API dependencies.
- **Backward Compatibility**: Fully backward-compatible. Existing configurations with no `device_certificate_provision_type`, `message_rate`, or `use_default_server_cert` set will continue to work as before (server provides the default values).
