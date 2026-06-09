## Why

用户需要通过 Terraform 对 EdgeOne (TEO) 站点进行认证操作。目前缺少相应的 Terraform 资源来执行站点认证，用户需要手动在控制台或通过其他方式完成站点验证。为 tencentcloud_teo_identify_zone 操作添加 Terraform 资源支持，可以将站点认证纳入 IaC 流程，实现基础设施即代码的完整覆盖。

## What Changes

- 新增 Terraform 资源 `tencentcloud_teo_identify_zone`，用于执行站点认证操作
- 资源类型为 RESOURCE_KIND_OPERATION（一次性操作）
- 实现创建接口：调用云 API `IdentifyZone` 进行站点认证
- 支持的入参：zone_name（站点名称）、domain（子域名，可选）
- 返回的认证信息：DNS 校验信息（ascription）和文件校验信息（file_ascription）
- 代码文件：resource_tc_teo_identify_zone_operation.go
- 对应的单元测试文件：resource_tc_teo_identify_zone_operation_test.go

## Capabilities

### New Capabilities

- `teo-identify-zone-operation`: EdgeOne 站点认证操作。提供通过 Terraform 执行站点认证的能力，支持 DNS 校验和文件校验两种认证方式的配置信息获取。

### Modified Capabilities

（无）

## Impact

- 新增文件：`tencentcloud/services/teo/resource_tc_teo_identify_zone_operation.go`
- 新增文件：`tencentcloud/services/teo/resource_tc_teo_identify_zone_operation_test.go`
- 新增文件：`website/docs/r/teo_identify_zone.html.markdown`（通过 `make doc` 自动生成）
- 依赖云 API：`github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901.IdentifyZone`
- 影响 TEO 服务的 Terraform Provider 功能扩展，不涉及已有资源的兼容性变更
