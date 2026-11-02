## Why

云防火墙纳管（Cloud Firewall Managed，fwm）产品目前缺少 Terraform 资源来管理「集团成员账号是否被云防火墙纳管」这一个开关。用户需要在 Terraform 中声明某个集团成员的纳管期望状态（开启/关闭），以便通过 IaC 统一管控成员纳管状态，避免手动在控制台逐个操作。该资源属于配置型资源（RESOURCE_KIND_CONFIG）：成员实体本就存在于集团中，资源只管理纳管开关的读取与更新，不存在创建/销毁实体的动作。

## What Changes

- 新增 Terraform 资源 `tencentcloud_fwm_cfw_managed`，用于管理单个集团成员的云防火墙纳管开关。
- Read：通过 `DescribeOrganMembers` 接口，使用 `Filters=[{Name:"AppId",Values:[appId],OperatorType:1}]` 按 AppId 过滤定位成员列表，再从中按 MemberId 匹配目标成员读取 `cfw_managed`。
- Update：通过 `ModifyCfwManaged` 接口设置 `CfwManaged`（0 关闭 / 1 开启）。写操作为同步接口，无需轮询。
- Create：empty 空实现，仅登记 ID（`d.SetId(member_id)`），随后复用 Update 调用设开关（参考配置型资源 `tencentcloud_cfw_edge_firewall_switch` 的 Create→Update 模式）。
- Delete：empty 空实现，从 state 移除即可，不调用云 API（不能「删除纳管」，只能关开关）。
- 资源标识 = 目标成员的 MemberId；支持 import，import 使用 `app_id:member_id` 两段式（Read 依赖 app_id 定位）。
- 在 connectivity 层新增 `UseFwmClient()` 方法与 `fwmConn` 字段（fwm 产品首次接入该 Provider）。
- 在 `provider.go` / `provider.md` 中注册该资源。

## Capabilities

### New Capabilities
- `fwm-cfw-managed-config`: 管理单个集团成员的云防火墙纳管开关配置（读取 + 更新），包含 schema 定义、CRUD 骨架、import 支持、connectivity 客户端接入、provider 注册与文档。

### Modified Capabilities
<!-- 无现有 spec 的 requirement 发生变更 -->

## Impact

- **新增代码**：
  - `tencentcloud/services/fwm/resource_tc_fwm_cfw_managed_config.go`（资源实现）
  - `tencentcloud/services/fwm/resource_tc_fwm_cfw_managed_config_test.go`（mock 单测，gomonkey）
  - `tencentcloud/services/fwm/resource_tc_fwm_cfw_managed_config.md`（文档样例）
  - `tencentcloud/services/fwm/service_tencentcloud_fwm.go`（服务层，封装 Describe/Modify 逻辑）
- **修改代码**：
  - `tencentcloud/connectivity/client.go`：新增 fwm 包导入、`fwmConn` 字段、`UseFwmClient()` 方法
  - `tencentcloud/provider.go`：注册 `tencentcloud_fwm_cfw_managed`
  - `tencentcloud/provider.md`：新增资源条目
- **依赖**：`github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/fwm/v20250611`（已在 go.mod / vendor / tools.go 中就绪）
- **云 API**：`DescribeOrganMembers`（查询）、`ModifyCfwManaged`（设置开关）
- **权限**：调用账号必须是集团管理员，否则返回 UnauthorizedOperation
