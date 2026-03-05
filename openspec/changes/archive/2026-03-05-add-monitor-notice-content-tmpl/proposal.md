# Change: Add Monitor Notice Content Template Resource

## Why
腾讯云监控服务需要支持自定义通知内容模板功能，允许用户通过 Terraform 管理告警通知的内容模板。当前 Provider 缺少该资源类型，无法满足用户对告警通知内容定制化的需求。

## What Changes
- 新增 `tencentcloud_monitor_notice_content_tmpl` resource 资源
- 实现完整的 CRUD 操作，支持创建、查询、更新和删除通知内容模板
- 支持多种通知渠道配置（邮件、短信、企业微信、钉钉、飞书等）
- 使用单一 TmplID 作为资源标识符，简化 ID 管理
- tmpl_name 支持更新（未设置 ForceNew）
- 优化 flatten 逻辑，只在有值时设置 template 子字段，避免资源漂移
- 添加相应的 service 层方法和重试机制

## Impact
- **新增文件**:
  - `tencentcloud/services/monitor/resource_tc_monitor_notice_content_tmpl.go` - 资源实现
  - `tencentcloud/services/monitor/resource_tc_monitor_notice_content_tmpl_test.go` - 验收测试
  - `examples/tencentcloud-monitor-notice-content-tmpl/main.tf` - 使用示例
  - `website/docs/r/monitor_notice_content_tmpl.html.markdown` - 资源文档
- **修改文件**:
  - `tencentcloud/provider.go` - 注册新资源
  - `tencentcloud/services/monitor/service_tencentcloud_monitor.go` - 添加 service 方法
- **新增 API 调用**:
  - CreateNoticeContentTmpl (创建)
  - DescribeNoticeContentTmpl (查询)
  - ModifyNoticeContentTmpl (修改)
  - DeleteNoticeContentTmpls (删除)
- **依赖变更**: 需要使用腾讯云 SDK monitor 服务 v20230616 版本的 API
