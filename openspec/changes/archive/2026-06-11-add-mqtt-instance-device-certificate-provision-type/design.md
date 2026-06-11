## Context

The `tencentcloud_mqtt_instance` resource already has the `device_certificate_provision_type` field defined in its schema as `Computed: true` and reads it from the `DescribeInstance` API response. The `ModifyInstance` API supports `DeviceCertificateProvisionType` as an input parameter, but the current Update function does not set it on the request despite the field being listed in `mutableArgs`. The `CreateInstance` API does not support this parameter.

Current state in `resource_tc_mqtt_instance.go`:
- Schema (line 96-100): `Computed: true` only
- Read (line 352-353): Already reads from `DescribeInstance` response
- Update (line 402): Listed in `mutableArgs` but missing the `request.DeviceCertificateProvisionType = ...` assignment

## Goals / Non-Goals

**Goals:**
- Make `device_certificate_provision_type` configurable by users via Terraform
- Support modification of the field through the `ModifyInstance` API
- Maintain backward compatibility with existing Terraform state and configurations

**Non-Goals:**
- Supporting `device_certificate_provision_type` at creation time (the `CreateInstance` API does not accept this parameter)
- Modifying the `DescribeInstance` read logic (already correctly implemented)
- Adding validation for allowed values (JITP/API) — defer to API-side validation

## Decisions

1. **Schema change: `Optional + Computed` instead of just `Optional`**
   - Rationale: The field has a server-side default value. Using `Optional + Computed` allows Terraform to read the server default when the user does not specify a value, while still allowing explicit configuration.
   - Alternative: Making it `Required` was rejected because existing resources already have this field set server-side and forcing users to specify it would be a breaking change.

2. **No `ForceNew` — use in-place update via `ModifyInstance`**
   - Rationale: The `ModifyInstance` API supports changing this field without recreating the instance. Using `ForceNew` would unnecessarily destroy and recreate instances.

3. **Use `d.GetOk` pattern for setting the request field**
   - Rationale: Consistent with how other string fields (e.g., `x509_mode`, `name`) are handled in the same Update function. Only set the field when the user has explicitly provided a value.

4. **No creation-time support**
   - Rationale: The `CreateInstance` API in the SDK (`github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/mqtt/v20240516`) does not include `DeviceCertificateProvisionType` in its request struct. Users who need a non-default value must set it after creation via an update (Terraform will handle this automatically on the first apply if the field is specified).

## Risks / Trade-offs

- [Risk] Users specify `device_certificate_provision_type` in their config but it only takes effect after creation (requires a second API call) → Mitigation: This is standard Terraform behavior for fields not supported at creation time. The resource will converge on the second `terraform apply` or within the same apply if Terraform detects drift during the read-after-create step.
- [Risk] Changing `Computed: true` to `Optional + Computed` could cause plan diffs for existing users → Mitigation: Since the field was already `Computed`, existing state already has the value. Adding `Optional` does not change the stored state value, so no spurious diffs will occur.
