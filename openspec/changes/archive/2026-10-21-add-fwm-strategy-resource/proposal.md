## Why

云防火墙管理（FWM，Firewall Manager）产品目前无法通过 Terraform 管理下发规则策略（Strategy）。下发策略是规则下发的载体：策略 = 规则组 × 执行区域（前区 pre / 后区 post）× 接收账号集合，规则组本身由其他资源创建，本资源负责把规则组按优先级下发到指定账号/账号组。用户需要在基础设施即代码中声明、变更并销毁这些下发策略，保证可重复部署与状态收敛。

## What Changes

- 新增通用资源 `tencentcloud_fwm_strategy`（RESOURCE_KIND_GENERAL），覆盖下发策略的完整 CRUD 生命周期：
  - **Create**：调用 `CreateStrategy` 接口，按 1:1 单实例建模——每次只传 1 个 `PreStrategy` 元素（或 1 个 `PostStrategy` 元素），取响应 `StrategyIds[0]` 作为单一资源 ID。若用户配置多个元素，provider 报错提示拆分为多个资源实例。
  - **Read**：优先调用单实例查询接口 `DescribeStrategy`（按 `StrategyId` 精确单查，不存在时返回 `ResourceNotFound`）；遇 `ResourceNotFound` 或查不到时 `d.SetId("")` 移除 state 不报错，保证 import 与 destroy 幂等。
  - **Update**：调用 `ModifyStrategy`，采用增量语义——只传发生变更的字段（`GroupId`/`Sequence`/`ReceiveAccount`/`ReceiveGroup`），未变更字段保持原值。
  - **Delete**：调用 `DeleteStrategy`，传 `StrategyIds` 数组（单个删除传 1 个元素）；不存在的 ID 静默成功，destroy 天然幂等无需吞错。
- 资源标识 = 服务端生成的 `StrategyId`（`fwms-xxxx`），**非复合 ID**，import 直接传 `StrategyId`。
- `product`/`sub_account_uin`/`exec_area` 标 `Required` + `ForceNew`：变更视为重建。其中 `DescribeStrategy` 响应由 `convertStrategyToStrategyResp` 组装，**不回传 `Product` 与 `ExecArea` 字段**，因此 `product`/`exec_area` 无法从响应回读，需用 state 保留值、不参与 diff（不 set 以保留用户配置）。
- `group_id`/`sequence`/`receive_account`/`receive_group` 可变，变更走 Update 原地修改。
- 四个接口全部同步（无任务 ID、无中间态），不生成 waitReady/async 轮询逻辑。
- 新增 FWM SDK 客户端访问方法 `UseFwmV20250611Client()`（在 `connectivity/client.go` 中注册，含 import 与连接字段）。
- 在 `provider.go` / `provider.md` 中注册 `tencentcloud_fwm_strategy`。

## Capabilities

### New Capabilities
- `fwm-strategy`: FWM 下发规则策略资源的全生命周期管理（Create/Read/Update/Delete/Import），单实例 1:1 建模，同步接口，增量更新语义。

### Modified Capabilities
<!-- 无 -->

## Impact

- **新增文件**：
  - `tencentcloud/services/fwm/resource_tc_fwm_strategy.go`（资源实现）
  - `tencentcloud/services/fwm/resource_tc_fwm_strategy_test.go`（单元测试，gomonkey mock 云 API）
  - `tencentcloud/services/fwm/resource_tc_fwm_strategy.md`（文档示例）
- **修改文件**：
  - `tencentcloud/connectivity/client.go`：新增 FWM client import、连接字段与 `UseFwmV20250611Client()` 方法
  - `tencentcloud/provider.go`：新增 fwm 服务 import 与 `tencentcloud_fwm_strategy` 资源注册
  - `tencentcloud/provider.md`：新增资源文档条目（由 `make doc` 生成）
- **依赖**：`vendor/github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/fwm/v20250611`（已存在于 vendor 中）
- **SDK 包路径**：`github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/fwm/v20250611`
- **API**：CreateStrategy / DescribeStrategy（单查）/ DescribeStrategies（备用列表）/ ModifyStrategy / DeleteStrategy