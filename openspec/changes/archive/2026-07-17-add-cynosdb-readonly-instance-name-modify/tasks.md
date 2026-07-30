## 1. 服务层方法

- [x] 1.1 在 `tencentcloud/services/cynosdb/service_tencentcloud_cynosdb.go` 中新增 `ModifyInstanceName(ctx, instanceId, instanceName)` 方法，使用 `cynosdb.NewModifyInstanceNameRequest()`，设置 `InstanceId` 与 `InstanceName`，内部使用 `resource.Retry(tccommon.WriteRetryTimeout, ...)` + `tccommon.RetryError(err)`，参考 `ModifyClusterName` 方法实现

## 2. 资源 Schema 与 Update 逻辑

- [x] 2.1 在 `tencentcloud/services/cynosdb/resource_tc_cynosdb_readonly_instance.go` 的 `ResourceTencentCloudCynosdbReadonlyInstance()` 中，移除 `instance_name` 字段的 `ForceNew: true`
- [x] 2.2 在 `resourceTencentCloudCynosdbReadonlyInstanceUpdate` 中，新增 `if d.HasChange("instance_name")` 分支，读取新值并调用 `cynosdbService.ModifyInstanceName(ctx, instanceId, instanceName)`，处理 error

## 3. 单元测试

- [x] 3.1 在 `tencentcloud/services/cynosdb/resource_tc_cynosdb_readonly_instance_test.go` 中新增 `TestUnitCynosdbReadonlyInstance_UpdateInstanceName`，使用 gomonkey mock `UseCynosdbClient`、`ModifyInstanceName`、`DescribeInstances`，验证 Update 路径调用 ModifyInstanceName 且不报错
- [x] 3.2 执行 `go test ./tencentcloud/services/cynosdb/ -run "TestUnitCynosdbReadonlyInstance_UpdateInstanceName" -v -count=1 -gcflags="all=-l"` 确保测试通过

## 4. 文档与 changelog

- [x] 4.1 更新 `website/docs/r/cynosdb_readonly_instance.html.markdown`，移除 `instance_name` 的 `ForceNew` 标记（由 `make doc` 在收尾阶段生成，本任务仅准备 .md 源文件）
- [x] 4.2 更新 `.changelog/4319.txt`，将变更描述改为 `resource/tencentcloud_cynosdb_readonly_instance: support instance_name modification via ModifyInstanceName API`（enhancement）

## 5. 验证

- [x] 5.1 确认所有修改的 Go 文件可正确编译（由后续流程的 go build 验证）

## 6. 额外修复

- [x] 6.1 修复 `resource_tc_cynosdb_cluster_v2_test.go` 中 sweeper 名称冲突（master 预先存在的问题：与 `resource_tc_cynosdb_cluster_test.go` 注册了同名 `tencentcloud_cynosdb` sweeper，导致 `log.Fatalf` 阻断整个包的测试运行），将 v2 的 sweeper 重命名为 `tencentcloud_cynosdb_v2`
