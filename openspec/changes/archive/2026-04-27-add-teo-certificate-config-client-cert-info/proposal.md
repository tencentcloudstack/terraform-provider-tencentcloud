## Why

The `tencentcloud_teo_certificate_config` resource currently supports server certificate configuration (`server_cert_info`), mode selection (`mode`), and upstream certificate configuration (`upstream_cert_info`), but lacks support for client certificate configuration (`client_cert_info`). The cloud API `ModifyHostsCertificate` already supports the `ClientCertInfo` parameter (type `MutualTLS`) for edge mutual TLS authentication, where client CA certificates are deployed on EO nodes for client-to-EO-node authentication. This parameter needs to be exposed in the Terraform resource so users can manage edge mutual TLS authentication through Infrastructure as Code.

## What Changes

- Add `client_cert_info` parameter (TypeList, MaxItems: 1, Optional+Computed) to the `tencentcloud_teo_certificate_config` resource schema, corresponding to the `ClientCertInfo` field of the `ModifyHostsCertificate` API
- The `client_cert_info` parameter contains:
  - `switch` (TypeString, Optional+Computed): Mutual TLS configuration switch, values: `on`/`off`
  - `cert_infos` (TypeList, Optional+Computed): Mutual TLS certificate list, with each item containing `cert_id` (Optional+Computed) and computed fields (`alias`, `type`, `expire_time`, `deploy_time`, `sign_algo`, `status`)
- All fields in `client_cert_info` are Optional+Computed because the struct is used both as input (ModifyHostsCertificate) and output (DescribeAccelerationDomains)
- Update the resource's update logic to include `client_cert_info` in the `ModifyHostsCertificate` API call
- Update the resource's read logic to parse `client_cert_info` from the `DescribeAccelerationDomains` API response
- Add `status` computed field to `upstream_cert_info.upstream_mutual_tls.cert_infos`
- Change `upstream_mutual_tls.switch` and `upstream_mutual_tls.cert_infos.cert_id` from Required to Optional+Computed
- Add `upstream_certificate_verify` sub-field to `upstream_cert_info`, containing `verification_mode` (Optional+Computed) and `custom_ca_certs` (list of CertificateInfo, Optional+Computed)
- Update the resource's `.md` documentation file with the new parameter example

## Capabilities

### New Capabilities
- `teo-certificate-config-client-cert-info`: Add client certificate (MutualTLS) configuration support, add `upstream_certificate_verify` to upstream_cert_info, add `status` field to cert_infos, and fix switch/cert_id to Optional+Computed

### Modified Capabilities

## Impact

- `tencentcloud/services/teo/resource_tc_teo_certificate_config.go`: Schema definition (client_cert_info + upstream status)
- `tencentcloud/services/teo/resource_tc_teo_certificate_config_extension.go`: Read/Update logic
- `tencentcloud/services/teo/resource_tc_teo_certificate_config.md`: Documentation
- `tencentcloud/services/teo/resource_tc_teo_certificate_config_test.go`: Unit tests (using `TestTeoCertificateConfig_` prefix)
- Cloud API: `ModifyHostsCertificate` (update), `DescribeAccelerationDomains` (read)
