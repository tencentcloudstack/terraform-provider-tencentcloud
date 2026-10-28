## Context

堡垒机（BH）产品提供了 `BindDeviceResource` 接口用于修改资产绑定的堡垒机服务。该接口同时承担绑定（Create）、更新绑定（Update）和解绑（Delete）的功能。Read 操作通过 `DescribeDevices` 接口查询设备详情获取绑定状态。

当前 vendor 中已包含 `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/bh/v20230418` 包，无需额外引入依赖。

## Goals / Non-Goals

**Goals:**
- 实现 `tencentcloud_bh_bind_device_resource` 资源的完整 CRUD 生命周期管理
- 支持将多个设备资产绑定到指定堡垒机服务实例
- 支持 K8S 集群托管场景的相关参数配置
- 在 provider 中注册新资源
- 提供单元测试（使用 gomonkey mock）

**Non-Goals:**
- 不实现数据源（data source）
- 不处理异步轮询（该接口为同步接口）
- 不修改现有 BH 资源的行为

## Decisions

### 1. CRUD 接口映射

| 操作 | 云API接口 | 说明 |
|------|-----------|------|
| Create | BindDeviceResource | 传入 resource_id 绑定设备到堡垒机服务 |
| Read | DescribeDevices | 通过 IdSet 查询设备详情，从 Device.Resource 获取绑定信息 |
| Update | BindDeviceResource | 更新绑定参数（同 Create 接口） |
| Delete | BindDeviceResource | 将 ResourceId 设为空字符串解绑 |

**理由**：该接口设计为幂等操作，CUD 使用同一接口，通过 ResourceId 是否为空来区分绑定/解绑。

### 2. 资源 ID 设计

使用 `device_id_set` 和 `resource_id` 作为联合 ID，格式为：`<device_ids>#<resource_id>`，其中多个 device_id 用逗号连接。

- 分隔符使用 `tccommon.FILED_SP`（即 `#`）
- device_id_set 中的多个 ID 用 `,` 连接

**理由**：该资源的唯一性由设备集合和目标堡垒机服务共同决定。

### 3. ForceNew 策略

由于 CUD 使用同一接口，所有字段均可在 Update 时传入，因此只将 `device_id_set` 设为 ForceNew（设备集合变更意味着绑定关系本身变化）。其余字段（resource_id、domain_id、manage_dimension 等）允许原地更新。

### 4. Read 实现

通过 `DescribeDevices` 接口传入 `IdSet`（取 device_id_set 中的第一个 ID）查询设备详情。从返回的 `Device` 结构体中读取：
- `Device.Resource` → 获取绑定的堡垒机服务信息（resource_id）
- `Device.DomainId` → domain_id
- `Device.ManageDimension` → manage_dimension
- `Device.ManageAccountId` → manage_account_id
- `Device.Namespace` → namespace
- `Device.Workload` → workload

注意：`manage_account`、`manage_kubeconfig` 为写入参数，DescribeDevices 返回中不包含，因此在 Read 中不设置这两个字段。

### 5. 单元测试策略

使用 gomonkey 对 `BindDeviceResource` 和 `DescribeDevices` 接口进行 mock，验证 CRUD 业务逻辑正确性。

## Risks / Trade-offs

- [Risk] DescribeDevices 返回的 Device 结构体中 Resource 字段可能为 nil（设备未绑定服务时） → 在 Read 中判断 nil，若为 nil 则认为资源已被删除，调用 `d.SetId("")` 移除 state
- [Risk] manage_account 和 manage_kubeconfig 为敏感字段且 Read 接口不返回 → 在 schema 中标记为 Sensitive，不在 Read 中覆盖已有值
- [Risk] device_id_set 为数组，联合 ID 中需要正确序列化/反序列化 → 使用逗号分隔，在 Read/Update/Delete 中正确解析
