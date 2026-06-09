## ADDED Requirements

### Requirement: Data source for querying TEO multi-path gateway regions
The system SHALL provide a Terraform data source `tencentcloud_teo_multi_path_gateway_region` that allows users to query available regions for TEO multi-path gateway by site ID. The data source SHALL call cloud API `DescribeMultiPathGatewayRegions` with `ZoneId` as input and return `GatewayRegions` list.

#### Scenario: Query available gateway regions with valid zone_id
- **WHEN** user provides a valid `zone_id` and reads the data source
- **THEN** the system SHALL call `DescribeMultiPathGatewayRegions` with the provided `zone_id` and return a `gateway_regions` list containing items with `region_id`, `cn_name`, and `en_name` fields

#### Scenario: Query with invalid zone_id
- **WHEN** user provides an invalid `zone_id` that does not exist
- **THEN** the system SHALL return the error from the cloud API, wrapped with retry error handling

#### Scenario: Query returns empty gateway regions
- **WHEN** the cloud API returns an empty `GatewayRegions` list
- **THEN** the system SHALL set `gateway_regions` to an empty list and set the data source ID based on an empty list hash

### Requirement: Schema definition for tencentcloud_teo_multi_path_gateway_region
The data source schema SHALL define the following fields:
- `zone_id` (Required, TypeString): The site ID to query gateway regions for
- `gateway_regions` (Computed, TypeList): List of available gateway regions, each containing:
  - `region_id` (Computed, TypeString): Region ID
  - `cn_name` (Computed, TypeString): Chinese name of the region
  - `en_name` (Computed, TypeString): English name of the region
- `result_output_file` (Optional, TypeString): File path to save query results

#### Scenario: Schema input field zone_id is required
- **WHEN** user creates a data source configuration without `zone_id`
- **THEN** terraform SHALL report a missing required argument error

#### Scenario: Schema output field gateway_regions contains nested attributes
- **WHEN** the data source is read successfully
- **THEN** each item in `gateway_regions` SHALL have `region_id`, `cn_name`, and `en_name` as computed string attributes

### Requirement: Service layer method for DescribeMultiPathGatewayRegions
The system SHALL provide a service method `DescribeTeoMultiPathGatewayRegionByFilter` in `service_tencentcloud_teo.go` that wraps the `DescribeMultiPathGatewayRegions` API call with retry handling using `tccommon.ReadRetryTimeout`.

#### Scenario: Service method calls API with zone_id
- **WHEN** the service method is called with a `paramMap` containing `ZoneId`
- **THEN** it SHALL create a `DescribeMultiPathGatewayRegionsRequest`, set the `ZoneId` field, call the API with retry, and return the `GatewayRegions` response list

#### Scenario: Service method handles API errors with retry
- **WHEN** the API call fails
- **THEN** the method SHALL retry with `tccommon.ReadRetryTimeout` and return a `tccommon.RetryError()` wrapped error if all retries fail

### Requirement: Provider registration for the new data source
The data source `tencentcloud_teo_multi_path_gateway_region` SHALL be registered in `provider.go` and `provider.md`.

#### Scenario: Data source is registered in provider
- **WHEN** the provider is initialized
- **THEN** the data source `tencentcloud_teo_multi_path_gateway_region` SHALL be available for use in Terraform configurations

### Requirement: Unit tests with gomonkey mock
The data source SHALL have unit tests in `data_source_tc_teo_multi_path_gateway_region_test.go` using gomonkey to mock cloud API calls, testing business logic without requiring real API access.

#### Scenario: Unit test for successful read
- **WHEN** the read function is called with valid parameters and the API is mocked to return gateway regions
- **THEN** the test SHALL verify that `gateway_regions` is correctly populated with the mocked data

#### Scenario: Unit test for API error
- **WHEN** the read function is called and the mocked API returns an error
- **THEN** the test SHALL verify that the error is properly propagated
