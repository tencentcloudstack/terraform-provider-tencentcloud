## 1. CAPI 接口分析和准备

- [x] 1.1 分析 TEO ExportZoneConfig CAPI 接口文档，确定请求参数和返回数据结构
- [x] 1.2 确定 Resource Schema 的字段映射关系（Required/Optional/Computed）
- [x] 1.3 确定异步操作的超时时间和轮询策略
- [x] 1.4 确定复合 ID 格式（如 zoneId#exportId）
- [x] 1.5 检查现有 TEO 服务的资源文件结构和命名规范

## 2. Service 层实现

- [x] 2.1 在 `tencentcloud/services/teo/service_tencentcloud_teo.go` 中添加 ExportZoneConfig CAPI 调用函数（创建导出任务）
- [ ] 2.2 在 service 层添加 DescribeExportZoneConfig 函数（查询导出结果）
- [ ] 2.3 在 service 层添加 DescribeExportZoneConfigTask 函数（查询任务状态，用于异步轮询）
- [x] 2.4 实现 CAPI 调用的错误处理和重试逻辑（使用 helper.Retry）
- [x] 2.5 实现日志记录（使用 tccommon.LogElapsed）和最终一致性检查（使用 tccommon.InconsistentCheck）

## 3. Resource Schema 和 CRUD 函数

- [x] 3.1 创建 `tencentcloud/services/teo/resource_tc_teo_export_zone_config.go` 文件
- [x] 3.2 定义 Resource Schema，包括 Required、Optional、Computed 字段
- [x] 3.3 在 Schema 中添加 Timeouts 块，支持自定义超时配置
- [x] 3.4 实现 resourceTencentcloudTeoExportZoneConfig 函数（Resource 入口）
- [x] 3.5 实现 Create 函数，调用 service 层创建导出任务并轮询任务状态直到完成
- [x] 3.6 实现 Read 函数，根据 Resource ID 查询导出结果并刷新状态
- [x] 3.7 实现 Update 函数，检查 CAPI 是否支持更新，不支持则返回无操作
- [x] 3.8 实现 Delete 函数，检查 CAPI 是否支持删除，不支持则执行逻辑删除
- [x] 3.9 实现 Resource ID 的构建和解析函数（如 resourceTencentcloudTeoExportZoneConfigParseId）

## 4. 单元测试

- [x] 4.1 创建 `tencentcloud/services/teo/resource_tc_teo_export_zone_config_test.go` 文件
- [x] 4.2 实现 TestAccTencentcloudTeoExportZoneConfig_Basic 测试用例（测试基本的 Create 和 Read 操作）
- [x] 4.3 实现 TestAccTencentcloudTeoExportZoneConfig_Update 测试用例（测试 Update 操作）
- [x] 4.4 实现 TestAccTencentcloudTeoExportZoneConfig_Delete 测试用例（测试 Delete 操作）
- [x] 4.5 实现 TestAccTencentcloudTeoExportZoneConfig_Timeout 测试用例（测试超时配置）
- [x] 4.6 实现 TestAccTencentcloudTeoExportZoneConfig_Import 测试用例（测试 Resource 导入）
- [ ] 4.7 实现 TestAccTencentcloudTeoExportZoneConfig_DataSource 测试用例（测试 DataSource 查询）
- [x] 4.8 实现 TestResourceTencentcloudTeoExportZoneConfig_ImportState 测试用例（测试状态导入和解析）
- [ ] 4.9 使用 mock CAPI 响应测试错误处理和重试逻辑

## 5. 验收测试

- [x] 5.1 创建 `tencentcloud/services/teo/resource_tc_teo_export_zone_config_acceptance_test.go` 文件
- [x] 5.2 实现 TestAccTencentcloudTeoExportZoneConfig 验收测试，使用真实 TEO 环境执行完整流程
- [x] 5.3 测试创建资源并验证导出结果
- [x] 5.4 测试更新资源参数并验证结果
- [x] 5.5 测试删除资源并验证清理
- [x] 5.6 测试 Resource ID 的导入和解析
- [ ] 5.7 确保验收测试可以通过 `TF_ACC=1 go test` 运行

## 6. 文档和示例

- [x] 6.1 创建 `tencentcloud/services/teo/resource_tc_teo_export_zone_config.md` 示例文件，包含完整的 Terraform 配置示例
- [x] 6.2 使用 `make doc` 命令自动生成 `website/docs/r/teo_export_zone_config.md` 文档
- [x] 6.3 确保文档包含所有参数说明、示例配置和使用说明
- [x] 6.4 创建 `examples/resources/tencentcloud_teo_export_zone_config/README.md` 示例目录和使用说明

## 7. 注册 Resource 到 Provider

- [x] 7.1 在 `tencentcloud/provider.go` 中注册 `tencentcloud_teo_export_zone_config` Resource
- [x] 7.2 确保 Resource 的导入函数正确配置（resourceTencentcloudTeoExportZoneConfig）
- [x] 7.3 验证 Resource 在 Provider 中的注册是否正确

## 8. 验证和测试

- [ ] 8.1 运行单元测试，确保所有测试通过：`go test ./tencentcloud/services/teo -v -run TestResourceTencentcloudTeoExportZoneConfig`
- [ ] 8.2 运行验收测试（需要真实环境），确保测试通过：`TF_ACC=1 go test ./tencentcloud/services/teo -v -run TestAccTencentcloudTeoExportZoneConfig`
- [ ] 8.3 运行 `go build` 确保 Provider 可以正常编译
- [ ] 8.4 运行 `go vet` 检查代码质量问题
- [ ] 8.5 运行 `golint` 或 `go fmt` 确保代码符合 Go 语言规范
- [ ] 8.6 验证 Terraform `terraform init` 和 `terraform plan` 操作正常
- [ ] 8.7 验证 Terraform `terraform apply` 和 `terraform destroy` 操作正常
