# 设计文档：tencentcloud_mysql_proxy_address_config

## 资源概览

- **资源名称**: `tencentcloud_mysql_proxy_address_config`
- **资源类型**: Config（无真实 Create/Delete API，仅有 Update）
- **Package**: `cdb`
- **文件位置**: `tencentcloud/services/cdb/`
- **唯一 ID**: `InstanceId#ProxyGroupId`

## Schema 设计

```hcl
resource "tencentcloud_mysql_proxy_address_config" "example" {
  instance_id      = "cdb-xxxxxxxx"
  proxy_group_id   = "proxy-xxxxxxxx"
  proxy_address_id = "proxyaddr-xxxxxxxx"
  weight_mode      = "system"
  is_kick_out      = false
  min_count        = 0
  max_delay        = 100
  fail_over        = true
  auto_add_ro      = true
  read_only        = false
  trans_split      = false
  connection_pool  = false
  auto_load_balance = false
  access_mode      = "nearby"
}
```

### 字段映射表（与 AdjustCdbProxyAddressRequest 一一对应）

| Schema 字段 | SDK 字段 | 类型 | Required/Optional | ForceNew | 说明 |
|-------------|----------|------|-------------------|----------|------|
| `instance_id` | — | String | Required | true | InstanceId（DescribeCdbProxyInfo 必传，用于 ID 组成和 Read 校验） |
| `proxy_group_id` | `ProxyGroupId` | String | Required | true | 代理组 ID |
| `proxy_address_id` | `ProxyAddressId` | String | Required | true | 代理地址 ID |
| `weight_mode` | `WeightMode` | String | Required | false | 权重分配模式：`system`/`custom` |
| `is_kick_out` | `IsKickOut` | Bool | Required | false | 是否开启延迟剔除 |
| `min_count` | `MinCount` | Int | Required | false | 最小保留数量 |
| `max_delay` | `MaxDelay` | Int | Required | false | 延迟剔除阈值（ms） |
| `fail_over` | `FailOver` | Bool | Required | false | 是否开启故障转移 |
| `auto_add_ro` | `AutoAddRo` | Bool | Required | false | 是否自动添加 RO |
| `read_only` | `ReadOnly` | Bool | Required | false | 是否只读 |
| `trans_split` | `TransSplit` | Bool | Optional | false | 是否开启事务分离 |
| `connection_pool` | `ConnectionPool` | Bool | Optional | false | 是否开启连接池 |
| `proxy_allocation` | `ProxyAllocation` | List | Optional | false | 读写权重分配 |
| `auto_load_balance` | `AutoLoadBalance` | Bool | Optional | false | 是否开启自适应负载均衡 |
| `access_mode` | `AccessMode` | String | Optional | false | 访问模式：`nearby`/`balance` |
| `ap_node_as_ro_node` | `ApNodeAsRoNode` | String | Optional | false | 是否将 libra 节点当作普通 RO |
| `ap_query_to_other_node` | `ApQueryToOtherNode` | String | Optional | false | libra 故障时是否转发给其他节点 |

### proxy_allocation 子结构

| Schema 字段 | SDK 字段 | 类型 | 说明 |
|-------------|----------|------|------|
| `region` | `Region` | String | 地域 |
| `zone` | `Zone` | String | 可用区 |
| `proxy_instance` | `ProxyInstance` | List | 实例列表 |

#### proxy_instance 子结构

| Schema 字段 | SDK 字段 | 类型 | 说明 |
|-------------|----------|------|------|
| `instance_id` | `InstanceId` | String | 实例 ID |
| `weight` | `Weight` | Int | 权重 |

## CRUD 实现

### Create

Config 型资源，Create 函数：
1. 从 schema 读取 `instance_id`、`proxy_group_id`
2. `d.SetId(strings.Join([]string{instanceId, proxyGroupId}, tccommon.FILED_SP))`
3. 直接调用 `resourceTencentCloudMysqlProxyAddressConfigUpdate`

### Read (`DescribeCdbProxyInfo`)

1. 解析 `d.Id()` 得到 `instanceId`、`proxyGroupId`
2. 调用 service 层 `DescribeMysqlProxyAddressConfig(ctx, instanceId, proxyGroupId, proxyAddressId)`
3. 若返回 nil，`d.SetId("")`，返回 nil
4. 回填所有 schema 字段

### Update (`AdjustCdbProxyAddress`)

1. 解析 `d.Id()` 得到 `instanceId`、`proxyGroupId`
2. **ID 校验**：调用 `DescribeCdbProxyInfo(instanceId, proxyGroupId)`，验证能唯一查询到目标代理组，否则返回错误
3. 构建 `AdjustCdbProxyAddressRequest`，填入所有必选/可选字段
4. `resource.Retry(WriteRetryTimeout, ...)` 调用 `AdjustCdbProxyAddressWithContext`
5. 调用 Read

### Delete（no-op）

Config 型资源，Delete 函数仅返回 nil，不调用任何 API。

## Service 层

在 `service_tencentcloud_mysql.go` 新增：

```go
func (me *MysqlService) DescribeMysqlProxyAddressConfig(
    ctx context.Context,
    instanceId, proxyGroupId, proxyAddressId string,
) (address *cdb.ProxyAddress, errRet error)
```

逻辑：
1. 调用 `DescribeCdbProxyInfo(instanceId, proxyGroupId)`
2. 遍历 `ProxyInfos[0].ProxyAddress`，找到 `ProxyAddressId == proxyAddressId` 的条目
3. 返回该 `*ProxyAddress`

## 代码风格

严格参考 `resource_tc_waf_owasp_rule_status_config.go`：
- Create = 设置 ID + 调用 Update
- Delete = 空函数（返回 nil）
- Update 中解析 ID、调用 service 校验、构建 request、Retry 调用
- 使用 `tccommon.NewResourceLifeCycleHandleFuncContext`
- 使用 `tccommon.LogElapsed` / `tccommon.InconsistentCheck`
- 使用 `resource.Retry(tccommon.WriteRetryTimeout, ...)`

## 文件清单

| 文件 | 类型 |
|------|------|
| `tencentcloud/services/cdb/service_tencentcloud_mysql.go` | 追加 service 方法 |
| `tencentcloud/services/cdb/resource_tc_mysql_proxy_address_config.go` | 新建 Resource 文件 |
| `tencentcloud/services/cdb/resource_tc_mysql_proxy_address_config.md` | 新建 Resource 文档 |
| `tencentcloud/services/cdb/resource_tc_mysql_proxy_address_config_test.go` | 新建测试文件 |
| `tencentcloud/provider.go` | 注册 Resource |
