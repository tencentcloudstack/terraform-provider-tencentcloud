## Why

TCR (Tencent Container Registry) 的 CreateInstance 接口已新增 `EnableCosMAZ` 和 `EnableCosVersioning` 两个字段，用于控制 COS 存储桶的多可用区和版本控制特性。当前 terraform provider 的 `tencentcloud_tcr_instance` 资源未同步这些字段，导致用户无法通过 Terraform 配置这些重要的存储特性，需要补齐这个能力缺失。

## What Changes

- **新增 schema 字段**: 在 `tencentcloud_tcr_instance` 资源中添加 `enable_cos_maz` 和 `enable_cos_versioning` 两个 Optional + Computed 字段
- **Create 逻辑**: 在 `resourceTencentCloudTcrInstanceCreate` 中支持从 schema 读取这两个字段并传递给 CreateInstance API
- **Read 逻辑**: 在 `resourceTencentCloudTcrInstanceRead` 中从 API 响应中读取这两个字段并 set 回 state
- **文档更新**: 更新 `website/docs/r/tcr_instance.html.markdown` 文档，添加新字段说明

## Capabilities

### New Capabilities
- `tcr-instance-cos-config`: TCR 实例 COS 存储桶配置能力，支持多可用区和版本控制特性

### Modified Capabilities
<!-- 无修改现有能力 -->

## Impact

**受影响代码:**
- `tencentcloud/services/tcr/resource_tc_tcr_instance.go` - 资源定义和 CRUD 逻辑
- `website/docs/r/tcr_instance.html.markdown` - 资源文档

**API 依赖:**
- TCR CreateInstance 接口 (已支持这两个字段)
- TCR DescribeInstances 接口 (已返回这两个字段)

**向后兼容性:**
- ✅ 完全兼容：新增字段为 Optional，不影响现有配置
- ✅ 现有资源的 state 会在下次 refresh 时自动补充这两个字段的值
