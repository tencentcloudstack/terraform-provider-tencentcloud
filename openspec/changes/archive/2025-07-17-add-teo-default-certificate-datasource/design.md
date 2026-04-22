## Context

Terraform Provider for TencentCloud supports TEO (TencentCloud EdgeOne) resources and datasources. Currently, there is no datasource to query default certificates for a TEO zone. The cloud API `DescribeDefaultCertificates` (package: `teo/v20220901`) supports querying default certificate information with zone ID and filters, but no corresponding Terraform datasource exists.

The existing TEO service layer already has a `DescribeTeoDefaultCertificate` method that queries a single certificate by zone ID and cert ID, but this does not support listing all default certificates or using the full filter capabilities of the API.

## Goals / Non-Goals

**Goals:**
- Provide a `tencentcloud_teo_default_certificate` datasource that allows users to query default certificates for a TEO zone
- Support the `zone_id` required parameter and optional `filters` parameter matching the `DescribeDefaultCertificates` API
- Return certificate details including `cert_id`, `alias`, `type`, `expire_time`, `effective_time`, `common_name`, `subject_alt_name`, `status`, `message`, `sign_algo`
- Follow the established pattern for TEO datasources (e.g., `data_source_tc_teo_environments.go`)
- Handle pagination automatically (API supports Offset/Limit with max Limit of 100)
- Register the datasource in `provider.go` and `provider.md`
- Add unit tests using gomonkey mock approach

**Non-Goals:**
- This is a read-only datasource; no create/update/delete operations
- No modification of existing TEO resources or datasources
- No changes to the existing `DescribeTeoDefaultCertificate` service method

## Decisions

1. **Datasource name**: `tencentcloud_teo_default_certificate` (singular, matching the API name `DescribeDefaultCertificates` which returns `DefaultServerCertInfo`). The output field `default_server_cert_info` will be a list type to accommodate multiple certificates.

2. **Service method**: Add a new `DescribeTeoDefaultCertificatesByFilter` method to the TEO service that supports the full `DescribeDefaultCertificates` API with filters and pagination, returning `[]*teo.DefaultServerCertInfo`. This follows the established pattern (e.g., `DescribeTeoEnvironmentsByFilter`).

3. **Pagination**: The API supports `Offset` and `Limit` parameters with max Limit of 100. The service method will automatically paginate to collect all results, using Limit=100 per request.

4. **Filters schema**: Use the standard TEO filter pattern with `name` (string), `values` (list of string) fields, matching the `teo.Filter` struct in the SDK.

5. **ID generation**: Use `helper.DataResourceIdsHash(certIds)` for the datasource ID, collecting `cert_id` from each returned certificate. This follows the list datasource pattern.

6. **zone_id as Required**: The API's ZoneId parameter maps to a Required schema field, consistent with other TEO datasources that query by zone.

7. **Unit testing**: Use gomonkey mock approach (not Terraform acceptance tests) per project requirements for new datasources.

## Risks / Trade-offs

- **API pagination**: If a zone has a very large number of default certificates, the auto-pagination could take time. Mitigation: The API max Limit is 100 and default certificates are typically few per zone.
- **Existing service method**: The existing `DescribeTeoDefaultCertificate` method queries a single cert. The new `DescribeTeoDefaultCertificatesByFilter` method is separate and does not affect the existing method. No risk of regression.
