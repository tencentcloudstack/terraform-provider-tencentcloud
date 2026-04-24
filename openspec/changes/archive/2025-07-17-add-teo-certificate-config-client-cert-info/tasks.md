## 1. Schema Definition

- [x] 1.1 Add `client_cert_info` parameter to the resource schema in `resource_tc_teo_certificate_config.go`, with TypeList, MaxItems: 1, Optional+Computed, containing `switch` (TypeString, Required) and `cert_infos` (TypeList, Optional+Computed, containing `cert_id` Required + `alias`/`type`/`expire_time`/`deploy_time`/`sign_algo`/`status` Computed)

## 2. CRUD Implementation

- [x] 2.1 Update `resourceTencentCloudTeoCertificateConfigUpdate` to include `ClientCertInfo` in `ModifyHostsCertificate` API request when `client_cert_info` is set, constructing `MutualTLS` object with `Switch` and `CertInfos` fields
- [x] 2.2 Update `resourceTencentCloudTeoCertificateConfigReadPostHandleResponse0` to parse `ClientCertInfo` from `AccelerationDomain.Certificate.ClientCertInfo` response, mapping `MutualTLS.Switch` → `switch` and `MutualTLS.CertInfos` → `cert_infos` with all computed fields

## 3. Unit Tests

- [x] 3.1 Add unit test cases for `client_cert_info` parameter in `resource_tc_teo_certificate_config_client_cert_info_test.go` using gomonkey mock, covering: create with client_cert_info enabled, update with client_cert_info, read with client_cert_info from API response, and backward compatibility without client_cert_info

## 4. Documentation

- [x] 4.1 Update `resource_tc_teo_certificate_config.md` to include `client_cert_info` parameter documentation with example usage showing edge mutual TLS configuration
