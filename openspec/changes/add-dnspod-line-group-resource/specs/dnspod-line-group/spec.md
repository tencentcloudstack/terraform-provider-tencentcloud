# Spec Delta: DNSPod 线路分组资源

## 概览

**能力名称**: DNSPod Line Group Management  
**资源名称**: `tencentcloud_dnspod_line_group`  
**变更类型**: ADDED  
**版本**: v1.0

---

## ADDED Requirements

### Requirement 1: 支持创建线路分组

资源必须支持通过 `CreateLineGroup` API 创建 DNSPod 线路分组。

#### Scenario: 使用域名创建基础线路分组

**Given**: 用户在 Terraform 配置中定义了线路分组资源，包含 `domain`、`name` 和 `lines` 字段  
**When**: 执行 `terraform apply`  
**Then**: 
- 调用 `CreateLineGroup` API
- 传递 `Domain`、`Name` 和 `Lines`（逗号分隔字符串）
- 成功创建线路分组
- 返回 `line_group_id`
- 资源 ID 格式为 `{domain}#{line_group_id}`

#### Scenario: 使用域名 ID 创建线路分组

**Given**: 用户配置中使用 `domain_id` 而非 `domain`  
**When**: 执行 `terraform apply`  
**Then**: 
- API 请求使用 `DomainId` 参数
- 优先级高于 `Domain` 参数
- 其他行为与基础场景一致

#### Scenario: 创建包含多个线路的分组

**Given**: 用户配置 `lines` 包含多个线路（如 ["电信", "移动", "联通"]）  
**When**: 执行 `terraform apply`  
**Then**: 
- 将 lines 列表转换为逗号分隔字符串 "电信,移动,联通"
- 传递给 API
- 所有线路正确保存到分组中

#### Scenario: 创建失败 - 分组名重复

**Given**: 域名下已存在同名线路分组  
**When**: 尝试创建同名分组  
**Then**: 
- API 返回错误 `InvalidParameter.GroupNameOccupied`
- Terraform 返回错误，包含清晰提示
- 不创建资源
- 不设置资源 ID

#### Scenario: 创建失败 - 线路已存在于其他分组

**Given**: 配置的线路已存在于另一个分组中  
**When**: 尝试创建包含该线路的分组  
**Then**: 
- API 返回错误 `InvalidParameter.LineInAnotherGroup`
- Terraform 返回错误，提示先从其他分组移除
- 不创建资源

---

### Requirement 2: 支持查询线路分组信息

资源必须通过 `DescribeLineGroupList` API 查询线路分组详情，并正确映射到 Terraform State。

#### Scenario: 查询现有线路分组

**Given**: 资源已创建，State 中存储了资源 ID  
**When**: 执行 `terraform refresh` 或 `terraform plan`  
**Then**: 
- 解析资源 ID 获取 `domain` 和 `line_group_id`
- 调用 `DescribeLineGroupList` API
- 遍历结果，找到匹配的分组
- 将字段映射到 State：
  - `domain`
  - `name`
  - `lines`（从逗号分隔字符串转换为列表）
  - `line_group_id`
  - `created_on`
  - `updated_on`

#### Scenario: 查询不存在的线路分组

**Given**: 线路分组已在云端被删除  
**When**: 执行 `terraform refresh`  
**Then**: 
- API 返回的列表中不包含该分组
- 设置 `d.SetId("")` 清空资源 ID
- 记录警告日志
- 触发 Terraform 重建资源（如果配置仍存在）

#### Scenario: 处理 API 返回的可选字段

**Given**: API 响应中某些字段为 nil  
**When**: 映射到 State  
**Then**: 
- 检查字段是否为 nil
- 仅设置非 nil 字段
- 不因 nil 值导致错误

---

### Requirement 3: 支持修改线路分组

资源必须支持通过 `ModifyLineGroup` API 更新线路分组的名称和线路列表。

#### Scenario: 修改线路分组名称

**Given**: 线路分组已存在  
**When**: 用户修改 `name` 字段并执行 `terraform apply`  
**Then**: 
- 调用 `ModifyLineGroup` API
- 传递新的 `Name` 和现有的 `Lines`
- 更新成功
- 调用 Read 函数刷新 State

#### Scenario: 修改线路列表

**Given**: 线路分组已存在  
**When**: 用户修改 `lines` 字段（添加、删除或替换线路）  
**Then**: 
- 调用 `ModifyLineGroup` API
- 传递现有的 `Name` 和新的 `Lines`（转换为逗号分隔字符串）
- 更新成功
- State 中 `lines` 字段更新为新值

#### Scenario: 同时修改名称和线路

**Given**: 线路分组已存在  
**When**: 用户同时修改 `name` 和 `lines` 字段  
**Then**: 
- 调用 `ModifyLineGroup` API 一次
- 传递新的 `Name` 和新的 `Lines`
- 更新成功

#### Scenario: 修改不可变字段 - Domain

**Given**: 线路分组已存在  
**When**: 用户修改 `domain` 字段  
**Then**: 
- Terraform 检测到 ForceNew 字段变更
- 销毁旧资源
- 创建新资源
- 新资源使用新的 domain

#### Scenario: 修改失败 - 新名称已被占用

**Given**: 目标名称已被其他分组使用  
**When**: 尝试修改为该名称  
**Then**: 
- API 返回错误 `InvalidParameter.GroupNameOccupied`
- Terraform 返回错误
- State 不变

---

### Requirement 4: 支持删除线路分组

资源必须通过 `DeleteLineGroup` API 删除线路分组，并正确处理幂等性。

#### Scenario: 删除现有线路分组

**Given**: 线路分组已存在  
**When**: 执行 `terraform destroy` 或移除配置并 apply  
**Then**: 
- 解析资源 ID
- 调用 `DeleteLineGroup` API
- 传递 `Domain` 和 `LineGroupId`
- 删除成功
- 从 State 中移除资源

#### Scenario: 删除不存在的线路分组（幂等性）

**Given**: 线路分组已在云端被删除  
**When**: Terraform 尝试删除  
**Then**: 
- API 返回错误 `InvalidParameter.LineNotExist`
- 视为成功（幂等性处理）
- 记录日志但不返回错误
- 从 State 中移除资源

#### Scenario: 删除失败 - 分组正在使用

**Given**: 线路分组正被 DNS 解析记录使用  
**When**: 尝试删除  
**Then**: 
- API 返回错误 `InvalidParameter.LineInUse`
- Terraform 返回错误，提示先移除依赖
- 资源保留在 State 中

---

### Requirement 5: 支持导入现有线路分组

资源必须支持通过 `terraform import` 导入云端已存在的线路分组到 Terraform 管理。

#### Scenario: 导入线路分组

**Given**: 云端存在线路分组，用户知道 `domain` 和 `line_group_id`  
**When**: 执行 `terraform import tencentcloud_dnspod_line_group.example example.com#123`  
**Then**: 
- 解析导入 ID 格式 `{domain}#{line_group_id}`
- 设置资源 ID
- 调用 Read 函数查询分组详情
- 将所有字段导入到 State
- 后续 `terraform plan` 显示 no changes

#### Scenario: 导入不存在的线路分组

**Given**: 用户提供的 `line_group_id` 不存在  
**When**: 执行 `terraform import`  
**Then**: 
- Read 函数查询失败
- 返回错误，提示分组不存在
- 不导入资源

#### Scenario: 导入 ID 格式错误

**Given**: 用户提供的导入 ID 格式不正确（如缺少 `#` 分隔符）  
**When**: 执行 `terraform import`  
**Then**: 
- ID 解析失败
- 返回错误，提示正确的导入格式
- 不导入资源

---

### Requirement 6: Lines 字段格式转换

资源必须正确处理 Terraform 列表格式和 API 逗号分隔字符串格式之间的转换。

#### Scenario: Create/Update - 列表转字符串

**Given**: Terraform 配置中 `lines = ["电信", "移动", "联通"]`  
**When**: 调用 Create 或 Update API  
**Then**: 
- 使用 `strings.Join(lines, ",")` 转换为 "电信,移动,联通"
- 传递给 API `Lines` 参数
- API 正确解析线路列表

#### Scenario: Read - 字符串转列表

**Given**: API 返回 `Lines: "电信,移动,联通"`  
**When**: 映射到 Terraform State  
**Then**: 
- 使用 `strings.Split(apiLines, ",")` 转换为 ["电信", "移动", "联通"]
- 设置到 State 的 `lines` 字段
- Terraform 识别为列表类型

#### Scenario: 空线路列表处理

**Given**: API 返回空字符串（虽然不应发生）  
**When**: 转换为列表  
**Then**: 
- 处理为空列表 `[]`
- 或根据业务逻辑返回错误

#### Scenario: 线路名称包含特殊字符

**Given**: 线路名称包含中文、空格或其他字符  
**When**: 进行格式转换  
**Then**: 
- 使用逗号 `,` 作为唯一分隔符
- 保留线路名称的原始字符
- 不进行转义或编码

---

### Requirement 7: 错误处理和重试机制

资源必须正确处理 API 错误，实现重试机制，并提供友好的错误信息。

#### Scenario: 速率限制重试

**Given**: API 请求超过频率限制（20次/秒）  
**When**: 调用任何 API  
**Then**: 
- 使用 `ratelimit.Check(request.GetAction())` 控制速率
- 如果超限，自动重试
- 使用 `tccommon.RetryError` 包装错误

#### Scenario: 写操作重试

**Given**: Create/Update/Delete 操作遇到临时错误  
**When**: 调用 API  
**Then**: 
- 使用 `resource.Retry` 包装调用
- 重试超时时间为 `tccommon.WriteRetryTimeout`（5分钟）
- 重试可恢复的错误（如网络错误）
- 不重试业务逻辑错误（如参数错误）

#### Scenario: 读操作重试

**Given**: Read 操作遇到临时错误  
**When**: 调用 API  
**Then**: 
- 使用 `resource.Retry` 包装调用
- 重试超时时间为 `tccommon.ReadRetryTimeout`（1分钟）
- 使用 `tccommon.RetryError(err, tccommon.InternalError)` 处理

#### Scenario: 业务错误友好提示

**Given**: API 返回业务逻辑错误（如 `InvalidParameter.GroupNameOccupied`）  
**When**: Terraform 捕获错误  
**Then**: 
- 记录详细日志（包含 request 和 error）
- 返回包含错误码和描述的错误信息
- 用户可根据错误信息定位问题

#### Scenario: 日志记录

**Given**: 任何 API 调用  
**When**: 执行操作  
**Then**: 
- 成功时记录 DEBUG 日志，包含 action、request、response
- 失败时记录 CRITICAL 日志，包含 action、request、reason
- 使用 LogId 关联日志

---

## Dependencies

### 内部依赖

- **Service 层**: 需要在 `service_tencentcloud_dnspod.go` 中添加 `DescribeDnspodLineGroupById()` 方法
- **Provider 注册**: 需要在 `provider.go` 中注册资源

### 外部依赖

- **腾讯云 SDK**: `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dnspod/v20210323`
  - 已验证 SDK 包含所需 API
- **API 版本**: `2021-03-23`
- **最低 Terraform 版本**: 0.13.x（项目要求）

---

## Non-functional Requirements

### 性能

- Create 操作在 10 秒内完成（正常网络条件下）
- Read 操作在 5 秒内完成
- Update 操作在 10 秒内完成
- Delete 操作在 10 秒内完成

### 可靠性

- 实现幂等性（Delete 操作可重复执行）
- 自动重试临时错误
- 速率限制保护

### 安全性

- 不在日志中暴露敏感信息
- 使用项目标准的认证机制

### 兼容性

- 遵循 Terraform Plugin SDK v2 规范
- 遵循项目代码规范和命名约定
- 向后兼容现有 DNSPod 资源

---

## Acceptance Criteria

资源满足以下所有标准即视为完成：

1. ✅ 所有 CRUD 操作功能正常
2. ✅ 导入功能正常工作
3. ✅ 所有 Scenarios 通过测试
4. ✅ 验收测试通过（`TF_ACC=1 go test`）
5. ✅ 代码通过 `golangci-lint` 检查
6. ✅ 代码格式化（`go fmt`）
7. ✅ 文档完整且格式正确
8. ✅ 资源已在 provider.go 中注册
9. ✅ Service 层方法实现并测试
10. ✅ 错误处理覆盖所有已知错误码
