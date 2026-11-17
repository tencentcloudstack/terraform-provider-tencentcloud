## Why

用户需要通过 Terraform 查询 CKafka 实例的路由信息。当前 Provider 已经支持 CKafka 路由资源的管理（`tencentcloud_ckafka_route`），但缺少查询路由信息的数据源，导致用户无法在 Terraform 配置中读取已有路由的详细信息（如接入方式、VIP 列表、域名、状态等）用于引用、校验或自动化编排。

## What Changes

- 新增 Data Source: `tencentcloud_ckafka_routes`
- 实现对 CKafka API `DescribeRoute` 接口的调用，查询指定实例的路由信息
- 支持以下输入参数：
  - `instance_id`（必填）：CKafka 实例 ID
  - `route_id`（可选）：路由 ID，用于精确查询单条路由
  - `main_route_flag`（可选）：是否显示主路由
- 返回路由信息列表 `routers`，每个路由包含：
  - `access_type`：实例接入方式
  - `route_id`：路由 ID
  - `vip_type`：路由网络类型
  - `vip_list`：虚拟 IP 列表（含 `vip`、`vport`）
  - `domain`：域名
  - `domain_port`：域名端口
  - `delete_timestamp`：删除时间戳
  - `subnet`：子网 ID
  - `broker_vip_list`：broker 虚拟 IP 列表（含 `vip`、`vport`）
  - `vpc_id`：私有网络 ID
  - `note`：备注信息
  - `status`：路由状态

## Capabilities

### New Capabilities
- `ckafka-routes-datasource`: 查询 CKafka 实例路由信息的数据源能力，封装 `DescribeRoute` 接口，支持按实例 ID、路由 ID、主路由标志查询路由列表

### Modified Capabilities
<!-- 无需修改现有 spec 级别行为 -->

## Impact

- **新增能力**: CKafka 实例路由信息查询
- **受影响的服务**: CKafka (tencentcloud/services/ckafka)
- **新增文件**:
  - `tencentcloud/services/ckafka/data_source_tc_ckafka_routes.go`
  - `tencentcloud/services/ckafka/data_source_tc_ckafka_routes.md`
  - `tencentcloud/services/ckafka/data_source_tc_ckafka_routes_test.go`
  - `tencentcloud/services/ckafka/service_tencentcloud_ckafka.go`（新增 `DescribeCkafkaRouteByFilter` 服务方法）
- **修改文件**:
  - `tencentcloud/provider.go`：注册新数据源
  - `tencentcloud/provider.md`：由 `make doc` 自动生成
- **API 依赖**:
  - CKafka API v20190819: `DescribeRoute`
  - 文档: https://cloud.tencent.com/document/api/597/40835
- **兼容性**: 无破坏性变更，纯新增功能
