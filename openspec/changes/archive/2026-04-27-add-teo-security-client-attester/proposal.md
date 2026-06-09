## Why

TEO (EdgeOne) 产品需要支持客户端认证选项（Security Client Attester）的 Terraform 管理。当前用户只能通过控制台或 API 来创建、查询、修改和删除客户端认证选项，无法通过 Terraform 进行基础设施即代码管理。新增此资源可以让用户在 Terraform 中完整管理 TEO 客户端认证选项的全生命周期。

## What Changes

- 新增 Terraform 资源 `tencentcloud_teo_security_client_attester`，支持 TEO 客户端认证选项的 CRUD 操作
  - Create: 调用 `CreateSecurityClientAttester` 接口创建客户端认证选项
  - Read: 调用 `DescribeSecurityClientAttester` 接口查询客户端认证选项
  - Update: 调用 `ModifySecurityClientAttester` 接口修改客户端认证选项
  - Delete: 调用 `DeleteSecurityClientAttester` 接口删除客户端认证选项
- 在 `provider.go` 和 `provider.md` 中注册新资源
- 新增资源文档 `.md` 文件

## Capabilities

### New Capabilities
- `teo-security-client-attester`: TEO 客户端认证选项资源的 Terraform 管理，支持创建、查询、修改和删除客户端认证选项，包含三种认证方式（TC-RCE、TC-CAPTCHA、TC-EO-CAPTCHA）的配置

### Modified Capabilities
（无已有能力需要修改）

## Impact

- 新增文件: `tencentcloud/services/teo/resource_tc_teo_security_client_attester.go`
- 新增文件: `tencentcloud/services/teo/resource_tc_teo_security_client_attester_test.go`
- 新增文件: `tencentcloud/services/teo/resource_tc_teo_security_client_attester.md`
- 修改文件: `tencentcloud/provider.go`（注册新资源）
- 修改文件: `tencentcloud/provider.md`（添加资源注释）
- 依赖云API: `teo/v20220901` 包中的 `CreateSecurityClientAttester`、`DescribeSecurityClientAttester`、`ModifySecurityClientAttester`、`DeleteSecurityClientAttester` 接口
