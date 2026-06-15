## Context

The `tencentcloud_mqtt_instance` resource already has the `device_certificate_provision_type` field defined in its schema as `Computed: true` and reads it from the `DescribeInstance` API response. The `ModifyInstance` API supports `DeviceCertificateProvisionType`, `MessageRate`, and `UseDefaultServerCert` as input parameters. The `DescribeInstance` API returns all three fields in its response. The `CreateInstance` API does not support any of these parameters.

Current state in `resource_tc_mqtt_instance.go`:
- `device_certificate_provision_type`: Changed to `Optional + Computed`, fully wired in Create (post-create modify), Read, and Update.
- `message_rate`: Added as `Optional + Computed` (TypeInt), wired in Create (post-create modify), Read, and Update.
- `use_default_server_cert`: Added as `Optional + Computed` (TypeBool), wired in Create (post-create modify), Read, and Update.

## Goals / Non-Goals

**Goals:**
- Make `device_certificate_provision_type`, `message_rate`, and `use_default_server_cert` configurable by users via Terraform
- Support modification of these fields through the `ModifyInstance` API
- Read these fields from the `DescribeInstance` API response
- Maintain backward compatibility with existing Terraform state and configurations

**Non-Goals:**
- Supporting these parameters at creation time (the `CreateInstance` API does not accept them)
- Adding validation for allowed values — defer to API-side validation

## Decisions

1. **Schema change: `Optional + Computed` for all three fields**
   - Rationale: These fields have server-side default values. Using `Optional + Computed` allows Terraform to read the server default when the user does not specify a value, while still allowing explicit configuration.
   - Alternative: Making them `Required` was rejected because existing resources already have these fields set server-side and forcing users to specify them would be a breaking change.

2. **No `ForceNew` — use in-place update via `ModifyInstance`**
   - Rationale: The `ModifyInstance` API supports changing these fields without recreating the instance. Using `ForceNew` would unnecessarily destroy and recreate instances.

3. **Use `d.GetOk` / `d.GetOkExists` pattern for setting request fields**
   - Rationale: Consistent with how other fields are handled in the same Update function. Only set the field when the user has explicitly provided a value. Use `d.GetOkExists` for bool/int fields to correctly detect zero-value settings.

4. **Post-create modify for initial configuration**
   - Rationale: The `CreateInstance` API does not support these parameters. After instance creation and status becomes RUNNING, a `ModifyInstance` call is made to set user-specified values. This is the same pattern already used for `automatic_activation`, `authorization_policy`, `x509_mode`, and `device_certificate_provision_type`.

## Risks / Trade-offs

- [Risk] Users specify these fields in their config but they only take effect after creation (requires a second API call) → Mitigation: This is standard Terraform behavior for fields not supported at creation time. The resource will converge within the same apply since the Create function already calls ModifyInstance after creation.
- [Risk] Adding new `Optional + Computed` fields could cause plan diffs for existing users → Mitigation: Since these are new fields not previously in state, Terraform will simply read the server values on the next refresh. No spurious diffs will occur.
