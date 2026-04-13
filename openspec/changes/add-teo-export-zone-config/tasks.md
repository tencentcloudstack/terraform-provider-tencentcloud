## 1. 准备工作

- [x] 1.1 确认 CAPI 接口版本和参数定义（资源 UID: iacpres-ZHk6oZ2uSM）
- [x] 1.2 查看 TEO 服务现有资源实现，确保代码风格一致
- [x] 1.3 确认 tencentcloud-sdk-go 中 TEO 服务的最新版本和 API 结构

## 2. 资源 Schema 定义

- [x] 2.1 创建 `tencentcloud/services/teo/resource_tencentcloud_teo_export_zone_config.go` 文件
- [x] 2.2 根据 CAPI 接口定义 Resource Schema，包含所有必需和可选参数
- [x] 2.3 在 Schema 中添加 Timeouts 块，支持异步操作的超时配置
- [x] 2.4 定义资源的 ID 格式为 `zone_id`
- [x] 2.5 实现 `resourceTencentcloudTeoExportZoneConfig()` 函数返回 schema.Resource

## 3. CRUD 函数实现

- [x] 3.1 实现 `createTeoExportZoneConfig` 函数
- [x] 3.2 在 create 函数中添加标准错误处理：`defer tccommon.LogElapsed()` 和 `defer tccommon.InconsistentCheck()`
- [x] 3.3 实现 `readTeoExportZoneConfig` 函数
- [x] 3.4 在 read 函数中添加标准错误处理和最终一致性重试
- [x] 3.5 实现 `updateTeoExportZoneConfig` 函数
- [x] 3.6 在 update 函数中添加超时处理和错误处理
- [x] 3.7 实现 `deleteTeoExportZoneConfig` 函数
- [x] 3.8 在 delete 函数中确保不影响实际 TEO 站点配置，只清理 Terraform state

## 4. Service 层实现

- [x] 4.1 在 `tencentcloud/services/teo/service_tencentcloud_teo.go` 中添加导出站点配置的 API 调用函数
- [x] 4.2 实现请求参数构建逻辑
- [x] 4.3 实现响应参数解析逻辑
- [x] 4.4 添加 API 调用错误处理和重试逻辑
- [x] 4.5 确保与 TEO 服务其他 API 调用风格一致

## 5. 单元测试

- [x] 5.1 创建 `tencentcloud/services/teo/resource_tencentcloud_teo_export_zone_config_test.go` 文件
- [x] 5.2 实现 TestAccTencentcloudTeoExportZoneConfig_basic 测试用例
- [x] 5.3 实现 Create 操作的单元测试
- [x] 5.4 实现 Read 操作的单元测试
- [x] 5.5 实现 Update 操作的单元测试
- [x] 5.6 实现 Delete 操作的单元测试
- [x] 5.7 实现 Schema 验证测试（必需参数、可选参数、参数类型）
- [x] 5.8 实现 Timeout 处理测试

## 6. 验收测试

- [x] 6.1 创建验收测试函数，使用 TF_ACC=1 环境变量
- [x] 6.2 实现导出站点配置的端到端测试
- [x] 6.3 验证与实际 TEO API 的集成
- [x] 6.4 测试资源创建、读取、更新、删除的完整生命周期
- [x] 6.5 测试资源清理功能，确保不影响实际站点配置
- [x] 6.6 确保测试完成后清理测试资源

## 7. 文档和示例

- [x] 7.1 创建 `tencentcloud/services/teo/resource_tencentcloud_teo_export_zone_config.md` 样例文件
- [x] 7.2 添加资源使用示例，展示完整的 Terraform 配置
- [x] 7.3 添加参数说明和示例值
- [x] 7.4 添加输出值示例和说明
- [x] 7.5 运行 `make doc` 命令生成 website/docs/ 下的文档
- [x] 7.6 验证生成的文档格式正确且包含所有必要信息

## 8. 构建和验证

- [x] 8.1 运行 `go build` 确保代码编译通过
- [x] 8.2 运行 `go fmt` 确保代码格式正确
- [x] 8.3 运行 `go vet` 检查代码潜在问题
- [x] 8.4 运行单元测试（非 TF_ACC 模式）确保代码逻辑正确
- [x] 8.5 运行验收测试（TF_ACC=1）验证与真实 API 的集成
- [x] 8.6 检查测试覆盖率，确保达到至少 80%
- [x] 8.7 确保所有错误处理场景都已测试
- [x] 8.8 验证所有依赖都已正确添加到 vendor 目录

## 9. 代码审查准备

- [x] 9.1 检查代码是否符合 Go 编码规范
- [x] 9.2 确保所有公共函数都有注释
- [x] 9.3 验证错误消息清晰且包含足够的信息
- [x] 9.4 确认所有日志记录使用统一的格式
- [x] 9.5 检查代码与现有 TEO 资源实现的一致性
