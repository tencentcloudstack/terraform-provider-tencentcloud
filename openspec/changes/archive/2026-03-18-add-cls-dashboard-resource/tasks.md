# Implementation Tasks

## 1. Schema 定义与资源创建
- [x] 1.1 创建 `resource_tc_cls_dashboard.go` 文件
- [x] 1.2 定义资源函数 `ResourceTencentCloudClsDashboard()`
- [x] 1.3 定义输入参数 Schema
  - [ ] `dashboard_name` (Required, String) - 仪表盘名称
  - [ ] `data` (Optional, String) - 仪表盘配置数据（JSON 字符串）
  - [ ] `tags` (Optional, Map) - 标签键值对
- [x] 1.4 定义输出参数 Schema
  - [ ] `dashboard_id` (Computed, String) - 仪表盘 ID
  - [ ] `create_time` (Computed, String) - 创建时间
  - [ ] `update_time` (Computed, String) - 更新时间
- [x] 1.5 添加 Schema 描述
  - [ ] 每个字段添加清晰的 Description
  - [ ] 标注必填/可选
- [x] 1.6 配置 Importer 支持
  - [ ] 使用 `schema.ImportStatePassthrough`
  - [ ] 资源 ID 格式：`{dashboard_id}`

## 2. Create 函数实现
- [x] 2.1 实现 `resourceTencentCloudClsDashboardCreate()` 函数
- [x] 2.2 添加日志和一致性检查
  - [ ] `defer tccommon.LogElapsed()`
  - [ ] `defer tccommon.InconsistentCheck()`
- [x] 2.3 获取 LogId 和 Context
- [x] 2.4 创建 API 请求对象
  - [ ] 初始化 `CreateDashboardRequest`
  - [ ] 从 schema 读取 `dashboard_name`
  - [ ] 从 schema 读取 `data`（如果提供）
  - [ ] 处理 `tags` 转换为 Tag 数组
- [x] 2.5 调用 API 创建仪表盘
  - [ ] 使用 `resource.Retry` 包装调用
  - [ ] 使用 `tccommon.WriteRetryTimeout` (5分钟)
  - [ ] 添加错误处理和日志
- [x] 2.6 处理 API 响应
  - [ ] 从响应中提取 `dashboard_id`
  - [ ] 使用 `d.SetId()` 设置资源 ID
- [x] 2.7 处理标签（如果使用 tag service）
  - [ ] 检查是否需要调用 tag service
  - [ ] 绑定标签到资源
- [x] 2.8 调用 Read 函数刷新状态
- [x] 2.9 错误处理
  - [ ] 捕获 `InvalidParameter.DashboardNameConflict` 错误
  - [ ] 捕获 `LimitExceeded` 错误
  - [ ] 捕获 `LimitExceeded.Tag` 错误
  - [ ] 返回友好的错误信息

## 3. Service 层方法（带分页和重试）
- [x] 3.1 打开 `service_tencentcloud_cls.go`
- [x] 3.2 实现 `DescribeClsDashboardById()` 方法
  - [ ] 接收参数：`ctx context.Context`, `dashboardId string`
  - [ ] 返回：`*cls.DashboardInfo`, `error`
- [x] 3.3 实现分页查询逻辑
  - [ ] 定义 offset 和 limit 变量（offset=0, limit=20）
  - [ ] 使用 for 循环实现分页
  - [ ] 设置 `request.Offset` 和 `request.Limit`
- [x] 3.4 实现重试机制
  - [ ] 使用 `resource.Retry(tccommon.ReadRetryTimeout, ...)`
  - [ ] 在 retry 函数内调用 `DescribeDashboards`
  - [ ] 使用 `tccommon.RetryError(e, tccommon.InternalError)` 处理错误
- [x] 3.5 实现查询逻辑
  - [ ] 使用 Filters 参数按 dashboardId 过滤
  - [ ] 遍历返回的 DashboardInfos 查找匹配项
  - [ ] 找到后立即返回
- [x] 3.6 实现分页判断
  - [ ] 如果返回结果数量 < limit，停止分页
  - [ ] 否则增加 offset 继续查询
- [x] 3.7 添加错误处理
  - [ ] 处理 API 调用错误
  - [ ] 添加日志记录
  - [ ] 资源不存在时返回 nil（不报错）

## 4. Read 函数实现
- [x] 4.1 实现 `resourceTencentCloudClsDashboardRead()` 函数
- [x] 4.2 添加日志和一致性检查
- [x] 4.3 获取 LogId 和 Context
- [x] 4.4 获取资源 ID
  - [ ] 从 `d.Id()` 获取 dashboard_id
- [x] 4.5 调用 Service 层方法查询仪表盘
  - [ ] 调用 `service.DescribeClsDashboardById(ctx, dashboardId)`
  - [ ] 处理查询错误
- [x] 4.6 处理查询结果
  - [ ] 如果仪表盘不存在，清空 ID 并记录警告日志
  - [ ] 如果仪表盘存在，设置 schema 字段
- [x] 4.7 设置 Terraform State
  - [ ] `d.Set("dashboard_name", dashboard.DashboardName)`
  - [ ] `d.Set("data", dashboard.Data)`
  - [ ] `d.Set("dashboard_id", dashboard.DashboardId)`
  - [ ] `d.Set("create_time", dashboard.CreateTime)`
  - [ ] `d.Set("update_time", dashboard.UpdateTime)`
- [x] 4.8 处理标签
  - [ ] 将 API 返回的 Tags 数组转换为 map
  - [ ] `d.Set("tags", tagsMap)`
- [x] 4.9 处理可选字段
  - [ ] 检查字段是否为 nil 再设置

## 5. Update 函数实现
- [x] 5.1 实现 `resourceTencentCloudClsDashboardUpdate()` 函数
- [x] 5.2 添加日志和一致性检查
- [x] 5.3 获取 LogId
- [x] 5.4 获取资源 ID
  - [ ] 从 `d.Id()` 获取 dashboard_id
- [x] 5.5 创建 API 请求对象
  - [ ] 初始化 `ModifyDashboardRequest`
  - [ ] 设置 `DashboardId`
- [x] 5.6 处理字段变更
  - [ ] 检查 `dashboard_name` 是否变更
  - [ ] 检查 `data` 是否变更
  - [ ] 检查 `tags` 是否变更
- [x] 5.7 更新基本字段
  - [ ] 如果 `dashboard_name` 变更，设置 `request.DashboardName`
  - [ ] 如果 `data` 变更，设置 `request.Data`
- [x] 5.8 处理标签更新
  - [ ] 如果 `tags` 变更，设置 `request.Tags`
  - [ ] 将 map 转换为 Tag 数组
- [x] 5.9 调用 API 更新仪表盘
  - [ ] 使用 `resource.Retry` 包装调用
  - [ ] 使用 `tccommon.WriteRetryTimeout`
  - [ ] 添加错误处理和日志
- [x] 5.10 错误处理
  - [ ] 捕获 `InvalidParameter.DashboardNameConflict` 错误
  - [ ] 返回友好的错误信息
- [x] 5.11 调用 Read 函数刷新状态

## 6. Delete 函数实现
- [x] 6.1 实现 `resourceTencentCloudClsDashboardDelete()` 函数
- [x] 6.2 添加日志和一致性检查
- [x] 6.3 获取 LogId
- [x] 6.4 获取资源 ID
  - [ ] 从 `d.Id()` 获取 dashboard_id
- [x] 6.5 创建 API 请求对象
  - [ ] 初始化 `DeleteDashboardRequest`
  - [ ] 设置 `DashboardId`
- [x] 6.6 调用 API 删除仪表盘
  - [ ] 使用 `resource.Retry` 包装调用
  - [ ] 使用 `tccommon.WriteRetryTimeout`
  - [ ] 添加错误处理和日志
- [x] 6.7 幂等性处理
  - [ ] 如果仪表盘已不存在，视为成功
  - [ ] 记录日志但不返回错误

## 7. 测试用例
- [x] 7.1 创建 `resource_tc_cls_dashboard_test.go` 文件
- [x] 7.2 实现基础测试 `TestAccTencentCloudClsDashboard_basic`
  - [ ] 创建仪表盘（空配置）
  - [ ] 验证字段正确性
  - [ ] 更新仪表盘名称
  - [ ] 更新配置数据
  - [ ] 删除仪表盘
- [x] 7.3 实现标签测试 `TestAccTencentCloudClsDashboard_tags`
  - [ ] 创建带标签的仪表盘
  - [ ] 验证标签正确绑定
  - [ ] 更新标签
  - [ ] 删除部分标签
- [x] 7.4 实现完整配置测试 `TestAccTencentCloudClsDashboard_withData`
  - [ ] 创建包含复杂配置的仪表盘
  - [ ] 验证配置正确保存
  - [ ] 更新配置
- [x] 7.5 实现导入测试
  - [ ] 创建仪表盘
  - [ ] 导入仪表盘
  - [ ] 验证状态一致性
- [x] 7.6 创建测试用的 Terraform 配置
  - [ ] 基础配置模板
  - [ ] 更新配置模板
  - [ ] 标签配置模板

## 8. 文档编写
- [x] 8.1 创建 `resource_tc_cls_dashboard.md` 文件
- [x] 8.2 编写资源描述
  - [ ] 功能说明
  - [ ] 适用场景
- [x] 8.3 编写使用示例
  - [ ] 基础示例（空配置）
  - [ ] 完整示例（包含配置和标签）
  - [ ] 导入示例
- [x] 8.4 编写参数说明（Argument Reference）
  - [ ] 每个输入参数的详细说明
  - [ ] 必填/可选标注
  - [ ] 参数约束和验证规则
- [x] 8.5 编写属性说明（Attributes Reference）
  - [ ] 所有输出字段的说明
  - [ ] 字段类型和格式
- [x] 8.6 编写导入说明（Import）
  - [ ] 导入命令格式
  - [ ] 示例命令
- [x] 8.7 添加注意事项
  - [ ] API 频率限制
  - [ ] 名称唯一性要求
  - [ ] 标签数量限制（最多 10 个）
  - [ ] Data 字段格式说明

## 9. 注册与集成
- [x] 9.1 打开 `tencentcloud/provider.go`
- [x] 9.2 找到 CLS 资源注册区域
- [x] 9.3 添加资源注册
  - [ ] 在 `ResourcesMap` 中添加：
    ```go
    "tencentcloud_cls_dashboard": cls.ResourceTencentCloudClsDashboard(),
    ```
  - [ ] 保持字母顺序排列
- [x] 9.4 验证导入语句存在
  - [ ] `cls "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/cls"`
- [x] 9.5 更新 `tencentcloud/provider.md`
  - [ ] 在 CLS Resource 列表中添加 `tencentcloud_cls_dashboard`
  - [ ] 保持字母顺序

## 10. 代码质量
- [x] 10.1 运行 `go fmt` 格式化所有新文件
  - [ ] `go fmt ./tencentcloud/services/cls/resource_tc_cls_dashboard.go`
  - [ ] `go fmt ./tencentcloud/services/cls/resource_tc_cls_dashboard_test.go`
  - [ ] `go fmt ./tencentcloud/services/cls/service_tencentcloud_cls.go`
- [x] 10.2 运行 `golangci-lint` 检查代码
  - [ ] `golangci-lint run ./tencentcloud/services/cls/...`
  - [ ] 修复所有 linter 警告
- [x] 10.3 检查代码规范
  - [ ] 函数命名符合规范
  - [ ] 导入别名正确（tccommon, helper）
  - [ ] 日志记录完整（logId, request, response）
  - [ ] 错误处理正确
- [x] 10.4 检查字段映射
  - [ ] 所有 API 字段正确映射到 schema
  - [ ] 指针安全解引用
  - [ ] nil 值处理

## 11. OpenSpec 规范
- [x] 11.1 创建 spec delta 文件 `specs/cls-dashboard/spec.md`
- [x] 11.2 定义 ADDED Requirements
  - [ ] Req-1: 支持创建仪表盘
  - [ ] Req-2: 支持查询仪表盘信息（带分页和重试）
  - [ ] Req-3: 支持修改仪表盘
  - [ ] Req-4: 支持删除仪表盘
  - [ ] Req-5: 支持导入现有仪表盘
  - [ ] Req-6: 支持标签管理
  - [ ] Req-7: 错误处理和重试机制
- [x] 11.3 为每个 Requirement 添加 Scenario
  - [ ] 每个 Requirement 至少 2 个 Scenario
  - [ ] 覆盖正常场景和边缘场景
- [x] 11.4 运行 `openspec validate add-cls-dashboard-resource --strict`
- [x] 11.5 解决所有验证错误

## 12. 验收测试
- [x] 12.1 设置测试环境变量
  - [ ] `TENCENTCLOUD_SECRET_ID`
  - [ ] `TENCENTCLOUD_SECRET_KEY`
  - [ ] `TF_ACC=1`
- [x] 12.2 运行验收测试
  - [ ] `TF_ACC=1 go test -v -run TestAccTencentCloudClsDashboard ./tencentcloud/services/cls/`
- [x] 12.3 验证测试场景
  - [ ] 基础 CRUD 操作通过
  - [ ] 标签管理通过
  - [ ] 完整配置测试通过
  - [ ] 导入功能通过
- [x] 12.4 测试真实 API 调用
  - [ ] 创建仪表盘成功
  - [ ] 查询仪表盘返回正确数据
  - [ ] 更新仪表盘生效
  - [ ] 删除仪表盘成功
- [x] 12.5 验证分页和重试逻辑
  - [ ] 在有多个仪表盘的账户中测试分页
  - [ ] 模拟临时错误验证重试

## 13. 最终验证
- [x] 13.1 编译整个 provider
  - [ ] `make build` 或 `go build`
  - [ ] 验证没有编译错误
- [x] 13.2 运行完整 lint 检查
  - [ ] `make lint`
  - [ ] 解决所有问题
- [x] 13.3 验证文档生成
  - [ ] `make doc`
  - [ ] 检查生成的文档格式正确
- [x] 13.4 创建完整示例配置
  - [ ] 在 `examples/` 目录创建示例（如果需要）
  - [ ] 手动运行 `terraform plan` 和 `terraform apply`
  - [ ] 验证资源创建、更新、删除流程
- [x] 13.5 验证 OpenSpec 最终状态
  - [ ] 运行 `openspec show add-cls-dashboard-resource`
  - [ ] 确认所有任务完成
- [x] 13.6 准备提交
  - [ ] 编写 commit message
  - [ ] 创建 changelog 条目（如果需要）
  - [ ] 准备 PR 说明

## 14. 错误场景测试
- [x] 14.1 测试名称冲突场景
  - [ ] 创建相同名称的仪表盘
  - [ ] 验证错误处理正确
- [x] 14.2 测试标签超限场景
  - [ ] 创建超过 10 个标签的仪表盘
  - [ ] 验证错误提示清晰
- [x] 14.3 测试无效配置数据
  - [ ] 传入非 JSON 格式的 data
  - [ ] 验证错误处理
- [x] 14.4 测试删除不存在的仪表盘
  - [ ] 删除已删除的仪表盘
  - [ ] 验证幂等性
- [x] 14.5 测试分页查询
  - [ ] 在有大量仪表盘的账户中查询
  - [ ] 验证能正确找到目标仪表盘

## 任务统计

**总任务数**: 132 个子任务  
**阶段数**: 14 个主要阶段

**预计工作量**: 6 天
- 开发: 3 天
- 测试: 2 天
- 文档与集成: 1 天
