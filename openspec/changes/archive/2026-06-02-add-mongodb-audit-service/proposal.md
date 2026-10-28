## Why

TencentCloud MongoDB 审计服务（Audit Service）允许用户对数据库实例开启审计日志功能，记录数据库操作以满足合规和安全需求。当前 Terraform Provider 缺少对 MongoDB 审计服务的管理能力，用户无法通过 IaC 方式开通、配置和关闭审计服务。

## What Changes

- 新增 Terraform 资源 `tencentcloud_mongodb_audit_service`，支持完整的 CRUD 生命周期管理：
  - **Create**: 调用 `OpenAuditService` 接口开通审计服务，支持配置日志保存时长、全审计/规则审计模式及过滤规则
  - **Read**: 调用 `DescribeAuditConfig` 接口查询审计服务配置状态
  - **Update**: 调用 `ModifyAuditService` 接口修改审计配置（日志保存时长、审计模式、过滤规则）
  - **Delete**: 调用 `CloseAuditService` 接口关闭审计服务
- 在 `provider.go` 和 `provider.md` 中注册新资源
- 新增资源文档 `.md` 文件
- 新增单元测试文件，使用 gomonkey mock 方式验证业务逻辑

## Capabilities

### New Capabilities
- `mongodb-audit-service`: MongoDB 审计服务资源的完整 CRUD 生命周期管理，包括开通审计、查询配置、修改配置和关闭审计

### Modified Capabilities

（无需修改现有能力）

## Impact

- 新增文件：`tencentcloud/services/mongodb/resource_tc_mongodb_audit_service.go`
- 新增文件：`tencentcloud/services/mongodb/resource_tc_mongodb_audit_service_test.go`
- 新增文件：`tencentcloud/services/mongodb/resource_tc_mongodb_audit_service.md`
- 修改文件：`tencentcloud/provider.go`（注册资源）
- 修改文件：`tencentcloud/provider.md`（添加资源文档引用）
- 依赖：`github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/mongodb/v20190725`（已在 vendor 中）
