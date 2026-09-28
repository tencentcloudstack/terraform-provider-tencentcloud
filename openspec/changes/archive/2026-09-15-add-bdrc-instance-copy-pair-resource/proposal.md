## Why

腾讯云备份与容灾（BDRC）产品目前在本 Provider 中没有任何 Terraform 资源覆盖。用户无法通过 Terraform 管理 CVM 容灾复制对（Instance Copy Pair）的完整生命周期，只能依赖控制台手动创建/修改/删除复制对，无法实现容灾配置的 IaC 化。新增 `tencentcloud_bdrc_instance_copy_pair` 资源填补这一空白，使用户可以在 Terraform 中声明式地创建 CVM 复制对、查询其状态、修改复制对名称并按需销毁，从而将 BDRC 容灾复制对纳入基础设施即代码管理。

## What Changes

- 新增资源 `tencentcloud_bdrc_instance_copy_pair`（RESOURCE_KIND_GENERAL），基于 bdrc v20260330 SDK 实现完整 CRUD 生命周期管理。
- Create：调用 `CreateInstanceCopyPair` 接口创建 CVM 复制对。入参包括保护组 ID（`ProtectGroupId`）、目标端 CVM 创建参数列表（`CreateTargetInstanceParameters`，含源实例 ID、计费模式、Placement、镜像、系统盘、数据盘、VPC、公网带宽、登录设置、增强服务等全部子参数）、复制对名称（`InstanceCopyPairName`）、幂等令牌（`ClientToken`）和期望 RPO（`RecoveryPointObjective`）。创建接口为异步接口，创建完成后需调用 Read 接口（`DescribeCopyPairs`）轮询直到复制对状态生效。
- Read：调用 `DescribeCopyPairs` 接口，按 `CopyPairIds` 过滤查询复制对详情，将返回的 `CopyPairSet` 中各项属性（复制对状态、生产/容灾地域与可用区、生产/容灾资源 ID、复制进度、RPO、数据方向、创建时间等）回填到 state。
- Update：调用 `ModifyCopyPairAttribute` 接口修改复制对名称（`CopyPairName`）。该接口仅支持修改名称，其余创建参数变更需重建（ForceNew）。
- Delete：调用 `DeleteCopyPairs` 接口删除复制对，支持通过 `DeleteTargetResource` 控制是否一并删除容灾站点云盘。
- 在 `tencentcloud/connectivity/client.go` 中新增 bdrc v20260330 客户端绑定方法 `UseBdrcV20260330Client` 及对应连接字段。
- 在 `tencentcloud/provider.go` 和 `tencentcloud/provider.md` 中注册 `tencentcloud_bdrc_instance_copy_pair` 资源。
- 创建 `tencentcloud/services/bdrc/` 服务目录，包含资源文件、服务层文件、资源样例 `.md` 文件和单元测试文件。

## Capabilities

### New Capabilities
- `bdrc-instance-copy-pair-resource`: BDRC CVM 容灾复制对资源的完整生命周期管理（创建/读取/更新/删除/导入），覆盖 `CreateInstanceCopyPair`、`DescribeCopyPairs`、`ModifyCopyPairAttribute`、`DeleteCopyPairs` 四个云 API 接口的全部参数映射与异步轮询逻辑。

### Modified Capabilities
<!-- 无。本变更为纯新增资源，不修改任何已有 spec 的 requirement 级行为。 -->

## Impact

- **新增代码**:
  - `tencentcloud/services/bdrc/resource_tc_bdrc_instance_copy_pair.go`：资源 schema 定义 + CRUD 函数 + 请求构建/响应扁平化辅助函数（单文件布局，参照 `tencentcloud_igtm_strategy` 风格）。
  - `tencentcloud/services/bdrc/service_tencentcloud_bdrc.go`：服务层，包含 `NewBdrcService`、`DescribeBdrcInstanceCopyPairById` 查询辅助方法。
  - `tencentcloud/services/bdrc/resource_tc_bdrc_instance_copy_pair.md`：资源文档示例（HCL 示例 + Import 说明）。
  - `tencentcloud/services/bdrc/resource_tc_bdrc_instance_copy_pair_test.go`：单元测试（使用 gomonkey mock 云 API，不使用 terraform 测试套件）。
- **修改代码**:
  - `tencentcloud/connectivity/client.go`：新增 bdrc v20260330 包导入、`TencentCloudClient` 结构体连接字段 `bdrcv20260330Conn`、客户端工厂方法 `UseBdrcV20260330Client`。
  - `tencentcloud/provider.go`：新增 bdrc 服务包导入、在 `ResourcesMap` 中注册 `tencentcloud_bdrc_instance_copy_pair`。
  - `tencentcloud/provider.md`：新增 BDRC 产品章节及 `tencentcloud_bdrc_instance_copy_pair` 资源条目。
- **APIs consumed**: `CreateInstanceCopyPair`、`DescribeCopyPairs`、`ModifyCopyPairAttribute`、`DeleteCopyPairs`（均已在 `vendor/github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/bdrc/v20260330/` 中存在）。
- **兼容性**: 纯新增，无破坏性变更，不影响任何已有资源或 state。
- **无需 SDK 升级**: 所有所需 API 已存在于 vendored SDK 中。
