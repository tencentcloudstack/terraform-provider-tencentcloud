## Context

The `tencentcloud_teo_certificate_config` resource is a RESOURCE_KIND_CONFIG resource that manages EdgeOne certificate configuration for acceleration domains. The resource currently supports `zone_id`, `host`, `mode`, `server_cert_info`, and `upstream_cert_info` parameters. The cloud API `ModifyHostsCertificate` already supports the `ClientCertInfo` parameter (type `MutualTLS`) for edge mutual TLS authentication, but it has not been exposed in the Terraform resource.

The read operation uses `DescribeAccelerationDomains` API, which returns `AccelerationDomain.Certificate.ClientCertInfo` of type `MutualTLS`. The `MutualTLS` struct contains `Switch` (string, on/off) and `CertInfos` ([]*CertificateInfo). Each `CertificateInfo` contains `CertId`, `Alias`, `Type`, `ExpireTime`, `DeployTime`, `SignAlgo`, `Status`.

When `ClientCertInfo` is used as an input parameter in `ModifyHostsCertificate`, only the `CertId` field of `CertificateInfo` needs to be provided for each certificate entry.

## Goals / Non-Goals

**Goals:**
- Add `client_cert_info` parameter to `tencentcloud_teo_certificate_config` resource schema
- Support creating, reading, and updating `client_cert_info` through the existing `ModifyHostsCertificate` API
- Maintain backward compatibility - existing configurations without `client_cert_info` must continue to work

**Non-Goals:**
- Not adding the deprecated `ApplyType` parameter (marked as deprecated in the SDK)
- Not modifying the existing `host` parameter to support the `Hosts` array format (that would be a separate change)
- Not adding support for any other missing API parameters beyond `client_cert_info`

## Decisions

1. **Schema design for `client_cert_info`**: Use `TypeList` with `MaxItems: 1` and `Optional+Computed`, mirroring the structure of `upstream_cert_info`. This allows the parameter to be optional while still being populated from the API response during read operations.

2. **Sub-fields of `client_cert_info`**:
   - `switch` (TypeString, Required): The mutual TLS on/off switch
   - `cert_infos` (TypeList, Optional+Computed): Certificate list. As input, only `cert_id` is needed; as output, all `CertificateInfo` fields are populated

3. **Input-only vs full schema for cert_infos**: When sending to `ModifyHostsCertificate`, only `cert_id` is required in `CertInfos`. However, the schema includes all `CertificateInfo` fields (`cert_id`, `alias`, `type`, `expire_time`, `deploy_time`, `sign_algo`, `status`) to properly read back the full response. Only `cert_id` is marked as Required; the rest are Computed.

4. **Update logic**: The `client_cert_info` should be included in the `ModifyHostsCertificate` request when the parameter is set. Since this is a RESOURCE_KIND_CONFIG resource, the update method calls `ModifyHostsCertificate` with all relevant parameters.

5. **Read logic**: Parse `ClientCertInfo` from `AccelerationDomain.Certificate.ClientCertInfo` in the existing `DescribeAccelerationDomains` response, mapping `MutualTLS.Switch` → `switch` and `MutualTLS.CertInfos` → `cert_infos`.

## Risks / Trade-offs

- [Backward compatibility] Adding `client_cert_info` as Optional+Computed ensures existing configurations are not broken → No migration needed
- [API behavior] Leaving `ClientCertInfo` blank in `ModifyHostsCertificate` means retaining the original configuration per API docs → Safe default behavior
- [Feature availability] The API docs note that edge mutual TLS is in beta testing → Users need to contact Tencent Cloud to enable the feature, but this doesn't affect the Terraform implementation
