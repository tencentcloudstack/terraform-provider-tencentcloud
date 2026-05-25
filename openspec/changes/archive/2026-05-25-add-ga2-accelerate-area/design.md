## Context

TencentCloud GA2（全球加速 2.0）产品提供了加速地域管理的云 API，包括 CreateAccelerateAreas、DescribeAccelerateAreas、ModifyAccelerateAreas、DeleteAccelerateAreas 四个接口。这些接口已在 vendor 中可用（`github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ga2/v20250115`）。

当前 provider 中尚无 GA2 相关资源，需要新建 `tencentcloud/services/ga2/` 目录并实现资源。

关键约束：
- Create/Modify/Delete 接口均为异步接口，返回 TaskId，需要轮询 DescribeAccelerateAreas 确认操作生效
- 加速地域挂载在全球加速实例下，使用 `global_accelerator_id` 作为资源 ID
- Delete 接口需要传入 `AcceleratorAreaIds`（从 Read 接口获取）

## Goals / Non-Goals

**Goals:**
- 实现 `tencentcloud_ga2_accelerate_area` 资源的完整 CRUD 生命周期管理
- 支持异步操作轮询等待机制
- 遵循 provider 现有代码风格（参考 tencentcloud_igtm_strategy）
- 提供单元测试（使用 gomonkey mock 云 API）
- 提供资源文档（.md 文件）

**Non-Goals:**
- 不实现 GA2 其他资源（如监听器、终端节点组等）
- 不实现数据源
- 不实现验收测试（使用 mock 单元测试替代）

## Decisions

### 1. 资源 ID 设计
**决策**: 使用 `global_accelerator_id` 作为 `d.SetId()` 的值。

**理由**: 加速地域是全球加速实例的子资源，一个实例下可以有多个加速地域。资源管理的是某个实例下的所有加速地域配置，而非单个加速地域。Create/Modify/Delete 接口都以 `GlobalAcceleratorId` 为主键操作。

### 2. Schema 设计
**决策**:
- `global_accelerator_id` (Required, ForceNew, String): 全球加速实例 ID，作为资源主键
- `accelerator_areas` (Required, List): 加速地域配置列表，包含子字段：
  - `accelerate_region` (Required, String): 加速地域
  - `bandwidth` (Required, Int): 带宽
  - `isp_type` (Optional, String): ISP 类型，支持 BGP/三网/精品，默认 BGP
  - `ip_version` (Optional, String): IP 版本，默认 IPv4
- `accelerate_area_set` (Computed, List): 查询返回的加速地域信息（含 AcceleratorAreaId、IpAddress 等只读字段）

**理由**: 将用户可配置的参数（accelerator_areas）与只读返回值（accelerate_area_set）分离，符合 Terraform 资源设计最佳实践。

### 3. 异步操作处理
**决策**: Create/Modify/Delete 操作完成后，调用 resourceTencentCloudGa2AccelerateAreaRead 轮询 DescribeAccelerateAreas 接口确认操作生效。

**理由**: 三个写操作接口都是异步的（返回 TaskId），需要等待操作完成后才能确认资源状态。通过在 Create/Update/Delete 末尾调用 Read 方法实现轮询。

### 4. Delete 实现
**决策**: Delete 时先调用 DescribeAccelerateAreas 获取当前所有 AcceleratorAreaId，然后调用 DeleteAccelerateAreas 传入所有 ID 进行删除。

**理由**: DeleteAccelerateAreas 接口需要 AcceleratorAreaIds 参数，而这些 ID 是由服务端生成的，需要先查询获取。

### 5. 测试策略
**决策**: 使用 gomonkey 对云 API 进行 mock，编写单元测试验证业务逻辑。

**理由**: 新增资源使用 mock 方式进行单元测试，避免依赖真实云环境。

## Risks / Trade-offs

- [异步操作超时] → 在 schema 中声明 Timeouts 块，使用合理的默认超时时间（Create/Update/Delete 各 10 分钟）
- [分页查询遗漏] → DescribeAccelerateAreas 支持分页，Read 方法中需要实现自动分页获取所有数据
- [并发操作冲突] → 依赖云 API 自身的并发控制（InstanceStateNotAllowedOperate 错误码），通过 retry 机制处理
