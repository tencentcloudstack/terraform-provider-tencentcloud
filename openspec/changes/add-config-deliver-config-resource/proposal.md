# Add tencentcloud_config_deliver_config Resource

## What

Add a new Terraform resource `tencentcloud_config_deliver_config` for managing the Tencent Cloud Config delivery settings (投递设置). This is a global singleton configuration resource — one per account/region — that controls how Config audit logs are delivered to COS or CLS.

## Why

Users of Tencent Cloud Config need to configure log delivery to COS buckets or CLS topics for compliance auditing and analysis. This global setting currently has no Terraform support. The resource allows users to declare and manage the delivery configuration as code.

## APIs Used

| Operation | API | Notes |
|---|---|---|
| Create (initial set) | `UpdateConfigDeliver` | No dedicated Create API; first-time apply calls Update |
| Read | `DescribeConfigDeliver` | No request params — returns global singleton |
| Update | `UpdateConfigDeliver` | Idempotent — updates the global deliver config |
| Delete | No-op | No delete API; set `Status=0` to disable |

## Resource ID

Uses `helper.BuildToken()` as the resource ID (singleton — no natural unique key).
