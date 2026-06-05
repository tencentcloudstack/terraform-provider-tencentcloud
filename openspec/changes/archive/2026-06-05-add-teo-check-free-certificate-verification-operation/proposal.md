## Why

TencentCloud EdgeOne (TEO) 用户需要在 Terraform 中验证免费证书申请结果。当前 provider 缺少对 `CheckFreeCertificateVerification` 接口的支持，用户无法通过 Terraform 自动化完成免费证书验证流程。

## What Changes

- 新增 RESOURCE_KIND_OPERATION 类型资源 `tencentcloud_teo_check_free_certificate_verification`，用于调用 TEO `CheckFreeCertificateVerification` 接口检查免费证书申请结果。
- 该资源为一次性操作资源，仅有 Create 方法，Read/Update/Delete 为空实现。
- Create 方法接收 `zone_id` 和 `domain` 作为入参，调用云 API 后将返回的 `common_name`、`signature_algorithm`、`expire_time` 设置为 Computed 属性。

## Capabilities

### New Capabilities
- `teo-check-free-certificate-verification-operation`: 提供一次性操作资源，调用 TEO CheckFreeCertificateVerification 接口验证免费证书并获取证书信息。

### Modified Capabilities

(无)

## Impact

- 新增文件: `tencentcloud/services/teo/resource_tc_teo_check_free_certificate_verification_operation.go`
- 新增测试文件: `tencentcloud/services/teo/resource_tc_teo_check_free_certificate_verification_operation_test.go`
- 新增文档: `tencentcloud/services/teo/resource_tc_teo_check_free_certificate_verification_operation.md`
- 修改 `tencentcloud/provider.go`: 注册新资源
- 修改 `tencentcloud/provider.md`: 添加资源文档引用
- 依赖: `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901`（已在 vendor 中）
