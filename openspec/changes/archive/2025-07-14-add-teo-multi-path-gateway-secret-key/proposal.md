## Why

目前 Terraform Provider for TencentCloud 不支持管理 EdgeOne (TEO) 多通道安全加速网关的接入密钥配置。用户需要通过控制台手动查看和修改密钥，无法通过基础设施即代码的方式统一管理。新增 `tencentcloud_teo_multi_path_gateway_secret_key` 资源可以让用户在 Terraform 中直接管理 TEO 多通道安全加速网关的接入密钥配置，实现配置的版本化和自动化管理。

## What Changes

- 新增 Terraform 资源 `tencentcloud_teo_multi_path_gateway_secret_key`，类型为 RESOURCE_KIND_CONFIG
- 资源文件: `tencentcloud/services/teo/resource_tc_teo_multi_path_gateway_secret_key_config.go`
- 实现创建（Create）: 设置 zone_id 为资源 ID，然后调用 Update 方法
- 实现读取（Read）: 调用 `DescribeMultiPathGatewaySecretKey` 接口查询密钥配置
- 实现更新（Update）: 调用 `ModifyMultiPathGatewaySecretKey` 接口修改密钥
- 实现删除（Delete）: 从 Terraform state 中移除资源（密钥配置无法删除，仅做资源清理）
- 支持 Import 导入已有配置
- 在 `tencentcloud/provider.go` 和 `tencentcloud/provider.md` 中注册新资源
- 新增资源文档: `tencentcloud/services/teo/resource_tc_teo_multi_path_gateway_secret_key_config.md`

## Capabilities

### New Capabilities
- `teo-multi-path-gateway-secret-key-config`: 管理 TEO 多通道安全加速网关接入密钥配置的 Terraform 资源，支持读取和更新密钥

### Modified Capabilities

## Impact

- 新增资源注册代码: `tencentcloud/provider.go` 和 `tencentcloud/provider.md`
- 新增资源实现文件: `tencentcloud/services/teo/resource_tc_teo_multi_path_gateway_secret_key_config.go`
- 新增单元测试文件: `tencentcloud/services/teo/resource_tc_teo_multi_path_gateway_secret_key_config_test.go`
- 新增资源文档: `tencentcloud/services/teo/resource_tc_teo_multi_path_gateway_secret_key_config.md`
- 依赖云 API 接口: `DescribeMultiPathGatewaySecretKey`、`ModifyMultiPathGatewaySecretKey`（teo v20220901）
