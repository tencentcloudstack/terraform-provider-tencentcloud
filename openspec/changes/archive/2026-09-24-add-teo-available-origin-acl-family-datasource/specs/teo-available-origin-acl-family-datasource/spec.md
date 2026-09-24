## ADDED Requirements

### Requirement: Data source for querying available origin ACL families
The provider SHALL provide a data source `tencentcloud_teo_available_origin_acl_family` that queries TEO origin protection available origin ACL control domains by calling the cloud API `DescribeAvailableOriginACLFamily`. The data source SHALL accept all input parameters of the API: `zone_id` (required) and `filters` (optional, filtering by `origin_acl_family`). Pagination parameters (`Offset`/`Limit`) SHALL NOT be exposed to users and SHALL be handled internally to fetch all available data.

#### Scenario: Query all available origin ACL families for a zone
- **WHEN** the user declares `tencentcloud_teo_available_origin_acl_family` with only `zone_id` set (no `filters`)
- **THEN** the provider SHALL call `DescribeAvailableOriginACLFamily` with `ZoneId` and internal pagination (`Limit=100`) to fetch all available control domains
- **AND** the `origin_acl_family_info_set` output SHALL contain every returned control domain with its `version`, `active_time`, `entire_addresses` (containing `ipv4`/`ipv6` lists), and `origin_acl_family`
- **AND** `total_count` SHALL reflect the total number of available control domains returned by the API

#### Scenario: Filter available origin ACL families by control domain
- **WHEN** the user declares `tencentcloud_teo_available_origin_acl_family` with `zone_id` and `filters` containing `name = "OriginACLFamily"` and one or more `values` (e.g. `["gaz", "mlc"]`)
- **THEN** the provider SHALL pass a `Filter` with `Name="OriginACLFamily"` and `Values` set to the provided values when calling `DescribeAvailableOriginACLFamily`
- **AND** the `origin_acl_family_info_set` SHALL contain only the control domains matching the filter

#### Scenario: Read returns empty response from cloud API
- **WHEN** `DescribeAvailableOriginACLFamily` returns an empty response (`response == nil` or `response.Response == nil`)
- **THEN** the provider SHALL return a `NonRetryableError` instead of clearing the data source id
- **AND** the provider SHALL log `[DATASOURCE] read empty, skip SetId` for troubleshooting

### Requirement: Output fields of the available origin ACL family data source
The data source `tencentcloud_teo_available_origin_acl_family` SHALL expose the control domain details as a flattened Computed list field `origin_acl_family_info_set`, where each element's parameters are flattened (not wrapped in an extra list layer). Each element SHALL include `version` (string), `active_time` (string), `entire_addresses` (list containing `ipv4` set and `ipv6` set), and `origin_acl_family` (string). The data source SHALL also expose `total_count` (int) as a top-level Computed field, and an optional `result_output_file` (string) to save query results.

#### Scenario: Output fields are populated from API response
- **WHEN** the API returns `OriginACLFamilyInfos` with one or more `OriginACLFamilyInfo` elements
- **THEN** the provider SHALL map each element's `Version` to `version`, `ActiveTime` to `active_time`, `OriginACLFamily` to `origin_acl_family`
- **AND** the provider SHALL map `EntireAddresses.IPv4` to `entire_addresses.ipv4` and `EntireAddresses.IPv6` to `entire_addresses.ipv6` when non-nil
- **AND** the provider SHALL only call `d.Set(...)` for a field when the corresponding API response field is non-nil

#### Scenario: Save results to output file
- **WHEN** the user sets `result_output_file` to a valid file path
- **THEN** the provider SHALL write the query result to that file
- **AND** the data source read SHALL still succeed when `result_output_file` is not set

### Requirement: Register the data source in provider
The provider SHALL register the data source `tencentcloud_teo_available_origin_acl_family` in `tencentcloud/provider.go` and `tencentcloud/provider.md` so it is available to Terraform users.

#### Scenario: Data source is available to Terraform
- **WHEN** the user runs `terraform plan` referencing `data.tencentcloud_teo_available_origin_acl_family`
- **THEN** the provider SHALL recognize the data source type without error
- **AND** the data source SHALL appear in the provider documentation