## ADDED Requirements

### Requirement: Data source schema for tencentcloud_teo_ip_region
The system SHALL provide a Terraform data source named `tencentcloud_teo_ip_region` with the following schema:
- `ips`: Required, TypeList of TypeString, maximum 100 items, representing the list of IP addresses (IPv4/IPv6) to query.
- `ip_region_info`: Computed, TypeList of TypeResource, containing the IP region information results.
- `result_output_file`: Optional, TypeString, used to save query results to a file.

Each element in `ip_region_info` SHALL contain:
- `ip`: Computed, TypeString, the IP address queried.
- `is_edge_one_ip`: Computed, TypeString, whether the IP belongs to an EdgeOne node ("yes" or "no").

#### Scenario: Query with valid IP list
- **WHEN** user provides a list of IP addresses via `ips` parameter
- **THEN** the data source SHALL call `DescribeIPRegion` API and return the `ip_region_info` list with each IP's `ip` and `is_edge_one_ip` fields populated

#### Scenario: IP belongs to EdgeOne
- **WHEN** a queried IP belongs to an EdgeOne node
- **THEN** the `is_edge_one_ip` field for that IP SHALL be "yes"

#### Scenario: IP does not belong to EdgeOne
- **WHEN** a queried IP does not belong to an EdgeOne node
- **THEN** the `is_edge_one_ip` field for that IP SHALL be "no"

### Requirement: Provider registration for tencentcloud_teo_ip_region
The system SHALL register the `tencentcloud_teo_ip_region` data source in `provider.go` data source map and document it in `provider.md`.

#### Scenario: Data source appears in provider
- **WHEN** the Terraform provider is initialized
- **THEN** the `tencentcloud_teo_ip_region` data source SHALL be available for use in Terraform configurations

### Requirement: Service layer method for DescribeIPRegion
The system SHALL implement a `DescribeTeoIPRegionByFilter` method in the `TeoService` struct that accepts a `paramMap` and calls the `DescribeIPRegion` SDK API with retry support.

#### Scenario: API call with retry
- **WHEN** the `DescribeIPRegion` API is called and encounters a transient error
- **THEN** the system SHALL retry the request using `tccommon.ReadRetryTimeout` and return `tccommon.RetryError` on failure

#### Scenario: Successful API response
- **WHEN** the `DescribeIPRegion` API returns successfully
- **THEN** the service method SHALL return the `IPRegionInfo` list from the response

### Requirement: Unit tests with gomonkey mock
The system SHALL provide unit tests for the `tencentcloud_teo_ip_region` data source using gomonkey to mock the cloud API calls, without relying on Terraform test framework.

#### Scenario: Mock successful query
- **WHEN** the unit test mocks a successful `DescribeIPRegion` response
- **THEN** the test SHALL verify that the data source correctly maps the response fields to the Terraform schema

### Requirement: Documentation for tencentcloud_teo_ip_region
The system SHALL provide a markdown documentation file `data_source_tc_teo_ip_region.md` under `gendoc/teo/` directory, following the project's documentation format.

#### Scenario: Documentation file exists
- **WHEN** the data source is implemented
- **THEN** a documentation file SHALL exist with description, example usage, and schema reference
