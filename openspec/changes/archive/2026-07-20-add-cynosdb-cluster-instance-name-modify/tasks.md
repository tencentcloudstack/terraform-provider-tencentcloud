## 1. Schema 调整

- [x] 1.1 在 `tencentcloud/services/cynosdb/extension_cynosdb.go` 的 `TencentCynosdbInstanceBaseInfo()` 中，将 `instance_name` 字段从 `Computed: true` 改为 `Optional: true, Computed: true`，描述更新为 "Name of instance. Only supported when modifying."

## 2. 资源 Update 逻辑

- [x] 2.1 在 `tencentcloud/services/cynosdb/resource_tc_cynosdb_cluster.go` 的 `resourceTencentCloudCynosdbClusterUpdate` 中，新增 `if d.HasChange("instance_name")` 分支，读取新值并调用 `cynosdbService.ModifyInstanceName(ctx, instanceId, instanceName)`，处理 error
- [x] 2.2 在 `tencentcloud/services/cynosdb/resource_tc_cynosdb_cluster_v2.go` 的 `resourceTencentCloudCynosdbClusterV2Update` 中，新增 `if d.HasChange("instance_name")` 分支，读取新值并调用 `cynosdbService.ModifyInstanceName(ctx, instanceId, instanceName)`，处理 error

## 3. 测试用例

- [x] 3.1 在 `tencentcloud/services/cynosdb/resource_tc_cynosdb_cluster_test.go` 中新增 `TestAccTencentCloudCynosdbClusterResourceUpdateInstanceName`，使用 terraform acc 测试套件，验证创建后修改 `instance_name` 触发 ModifyInstanceName 且 state 正确回填
- [x] 3.2 在 `tencentcloud/services/cynosdb/resource_tc_cynosdb_cluster_v2_test.go` 中新增 `TestAccTencentCloudCynosdbClusterV2ResourceUpdateInstanceName`，使用 terraform acc 测试套件，验证创建后修改 `instance_name` 触发 ModifyInstanceName 且 state 正确回填

> 注：cluster 与 cluster_v2 为现有资源，按代码生成要求使用 terraform acc 测试套件补充用例，不使用 go test 执行。

## 4. 文档与 changelog

- [x] 4.1 更新 `.changelog/4323.txt`，将变更描述改为 `resource/tencentcloud_cynosdb_cluster: support instance_name modification via ModifyInstanceName API`（enhancement），同时保留 cluster_v2 描述
- [x] 4.2 更新 `website/docs/r/cynosdb_cluster.html.markdown` 与 `website/docs/r/cynosdb_cluster_v2.html.markdown` 中 `instance_name` 的描述（由 `make doc` 在收尾阶段生成）

## 5. 验证

- [x] 5.1 确认所有修改的 Go 文件可正确编译（由后续流程的 go build 验证）
