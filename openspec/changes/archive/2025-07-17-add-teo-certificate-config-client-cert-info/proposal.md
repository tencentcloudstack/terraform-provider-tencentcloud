## Why

The `tencentcloud_teo_certificate_config` resource currently supports server certificate configuration (`server_cert_info`), mode selection (`mode`), and upstream certificate configuration (`upstream_cert_info`), but lacks support for client certificate configuration (`client_cert_info`). The cloud API `ModifyHostsCertificate` already supports the `ClientCertInfo` parameter (type `MutualTLS`) for edge mutual TLS authentication, where client CA certificates are deployed on EO nodes for client-to-EO-node authentication. This parameter needs to be exposed in the Terraform resource so users can manage edge mutual TLS authentication through Infrastructure as Code.

## What Changes

- Add `client_cert_info` parameter (TypeList, MaxItems: 1, Optional+Computed) to the `tencentcloud_teo_certificate_config` resource schema, corresponding to the `ClientCertInfo` field of the `ModifyHostsCertificate` API
- The `client_cert_info` parameter contains:
  - `switch` (TypeString, Required): Mutual TLS configuration switch, values: `on`/`off`
  - `cert_infos` (TypeList, Optional+Computed): Mutual TLS certificate list, with each item containing `cert_id` (Required) and computed fields (`alias`, `type`, `expire_time`, `deploy_time`, `sign_algo`, `status`)
- Update the resource's update logic to include `client_cert_info` in the `ModifyHostsCertificate` API call
- Update the resource's read logic to parse `client_cert_info` from the `DescribeAccelerationDomains` API response
- Update the resource's `.md` documentation file with the new parameter example

## Capabilities

### New Capabilities
- `teo-certificate-config-client-cert-info`: Add client certificate (MutualTLS) configuration support to the teo_certificate_config resource, enabling edge mutual TLS authentication management

### Modified Capabilities

## Impact

- `tencentcloud/services/teo/resource_tc_teo_certificate_config.go`: Schema definition, CRUD logic
- `tencentcloud/services/teo/resource_tc_teo_certificate_config.md`: Documentation
- `tencentcloud/services/teo/resource_tc_teo_certificate_config_test.go`: Unit tests
- Cloud API: `ModifyHostsCertificate` (update), `DescribeAccelerationDomains` (read)
