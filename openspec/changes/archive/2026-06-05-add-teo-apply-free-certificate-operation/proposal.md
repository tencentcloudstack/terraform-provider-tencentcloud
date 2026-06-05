## Why

TencentCloud EdgeOne (TEO) 用户需要通过 Terraform 申请免费证书，以便在 CNAME 接入模式下为域名配置 HTTPS。当前 Terraform Provider 缺少该操作资源，用户无法通过 IaC 方式发起免费证书申请并获取验证信息。

## What Changes

- 新增 RESOURCE_KIND_OPERATION 类型资源 `tencentcloud_teo_apply_free_certificate`，调用 TEO ApplyFreeCertificate 接口申请免费证书
- 资源为一次性操作类型，只有 Create 方法有实际逻辑，Read/Update/Delete 为空实现
- Create 方法接收 zone_id、domain、verification_method 三个必填参数，调用接口后将返回的 dns_verification 或 file_verification 验证信息写入 state
- 在 provider.go 和 provider.md 中注册该资源

## Capabilities

### New Capabilities
- `teo-apply-free-certificate`: 提供 TEO 免费证书申请操作资源，支持指定验证方式（DNS/HTTP）并返回验证信息

### Modified Capabilities

## Impact

- 新增文件: `tencentcloud/services/teo/resource_tc_teo_apply_free_certificate_operation.go`
- 新增文件: `tencentcloud/services/teo/resource_tc_teo_apply_free_certificate_operation_test.go`
- 新增文件: `tencentcloud/services/teo/resource_tc_teo_apply_free_certificate_operation.md`
- 修改文件: `tencentcloud/provider.go`（注册资源）
- 修改文件: `tencentcloud/provider.md`（添加资源文档引用）
- 依赖: `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901`（已在 vendor 中）
