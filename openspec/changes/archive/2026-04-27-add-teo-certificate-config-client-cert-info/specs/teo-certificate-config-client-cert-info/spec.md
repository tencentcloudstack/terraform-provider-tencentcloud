## ADDED Requirements

### Requirement: client_cert_info schema parameter
The `tencentcloud_teo_certificate_config` resource SHALL include a `client_cert_info` parameter of type `TypeList` with `MaxItems: 1`, `Optional: true`, and `Computed: true`. This parameter SHALL contain:
- `switch` (TypeString, Optional+Computed): Mutual TLS configuration switch, values SHALL be `on` or `off`
- `cert_infos` (TypeList, Optional+Computed): Mutual TLS certificate list, where each item SHALL contain:
  - `cert_id` (TypeString, Optional+Computed): Certificate ID from SSL
  - `alias` (TypeString, Computed): Certificate alias
  - `type` (TypeString, Computed): Certificate type (default/upload/managed)
  - `expire_time` (TypeString, Computed): Certificate expiration time
  - `deploy_time` (TypeString, Computed): Certificate deployment time
  - `sign_algo` (TypeString, Computed): Signature algorithm
  - `status` (TypeString, Computed): Certificate status

All fields within `client_cert_info` that are used for both input (ModifyHostsCertificate) and output (DescribeAccelerationDomains) SHALL be Optional+Computed.

#### Scenario: Resource schema includes client_cert_info parameter
- **WHEN** the resource schema is defined
- **THEN** `client_cert_info` SHALL be present as a TypeList with MaxItems 1, Optional and Computed

### Requirement: client_cert_info update via ModifyHostsCertificate
When `client_cert_info` is set in the resource configuration, the update method SHALL include `ClientCertInfo` in the `ModifyHostsCertificate` API request, constructing a `MutualTLS` object with `Switch` and `CertInfos` fields.

#### Scenario: Update resource with client_cert_info enabled
- **WHEN** `client_cert_info` is configured with `switch = "on"` and `cert_infos` containing certificate IDs
- **THEN** the `ModifyHostsCertificate` API SHALL be called with `ClientCertInfo` containing `Switch: "on"` and `CertInfos` with the provided `CertId` values

#### Scenario: client_cert_info not specified
- **WHEN** `client_cert_info` is not provided in the resource configuration
- **THEN** the `ModifyHostsCertificate` API SHALL NOT include `ClientCertInfo` in the request, retaining the original configuration

### Requirement: client_cert_info read from DescribeAccelerationDomains
The read method SHALL parse `ClientCertInfo` from the `AccelerationDomain.Certificate.ClientCertInfo` field in the `DescribeAccelerationDomains` API response.

#### Scenario: Read client_cert_info from API response
- **WHEN** the `DescribeAccelerationDomains` API returns `Certificate.ClientCertInfo` with data
- **THEN** the `client_cert_info` parameter SHALL be populated with all fields

#### Scenario: Read response with no client_cert_info
- **WHEN** the `DescribeAccelerationDomains` API returns `Certificate.ClientCertInfo` as nil
- **THEN** the `client_cert_info` parameter SHALL be empty in the Terraform state

### Requirement: upstream_mutual_tls field attributes
The `upstream_mutual_tls.switch` and `upstream_mutual_tls.cert_infos.cert_id` fields SHALL be Optional+Computed (not Required), since the struct is used both as API input and output.

### Requirement: upstream_cert_info.upstream_mutual_tls.cert_infos.status field
The `upstream_cert_info.upstream_mutual_tls.cert_infos` nested schema SHALL include a `status` field (TypeString, Computed).

### Requirement: upstream_certificate_verify sub-field
The `upstream_cert_info` schema SHALL include an `upstream_certificate_verify` sub-field (TypeList, MaxItems: 1, Optional+Computed), containing:
- `verification_mode` (TypeString, Optional+Computed): Origin certificate verification mode (`disable` or `custom_ca`)
- `custom_ca_certs` (TypeList, Optional+Computed): List of trusted CA certificates, each containing CertificateInfo fields (`cert_id` Optional+Computed, plus computed fields `alias`, `type`, `expire_time`, `deploy_time`, `sign_algo`, `status`)

The read method SHALL parse `UpstreamCertificateVerify` from `UpstreamCertInfo` in the API response. The update method SHALL include `UpstreamCertificateVerify` in the `ModifyHostsCertificate` request when set.

### Requirement: Backward compatibility
Adding these new fields SHALL NOT break existing Terraform configurations.
