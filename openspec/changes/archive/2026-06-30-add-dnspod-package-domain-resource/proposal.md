## Why

DNSPod 提供了套餐与域名绑定管理功能，允许用户将购买的付费套餐绑定到指定域名以获得增值 DNS 服务。目前 Terraform Provider 缺少对套餐域名绑定关系的管理能力，用户无法通过 IaC 方式对套餐进行域名绑定、换绑和解绑操作，只能通过控制台手动操作，效率低且不可追溯。

## What Changes

- 新增 `tencentcloud_dnspod_package_domain` 资源，支持套餐与域名的完整生命周期管理
- Create: 通过 `ModifyPackageDomain` API（Operation="bind"）实现套餐绑定域名
- Read: 通过 `DescribeDomainVipList` API 查询套餐绑定的域名信息
- Update: 通过 `ModifyPackageDomain` API（Operation="change"）实现套餐更换域名
- Delete: 通过 `ModifyPackageDomain` API（Operation="unbind"）实现套餐解绑域名
- Import: 支持导入现有套餐域名绑定关系

## Capabilities

### New Capabilities
- `dnspod-package-domain-resource`: 提供 DNSPod 套餐域名绑定资源的 Terraform 管理能力，包括绑定、换绑、解绑和查询套餐域名关系

### Modified Capabilities
<!-- No existing capabilities are modified -->

## Impact

- **新增文件**: `tencentcloud/services/dnspod/resource_tc_dnspod_package_domain.go` — 资源实现
- **新增文件**: `tencentcloud/services/dnspod/resource_tc_dnspod_package_domain_test.go` — 单元测试
- **新增文件**: `tencentcloud/services/dnspod/resource_tc_dnspod_package_domain.md` — 资源文档
- **新增文件**: `openspec/changes/add-dnspod-package-domain-resource/specs/dnspod-package-domain-resource/spec.md` — 规范文件
- **修改文件**: `tencentcloud/provider.go` — 注册新资源到 Provider
- **修改文件**: `tencentcloud/provider.md` — 添加资源文档索引
- **依赖**: `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dnspod/v20210323`（已存在）
