## 1. Schema 定义

- [x] 1.1 创建 `tencentcloud/services/teo/resource_tc_teo_check_cname_status_operation.go` 文件
- [x] 1.2 定义 Resource Schema，包含必需参数 `zone_id` (string) 和 `record_names` (list of string)
- [x] 1.3 定义输出参数 `cname_status`，为嵌套对象列表，包含 `record_name`、`cname`、`status` 字段
- [x] 1.4 在 schema 中设置资源的 `Timeouts` 配置（虽然是一次性操作，但仍需声明以支持 API 调用超时）
- [x] 1.5 注册资源到 `tencentcloud/services/teo/service_tencentcloud_teo.go`（如需要）

## 2. Create 函数实现

- [x] 2.1 实现 `createCheckCnameStatus` 函数，调用 TEO CheckCnameStatus API
- [x] 2.2 在 service 层实现 `teoService.checkCnameStatus` 函数，封装 API 调用逻辑
- [x] 2.3 将云 API 响应映射到 Terraform schema 输出参数
- [x] 2.4 添加错误处理和日志记录（使用 `defer tccommon.LogElapsed()`）
- [x] 2.5 添加参数验证（确保 `zone_id` 和 `record_names` 不为空）

## 3. Service 层实现

- [x] 3.1 在 `tencentcloud/services/teo/service_tencentcloud_teo.go` 中添加 `checkCnameStatus` 方法
- [x] 3.2 实现云 API 请求参数构建逻辑（将 Terraform 参数映射到云 API 参数）
- [x] 3.3 实现云 API 响应解析逻辑（将云 API 响应映射到 Terraform 输出）
- [x] 3.4 添加 API 调用错误处理

## 4. 单元测试

- [x] 4.1 创建 `tencentcloud/services/teo/resource_tc_teo_check_cname_status_operation_test.go` 文件
- [x] 4.2 使用 gomonkey mock 云 API 的 `CheckCnameStatus` 方法
- [x] 4.3 添加测试用例：检查多个域名的 CNAME 状态
- [x] 4.4 添加测试用例：检查单个域名的 CNAME 状态
- [x] 4.5 添加测试用例：处理空 `record_names` 列表
- [x] 4.6 添加测试用例：处理 API 错误
- [x] 4.7 添加测试用例：验证必需参数缺失时的错误处理

## 5. 文档

- [x] 5.1 创建资源示例文件 `tencentcloud/services/teo/resource_tc_teo_check_cname_status_operation.md`
- [x] 5.2 在示例文件中添加基本使用示例
- [x] 5.3 在示例文件中添加参数说明
- [x] 5.4 在示例文件中添加输出字段说明

## 6. 验证

- [x] 6.1 运行单元测试 `go test -v ./tencentcloud/services/teo/... -run TestAccTeoCheckCnameStatus`
- [x] 6.2 使用 `make doc` 命令生成 website 文档，验证文档生成正确

## 7. 资源注册（补充）

- [x] 7.1 在 `tencentcloud/provider.go` 中注册 `tencentcloud_teo_check_cname_status_operation` 资源（已存在）
- [x] 7.2 在 `tencentcloud/provider.md` 中添加 `tencentcloud_teo_check_cname_status_operation` 索引条目
- [x] 7.3 修复 `.md` 源文件中不必要的 Import 部分
- [x] 7.4 重新执行 `make doc` 生成更新后的文档
