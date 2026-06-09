## Why

用户需要通过 Terraform 检查 TEO（TencentCloud EdgeOne）域名的 CNAME 配置状态，以验证配置是否正确生效。这是域名接入和配置验证过程中的关键步骤，目前缺少对应的 Terraform 操作资源支持。

## What Changes

- 新增 Terraform 操作资源 `tencentcloud_teo_check_cname_status`，用于检查 TEO 域名的 CNAME 配置状态
- 支持通过 `zone_id` 和 `record_names` 参数检查一个或多个域名的 CNAME 状态
- 返回 `cname_status` 结果，包含每个域名的记录名称（record_name）、CNAME地址（cname）和状态（status）
- 资源类型为 RESOURCE_KIND_OPERATION（一次性操作，不需要持久化状态，仅需实现 Create 接口，RUD 接口为空）

## Capabilities

### New Capabilities
- `teo-check-cname-status`: 新增检查 TEO 域名 CNAME 状态的能力，支持批量检查多个域名，返回每个域名的 CNAME 配置状态（active/moved）

### Modified Capabilities
（无现有能力变更）

## Impact

**新增文件：**
- `tencentcloud/services/teo/resource_tc_teo_check_cname_status_operation.go` - 操作资源主文件
- `tencentcloud/services/teo/resource_tc_teo_check_cname_status_operation_test.go` - 操作资源单元测试文件
- `website/docs/r/teo_check_cname_status.html.md` - 资源文档

**修改文件：**
- `tencentcloud/services/teo/service_tencentcloud_teo.go` - 可能需要注册新资源到服务

**依赖：**
- 使用云 API: `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901` 的 `CheckCnameStatus` 接口
- 云 API 入参: `ZoneId` (string), `RecordNames` ([]string)
- 云 API 出参: `CnameStatus` ([]CnameStatus)，其中 CnameStatus 包含 `RecordName`, `Cname`, `Status`
