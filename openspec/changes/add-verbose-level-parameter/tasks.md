## 1. 代码实现

- [x] 1.1 在 `tencentcloud/services/vpc/data_source_tc_nats.go` 的 schema 中添加 `VerboseLevel` 可选参数（类型：Int）
- [x] 1.2 修改 `data_source_tc_nats.go` 的 Read 函数，从 schema 中读取 `VerboseLevel` 参数值
- [x] 1.3 修改 API 调用逻辑，将 `VerboseLevel` 参数传递给 `DescribeNatGateways` API
- [x] 1.4 更新 `tencentcloud/services/vpc/data_source_tc_nats_test.go`，添加 `VerboseLevel` 参数的测试用例

## 2. 文档更新

- [x] 2.1 更新 `tencentcloud/services/vpc/data_source_tc_nats.md` 示例文件，添加 `VerboseLevel` 参数的使用示例
- [ ] 2.2 运行 `make doc` 命令，自动生成 `website/docs/data-sources/tencentcloud_nats.html.markdown` 文档（跳过：环境中缺少必要的构建工具）

## 3. 代码验证

- [ ] 3.1 运行 `make build` 验证代码编译成功（跳过：环境中缺少必要的构建工具）
- [ ] 3.2 运行 `make lint` 检查代码规范（跳过：环境中缺少必要的构建工具）
- [ ] 3.3 运行单元测试（仅运行 `tencentcloud/services/vpc/data_source_tc_nats_test.go` 相关测试）（跳过：环境中缺少必要的构建工具）
- [ ] 3.4 设置 `TF_ACC=1` 环境变量，运行验收测试验证 `tencentcloud_nats` 数据源功能（跳过：环境中缺少必要的构建工具）
