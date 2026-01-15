# Change: 为 CCN 资源添加 InstanceMeteringType 参数支持

## Why

目前 Terraform Provider 的 `tencentcloud_ccn` 资源支持配置云联网的计费模式（`charge_type`：预付费/后付费）和限速类型（`bandwidth_limit_type`：地域间限速/地域出口限速），但缺少对计量模式（`InstanceMeteringType`）的支持。

根据腾讯云官方文档：
- **CreateCcn API**: https://cloud.tencent.com/document/api/215/19204
- **DescribeCcns API**: https://cloud.tencent.com/document/api/215/19199

`InstanceMeteringType` 参数用于指定云联网实例的计量类型，这是创建云联网时的重要配置项。该参数决定了实例的计费方式和带宽计量规则。

**业务场景：**
- 用户需要根据不同的业务需求选择不同的计量模式
- 计量类型影响成本核算和带宽管理策略
- 与计费模式（charge_type）配合使用，提供完整的计费配置能力

**当前限制：**
- 用户无法通过 Terraform 配置云联网的计量类型
- SDK 支持该参数，但 Provider 未暴露
- 用户只能使用云平台默认的计量模式

通过添加 `instance_metering_type` 参数支持，用户可以：
- 在创建时指定计量模式
- 通过基础设施即代码管理完整的云联网配置
- 确保 Terraform 配置与云平台实际状态一致

## What Changes

在 `tencentcloud_ccn` 资源中添加 `instance_metering_type` 可选参数，该参数：
- **类型**: String
- **可选**: Optional
- **创建时**: 传递给 CreateCcn API
- **不可修改**: ForceNew: true（修改需要重建资源）
- **查询时**: 从 DescribeCcns API 响应中读取并同步到 Terraform 状态

### 修改文件
- `tencentcloud/services/ccn/resource_tc_ccn.go` - 添加 Schema 定义和 CRUD 逻辑
- `tencentcloud/services/ccn/service_tencentcloud_ccn.go` - 更新 CreateCcn 方法签名，添加 instanceMeteringType 参数，更新 CcnBasicInfo 结构体
- `tencentcloud/services/ccn/resource_tc_ccn_test.go` - 添加测试用例覆盖

### 资源 Schema 变更
```hcl
resource "tencentcloud_ccn" "example" {
  name                    = "example-ccn"
  description             = "Example CCN"
  qos                     = "AU"
  charge_type             = "POSTPAID"
  bandwidth_limit_type    = "OUTER_REGION_LIMIT"
  instance_metering_type  = "INTER_REGION_BANDWIDTH"  # 新增参数
}
```

### API 集成
- **创建时**: 将 `instance_metering_type` 传递给 `CreateCcnRequest.InstanceMeteringType`
- **读取时**: 从 `CCN.InstanceMeteringType` 字段读取并同步到状态
- **更新时**: 不支持修改（ForceNew: true，修改会触发重建）

## Impact

### 受影响的规范
- 新增规范：`ccn-resource` - 云联网资源管理（InstanceMeteringType 参数）

### 受影响的代码
- `tencentcloud/services/ccn/resource_tc_ccn.go` - Schema 定义，Create/Read 逻辑
- `tencentcloud/services/ccn/service_tencentcloud_ccn.go` - CreateCcn 方法，CcnBasicInfo 结构体，DescribeCcns 解析逻辑
- `tencentcloud/services/ccn/resource_tc_ccn_test.go` - 添加测试覆盖

### 向后兼容性
- ✅ **完全向后兼容** - 参数为 Optional，现有配置无需修改
- ✅ **不影响现有资源** - 不修改已有资源的行为
- ✅ **SDK 已支持** - 腾讯云 SDK 已包含该字段，无需升级依赖

### API 兼容性
- ✅ `CreateCcn` API 支持 `InstanceMeteringType` 参数（可选）
- ✅ `DescribeCcns` API 返回结果包含 `InstanceMeteringType` 字段
- ✅ 参数为可选，不传递时使用云平台默认值

### 测试影响
- 需要在验收测试中覆盖 `instance_metering_type` 参数的设置和读取
- 测试用例需要验证参数值正确传递到云端并同步回状态

### 依赖关系
- 无新增依赖
- 使用现有的 SDK 版本（已包含该字段）
