## 1. Service 层 (service_tencentcloud_vdb.go)

- [x] 1.1 定义 `VdbClientInterface` 接口，抽象所需 SDK client 方法: `DescribeInstancesWithContext`, `DescribeInstanceNodesWithContext`, `DescribeDBSecurityGroupsWithContext`, `CreateInstanceWithContext`, `ScaleUpInstanceWithContext`, `ScaleOutInstanceWithContext`, `IsolateInstanceWithContext`, `DestroyInstancesWithContext`, `ModifyDBInstanceSecurityGroupsWithContext`
- [x] 1.2 定义 `VdbService` struct，持有 `client *connectivity.TencentCloudClient` 和 `vdbClient VdbClientInterface`（用于 mock 注入），实现 `getVdbClient()` 方法
- [x] 1.3 实现 `DescribeVdbInstanceById` — 调用 `DescribeInstances`，Limit 设为最大值 100，返回 `*vdb.InstanceInfo`（空列表返回 nil）
- [x] 1.4 实现 `DescribeVdbInstanceNodesById` — 调用 `DescribeInstanceNodes`，Limit 设为最大值 100
- [x] 1.5 实现 `DescribeDBSecurityGroupsByInstanceId` — 调用 `DescribeDBSecurityGroups`
- [x] 1.6 实现轮询方法: `WaitForInstanceStatus`, `WaitForInstanceScaleUp`, `WaitForInstanceScaleOut`, `WaitForInstanceNotFound`, `WaitForSecurityGroupsMatch`

## 2. 资源 CRUD 实现 (resource_tc_vdb_instance.go)

- [x] 2.1 定义 `ResourceTencentCloudVdbInstance()` schema，包含: 必填字段 (vpc_id/subnet_id/pay_mode/security_group_ids)、可选字段 (instance_name/pay_period/auto_renew/params/resource_tags/instance_type/mode/product_type/node_type/cpu/memory/disk_size/worker_node_num/force_delete)、Computed 字段 (status/region/zone/product/shard_num/api_version/extend/expired_at/is_no_expired/wan_address/isolate_at/task_status/networks/nodes/created_at/engine_name/engine_version)、Timeouts 块 (Create 30min, Update 30min, Delete 10min)、Importer
- [x] 2.2 实现 `resourceTencentCloudVdbInstanceCreate` — 构建 CreateInstanceRequest，retry 调用，检查 InstanceIds 非空，`d.SetId()` 后轮询等待 online 状态
- [x] 2.3 实现 `resourceTencentCloudVdbInstanceRead` — 调用 service 查询，nil 时打印 `[CRUD]` 日志后 SetId("")，逐字段 nil 检查后 set，回填 networks/resource_tags/security_group_ids/nodes 子结构
- [x] 2.4 实现 `resourceTencentCloudVdbInstanceUpdate` — immutableFields 检查、ScaleUp (cpu/memory/disk_size)、ScaleOut (worker_node_num)、ModifyDBInstanceSecurityGroups (security_group_ids)，各步骤间轮询等待
- [x] 2.5 实现 `resourceTencentCloudVdbInstanceDelete` — IsolateInstance + 轮询等待 isolated，force_delete 时 DestroyInstances + 轮询等待 not found

## 3. Provider 注册

- [x] 3.1 确认 `tencentcloud/provider.go` 中已导入 vdb service 包并注册 `tencentcloud_vdb_instance` 资源
- [x] 3.2 确认 `tencentcloud/provider.md` 中已包含资源条目

## 4. 资源文档 (resource_tc_vdb_instance.md)

- [x] 4.1 创建/更新资源文档，参考 `resource_tc_config_compliance_pack.md` 命名格式，包含: 一句话描述（带 VDB 产品名）、Example Usage（按量付费和包月两种场景）、Import 部分

## 5. 单元测试 (resource_tc_vdb_instance_test.go)

- [x] 5.1 创建 MockVdbClient (实现 VdbClientInterface, 使用 testify/mock)，包含所有接口方法的 mock 实现
- [x] 5.2 编写 helper 函数: `newTestService`, `strPtr`, `uint64Ptr`, `float64Ptr`, `int64Ptr`, `makeInstanceInfo`, `makeDescribeResponse`
- [x] 5.3 编写 `DescribeVdbInstanceById` 测试用例: Found, NotFound, NilItems, ApiError, NilResponse
- [x] 5.4 编写 `DescribeVdbInstanceNodesById` 测试用例: Found, Empty
- [x] 5.5 编写 `DescribeDBSecurityGroupsByInstanceId` 测试用例: Found, Empty
- [x] 5.6 编写 `WaitForInstanceStatus` 测试用例: AlreadyOnline, CaseInsensitive, IsolatedNotFound, DestroyedNotFound
- [x] 5.7 编写 `WaitForInstanceScaleUp` 测试用例: AlreadyMatched, ValuesConverge, NilCpuMemoryDisk
- [x] 5.8 编写 `WaitForInstanceScaleOut` 测试用例: AlreadyMatched, ValuesConverge
- [x] 5.9 编写 `WaitForInstanceNotFound` 测试用例: AlreadyGone, EventuallyGone
- [x] 5.10 编写 `WaitForSecurityGroupsMatch` 测试用例: AlreadyMatched, EventuallyMatch
- [x] 5.11 编写 Schema 验证测试: 必填字段、可选字段、Computed 字段、ForceNew 字段、可更新字段、force_delete 默认值、CRUD 函数非 nil、Timeouts 非 nil
- [x] 5.12 编写 networks/nodes 子字段验证测试

## 6. 代码正确性检查

- [x] 6.1 检查 CreateInstance 的所有入参在云API request struct 中存在
- [x] 6.2 检查 ScaleUpInstance 的入参 (InstanceId/Cpu/Memory/StorageSize/RunNow) 在云API request struct 中存在
- [x] 6.3 检查 ScaleOutInstance 的入参 (InstanceId/ReplicaNum/RunNow) 在云API request struct 中存在
- [x] 6.4 检查 IsolateInstance 的入参 (InstanceId) 在云API request struct 中存在
- [x] 6.5 检查 DestroyInstances 的入参 (InstanceIds) 在云API request struct 中存在
- [x] 6.6 检查 ModifyDBInstanceSecurityGroups 的入参 (SecurityGroupIds/InstanceIds) 在云API request struct 中存在
- [x] 6.7 检查所有函数返回的 error 均被处理，不会出错的用 `_ =` 赋值