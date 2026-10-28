## 1. Schema Definition

- [x] 1.1 Add `add_node_list` TypeList schema field with sub-fields `role` (Required, TypeString) and `zone` (Required, TypeString) to `ResourceTencentCloudMongodbShardingInstance`
- [x] 1.2 Add `remove_node_list` TypeList schema field with sub-fields `role` (Required, TypeString), `node_name` (Required, TypeString), and `zone` (Required, TypeString) to `ResourceTencentCloudMongodbShardingInstance`

## 2. Update Function Wiring

- [x] 2.1 Add change detection for `add_node_list` and `remove_node_list` in `resourceMongodbShardingInstanceUpdate`
- [x] 2.2 Pass `add_node_list` and `remove_node_list` to the existing `UpgradeInstance` service call via the `params` map

## 3. Documentation

- [x] 3.1 Update `resource_tc_mongodb_sharding_instance.md` with example usage of `add_node_list` and `remove_node_list`

## 4. Validation

- [x] 4.1 Verify the code compiles correctly
- [x] 4.2 Review all generated code for correctness and consistency with existing patterns
