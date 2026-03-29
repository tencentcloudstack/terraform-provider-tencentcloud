## 1. Schema 修改

- [x] 1.1 在 `tencentcloud/services/teo/data_source_tc_teo_l7_acc_rule.go` 的 schema 中添加 `total_count` 字段，类型为 `TypeInt`，设置 `Computed: true`
- [x] 1.2 确保新增字段不影响现有 schema 的向后兼容性

## 2. Read 函数实现

- [x] 2.1 在 `data_source_tc_teo_l7_acc_rule.go` 的 Read 函数中，从 DescribeL7AccRules API 响应中读取 `TotalCount` 字段
- [x] 2.2 将 `TotalCount` 值设置到 state 的 `total_count` 字段
- [x] 2.3 添加 nil 检查逻辑，当 `TotalCount` 为 nil 时设置默认值 0

## 3. 测试更新

- [x] 3.1 在 `data_source_tc_teo_l7_acc_rule_test.go` 中添加验收测试用例，验证 `total_count` 字段正确返回
- [x] 3.2 确保测试覆盖 TotalCount 正常值和 nil 值的场景
- [ ] 3.3 运行测试验证功能正确性（需要在有 Go 环境的情况下手动执行）

## 4. 文档更新

- [x] 4.1 更新 `data_source_tc_teo_l7_acc_rule.md` 示例文件，添加 `total_count` 字段的示例
- [ ] 4.2 运行 `make doc` 命令生成 `website/docs/` 下的 markdown 文档（需要在有 Go 环境的情况下手动执行）
- [ ] 4.3 确认生成的文档中包含 `total_count` 字段的说明（需要在有 Go 环境的情况下手动执行）

## 5. 代码验证

- [ ] 5.1 运行 `make fmt` 进行代码格式化（需要在有 Go 环境的情况下手动执行）
- [ ] 5.2 运行 `make lint` 进行代码检查（需要在有 Go 环境的情况下手动执行）
- [ ] 5.3 运行 `make test` 运行单元测试（需要在有 Go 环境的情况下手动执行）
- [ ] 5.4 运行 `make testacc` 运行验收测试（需要在有 Go 环境的情况下手动执行，需要配置 TENCENTCLOUD_SECRET_ID/KEY 环境变量）
