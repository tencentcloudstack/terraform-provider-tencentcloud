## 1. 数据源实现

- [x] 1.1 创建 `tencentcloud/services/cvm/data_source_tc_cvm_account_quota.go` 文件
- [x] 1.2 定义数据源 Schema，包含 `zone` 和 `quota_type` 过滤参数
- [x] 1.3 定义输出属性 Schema：`app_id`, `account_quota_overview`, `result_output_file`
- [x] 1.4 实现 `dataSourceTencentCloudCvmAccountQuota()` 函数返回 schema.Resource
- [x] 1.5 实现 `dataSourceTencentCloudCvmAccountQuotaRead()` 函数
- [x] 1.6 实现 Filter 参数构建逻辑，将 `zone` 和 `quota_type` 转换为 API Filter 格式
- [x] 1.7 调用 DescribeAccountQuota API 获取配额数据
- [x] 1.8 实现 `flattenAccountQuotaOverview()` 函数转换 AccountQuotaOverview 数据
- [x] 1.9 实现 `flattenAccountQuota()` 函数转换 AccountQuota 数据
- [x] 1.10 实现 `flattenPostPaidQuotaSet()` 函数转换后付费配额列表
- [x] 1.11 实现 `flattenPrePaidQuotaSet()` 函数转换预付费配额列表
- [x] 1.12 实现 `flattenSpotPaidQuotaSet()` 函数转换竞价实例配额列表
- [x] 1.13 实现 `flattenImageQuotaSet()` 函数转换镜像配额列表
- [x] 1.14 实现 `flattenDisasterRecoverGroupQuotaSet()` 函数转换置放群组配额列表
- [x] 1.15 设置输出属性到 Terraform state
- [x] 1.16 实现 `result_output_file` 功能，将结果保存为 JSON 文件
- [x] 1.17 添加错误处理：defer tccommon.LogElapsed() 和 defer tccommon.InconsistentCheck()

## 2. Provider 注册

- [x] 2.1 在 `tencentcloud/provider.go` 的 DataSourcesMap 中注册 `tencentcloud_cvm_account_quota`
- [x] 2.2 验证数据源注册名称与文档一致

## 3. 测试文件创建

- [x] 3.1 创建 `tencentcloud/services/cvm/data_source_tc_cvm_account_quota_test.go` 文件
- [x] 3.2 实现 `TestAccTencentCloudCvmAccountQuotaDataSource_basic` 基本查询测试
- [x] 3.3 实现 `TestAccTencentCloudCvmAccountQuotaDataSource_filterByZone` 按可用区过滤测试
- [x] 3.4 实现 `TestAccTencentCloudCvmAccountQuotaDataSource_filterByQuotaType` 按配额类型过滤测试
- [x] 3.5 添加测试数据验证：检查返回的 app_id 和 account_quota_overview 结构

## 4. 文档模板创建

- [x] 4.1 创建 `tencentcloud/services/cvm/data_source_tc_cvm_account_quota.md` 文档模板
- [x] 4.2 添加数据源描述：查询 CVM 账户配额详情
- [x] 4.3 添加基本用法示例：不带过滤条件的查询
- [x] 4.4 添加按可用区过滤示例
- [x] 4.5 添加按配额类型过滤示例
- [x] 4.6 添加综合过滤示例（zone + quota_type）
- [x] 4.7 添加使用 result_output_file 保存结果的示例

## 5. Provider 文档更新

- [x] 5.1 在 `tencentcloud/provider.md` 的数据源列表中添加 `tencentcloud_cvm_account_quota`
- [x] 5.2 确保按字母序插入到正确位置（CVM 相关数据源区域）

## 6. 代码质量检查

- [x] 6.1 运行 `go build` 确保代码编译通过
- [x] 6.2 运行 `go fmt` 格式化代码
- [x] 6.3 运行 `go vet` 检查代码问题
- [x] 6.4 检查是否符合 Provider 代码规范（错误处理、命名规范等）

## 7. 单元测试

- [x] 7.1 运行 `go test ./tencentcloud/services/cvm/data_source_tc_cvm_account_quota_test.go` 验证测试通过
- [x] 7.2 检查测试覆盖率，确保主要逻辑路径被覆盖

## 8. 验收测试

- [ ] 8.1 设置环境变量 `TF_ACC=1`, `TENCENTCLOUD_SECRET_ID`, `TENCENTCLOUD_SECRET_KEY`
- [ ] 8.2 运行 `go test -v ./tencentcloud/services/cvm/data_source_tc_cvm_account_quota_test.go`
- [ ] 8.3 验证基本查询测试通过
- [ ] 8.4 验证按可用区过滤测试通过
- [ ] 8.5 验证按配额类型过滤测试通过
- [ ] 8.6 手动创建 Terraform 配置测试实际使用场景

## 9. 文档生成

- [x] 9.1 运行 `make doc` 命令生成 `website/docs/d/cvm_account_quota.html.markdown` 文档
- [x] 9.2 验证生成的文档格式正确，包含所有参数说明和示例

## 10. 最终验证

- [x] 10.1 运行完整的编译：`go build`
- [x] 10.2 检查是否有 linter 错误或警告
- [x] 10.3 验证所有测试通过
- [x] 10.4 检查文档完整性和准确性

## 11. 提交准备

- [ ] 11.1 在 `.changelog/` 目录创建 changelog 文件，描述新增的数据源
- [x] 11.2 检查所有文件符合项目编码规范
- [x] 11.3 确认没有破坏性变更
- [ ] 11.4 准备 PR 描述：说明新增的数据源功能、API 映射、测试覆盖情况
