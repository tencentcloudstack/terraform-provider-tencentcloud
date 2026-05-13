## Why

TEO (EdgeOne) 产品当前在 Terraform Provider 中缺少对负载均衡实例(Load Balancer)资源的管理能力。用户需要通过 Terraform 来创建、读取、更新和删除 TEO 负载均衡实例，以实现站点加速场景下的源站流量调度、健康检查和故障转移等能力。

## What Changes

- 新增 Terraform 通用资源 `tencentcloud_teo_load_balancer`，支持 TEO 负载均衡实例的完整 CRUD 生命周期管理
  - Create: 调用 `CreateLoadBalancer` 接口创建负载均衡实例
  - Read: 调用 `DescribeLoadBalancerList` 接口按 InstanceId 过滤查询实例详情
  - Update: 调用 `ModifyLoadBalancer` 接口更新实例配置
  - Delete: 调用 `DeleteLoadBalancer` 接口删除实例
- 资源支持以下核心配置：
  - 站点 ID (zone_id)：指定所属站点
  - 实例名称 (name)：负载均衡实例名称
  - 实例类型 (type)：HTTP 专用型 / 通用型
  - 源站组列表 (origin_groups)：源站组及优先级配置
  - 健康检查策略 (health_checker)：支持 HTTP/HTTPS/TCP/UDP/ICMP Ping 等多种检查方式
  - 流量调度策略 (steering_policy)：按优先级故障转移
  - 请求重试策略 (failover_policy)：跨源站组重试 / 同源站组内重试
- 在 `provider.go` 和 `provider.md` 中注册该资源
- 生成资源文档 `.md` 文件

## Capabilities

### New Capabilities
- `teo-load-balancer-resource`: TEO 负载均衡实例资源的 CRUD 管理，包括创建、读取、更新、删除负载均衡实例，支持源站组配置、健康检查策略、流量调度策略和请求重试策略等核心功能

### Modified Capabilities

## Impact

- 新增文件：`tencentcloud/services/teo/resource_tc_teo_load_balancer.go`
- 新增文件：`tencentcloud/services/teo/resource_tc_teo_load_balancer_test.go`
- 新增文件：`tencentcloud/services/teo/resource_tc_teo_load_balancer.md`
- 修改文件：`tencentcloud/provider.go`（注册新资源）
- 修改文件：`tencentcloud/provider.md`（文档更新）
- 依赖云 API：`teo/v20220901` 包中的 CreateLoadBalancer、DeleteLoadBalancer、DescribeLoadBalancerList、ModifyLoadBalancer 接口
