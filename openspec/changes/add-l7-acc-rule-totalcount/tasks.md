## 1. 代码修改

- [x] 1.1 在 ResourceTencentCloudTeoL7AccRule() 函数的 Schema 中添加 total_count 字段定义
- [x] 1.2 在 resourceTencentCloudTeoL7AccRuleRead() 函数中设置 total_count 字段值

## 2. 测试

- [x] 2.1 添加单元测试，验证 total_count 字段正确读取和设置
- [ ] 2.2 运行 acceptance test，验证集成功能正常
- [x] 2.3 更新 resource_tc_teo_l7_acc_rule_test.go 测试文件，添加 total_count 相关测试用例

## 3. 文档

- [x] 3.1 更新 resource_tc_teo_l7_acc_rule.md 文档，添加 total_count 字段的说明
- [x] 3.2 验证文档完整性，确保示例代码包含 total_count 字段

## 4. 验证

- [ ] 4.1 确认所有代码编译通过
- [ ] 4.2 确认所有测试通过
- [x] 4.3 确认向后兼容性，现有配置无需修改
- [ ] 4.4 执行 make lint 确保代码风格符合规范
- [ ] 4.5 执行 make fmt 确保代码格式正确
