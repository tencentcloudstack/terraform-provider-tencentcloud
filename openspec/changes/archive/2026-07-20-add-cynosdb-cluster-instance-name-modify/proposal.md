## Why

`tencentcloud_cynosdb_cluster` 与 `tencentcloud_cynosdb_cluster_v2` 资源中的 `instance_name` 字段当前为 `Computed`，无法通过 Terraform 修改实例名称。腾讯云 TDSQL-C MySQL 提供了 `ModifyInstanceName` 接口用于修改实例名称，应在 Terraform 中接入该接口，使 `instance_name` 支持就地更新，避免通过控制台单独修改后 state 漂移。

## What Changes

- 将 `tencentcloud_cynosdb_cluster` 与 `tencentcloud_cynosdb_cluster_v2` 资源中 `instance_name` 字段从 `Computed` 改为 `Optional` + `Computed`，使其支持用户配置与更新。
- 在 `resourceTencentCloudCynosdbClusterUpdate` 与 `resourceTencentCloudCynosdbClusterV2Update` 中新增 `d.HasChange("instance_name")` 分支，调用云 API `ModifyInstanceName` 完成实例名称修改。
- 复用 `service_tencentcloud_cynosdb.go` 中已有的 `ModifyInstanceName` 服务方法（`WriteRetryTimeout` 重试）。
- 更新文档 `website/docs/r/cynosdb_cluster.html.markdown` 与 `website/docs/r/cynosdb_cluster_v2.html.markdown`，将 `instance_name` 标记为可修改（由 `make doc` 在收尾阶段生成）。
- 在 `.changelog/4323.txt` 中追加 `instance_name` 支持修改的增强说明。
- 在 `resource_tc_cynosdb_cluster_test.go` 与 `resource_tc_cynosdb_cluster_v2_test.go` 中补充 `instance_name` 更新的 terraform acc 测试用例（cluster 与 cluster_v2 为现有资源，使用 acc 测试套件，不使用 go test 执行）。

## Capabilities

### New Capabilities
- `cynosdb-cluster-instance-name-modify`: 支持通过 `ModifyInstanceName` 接口就地修改 `tencentcloud_cynosdb_cluster` 与 `tencentcloud_cynosdb_cluster_v2` 的 `instance_name` 字段。

### Modified Capabilities
<!-- 无需修改已有 spec 层面的行为契约 -->

## Impact

- **代码文件**：
  - `tencentcloud/services/cynosdb/extension_cynosdb.go`（schema：`instance_name` 改为 `Optional` + `Computed`）
  - `tencentcloud/services/cynosdb/resource_tc_cynosdb_cluster.go`（update：新增 `instance_name` 分支）
  - `tencentcloud/services/cynosdb/resource_tc_cynosdb_cluster_v2.go`（update：新增 `instance_name` 分支）
  - `tencentcloud/services/cynosdb/resource_tc_cynosdb_cluster_test.go`（acc 测试）
  - `tencentcloud/services/cynosdb/resource_tc_cynosdb_cluster_v2_test.go`（acc 测试）
- **文档文件**：`website/docs/r/cynosdb_cluster.html.markdown`、`website/docs/r/cynosdb_cluster_v2.html.markdown`
- **changelog**：`.changelog/4323.txt`
- **API**：`ModifyInstanceName`（cynosdb v20190107，已有服务层封装，直接复用）
- **向后兼容**：仅将 `Computed` 字段改为 `Optional` + `Computed`，原有 state 与配置完全兼容，不填写 `instance_name` 时行为不变。
