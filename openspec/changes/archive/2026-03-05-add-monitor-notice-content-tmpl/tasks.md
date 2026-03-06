# Implementation Tasks

## 1. Service Layer Implementation
- [x] 1.1 在 `service_tencentcloud_monitor.go` 中添加 `CreateNoticeContentTmpl` 方法
- [x] 1.2 在 `service_tencentcloud_monitor.go` 中添加 `DescribeNoticeContentTmplByFilter` 方法
- [x] 1.3 在 `service_tencentcloud_monitor.go` 中添加 `ModifyNoticeContentTmpl` 方法
- [x] 1.4 在 `service_tencentcloud_monitor.go` 中添加 `DeleteNoticeContentTmpl` 方法
- [x] 1.5 实现 API 调用的重试机制和错误处理

## 2. Resource Implementation
- [x] 2.1 创建 `resource_tc_monitor_notice_content_tmpl.go` 文件
- [x] 2.2 实现资源 Schema 定义（包含 tmpl_name, monitor_type, tmpl_language, tmpl_contents）
- [x] 2.3 实现 `resourceTencentCloudMonitorNoticeContentTmplCreate` 方法
- [x] 2.4 实现 `resourceTencentCloudMonitorNoticeContentTmplRead` 方法
- [x] 2.5 实现 `resourceTencentCloudMonitorNoticeContentTmplUpdate` 方法
- [x] 2.6 实现 `resourceTencentCloudMonitorNoticeContentTmplDelete` 方法
- [x] 2.7 使用单一 TmplID 作为资源 ID（已优化）
- [x] 2.8 实现 tmpl_contents 复杂嵌套结构的 Schema 映射
- [x] 2.9 优化 flattenNoticeContentTmplItem 函数，避免资源漂移
- [x] 2.10 移除 tmpl_name 的 ForceNew 属性，支持更新

## 3. Provider Registration
- [x] 3.1 在 `provider.go` 的 ResourcesMap 中注册新资源

## 4. Testing
- [x] 4.1 创建 `resource_tc_monitor_notice_content_tmpl_test.go` 文件
- [x] 4.2 实现基础 CRUD 验收测试（创建、读取、更新、删除）
- [x] 4.3 实现不同通知渠道配置的测试用例
- [x] 4.4 运行并通过所有验收测试

## 5. Examples and Documentation
- [x] 5.1 创建 `examples/tencentcloud-monitor-notice-content-tmpl/main.tf` 示例文件
- [x] 5.2 创建资源文档 `website/docs/r/monitor_notice_content_tmpl.html.markdown`
- [x] 5.3 文档中包含完整的参数说明和使用示例

## 6. Code Quality
- [x] 6.1 运行 `make fmt` 进行代码格式化
- [x] 6.2 运行 `make lint` 通过代码检查
- [x] 6.3 确保所有日志和错误信息完整清晰
