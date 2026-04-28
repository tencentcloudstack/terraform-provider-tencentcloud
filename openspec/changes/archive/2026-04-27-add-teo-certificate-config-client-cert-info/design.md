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
   - `switch` (TypeString, Optional+Computed): The mutual TLS on/off switch
   - `cert_infos` (TypeList, Optional+Computed): Certificate list. As input, only `cert_id` is needed; as output, all `CertificateInfo` fields are populated

3. **All fields Optional+Computed**: Since `client_cert_info` (type `MutualTLS`) is used both as input in `ModifyHostsCertificate` and output in `DescribeAccelerationDomains`, `switch` and `cert_id` are set to Optional+Computed rather than Required. The same applies to `upstream_mutual_tls.switch` and `upstream_mutual_tls.cert_infos.cert_id`.

4. **Update logic**: The `client_cert_info` should be included in the `ModifyHostsCertificate` request when the parameter is set. Since this is a RESOURCE_KIND_CONFIG resource, the update method calls `ModifyHostsCertificate` with all relevant parameters.

5. **Read logic**: Parse `ClientCertInfo` from `AccelerationDomain.Certificate.ClientCertInfo` in the existing `DescribeAccelerationDomains` response, mapping `MutualTLS.Switch` → `switch` and `MutualTLS.CertInfos` → `cert_infos`.

6. **Upstream cert_infos status field**: Add `status` as a Computed field to `upstream_cert_info.upstream_mutual_tls.cert_infos` schema and read it from the API response, matching the `CertificateInfo.Status` field in the SDK. Values include `deployed`, `processing`, `applying`, `failed`, `issued`.

7. **Add `upstream_certificate_verify`**: The `UpstreamCertInfo` SDK struct also contains `UpstreamCertificateVerify *OriginCertificateVerify`. Add `upstream_certificate_verify` sub-field to `upstream_cert_info` schema with `verification_mode` (Optional+Computed) and `custom_ca_certs` (TypeList of CertificateInfo, Optional+Computed). Read and write logic mirrors the existing `upstream_mutual_tls` pattern.

8. **Test organization**: All unit tests are placed in `resource_tc_teo_certificate_config_test.go` using the `TestTeoCertificateConfig_` prefix.

## Risks / Trade-offs

- [Backward compatibility] Adding `client_cert_info` as Optional+Computed ensures existing configurations are not broken → No migration needed
- [API behavior] Leaving `ClientCertInfo` blank in `ModifyHostsCertificate` means retaining the original configuration per API docs → Safe default behavior
- [Feature availability] The API docs note that edge mutual TLS is in beta testing → Users need to contact Tencent Cloud to enable the feature, but this doesn't affect the Terraform implementation
