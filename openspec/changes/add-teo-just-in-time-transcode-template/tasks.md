## 1. 资源文件创建

- [x] 1.1 创建 `tencentcloud/services/teo/resource_tc_teo_just_in_time_transcode_template.go` 文件
- [x] 1.2 创建 `tencentcloud/services/teo/resource_tc_teo_just_in_time_transcode_template_test.go` 文件
- [x] 1.3 创建 `tencentcloud/services/teo/resource_tc_teo_just_in_time_transcode_template.md` 示例文件

## 2. Schema 定义实现

- [x] 2.1 在资源文件中定义 Resource 函数，注册资源 schema
- [x] 2.2 定义基础参数 schema（zone_id, template_name, comment）
- [x] 2.3 定义开关参数 schema（video_stream_switch, audio_stream_switch）
- [x] 2.4 定义 video_template block schema 及其嵌套字段
- [x] 2.5 定义 audio_template block schema 及其嵌套字段
- [x] 2.6 定义 computed 参数（template_id, create_time, update_time）
- [x] 2.7 定义 Timeouts block（create, delete）
- [x] 2.8 设置 ForceNew 标记为所有输入参数
- [x] 2.9 添加参数验证函数（ValidateFunc）

## 3. Create 函数实现

- [x] 3.1 实现 resourceTencentcloudTeoJustInTimeTranscodeTemplateCreate 函数
- [x] 3.2 实现 resourceTencentcloudTeoJustInTimeTranscodeTemplateCreateCtx 函数（支持 context）
- [x] 3.3 调用 CreateJustInTimeTranscodeTemplate API
- [x] 3.4 映射 Terraform 参数到 API 请求参数
- [x] 3.5 解析 API 响应，获取 template_id
- [x] 3.6 设置资源 ID 为 zone_id#template_id
- [x] 3.7 添加错误处理和重试逻辑（使用 tccommon.WriteRetryTimeout）
- [x] 3.8 添加日志记录（使用 tccommon.LogElapsed）
- [x] 3.9 添加一致性检查（使用 tccommon.InconsistentCheck）
- [x] 3.10 实现异步创建验证逻辑，轮询 DescribeJustInTimeTranscodeTemplates

## 4. Read 函数实现

- [x] 4.1 实现 resourceTencentcloudTeoJustInTimeTranscodeTemplateRead 函数
- [x] 4.2 实现 resourceTencentcloudTeoJustInTimeTranscodeTemplateReadCtx 函数（支持 context）
- [x] 4.3 解析资源 ID，提取 zone_id 和 template_id
- [x] 4.4 调用 DescribeJustInTimeTranscodeTemplates API
- [x] 4.5 设置过滤条件（Filters），按 template_id 过滤
- [x] 4.6 映射 API 响应参数到 Terraform schema
- [x] 4.7 处理模板不存在的情况（返回 nil）
- [x] 4.8 映射嵌套的 video_template 和 audio_template 参数
- [x] 4.9 添加错误处理和重试逻辑（使用 tccommon.Retry）
- [x] 4.10 添加日志记录（使用 tccommon.LogElapsed）
- [x] 4.11 添加一致性检查（使用 tccommon.InconsistentCheck）

## 5. Delete 函数实现

- [x] 5.1 实现 resourceTencentcloudTeoJustInTimeTranscodeTemplateDelete 函数
- [x] 5.2 实现 resourceTencentcloudTeoJustInTimeTranscodeTemplateDeleteCtx 函数（支持 context）
- [x] 5.3 解析资源 ID，提取 zone_id 和 template_id
- [x] 5.4 调用 DeleteJustInTimeTranscodeTemplates API
- [x] 5.5 添加错误处理和重试逻辑（使用 tccommon.WriteRetryTimeout）
- [x] 5.6 添加日志记录（使用 tccommon.LogElapsed）
- [x] 5.7 添加一致性检查（使用 tccommon.InconsistentCheck）
- [x] 5.8 实现异步删除验证逻辑，轮询 DescribeJustInTimeTranscodeTemplates

## 6. 资源 ID 处理

- [x] 6.1 实现资源 ID 编码逻辑（zone_id#template_id）
- [x] 6.2 实现资源 ID 解析逻辑（split on '#'）
- [x] 6.3 添加 malformed ID 的错误处理
- [x] 6.4 处理 zone_id 或 template_id 包含 '#' 的情况

## 7. 单元测试实现

- [x] 7.1 测试成功创建场景（最小参数）
- [x] 7.2 测试成功创建场景（完整参数）
- [x] 7.3 测试缺失必需参数的验证
- [x] 7.4 测试参数验证（video_stream_switch 和 audio_stream_switch）
- [x] 7.5 测试空参数的验证
- [x] 7.6 测试参数长度限制验证
- [x] 7.7 测试 API 错误处理
- [x] 7.8 测试速率限制错误的重试逻辑
- [x] 7.9 测试认证错误处理（不重试）
- [x] 7.10 测试网络超时的重试逻辑
- [x] 7.11 测试异步创建验证逻辑
- [x] 7.12 测试成功读取场景
- [x] 7.13 测试读取不存在的模板
- [x] 7.14 测试读取过程中的速率限制处理
- [x] 7.15 测试成功删除场景
- [x] 7.16 测试删除不存在的模板
- [x] 7.17 测试异步删除验证逻辑
- [x] 7.18 测试超时处理（创建和删除）
- [x] 7.19 测试资源 ID 解析逻辑
- [x] 7.20 测试 malformed 资源 ID 处理

## 8. 资源注册

- [x] 8.1 在 `tencentcloud/services/teo/tencentcloud_teo_schema.go` 中注册新资源
- [x] 8.2 确保资源名称符合命名规范

## 9. 示例文件更新

- [x] 9.1 编写最小配置示例（仅必需参数）
- [x] 9.2 编写完整配置示例（所有参数）
- [x] 9.3 编写视频流仅启用示例
- [x] 9.4 编写音频流仅启用示例
- [x] 9.5 添加参数说明和注释

## 10. 文档生成

- [x] 10.1 运行 `make doc` 命令生成 `website/docs/r/teo_just_in_time_transcode_template.md`
- [x] 10.2 验证生成的文档内容完整
- [x] 10.3 确认文档包含所有参数说明和示例

## 11. 验证任务

- [x] 11.1 运行 `gofmt` 格式化代码
- [x] 11.2 检查代码编译是否通过（检查导入和语法）
- [x] 11.3 运行单元测试验证功能正确性
- [x] 11.4 检查资源注册是否成功
- [x] 11.5 验证示例文件的配置正确性

## 12. 代码审查准备

- [x] 12.1 确保代码符合项目编码规范
- [x] 12.2 添加必要的注释和文档
- [x] 12.3 确认错误消息清晰明确
- [x] 12.4 确认日志记录充分
- [x] 12.5 确认所有测试通过
