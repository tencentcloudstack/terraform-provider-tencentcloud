## Why

The `tencentcloud_mqtt_instance` resource currently exposes `device_certificate_provision_type` as a read-only (`Computed`) attribute. However, the `ModifyInstance` API supports setting this parameter, allowing users to configure the client certificate registration method (JITP for automatic registration or API for manual registration). Making this field writable enables Terraform users to manage the device certificate provisioning strategy as part of their infrastructure-as-code workflow.

## What Changes

- Change the `device_certificate_provision_type` schema field from `Computed: true` (read-only) to `Optional: true, Computed: true` (user-configurable with server default).
- Add logic in the Update function to set `request.DeviceCertificateProvisionType` when the field is provided, enabling modification via the `ModifyInstance` API.
- Update unit tests to cover the new writable behavior.
- Update the resource documentation (.md) to reflect the field is now configurable.

## Capabilities

### New Capabilities
- `mqtt-instance-device-cert-provision-type`: Allow users to configure the `device_certificate_provision_type` parameter on `tencentcloud_mqtt_instance` resource via the `ModifyInstance` API.

### Modified Capabilities

## Impact

- **Code**: `tencentcloud/services/mqtt/resource_tc_mqtt_instance.go` — schema change and Update function modification.
- **Tests**: `tencentcloud/services/mqtt/resource_tc_mqtt_instance_test.go` — add unit test coverage for the new writable parameter.
- **Documentation**: `tencentcloud/services/mqtt/resource_tc_mqtt_instance.md` — update to show the field as configurable.
- **APIs**: Uses existing `ModifyInstance` API (`github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/mqtt/v20240516`). No new API dependencies.
- **Backward Compatibility**: Fully backward-compatible. Existing configurations with no `device_certificate_provision_type` set will continue to work as before (server provides the default value).
