# Change: Extend instance_types Data Source with All API Response Fields

## Why

The `tencentcloud_instance_types` data source currently only exposes a limited subset of fields returned by the `DescribeZoneInstanceConfigInfos` API. Users need access to additional instance configuration details such as network card type, instance bandwidth, CPU type, frequency, local disk specifications, pricing information, and GPU count to make informed decisions when selecting instance types for their infrastructure.

The API returns comprehensive instance configuration data through the `InstanceTypeQuotaItem` structure, but the current Terraform data source only maps 8 fields out of 25+ available fields, limiting the usefulness of this data source for capacity planning and cost optimization.

## What Changes

- Add missing computed fields to the `instance_types` data source schema to match all fields returned by the `DescribeZoneInstanceConfigInfos` API
- Map all available fields from `InstanceTypeQuotaItem` SDK structure to the Terraform schema
- Update data source documentation with new fields and examples
- Ensure backward compatibility - all new fields are optional computed fields

### New Computed Fields to Add

From `InstanceTypeQuotaItem` structure:
- `network_card` - Network card type (e.g., 25 represents 25G network card)
- `type_name` - Instance type display name
- `local_disk_type_list` - List of local disk specifications
- `sold_out_reason` - Reason for sold out status (if applicable)
- `instance_bandwidth` - Internal network bandwidth in Gbps
- `instance_pps` - Network packet forwarding capacity in 10K PPS
- `storage_block_amount` - Number of local storage blocks
- `cpu_type` - Processor model
- `fpga` - Number of FPGA cores
- `gpu_count` - Physical GPU card count mapped to instance
- `frequency` - CPU frequency information
- `status_category` - Stock status category (EnoughStock/NormalStock/UnderStock/WithoutStock)
- `remark` - Instance remarks
- `price` - Instance pricing information (nested structure)
- `externals` - Extended attributes (nested structure)

## Impact

### Affected Specs
- `datasource-cvm-instance-types` (to be created as this is the first spec for this data source)

### Affected Code
- `tencentcloud/services/cvm/data_source_tc_instance_types.go` - Add new fields to schema and mapping logic
- `tencentcloud/services/cvm/data_source_tc_instance_types.md` - Update documentation with new fields
- `tencentcloud/services/cvm/data_source_tc_instance_types_test.go` - Add test coverage for new fields (if tests exist)
- `website/docs/d/instance_types.html.markdown` - Update public documentation

### Breaking Changes
None - all changes are additive (new computed fields only)

### Dependencies
None - uses existing `DescribeZoneInstanceConfigInfos` API call

### Testing Requirements
- Verify all new fields are correctly populated from API response
- Test with different instance types (standard, GPU, local disk types)
- Test with sold out instances to verify sold_out_reason field
- Verify backward compatibility with existing configurations
