## Why

Terraform Provider for TencentCloud currently lacks a datasource to query TEO (TencentCloud EdgeOne) default certificates. Users need to read default certificate information (such as cert ID, status, expiry, algorithm, etc.) for a given zone to reference in other Terraform resources or for operational visibility. The cloud API `DescribeDefaultCertificates` already supports this query, but no corresponding Terraform datasource exists.

## What Changes

- Add a new datasource `tencentcloud_teo_default_certificate` that calls the `DescribeDefaultCertificates` API to query default certificate list for a TEO zone
- The datasource supports filtering by `zone_id` (required) and optional `filters` for advanced filtering
- Output includes `default_server_cert_info` list with certificate details: `cert_id`, `alias`, `type`, `expire_time`, `effective_time`, `common_name`, `subject_alt_name`, `status`, `message`, `sign_algo`
- Register the new datasource in `provider.go` and `provider.md`
- Add corresponding service layer method `DescribeTeoDefaultCertificatesByFilter` in the TEO service
- Add unit tests using gomonkey mock approach
- Add example `.md` documentation file

## Capabilities

### New Capabilities
- `teo-default-certificate-datasource`: Provides a Terraform datasource to query TEO default certificates via `DescribeDefaultCertificates` API, enabling users to read certificate attributes for a given zone

### Modified Capabilities
<!-- No existing capabilities are being modified -->

## Impact

- **New files**: `tencentcloud/services/teo/data_source_tc_teo_default_certificate.go`, `tencentcloud/services/teo/data_source_tc_teo_default_certificate_test.go`, `tencentcloud/services/teo/data_source_tc_teo_default_certificate.md`
- **Modified files**: `tencentcloud/services/teo/service_tencentcloud_teo.go` (add service method), `tencentcloud/provider.go` (register datasource), `tencentcloud/provider.md` (document datasource)
- **Cloud API**: Uses `DescribeDefaultCertificates` from `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901`
- **No breaking changes**: This is a purely additive change
