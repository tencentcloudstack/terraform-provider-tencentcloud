## Why

腾讯云防火墙（CFW）集群防火墙支持 Bypass 模式配置，允许用户控制流量是否绕过防火墙。目前 Terraform Provider 缺少对该配置的管理能力，用户无法通过 IaC 方式管理集群防火墙的 Bypass 状态。

## What Changes

- 新增 `tencentcloud_cfw_cluster_fw_bypass_config` CONFIG 类型资源，支持读取和更新集群防火墙的 Bypass 状态配置
- 资源通过 `DescribeClusterNatCcnFwSwitchList` 接口读取当前 Bypass 状态
- 资源通过 `ModifyClusterFwBypass` 接口更新 Bypass 开关状态
- Schema 仅展示 `ModifyClusterFwBypass` 接口入参（`fw_type`、`ccn_id`、`enable`、`nat_ins_id`），不展示 Read 接口入参

## Capabilities

### New Capabilities
- `cfw-cluster-fw-bypass-config`: 管理腾讯云 CFW 集群防火墙 Bypass 状态配置，支持读取 NAT CCN 集群模式防火墙开关列表，并修改集群防火墙的 Bypass 开关（true-开启Bypass绕过防火墙，false-关闭Bypass经过防火墙）

### Modified Capabilities

## Impact

- 新增文件: `tencentcloud/services/cfw/resource_tc_cfw_cluster_fw_bypass_config.go`
- 新增文件: `tencentcloud/services/cfw/resource_tc_cfw_cluster_fw_bypass_config_test.go`
- 新增文件: `tencentcloud/services/cfw/resource_tc_cfw_cluster_fw_bypass_config.md`
- 修改文件: `tencentcloud/provider.go`（注册新资源）
- 修改文件: `tencentcloud/provider.md`（添加资源文档引用）
- 依赖云 API: `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cfw/v20190904`
