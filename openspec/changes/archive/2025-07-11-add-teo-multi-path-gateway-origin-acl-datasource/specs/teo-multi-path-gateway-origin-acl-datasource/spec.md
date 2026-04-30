## ADDED Requirements

### Requirement: Data source schema for multi-path gateway origin ACL
The data source `tencentcloud_teo_multi_path_gateway_origin_acl` SHALL define the following schema:

**Input parameters:**
- `zone_id` (TypeString, Required): Zone ID
- `gateway_id` (TypeString, Required): Gateway ID
- `result_output_file` (TypeString, Optional): Used to save results

**Computed output:**
- `multi_path_gateway_origin_acl_info` (TypeList, Computed, MaxItems:1): Origin ACL info block containing:
  - `multi_path_gateway_current_origin_acl` (TypeList, Computed, MaxItems:1): Current effective origin ACL
    - `entire_addresses` (TypeList, Computed, MaxItems:1): IP CIDR details
      - `ipv4` (TypeSet, Computed): IPv4 CIDR list
      - `ipv6` (TypeSet, Computed): IPv6 CIDR list
    - `version` (TypeString, Computed): Version number
    - `is_planed` (TypeString, Computed): Whether update confirmation is completed
  - `multi_path_gateway_next_origin_acl` (TypeList, Computed, MaxItems:1): Next version origin ACL
    - `version` (TypeString, Computed): Version number
    - `entire_addresses` (TypeList, Computed, MaxItems:1): IP CIDR details
      - `ipv4` (TypeSet, Computed): IPv4 CIDR list
      - `ipv6` (TypeSet, Computed): IPv6 CIDR list
    - `added_addresses` (TypeList, Computed, MaxItems:1): Added IP CIDRs compared to current
      - `ipv4` (TypeSet, Computed): IPv4 CIDR list
      - `ipv6` (TypeSet, Computed): IPv6 CIDR list
    - `removed_addresses` (TypeList, Computed, MaxItems:1): Removed IP CIDRs compared to current
      - `ipv4` (TypeSet, Computed): IPv4 CIDR list
      - `ipv6` (TypeSet, Computed): IPv6 CIDR list
    - `no_change_addresses` (TypeList, Computed, MaxItems:1): Unchanged IP CIDRs compared to current
      - `ipv4` (TypeSet, Computed): IPv4 CIDR list
      - `ipv6` (TypeSet, Computed): IPv6 CIDR list

#### Scenario: Data source with required parameters
- **WHEN** a user declares the data source with `zone_id` and `gateway_id`
- **THEN** Terraform SHALL call `DescribeMultiPathGatewayOriginACL` API with the provided parameters and populate the computed output fields

#### Scenario: Data source with result output file
- **WHEN** a user provides `result_output_file` parameter
- **THEN** the data source SHALL write the query results to the specified file path

### Requirement: Read function calls DescribeMultiPathGatewayOriginACL API
The Read function SHALL call the `DescribeMultiPathGatewayOriginACL` API from the `teo/v20220901` SDK package with `ZoneId` and `GatewayId` as request parameters.

#### Scenario: Successful API call
- **WHEN** the Read function is invoked with valid `zone_id` and `gateway_id`
- **THEN** it SHALL build a `DescribeMultiPathGatewayOriginACLRequest` with the provided parameters and call the API with `tccommon.ReadRetryTimeout` retry logic

#### Scenario: API call failure
- **WHEN** the API call fails with a retryable error
- **THEN** the Read function SHALL return a retry error using `tccommon.RetryError()`

### Requirement: Nil safety for nested response fields
The Read function SHALL check nil at each level of the nested response before accessing child fields. This includes checking `MultiPathGatewayOriginACLInfo`, `MultiPathGatewayCurrentOriginACL`, `MultiPathGatewayNextOriginACL`, and each `Addresses` struct before reading their fields.

#### Scenario: Response with nil current origin ACL
- **WHEN** the API returns a response where `MultiPathGatewayCurrentOriginACL` is nil
- **THEN** the Read function SHALL NOT set `multi_path_gateway_current_origin_acl` and SHALL NOT panic

#### Scenario: Response with nil next origin ACL
- **WHEN** the API returns a response where `MultiPathGatewayNextOriginACL` is nil
- **THEN** the Read function SHALL NOT set `multi_path_gateway_next_origin_acl` and SHALL NOT panic

### Requirement: Data source ID uses composite key
The data source SHALL set its ID as a composite of `zone_id` and `gateway_id` joined by `tccommon.FILED_SP`.

#### Scenario: Setting data source ID
- **WHEN** the Read function successfully retrieves data
- **THEN** the data source ID SHALL be set to `zone_id + FILED_SP + gateway_id`

### Requirement: Service method for DescribeMultiPathGatewayOriginACL
A service method `DescribeTeoMultiPathGatewayOriginAclByFilter` SHALL be added to `service_tencentcloud_teo.go` that accepts a `paramMap` with `ZoneId` and `GatewayId`, constructs the SDK request, calls the API, and returns the response params.

#### Scenario: Service method with valid parameters
- **WHEN** the service method is called with `ZoneId` and `GatewayId` in the paramMap
- **THEN** it SHALL create a `DescribeMultiPathGatewayOriginACLRequest`, set the fields from paramMap, call the API, and return `*teov20220901.DescribeMultiPathGatewayOriginACLResponseParams`

### Requirement: Provider registration
The data source SHALL be registered in `provider.go` under the dataSources map with key `tencentcloud_teo_multi_path_gateway_origin_acl` and value `teo.DataSourceTencentCloudTeoMultiPathGatewayOriginAcl()`.

#### Scenario: Data source available in provider
- **WHEN** the provider is initialized
- **THEN** the data source `tencentcloud_teo_multi_path_gateway_origin_acl` SHALL be available for use in Terraform configurations

### Requirement: Documentation file
A documentation file `data_source_tc_teo_multi_path_gateway_origin_acl.md` SHALL be created following the TEO documentation pattern, including a one-line description mentioning the cloud product name (TEO), example usage with `zone_id` and `gateway_id`, and no Argument/Attribute Reference sections.

#### Scenario: Documentation file exists
- **WHEN** the data source is implemented
- **THEN** a `.md` documentation file SHALL exist alongside the `.go` file with proper description and example usage
