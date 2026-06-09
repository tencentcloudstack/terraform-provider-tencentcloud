## ADDED Requirements

### Requirement: Datasource schema definition
The `tencentcloud_teo_default_certificate` datasource SHALL define the following schema:
- `zone_id` (TypeString, Optional): Zone ID for querying default certificates
- `filters` (TypeList, Optional): Filter criteria with nested schema containing `name` (TypeString, Required), `values` (TypeSet of TypeString, Required)
- `default_server_cert_info` (TypeList, Computed): List of default certificate information with nested schema containing `cert_id` (TypeString), `alias` (TypeString), `type` (TypeString), `expire_time` (TypeString), `effective_time` (TypeString), `common_name` (TypeString), `subject_alt_name` (TypeSet of TypeString), `status` (TypeString), `message` (TypeString), `sign_algo` (TypeString)
- `result_output_file` (TypeString, Optional): File path to save query results

#### Scenario: Schema with required zone_id only
- **WHEN** a user defines the datasource with only `zone_id` set
- **THEN** the datasource SHALL query all default certificates for that zone without additional filters

#### Scenario: Schema with filters
- **WHEN** a user defines the datasource with `zone_id` and `filters`
- **THEN** the datasource SHALL pass both zone_id and filters to the `DescribeDefaultCertificates` API

### Requirement: Read operation calls DescribeDefaultCertificates API
The datasource Read function SHALL call the `DescribeDefaultCertificates` API from the `teo/v20220901` package with retry logic using `tccommon.ReadRetryTimeout`. If the API call fails, the error SHALL be wrapped with `tccommon.RetryError()` and returned.

#### Scenario: Successful API call
- **WHEN** the Read function calls `DescribeDefaultCertificates` with valid `zone_id`
- **THEN** the response `DefaultServerCertInfo` SHALL be flattened into the `default_server_cert_info` schema field

#### Scenario: API call failure with retry
- **WHEN** the `DescribeDefaultCertificates` API call fails with a retryable error
- **THEN** the Read function SHALL retry the API call within `tccommon.ReadRetryTimeout`

#### Scenario: API call permanent failure
- **WHEN** the `DescribeDefaultCertificates` API call fails with a non-retryable error
- **THEN** the Read function SHALL return the wrapped error

### Requirement: Pagination handling
The service method SHALL handle pagination for the `DescribeDefaultCertificates` API. The API supports `Offset` and `Limit` parameters with a maximum `Limit` of 100. The service method SHALL automatically paginate to collect all results.

#### Scenario: Results fit in single page
- **WHEN** the total number of default certificates is less than or equal to 100
- **THEN** the service method SHALL return results from a single API call with Limit=100

#### Scenario: Results span multiple pages
- **WHEN** the total number of default certificates exceeds 100
- **THEN** the service method SHALL make subsequent API calls with incremented Offset until all results are collected

### Requirement: Response flattening
The Read function SHALL flatten each `DefaultServerCertInfo` object from the API response into a map with snake_case keys: `cert_id`, `alias`, `type`, `expire_time`, `effective_time`, `common_name`, `subject_alt_name`, `status`, `message`, `sign_algo`.

#### Scenario: All certificate fields populated
- **WHEN** the API returns a `DefaultServerCertInfo` with all fields populated
- **THEN** the flattened map SHALL contain all fields with their corresponding values

#### Scenario: Some certificate fields nil
- **WHEN** the API returns a `DefaultServerCertInfo` with some nil fields
- **THEN** the flattened map SHALL omit the nil fields from the map entry

### Requirement: Datasource ID generation
The datasource SHALL set its ID using `helper.DataResourceIdsHash()` with the list of `cert_id` values collected from the response. If no certificates are found, the ID SHALL be set using `helper.DataResourceIdsHash()` with an empty list.

#### Scenario: Certificates found
- **WHEN** the API returns one or more default certificates
- **THEN** the datasource ID SHALL be a hash of all `cert_id` values

#### Scenario: No certificates found
- **WHEN** the API returns no default certificates
- **THEN** the datasource ID SHALL be a hash of an empty list

### Requirement: Result output file
The datasource SHALL support the `result_output_file` parameter. When specified, the query results SHALL be written to the specified file path using `tccommon.WriteToFile()`.

#### Scenario: Output file specified
- **WHEN** the `result_output_file` parameter is set to a valid file path
- **THEN** the datasource SHALL write the certificate list to that file

#### Scenario: No output file specified
- **WHEN** the `result_output_file` parameter is not set
- **THEN** the datasource SHALL skip file writing

### Requirement: Provider registration
The `tencentcloud_teo_default_certificate` datasource SHALL be registered in `provider.go` in the datasources map and listed in `provider.md` under the Data Sources section.

#### Scenario: Provider recognizes the datasource
- **WHEN** a Terraform configuration references `data.tencentcloud_teo_default_certificate`
- **THEN** the provider SHALL resolve and execute the datasource

### Requirement: Service layer method
A new service method `DescribeTeoDefaultCertificatesByFilter` SHALL be added to the `TeoService` struct in `service_tencentcloud_teo.go`. This method SHALL accept a `paramMap` parameter, construct the `DescribeDefaultCertificates` request with zone_id, filters, offset, and limit, handle pagination, and return `[]*teo.DefaultServerCertInfo`.

#### Scenario: Service method with zone_id only
- **WHEN** the paramMap contains only `ZoneId`
- **THEN** the service method SHALL query all default certificates for that zone

#### Scenario: Service method with zone_id and filters
- **WHEN** the paramMap contains `ZoneId` and `Filters`
- **THEN** the service method SHALL pass both to the API request

### Requirement: Unit tests with gomonkey mock
Unit tests SHALL be created in `data_source_tc_teo_default_certificate_test.go` using gomonkey to mock the cloud API calls. Tests SHALL cover the Read function with successful API response and empty API response scenarios.

#### Scenario: Successful read test
- **WHEN** the Read function is called and the mocked API returns certificate data
- **THEN** the test SHALL verify that `default_server_cert_info` is correctly populated

#### Scenario: Empty result test
- **WHEN** the Read function is called and the mocked API returns no certificates
- **THEN** the test SHALL verify that `default_server_cert_info` is empty
