# Tasks: Expose All DescribeApmInstances Response Fields

## Phase 1: Schema 扩展
- [x] 1.1 在 `data_source_tc_apm_instances.go` 的 `instance_list` Elem Schema 中新增所有缺失字段定义（约 30 个 Computed 字段）
- [x] 1.2 在 `dataSourceTencentCloudApmInstancesRead` 函数中补充所有缺失字段的数据映射逻辑

## Phase 2: 测试更新
- [x] 2.1 更新 `data_source_tc_apm_instances_test.go`，验证新增字段存在且类型正确

## Phase 3: 文档与代码质量
- [x] 3.1 更新 `data_source_tc_apm_instances.md`，补充新增字段说明和示例
- [x] 3.2 运行 `gofmt` 格式化代码，编译通过
