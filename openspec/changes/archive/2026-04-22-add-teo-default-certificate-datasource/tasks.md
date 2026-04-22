## 1. Service Layer

- ~~1.1 Add `DescribeTeoDefaultCertificatesByFilter` method to `TeoService`~~ Removed: API is called directly in the datasource Read function.

## 2. Datasource Schema and Read Function

- [x] 2.1 Create `tencentcloud/services/teo/data_source_tc_teo_default_certificate.go` with `DataSourceTencentCloudTeoDefaultCertificate()` function defining schema: `zone_id` (Optional), `filters` (Optional), `default_server_cert_info` (Computed), `result_output_file` (Optional). Both `zone_id` and `filters` are Optional with runtime validation requiring at least one.
- [x] 2.2 Implement `dataSourceTencentCloudTeoDefaultCertificateRead()` function that validates at least one of `zone_id` or `filters` is specified, calls `DescribeDefaultCertificates` API directly with pagination (Limit=100), flattens response into `default_server_cert_info` list, sets ID with `helper.DataResourceIdsHash()`, and handles `result_output_file`

## 3. Provider Registration

- [x] 3.1 Register `tencentcloud_teo_default_certificate` datasource in `tencentcloud/provider.go` datasources map
- [x] 3.2 Add `tencentcloud_teo_default_certificate` entry in `tencentcloud/provider.md` Data Sources section

## 4. Documentation

- [x] 4.1 Create `tencentcloud/services/teo/data_source_tc_teo_default_certificate.md` with description, example usage showing zone_id and filters, and output attributes

## 5. Unit Tests

- [x] 5.1 Create `tencentcloud/services/teo/data_source_tc_teo_default_certificate_test.go` with gomonkey mock tests covering successful read and empty result scenarios
- [x] 5.2 Run unit tests to verify correctness
