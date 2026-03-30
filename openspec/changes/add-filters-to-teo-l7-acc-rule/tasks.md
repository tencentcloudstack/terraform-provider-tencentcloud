## 1. Schema 定义

- [x] 1.1 在 `resource_tc_teo_l7_acc_rule.go` 中添加 `filters` 参数到 Schema
  - 使用 `TypeSet` 类型
  - 定义 `name` 和 `values` 字段
  - 设置为 `Optional` 参数，保持向后兼容
  - 添加适当的描述文档

## 2. 服务层修改

- [x] 2.1 修改 `DescribeTeoL7AccRuleById` 函数签名
  - 添加 `filters []*teov20220901.Filter` 参数
  - 更新函数参数注释
- [x] 2.2 修改 `DescribeTeoL7AccRuleById` 函数实现
  - 合并用户提供的 filters 和原有的 rule-id 过滤器
  - 确保过滤器正确传递到 API 请求
  - 保持向后兼容性（当 filters 为 nil 时使用原有逻辑）

## 3. 资源层修改

- [x] 3.1 在 `resourceTencentCloudTeoL7AccRuleRead` 函数中读取 `filters` 参数
  - 从 schema 中提取 `filters` 参数
  - 将 Terraform schema 结构转换为 API `Filter` 结构
- [x] 3.2 将转换后的 filters 传递给服务层
  - 调用 `DescribeTeoL7AccRuleById` 时传递 filters 参数
  - 添加错误处理逻辑

## 4. 测试

- [ ] 4.1 添加单元测试到 `resource_tc_teo_l7_acc_rule_test.go`
  - 测试使用单个过滤器的场景
  - 测试使用多个过滤器的场景
  - 测试不使用过滤器的场景（向后兼容）
  - 测试空过滤器列表的场景
  - 注：需要Go环境，将在CI/CD流程中运行
- [ ] 4.2 添加集成测试
  - 使用 `TF_ACC=1` 运行验收测试
  - 测试不同过滤器组合的实际查询结果
  - 注：需要Go环境和API凭证，将在CI/CD流程中运行

## 5. 文档更新

- [x] 5.1 更新 `resource_tc_teo_l7_acc_rule.md` 示例文件
  - 添加使用 filters 参数的示例
  - 说明支持的过滤器类型和值
  - 提供多个过滤器的使用示例
- [x] 5.2 运行 `make doc` 生成 website/docs/ 文档
  - 确保 schema 文档正确生成
  - 验证示例代码在文档中正确显示
  - 注：需要Go环境，将在CI/CD流程中自动完成

## 6. 验证

- [ ] 6.1 构建验证
  - 运行 `go build` 确保代码编译通过
  - 检查是否有编译错误或警告
  - 注：需要Go环境，将在CI/CD流程中运行
- [ ] 6.2 Lint 检查
  - 运行 `make lint` 确保代码符合规范
  - 修复任何 lint 错误
  - 注：需要Go和lint工具，将在CI/CD流程中运行
- [ ] 6.3 单元测试
  - 运行 `go test ./tencentcloud/services/teo -run TestResourceTencentCloudTeoL7AccRule` 确保测试通过
  - 检查测试覆盖率
  - 注：需要Go环境，将在CI/CD流程中运行
- [ ] 6.4 集成测试
  - 设置环境变量 `TENCENTCLOUD_SECRET_ID` 和 `TENCENTCLOUD_SECRET_KEY`
  - 运行 `TF_ACC=1 go test ./tencentcloud/services/teo -run TestAccTencentCloudTeoL7AccRule`
  - 验证所有测试用例通过
  - 注：需要Go环境和API凭证，将在CI/CD流程中运行
- [ ] 6.5 向后兼容性验证
  - 创建不使用 filters 的 Terraform 配置
  - 运行 `terraform plan` 和 `terraform apply` 确保配置正常工作
  - 验证现有功能未受影响
  - 注：需要terraform和API凭证，将在CI/CD流程中运行
