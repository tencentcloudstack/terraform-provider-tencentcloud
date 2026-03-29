## 1. Schema 定义更新

- [x] 1.1 在 `tencentcloud/services/teo/data_source_tc_teo_l7_acc_rule.go` 的 Schema 函数中添加 TotalCount 字段定义，类型为 TypeInt，设置为 Computed

## 2. 数据映射实现

- [x] 2.1 在 `tencentcloud/services/teo/data_source_tc_teo_l7_acc_rule.go` 的 Read 函数中，从 DescribeL7AccRules API 响应中读取 TotalCount 字段并使用 d.Set() 设置到 state
- [x] 2.2 添加 nil 检查确保安全，只有当 TotalCount 字段非 nil 时才设置

## 3. 测试更新

- [x] 3.1 在 `tencentcloud/services/teo/data_source_tc_teo_l7_acc_rule_test.go` 中添加测试用例，验证 TotalCount 字段正确设置
- [x] 3.2 运行测试确保现有测试用例仍然通过（TF_ACC=1）

## 4. 文档更新

- [x] 4.1 运行 `make doc` 自动生成 website/docs/ 下的 markdown 文档
- [x] 4.2 验证生成的文档中包含 TotalCount 字段的正确描述

## 5. 验证与构建

- [x] 5.1 运行 `go build` 确保代码编译通过
- [x] 5.2 运行 `make lint` 确保代码符合规范
- [x] 5.3 运行完整的测试套件确保功能正确
