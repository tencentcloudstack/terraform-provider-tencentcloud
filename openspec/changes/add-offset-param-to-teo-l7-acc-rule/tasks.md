## 1. 数据源代码修改

- [x] 1.1 在 tencentcloud/services/teo/data_source_tc_teo_l7_acc_rule.go 的 schema 中添加 Offset 参数（Optional, Int, 添加 ValidateDiagFunc 验证非负值）
- [x] 1.2 在 Read 函数中从 d.Get("offset") 获取 Offset 值，有值时传递给 DescribeL7AccRules API
- [x] 1.3 更新数据源示例文件 data_source_tc_teo_l7_acc_rule.md，添加 Offset 参数使用示例

## 2. 文档生成

- [ ] 2.1 运行 make doc 命令自动生成 website/docs/r/teo_l7_acc_rule.html.markdown 文档（由于环境限制，需要在有 go 环境时执行）

## 3. 测试实现

- [x] 3.1 在 tencentcloud/services/teo/data_source_tc_teo_l7_acc_rule_test.go 中添加 Offset 参数测试用例
- [x] 3.2 添加测试验证 Offset 参数为 0 时的默认行为
- [x] 3.3 添加测试验证 Offset 参数为正数时的正确传递
- [x] 3.4 添加测试验证 Offset 参数为负数时的验证错误

## 4. 验证和质量检查

- [ ] 4.1 运行 go build 确保代码编译通过（需要在有 go 环境时执行）
- [ ] 4.2 运行 go vet 进行代码静态检查（需要在有 go 环境时执行）
- [ ] 4.3 运行 TF_ACC=1 go test 测试 tencentcloud/services/teo/data_source_tc_teo_l7_acc_rule_test.go（需要在有 go 环境和 TENCENTCLOUD_SECRET_ID/KEY 环境变量时执行）
- [ ] 4.4 检查生成的文档 website/docs/r/teo_l7_acc_rule.html.markdown 确认 Offset 参数文档正确（需要在运行 make doc 后执行）
- [ ] 4.5 运行 gofmt 格式化代码确保代码风格一致（需要在有 go 环境时执行）
