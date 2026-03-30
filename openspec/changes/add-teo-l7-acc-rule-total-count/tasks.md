## 1. 修改数据源 Schema

- [x] 1.1 在 tencentcloud/services/teo/data_source_teo_l7_acc_rule.go 的 Schema 定义中添加 total_count 字段
- [x] 1.2 设置 total_count 字段为 Computed 类型，使用 schema.TypeInt

## 2. 修改数据源 Read 函数

- [x] 2.1 在 dataSourceTencentcloudTeoL7AccRuleRead 函数中，从 DescribeL7AccRules API 响应中解析 TotalCount 字段
- [x] 2.2 使用 d.Set("total_count", *response.TotalCount) 将值设置到 schema 中
- [x] 2.3 添加 nil 检查，处理 TotalCount 指针可能为 nil 的情况

## 3. 更新单元测试

- [x] 3.1 在 tencentcloud/services/teo/data_source_teo_l7_acc_rule_test.go 中添加 TotalCount 字段的验证
- [x] 3.2 确保测试用例覆盖 TotalCount 为 0 和非 0 的情况

## 4. 更新文档

- [x] 4.1 在 website/docs/t/teo_l7_acc_rule.html.markdown 中添加 total_count 字段的说明
- [x] 4.2 确保 total_count 字段标记为只读（Computed）

## 5. 代码验证

- [x] 5.1 执行 make build 确保代码编译通过
- [x] 5.2 执行 make lint 确保代码符合 lint 规则
- [x] 5.3 执行 TF_ACC=1 go test -v -run TestAccTencentcloudTeoL7AccRuleDataSource 确保测试通过
