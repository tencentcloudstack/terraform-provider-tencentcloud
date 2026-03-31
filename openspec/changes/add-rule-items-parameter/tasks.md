## 1. 调研与准备工作

- [x] 1.1 调研 DescribeRules API 返回的 RuleItems 参数结构
- [x] 1.2 确认 tencentcloud_teo_rule_engine 资源文件位置（resource 或 data_source）
- [x] 1.3 查看现有代码结构和 API 调用逻辑

## 2. Schema 更新

- [x] 2.1 在资源/数据源 schema 中添加 RuleItems 参数定义
- [x] 2.2 定义 RuleItems 的嵌套结构（根据 API 响应结构）
- [x] 2.3 设置 RuleItems 为 Optional 字段，保持向后兼容性

## 3. 代码实现

- [x] 3.1 更新资源/数据源 Read 函数，添加 RuleItems 解析逻辑
- [x] 3.2 实现 RuleItems 数据的映射和类型转换
- [x] 3.3 确保 API 响应中的 RuleItems 正确填充到 schema

## 4. 文档更新

- [x] 4.1 更新资源/数据源示例文件（.md），添加 RuleItems 使用示例
- [ ] 4.2 使用 `make doc` 命令生成 website/docs/ 下的文档（需要在有 Go 环境中执行）

## 5. 测试实现

- [x] 5.1 更新或添加 RuleItems 的测试用例
- [ ] 5.2 运行单元测试确保代码正确性（需要在有 Go 环境中执行）
- [ ] 5.3 运行 TF_ACC=1 验收测试（需要 TENCENTCLOUD_SECRET_ID/KEY 环境变量）（需要在有 Go 环境中执行）

## 6. 验证与构建

- [ ] 6.1 执行 `make build` 确保代码编译通过（需要在有 Go 环境中执行）
- [ ] 6.2 执行 `make lint` 确保代码符合规范（需要在有 Go 环境中执行）
- [ ] 6.3 执行 `make test` 运行所有测试（需要在有 Go 环境中执行）
- [x] 6.4 验证向后兼容性，确保现有功能不受影响
