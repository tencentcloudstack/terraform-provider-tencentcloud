# 变更提案：tencentcloud_bh_resource 新增 public_ip_set 和 private_ip_set 字段

## 变更类型

**功能增强** — 在 `tencentcloud_bh_resource` 资源中补充 `DescribeResources` 接口返回的 `PublicIpSet` 和 `PrivateIpSet` 字段。

## Why

`DescribeResources` 接口（https://cloud.tencent.com/document/api/1025/74801）在 `ResourceSet` 中返回以下字段，当前 Read 模块均未读取：

| SDK 字段 | 类型 | 说明 |
|---------|------|------|
| `PublicIpSet` | `[]*string` | 堡垒机实例的公网 IP 地址集合 |
| `PrivateIpSet` | `[]*string` | 堡垒机实例的内网 IP 地址集合 |

用户无法通过 Terraform 获取堡垒机的 IP 地址用于后续资源配置（如安全组规则、DNS 解析等）。

## What Changes

| 文件 | 修改内容 |
|------|---------|
| `tencentcloud/services/bh/resource_tc_bh_resource.go` | schema 新增 `public_ip_set`、`private_ip_set`（均为 Computed TypeList）；Read 模块新增对应赋值 |

### 字段设计

| Terraform 字段 | SDK 字段 | 类型 | 属性 | 说明 |
|--------------|---------|------|------|------|
| `public_ip_set` | `PublicIpSet` | `TypeList[TypeString]` | Computed | 堡垒机公网 IP 列表 |
| `private_ip_set` | `PrivateIpSet` | `TypeList[TypeString]` | Computed | 堡垒机内网 IP 列表 |

### 向后兼容性

✅ 完全向后兼容：两个字段均为纯 Computed，不影响已有配置。
