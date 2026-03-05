# Monitor Notice Content Template Resource

## ADDED Requirements

### Requirement: Resource Schema Definition
资源 MUST 提供完整的 Schema 定义，支持腾讯云监控通知内容模板的所有核心参数。

#### Scenario: 定义基础属性
- **WHEN** 用户配置 `tencentcloud_monitor_notice_content_tmpl` 资源
- **THEN** 资源支持以下属性:
  - `tmpl_name` (String, Required) - 模板名称
  - `monitor_type` (String, Required) - 监控类型（如 MT_QCE）
  - `tmpl_language` (String, Required) - 模板语言（zh 或 en）
  - `tmpl_contents` (Map, Required) - 模板内容配置

#### Scenario: 定义复杂嵌套结构
- **WHEN** 用户配置 `tmpl_contents` 参数
- **THEN** 系统支持以下通知渠道的嵌套配置:
  - `qcloud_yehe` - 腾讯云渠道（包含 email, sms, voice, wechat, qywx, site, andon 子渠道）
  - `wework_robot` - 企业微信机器人
  - `dingding_robot` - 钉钉机器人
  - `feishu_robot` - 飞书机器人
- **AND** 每个渠道包含 `matching_status` 列表和 `template` 配置

### Requirement: Resource ID Management
资源 MUST 使用复合 ID 格式以确保唯一性和可追溯性。

#### Scenario: 创建资源时生成复合 ID
- **WHEN** 用户执行 `terraform apply` 创建资源
- **THEN** 系统调用 CreateNoticeContentTmpl API
- **AND** API 返回 `TmplID`
- **AND** 系统生成复合 ID 格式为 `{tmplID}#{tmplName}`
- **AND** 设置资源的 Terraform ID 为该复合 ID

#### Scenario: 读取资源时解析复合 ID
- **WHEN** Terraform 执行 refresh 或 read 操作
- **THEN** 系统从资源 ID 中解析出 `tmplID` 和 `tmplName`
- **AND** 将 `tmplID` 转换为字符串列表传递给 DescribeNoticeContentTmpl 的 `TmplIDs` 参数
- **AND** 将 `tmplName` 传递给 `TmplName` 参数进行查询

### Requirement: Create Operation
资源 MUST 实现标准的 Terraform Create 操作，调用腾讯云 CreateNoticeContentTmpl API。

#### Scenario: 成功创建通知内容模板
- **WHEN** 用户定义资源并执行 `terraform apply`
- **THEN** 系统从 Schema 中读取 `tmpl_name`, `monitor_type`, `tmpl_language`, `tmpl_contents` 参数
- **AND** 构建 CreateNoticeContentTmplRequest 请求对象
- **AND** 将 `tmpl_contents` 的复杂嵌套结构转换为 API 要求的格式
- **AND** 调用 CreateNoticeContentTmpl API
- **AND** API 返回 `TmplID`
- **AND** 设置资源 ID 为 `TmplID`
- **AND** 执行 Read 操作更新状态

#### Scenario: 创建失败时返回错误
- **WHEN** API 调用失败（如参数不合法、权限不足）
- **THEN** 系统返回包含详细错误信息的 Terraform error
- **AND** 资源未被添加到状态文件中

### Requirement: Read Operation
资源 MUST 实现标准的 Terraform Read 操作，查询模板的最新状态。

#### Scenario: 成功读取模板信息
- **WHEN** Terraform 执行 refresh 或 read 操作
- **THEN** 系统从资源 ID 解析 `tmplID` 和 `tmplName`
- **AND** 构建 DescribeNoticeContentTmplRequest，将 `tmplID` 放入 `TmplIDs` 字符串列表参数
- **AND** 将 `tmplName` 设置到 `TmplName` 参数
- **AND** 调用 DescribeNoticeContentTmpl API 进行查询
- **AND** 如果返回结果不为空，将 API 响应的字段更新到 Terraform state
- **AND** 将 `TmplContents` 反序列化为 Schema 定义的嵌套结构

#### Scenario: 资源不存在时移除状态
- **WHEN** DescribeNoticeContentTmpl API 返回空结果
- **THEN** 系统调用 `d.SetId("")` 将资源从状态中移除
- **AND** 记录日志表明资源已被外部删除

#### Scenario: 查询失败时重试
- **WHEN** API 调用由于临时错误失败（如网络问题）
- **THEN** 系统使用 `resource.Retry` 和 `tccommon.ReadRetryTimeout` 进行重试
- **AND** 在重试函数内部执行 `ratelimit.Check` 限流检查
- **AND** 返回 `tccommon.RetryError(e)` 触发重试
- **AND** 响应为 nil 时返回 `resource.NonRetryableError`

### Requirement: Update Operation
资源 MUST 实现标准的 Terraform Update 操作，调用 ModifyNoticeContentTmpl API。

#### Scenario: 成功更新模板配置
- **WHEN** 用户修改资源配置并执行 `terraform apply`
- **THEN** 系统检测到 `d.HasChange()` 返回 true
- **AND** 从资源 ID 解析出 `tmplID`
- **AND** 构建 ModifyNoticeContentTmplRequest，包含 `TmplID`, `TmplName`, `TmplContents`
- **AND** 将更新后的 `tmpl_contents` 转换为 API 要求的格式
- **AND** 调用 ModifyNoticeContentTmpl API
- **AND** API 成功后，执行 Read 操作刷新状态

#### Scenario: 更新失败时返回错误
- **WHEN** ModifyNoticeContentTmpl API 调用失败
- **THEN** 系统返回错误信息
- **AND** Terraform 状态保持不变

### Requirement: Delete Operation
资源 MUST 实现标准的 Terraform Delete 操作，调用 DeleteNoticeContentTmpls API。

#### Scenario: 成功删除模板
- **WHEN** 用户执行 `terraform destroy` 删除资源
- **THEN** 系统从资源 ID 解析出 `tmplID`
- **AND** 构建 DeleteNoticeContentTmplsRequest，将 `tmplID` 放入 `TmplIDs` 列表参数
- **AND** 调用 DeleteNoticeContentTmpls API
- **AND** API 成功返回后，资源从状态中移除

#### Scenario: 删除已不存在的资源
- **WHEN** 资源已被外部删除
- **THEN** 系统忽略 ResourceNotFound 类型错误
- **AND** 正常完成删除操作

### Requirement: Service Layer Implementation
Service 层 MUST 提供封装良好的方法，处理 API 调用、重试机制和错误处理。

#### Scenario: CreateNoticeContentTmpl 方法
- **WHEN** Resource 层调用 service 的 CreateNoticeContentTmpl 方法
- **THEN** service 方法构建 API request
- **AND** 执行 ratelimit.Check
- **AND** 调用腾讯云 SDK 的 CreateNoticeContentTmpl
- **AND** 记录 Debug 日志包含 action 和 request ID
- **AND** 返回响应的 TmplID

#### Scenario: DescribeNoticeContentTmplByFilter 方法
- **WHEN** Resource 层调用 service 的 DescribeNoticeContentTmplByFilter 方法
- **THEN** service 方法支持通过 `tmplIDs` 和 `tmplName` 过滤
- **AND** 实现重试机制，在重试函数内部执行 ratelimit.Check
- **AND** 响应为空时返回 NonRetryableError
- **AND** 成功后在重试函数内部记录日志
- **AND** 返回匹配的模板列表

#### Scenario: ModifyNoticeContentTmpl 方法
- **WHEN** Resource 层调用 service 的 ModifyNoticeContentTmpl 方法
- **THEN** service 方法构建包含 TmplID, TmplName, TmplContents 的 request
- **AND** 执行 ratelimit.Check
- **AND** 调用腾讯云 SDK 的 ModifyNoticeContentTmpl
- **AND** 记录操作日志

#### Scenario: DeleteNoticeContentTmpl 方法
- **WHEN** Resource 层调用 service 的 DeleteNoticeContentTmpl 方法
- **THEN** service 方法接受 tmplID 参数
- **AND** 构建 DeleteNoticeContentTmplsRequest，将 tmplID 放入列表
- **AND** 执行 ratelimit.Check
- **AND** 调用腾讯云 SDK 的 DeleteNoticeContentTmpls
- **AND** 记录操作日志

### Requirement: Error Handling and Logging
系统 MUST 提供完善的错误处理和日志记录机制。

#### Scenario: 记录操作耗时
- **WHEN** 任何 Resource 操作开始执行
- **THEN** 使用 `defer tccommon.LogElapsed()` 记录操作耗时

#### Scenario: 状态一致性检查
- **WHEN** Create 或 Update 操作执行
- **THEN** 使用 `defer tccommon.InconsistentCheck()` 进行状态一致性检查

#### Scenario: 详细错误日志
- **WHEN** API 调用失败
- **THEN** 日志包含 action 名称、request ID、错误详情
- **AND** 错误信息清晰易于排查

### Requirement: Testing Coverage
资源 MUST 包含完整的验收测试，覆盖所有 CRUD 操作。

#### Scenario: 基础 CRUD 测试
- **WHEN** 执行验收测试
- **THEN** 测试用例包含:
  - 创建资源并验证状态
  - 读取资源并验证所有字段
  - 更新资源配置并验证变更
  - 删除资源并验证清理

#### Scenario: 多渠道配置测试
- **WHEN** 执行测试
- **THEN** 测试覆盖不同通知渠道的配置（邮件、短信、企业微信、钉钉、飞书）
- **AND** 验证复杂嵌套结构的正确性

### Requirement: Documentation
资源 MUST 提供完整的文档，包括使用示例和参数说明。

#### Scenario: 资源文档内容
- **WHEN** 用户查看资源文档
- **THEN** 文档包含:
  - 资源描述和用途
  - 完整的参数列表和类型说明
  - 各参数的必填性和默认值
  - 至少一个完整的使用示例
  - 复杂嵌套结构的示例

#### Scenario: 使用示例
- **WHEN** 用户查看 examples 目录
- **THEN** 示例展示典型的配置场景
- **AND** 示例代码可以直接运行
- **AND** 包含必要的注释说明
