## Why

`tencentcloud_cynosdb_readonly_instance` 资源当前的 `instance_name` 字段被标记为 `ForceNew`，用户修改实例名称时会导致资源重建。腾讯云 TDSQL-C MySQL 提供了 `ModifyInstanceName` 接口用于修改实例名称，应在 Terraform 中接入该接口，使 `instance_name` 支持就地更新，避免不必要的重建。

## What Changes

- 移除 `tencentcloud_cynosdb_readonly_instance` 资源中 `instance_name` 字段的 `ForceNew`，使其支持更新。
- 在 `resourceTencentCloudCynosdbReadonlyInstanceUpdate` 中新增 `d.HasChange("instance_name")` 分支，调用云 API `ModifyInstanceName` 完成实例名称修改。
- 在 `service_tencentcloud_cynosdb.go` 中新增 `ModifyInstanceName` 服务方法，封装 `ModifyInstanceName` API 调用（含 `WriteRetryTimeout` 重试）。
- 更新文档 `website/docs/r/cynosdb_readonly_instance.html.markdown`，将 `instance_name` 的 `ForceNew` 标记移除。
- 更新 `.changelog/4319.txt`，将变更描述从"update instance_name field description"改为支持 `instance_name` 修改的增强说明。
- 在 `resource_tc_cynosdb_readonly_instance_test.go` 中补充 `instance_name` 更新的单元测试用例���使用 gomonkey mock 云 API）。

## Capabilities

### New Capabilities
- `cynosdb-readonly-instance-name-modify`: 支持通过 `ModifyInstanceName` 接口就地修改 `tencentcloud_cynosdb_readonly_instance` 的 `instance_name` 字段。

### Modified Capabilities
<!-- 无需修改已有 spec 层面的行为契约 -->

## Impact

- **代码文件**：
  - `tencentcloud/services/cynosdb/resource_tc_cynosdb_readonly_instance.go`（schema、update）
  - `tencentcloud/services/cynosdb/service_tencentcloud_cynosdb.go`（新增服务方法）
  - `tencentcloud/services/cynosdb/resource_tc_cynosdb_readonly_instance_test.go`（单元测试）
- **文档文件**：`website/docs/r/cynosdb_readonly_instance.html.markdown`
- **changelog**：`.changelog/4319.txt`
- **API**：`ModifyInstanceName`（cynosdb v20190107）
- **向后兼容**：仅将 `ForceNew` 字段改为可更新，原有 state 与配置完全兼容，不会触发重建。
