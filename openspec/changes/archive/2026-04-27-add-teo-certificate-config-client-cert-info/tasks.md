## 1. Schema Definition

- [x] 1.1 Add `client_cert_info` parameter to the resource schema with TypeList, MaxItems: 1, Optional+Computed, containing `switch` (Optional+Computed) and `cert_infos` (TypeList, Optional+Computed, containing `cert_id` Optional+Computed + `alias`/`type`/`expire_time`/`deploy_time`/`sign_algo`/`status` Computed)
- [x] 1.2 Add `status` computed field to `upstream_cert_info.upstream_mutual_tls.cert_infos` schema
- [x] 1.3 Change `upstream_mutual_tls.switch` and `upstream_mutual_tls.cert_infos.cert_id` from Required to Optional+Computed
- [x] 1.4 Add `upstream_certificate_verify` sub-field to `upstream_cert_info` schema, containing `verification_mode` (Optional+Computed) and `custom_ca_certs` (TypeList of CertificateInfo with cert_id Optional+Computed + computed fields)

## 2. CRUD Implementation

- [x] 2.1 Update `resourceTencentCloudTeoCertificateConfigUpdateOnStart` to include `ClientCertInfo` in `ModifyHostsCertificate` API request
- [x] 2.2 Update `resourceTencentCloudTeoCertificateConfigReadPostHandleResponse0` to parse `ClientCertInfo` from API response
- [x] 2.3 Update read handler to read `status` field from `upstream_cert_info.upstream_mutual_tls.cert_infos`
- [x] 2.4 Add read logic for `upstream_certificate_verify` from `UpstreamCertInfo.UpstreamCertificateVerify`
- [x] 2.5 Add write logic for `upstream_certificate_verify` in `ModifyHostsCertificate` request

## 3. Unit Tests

- [x] 3.1 Unit tests in `resource_tc_teo_certificate_config_test.go` with `TestTeoCertificateConfig_` prefix, covering: schema validation (Optional+Computed for switch/cert_id), create/update/read with client_cert_info, backward compatibility

## 4. Documentation

- [x] 4.1 Update `resource_tc_teo_certificate_config.md` with client_cert_info example
