# Change: Add Monitor Notice Content Template Resource

## Why

Users need to manage custom notification content templates for Tencent Cloud Monitor service. Currently, there is no Terraform resource to create, read, update, or delete notification content templates, requiring manual configuration through the console or API calls.

## What Changes

- Add new resource `tencentcloud_monitor_notice_content_tmpl` to manage notification content templates
- Implement full CRUD operations using Monitor API endpoints
- Support complex nested template configurations for different notification channels (WeWorkRobot, DingDingRobot, FeiShuRobot, etc.)
- Use composite ID format `tmplID#tmplName` to uniquely identify resources

## Impact

- **Affected specs**: monitor-notice-content-tmpl (new capability)
- **Affected code**: 
  - `tencentcloud/services/monitor/resource_tc_monitor_notice_content_tmpl.go` (new)
  - `tencentcloud/services/monitor/service_tencentcloud_monitor.go` (add service methods)
  - `tencentcloud/provider.go` (register new resource)

## Breaking Changes

None - this is a new resource addition.

## Dependencies

- Tencent Cloud Monitor API version: 2023-06-16
- SDK: `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/monitor/v20230616`
