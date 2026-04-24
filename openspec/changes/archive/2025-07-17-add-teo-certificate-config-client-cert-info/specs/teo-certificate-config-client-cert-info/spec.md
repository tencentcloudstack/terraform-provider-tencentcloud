## ADDED Requirements

### Requirement: client_cert_info schema parameter
The `tencentcloud_teo_certificate_config` resource SHALL include a `client_cert_info` parameter of type `TypeList` with `MaxItems: 1`, `Optional: true`, and `Computed: true`. This parameter SHALL contain:
- `switch` (TypeString, Required): Mutual TLS configuration switch, values SHALL be `on` or `off`
- `cert_infos` (TypeList, Optional+Computed): Mutual TLS certificate list, where each item SHALL contain:
  - `cert_id` (TypeString, Required): Certificate ID from SSL
  - `alias` (TypeString, Computed): Certificate alias
  - `type` (TypeString, Computed): Certificate type (default/upload/managed)
  - `expire_time` (TypeString, Computed): Certificate expiration time
  - `deploy_time` (TypeString, Computed): Certificate deployment time
  - `sign_algo` (TypeString, Computed): Signature algorithm
  - `status` (TypeString, Computed): Certificate status

#### Scenario: Resource schema includes client_cert_info parameter
- **WHEN** the resource schema is defined
- **THEN** `client_cert_info` SHALL be present as a TypeList with MaxItems 1, Optional and Computed

#### Scenario: client_cert_info switch is required within the block
- **WHEN** a user specifies `client_cert_info`
- **THEN** the `switch` field within `client_cert_info` MUST be provided with value `on` or `off`

### Requirement: client_cert_info update via ModifyHostsCertificate
When `client_cert_info` is set in the resource configuration, the update method SHALL include `ClientCertInfo` in the `ModifyHostsCertificate` API request, constructing a `MutualTLS` object with `Switch` and `CertInfos` fields. When used as input, each `CertificateInfo` entry SHALL only require `CertId`.

#### Scenario: Update resource with client_cert_info enabled
- **WHEN** `client_cert_info` is configured with `switch = "on"` and `cert_infos` containing certificate IDs
- **THEN** the `ModifyHostsCertificate` API SHALL be called with `ClientCertInfo` containing `Switch: "on"` and `CertInfos` with the provided `CertId` values

#### Scenario: Update resource with client_cert_info disabled
- **WHEN** `client_cert_info` is configured with `switch = "off"`
- **THEN** the `ModifyHostsCertificate` API SHALL be called with `ClientCertInfo` containing `Switch: "off"`

#### Scenario: client_cert_info not specified
- **WHEN** `client_cert_info` is not provided in the resource configuration
- **THEN** the `ModifyHostsCertificate` API SHALL NOT include `ClientCertInfo` in the request, retaining the original configuration

### Requirement: client_cert_info read from DescribeAccelerationDomains
The read method SHALL parse `ClientCertInfo` from the `AccelerationDomain.Certificate.ClientCertInfo` field in the `DescribeAccelerationDomains` API response, mapping `MutualTLS.Switch` to `switch` and `MutualTLS.CertInfos` to `cert_infos` with all computed fields populated.

#### Scenario: Read client_cert_info from API response
- **WHEN** the `DescribeAccelerationDomains` API returns `Certificate.ClientCertInfo` with `Switch: "on"` and `CertInfos` containing certificate details
- **THEN** the `client_cert_info` parameter SHALL be populated in the Terraform state with `switch` and `cert_infos` including all computed fields (alias, type, expire_time, deploy_time, sign_algo, status)

#### Scenario: Read response with no client_cert_info
- **WHEN** the `DescribeAccelerationDomains` API returns `Certificate.ClientCertInfo` as nil
- **THEN** the `client_cert_info` parameter SHALL be empty in the Terraform state

### Requirement: Backward compatibility
Adding `client_cert_info` SHALL NOT break existing Terraform configurations that do not use this parameter. Existing resources without `client_cert_info` SHALL continue to function without any changes.

#### Scenario: Existing configuration without client_cert_info
- **WHEN** an existing Terraform configuration does not include `client_cert_info`
- **THEN** the resource SHALL work as before, with no changes to create, read, update, or delete behavior
