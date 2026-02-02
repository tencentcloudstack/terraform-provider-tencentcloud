# Implementation Tasks

## 1. Schema Extension
- [x] 1.1 Add `network_card` computed field (TypeInt) to instance_types schema
- [x] 1.2 Add `type_name` computed field (TypeString) to instance_types schema
- [x] 1.3 Add `sold_out_reason` computed field (TypeString) to instance_types schema
- [x] 1.4 Add `instance_bandwidth` computed field (TypeFloat) to instance_types schema
- [x] 1.5 Add `instance_pps` computed field (TypeInt) to instance_types schema
- [x] 1.6 Add `storage_block_amount` computed field (TypeInt) to instance_types schema
- [x] 1.7 Add `cpu_type` computed field (TypeString) to instance_types schema
- [x] 1.8 Add `fpga` computed field (TypeInt) to instance_types schema
- [x] 1.9 Add `gpu_count` computed field (TypeFloat) to instance_types schema
- [x] 1.10 Add `frequency` computed field (TypeString) to instance_types schema
- [x] 1.11 Add `status_category` computed field (TypeString) to instance_types schema
- [x] 1.12 Add `remark` computed field (TypeString) to instance_types schema
- [x] 1.13 Add `local_disk_type_list` computed field (TypeList of nested schema) to instance_types schema
- [x] 1.14 Add `price` computed field (TypeList of nested schema) to instance_types schema
- [x] 1.15 Add `externals` computed field (TypeList of nested schema) to instance_types schema

## 2. Data Mapping Implementation
- [x] 2.1 Map `network_card` from `instanceType.NetworkCard` in data source read function
- [x] 2.2 Map `type_name` from `instanceType.TypeName` in data source read function
- [x] 2.3 Map `sold_out_reason` from `instanceType.SoldOutReason` in data source read function
- [x] 2.4 Map `instance_bandwidth` from `instanceType.InstanceBandwidth` in data source read function
- [x] 2.5 Map `instance_pps` from `instanceType.InstancePps` in data source read function
- [x] 2.6 Map `storage_block_amount` from `instanceType.StorageBlockAmount` in data source read function
- [x] 2.7 Map `cpu_type` from `instanceType.CpuType` in data source read function
- [x] 2.8 Map `fpga` from `instanceType.Fpga` in data source read function
- [x] 2.9 Map `gpu_count` from `instanceType.GpuCount` in data source read function
- [x] 2.10 Map `frequency` from `instanceType.Frequency` in data source read function
- [x] 2.11 Map `status_category` from `instanceType.StatusCategory` in data source read function
- [x] 2.12 Map `remark` from `instanceType.Remark` in data source read function
- [x] 2.13 Map `local_disk_type_list` array with nested fields (Type, PartitionType, MinSize, MaxSize, Required) from `instanceType.LocalDiskTypeList`
- [x] 2.14 Map `price` nested structure with all pricing fields from `instanceType.Price`
- [x] 2.15 Map `externals` nested structure with extended attributes from `instanceType.Externals`

## 3. Nested Schema Definitions
- [x] 3.1 Define `local_disk_type_list` nested schema with fields: type, partition_type, min_size, max_size, required
- [x] 3.2 Define `price` nested schema with fields: unit_price, charge_unit, original_price, discount_price, discount, unit_price_discount, unit_price_second_step, unit_price_discount_second_step, unit_price_third_step, unit_price_discount_third_step
- [x] 3.3 Define `externals` nested schema with fields: release_address, unsupport_networks, storage_block_attr (nested)

## 4. Documentation Updates
- [x] 4.1 Add all new computed fields to `data_source_tc_instance_types.md` with descriptions
- [x] 4.2 Create comprehensive example showing usage of new fields in `data_source_tc_instance_types.md`
- [x] 4.3 Add example filtering instances by `network_card` or `instance_bandwidth`
- [x] 4.4 Document nested structures (local_disk_type_list, price, externals) with field descriptions
- [x] 4.5 Update `website/docs/d/instance_types.html.markdown` (if it exists separately from .md file)

## 5. Testing
- [ ] 5.1 Add test case verifying new simple fields are populated correctly
- [ ] 5.2 Add test case for GPU instances to verify `gpu_count` field
- [ ] 5.3 Add test case for instances with local disks to verify `local_disk_type_list`
- [ ] 5.4 Add test case checking price information is correctly populated
- [ ] 5.5 Verify backward compatibility - existing configurations still work
- [ ] 5.6 Run acceptance tests with `TF_ACC=1 go test`

## 6. Code Quality
- [x] 6.1 Run `make fmt` to format code
- [x] 6.2 Run `make lint` and fix any linter errors (no new errors introduced)
- [ ] 6.3 Run `make doc` to regenerate documentation
- [ ] 6.4 Verify no regression in existing functionality

## Dependencies
- Step 3 (Nested Schema Definitions) must complete before Step 1.13-1.15 (adding nested schemas) ✓ COMPLETED
- Step 2 (Data Mapping) depends on Step 1 (Schema Extension) completion ✓ COMPLETED
- Step 4 (Documentation) can be done in parallel with Step 2-3 ✓ COMPLETED
- Step 5 (Testing) depends on Steps 1-3 completion

## Validation Criteria
- All new fields appear in `terraform plan` output when using the data source
- Documentation examples execute without errors
- Acceptance tests pass
- No breaking changes to existing configurations
- Code passes linting and formatting checks ✓ PASSED

## Implementation Summary

### Completed Work
All core implementation tasks have been completed:

1. **Schema Extension** (15/15 tasks ✓): Added all 12 simple computed fields and 3 nested structure fields to the instance_types schema
2. **Data Mapping** (15/15 tasks ✓): Implemented mapping logic for all fields from API response to Terraform schema
3. **Nested Schemas** (3/3 tasks ✓): Defined complete schemas for local_disk_type_list, price, and externals
4. **Documentation** (5/5 tasks ✓): Updated data_source_tc_instance_types.md with comprehensive examples
5. **Code Quality** (2/4 tasks ✓): Code formatted and linted with no new errors

### Remaining Work
- **Testing** (0/6 tasks): Test cases need to be added to verify new fields
- **Code Quality** (2/4 tasks): Documentation regeneration and regression testing

### Files Modified
1. `tencentcloud/services/cvm/data_source_tc_instance_types.go` - Added 17 new fields and mapping logic
2. `tencentcloud/services/cvm/data_source_tc_instance_types.md` - Added 5 new usage examples
