## 1. Service Layer

- [x] 1.1 Add `DescribeTeoDefaultCertificatesByFilter` method to `TeoService` in `tencentcloud/services/teo/service_tencentcloud_teo.go` that calls `DescribeDefaultCertificates` API with pagination support (Limit=100), accepting paramMap with ZoneId and Filters, returning `[]*teo.DefaultServerCertInfo`

## 2. Datasource Schema and Read Function

- [x] 2.1 Create `tencentcloud/services/teo/data_source_tc_teo_default_certificate.go` with `DataSourceTencentCloudTeoDefaultCertificate()` function defining schema: `zone_id` (Required), `filters` (Optional), `default_server_cert_info` (Computed), `result_output_file` (Optional)
- [x] 2.2 Implement `dataSourceTencentCloudTeoDefaultCertificateRead()` function that builds paramMap from schema, calls service method with retry, flattens `DefaultServerCertInfo` response into `default_server_cert_info` list, sets ID with `helper.DataResourceIdsHash()`, and handles `result_output_file`

## 3. Provider Registration

- [x] 3.1 Register `tencentcloud_teo_default_certificate` datasource in `tencentcloud/provider.go` datasources map
- [x] 3.2 Add `tencentcloud_teo_default_certificate` entry in `tencentcloud/provider.md` Data Sources section

## 4. Documentation

- [x] 4.1 Create `tencentcloud/services/teo/data_source_tc_teo_default_certificate.md` with description, example usage showing zone_id and filters, and output attributes

## 5. Unit Tests

- [x] 5.1 Create `tencentcloud/services/teo/data_source_tc_teo_default_certificate_test.go` with gomonkey mock tests covering successful read and empty result scenarios
- [x] 5.2 Run unit tests to verify correctness
