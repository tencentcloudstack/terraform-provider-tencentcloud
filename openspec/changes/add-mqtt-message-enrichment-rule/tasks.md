# Tasks: Add MQTT Message Enrichment Rule Resource

## Phase 1: Service Layer Implementation (3 tasks)

- [x] 1.1 在 `service_tencentcloud_mqtt.go` 中添加 `DescribeMqttMessageEnrichmentRuleById` 方法
  - 调用 `DescribeMessageEnrichmentRules` API
  - 根据 `InstanceId` 和 `Id` 过滤返回单个规则
  - 处理规则不存在的情况
  - **验证**: 方法能正确返回指定规则或 nil

- [x] 1.2 在 `service_tencentcloud_mqtt.go` 中添加 `DeleteMqttMessageEnrichmentRuleById` 方法
  - 调用 `DeleteMessageEnrichmentRule` API
  - 使用重试机制处理删除操作
  - **验证**: 能成功删除规则并返回 nil 错误

- [x] 1.3 确保 SDK 版本支持
  - 验证 `tencentcloud-sdk-go/tencentcloud/mqtt/v20240516` 包含所需 API
  - 检查 `UseMqttV20240516Client()` 方法可用
  - **验证**: 能成功编译并调用 SDK 方法

## Phase 2: Resource Schema Definition (1 task)

- [x] 2.1 定义资源 Schema
  - Required 字段: `instance_id`, `rule_name`, `condition`, `actions`
  - Optional 字段: `status`, `remark`
  - Computed 字段: `priority`, `id` (规则ID), `created_time`, `update_time`
  - ForceNew 字段: `instance_id`
  - **验证**: Schema 定义完整且符合 API 规范

## Phase 3: CRUD Operations Implementation (4 tasks)

- [x] 3.1 实现 Create 操作
  - 解析输入参数构建 `CreateMessageEnrichmentRuleRequest`
  - `Priority` 默认设置为 1
  - Base64 编码验证 `condition` 和 `actions` 字段
  - 调用 API 创建规则
  - 设置复合 ID: `{instance_id}#{id}`
  - **验证**: 能成功创建规则并返回正确的 ID

- [x] 3.2 实现 Read 操作
  - 解析复合 ID 获取 `instance_id` 和 `id`
  - 调用服务层方法查询规则
  - 将 API 响应映射到 Terraform 状态
  - 处理规则不存在场景 (设置 ID 为空)
  - Base64 解码验证返回的 `condition` 和 `actions`
  - **验证**: Read 能正确获取和刷新资源状态

- [x] 3.3 实现 Update 操作
  - 解析复合 ID
  - 构建 `ModifyMessageEnrichmentRuleRequest`
  - 注意: 修改时需提交所有字段 (完整更新)
  - 处理 `Priority` 字段 (不可修改,使用现有值)
  - 调用 API 更新规则
  - 调用 Read 刷新状态
  - **验证**: Update 能正确修改规则属性

- [x] 3.4 实现 Delete 操作
  - 解析复合 ID
  - 调用服务层删除方法
  - 使用重试机制确保删除成功
  - **验证**: Delete 能成功删除规则

## Phase 4: Import Support (1 task)

- [x] 4.1 实现 Import 功能
  - 支持导入格式: `{instance_id}#{id}`
  - 实现 `ImportStatePassthrough`
  - **验证**: 能通过 `terraform import` 导入现有规则

## Phase 5: Testing (3 tasks)

- [x] 5.1 编写基础验收测试
  - 测试创建、读取、更新、删除流程
  - 测试字段设置和状态更新
  - 使用测试 MQTT 实例
  - **验证**: `make testacc` 测试通过

- [x] 5.2 编写 Import 测试
  - 测试导入功能
  - 验证导入后的状态正确性
  - **验证**: Import 测试通过

- [x] 5.3 添加边界条件测试
  - 测试 Base64 编码/解码
  - 测试规则不存在的场景
  - 测试必填字段验证
  - **验证**: 边界条件处理正确

## Phase 6: Documentation (2 tasks)

- [x] 6.1 编写资源文档 Markdown
  - 创建 `resource_tc_mqtt_message_enrichment_rule.md`
  - 包含完整的参数说明
  - 提供实际使用示例
  - 说明 Base64 编码要求
  - **验证**: 文档清晰易懂

- [x] 6.2 生成 Provider 文档
  - 运行 `make doc` 生成文档
  - 验证生成的文档格式正确
  - **验证**: 文档生成成功无错误

## Phase 7: Code Quality (2 tasks)

- [x] 7.1 代码格式化和静态检查
  - 运行 `make fmt` 格式化代码
  - 运行 `make lint` 检查代码质量
  - 修复所有 linting 错误和警告
  - **验证**: Lint 检查全部通过

- [x] 7.2 最终审查
  - 审查所有代码变更
  - 确认符合项目编码规范
  - 确认日志记录完整
  - 确认错误处理恰当
  - **验证**: 代码质量达到生产标准

## Summary

**总任务数**: 16 tasks
**已完成**: 16 tasks ✅
**预估时间**: 1-2 个工作日

**关键依赖**:
- 需要腾讯云 MQTT 实例用于测试
- SDK 版本已支持消息增强规则 API
- 所有任务按顺序执行,前置任务完成后才能进行后续任务

**成功标准**:
- ✅ 所有 CRUD 操作正常工作
- ✅ 验收测试全部通过
- ✅ 文档完整且准确
- ✅ 代码符合项目规范
- ✅ 无 breaking changes
